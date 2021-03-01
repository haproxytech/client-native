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

package configuration

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/haproxytech/client-native/v2/misc"
	"github.com/haproxytech/client-native/v2/models"
)

func TestGetSites(t *testing.T) { //nolint:gocognit,gocyclo
	v, sites, err := client.GetSites("")
	if err != nil {
		t.Error(err.Error())
	}

	if len(sites) != 2 {
		t.Errorf("%v sites returned, expected 2", len(sites))
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	for _, s := range sites {
		switch s.Name {
		case "test":
			if *s.Service.Maxconn != 2000 {
				t.Errorf("%v: Maxconn not 2000: %v", s.Name, *s.Service.Maxconn)
			}
			if s.Service.Mode != "http" {
				t.Errorf("%v: Mode not http: %v", s.Name, s.Service.Mode)
			}
			if s.Service.HTTPConnectionMode != "httpclose" {
				t.Errorf("%v: HTTPConnectionMode not httpclose: %v", s.Name, s.Service.HTTPConnectionMode)
			}
			if len(s.Service.Listeners) != 2 {
				t.Errorf("%v: Got %v listeners, expected 2", s.Name, len(s.Service.Listeners))
			}
			for _, l := range s.Service.Listeners {
				if l.Name != "webserv" && l.Name != "webserv2" {
					t.Errorf("Expected only webserv or webserv2 listeners, %v found", l.Name)
				}
				if l.Address != "192.168.1.1" {
					t.Errorf("%v: Address not 192.168.1.1: %v", l.Name, l.Address)
				}
				if *l.Port != 80 && *l.Port != 8080 {
					t.Errorf("%v: Port not 80 or 8080: %v", l.Name, *l.Port)
				}
			}
			for _, b := range s.Farms {
				switch b.Name {
				case "test":
					if b.UseAs != "default" {
						t.Errorf("%v: %v: UseAs not default: %v", s.Name, b.Name, b.UseAs)
					}
					if *b.Balance.Algorithm != "roundrobin" {
						t.Errorf("%v: %v: Balance not roundrobin: %v", s.Name, b.Name, b.Balance.Algorithm)
					}
					if *b.Forwardfor.Enabled != "enabled" {
						t.Errorf("%v: %v: Forwardfor not enabled: %v", s.Name, b.Name, b.Forwardfor.Enabled)
					}
					if b.Mode != "http" {
						t.Errorf("%v: %v: Protocol not http: %v", s.Name, b.Name, b.Mode)
					}
					if len(b.Servers) != 2 {
						t.Errorf("%v: %v: Got %v servers, expected 2", s.Name, b.Name, len(b.Servers))
					}
					for _, srv := range b.Servers {
						if srv.Name != "webserv" && srv.Name != "webserv2" {
							t.Errorf("Expected only webserv or webserv2 servers, %v found", srv.Name)
						}
						if srv.Address != "192.168.1.1" {
							t.Errorf("%v: %v: %v: Address not 192.168.1.1: %v", s.Name, b.Name, srv.Name, srv.Address)
						}
						if *srv.Port != 9300 && *srv.Port != 9200 {
							t.Errorf("%v: %v: %v: Port not 9300 or 9200: %v", s.Name, b.Name, srv.Name, *srv.Port)
						}
						if srv.Ssl != "enabled" {
							t.Errorf("%v: %v: %v: Ssl not enabled: %v", s.Name, b.Name, srv.Name, srv.Ssl)
						}
						if *srv.Weight != 10 {
							t.Errorf("%v: %v: %v: Weight not 10: %v", s.Name, b.Name, srv.Name, *srv.Weight)
						}
					}
				case "test_2":
					if b.UseAs != "conditional" {
						t.Errorf("%v: %v: UseAs not conditional: %v", s.Name, b.Name, b.UseAs)
					}
					if b.Cond != "if" {
						t.Errorf("%v: %v: Cond not if: %v", s.Name, b.Name, b.Cond)
					}
					if b.CondTest != "TRUE" {
						t.Errorf("%v: %v: CondTest not TRUE: %v", s.Name, b.Name, b.CondTest)
					}
					if *b.Balance.Algorithm != "roundrobin" {
						t.Errorf("%v: %v: Balance not roundrobin: %v", s.Name, b.Name, b.Balance.Algorithm)
					}
					if *b.Forwardfor.Enabled != "enabled" {
						t.Errorf("%v: %v: Forwardfor not enabled: %v", s.Name, b.Name, b.Forwardfor.Enabled)
					}
					if b.Mode != "http" {
						t.Errorf("%v: %v: Mode not http: %v", s.Name, b.Name, b.Mode)
					}
					if len(b.Servers) != 0 {
						t.Errorf("%v: %v: Got %v servers, expected 0", s.Name, b.Name, len(b.Servers))
					}
				default:
					t.Errorf("%v: Expected only test or test_2 backends, %v found", s.Name, b.Name)
				}
			}
		case "test_2":
			if *s.Service.Maxconn != 2000 {
				t.Errorf("%v: MaxConnections not 2000: %v", s.Name, *s.Service.Maxconn)
			}
			if s.Service.Mode != "http" {
				t.Errorf("%v: Protocol not http: %v", s.Name, s.Service.Mode)
			}
			if s.Service.HTTPConnectionMode != "httpclose" {
				t.Errorf("%v: HTTPConnectionMode not httpclose: %v", s.Name, s.Service.HTTPConnectionMode)
			}
			if len(s.Service.Listeners) != 0 {
				t.Errorf("%v: Got %v listeners, expected 0", s.Name, len(s.Service.Listeners))
			}
			for _, b := range s.Farms {
				if b.Name == "test_2" {
					if b.UseAs != "default" {
						t.Errorf("%v: %v: UseAs not default: %v", s.Name, b.Name, b.UseAs)
					}
					if *b.Balance.Algorithm != "roundrobin" {
						t.Errorf("%v: %v: Balance not roundrobin: %v", s.Name, b.Name, b.Balance.Algorithm)
					}
					if *b.Forwardfor.Enabled != "enabled" {
						t.Errorf("%v: %v: Forwardfor not enabled: %v", s.Name, b.Name, *b.Forwardfor.Enabled)
					}
					if b.Mode != "http" {
						t.Errorf("%v: %v: Protocol not http: %v", s.Name, b.Name, b.Mode)
					}
					if len(b.Servers) != 0 {
						t.Errorf("%v: %v: Got %v servers, expected 0", s.Name, b.Name, len(b.Servers))
					}
				} else {
					t.Errorf("%v: Expected only test_2 backend, %v found", s.Name, b.Name)
				}
			}
		default:
			t.Errorf("Expected only test or test_2 sites, %v found", s.Name)
		}
	}
}

func TestGetSite(t *testing.T) { //nolint:gocognit,gocyclo
	v, s, err := client.GetSite("test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	if s.Name != "test" {
		t.Errorf("Name not test: %v", s.Name)
	}
	if *s.Service.Maxconn != 2000 {
		t.Errorf("%v: MaxConnections not 2000: %v", s.Name, *s.Service.Maxconn)
	}
	if s.Service.Mode != "http" {
		t.Errorf("%v: Protocol not http: %v", s.Name, s.Service.Mode)
	}
	if s.Service.HTTPConnectionMode != "httpclose" {
		t.Errorf("%v: HTTPConnectionMode not httpclose: %v", s.Name, s.Service.HTTPConnectionMode)
	}
	if len(s.Service.Listeners) != 2 {
		t.Errorf("%v: Got %v listeners, expected 2", s.Name, len(s.Service.Listeners))
	}
	for _, l := range s.Service.Listeners {
		if l.Name != "webserv" && l.Name != "webserv2" {
			t.Errorf("Expected only webserv or webserv2 listeners, %v found", l.Name)
		}
		if l.Address != "192.168.1.1" {
			t.Errorf("%v: Address not 192.168.1.1: %v", l.Name, l.Address)
		}
		if *l.Port != 80 && *l.Port != 8080 {
			t.Errorf("%v: Port not 80 or 8080: %v", l.Name, *l.Port)
		}
	}
	for _, b := range s.Farms {
		if b.Name == "test" {
			if b.UseAs != "default" {
				t.Errorf("%v: %v: UseAs not default: %v", s.Name, b.Name, b.UseAs)
			}
			if *b.Balance.Algorithm != "roundrobin" {
				t.Errorf("%v: %v: Balance not roundrobin: %v", s.Name, b.Name, b.Balance.Algorithm)
			}
			if *b.Forwardfor.Enabled != "enabled" {
				t.Errorf("%v: %v: HTTPXffHeaderInsert not enabled: %v", s.Name, b.Name, *b.Forwardfor.Enabled)
			}
			if b.Mode != "http" {
				t.Errorf("%v: %v: Protocol not http: %v", s.Name, b.Name, b.Mode)
			}
			if len(b.Servers) != 2 {
				t.Errorf("%v: %v: Got %v servers, expected 2", s.Name, b.Name, len(b.Servers))
			}
			for _, srv := range b.Servers {
				if srv.Name != "webserv" && srv.Name != "webserv2" {
					t.Errorf("Expected only webserv or webserv2 servers, %v found", srv.Name)
				}
				if srv.Address != "192.168.1.1" {
					t.Errorf("%v: %v: %v: Address not 192.168.1.1: %v", s.Name, b.Name, srv.Name, srv.Address)
				}
				if *srv.Port != 9300 && *srv.Port != 9200 {
					t.Errorf("%v: %v: %v: Port not 9300 or 9200: %v", s.Name, b.Name, srv.Name, *srv.Port)
				}
				if srv.Ssl != "enabled" {
					t.Errorf("%v: %v: %v: Ssl not enabled: %v", s.Name, b.Name, srv.Name, srv.Ssl)
				}
				if *srv.Weight != 10 {
					t.Errorf("%v: %v: %v: Weight not 10: %v", s.Name, b.Name, srv.Name, *srv.Weight)
				}
			}
		} else if b.Name == "test_2" {
			if b.UseAs != "conditional" {
				t.Errorf("%v: %v: UseAs not conditional: %v", s.Name, b.Name, b.UseAs)
			}
			if b.Cond != "if" {
				t.Errorf("%v: %v: Cond not if: %v", s.Name, b.Name, b.Cond)
			}
			if b.CondTest != "TRUE" {
				t.Errorf("%v: %v: CondTest not TRUE: %v", s.Name, b.Name, b.CondTest)
			}
			if *b.Balance.Algorithm != "roundrobin" {
				t.Errorf("%v: %v: Balance not roundrobin: %v", s.Name, b.Name, b.Balance.Algorithm)
			}
			if *b.Forwardfor.Enabled != "enabled" {
				t.Errorf("%v: %v: Forwardfor not enabled: %v", s.Name, b.Name, *b.Forwardfor.Enabled)
			}
			if b.Mode != "http" {
				t.Errorf("%v: %v: Mode not http: %v", s.Name, b.Name, b.Mode)
			}
			if len(b.Servers) != 0 {
				t.Errorf("%v: %v: Got %v servers, expected 0", s.Name, b.Name, len(b.Servers))
			}
		}
	}

	_, err = s.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}
}

func TestCreateEditDeleteSite(t *testing.T) {
	// TestCreateSite
	mConn := int64(2000)
	port := int64(5000)
	enabled := "enabled"
	balanceAlgorithm := "uri"
	s := &models.Site{
		Name: "created",
		Service: &models.SiteService{
			Mode:    "http",
			Maxconn: &mConn,
			Listeners: []*models.Bind{
				&models.Bind{
					Name:    "created1",
					Address: "127.0.0.1",
					Port:    &port,
				},
				&models.Bind{
					Name:    "created2",
					Address: "127.0.0.2",
					Port:    &port,
				},
			},
		},
		Farms: []*models.SiteFarm{
			&models.SiteFarm{
				Name:       "createdBck",
				Balance:    &models.Balance{Algorithm: &balanceAlgorithm},
				UseAs:      "default",
				Forwardfor: &models.Forwardfor{Enabled: &enabled},
				Servers: []*models.Server{
					&models.Server{
						Name:    "created1",
						Address: "127.0.1.1",
						Port:    &port,
						Ssl:     "enabled",
					},
					&models.Server{
						Name:    "created2",
						Address: "127.0.1.2",
						Port:    &port,
						Ssl:     "enabled",
					},
				},
			},
			&models.SiteFarm{
				Name:       "createdBck2",
				Balance:    &models.Balance{Algorithm: &balanceAlgorithm},
				UseAs:      "conditional",
				Cond:       "if",
				CondTest:   "TRUE",
				Forwardfor: &models.Forwardfor{Enabled: &enabled},
			},
		},
	}

	err := client.CreateSite(s, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, site, err := client.GetSite("created", "")
	if err != nil {
		t.Error(err.Error())
	}

	if !siteDeepEqual(site, s, t) {
		fmt.Printf("Created site: %v\n", *site)
		fmt.Printf("Given site: %v\n", *s)
		t.Error("Created site not equal to given site")
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	err = client.CreateSite(s, "", version)
	if err == nil {
		t.Error("Should throw error site already exists")
		version++
	}

	// TestEditSite
	editBalanceAlgorithm := "roundrobin"
	s = &models.Site{
		Name: "created",
		Service: &models.SiteService{
			Mode:    "tcp",
			Maxconn: &mConn,
			Listeners: []*models.Bind{
				&models.Bind{
					Name:    "created1",
					Address: "127.0.0.2",
					Port:    &port,
				},
			},
		},
		Farms: []*models.SiteFarm{
			&models.SiteFarm{
				Name:       "createdBck3",
				Balance:    &models.Balance{Algorithm: &balanceAlgorithm},
				UseAs:      "conditional",
				Cond:       "if",
				CondTest:   "TRUE",
				Forwardfor: &models.Forwardfor{Enabled: &enabled},
				Servers: []*models.Server{
					&models.Server{
						Name:    "created3",
						Address: "127.0.1.2",
						Port:    &port,
						Ssl:     "enabled",
					},
				},
			},
			&models.SiteFarm{
				Name:    "createdBck2",
				Balance: &models.Balance{Algorithm: &editBalanceAlgorithm},
				UseAs:   "default",
				Servers: []*models.Server{
					&models.Server{
						Name:    "created2",
						Address: "127.0.1.2",
						Port:    &port,
					},
				},
			},
		},
	}

	err = client.EditSite("created", s, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, site, err = client.GetSite("created", "")
	if err != nil {
		t.Error(err.Error())
	}

	if !siteDeepEqual(site, s, t) {
		fmt.Printf("Edited site: %v\n", *site)
		fmt.Printf("Given site: %v\n", *s)
		t.Error("Edited site not equal to given site")
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	// TestDeleteSite
	err = client.DeleteSite("created", "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	if v, _ := client.GetVersion(""); v != version {
		t.Error("Version not incremented")
	}

	_, _, err = client.GetSite("created", "")
	if err == nil {
		t.Error("DeleteSite failed, site test still exists")
	}

	err = client.DeleteSite("doesnotexist", "", version)
	if err == nil {
		t.Error("Should throw error, non existant site")
	}
}

func siteDeepEqual(x, y *models.Site, t *testing.T) bool { //nolint:gocognit
	if x.Name != y.Name {
		return false
	}

	// check frontend listeners
	if len(x.Service.Listeners) != len(y.Service.Listeners) {
		return false
	}
	if !assert.ElementsMatch(t, x.Service.Listeners, y.Service.Listeners) {
		return false
	}
	// Check Service
	if x.Service.HTTPConnectionMode != y.Service.HTTPConnectionMode {
		return false
	}
	if *x.Service.Maxconn != *y.Service.Maxconn {
		return false
	}
	if x.Service.Mode != y.Service.Mode {
		return false
	}

	// Check Farms
	if len(x.Farms) != len(y.Farms) {
		return false
	}
	backends := make([]interface{}, len(y.Farms))
	for i := range x.Farms {
		backends[i] = y.Farms[i]
	}
	for _, b := range x.Farms {
		b2Interface := misc.GetObjByField(backends, "Name", b.Name)
		if b2Interface == nil {
			return false
		}
		b2 := b2Interface.(*models.SiteFarm)
		// Compare backends
		if !reflect.DeepEqual(b.Forwardfor, b2.Forwardfor) {
			return false
		}
		b.Forwardfor = nil
		b2.Forwardfor = nil
		if *b.Balance.Algorithm != *b2.Balance.Algorithm {
			return false
		}
		if b.Balance.HdrName != b2.Balance.HdrName {
			return false
		}
		if b.Balance.HdrUseDomainOnly != b2.Balance.HdrUseDomainOnly {
			return false
		}
		if b.Balance.RandomDraws != b2.Balance.RandomDraws {
			return false
		}
		if b.Balance.RdpCookieName != b2.Balance.RdpCookieName {
			return false
		}
		if b.Balance.URIDepth != b2.Balance.URIDepth {
			return false
		}
		if b.Balance.URILen != b2.Balance.URILen {
			return false
		}
		if b.Balance.URIWhole != b2.Balance.URIWhole {
			return false
		}
		if b.Balance.URLParam != b2.Balance.URLParam {
			return false
		}
		if b.Balance.URLParamCheckPost != b2.Balance.URLParamCheckPost {
			return false
		}
		if b.Balance.URLParamMaxWait != b2.Balance.URLParamMaxWait {
			return false
		}
		if b.Mode != b2.Mode {
			return false
		}
		if b.Cond != b2.Cond {
			return false
		}
		if b.CondTest != b2.CondTest {
			return false
		}
		if b.UseAs != b2.UseAs {
			return false
		}

		// Compare Servers
		if len(b.Servers) != len(b2.Servers) {
			return false
		}
		if !assert.ElementsMatch(t, b.Servers, b2.Servers) {
			return false
		}
	}
	return true
}
