package configuration

import (
	"fmt"
	"reflect"
	"strconv"

	strfmt "github.com/go-openapi/strfmt"
	"github.com/haproxytech/client-native/misc"
	parser "github.com/haproxytech/config-parser"
	"github.com/haproxytech/models"
)

// GetSites returns a struct with configuration version and an array of
// configured sites. Returns error on fail.
func (c *Client) GetSites(transactionID string) (*models.GetSitesOKBody, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return nil, err
	}

	sites, err := c.parseSites(p)
	if err != nil {
		return nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return nil, err
	}

	return &models.GetSitesOKBody{Version: v, Data: sites}, nil
}

// GetSite returns a struct with configuration version and a requested site.
// Returns error on fail or if backend does not exist.
func (c *Client) GetSite(name string, transactionID string) (*models.GetSiteOKBody, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return nil, err
	}

	if !c.checkSectionExists(parser.Frontends, name, p) {
		return nil, NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("Site %s does not exist", name))
	}

	site := c.parseSite(name, p)
	if site == nil {
		return nil, NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("Site %s does not exist", name))
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return nil, err
	}

	return &models.GetSiteOKBody{Version: v, Data: site}, nil
}

// CreateSite creates a site in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *Client) CreateSite(data *models.Site, transactionID string, version int64) error {
	var res []error
	var err error

	if c.UseValidation {
		validationErr := data.Validate(strfmt.Default)
		if validationErr != nil {
			return NewConfError(ErrValidationError, validationErr.Error())
		}
	}
	// start an implicit transaction for create site (multiple operations required) if not already given
	p, t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	//create frontend
	frontend := serializeServiceToFrontend(data.Service, data.Name)

	if frontend != nil {
		err = c.CreateFrontend(frontend, t, 0)
		if err != nil {
			res = append(res, err)
		}
	}

	//create listeners
	for _, l := range data.Service.Listeners {
		//sanitize name
		if l.Name == "" {
			l.Name = l.Address + ":" + strconv.FormatInt(*l.Port, 10)
		}
		bind := serializeListenerToBind(l)
		if bind != nil {
			err = c.CreateBind(data.Name, bind, t, 0)
			if err != nil {
				res = append(res, err)
			}
		}
	}

	//create backends
	for _, b := range data.Farms {
		backend := serializeFarmToBackend(b)
		if backend == nil {
			continue
		}
		err = c.CreateBackend(backend, t, 0)
		if err != nil {
			res = append(res, err)
		}
		//create servers
		for _, s := range b.Servers {
			//sanitize name
			if s.Name == "" {
				s.Name = s.Address + ":" + strconv.FormatInt(*s.Port, 10)
			}
			server := serializeSiteServer(s)
			err = c.CreateServer(b.Name, server, t, 0)
			if err != nil {
				res = append(res, err)
			}
		}
		//create bck-frontend relations
		err = c.createBckFrontendRels(data.Name, b, false, t, p)
		if err != nil {
			res = append(res, err)
		}
	}
	if len(res) > 0 {
		return c.handleError(data.Name, "", "", t, transactionID == "", CompositeTransactionError(res...))
	}

	if err := c.saveData(p, t, transactionID == ""); err != nil {
		return err
	}

	return nil
}

// EditSite edits a site in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *Client) EditSite(name string, data *models.Site, transactionID string, version int64) error {
	var res []error
	var err error

	if c.UseValidation {
		validationErr := data.Validate(strfmt.Default)
		if validationErr != nil {
			return NewConfError(ErrValidationError, validationErr.Error())
		}
	}
	// start an implicit transaction for create site (multiple operations required) if not already given
	p, t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	site, err := c.GetSite(name, transactionID)
	if err != nil {
		return err
	}
	confS := site.Data

	//edit frontend
	if !reflect.DeepEqual(data.Service, confS.Service) {
		err := c.editService(data.Name, data.Service, t, p)
		if err != nil {
			res = append(res, err)
		}
		//compare listeners
		if !reflect.DeepEqual(confS.Service.Listeners, data.Service.Listeners) {
			//add missing listeners by name, edit existing
			for _, l := range data.Service.Listeners {
				listeners := make([]interface{}, len(confS.Service.Listeners))
				for i := range confS.Service.Listeners {
					listeners[i] = confS.Service.Listeners[i]
				}

				confLIface := misc.GetObjByField(listeners, "Name", l.Name)
				if confLIface == nil {
					// create
					bind := serializeListenerToBind(l)
					if bind != nil {
						err = c.CreateBind(data.Name, bind, t, 0)
						if err != nil {
							res = append(res, err)
						}
					}
				} else {
					confL := confLIface.(*models.SiteServiceListenersItems)
					if !reflect.DeepEqual(l, confL) {
						err := c.editListener(l.Name, data.Name, l, t)
						if err != nil {
							res = append(res, err)
						}
					} else {
						continue
					}
				}
			}
			//delete non existing listeners
			for _, l := range confS.Service.Listeners {
				listeners := make([]interface{}, len(data.Service.Listeners))
				for i := range data.Service.Listeners {
					listeners[i] = data.Service.Listeners[i]
				}
				if misc.GetObjByField(listeners, "Name", l.Name) == nil {
					err = c.DeleteBind(l.Name, data.Name, t, 0)
					if err != nil {
						res = append(res, err)
					}
				}
			}
		}
	}
	bcks := make([]interface{}, len(confS.Farms))
	for i := range confS.Farms {
		bcks[i] = confS.Farms[i]
	}
	defaultBck := ""
	// check if backends changed
	if !reflect.DeepEqual(confS.Farms, data.Farms) {
		for _, b := range data.Farms {
			// add missing backends
			confBIface := misc.GetObjByField(bcks, "Name", b.Name)
			if confBIface == nil {
				backend := serializeFarmToBackend(b)
				if b != nil {
					err = c.CreateBackend(backend, t, 0)
					if err != nil {
						res = append(res, err)
					}
					for _, s := range b.Servers {
						//sanitize name
						if s.Name == "" {
							s.Name = s.Address + ":" + strconv.FormatInt(*s.Port, 10)
						}
						server := serializeSiteServer(s)
						if server != nil {
							err := c.CreateServer(b.Name, server, t, 0)
							if err != nil {
								res = append(res, err)
							}
						}
					}
					if b.UseAs == "default" && defaultBck != "" {
						return NewConfError(ErrValidationError, fmt.Sprintf("Multiple default backends found in site: %v", name))
					} else if b.UseAs == "default" && defaultBck == "" {
						defaultBck = b.Name
					}
					//create bck-frontend relations
					err = c.createBckFrontendRels(name, b, false, t, p)
					if err != nil {
						res = append(res, err)
					}
				}
			} else {
				if b.UseAs == "default" && defaultBck != "" {
					return NewConfError(ErrValidationError, fmt.Sprintf("Multiple default backends found in site: %v", name))
				} else if b.UseAs == "default" && defaultBck == "" {
					defaultBck = b.Name
				}
				confB := confBIface.(*models.SiteFarmsItems)
				if !reflect.DeepEqual(b, confB) {
					// check if use as has changed
					if b.UseAs != confB.UseAs {
						err := c.createBckFrontendRels(name, b, true, t, p)
						if err != nil {
							res = append(res, err)
						}
					}
					err := c.editFarm(b.Name, b, t, p)
					if err != nil {
						res = append(res, err)
					}
					servers := make([]interface{}, len(confB.Servers))
					for i := range confB.Servers {
						servers[i] = confB.Servers[i]
					}
					for _, srv := range b.Servers {
						confSrvIFace := misc.GetObjByField(servers, "Name", srv.Name)
						if confSrvIFace == nil {
							// create
							server := serializeSiteServer(srv)
							if server != nil {
								err := c.CreateServer(b.Name, server, t, 0)
								if err != nil {
									res = append(res, err)
								}
							}
						} else {
							confSrv := confSrvIFace.(*models.SiteFarmsItemsServersItems)
							if !reflect.DeepEqual(srv, confSrv) {
								//edit
								err := c.editSiteServer(srv.Name, b.Name, srv, t)
								if err != nil {
									res = append(res, err)
								}
							} else {
								continue
							}
						}
					}
					servers = make([]interface{}, len(b.Servers))
					for i := range b.Servers {
						bcks[i] = b.Servers[i]
					}
					//delete non existing servers
					for _, srv := range confB.Servers {
						if misc.GetObjByField(servers, "Name", srv.Name) == nil {
							err := c.DeleteServer(srv.Name, b.Name, t, 0)
							if err != nil {
								res = append(res, err)
							}
						}
					}

				}
			}
		}
		bcks = make([]interface{}, len(data.Farms))
		for i := range data.Farms {
			bcks[i] = data.Farms[i]
		}
		// delete non existing backends and remove uses in frontends
		for _, b := range confS.Farms {
			if misc.GetObjByField(bcks, "Name", b.Name) == nil {
				// default_bck
				if b.UseAs == "conditional" {
					// find the correct usefarm and remove it
					err := c.removeUseFarm(name, b.Name, t, p)
					if err != nil {
						res = append(res, err)
					}
				}
				err := c.DeleteBackend(b.Name, t, 0)
				if err != nil {
					res = append(res, err)
				}
			}
		}
	}
	// remove default backend if no default backends specified
	if defaultBck == "" {
		err = c.removeDefaultBckToFrontend(name, t, p)
		if err != nil {
			res = append(res, err)
		}
	}

	if len(res) > 0 {
		return c.handleError(data.Name, "", "", t, transactionID == "", CompositeTransactionError(res...))
	}

	if err := c.saveData(p, t, transactionID == ""); err != nil {
		return err
	}

	return nil
}

// DeleteSite deletes a site in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *Client) DeleteSite(name string, transactionID string, version int64) error {
	var res []error
	var err error

	// start an implicit transaction for delete site (multiple operations required) if not already given
	p, t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	site, err := c.GetSite(name, t)
	if err != nil {
		return err
	}

	err = c.DeleteFrontend(site.Data.Name, t, 0)
	if err != nil {
		res = append(res, err)
	}

	for _, b := range site.Data.Farms {
		err = c.DeleteBackend(b.Name, t, 0)
		if err != nil {
			res = append(res, err)
		}
	}

	if len(res) > 0 {
		return c.handleError(name, "", "", t, transactionID == "", CompositeTransactionError(res...))
	}

	if err := c.saveData(p, t, transactionID == ""); err != nil {
		return err
	}

	return nil
}

func (c *Client) parseSites(p *parser.Parser) (models.Sites, error) {
	sites := models.Sites{}
	fNames, err := p.SectionsGet(parser.Frontends)
	if err != nil {
		return nil, err
	}

	for _, s := range fNames {
		site := c.parseSite(s, p)
		if site != nil {
			sites = append(sites, site)
		}
	}
	return sites, nil
}

func (c *Client) parseSite(s string, p *parser.Parser) *models.Site {
	frontend := &models.Frontend{Name: s}
	if err := c.parseSection(frontend, parser.Frontends, s, p); err != nil {
		return nil
	}
	site := &models.Site{
		Name: s,
		Service: &models.SiteService{
			HTTPConnectionMode: frontend.HTTPConnectionMode,
			Maxconn:            frontend.Maxconn,
			Mode:               frontend.Mode,
			Listeners:          c.parseServiceListeners(s, p),
		},
		Farms: []*models.SiteFarmsItems{},
	}

	// Find backends using default_backend and use_backends
	if frontend.DefaultBackend != "" {
		// parse default backend
		farm := c.parseFarm(frontend.DefaultBackend, "default", "", "", p)
		if farm != nil {
			site.Farms = append(site.Farms, farm)
		}
	}
	ubs, err := c.parseBackendSwitchingRules(s, p)
	if err == nil {
		for _, ub := range ubs {
			farm := c.parseFarm(ub.Name, "conditional", ub.Cond, ub.CondTest, p)
			if farm != nil {
				site.Farms = append(site.Farms, farm)
			}
		}
	}
	return site
}

func (c *Client) parseServiceListeners(service string, p *parser.Parser) []*models.SiteServiceListenersItems {
	listeners := []*models.SiteServiceListenersItems{}
	binds, err := c.parseBinds(service, p)
	if err == nil {
		for _, b := range binds {
			li := &models.SiteServiceListenersItems{
				Address:        b.Address,
				Name:           b.Name,
				Port:           b.Port,
				Ssl:            b.Ssl,
				SslCertificate: b.SslCertificate,
			}
			listeners = append(listeners, li)
		}
	}
	return listeners
}

func (c *Client) parseFarm(name string, useAs string, cond string, condTest string, p *parser.Parser) *models.SiteFarmsItems {
	backend := &models.Backend{Name: name}
	if err := c.parseSection(backend, parser.Backends, name, p); err == nil {
		farm := &models.SiteFarmsItems{
			UseAs:    useAs,
			Cond:     cond,
			CondTest: condTest,
			Mode:     backend.Mode,
			Name:     backend.Name,
			Servers:  c.parseFarmServers(backend.Name, p),
		}
		if backend.Forwardfor == "enabled" {
			farm.Forwardfor = true
		}
		if backend.Balance != nil {
			farm.Balance = &models.SiteFarmsItemsBalance{
				Algorithm: backend.Balance.Algorithm,
				Arguments: backend.Balance.Arguments,
			}
		}
		return farm
	}
	return nil
}

func (c *Client) parseFarmServers(farm string, p *parser.Parser) []*models.SiteFarmsItemsServersItems {
	servers := []*models.SiteFarmsItemsServersItems{}

	srvs, err := c.parseServers(farm, p)
	if err != nil {
		return servers
	}

	for _, s := range srvs {
		server := &models.SiteFarmsItemsServersItems{
			Name:           s.Name,
			Address:        s.Address,
			Port:           s.Port,
			SslCertificate: s.SslCertificate,
			Weight:         s.Weight,
		}
		if s.Ssl == "enabled" {
			server.Ssl = true
		}
		servers = append(servers, server)
	}
	return servers
}

func serializeServiceToFrontend(service *models.SiteService, name string) *models.Frontend {
	return &models.Frontend{
		Name:               name,
		Mode:               service.Mode,
		Maxconn:            service.Maxconn,
		HTTPConnectionMode: service.HTTPConnectionMode,
	}
}

func serializeFarmToBackend(farm *models.SiteFarmsItems) *models.Backend {
	backend := &models.Backend{
		Name: farm.Name,
		Mode: farm.Mode,
	}
	if farm.Forwardfor {
		backend.Forwardfor = "enabled"
	}
	if farm.Balance != nil {
		backend.Balance = &models.BackendBalance{Algorithm: farm.Balance.Algorithm, Arguments: farm.Balance.Arguments}
	}
	return backend
}

func serializeListenerToBind(listener *models.SiteServiceListenersItems) *models.Bind {
	return &models.Bind{
		Name:           listener.Name,
		Address:        listener.Address,
		Port:           listener.Port,
		Ssl:            listener.Ssl,
		SslCertificate: listener.SslCertificate,
	}
}

func serializeSiteServer(srv *models.SiteFarmsItemsServersItems) *models.Server {
	server := &models.Server{
		Address:        srv.Address,
		Name:           srv.Name,
		Port:           srv.Port,
		SslCertificate: srv.SslCertificate,
		Weight:         srv.Weight,
	}
	if srv.Ssl {
		server.Ssl = "enabled"
	}
	return server
}

// frontend backend relation helper methods
func (c *Client) removeUseFarm(frontend string, backend string, t string, p *parser.Parser) error {
	ufs, err := c.parseBackendSwitchingRules(frontend, p)
	if err != nil {
		return err
	}
	for i, uf := range ufs {
		if uf.Name == backend {
			return c.DeleteBackendSwitchingRule(int64(i), frontend, t, 0)
		}
	}
	return nil
}

func (c *Client) createBckFrontendRels(name string, b *models.SiteFarmsItems, edit bool, t string, p *parser.Parser) error {
	var res []error
	var err error
	if b.UseAs == "default" {
		if edit {
			err = c.removeUseFarm(name, b.Name, t, p)
			if err != nil {
				res = append(res, err)
			}
		}
		err = c.addDefaultBckToFrontend(name, b.Name, t, p)
		if err != nil {
			res = append(res, err)
		}
	} else {
		if b.Cond == "" || b.CondTest == "" {
			res = append(res, fmt.Errorf("Backend %s set as conditional but no conditions provided", b.Name))
		} else {
			i := int64(0)
			uf := &models.BackendSwitchingRule{
				ID:       &i,
				Name:     b.Name,
				Cond:     b.Cond,
				CondTest: b.CondTest,
			}
			err = c.CreateBackendSwitchingRule(name, uf, t, 0)
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

func (c *Client) addDefaultBckToFrontend(fName string, bName string, t string, p *parser.Parser) error {
	frontend := &models.Frontend{Name: fName}

	if err := c.parseSection(frontend, parser.Frontends, fName, p); err != nil {
		return err
	}
	frontend.DefaultBackend = bName
	if err := c.EditFrontend(fName, frontend, t, 0); err != nil {
		return err
	}
	return nil
}

func (c *Client) removeDefaultBckToFrontend(fName string, t string, p *parser.Parser) error {
	frontend := &models.Frontend{Name: fName}
	if err := c.parseSection(frontend, parser.Frontends, fName, p); err != nil {
		return err
	}
	frontend.DefaultBackend = ""
	if err := c.EditFrontend(fName, frontend, t, 0); err != nil {
		return err
	}
	return nil
}

func (c *Client) editService(name string, service *models.SiteService, t string, p *parser.Parser) error {
	frontend := &models.Frontend{Name: name}
	if err := c.parseSection(frontend, parser.Frontends, name, p); err != nil {
		return err
	}

	frontend.HTTPConnectionMode = service.HTTPConnectionMode
	frontend.Maxconn = service.Maxconn
	frontend.Mode = service.Mode

	if err := c.EditFrontend(name, frontend, t, 0); err != nil {
		return err
	}
	return nil
}

func (c *Client) editFarm(name string, farm *models.SiteFarmsItems, t string, p *parser.Parser) error {
	backend := &models.Backend{Name: name}
	if err := c.parseSection(backend, parser.Backends, name, p); err != nil {
		return err
	}

	backend.Mode = farm.Mode
	if farm.Forwardfor {
		backend.Forwardfor = "enabled"
	} else {
		backend.Forwardfor = ""
	}
	if farm.Balance != nil {
		backend.Balance = &models.BackendBalance{Algorithm: farm.Balance.Algorithm, Arguments: farm.Balance.Arguments}
	} else {
		backend.Balance = nil
	}
	if err := c.EditBackend(name, backend, t, 0); err != nil {
		return err
	}
	return nil
}

func (c *Client) editListener(name string, frontend string, listener *models.SiteServiceListenersItems, t string) error {
	bind, err := c.GetBind(name, frontend, t)
	if err != nil {
		return err
	}
	bind.Data.Address = listener.Address
	bind.Data.Port = listener.Port
	bind.Data.Ssl = listener.Ssl
	bind.Data.SslCertificate = listener.SslCertificate

	if err := c.EditBind(name, frontend, bind.Data, t, 0); err != nil {
		return err
	}
	return nil
}

func (c *Client) editSiteServer(name string, backend string, server *models.SiteFarmsItemsServersItems, t string) error {
	srv, err := c.GetServer(name, backend, t)
	if err != nil {
		return err
	}
	srv.Data.Address = server.Address
	srv.Data.Port = server.Port
	srv.Data.SslCertificate = server.SslCertificate
	srv.Data.Weight = server.Weight
	if server.Ssl {
		srv.Data.Ssl = "enabled"
	} else {
		srv.Data.Ssl = ""
	}

	if err := c.EditServer(name, backend, srv.Data, t, 0); err != nil {
		return err
	}
	return nil
}
