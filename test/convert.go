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

	"github.com/haproxytech/client-native/v6/configuration"
	"github.com/haproxytech/client-native/v6/models"
)

//go:embed expected/structured.json
var expectedStructuredV2JSON []byte
var expectedStructuredV2 map[string]interface{}

//go:embed expected/structured.json
var expectedStructuredJSON []byte
var expectedStructured map[string]interface{}

var onceStrutructured sync.Once

func initStructuredExpected() {
	onceStrutructured.Do(func() {
		err := json.Unmarshal(expectedStructuredV2JSON, &expectedStructuredV2)
		if err != nil {
			panic(err)
		}
		err = json.Unmarshal(expectedStructuredJSON, &expectedStructured)
		if err != nil {
			panic(err)
		}
	})
}

func expectedResources[T any](elementKey string) (map[string]T, error) {
	v := expectedStructured[elementKey]
	j, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	var elems map[string]T
	err = json.Unmarshal(j, &elems)
	if err != nil {
		// Case Defaults, Globals, Traces
		var elem T
		err = json.Unmarshal(j, &elem)
		if err != nil {
			return nil, err
		}
		return map[string]T{"": elem}, err
	}

	return elems, nil
}

func expectedChildResources[P, T any](res map[string][]T, parentKey, parentNameKey, elementKey string) error {
	v := expectedStructured[parentKey]
	parentMap, ok := v.(map[string]interface{})
	if !ok {
		return fmt.Errorf("expectedStructuredV3[%s] is not a map[string]interface{}", parentKey)
	}

	for pname, pmvalue := range parentMap {
		var pkey string
		switch parentKey {
		case "frontends":
			pkey = configuration.FrontendParentName
		case "backends":
			pkey = configuration.BackendParentName
		case "fcgi_apps":
			pkey = configuration.FCGIAppParentName
		case "log_profiles":
			pkey = configuration.LogProfileParentName
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
		case "defaults":
			pkey = configuration.DefaultsParentName
		case "crt_stores":
			pkey = configuration.CrtStoreParentName
		}
		key := fmt.Sprintf("%s/%s", pkey, pname)
		if pname == "dynamic_update_rule_list" {
			key = "dynamic_update"
		}

		switch pmap := pmvalue.(type) {
		case map[string]interface{}:
			e, ok := pmap[elementKey]
			if !ok {
				res[key] = []T{}
				continue
			}
			// Case list of T
			var resources []T
			elistj, err := json.Marshal(e)
			if err != nil {
				return err
			}
			err = json.Unmarshal(elistj, &resources)
			if err == nil {
				res[key] = append(res[key], resources...)
			} else {
				// Case map[string]T
				var resources map[string]T
				err = json.Unmarshal(elistj, &resources)
				if err != nil {
					return err
				}
				for _, v := range resources {
					res[key] = append(res[key], v)
				}
			}
		case []interface{}:
			// Case dynnamic update rule
			var resources []T
			elistj, err := json.Marshal(pmap)
			if err != nil {
				continue
			}
			err = json.Unmarshal(elistj, &resources)
			if err == nil {
				res[key] = append(res[key], resources...)
			}
		default:
			// fmt.Printf("pmap type not handled: %+v\n", pmap)
			continue
		}

	}

	return nil
}

func expectedRootChildResources[P, T any](res map[string][]T, parentKey, parentNameKey, elementKey string) error {
	v := expectedStructured[parentKey]
	parentMap, ok := v.(map[string]interface{})
	if !ok {
		return fmt.Errorf("expectedStructuredV3[%s] is not a map[string]interface{}", parentKey)
	}
	key := parentKey

	e, ok := parentMap[elementKey]
	if !ok {
		res[key] = []T{}
	}
	var resources []T
	elistj, err := json.Marshal(e)
	if err != nil {
		return err
	}

	err = json.Unmarshal(elistj, &resources)
	if err == nil {
		res[key] = append(res[key], resources...)
	}

	return nil
}

func toResMap[T any](keyRoot string, resources map[string]T) map[string][]*T {
	res := make(map[string][]*T)
	for _, v := range resources {
		currentv := v
		res[keyRoot] = append(res[keyRoot], &currentv)
	}
	return res
}

func toSliceOfPtrs[T any](resources []T) []*T {
	res := make([]*T, 0)
	for _, v := range resources {
		currentv := v
		res = append(res, &currentv)
	}
	return res
}

func StructuredToBackendMap() map[string]models.Backends {
	resources, _ := expectedResources[models.Backend]("backends")
	res := make(map[string]models.Backends)
	keyRoot := ""
	t := toResMap(keyRoot, resources)
	res[keyRoot] = t[keyRoot]
	return res
}

func StructuredToFrontendMap() map[string]models.Frontends {
	resources, _ := expectedResources[models.Frontend]("frontends")
	res := make(map[string]models.Frontends)
	keyRoot := ""
	t := toResMap(keyRoot, resources)
	res[keyRoot] = t[keyRoot]
	return res
}

func StructuredToCacheMap() map[string]models.Caches {
	resources, _ := expectedResources[models.Cache]("caches")
	res := make(map[string]models.Caches)
	keyRoot := ""
	t := toResMap(keyRoot, resources)
	res[keyRoot] = t[keyRoot]
	return res
}

func StructuredToACLMap() map[string]models.Acls {
	res := make(map[string]models.Acls)
	resources := make(map[string][]models.ACL)
	_ = expectedChildResources[models.Defaults, models.ACL](resources, "defaults", "name", "acl_list")
	_ = expectedChildResources[models.Frontend, models.ACL](resources, "frontends", "name", "acl_list")
	_ = expectedChildResources[models.Backend, models.ACL](resources, "backends", "name", "acl_list")
	_ = expectedChildResources[models.FCGIApp, models.ACL](resources, "fcgi_apps", "name", "acl_list")
	for k, v := range resources {
		res[k] = toSliceOfPtrs(v)
	}
	return res
}

func StructuredToBackendSwitchingRuleMap() map[string]models.BackendSwitchingRules {
	res := make(map[string]models.BackendSwitchingRules)
	resources := make(map[string][]models.BackendSwitchingRule)
	_ = expectedChildResources[models.Frontend, models.BackendSwitchingRule](resources, "frontends", "name", "backend_switching_rule_list")
	for k, v := range resources {
		res[k] = toSliceOfPtrs(v)
	}
	return res
}

func StructuredToBindMap() map[string]models.Binds {
	res := make(map[string]models.Binds)
	resources := make(map[string][]models.Bind)
	_ = expectedChildResources[models.Frontend, models.Bind](resources, "frontends", "name", "binds")
	_ = expectedChildResources[models.LogForward, models.Bind](resources, "log_forwards", "name", "binds")
	_ = expectedChildResources[models.PeerSection, models.Bind](resources, "peers", "name", "bindss")
	for k, v := range resources {
		res[k] = toSliceOfPtrs(v)
	}
	return res
}

func StructuredToCaptureMap() map[string]models.Captures {
	res := make(map[string]models.Captures)
	resources := make(map[string][]models.Capture)
	_ = expectedChildResources[models.Frontend, models.Capture](resources, "frontends", "name", "captures")
	for k, v := range resources {
		res[k] = toSliceOfPtrs(v)
	}
	return res
}

func StructuredToGlobalMap() models.Global {
	resources, _ := expectedResources[models.Global]("global")
	return resources[""]
}

func StructuredToTracesMap() models.Traces {
	resources, _ := expectedResources[models.Traces]("traces")
	return resources[""]
}

func StructuredToDefaultsMap() map[string][]*models.Defaults {
	resources, _ := expectedResources[models.Defaults]("defaults")
	res := make(map[string][]*models.Defaults)
	for _, v := range resources {
		currentv := v
		res[v.Name] = append(res[v.Name], &currentv)
	}
	return res
}

func StructuredToDgramBindMap() map[string]models.DgramBinds {
	res := make(map[string]models.DgramBinds)
	resources := make(map[string][]models.DgramBind)
	_ = expectedChildResources[models.LogForward, models.DgramBind](resources, "log_forwards", "name", "dgram_binds")
	for k, v := range resources {
		res[k] = toSliceOfPtrs(v)
	}
	return res
}

func StructuredToFCGIAppMap() map[string]models.FCGIApps {
	resources, _ := expectedResources[models.FCGIApp]("fcgi_apps")
	res := make(map[string]models.FCGIApps)
	keyRoot := ""
	t := toResMap(keyRoot, resources)
	res[keyRoot] = t[keyRoot]
	return res
}

func StructuredToFilterMap() map[string]models.Filters {
	res := make(map[string]models.Filters)
	resources := make(map[string][]models.Filter)
	_ = expectedChildResources[models.Frontend, models.Filter](resources, "frontends", "name", "filter_list")
	_ = expectedChildResources[models.Backend, models.Filter](resources, "backends", "name", "filter_list")
	for k, v := range resources {
		res[k] = toSliceOfPtrs(v)
	}
	return res
}

func StructuredToGroupMap() map[string]models.Groups {
	res := make(map[string]models.Groups)
	resources := make(map[string][]models.Group)
	_ = expectedChildResources[models.Userlist, models.Group](resources, "userlists", "name", "groups")
	for k, v := range resources {
		res[k] = toSliceOfPtrs(v)
	}
	return res
}

func StructuredToHTTPAfterResponseRuleMap() map[string]models.HTTPAfterResponseRules {
	res := make(map[string]models.HTTPAfterResponseRules)
	resources := make(map[string][]models.HTTPAfterResponseRule)
	_ = expectedChildResources[models.Frontend, models.HTTPAfterResponseRule](resources, "frontends", "name", "http_after_response_rule_list")
	_ = expectedChildResources[models.Backend, models.HTTPAfterResponseRule](resources, "backends", "name", "http_after_response_rule_list")
	_ = expectedChildResources[models.Defaults, models.HTTPAfterResponseRule](resources, "defaults", "name", "http_after_response_rule_list")
	for k, v := range resources {
		res[k] = toSliceOfPtrs(v)
	}
	return res
}

func StructuredToHTTPCheckMap() map[string]models.HTTPChecks {
	res := make(map[string]models.HTTPChecks)
	resources := make(map[string][]models.HTTPCheck)
	_ = expectedChildResources[models.Backend, models.HTTPCheck](resources, "backends", "name", "http_check_list")
	_ = expectedChildResources[models.Defaults, models.HTTPCheck](resources, "defaults", "name", "http_check_list")
	for k, v := range resources {
		res[k] = toSliceOfPtrs(v)
	}
	return res
}

func StructuredToHTTPErrorRuleMap() map[string]models.HTTPErrorRules {
	res := make(map[string]models.HTTPErrorRules)
	resources := make(map[string][]models.HTTPErrorRule)
	_ = expectedChildResources[models.Backend, models.HTTPErrorRule](resources, "backends", "name", "http_error_rule_list")
	_ = expectedChildResources[models.Frontend, models.HTTPErrorRule](resources, "frontends", "name", "http_error_rule_list")
	_ = expectedChildResources[models.Defaults, models.HTTPErrorRule](resources, "defaults", "name", "http_error_rule_list")
	for k, v := range resources {
		res[k] = toSliceOfPtrs(v)
	}
	return res
}

func StructuredToHTTPErrorSectionMap() map[string]models.HTTPErrorsSections {
	resources, _ := expectedResources[models.HTTPErrorsSection]("http_errors")
	res := make(map[string]models.HTTPErrorsSections)
	keyRoot := ""
	t := toResMap(keyRoot, resources)
	res[keyRoot] = t[keyRoot]
	return res
}

func StructuredToHTTPRequestRuleMap() map[string]models.HTTPRequestRules {
	res := make(map[string]models.HTTPRequestRules)
	resources := make(map[string][]models.HTTPRequestRule)
	_ = expectedChildResources[models.Backend, models.HTTPRequestRule](resources, "backends", "name", "http_request_rule_list")
	_ = expectedChildResources[models.Frontend, models.HTTPRequestRule](resources, "frontends", "name", "http_request_rule_list")
	_ = expectedChildResources[models.Defaults, models.HTTPRequestRule](resources, "defaults", "name", "http_request_rule_list")
	for k, v := range resources {
		res[k] = toSliceOfPtrs(v)
	}
	return res
}

func StructuredToHTTPResponseRuleMap() map[string]models.HTTPResponseRules {
	res := make(map[string]models.HTTPResponseRules)
	resources := make(map[string][]models.HTTPResponseRule)
	_ = expectedChildResources[models.Backend, models.HTTPResponseRule](resources, "backends", "name", "http_response_rule_list")
	_ = expectedChildResources[models.Frontend, models.HTTPResponseRule](resources, "frontends", "name", "http_response_rule_list")
	_ = expectedChildResources[models.Defaults, models.HTTPResponseRule](resources, "defaults", "name", "http_response_rule_list")
	for k, v := range resources {
		res[k] = toSliceOfPtrs(v)
	}
	return res
}

func StructuredToQUICInitialRuleMap() map[string]models.QUICInitialRules {
	res := make(map[string]models.QUICInitialRules)
	resources := make(map[string][]models.QUICInitialRule)
	_ = expectedChildResources[models.Backend, models.QUICInitialRule](resources, "frontends", "name", "quic_initial_rule_list")
	_ = expectedChildResources[models.Defaults, models.QUICInitialRule](resources, "defaults", "name", "quic_initial_rule_list")
	for k, v := range resources {
		res[k] = toSliceOfPtrs(v)
	}
	return res
}

func StructuredToLogForwardMap() map[string]models.LogForwards {
	resources, _ := expectedResources[models.LogForward]("log_forwards")
	res := make(map[string]models.LogForwards)
	keyRoot := ""
	t := toResMap(keyRoot, resources)
	res[keyRoot] = t[keyRoot]
	return res
}

func StructuredToLogTargetMap() map[string]models.LogTargets {
	res := make(map[string]models.LogTargets)
	resources := make(map[string][]models.LogTarget)
	_ = expectedChildResources[models.Backend, models.LogTarget](resources, "backends", "name", "log_target_list")
	_ = expectedChildResources[models.Frontend, models.LogTarget](resources, "frontends", "name", "log_target_list")
	_ = expectedChildResources[models.LogForward, models.LogTarget](resources, "log_forwards", "name", "log_target_list")
	_ = expectedChildResources[models.PeerSection, models.LogTarget](resources, "peers", "name", "log_target_list")
	_ = expectedRootChildResources[models.Defaults](resources, "defaults", "name", "log_target_list")
	_ = expectedRootChildResources[models.Global](resources, "global", "name", "log_target_list")
	for k, v := range resources {
		res[k] = toSliceOfPtrs(v)
	}
	return res
}

func StructuredToMailerEntryMap() map[string]models.MailerEntries {
	res := make(map[string]models.MailerEntries)
	resources := make(map[string][]models.MailerEntry)
	_ = expectedChildResources[models.MailersSection, models.MailerEntry](resources, "mailers_sections", "name", "mailer_entries")
	for k, v := range resources {
		res[k] = toSliceOfPtrs(v)
	}
	return res
}

func StructuredToMailersSectionMap() map[string]models.MailersSections {
	resources, _ := expectedResources[models.MailersSection]("mailers_sections")
	res := make(map[string]models.MailersSections)
	keyRoot := ""
	t := toResMap(keyRoot, resources)
	res[keyRoot] = t[keyRoot]
	return res
}

func StructuredToCrtStoreMap() map[string]models.CrtStores {
	resources, _ := expectedResources[models.CrtStore]("crt_stores")
	res := make(map[string]models.CrtStores)
	keyRoot := ""
	t := toResMap(keyRoot, resources)
	res[keyRoot] = t[keyRoot]
	return res
}

func StructuredToNameserverMap() map[string]models.Nameservers {
	res := make(map[string]models.Nameservers)
	resources := make(map[string][]models.Nameserver)
	_ = expectedChildResources[models.Resolver, models.Nameserver](resources, "resolvers", "name", "nameservers")
	for k, v := range resources {
		res[k] = toSliceOfPtrs(v)
	}
	return res
}

func StructuredToPeerEntryMap() map[string]models.PeerEntries {
	res := make(map[string]models.PeerEntries)
	resources := make(map[string][]models.PeerEntry)
	_ = expectedChildResources[models.PeerSection, models.PeerEntry](resources, "peers", "name", "peer_entries")
	for k, v := range resources {
		res[k] = toSliceOfPtrs(v)
	}
	return res
}

func StructuredToPeerSectionMap() map[string]models.PeerSections {
	resources, _ := expectedResources[models.PeerSection]("peers")
	res := make(map[string]models.PeerSections)
	keyRoot := ""
	t := toResMap(keyRoot, resources)
	res[keyRoot] = t[keyRoot]
	return res
}

func StructuredToProgramMap() map[string]models.Programs {
	resources, _ := expectedResources[models.Program]("programs")
	res := make(map[string]models.Programs)
	keyRoot := ""
	t := toResMap(keyRoot, resources)
	res[keyRoot] = t[keyRoot]
	return res
}

func StructuredToResolverMap() map[string]models.Resolvers {
	resources, _ := expectedResources[models.Resolver]("resolvers")
	res := make(map[string]models.Resolvers)
	keyRoot := ""
	t := toResMap(keyRoot, resources)
	res[keyRoot] = t[keyRoot]
	return res
}

func StructuredToRingMap() map[string]models.Rings {
	resources, _ := expectedResources[models.Ring]("rings")
	res := make(map[string]models.Rings)
	keyRoot := ""
	t := toResMap(keyRoot, resources)
	res[keyRoot] = t[keyRoot]
	return res
}

func StructuredToServerSwitchingRuleMap() map[string]models.ServerSwitchingRules {
	res := make(map[string]models.ServerSwitchingRules)
	resources := make(map[string][]models.ServerSwitchingRule)
	_ = expectedChildResources[models.Backend](resources, "backends", "name", "server_switching_rule_list")
	for k, v := range resources {
		res[k] = toSliceOfPtrs(v)
	}
	return res
}

func StructuredToServerTemplateMap() map[string]models.ServerTemplates {
	res := make(map[string]models.ServerTemplates)
	resources := make(map[string][]models.ServerTemplate)
	_ = expectedChildResources[models.Backend, models.ServerTemplate](resources, "backends", "name", "server_templates")
	for k, v := range resources {
		res[k] = toSliceOfPtrs(v)
	}
	return res
}

func StructuredToServerMap() map[string]models.Servers { //nolint:dupl
	res := make(map[string]models.Servers)
	resources := make(map[string][]models.Server)
	_ = expectedChildResources[models.Backend, models.Server](resources, "backends", "name", "servers")
	_ = expectedChildResources[models.Ring, models.Server](resources, "rings", "name", "servers")
	_ = expectedChildResources[models.PeerSection, models.Server](resources, "peers", "name", "servers")
	for k, v := range resources {
		res[k] = toSliceOfPtrs(v)
	}
	return res
}

func StructuredToStickRuleMap() map[string]models.StickRules {
	res := make(map[string]models.StickRules)
	resources := make(map[string][]models.StickRule)
	_ = expectedChildResources[models.Backend, models.StickRule](resources, "backends", "name", "stick_rule_list")
	for k, v := range resources {
		res[k] = toSliceOfPtrs(v)
	}
	return res
}

func StructuredToTCPRequestRuleMap() map[string]models.TCPRequestRules {
	res := make(map[string]models.TCPRequestRules)
	resources := make(map[string][]models.TCPRequestRule)
	_ = expectedChildResources[models.Backend, models.TCPRequestRule](resources, "backends", "name", "tcp_request_rule_list")
	_ = expectedChildResources[models.Frontend](resources, "frontends", "name", "tcp_request_rule_list")
	_ = expectedChildResources[models.Defaults](resources, "defaults", "name", "tcp_request_rule_list")

	for k, v := range resources {
		res[k] = toSliceOfPtrs(v)
	}
	return res
}

func StructuredToTCPResponseRuleMap() map[string]models.TCPResponseRules {
	res := make(map[string]models.TCPResponseRules)
	resources := make(map[string][]models.TCPResponseRule)
	_ = expectedChildResources[models.Backend](resources, "backends", "name", "tcp_response_rule_list")
	for k, v := range resources {
		res[k] = toSliceOfPtrs(v)
	}
	return res
}

func StructuredToTCPCheckMap() map[string]models.TCPChecks {
	res := make(map[string]models.TCPChecks)
	resources := make(map[string][]models.TCPCheck)
	_ = expectedChildResources[models.Backend, models.TCPCheck](resources, "backends", "name", "tcp_check_rule_list")
	_ = expectedRootChildResources[models.Defaults, models.TCPCheck](resources, "defaults", "name", "tcp_check_rule_list")
	for k, v := range resources {
		res[k] = toSliceOfPtrs(v)
	}
	return res
}

func StructuredToUserMap() map[string]models.Users {
	res := make(map[string]models.Users)
	resources := make(map[string][]models.User)
	_ = expectedChildResources[models.Defaults, models.User](resources, "userlists", "name", "users")
	for k, v := range resources {
		res[k] = toSliceOfPtrs(v)
	}
	return res
}

func StructuredToLogProfileMap() map[string]models.LogProfiles {
	resources, _ := expectedResources[models.LogProfile]("log_profiles")
	res := make(map[string]models.LogProfiles)
	keyRoot := ""
	t := toResMap(keyRoot, resources)
	res[keyRoot] = t[keyRoot]
	return res
}
