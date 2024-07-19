package test

import (
	_ "embed"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/haproxytech/client-native/v6/misc"
	"github.com/haproxytech/client-native/v6/models"
	"github.com/stretchr/testify/require"
)

func namedDefaultsExpectation() map[string][]*models.Defaults {
	initStructuredExpected()
	res := StructuredToDefaultsMap()
	return res
}

func TestPushDefaults(t *testing.T) {
	tOut := int64(6000)
	tOutS := int64(200)
	balanceAlgorithm := "leastconn"
	cpkaCnt := int64(10)
	cpkaTimeout := int64(10000)
	statsRealm := "Haproxy Stats"
	d := &models.Defaults{
		DefaultsBase: models.DefaultsBase{
			Clitcpka:       "disabled",
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
				Type: "send-state",
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
	m := make(map[string][]*models.Defaults)
	v, defaults, err := clientTest.GetDefaultsSections("")
	if err != nil {
		t.Error(err.Error())
	}

	for _, v := range defaults {
		d := *v
		m[d.Name] = []*models.Defaults{&d}
	}

	if len(defaults) != 3 {
		t.Errorf("%v defaults returned, expected 2", len(defaults))
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	checkNamedDefaults(t, m)
}

func checkNamedDefaults(t *testing.T, got map[string][]*models.Defaults) {
	exp := namedDefaultsExpectation()
	for k, v := range got {
		want, ok := exp[k]
		require.True(t, ok, "k=%s", k)
		require.Equal(t, len(want), len(v), "k=%s", k)
		for _, g := range v {
			for _, w := range want {
				if g.Name == w.Name {
					// This is due to the fact the unnamed defaults is modified here in TestEditCreateDeleteDefaultsSection
					// So value is not equal to what was in configuration_test.go is the test runs after the edit one.
					if g.Name != "unnamed_defaults_1" {
						require.True(t, g.DefaultsBase.Equal(w.DefaultsBase), "k=%s - diff %v", k, cmp.Diff(*g, *w))
						break
					}
				}
			}
		}
	}
}

func TestGetDefaultsSection(t *testing.T) {
	m := make(map[string][]*models.Defaults)

	v, d, err := clientTest.GetDefaultsSection("test_defaults", "")
	if err != nil {
		t.Error(err.Error())
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}
	m["test_defaults"] = append(m[""], d)

	checkNamedDefaults(t, m)
}

func TestEditCreateDeleteDefaultsSection(t *testing.T) {
	// test creating a new section
	d := &models.Defaults{
		DefaultsBase: models.DefaultsBase{
			Name:           "created",
			Clitcpka:       "disabled",
			DefaultBackend: "test2",
		}}
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
		DefaultsBase: models.DefaultsBase{
			From:           "unnamed_defaults_1",
			Name:           "created",
			Clitcpka:       "enabled",
			DefaultBackend: "test2",
		}}
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
