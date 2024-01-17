// Copyright 2019 HAProxy Technologies
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package test

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/haproxytech/client-native/v5/configuration"
	"github.com/haproxytech/client-native/v5/models"
)

//go:embed expected/structured.json
var expectedStructuredJSON []byte
var expectedStructured map[string]interface{}

var onceStrutructured sync.Once

func initStructuredExpected() {
	onceStrutructured.Do(func() {
		err := json.Unmarshal(expectedStructuredJSON, &expectedStructured)
		if err != nil {
			panic(err)
		}
	})
}

func expectedResources[T any](res T, elementKey string) error {
	v := expectedStructured[elementKey]
	j, err := json.Marshal(v)
	if err != nil {
		return err
	}
	err = json.Unmarshal(j, &res)
	if err != nil {
		return err
	}
	return nil
}

func expectedChildResources[P, T any](res map[string]T, parentKey, parentNameKey, elementKey string) error {
	v := expectedStructured[parentKey]
	vlist, ok := v.([]interface{})
	if ok {
		for _, p := range vlist {
			pj, err := json.Marshal(p)
			if err != nil {
				return err
			}
			var parent P
			err = json.Unmarshal(pj, &parent)
			if err != nil {
				return err
			}

			pmap := p.(map[string]interface{})
			pName, ok := pmap[parentNameKey]
			if !ok {
				continue
			}
			var pkey string
			switch parentKey {
			case "frontends":
				pkey = configuration.FrontendParentName
			case "backends":
				pkey = configuration.BackendParentName
			case "fcgi_apps":
				pkey = configuration.FCGIAppParentName
			case "log_forwards":
				pkey = configuration.LogForwardParentName
			case "userlists":
				pkey = "userlist"
			case "mailers_sections":
				pkey = "mailers_sections"
			case "resolvers":
				pkey = configuration.ResolverParentName
			case "peers":
				pkey = configuration.PeersParentName
			case "rings":
				pkey = configuration.RingParentName
			}
			key := fmt.Sprintf("%s/%s", pkey, pName)

			ellist, ok := pmap[elementKey]
			if !ok {
				res[key] = *new(T)
				continue
			}
			elistj, err := json.Marshal(ellist)
			if err != nil {
				return err
			}
			var resources T
			err = json.Unmarshal(elistj, &resources)
			if err != nil {
				return err
			}
			res[key] = resources
		}
	} else {
		pmap := v.(map[string]interface{})
		key := parentKey

		ellist, ok := pmap[elementKey]
		if !ok {
			res[key] = *new(T)
		}
		elistj, err := json.Marshal(ellist)
		if err != nil {
			return err
		}
		var resources T
		err = json.Unmarshal(elistj, &resources)
		if err != nil {
			return err
		}
		res[key] = resources
	}

	return nil
}

func StructuredToBackendMap() map[string]models.Backends {
	var l models.Backends
	_ = expectedResources(&l, "backends")
	res := make(map[string]models.Backends)
	keyRoot := ""
	res[keyRoot] = l
	return res
}

func StructuredToFrontendMap() map[string]models.Frontends {
	var l models.Frontends
	_ = expectedResources(&l, "frontends")
	res := make(map[string]models.Frontends)
	keyRoot := ""
	res[keyRoot] = l
	return res
}

func StructuredToCacheMap() map[string]models.Caches {
	var l models.Caches
	_ = expectedResources(&l, "caches")
	res := make(map[string]models.Caches)
	keyRoot := ""
	res[keyRoot] = l
	return res
}

func StructuredToACLMap() map[string]models.Acls {
	res := make(map[string]models.Acls)
	_ = expectedChildResources[models.Frontend](res, "frontends", "name", "acls")
	_ = expectedChildResources[models.Backend](res, "backends", "name", "acls")
	_ = expectedChildResources[models.FCGIApp](res, "fcgi_apps", "name", "acls")
	return res
}

func StructuredToBackendSwitchingRuleMap() map[string]models.BackendSwitchingRules {
	res := make(map[string]models.BackendSwitchingRules)
	_ = expectedChildResources[models.Frontend](res, "frontends", "name", "backend_switching_rules")
	return res
}

func StructuredToBindMap() map[string]models.Binds {
	res := make(map[string]models.Binds)
	_ = expectedChildResources[models.Frontend](res, "frontends", "name", "binds")
	_ = expectedChildResources[models.LogForward](res, "log_forwards", "name", "binds")
	_ = expectedChildResources[models.PeerSection](res, "peers", "name", "binds")
	return res
}

func StructuredToCaptureMap() map[string]models.Captures {
	res := make(map[string]models.Captures)
	_ = expectedChildResources[models.Frontend](res, "frontends", "name", "captures")
	return res
}

func StructuredToDefaultsMap() models.Defaults {
	var l models.Defaults
	_ = expectedResources(&l, "defaults")
	return l
}

func StructuredToGlobalMap() models.Global {
	var l models.Global
	_ = expectedResources(&l, "global")
	return l
}

func StructuredToNamedDefaultsMap() map[string][]*models.Defaults {
	var l []*models.Defaults
	res := make(map[string][]*models.Defaults)
	_ = expectedResources(&l, "named_defaults")
	for _, v := range l {
		res[v.Name] = append(res[v.Name], v)
	}
	return res
}

func StructuredToDgramBindMap() map[string]models.DgramBinds {
	res := make(map[string]models.DgramBinds)
	_ = expectedChildResources[models.LogForward](res, "log_forwards", "name", "dgram_binds")
	return res
}

func StructuredToFCGIAppMap() map[string]models.FCGIApps {
	res := make(map[string]models.FCGIApps)
	var l models.FCGIApps
	_ = expectedResources(&l, "fcgi_apps")
	keyRoot := ""
	res[keyRoot] = l
	return res
}

func StructuredToFilterMap() map[string]models.Filters {
	res := make(map[string]models.Filters)
	_ = expectedChildResources[models.Frontend](res, "frontends", "name", "filters")
	_ = expectedChildResources[models.Backend](res, "backends", "name", "filters")
	return res
}

func StructuredToGroupMap() map[string]models.Groups {
	res := make(map[string]models.Groups)
	_ = expectedChildResources[models.Userlist](res, "userlists", "name", "groups")
	return res
}

func StructuredToHTTPAfterResponseRuleMap() map[string]models.HTTPAfterResponseRules {
	res := make(map[string]models.HTTPAfterResponseRules)
	_ = expectedChildResources[models.Frontend](res, "frontends", "name", "http_after_response_rules")
	_ = expectedChildResources[models.Backend](res, "backends", "name", "http_after_response_rules")
	return res
}

func StructuredToHTTPCheckMap() map[string]models.HTTPChecks {
	res := make(map[string]models.HTTPChecks)
	_ = expectedChildResources[models.Backend](res, "backends", "name", "http_checks")
	_ = expectedChildResources[models.Defaults](res, "defaults", "name", "http_checks")
	return res
}

func StructuredToHTTPErrorRuleMap() map[string]models.HTTPErrorRules {
	res := make(map[string]models.HTTPErrorRules)
	_ = expectedChildResources[models.Backend](res, "backends", "name", "http_error_rules")
	_ = expectedChildResources[models.Frontend](res, "frontends", "name", "http_error_rules")
	_ = expectedChildResources[models.Defaults](res, "defaults", "name", "http_error_rules")
	return res
}

func StructuredToHTTPErrorSectionMap() map[string]models.HTTPErrorsSections {
	res := make(map[string]models.HTTPErrorsSections)
	var l models.HTTPErrorsSections
	_ = expectedResources(&l, "http_errors")
	keyRoot := ""
	res[keyRoot] = l
	return res
}

func StructuredToHTTPRequestRuleMap() map[string]models.HTTPRequestRules {
	res := make(map[string]models.HTTPRequestRules)
	_ = expectedChildResources[models.Backend](res, "backends", "name", "http_request_rules")
	_ = expectedChildResources[models.Frontend](res, "frontends", "name", "http_request_rules")
	return res
}

func StructuredToHTTPResponseRuleMap() map[string]models.HTTPResponseRules {
	res := make(map[string]models.HTTPResponseRules)
	_ = expectedChildResources[models.Backend](res, "backends", "name", "http_response_rules")
	_ = expectedChildResources[models.Frontend](res, "frontends", "name", "http_response_rules")
	return res
}

func StructuredToLogForwardMap() map[string]models.LogForwards {
	res := make(map[string]models.LogForwards)
	var l models.LogForwards
	_ = expectedResources(&l, "log_forwards")
	keyRoot := ""
	res[keyRoot] = l
	return res
}

func StructuredToLogTargetMap() map[string]models.LogTargets {
	res := make(map[string]models.LogTargets)
	_ = expectedChildResources[models.Backend](res, "backends", "name", "log_targets")
	_ = expectedChildResources[models.Frontend](res, "frontends", "name", "log_targets")
	_ = expectedChildResources[models.LogForward](res, "log_forwards", "name", "log_targets")
	_ = expectedChildResources[models.PeerSection](res, "peers", "name", "log_targets")
	_ = expectedChildResources[models.Defaults](res, "defaults", "name", "log_targets")
	_ = expectedChildResources[models.Global](res, "global", "name", "log_targets")
	return res
}

func StructuredToMailerEntryMap() map[string]models.MailerEntries {
	res := make(map[string]models.MailerEntries)
	_ = expectedChildResources[models.MailersSection](res, "mailers_sections", "name", "mailer_entries")
	return res
}

func StructuredToMailersSectionMap() map[string]models.MailersSections {
	res := make(map[string]models.MailersSections)
	var l models.MailersSections
	_ = expectedResources(&l, "mailers_sections")
	keyRoot := ""
	res[keyRoot] = l
	return res
}

func StructuredToNameserverMap() map[string]models.Nameservers {
	res := make(map[string]models.Nameservers)
	_ = expectedChildResources[models.Resolver](res, "resolvers", "name", "nameservers")
	return res
}

func StructuredToPeerEntryMap() map[string]models.PeerEntries {
	res := make(map[string]models.PeerEntries)
	_ = expectedChildResources[models.PeerSection](res, "peers", "name", "peer_entries")
	return res
}

func StructuredToPeerSectionMap() map[string]models.PeerSections {
	res := make(map[string]models.PeerSections)
	var l models.PeerSections
	_ = expectedResources(&l, "peers")
	keyRoot := ""
	res[keyRoot] = l
	return res
}

func StructuredToProgramMap() map[string]models.Programs {
	res := make(map[string]models.Programs)
	var l models.Programs
	_ = expectedResources(&l, "programs")
	keyRoot := ""
	res[keyRoot] = l
	return res
}

func StructuredToResolverMap() map[string]models.Resolvers {
	res := make(map[string]models.Resolvers)
	var l models.Resolvers
	_ = expectedResources(&l, "resolvers")
	keyRoot := ""
	res[keyRoot] = l
	return res
}

func StructuredToRingMap() map[string]models.Rings {
	res := make(map[string]models.Rings)
	var l models.Rings
	_ = expectedResources(&l, "rings")
	keyRoot := ""
	res[keyRoot] = l
	return res
}

func StructuredToServerSwitchingRuleMap() map[string]models.ServerSwitchingRules {
	res := make(map[string]models.ServerSwitchingRules)
	_ = expectedChildResources[models.Backend](res, "backends", "name", "server_switching_rules")
	return res
}

func StructuredToServerTemplateMap() map[string]models.ServerTemplates {
	res := make(map[string]models.ServerTemplates)
	_ = expectedChildResources[models.Backend](res, "backends", "name", "server_templates")

	return res
}

func StructuredToServerMap() map[string]models.Servers { //nolint:dupl
	res := make(map[string]models.Servers)
	_ = expectedChildResources[models.Backend](res, "backends", "name", "servers")
	_ = expectedChildResources[models.Ring](res, "rings", "name", "servers")
	_ = expectedChildResources[models.PeerSection](res, "peers", "name", "servers")
	return res
}

func StructuredToStickRuleMap() map[string]models.StickRules {
	res := make(map[string]models.StickRules)
	_ = expectedChildResources[models.Backend](res, "backends", "name", "stick_rules")
	return res
}

func StructuredToTCPRequestRuleMap() map[string]models.TCPRequestRules {
	res := make(map[string]models.TCPRequestRules)
	_ = expectedChildResources[models.Backend](res, "backends", "name", "tcp_request_rules")
	_ = expectedChildResources[models.Frontend](res, "frontends", "name", "tcp_request_rules")
	return res
}

func StructuredToTCPResponseRuleMap() map[string]models.TCPResponseRules {
	res := make(map[string]models.TCPResponseRules)
	_ = expectedChildResources[models.Backend](res, "backends", "name", "tcp_response_rules")
	return res
}

func StructuredToTCPCheckMap() map[string]models.TCPChecks {
	res := make(map[string]models.TCPChecks)
	_ = expectedChildResources[models.Backend](res, "backends", "name", "tcp_checks")
	_ = expectedChildResources[models.Defaults](res, "defaults", "name", "tcp_checks")
	return res
}

func StructuredToUserMap() map[string]models.Users {
	res := make(map[string]models.Users)
	_ = expectedChildResources[models.Defaults](res, "userlists", "name", "users")
	return res
}
