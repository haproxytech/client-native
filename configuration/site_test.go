package configuration

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/haproxytech/client-native/misc"
	"github.com/haproxytech/models"
)

func TestGetSites(t *testing.T) {
	sites, err := client.GetSites("")
	if err != nil {
		t.Error(err.Error())
	}

	if len(sites.Data) != 2 {
		t.Errorf("%v sites returned, expected 2", len(sites.Data))
	}

	if sites.Version != version {
		t.Errorf("Version %v returned, expected %v", sites.Version, version)
	}

	for _, s := range sites.Data {
		if s.Name == "test" {
			if *s.Frontend.MaxConnections != 2000 {
				t.Errorf("%v: MaxConnections not 2000: %v", s.Name, *s.Frontend.MaxConnections)
			}
			if s.Frontend.Log != "enabled" {
				t.Errorf("%v: Log not enabled: %v", s.Name, s.Frontend.Log)
			}
			if s.Frontend.Protocol != "http" {
				t.Errorf("%v: Protocol not http: %v", s.Name, s.Frontend.Protocol)
			}
			if s.Frontend.HTTPConnectionMode != "passive-close" {
				t.Errorf("%v: HTTPConnectionMode not passive-close: %v", s.Name, s.Frontend.HTTPConnectionMode)
			}
			if len(s.Frontend.Listeners) != 2 {
				t.Errorf("%v: Got %v listeners, expected 2", s.Name, len(s.Frontend.Listeners))
			}
			for _, l := range s.Frontend.Listeners {
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
			for _, b := range s.Backends {
				if b.Name == "test" {
					if b.UseAs != "default" {
						t.Errorf("%v: %v: UseAs not default: %v", s.Name, b.Name, b.UseAs)
					}
					if b.Balance != "roundrobin" {
						t.Errorf("%v: %v: Balance not roundrobin: %v", s.Name, b.Name, b.Balance)
					}
					if b.HTTPXffHeaderInsert != "enabled" {
						t.Errorf("%v: %v: HTTPXffHeaderInsert not enabled: %v", s.Name, b.Name, b.HTTPXffHeaderInsert)
					}
					if b.Log != "enabled" {
						t.Errorf("%v: %v: Log not enabled: %v", s.Name, b.Name, b.Log)
					}
					if b.Protocol != "http" {
						t.Errorf("%v: %v: Protocol not http: %v", s.Name, b.Name, b.Protocol)
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
					if b.Balance != "roundrobin" {
						t.Errorf("%v: %v: Balance not roundrobin: %v", s.Name, b.Name, b.Balance)
					}
					if b.HTTPXffHeaderInsert != "enabled" {
						t.Errorf("%v: %v: HTTPXffHeaderInsert not enabled: %v", s.Name, b.Name, b.HTTPXffHeaderInsert)
					}
					if b.Log != "enabled" {
						t.Errorf("%v: %v: Log not enabled: %v", s.Name, b.Name, b.Log)
					}
					if b.Protocol != "http" {
						t.Errorf("%v: %v: Protocol not http: %v", s.Name, b.Name, b.Protocol)
					}
					if len(b.Servers) != 0 {
						t.Errorf("%v: %v: Got %v servers, expected 0", s.Name, b.Name, len(b.Servers))
					}
				} else {
					t.Errorf("%v: Expected only test or test_2 backends, %v found", s.Name, b.Name)
				}
			}
		} else if s.Name == "test_2" {
			if *s.Frontend.MaxConnections != 2000 {
				t.Errorf("%v: MaxConnections not 2000: %v", s.Name, *s.Frontend.MaxConnections)
			}
			if s.Frontend.Log != "enabled" {
				t.Errorf("%v: Log not enabled: %v", s.Name, s.Frontend.Log)
			}
			if s.Frontend.Protocol != "http" {
				t.Errorf("%v: Protocol not http: %v", s.Name, s.Frontend.Protocol)
			}
			if s.Frontend.HTTPConnectionMode != "passive-close" {
				t.Errorf("%v: HTTPConnectionMode not passive-close: %v", s.Name, s.Frontend.HTTPConnectionMode)
			}
			if len(s.Frontend.Listeners) != 0 {
				t.Errorf("%v: Got %v listeners, expected 0", s.Name, len(s.Frontend.Listeners))
			}
			for _, b := range s.Backends {
				if b.Name == "test_2" {
					if b.UseAs != "default" {
						t.Errorf("%v: %v: UseAs not default: %v", s.Name, b.Name, b.UseAs)
					}
					if b.Balance != "roundrobin" {
						t.Errorf("%v: %v: Balance not roundrobin: %v", s.Name, b.Name, b.Balance)
					}
					if b.HTTPXffHeaderInsert != "enabled" {
						t.Errorf("%v: %v: HTTPXffHeaderInsert not enabled: %v", s.Name, b.Name, b.HTTPXffHeaderInsert)
					}
					if b.Log != "enabled" {
						t.Errorf("%v: %v: Log not enabled: %v", s.Name, b.Name, b.Log)
					}
					if b.Protocol != "http" {
						t.Errorf("%v: %v: Protocol not http: %v", s.Name, b.Name, b.Protocol)
					}
					if len(b.Servers) != 0 {
						t.Errorf("%v: %v: Got %v servers, expected 0", s.Name, b.Name, len(b.Servers))
					}
				} else {
					t.Errorf("%v: Expected only test_2 backend, %v found", s.Name, b.Name)
				}
			}
		} else {
			t.Errorf("Expected only test or test_2 sites, %v found", s.Name)
		}
	}

	sJSON, err := sites.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	if !t.Failed() {
		fmt.Println("GetSites succesful\nResponse: \n" + string(sJSON) + "\n")
	}
}

func TestGetSite(t *testing.T) {
	site, err := client.GetSite("test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if site.Version != version {
		t.Errorf("Version %v returned, expected %v", site.Version, version)
	}

	s := site.Data

	if s.Name != "test" {
		t.Errorf("Name not test: %v", s.Name)
	}
	if *s.Frontend.MaxConnections != 2000 {
		t.Errorf("%v: MaxConnections not 2000: %v", s.Name, *s.Frontend.MaxConnections)
	}
	if s.Frontend.Log != "enabled" {
		t.Errorf("%v: Log not enabled: %v", s.Name, s.Frontend.Log)
	}
	if s.Frontend.Protocol != "http" {
		t.Errorf("%v: Protocol not http: %v", s.Name, s.Frontend.Protocol)
	}
	if s.Frontend.HTTPConnectionMode != "passive-close" {
		t.Errorf("%v: HTTPConnectionMode not passive-close: %v", s.Name, s.Frontend.HTTPConnectionMode)
	}
	if len(s.Frontend.Listeners) != 2 {
		t.Errorf("%v: Got %v listeners, expected 2", s.Name, len(s.Frontend.Listeners))
	}
	for _, l := range s.Frontend.Listeners {
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
	for _, b := range s.Backends {
		if b.Name == "test" {
			if b.UseAs != "default" {
				t.Errorf("%v: %v: UseAs not default: %v", s.Name, b.Name, b.UseAs)
			}
			if b.Balance != "roundrobin" {
				t.Errorf("%v: %v: Balance not roundrobin: %v", s.Name, b.Name, b.Balance)
			}
			if b.HTTPXffHeaderInsert != "enabled" {
				t.Errorf("%v: %v: HTTPXffHeaderInsert not enabled: %v", s.Name, b.Name, b.HTTPXffHeaderInsert)
			}
			if b.Log != "enabled" {
				t.Errorf("%v: %v: Log not enabled: %v", s.Name, b.Name, b.Log)
			}
			if b.Protocol != "http" {
				t.Errorf("%v: %v: Protocol not http: %v", s.Name, b.Name, b.Protocol)
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
			if b.Balance != "roundrobin" {
				t.Errorf("%v: %v: Balance not roundrobin: %v", s.Name, b.Name, b.Balance)
			}
			if b.HTTPXffHeaderInsert != "enabled" {
				t.Errorf("%v: %v: HTTPXffHeaderInsert not enabled: %v", s.Name, b.Name, b.HTTPXffHeaderInsert)
			}
			if b.Log != "enabled" {
				t.Errorf("%v: %v: Log not enabled: %v", s.Name, b.Name, b.Log)
			}
			if b.Protocol != "http" {
				t.Errorf("%v: %v: Protocol not http: %v", s.Name, b.Name, b.Protocol)
			}
			if len(b.Servers) != 0 {
				t.Errorf("%v: %v: Got %v servers, expected 0", s.Name, b.Name, len(b.Servers))
			}
		}
	}

	sJSON, err := site.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	if !t.Failed() {
		fmt.Println("GetSite succesful\nResponse: \n" + string(sJSON) + "\n")
	}
}

func TestCreateEditDeleteSite(t *testing.T) {
	// TestCreateSite
	mConn := int64(2000)
	port := int64(5000)
	s := &models.Site{
		Name: "created",
		Frontend: &models.SiteFrontend{
			Protocol:       "http",
			Log:            "enabled",
			MaxConnections: &mConn,
			Listeners: []*models.SiteFrontendListenersItems{
				&models.SiteFrontendListenersItems{
					Name:    "created1",
					Address: "127.0.0.1",
					Port:    &port,
				},
				&models.SiteFrontendListenersItems{
					Name:    "created2",
					Address: "127.0.0.2",
					Port:    &port,
				},
			},
		},
		Backends: []*models.SiteBackendsItems{
			&models.SiteBackendsItems{
				Name:                "createdBck",
				Balance:             "hash-uri",
				UseAs:               "default",
				Log:                 "enabled",
				HTTPXffHeaderInsert: "enabled",
				Servers: []*models.SiteBackendsItemsServersItems{
					&models.SiteBackendsItemsServersItems{
						Name:    "created1",
						Address: "127.0.1.1",
						Port:    &port,
						Ssl:     "enabled",
					},
					&models.SiteBackendsItemsServersItems{
						Name:    "created2",
						Address: "127.0.1.2",
						Port:    &port,
						Ssl:     "enabled",
					},
				},
			},
			&models.SiteBackendsItems{
				Name:                "createdBck2",
				Balance:             "hash-uri",
				UseAs:               "conditional",
				Cond:                "if",
				CondTest:            "TRUE",
				Log:                 "enabled",
				HTTPXffHeaderInsert: "enabled",
			},
		},
	}

	err := client.CreateSite(s, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	site, err := client.GetSite("created", "")
	if err != nil {
		t.Error(err.Error())
	}

	if !siteDeepEqual(site.Data, s, t) {
		fmt.Printf("Created site: %v\n", *site.Data)
		fmt.Printf("Given site: %v\n", *s)
		t.Error("Created site not equal to given site")
	}

	if site.Version != version {
		t.Errorf("Version %v returned, expected %v", site.Version, version)
	}

	err = client.CreateSite(s, "", version)
	if err == nil {
		t.Error("Should throw error site already exists")
		version++
	}

	if !t.Failed() {
		fmt.Println("CreateSite successful")
	}

	// TestEditSite
	s = &models.Site{
		Name: "created",
		Frontend: &models.SiteFrontend{
			Protocol:       "tcp",
			Log:            "enabled",
			MaxConnections: &mConn,
			Listeners: []*models.SiteFrontendListenersItems{
				&models.SiteFrontendListenersItems{
					Name:    "created1",
					Address: "127.0.0.2",
					Port:    &port,
				},
			},
		},
		Backends: []*models.SiteBackendsItems{
			&models.SiteBackendsItems{
				Name:                "createdBck3",
				Balance:             "hash-uri",
				UseAs:               "conditional",
				Cond:                "if",
				CondTest:            "TRUE",
				Log:                 "enabled",
				HTTPXffHeaderInsert: "enabled",
				Servers: []*models.SiteBackendsItemsServersItems{
					&models.SiteBackendsItemsServersItems{
						Name:    "created3",
						Address: "127.0.1.2",
						Port:    &port,
						Ssl:     "enabled",
					},
				},
			},
			&models.SiteBackendsItems{
				Name:    "createdBck2",
				Balance: "roundrobin",
				UseAs:   "default",
				Servers: []*models.SiteBackendsItemsServersItems{
					&models.SiteBackendsItemsServersItems{
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

	site, err = client.GetSite("created", "")
	if err != nil {
		t.Error(err.Error())
	}

	if !siteDeepEqual(site.Data, s, t) {
		fmt.Printf("Edited site: %v\n", *site.Data)
		fmt.Printf("Given site: %v\n", *s)
		t.Error("Edited site not equal to given site")
	}

	if site.Version != version {
		t.Errorf("Version %v returned, expected %v", site.Version, version)
	}

	if !t.Failed() {
		fmt.Println("EditSite successful")
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

	_, err = client.GetSite("created", "")
	if err == nil {
		t.Error("DeleteSite failed, site test still exists")
	}

	err = client.DeleteSite("doesnotexist", "", version)
	if err == nil {
		t.Error("Should throw error, non existant site")
	}

	if !t.Failed() {
		fmt.Println("DeleteSite successful")
	}
}

func siteDeepEqual(x, y *models.Site, t *testing.T) bool {
	if x.Name != y.Name {
		return false
	}

	// check frontend listeners
	if len(x.Frontend.Listeners) != len(x.Frontend.Listeners) {
		return false
	}
	if !assert.ElementsMatch(t, x.Frontend.Listeners, y.Frontend.Listeners) {
		return false
	}
	// Check Frontend
	if x.Frontend.HTTPConnectionMode != y.Frontend.HTTPConnectionMode {
		return false
	}
	if x.Frontend.Log != y.Frontend.Log {
		return false
	}
	if *x.Frontend.MaxConnections != *y.Frontend.MaxConnections {
		return false
	}
	if x.Frontend.Protocol != y.Frontend.Protocol {
		return false
	}

	// Check Backends
	if len(x.Backends) != len(y.Backends) {
		return false
	}
	backends := make([]interface{}, len(y.Backends))
	for i := range x.Backends {
		backends[i] = y.Backends[i]
	}
	for _, b := range x.Backends {
		b2Interface := misc.GetObjByField(backends, "Name", b.Name)
		if b2Interface == nil {
			return false
		}
		b2 := b2Interface.(*models.SiteBackendsItems)
		// Compare backends
		if b.Balance != b2.Balance {
			return false
		}
		if b.Protocol != b2.Protocol {
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
		if b.Log != b2.Log {
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

// func TestEditServer(t *testing.T) {
// 	port := int64(5300)
// 	s := &models.Server{
// 		Name:    "created",
// 		Address: "192.168.3.1",
// 		Port:    &port,
// 	}

// 	err := client.EditServer("created", "test", s, "", version)
// 	if err != nil {
// 		t.Error(err.Error())
// 	} else {
// 		version = version + 1
// 	}

// 	server, err := client.GetServer("created", "test")
// 	if err != nil {
// 		t.Error(err.Error())
// 	}
// 	sEdited := server.Data

// 	if !reflect.DeepEqual(sEdited, s) {
// 		fmt.Printf("Edited server: %v\n", sEdited)
// 		fmt.Printf("Given server: %v\n", s)
// 		t.Error("Edited server not equal to given server")
// 	}

// 	if server.Version != version {
// 		t.Errorf("Version %v returned, expected %v", server.Version, version)
// 	}

// 	if !t.Failed() {
// 		fmt.Println("EditServer successful")
// 	}
// }
