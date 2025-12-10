package test

import (
	_ "embed"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/haproxytech/client-native/v6/misc"
	"github.com/haproxytech/client-native/v6/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPushStructuredDefaults(t *testing.T) {
	clientTest, filename, err := getTestClient()
	require.NoError(t, err)
	defer os.Remove(filename)
	version := int64(1)

	tOut := int64(6000)
	tOutS := int64(200)
	balanceAlgorithm := "leastconn"
	cpkaCnt := int64(10)
	cpkaTimeout := int64(10000)
	statsRealm := "Haproxy Stats"
	d := &models.Defaults{
		DefaultsBase: models.DefaultsBase{
			Name:           "unnamed_defaults1",
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
			ExternalCheck:                        "",
			ExternalCheckPath:                    "/bin",
			ExternalCheckCommand:                 "/bin/false",
			Logasap:                              "disabled",
			Allbackups:                           "enabled",
			AcceptInvalidHTTPRequest:             "disabled",
			AcceptInvalidHTTPResponse:            "disabled",
			AcceptUnsafeViolationsInHTTPRequest:  "disabled",
			AcceptUnsafeViolationsInHTTPResponse: "disabled",
			DisableH2Upgrade:                     "disabled",
			LogHealthChecks:                      "disabled",
			ClitcpkaCnt:                          &cpkaCnt,
			ClitcpkaIdle:                         &cpkaTimeout,
			ClitcpkaIntvl:                        &cpkaTimeout,
			SrvtcpkaCnt:                          &cpkaCnt,
			SrvtcpkaIdle:                         &cpkaTimeout,
			SrvtcpkaIntvl:                        &cpkaTimeout,
			Checkcache:                           "enabled",
			HTTPIgnoreProbes:                     "enabled",
			HTTPUseProxyHeader:                   "enabled",
			Httpslog:                             "enabled",
			IndependentStreams:                   "enabled",
			Nolinger:                             "enabled",
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
		TCPResponseRuleList: models.TCPResponseRules{
			&models.TCPResponseRule{
				Type:     "content",
				Action:   "accept",
				Cond:     "if",
				CondTest: "FALSE",
			},
		},
		HTTPRequestRuleList: models.HTTPRequestRules{
			&models.HTTPRequestRule{
				Type:     "allow",
				Cond:     "if",
				CondTest: "TRUE",
			},
		},
		HTTPCheckList: models.HTTPChecks{
			&models.HTTPCheck{
				Type: "send-state",
			},
		},
		HTTPErrorRuleList: models.HTTPErrorRules{
			&models.HTTPErrorRule{
				ReturnContent:       "/test/503",
				ReturnContentFormat: "file",
				ReturnContentType:   misc.Ptr("\"application/json\""),
				Status:              503,
				Type:                "status",
			},
		},
		TCPCheckRuleList: models.TCPChecks{
			&models.TCPCheck{
				Action: "send",
				Data:   "GET\\ /\\ HTTP/2.0\\r\\n",
			},
		},
		LogTargetList: models.LogTargets{
			&models.LogTarget{
				Address:  "192.169.0.1",
				Facility: "mail",
				Global:   true,
			},
		},
	}

	err = clientTest.PushStructuredDefaultsConfiguration(d, "", version)

	require.NoError(t, err)
	version++

	ver, defaults, err := clientTest.GetStructuredDefaultsConfiguration("")
	require.NoError(t, err)

	// assert.EqualValues(t, d, defaults)
	require.True(t, defaults.Equal(*d), "k=%s - diff %v", "defaults", cmp.Diff(*defaults, *d))

	if ver != version {
		t.Error("Version not incremented!")
	}

	err = clientTest.PushStructuredDefaultsConfiguration(d, "", 1055)

	if err == nil {
		t.Error("Should have returned version conflict.")
	}
}

func TestGetStructuredDefaultsSections(t *testing.T) {
	clientTest, filename, err := getTestClient()
	require.NoError(t, err)
	defer os.Remove(filename)
	version := int64(1)

	m := make(map[string]*models.Defaults)
	v, defaults, err := clientTest.GetStructuredDefaultsSections("")
	require.NoError(t, err)

	for _, v := range defaults {
		m[v.Name] = v
	}

	require.Equal(t, 3, len(defaults), "%v defaults returned, expected 2", len(defaults))
	require.Equal(t, version, v, "Version %v returned, expected %v", v, version)

	checkNamedStructuredDefaults(t, m)
}

func checkNamedStructuredDefaults(t *testing.T, got map[string]*models.Defaults) {
	exp := namedDefaultsExpectation()
	for k, v := range got {
		w, ok := exp[k]
		require.True(t, ok, "k=%s", k)
		if v.Name != "unnamed_defaults_1" {
			require.True(t, v.Equal(*w), "k=%s - diff %v", k, cmp.Diff(*v, *w))
			break
		}
	}
}

func TestGetStructuredDefaultsSection(t *testing.T) {
	clientTest, filename, err := getTestClient()
	require.NoError(t, err)
	defer os.Remove(filename)
	version := int64(1)

	m := make(map[string]*models.Defaults)

	v, d, err := clientTest.GetStructuredDefaultsSection("test_defaults", "")
	require.NoError(t, err)
	require.Equal(t, version, v, "Version %v returned, expected %v", v, version)

	m["test_defaults"] = d
	checkNamedStructuredDefaults(t, m)
}

func TestEditCreateDeleteStructuredDefaultsSection(t *testing.T) {
	clientTest, filename, err := getTestClient()
	require.NoError(t, err)
	defer os.Remove(filename)
	version := int64(1)

	// test creating a new section
	d := &models.Defaults{
		DefaultsBase: models.DefaultsBase{
			Name:           "created",
			Clitcpka:       "disabled",
			DefaultBackend: "test2",
		},
	}
	err = clientTest.CreateStructuredDefaultsSection(d, "", version)
	require.NoError(t, err)
	version++

	v, defaults, err := clientTest.GetStructuredDefaultsSection("created", "")
	require.NoError(t, err)
	assert.EqualValues(t, d, defaults)
	require.Equal(t, version, v, "Version %v returned, expected %v", v, version)

	err = clientTest.CreateStructuredDefaultsSection(d, "", version)
	require.Error(t, err, "Should throw error defaults section already exists")

	d = &models.Defaults{
		DefaultsBase: models.DefaultsBase{
			From:           "unnamed_defaults_1",
			Name:           "created",
			Clitcpka:       "enabled",
			DefaultBackend: "test2",
		},
	}
	err = clientTest.EditStructuredDefaultsSection("created", d, "", version)
	require.NoError(t, err)
	version++

	v, defaults, err = clientTest.GetStructuredDefaultsSection("created", "")
	require.NoError(t, err)
	assert.EqualValues(t, d, defaults)
	require.Equal(t, version, v, "Version %v returned, expected %v", v, version)

	err = clientTest.EditStructuredDefaultsSection("i_dont_exist", d, "", version)
	require.Error(t, err, "editing non-existing defaults section succeeded")

	// TestDeleteDefaultsSection
	err = clientTest.DeleteDefaultsSection("created", "", version)
	require.NoError(t, err)
	version++

	v, _ = clientTest.GetVersion("")
	require.Equal(t, version, v, "Version %v returned, expected %v", v, version)

	_, _, err = clientTest.GetStructuredDefaultsSection("created", "")
	require.Error(t, err, "DeleteDefaultsSection failed, defaults section created still exists")

	err = clientTest.DeleteDefaultsSection("i_dont_exist", "", version)
	require.Error(t, err, "Should throw error, non existent defaults section")
}
