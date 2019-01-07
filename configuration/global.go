package configuration

import (
	"fmt"
	"strings"

	strfmt "github.com/go-openapi/strfmt"
	bindoptions "github.com/haproxytech/config-parser/bind-options"
	"github.com/haproxytech/config-parser/parsers"
	"github.com/haproxytech/config-parser/parsers/global"
	"github.com/haproxytech/config-parser/parsers/simple"
	"github.com/haproxytech/config-parser/parsers/stats"
	"github.com/haproxytech/models"
)

// GetGlobalConfiguration returns a struct with configuration version and a
// struct representing Global configuration
func (c *Client) GetGlobalConfiguration() (*models.GetGlobalOKBody, error) {
	err := c.GlobalParser.LoadData(c.GlobalConfigurationFile)
	if err != nil {
		return nil, err
	}

	data, err := c.GlobalParser.GetGlobalAttr("daemon")
	d := ""
	if err == nil {
		daemon := data.(*parsers.Daemon)
		if daemon.Valid() {
			d = "enabled"
		} else {
			d = "disabled"
		}
	}

	data, err = c.GlobalParser.GetGlobalAttr("maxconn")
	mConn := int64(0)
	if err == nil {
		maxConn := data.(*parsers.MaxConn)
		mConn = maxConn.Value
	}

	data, err = c.GlobalParser.GetGlobalAttr("nbproc")
	nbproc := int64(0)
	if err == nil {
		nbProcParser := data.(*global.NbProc)
		nbproc = nbProcParser.Value
	}

	data, err = c.GlobalParser.GetGlobalAttr("stats socket")
	rAPI := ""
	rLevel := ""
	rMode := ""
	if err == nil {
		sockets := data.(*stats.SocketLines)
		if len(sockets.SocketLines) > 0 {
			s := sockets.SocketLines[0]
			rAPI = s.Path
			for _, p := range s.Params {
				d := p.(*bindoptions.BindOptionValue)
				if d.Name == "level" {
					rLevel = d.Value
				} else if d.Name == "mode" {
					rMode = d.Value
				}
			}
		}
	}

	data, err = c.GlobalParser.GetGlobalAttr("ssl-default-bind-ciphers")
	sslCiphers := ""
	if err == nil {
		sslCiphersParser := data.(*simple.SimpleString)
		sslCiphers = sslCiphersParser.Value
	}

	data, err = c.GlobalParser.GetGlobalAttr("ssl-default-bind-options")
	sslOptions := ""
	if err == nil {
		sslOptionsParser := data.(*simple.SimpleStringMultiple)
		sslOptions = strings.Join(sslOptionsParser.Value, " ")
	}

	data, err = c.GlobalParser.GetGlobalAttr("tune.ssl.default-dh-param")
	dhParam := int64(0)
	if err == nil {
		dhParamsParser := data.(*simple.SimpleNumber)
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

	v, err := c.GetGlobalVersion()
	if err != nil {
		return nil, err
	}

	return &models.GetGlobalOKBody{Version: v, Data: g}, nil
}

// PushGlobalConfiguration pushes a Global config struct to global
// config gile
func (c *Client) PushGlobalConfiguration(data *models.Global, version int64) error {
	ondiskV, _ := c.GetGlobalVersion()
	if ondiskV != version {
		return NewConfError(ErrVersionMismatch, fmt.Sprintf("Version in configuration file is %v, given version is %v", ondiskV, version))
	}

	validationErr := data.Validate(strfmt.Default)
	if validationErr != nil {
		return NewConfError(ErrValidationError, validationErr.Error())
	}

	err := c.GlobalParser.LoadData(c.GlobalConfigurationFile)
	if err != nil {
		return nil
	}

	pDaemon := &parsers.Daemon{}
	pDaemon.Init()
	if data.Daemon == "enabled" {
		pDaemon.Enabled = true
	}
	c.GlobalParser.Global.Set(pDaemon)

	pMaxConn := &parsers.MaxConn{}
	pMaxConn.Init()
	if data.Maxconn > 0 {
		pMaxConn.Value = data.Maxconn
	}
	c.GlobalParser.Global.Set(pMaxConn)

	pNbProc := &global.NbProc{}
	pNbProc.Init()
	if data.Nbproc > 0 {
		pNbProc.Value = data.Nbproc
		pNbProc.Enabled = true
	}
	c.GlobalParser.Global.Set(pNbProc)

	ondisk, _ := c.GlobalParser.GetGlobalAttr("stats socket")
	if ondisk == nil {
		if data.RuntimeAPI != "" {
			pStatsSocket := &stats.SocketLines{}
			pStatsSocket.Init()

			s := &stats.Socket{}
			s.Path = data.RuntimeAPI
			s.Params = []bindoptions.BindOption{}
			s.Params = append(s.Params, &bindoptions.BindOptionValue{Name: "level", Value: data.RuntimeAPILevel})
			s.Params = append(s.Params, &bindoptions.BindOptionValue{Name: "mode", Value: data.RuntimeAPIMode})

			pStatsSocket.SocketLines = append(pStatsSocket.SocketLines, s)
			c.GlobalParser.Global.Set(pStatsSocket)
		}
	} else {
		sockets := ondisk.(*stats.SocketLines)
		if data.RuntimeAPI == "" {
			sockets.Init()
		} else {
			s := sockets.SocketLines[0]
			s.Path = data.RuntimeAPI
			for _, p := range s.Params {
				d := p.(*bindoptions.BindOptionValue)
				if d.Name == "level" {
					d.Value = data.RuntimeAPILevel
				} else if d.Name == "mode" {
					d.Value = data.RuntimeAPIMode
				}
			}
		}
		c.GlobalParser.Global.Set(sockets)
	}

	pSSLCiphers := &simple.SimpleString{}
	pSSLCiphers.Init()
	if data.SslDefaultBindCiphers != "" {
		pSSLCiphers.Name = "ssl-default-bind-ciphers"
		pSSLCiphers.SearchName = pSSLCiphers.Name
		pSSLCiphers.Value = data.SslDefaultBindCiphers
		pSSLCiphers.Enabled = true
	}
	c.GlobalParser.Global.Set(pSSLCiphers)

	pSSLOptions := &simple.SimpleStringMultiple{}
	pSSLOptions.Init()
	if data.SslDefaultBindOptions != "" {
		pSSLOptions.Name = "ssl-default-bind-options"
		pSSLOptions.SearchName = pSSLOptions.Name
		pSSLOptions.Value = strings.Split(data.SslDefaultBindOptions, " ")
		pSSLOptions.Enabled = true
	}
	c.GlobalParser.Global.Set(pSSLOptions)

	pDhParams := &simple.SimpleNumber{}
	pDhParams.Init()
	if data.TuneSslDefaultDhParam > 0 {
		pDhParams.Name = "tune.ssl.default-dh-param"
		pDhParams.SearchName = pDhParams.Name
		pDhParams.Value = data.TuneSslDefaultDhParam
		pDhParams.Enabled = true
	}
	c.GlobalParser.Global.Set(pDhParams)

	fmt.Println(c.GlobalParser.String())
	err = c.GlobalParser.Save(c.GlobalConfigurationFile)
	if err != nil {
		return err
	}
	err = c.incrementGlobalVersion()
	if err != nil {
		return err
	}
	return nil
}
