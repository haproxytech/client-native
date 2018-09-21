package configuration

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	strfmt "github.com/go-openapi/strfmt"
	"github.com/haproxytech/client-native/misc"
	"github.com/haproxytech/models"
)

// GetSites returns a struct with configuration version and an array of
// configured sites. Returns error on fail.
func (c *LBCTLConfigurationClient) GetSites() (*models.GetSitesOKBody, error) {
	response, err := c.executeLBCTL("l7-dump", "")
	if err != nil {
		return nil, err
	}

	sites := c.parseSites(response)

	v, err := c.GetVersion()
	if err != nil {
		return nil, err
	}

	return &models.GetSitesOKBody{Version: v, Data: sites}, nil
}

// GetSite returns a struct with configuration version and a requested site.
// Returns error on fail or if backend does not exist.
func (c *LBCTLConfigurationClient) GetSite(name string) (*models.GetSiteOKBody, error) {
	response, err := c.executeLBCTL("l7-dump", "")
	if err != nil {
		return nil, err
	}

	site := c.parseSite(response, name)

	v, err := c.GetVersion()
	if err != nil {
		return nil, err
	}

	return &models.GetSiteOKBody{Version: v, Data: site}, nil
}

// CreateSite creates a site in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *LBCTLConfigurationClient) CreateSite(data *models.Site, transactionID string, version int64) error {
	var res []error
	var err error

	validationErr := data.Validate(strfmt.Default)
	if validationErr != nil {
		return validationErr
	}
	// start an implicit transaction for create site (multiple operations required) if not already given
	t, err := c.checkTransactionOrVersion(transactionID, version, true)
	if err != nil {
		return err
	}

	//create frontend
	skip := []string{"Listeners"}
	err = c.createObject(*data.Name, "service", "", "", data, skip, t, version)
	if err != nil {
		res = append(res, err)
	}

	//create listeners
	for _, l := range data.Frontend.Listeners {
		//sanitize name
		if l.Name == "" {
			l.Name = l.Address + ":" + string(*l.Port)
		}
		err = c.createObject(l.Name, "listener", *data.Name, "", l, nil, t, version)
		if err != nil {
			res = append(res, err)
		}
	}

	//create backends
	for _, b := range data.Backends {
		skip = []string{"Servers", "UseAs", "Cond", "CondTest"}
		err = c.createObject(b.Name, "farm", "", "", b, nil, t, version)
		if err != nil {
			res = append(res, err)
		}
		//create servers
		for _, s := range b.Servers {
			//sanitize name
			if s.Name == "" {
				s.Name = s.Address + ":" + string(*s.Port)
			}
			err = c.createObject(s.Name, "server", b.Name, "", s, nil, t, version)
			if err != nil {
				res = append(res, err)
			}
		}
		//create bck-frontend relations
		err = c.createBckFrontendRels(*data.Name, b, false, t)
		if err != nil {
			res = append(res, err)
		}
	}
	if len(res) > 0 {
		return CompositeTransactionError(res...)
	}

	err = c.CommitTransaction(t)
	if err != nil {
		return err
	}

	return nil
}

// EditSite edits a site in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *LBCTLConfigurationClient) EditSite(name string, data *models.Site, transactionID string, version int64) error {
	var res []error
	var err error

	t, err := c.checkTransactionOrVersion(transactionID, version, true)
	if err != nil {
		return err
	}

	site, err := c.GetSite(name)
	if err != nil {
		return err
	}
	confS := site.Data

	//edit frontend
	if !reflect.DeepEqual(data.Frontend, confS.Frontend) {
		err := c.editObject(name, "service", "", "", data.Frontend, confS.Frontend, []string{"Listeners"}, transactionID, version)
		if err != nil {
			res = append(res, err)
		}
		//compare listeners
		if !reflect.DeepEqual(confS.Frontend.Listeners, data.Frontend.Listeners) {
			//add missing listeners by name, edit existing
			for _, l := range data.Frontend.Listeners {
				listeners := make([]interface{}, len(confS.Frontend.Listeners))
				for i := range confS.Frontend.Listeners {
					listeners[i] = confS.Frontend.Listeners[i]
				}

				confLIface := misc.GetObjByField(listeners, "Name", l.Name)
				if confLIface == nil {
					// create
					err = c.createObject(l.Name, "listener", *data.Name, "", l, nil, t, version)
					if err != nil {
						res = append(res, err)
					}
				} else {
					confL := confLIface.(*models.SiteFrontendListenersItems)
					if !reflect.DeepEqual(l, confL) {
						//edit
						err = c.editObject(l.Name, "listener", *data.Name, "", l, confL, nil, t, version)
						if err != nil {
							res = append(res, err)
						}
					} else {
						continue
					}
				}
			}
			//delete non existing listeners
			for _, l := range confS.Frontend.Listeners {
				listeners := make([]interface{}, len(data.Frontend.Listeners))
				for i := range data.Frontend.Listeners {
					listeners[i] = data.Frontend.Listeners[i]
				}
				if misc.GetObjByField(listeners, "Name", l.Name) == nil {
					err = c.deleteObject(l.Name, "listener", *data.Name, "", t, version)
					if err != nil {
						res = append(res, err)
					}
				}
			}
		}
	}
	// check if backends changed
	if !reflect.DeepEqual(confS.Backends, data.Backends) {
		for _, b := range data.Backends {
			// add missing backends
			bcks := make([]interface{}, len(confS.Backends))
			for i := range confS.Backends {
				bcks[i] = confS.Backends[i]
			}
			confBIface := misc.GetObjByField(bcks, "Name", b.Name)
			if confBIface == nil {
				err = c.createObject(b.Name, "farm", "", "", b, []string{"Servers", "UseAs", "Cond", "CondTest"}, t, version)
				if err != nil {
					res = append(res, err)
				}
				for _, s := range b.Servers {
					//sanitize name
					if s.Name == "" {
						s.Name = s.Address + ":" + string(*s.Port)
					}
					err = c.createObject(s.Name, "server", b.Name, "", s, nil, t, version)
					if err != nil {
						res = append(res, err)
					}
				}
				//create bck-frontend relations
				err = c.createBckFrontendRels(name, b, false, t)
				if err != nil {
					res = append(res, err)
				}
			} else {
				confB := confBIface.(*models.SiteBackendsItems)
				if !reflect.DeepEqual(b, confB) {
					// check if use as has changed
					if b.UseAs != confB.UseAs {
						err = c.createBckFrontendRels(name, b, true, t)
						if err != nil {
							res = append(res, err)
						}
					}
					err = c.editObject(b.Name, "farm", "", "", b, confB, []string{"Servers", "UseAs", "Cond", "CondTest"}, t, version)
					if err != nil {
						res = append(res, err)
					}
					for _, srv := range b.Servers {
						servers := make([]interface{}, len(b.Servers))
						for i := range b.Servers {
							servers[i] = b.Servers[i]
						}
						confSrvIFace := misc.GetObjByField(servers, "Name", srv.Name)
						if confSrvIFace == nil {
							// create
							err = c.createObject(srv.Name, "server", b.Name, "", srv, nil, t, version)
							if err != nil {
								res = append(res, err)
							}
						} else {
							confSrv := confSrvIFace.(*models.SiteBackendsItemsServersItems)
							if !reflect.DeepEqual(srv, confSrv) {
								//edit
								err = c.editObject(srv.Name, "server", b.Name, "", srv, confSrv, nil, t, version)
								if err != nil {
									res = append(res, err)
								}
							} else {
								continue
							}
						}
					}
					//delete non existing servers
					for _, srv := range b.Servers {
						servers := make([]interface{}, len(confB.Servers))
						for i := range confB.Servers {
							servers[i] = confB.Servers[i]
						}
						if misc.GetObjByField(servers, "Name", srv.Name) == nil {
							err = c.deleteObject(srv.Name, "server", b.Name, "", t, version)
							if err != nil {
								res = append(res, err)
							}
						}
					}
				}
			}
		}
		// delete non existing backends and remove uses in frontends
		for _, b := range confS.Backends {
			bcks := make([]interface{}, len(data.Backends))
			for i := range data.Backends {
				bcks[i] = data.Backends[i]
			}
			if misc.GetObjByField(bcks, "Name", b.Name) == nil {
				// default_bck
				if b.UseAs == "default" {
					// remove default bck option on frontend
					err = c.removeDefaultBckToFrontend(name, t)
					if err != nil {
						res = append(res, err)
					}
				} else {
					// find the correct usefarm and remove it
					err = c.removeUseFarm(name, b.Name, t)
					if err != nil {
						res = append(res, err)
					}
				}
				err = c.deleteObject(b.Name, "farm", "", "", t, version)
				if err != nil {
					res = append(res, err)
				}
			}
		}
	}

	if len(res) > 0 {
		return CompositeTransactionError(res...)
	}
	err = c.CommitTransaction(t)
	if err != nil {
		return err
	}

	return nil
}

// DeleteSite deletes a site in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *LBCTLConfigurationClient) DeleteSite(name string, transactionID string, version int64) error {
	var res []error
	var err error

	// start an implicit transaction for delete site (multiple operations required) if not already given
	t, err := c.checkTransactionOrVersion(transactionID, version, true)
	if err != nil {
		return err
	}

	site, err := c.GetSite(name)
	if err != nil {
		return err
	}

	err = c.DeleteFrontend(*site.Data.Name, t, 0)
	if err != nil {
		res = append(res, err)
	}
	for _, b := range site.Data.Backends {
		err = c.DeleteBackend(b.Name, t, version)
		if err != nil {
			res = append(res, err)
		}
	}

	if len(res) > 0 {
		return CompositeTransactionError(res...)
	}

	err = c.CommitTransaction(t)
	if err != nil {
		return err
	}
	return nil
}

func (c *LBCTLConfigurationClient) parseSite(response string, name string) *models.Site {
	bckCache := make(map[string]*models.SiteBackendsItems)
	lCache := make([]*models.SiteFrontendListenersItems, 0, 1)
	sCache := make(map[string][]*models.SiteBackendsItemsServersItems)
	frBckRelsCache := make(map[string][]map[string]string)

	site := &models.Site{}

	for _, obj := range strings.Split(response, "\n\n") {
		if strings.TrimSpace(obj) == "" {
			continue
		}
		if strings.HasPrefix(obj, ".service") {
			f := &models.SiteFrontend{}
			n := strings.TrimSpace(obj[strings.Index(obj, ".service ")+9 : strings.Index(obj, "\n")])
			if n != name {
				continue
			}

			site.Name = &n
			c.parseObject(obj, f)
			site.Frontend = f

			frBckRelsCache[*site.Name] = make([]map[string]string, 0, 1)
			// parse frontend-backend relations
			if strings.Contains(obj, "+default_farm") {
				bckName := strings.TrimSpace(obj[strings.Index(obj, "+default_farm ")+14:])
				bckName = bckName[:strings.Index(bckName, "\n")]
				frBckRelsCache[*site.Name] = append(frBckRelsCache[*site.Name], map[string]string{bckName: "default"})
			}
		}

		if strings.HasPrefix(obj, ".farm") {
			b := &models.SiteBackendsItems{}
			b.Name = strings.TrimSpace(obj[strings.Index(obj, ".farm ")+6 : strings.Index(obj, "\n")])
			c.parseObject(obj, b)
			bckCache[b.Name] = b
		}

		if strings.HasPrefix(obj, ".listener") {
			name, parent := splitHeaderLine(obj)
			l := &models.SiteFrontendListenersItems{}
			if parent != name {
				continue
			}
			c.parseObject(obj, l)
			lCache = append(lCache, l)
		}

		if strings.HasPrefix(obj, ".server") {
			name, parent := splitHeaderLine(obj)
			if name == "" {
				continue
			}
			s := &models.SiteBackendsItemsServersItems{}
			c.parseObject(obj, s)
			if a, ok := sCache[parent]; ok {
				sCache[parent] = append(a, s)
			} else {
				a := make([]*models.SiteBackendsItemsServersItems, 0, 1)
				sCache[parent] = append(a, s)
			}
		}

		if strings.HasPrefix(obj, ".usefarm") {
			f := ""
			b := ""
			val := ""
			for _, line := range strings.Split(obj, "\n") {
				if strings.HasPrefix(line, ".usefarm") {
					w := strings.Split(line, " ")
					f = w[len(w)-1]
				}
				if strings.HasPrefix(line, "+target_farm") {
					b = strings.TrimSpace(strings.SplitN(line, " ", 2)[1])
				}
				if strings.HasPrefix(line, "+cond") {
					w := strings.Split(line, " ")
					val = w[1] + " " + val
				}
				if strings.HasPrefix(line, "+cond_test") {
					w := strings.Split(line, " ")
					val = val + " " + w[1]
				}
			}

			if _, ok := frBckRelsCache[f]; !ok {
				frBckRelsCache[f] = make([]map[string]string, 0, 1)
			}

			frBckRelsCache[f] = append(frBckRelsCache[f], map[string]string{b: val})
		}
	}

	// get listeners for frontend
	site.Frontend.Listeners = lCache

	// add backends
	for _, bckObj := range frBckRelsCache[*site.Name] {
		for bck, val := range bckObj {
			b := bckCache[bck]
			if val == "default" {
				b.UseAs = val
			} else {
				b.UseAs = "conditional"
				split := strings.SplitN(val, " ", 1)
				b.Cond = strings.TrimSpace(split[0])
				b.CondTest = strings.TrimSpace(split[1])
			}

			site.Backends = append(site.Backends, b)
		}
	}

	// get servers for backends
	for _, bck := range site.Backends {
		bck.Servers = sCache[bck.Name]
	}
	return site
}

func (c *LBCTLConfigurationClient) parseSites(response string) models.Sites {
	bckCache := make(map[string]*models.SiteBackendsItems)
	lCache := make(map[string][]*models.SiteFrontendListenersItems)
	sCache := make(map[string][]*models.SiteBackendsItemsServersItems)
	sites := make([]*models.Site, 0, 1)
	frBckRelsCache := make(map[string][]map[string]string)

	for _, obj := range strings.Split(response, "\n\n") {
		site := &models.Site{}
		if strings.TrimSpace(obj) == "" {
			continue
		}
		if strings.HasPrefix(obj, ".service") {
			f := &models.SiteFrontend{}
			n := strings.TrimSpace(obj[strings.Index(obj, ".service ")+9 : strings.Index(obj, "\n")])
			if n == "" {
				continue
			}
			site.Name = &n
			c.parseObject(obj, f)
			site.Frontend = f

			frBckRelsCache[*site.Name] = make([]map[string]string, 0, 1)
			// parse frontend-backend relations
			if strings.Contains(obj, "+default_farm") {
				bckName := strings.TrimSpace(obj[strings.Index(obj, "+default_farm ")+14:])
				bckName = bckName[:strings.Index(bckName, "\n")]
				frBckRelsCache[*site.Name] = append(frBckRelsCache[*site.Name], map[string]string{bckName: "default"})
			}
		}

		if strings.HasPrefix(obj, ".farm") {
			b := &models.SiteBackendsItems{}
			b.Name = strings.TrimSpace(obj[strings.Index(obj, ".farm ")+6 : strings.Index(obj, "\n")])
			c.parseObject(obj, b)
			bckCache[b.Name] = b
		}

		if strings.HasPrefix(obj, ".listener") {
			name, parent := splitHeaderLine(obj)
			if name == "" {
				continue
			}
			l := &models.SiteFrontendListenersItems{}
			c.parseObject(obj, l)
			if a, ok := lCache[parent]; ok {
				lCache[parent] = append(a, l)
			} else {
				a := make([]*models.SiteFrontendListenersItems, 0, 1)
				lCache[parent] = append(a, l)
			}
		}

		if strings.HasPrefix(obj, ".server") {
			name, parent := splitHeaderLine(obj)
			if name == "" {
				continue
			}
			s := &models.SiteBackendsItemsServersItems{}
			c.parseObject(obj, s)
			if a, ok := sCache[parent]; ok {
				sCache[parent] = append(a, s)
			} else {
				a := make([]*models.SiteBackendsItemsServersItems, 0, 1)
				sCache[parent] = append(a, s)
			}
		}

		if strings.HasPrefix(obj, ".usefarm") {
			f := ""
			b := ""
			val := ""
			for _, line := range strings.Split(obj, "\n") {
				if strings.HasPrefix(line, ".usefarm") {
					w := strings.Split(line, " ")
					f = w[len(w)-1]
				}
				if strings.HasPrefix(line, "+target_farm") {
					b = strings.TrimSpace(strings.SplitN(line, " ", 2)[1])
				}
				if strings.HasPrefix(line, "+cond") {
					w := strings.Split(line, " ")
					val = w[1] + " " + val
				}
				if strings.HasPrefix(line, "+cond_test") {
					w := strings.Split(line, " ")
					val = val + " " + w[1]
				}
			}

			if _, ok := frBckRelsCache[f]; !ok {
				frBckRelsCache[f] = make([]map[string]string, 0, 1)
			}

			frBckRelsCache[f] = append(frBckRelsCache[f], map[string]string{b: val})
		}

		if site.Name != nil {
			sites = append(sites, site)
		}
	}
	for _, site := range sites {
		// get listeners for frontend
		site.Frontend.Listeners = lCache[*site.Name]

		// add backends
		for _, bckObj := range frBckRelsCache[*site.Name] {
			for bck, val := range bckObj {
				b := bckCache[bck]
				if val == "default" {
					b.UseAs = val
				} else {
					b.UseAs = "conditional"
					split := strings.SplitN(val, " ", 2)
					b.Cond = strings.TrimSpace(split[0])
					b.CondTest = strings.TrimSpace(split[1])
				}

				site.Backends = append(site.Backends, b)
			}
		}

		// get servers for backends
		for _, bck := range site.Backends {
			bck.Servers = sCache[bck.Name]
		}
	}

	return sites
}

// frontend backend relation helper methods
func (c *LBCTLConfigurationClient) removeUseFarm(frontend string, backend string, t string) error {
	ufs, err := c.getUseFarms(frontend)
	if err != nil {
		return err
	}
	for _, uf := range ufs {
		if uf.TargetFarm == backend {
			return c.deleteObject(string(*uf.ID), "usefarm", frontend, "service", t, 0)
		}
	}
	return nil
}

func (c *LBCTLConfigurationClient) createBckFrontendRels(name string, b *models.SiteBackendsItems, edit bool, t string) error {
	var res []error
	var err error
	id := int64(0)
	if b.UseAs == "default" {
		if edit {
			err = c.removeUseFarm(name, b.Name, t)
			if err != nil {
				res = append(res, err)
			}
		}
		err = c.addDefaultBckToFrontend(name, b.Name, t)
		if err != nil {
			res = append(res, err)
		}
	} else {
		if edit {
			err = c.removeDefaultBckToFrontend(name, t)
			if err != nil {
				res = append(res, err)
			}
		}
		if b.Cond == "" || b.CondTest == "" {
			res = append(res, fmt.Errorf("Backend %s set as conditional but no conditions provided", b.Name))
		} else {
			uf := &misc.UseFarm{
				ID:         &id,
				TargetFarm: b.Name,
				Cond:       b.Cond,
				CondTest:   b.CondTest,
			}
			err = c.createObject(string(*uf.ID), "usefarm", name, "frontend", uf, nil, t, 0)
			if err != nil {
				res = append(res, err)
			}
		}
	}
	if len(res) > 0 {
		return CompositeTransactionError(res...)
	}
	return nil
}

func (c *LBCTLConfigurationClient) addDefaultBckToFrontend(fName string, bName string, t string) error {
	response, err := c.executeLBCTL("l7-service-update", t, fName, "--default_farm", bName)
	if err != nil {
		return errors.New(err.Error() + ": " + response)
	}
	return nil
}

func (c *LBCTLConfigurationClient) removeDefaultBckToFrontend(fName string, t string) error {
	response, err := c.executeLBCTL("l7-service-update", t, fName, "--reset-default_farm")
	if err != nil {
		return errors.New(err.Error() + ": " + response)
	}
	return nil
}

func (c *LBCTLConfigurationClient) getUseFarms(parent string) ([]*misc.UseFarm, error) {
	response, err := c.executeLBCTL("l7-service-usefarm-dump", "", parent)
	if err != nil {
		return nil, err
	}

	useFarms := make([]*misc.UseFarm, 0, 1)
	for _, ufString := range strings.Split(response, "\n\n") {
		if strings.TrimSpace(ufString) == "" {
			continue
		}

		uf := &misc.UseFarm{}
		for _, line := range strings.Split(ufString, "\n") {
			if strings.HasPrefix(line, ".usefarm") {
				w := strings.Split(line, " ")
				idStr := w[len(w)-3]
				id, err := strconv.ParseInt(idStr, 10, 64)
				if err != nil {
					continue
				}
				uf.ID = &id
			}
			if strings.HasPrefix(line, "+target_farm") {
				w := strings.Split(line, " ")
				uf.TargetFarm = w[1]
			}
			if strings.HasPrefix(line, "+cond") {
				w := strings.Split(line, " ")
				uf.Cond = w[1]
			}
			if strings.HasPrefix(line, "+cond_test") {
				w := strings.Split(line, " ")
				uf.CondTest = w[1]
			}
			useFarms = append(useFarms, uf)
		}
	}
	return useFarms, nil
}
