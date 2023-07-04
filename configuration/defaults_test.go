package configuration

import (
	"fmt"
	"testing"

	"github.com/haproxytech/client-native/v5/misc"
	"github.com/haproxytech/client-native/v5/models"
)

func TestGetDefaults(t *testing.T) { //nolint:gocognit,gocyclo
	v, d, err := clientTest.GetDefaultsConfiguration("")
	if err != nil {
		t.Error(err.Error())
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	if *d.Balance.Algorithm != "roundrobin" {
		t.Errorf("Balance.Algorithm not roundrobin: %v", d.Balance.Algorithm)
	}
	if d.Mode != "http" {
		t.Errorf("Mode not http: %v", d.Mode)
	}
	if *d.Backlog != 1024 {
		t.Errorf("Backlog not 2048: %v", *d.Backlog)
	}
	if d.MonitorURI != "/monitor" {
		t.Errorf("MonitorURI not /monitor: %v", d.MonitorURI)
	}
	if d.BindProcess != "1-4" {
		t.Errorf("BindProcess not 1-4: %v", d.BindProcess)
	}
	if d.Clitcpka != "enabled" {
		t.Errorf("Clitcpka not enabled: %v", d.Clitcpka)
	}
	if d.Srvtcpka != "enabled" {
		t.Errorf("Srvtcpka not enabled: %v", d.Srvtcpka)
	}
	if d.Tcpka != "enabled" {
		t.Errorf("Tcpka not enabled: %v", d.Tcpka)
	}
	if d.Dontlognull != "enabled" {
		t.Errorf("Dontlognull not enabled: %v", d.Dontlognull)
	}
	if d.DisableH2Upgrade != "enabled" {
		t.Errorf("DisableH2Upgrade not enabled: %v", d.DisableH2Upgrade)
	}
	if d.HTTPUseHtx != "enabled" {
		t.Errorf("HTTPUseHtx not enabled: %v", d.HTTPUseHtx)
	}
	if !d.Httplog {
		t.Errorf("Httplog not enabled: %v", d.Httplog)
	}
	if d.LogHealthChecks != "enabled" {
		t.Errorf("LogHealthChecks not enabled: %v", d.LogHealthChecks)
	}
	if d.HTTPConnectionMode != "httpclose" {
		t.Errorf("HTTPConnectionMode not httpclose: %v", d.HTTPConnectionMode)
	}
	if d.HTTPRestrictReqHdrNames != "reject" {
		t.Errorf("HTTPRestrictReqHdrNames not reject: %v", d.HTTPRestrictReqHdrNames)
	}
	if d.DefaultBackend != "test" {
		t.Errorf("DefaultBackend not test: %v", d.DefaultBackend)
	}
	if *d.Maxconn != 2000 {
		t.Errorf("Maxconn not 2000: %v", *d.Maxconn)
	}
	if *d.ClientTimeout != 4000 {
		t.Errorf("ClientTimeout not 4000: %v", *d.ClientTimeout)
	}
	if *d.CheckTimeout != 2000 {
		t.Errorf("CheckTimeout not 2000: %v", *d.CheckTimeout)
	}
	if *d.ConnectTimeout != 5000 {
		t.Errorf("ConnectTimeout not 5000: %v", *d.ConnectTimeout)
	}
	if *d.QueueTimeout != 900 {
		t.Errorf("QueueTimeout not 900: %v", *d.QueueTimeout)
	}
	if *d.ServerTimeout != 2000 {
		t.Errorf("ServerTimeout not 2000: %v", *d.ServerTimeout)
	}
	if *d.ServerFinTimeout != 1000 {
		t.Errorf("ServerFinTimeout not 1000: %v", *d.ServerFinTimeout)
	}
	if *d.ClientFinTimeout != 1000 {
		t.Errorf("ServerFinTimeout not 1000: %v", *d.ClientFinTimeout)
	}
	if *d.TarpitTimeout != 2000 {
		t.Errorf("TarpitTimeout not 2000: %v", *d.TarpitTimeout)
	}
	if *d.HTTPRequestTimeout != 2000 {
		t.Errorf("HTTPRequestTimeout not 2000: %v", *d.HTTPRequestTimeout)
	}
	if *d.HTTPKeepAliveTimeout != 3000 {
		t.Errorf("HTTPKeepAliveTimeout not 3000: %v", *d.HTTPKeepAliveTimeout)
	}
	if *d.DefaultServer.Fall != 2000 {
		t.Errorf("DefaultServer.Fall not 2000: %v", *d.DefaultServer.Fall)
	}
	if *d.DefaultServer.Rise != 4000 {
		t.Errorf("DefaultServer.Rise not 4000: %v", *d.DefaultServer.Rise)
	}
	if *d.DefaultServer.Inter != 5000 {
		t.Errorf("DefaultServer.Inter not 5000: %v", *d.DefaultServer.Inter)
	}
	if *d.DefaultServer.HealthCheckPort != 8888 {
		t.Errorf("DefaultServer.HealthCheckPort not 8888: %v", *d.DefaultServer.HealthCheckPort)
	}
	if len(d.ErrorFiles) != 3 {
		t.Errorf("ErrorFiles not 3: %v", len(d.ErrorFiles))
	} else {
		for _, ef := range d.ErrorFiles {
			if ef.Code == 403 {
				if ef.File != "/test/403.html" {
					t.Errorf("File for %v not 403: %v", ef.Code, ef.File)
				}
			}
			if ef.Code == 500 {
				if ef.File != "/test/500.html" {
					t.Errorf("File for %v not 500: %v", ef.Code, ef.File)
				}
			}
			if ef.Code == 429 {
				if ef.File != "/test/429.html" {
					t.Errorf("File for %v not 429: %v", ef.Code, ef.File)
				}
			}
		}
	}
	if d.ExternalCheck != "enabled" {
		t.Errorf("ExternalCheck not enabled: %v", d.ExternalCheck)
	}
	if d.ExternalCheckPath != "/bin" {
		t.Errorf("ExternalCheckPath not /bin: %v", d.ExternalCheckPath)
	}
	if d.ExternalCheckCommand != "/bin/true" {
		t.Errorf("ExternalCheckCommand not /bin/true: %v", d.ExternalCheckCommand)
	}
	if d.AcceptInvalidHTTPRequest != "enabled" {
		t.Errorf("AcceptInvalidHTTPRequest not enabled: %v", d.AcceptInvalidHTTPRequest)
	}
	if d.AcceptInvalidHTTPResponse != "enabled" {
		t.Errorf("AcceptInvalidHTTPResponse not enabled: %v", d.AcceptInvalidHTTPResponse)
	}
	if d.H1CaseAdjustBogusClient != "enabled" {
		t.Errorf("H1CaseAdjustBogusClient not enabled: %v", d.H1CaseAdjustBogusClient)
	}
	if d.H1CaseAdjustBogusServer != "enabled" {
		t.Errorf("H1CaseAdjustBogusServer not enabled: %v", d.H1CaseAdjustBogusServer)
	}
	if d.Compression == nil {
		t.Errorf("Compression is nil")
	} else {
		if !d.Compression.Offload {
			t.Errorf("Compression.Offload wrong: %v", d.Compression.Offload)
		}
	}
	if d.ClitcpkaCnt == nil {
		t.Errorf("ClitcpkaCnt is nil")
	} else if *d.ClitcpkaCnt != 10 {
		t.Errorf("ClitcpkaCnt not 10: %v", *d.ClitcpkaCnt)
	}
	if d.ClitcpkaIdle == nil {
		t.Errorf("ClitcpkaIdle is nil")
	} else if *d.ClitcpkaIdle != 10000 {
		t.Errorf("ClitcpkaIdle not 10000: %v", *d.ClitcpkaIdle)
	}
	if d.ClitcpkaIntvl == nil {
		t.Errorf("ClitcpkaIntvl is nil")
	} else if *d.ClitcpkaIntvl != 10000 {
		t.Errorf("ClitcpkaIntvl not 10000: %v", *d.ClitcpkaIntvl)
	}
	if d.Checkcache != "disabled" {
		t.Errorf("Checkcache not disabled: %v", d.Checkcache)
	}
	if d.HTTPIgnoreProbes != "disabled" {
		t.Errorf("HTTPIgnoreProbes not disabled: %v", d.HTTPIgnoreProbes)
	}
	if d.HTTPUseProxyHeader != "disabled" {
		t.Errorf("HTTPUseProxyHeader not disabled: %v", d.HTTPUseProxyHeader)
	}
	if d.Httpslog != "disabled" {
		t.Errorf("Httpslog not disabled: %v", d.Httpslog)
	}
	if d.IndependentStreams != "disabled" {
		t.Errorf("IndependentStreams not disabled: %v", d.IndependentStreams)
	}
	if d.Nolinger != "disabled" {
		t.Errorf("Nolinger not disabled: %v", d.Nolinger)
	}
	if d.Originalto == nil {
		t.Error("Originalto is nil, expected not nil")
	} else {
		if *d.Originalto.Enabled != "enabled" {
			t.Errorf("Originalto.Enabled is not enabled: %v", *d.Originalto.Enabled)
		}
		if d.Originalto.Except != "" {
			t.Errorf("Originalto.Except is not empty: %v", d.Originalto.Except)
		}
		if d.Originalto.Header != "" {
			t.Errorf("Originalto.Header is not empty: %v", d.Originalto.Header)
		}
	}
	if d.Persist != "enabled" {
		t.Errorf("Persist not enabled: %v", d.Persist)
	}
	if d.PreferLastServer != "enabled" {
		t.Errorf("PreferLastServer not enabled: %v", d.PreferLastServer)
	}
	if d.SocketStats != "enabled" {
		t.Errorf("SocketStats not enabled: %v", d.SocketStats)
	}
	if d.TCPSmartAccept != "enabled" {
		t.Errorf("TCPSmartAccept not enabled: %v", d.TCPSmartAccept)
	}
	if d.TCPSmartConnect != "enabled" {
		t.Errorf("TCPSmartConnect not enabled: %v", d.TCPSmartConnect)
	}
	if d.Transparent != "enabled" {
		t.Errorf("Transparent not enabled: %v", d.Transparent)
	}
	if d.DontlogNormal != "enabled" {
		t.Errorf("DontlogNormal not enabled: %v", d.DontlogNormal)
	}
	if d.HTTPNoDelay != "enabled" {
		t.Errorf("HTTPNoDelay not enabled: %v", d.HTTPNoDelay)
	}
	if d.SpliceAuto != "enabled" {
		t.Errorf("SpliceAuto not enabled: %v", d.SpliceAuto)
	}
	if d.SpliceRequest != "enabled" {
		t.Errorf("SpliceRequest not enabled: %v", d.SpliceRequest)
	}
	if d.SpliceResponse != "enabled" {
		t.Errorf("SpliceResponse not enabled: %v", d.SpliceResponse)
	}
	if d.IdleCloseOnResponse != "enabled" {
		t.Errorf("IdleCloseOnResponse not enabled: %v", d.IdleCloseOnResponse)
	}

	if d.SrvtcpkaCnt == nil {
		t.Errorf("SrvtcpkaCnt is nil")
	} else if *d.SrvtcpkaCnt != 10 {
		t.Errorf("SrvtcpkaCnt not 10: %v", *d.SrvtcpkaCnt)
	}
	if d.SrvtcpkaIdle == nil {
		t.Errorf("SrvtcpkaIdle is nil")
	} else if *d.SrvtcpkaIdle != 10000 {
		t.Errorf("SrvtcpkaIdle not 10000: %v", *d.SrvtcpkaIdle)
	}
	if d.SrvtcpkaIntvl == nil {
		t.Errorf("SrvtcpkaIntvl is nil")
	} else if *d.SrvtcpkaIntvl != 10000 {
		t.Errorf("SrvtcpkaIntvl not 10000: %v", *d.SrvtcpkaIntvl)
	}
	if d.StatsOptions == nil {
		t.Errorf("StatsOptions is nil")
	}
	if d.StatsOptions.StatsShowModules != true {
		t.Error("StatsShowModules not set")
	}
	if d.StatsOptions.StatsRealm != true {
		t.Error("StatsRealm not set")
	}
	if d.StatsOptions.StatsRealmRealm == nil {
		t.Errorf("StatsRealmRealm is nil")
	} else if *d.StatsOptions.StatsRealmRealm != `HAProxy\\ Statistics` {
		t.Errorf("StatsRealmRealm not 'HAProxy Statistics': %v", *d.StatsOptions.StatsRealmRealm)
	}
	if len(d.StatsOptions.StatsAuths) != 2 {
		t.Errorf("StatsAuths expected 2 instances got: %v", len(d.StatsOptions.StatsAuths))
	}
	if d.StatsOptions.StatsAuths[0].User == nil {
		t.Errorf("StatsAuths 0 User is nil")
	} else if *d.StatsOptions.StatsAuths[0].User != "admin" {
		t.Errorf("StatsAuths 0 User not admin: %v", *d.StatsOptions.StatsAuths[0].User)
	}
	if d.StatsOptions.StatsAuths[0].Passwd == nil {
		t.Errorf("StatsAuths 0 Passwd is nil")
	} else if *d.StatsOptions.StatsAuths[0].Passwd != "AdMiN123" {
		t.Errorf("StatsAuths 0 Passwd not AdMiN123: %v", *d.StatsOptions.StatsAuths[0].Passwd)
	}
	if d.StatsOptions.StatsAuths[1].User == nil {
		t.Errorf("StatsAuths 1 User is nil")
	} else if *d.StatsOptions.StatsAuths[1].User != "admin2" {
		t.Errorf("StatsAuths 1 User not admin2: %v", *d.StatsOptions.StatsAuths[1].User)
	}
	if d.StatsOptions.StatsAuths[1].Passwd == nil {
		t.Errorf("StatsAuths 1 Passwd is nil")
	} else if *d.StatsOptions.StatsAuths[1].Passwd != "AdMiN1234" {
		t.Errorf("StatsAuths 1 Passwd not AdMiN1234: %v", *d.StatsOptions.StatsAuths[1].Passwd)
	}
	if d.LoadServerStateFromFile != "global" {
		t.Errorf("LoadServerStateFromFile not global: %v", d.LoadServerStateFromFile)
	}

	if d.EmailAlert == nil {
		t.Error("EmailAlert is nil")
	} else if *d.EmailAlert.From != "srv01@example.com" {
		t.Errorf("EmailAlert.From is not srv01@example.com: %v", *d.EmailAlert.From)
	} else if *d.EmailAlert.To != "support@example.com" {
		t.Errorf("EmailAlert.To is not support@example.com: %v", *d.EmailAlert.To)
	} else if d.EmailAlert.Level != "err" {
		t.Errorf("EmailAlert.Level is not err: %v", d.EmailAlert.Level)
	} else if d.EmailAlert.Myhostname != "srv01" {
		t.Errorf("EmailAlert.Myhostname is not srv01: %v", d.EmailAlert.Myhostname)
	} else if *d.EmailAlert.Mailers != "localmailer1" {
		t.Errorf("EmailAlert.Mailers is not localmailer1: %v", *d.EmailAlert.Mailers)
	}

	if *d.Fullconn != 10 {
		t.Errorf("Fullconn not 10: %v", *d.Fullconn)
	}
	if *d.HTTPSendNameHeader != "" {
		t.Errorf("HTTPSendNameHeader not empty: %v", *d.HTTPSendNameHeader)
	}
	if *d.MaxKeepAliveQueue != 100 {
		t.Errorf("MaxKeepAliveQueue not 100: %v", *d.MaxKeepAliveQueue)
	}
	if d.RetryOn != "503 504" {
		t.Errorf("RetryOn not 503 504: %v", d.RetryOn)
	}
	if d.PersistRule.RdpCookieName != "" {
		t.Errorf("PersistRule.RdpCookieName not empty: %v", d.PersistRule.RdpCookieName)
	}
	if *d.Source.Address != "192.168.1.200" {
		t.Errorf("Source Address not 192.168.1.200: %v", *d.Source.Address)
	}
	if d.Source.Port != 80 {
		t.Errorf("Source Port not 80: %v", d.Source.Port)
	}
	if d.Source.Usesrc != "address" {
		t.Errorf("Source Usesrc not address: %v", d.Source.Usesrc)
	}
	if d.Source.Hdr != "" {
		t.Errorf("Source Hdr not empty: %v", d.Source.Hdr)
	}
	if d.Source.Occ != "" {
		t.Errorf("Source Occ not empty: %v", d.Source.Occ)
	}
	if d.Source.AddressSecond != "192.168.1.201" {
		t.Errorf("Source Address not 192.168.1.201 %v", d.Source.AddressSecond)
	}
	if d.Source.PortSecond != 443 {
		t.Errorf("Source PortSecond not 443: %v", d.Source.PortSecond)
	}
	if d.Source.Interface != "" {
		t.Errorf("Source Interface not empty: %v", d.Source.Interface)
	}
}

func TestPushDefaults(t *testing.T) {
	tOut := int64(6000)
	tOutS := int64(200)
	balanceAlgorithm := "leastconn"
	cpkaCnt := int64(10)
	cpkaTimeout := int64(10000)
	statsRealm := "Haproxy Stats"
	d := &models.Defaults{
		Clitcpka:       "disabled",
		BindProcess:    "1-4",
		DefaultBackend: "test2",
		ErrorFiles: []*models.Errorfile{
			{
				Code: 400,
				File: "/test/400.html",
			},
			{
				Code: 403,
				File: "/test/403.html",
			},
			{
				Code: 429,
				File: "/test/429.html",
			},
			{
				Code: 500,
				File: "/test/500.html",
			},
		},
		CheckTimeout:   &tOutS,
		ConnectTimeout: &tOut,
		ServerTimeout:  &tOutS,
		QueueTimeout:   &tOutS,
		Mode:           "tcp",
		MonitorURI:     "/healthz",
		HTTPUseHtx:     "enabled",
		Balance: &models.Balance{
			Algorithm: &balanceAlgorithm,
		},
		ExternalCheck:        "",
		ExternalCheckPath:    "/bin",
		ExternalCheckCommand: "/bin/false",
		Logasap:              "disabled",
		Allbackups:           "enabled",
		HTTPCheck: &models.HTTPCheck{
			Index: misc.Int64P(0),
			Type:  "send-state",
		},
		AcceptInvalidHTTPRequest:  "disabled",
		AcceptInvalidHTTPResponse: "disabled",
		DisableH2Upgrade:          "disabled",
		LogHealthChecks:           "disabled",
		ClitcpkaCnt:               &cpkaCnt,
		ClitcpkaIdle:              &cpkaTimeout,
		ClitcpkaIntvl:             &cpkaTimeout,
		SrvtcpkaCnt:               &cpkaCnt,
		SrvtcpkaIdle:              &cpkaTimeout,
		SrvtcpkaIntvl:             &cpkaTimeout,
		Checkcache:                "enabled",
		HTTPIgnoreProbes:          "enabled",
		HTTPUseProxyHeader:        "enabled",
		Httpslog:                  "enabled",
		IndependentStreams:        "enabled",
		Nolinger:                  "enabled",
		Originalto: &models.Originalto{
			Enabled: misc.StringP("enabled"),
			Except:  "127.0.0.1",
			Header:  "X-Client-Dst",
		},
		Persist:             "disabled",
		PreferLastServer:    "disabled",
		SocketStats:         "disabled",
		TCPSmartAccept:      "disabled",
		TCPSmartConnect:     "disabled",
		Transparent:         "disabled",
		DontlogNormal:       "disabled",
		HTTPNoDelay:         "disabled",
		SpliceAuto:          "disabled",
		SpliceRequest:       "disabled",
		SpliceResponse:      "disabled",
		IdleCloseOnResponse: "disabled",
		StatsOptions: &models.StatsOptions{
			StatsShowModules: true,
			StatsRealm:       true,
			StatsRealmRealm:  &statsRealm,
			StatsAuths: []*models.StatsAuth{
				{User: misc.StringP("user1"), Passwd: misc.StringP("pwd1")},
				{User: misc.StringP("user2"), Passwd: misc.StringP("pwd2")},
			},
		},
		EmailAlert: &models.EmailAlert{
			From:       misc.StringP("srv01@example.com"),
			To:         misc.StringP("support@example.com"),
			Level:      "err",
			Myhostname: "srv01",
			Mailers:    misc.StringP("localmailer1"),
		},
		HTTPSendNameHeader: misc.StringP(""),
		Source: &models.Source{
			Address:   misc.StringP("127.0.0.1"),
			Interface: "lo",
		},
	}

	err := clientTest.PushDefaultsConfiguration(d, "", version)

	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	ver, defaults, err := clientTest.GetDefaultsConfiguration("")
	if err != nil {
		t.Error(err.Error())
	}

	var givenJSON []byte
	givenJSON, err = d.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	var ondiskJSON []byte
	ondiskJSON, err = defaults.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	if string(givenJSON) != string(ondiskJSON) {
		fmt.Printf("Created defaults: %v\n", string(ondiskJSON))
		fmt.Printf("Given defaults: %v\n", string(givenJSON))
		t.Error("Created defaults not equal to given defaults")
	}

	if ver != version {
		t.Error("Version not incremented!")
	}

	err = clientTest.PushDefaultsConfiguration(d, "", 1055)

	if err == nil {
		t.Error("Should have returned version conflict.")
	}
}

func TestGetDefaultsSections(t *testing.T) {
	v, defaults, err := clientTest.GetDefaultsSections("")
	if err != nil {
		t.Error(err.Error())
	}

	if len(defaults) != 3 {
		t.Errorf("%v defaults returned, expected 2", len(defaults))
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	unnamedSectionFound := false
	for _, d := range defaults {
		if d.Name == "test_defaults" {
			if d.From != "" {
				t.Errorf("%s: From not empty string: %s", d.Name, d.From)
			}
		}
		if d.Name == "test_defaults_2" {
			if d.From != "test_defaults" {
				t.Errorf("%s: From not test_defaults: %s", d.Name, d.From)
			}
		}
		if d.Name == "unnamed_defaults_1" {
			unnamedSectionFound = true
		}
	}

	if !unnamedSectionFound {
		t.Errorf("Unnamed section not found")
	}
}

func TestGetDefaultsSection(t *testing.T) {
	v, d, err := clientTest.GetDefaultsSection("test_defaults", "")
	if err != nil {
		t.Error(err.Error())
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}
	if d.Name != "test_defaults" {
		t.Errorf("%s: Name not test_defaults: %s", d.Name, d.Name)
	}
	if d.From != "" {
		t.Errorf("%s: From not empty string: %s", d.Name, d.From)
	}
}

func TestEditCreateDeleteDefaultsSection(t *testing.T) {
	// test creating a new section
	d := &models.Defaults{
		Name:           "created",
		Clitcpka:       "disabled",
		BindProcess:    "1-4",
		DefaultBackend: "test2",
	}
	err := clientTest.CreateDefaultsSection(d, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, defaults, err := clientTest.GetDefaultsSection("created", "")
	if err != nil {
		t.Error(err.Error())
	}

	var givenJSON []byte
	givenJSON, err = d.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	var ondiskJSON []byte
	ondiskJSON, err = defaults.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	if string(givenJSON) != string(ondiskJSON) {
		fmt.Printf("Created defaults: %v\n", string(ondiskJSON))
		fmt.Printf("Given defaults: %v\n", string(givenJSON))
		t.Error("Created defaults not equal to given defaults")
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	err = clientTest.CreateDefaultsSection(d, "", version)
	if err == nil {
		t.Error("Should throw error defaults section already exists")
		version++
	}

	d = &models.Defaults{
		From:           "unnamed_defaults_1",
		Name:           "created",
		Clitcpka:       "enabled",
		BindProcess:    "1-4",
		DefaultBackend: "test2",
	}
	err = clientTest.EditDefaultsSection("created", d, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, defaults, err = clientTest.GetDefaultsSection("created", "")
	if err != nil {
		t.Error(err.Error())
	}

	givenJSON, err = d.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	ondiskJSON, err = defaults.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	if string(givenJSON) != string(ondiskJSON) {
		fmt.Printf("Created defaults: %v\n", string(ondiskJSON))
		fmt.Printf("Given defaults: %v\n", string(givenJSON))
		t.Error("Created defaults not equal to given defaults")
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	err = clientTest.EditDefaultsSection("i_dont_exist", d, "", version)
	if err == nil {
		t.Error("editing non-existing defaults section succeeded")
	}

	// TestDeleteDefaultsSection
	err = clientTest.DeleteDefaultsSection("created", "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	if v, _ := clientTest.GetVersion(""); v != version {
		t.Error("Version not incremented")
	}

	_, _, err = clientTest.GetDefaultsSection("created", "")
	if err == nil {
		t.Error("DeleteDefaultsSection failed, defaults section created still exists")
	}

	err = clientTest.DeleteDefaultsSection("i_dont_exist", "", version)
	if err == nil {
		t.Error("Should throw error, non existent defaults section")
		version++
	}
}
