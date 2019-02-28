package configuration

import (
	"strings"

	strfmt "github.com/go-openapi/strfmt"
	parser "github.com/haproxytech/config-parser"
	"github.com/haproxytech/config-parser/errors"
	"github.com/haproxytech/config-parser/params"
	"github.com/haproxytech/config-parser/types"
	"github.com/haproxytech/models"
)

// GetGlobalConfiguration returns a struct with configuration version and a
// struct representing Global configuration
func (c *Client) GetGlobalConfiguration(transactionID string) (*models.GetGlobalOKBody, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return nil, err
	}

	_, err = p.Get(parser.Global, parser.GlobalSectionName, "daemon")
	d := "enabled"
	if err == errors.FetchError {
		d = "disabled"
	}

	data, err := p.Get(parser.Global, parser.GlobalSectionName, "maxconn")
	mConn := int64(0)
	if err == nil {
		maxConn := data.(*types.Int64C)
		mConn = maxConn.Value
	}

	data, err = p.Get(parser.Global, parser.GlobalSectionName, "nbproc")
	nbproc := int64(0)
	if err == nil {
		nbProcParser := data.(*types.Int64C)
		nbproc = nbProcParser.Value
	}

	data, err = p.Get(parser.Global, parser.GlobalSectionName, "stats socket")
	rAPI := ""
	rLevel := ""
	rMode := ""
	if err == nil {
		sockets := data.([]types.Socket)
		if len(sockets) > 0 {
			s := sockets[0]
			rAPI = s.Path
			for _, p := range s.Params {
				d := p.(*params.BindOptionValue)
				if d.Name == "level" {
					rLevel = d.Value
				} else if d.Name == "mode" {
					rMode = d.Value
				}
			}
		}
	}

	data, err = p.Get(parser.Global, parser.GlobalSectionName, "ssl-default-bind-ciphers")
	sslCiphers := ""
	if err == nil {
		sslCiphersParser := data.(*types.StringC)
		sslCiphers = sslCiphersParser.Value
	}

	data, err = p.Get(parser.Global, parser.GlobalSectionName, "ssl-default-bind-options")
	sslOptions := ""
	if err == nil {
		sslOptionsParser := data.(*types.StringSliceC)
		sslOptions = strings.Join(sslOptionsParser.Value, " ")
	}

	data, err = p.Get(parser.Global, parser.GlobalSectionName, "tune.ssl.default-dh-param")
	dhParam := int64(0)
	if err == nil {
		dhParamsParser := data.(*types.Int64C)
		dhParam = dhParamsParser.Value
	}

	g := &models.Global{
		Daemon:                d,
		Maxconn:               mConn,
		Nbproc:                nbproc,
		RuntimeAPI:            rAPI,
		RuntimeAPILevel:       rLevel,
		RuntimeAPIMode:        rMode,
		SslDefaultBindCiphers: sslCiphers,
		SslDefaultBindOptions: sslOptions,
		TuneSslDefaultDhParam: dhParam,
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return nil, err
	}

	return &models.GetGlobalOKBody{Version: v, Data: g}, nil
}

// PushGlobalConfiguration pushes a Global config struct to global
// config gile
func (c *Client) PushGlobalConfiguration(data *models.Global, transactionID string, version int64) error {
	if c.UseValidation {
		validationErr := data.Validate(strfmt.Default)
		if validationErr != nil {
			return NewConfError(ErrValidationError, validationErr.Error())
		}
	}

	p, t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	pDaemon := &types.Enabled{}
	if data.Daemon != "enabled" {
		pDaemon = nil
	}
	p.Set(parser.Global, parser.GlobalSectionName, "daemon", pDaemon)

	pMaxConn := &types.Int64C{
		Value: data.Maxconn,
	}
	if data.Maxconn == 0 {
		pMaxConn = nil
	}
	p.Set(parser.Global, parser.GlobalSectionName, "maxconn", pMaxConn)

	pNbProc := &types.Int64C{
		Value: data.Nbproc,
	}
	if data.Nbproc == 0 {
		pNbProc = nil
	}
	p.Set(parser.Global, parser.GlobalSectionName, "nbproc", pNbProc)

	ondisk, err := p.Get(parser.Global, parser.GlobalSectionName, "stats socket")
	if err != nil {
		if data.RuntimeAPI != "" {
			pStatsSocket := types.Socket{
				Path: data.RuntimeAPI,
				Params: []params.BindOption{
					&params.BindOptionValue{Name: "level", Value: data.RuntimeAPILevel},
					&params.BindOptionValue{Name: "mode", Value: data.RuntimeAPIMode},
				},
			}
			p.Set(parser.Global, parser.GlobalSectionName, "stats socket", pStatsSocket)
		} else {
			p.Set(parser.Global, parser.GlobalSectionName, "stats socket", nil)
		}
	} else {
		sockets := ondisk.([]types.Socket)
		if data.RuntimeAPI != "" {
			pStatsSocket := types.Socket{
				Path: data.RuntimeAPI,
				Params: []params.BindOption{
					&params.BindOptionValue{Name: "level", Value: data.RuntimeAPILevel},
					&params.BindOptionValue{Name: "mode", Value: data.RuntimeAPIMode},
				},
			}
			(sockets)[0] = pStatsSocket
		}
		p.Set(parser.Global, parser.GlobalSectionName, "stats socket", sockets)
	}

	pSSLCiphers := &types.StringC{
		Value: data.SslDefaultBindCiphers,
	}
	if data.SslDefaultBindCiphers == "" {
		pSSLCiphers = nil
	}
	p.Set(parser.Global, parser.GlobalSectionName, "ssl-default-bind-ciphers", pSSLCiphers)

	pSSLOptions := &types.StringSliceC{
		Value: strings.Split(data.SslDefaultBindOptions, " "),
	}
	if data.SslDefaultBindCiphers == "" {
		pSSLOptions = nil
	}
	p.Set(parser.Global, parser.GlobalSectionName, "ssl-default-bind-options", pSSLOptions)

	pDhParams := &types.Int64C{
		Value: data.TuneSslDefaultDhParam,
	}
	if data.TuneSslDefaultDhParam == 0 {
		pDhParams = nil
	}
	p.Set(parser.Global, parser.GlobalSectionName, "tune.ssl.default-dh-param", pDhParams)

	if err := c.saveData(p, t, transactionID == ""); err != nil {
		return err
	}
	return nil
}
