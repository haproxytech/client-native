/*
Copyright 2026 HAProxy Technologies

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package parsers

import (
	"errors"
	"fmt"
	"strings"

	"github.com/haproxytech/client-native/v6/config-parser/common"
	parsererrors "github.com/haproxytech/client-native/v6/config-parser/errors"
	"github.com/haproxytech/client-native/v6/misc"
	"github.com/haproxytech/client-native/v6/models"
)

type FilterSequence struct {
	data        []models.FilterSequence
	preComments []string // comments that appear before the actual line
}

func (p *FilterSequence) Parse(line string, parts []string, comment string) (string, error) {
	sequence, err := parseFilterSequence(parts, comment)
	if err != nil {
		return "", &parsererrors.ParseError{Parser: "FilterSequence", Line: line}
	}
	if err := p.Insert(*sequence, -1); err != nil {
		return "", err
	}
	return "", nil
}

func parseFilterSequence(parts []string, comment string) (*models.FilterSequence, error) {
	if len(parts) < 3 || parts[0] != "filter-sequence" {
		return nil, parsererrors.ErrInvalidData
	}

	sequence := &models.FilterSequence{Side: parts[1]}
	switch sequence.Side {
	case "global", "request", "response":
	default:
		return nil, parsererrors.ErrInvalidData
	}

	for _, token := range parts[2:] {
		rule, err := parseFilterSequenceRule(token)
		if err != nil {
			return nil, err
		}
		sequence.Rules = append(sequence.Rules, rule)
	}

	if comment != "" {
		sequence.Metadata = misc.ParseMetadata(comment)
	}
	return sequence, nil
}

func parseFilterSequenceRule(value string) (*models.FilterSequenceRule, error) {
	instanceText, rest, found := strings.Cut(value, ":")
	if !found {
		return nil, parsererrors.ErrInvalidData
	}
	direction, targetText, found := strings.Cut(rest, "(")
	if !found || !strings.HasSuffix(targetText, ")") {
		return nil, parsererrors.ErrInvalidData
	}
	switch direction {
	case "before", "after":
	default:
		return nil, parsererrors.ErrInvalidData
	}
	targetText = strings.TrimSuffix(targetText, ")")
	if strings.ContainsAny(targetText, "()") {
		return nil, parsererrors.ErrInvalidData
	}

	instance, err := parseFilterSelector(instanceText)
	if err != nil {
		return nil, err
	}
	target, err := parseFilterSelector(targetText)
	if err != nil {
		return nil, err
	}
	return &models.FilterSequenceRule{
		Instance:  instance,
		Direction: direction,
		Target:    target,
	}, nil
}

// ValidateFilterSequence rejects an explicit cycle within a single sequence
// line. Cycles spanning several filter-sequence lines are left to HAProxy:
// repeated directives are cumulative and the binary already validates them.
func ValidateFilterSequence(sequence models.FilterSequence) error {
	graph := map[string][]string{}
	for _, rule := range sequence.Rules {
		if rule == nil || rule.Instance == nil || rule.Target == nil {
			return errors.New("filter sequence rule is empty")
		}
		from, to := filterSelectorString(rule.Instance), filterSelectorString(rule.Target)
		if rule.Direction == "after" {
			from, to = to, from
		}
		graph[from] = append(graph[from], to)
	}

	const (
		visiting = 1
		visited  = 2
	)
	colors := map[string]int{}
	var visit func(node string) error
	visit = func(node string) error {
		switch colors[node] {
		case visiting:
			return fmt.Errorf("filter sequence cycle detected on side %q", sequence.Side)
		case visited:
			return nil
		}
		colors[node] = visiting
		for _, next := range graph[node] {
			if err := visit(next); err != nil {
				return err
			}
		}
		colors[node] = visited
		return nil
	}

	for node := range graph {
		if err := visit(node); err != nil {
			return err
		}
	}
	return nil
}

func (p *FilterSequence) Result() ([]common.ReturnResultLine, error) {
	if len(p.data) == 0 {
		return nil, parsererrors.ErrFetch
	}
	result := make([]common.ReturnResultLine, len(p.data))
	for index, sequence := range p.data {
		var line strings.Builder
		line.WriteString("filter-sequence ")
		line.WriteString(sequence.Side)
		for _, rule := range sequence.Rules {
			line.WriteString(" ")
			line.WriteString(filterSelectorString(rule.Instance))
			line.WriteString(":")
			line.WriteString(rule.Direction)
			line.WriteString("(")
			line.WriteString(filterSelectorString(rule.Target))
			line.WriteString(")")
		}

		comment := ""
		if sequence.Metadata != nil {
			serializedComment, err := misc.SerializeMetadata(sequence.Metadata)
			if err == nil {
				comment = serializedComment
			}
		}
		result[index] = common.ReturnResultLine{
			Data:    line.String(),
			Comment: comment,
		}
	}
	return result, nil
}
