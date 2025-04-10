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

package parents

type Parent struct {
	PathParentType    string
	ParentType        string
	IsGenericParent   bool // only one of the parents for a child type can be set as Geenric - used in Dapi
	GenericParentType string
	IsUnnamedParent   bool
}

func Parents(childType string) []Parent {
	switch childType {
	case ServerChildType:
		return []Parent{
			{PathParentType: "backends", ParentType: "Backend", IsGenericParent: true},
			{PathParentType: "peers", ParentType: "Peer", GenericParentType: "Backend"},
			{PathParentType: "rings", ParentType: "Ring", GenericParentType: "Backend"},
		}
	case HTTPAfterResponseRuleChildType:
		return []Parent{
			{PathParentType: "backends", ParentType: "Backend", IsGenericParent: true},
			{PathParentType: "frontends", ParentType: "Frontend", GenericParentType: "Backend"},
		}
	case HTTPCheckChildType:
		return []Parent{
			{PathParentType: "backends", ParentType: "Backend", IsGenericParent: true},
			{PathParentType: "defaults", ParentType: "Defaults", GenericParentType: "Backend"},
		}
	case HTTPErrorRuleChildType:
		return []Parent{
			{PathParentType: "backends", ParentType: "Backend", IsGenericParent: true},
			{PathParentType: "frontends", ParentType: "Frontend", GenericParentType: "Backend"},
			{PathParentType: "defaults", ParentType: "Defaults", GenericParentType: "Backend"},
		}
	case HTTPRequestRuleChildType:
		return []Parent{
			{PathParentType: "backends", ParentType: "Backend", IsGenericParent: true},
			{PathParentType: "frontends", ParentType: "Frontend", GenericParentType: "Backend"},
		}
	case HTTPResponseRuleChildType:
		return []Parent{
			{PathParentType: "backends", ParentType: "Backend", IsGenericParent: true},
			{PathParentType: "frontends", ParentType: "Frontend", GenericParentType: "Backend"},
		}
	case TCPCheckChildType:
		return []Parent{
			{PathParentType: "backends", ParentType: "Backend", IsGenericParent: true},
			{PathParentType: "defaults", ParentType: "Defaults", GenericParentType: "Backend"},
		}
	case TCPRequestRuleChildType:
		return []Parent{
			{PathParentType: "backends", ParentType: "Backend", IsGenericParent: true},
			{PathParentType: "frontends", ParentType: "Frontend", GenericParentType: "Backend"},
		}
	case TCPResponseRuleChildType:
		return []Parent{
			{PathParentType: "backends", ParentType: "Backend", IsGenericParent: true},
		}
	case QUICInitialRuleType:
		return []Parent{
			{PathParentType: "frontends", ParentType: "Frontend", IsGenericParent: true},
			{PathParentType: "defaults", ParentType: "Defaults", GenericParentType: "Frontend"},
		}
	case ACLChildType:
		return []Parent{
			{PathParentType: "backends", ParentType: "Backend", IsGenericParent: true},
			{PathParentType: "frontends", ParentType: "Frontend", GenericParentType: "Backend"},
			{PathParentType: "fcgi_apps", ParentType: "FCGIApp", GenericParentType: "Backend"},
		}
	case BindChildType:
		return []Parent{
			{PathParentType: "frontends", ParentType: "Frontend", IsGenericParent: true},
			{PathParentType: "log_forwards", ParentType: "LogForward", GenericParentType: "Frontend"},
			{PathParentType: "peers", ParentType: "Peer", GenericParentType: "Frontend"},
		}
	case FilterChildType:
		return []Parent{
			{PathParentType: "backends", ParentType: "Backend", IsGenericParent: true},
			{PathParentType: "frontends", ParentType: "Frontend", GenericParentType: "Backend"},
		}
	case LogTargetChildType:
		return []Parent{
			{PathParentType: "backends", ParentType: "Backend", IsGenericParent: true},
			{PathParentType: "frontends", ParentType: "Frontend", GenericParentType: "Backend"},
			{PathParentType: "defaults", ParentType: "Defaults", GenericParentType: "Backend"},
			{PathParentType: "peers", ParentType: "Peer", GenericParentType: "Backend"},
			{PathParentType: "log_forwards", ParentType: "LogForward", GenericParentType: "Backend"},
		}
	case SSLFrontUseChildType:
		return []Parent{
			{PathParentType: "frontends", ParentType: "Frontend", IsGenericParent: true},
		}
	}
	return nil
}
