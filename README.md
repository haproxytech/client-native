# ![HAProxy](assets/images/haproxy-weblogo-210x49.png "HAProxy")

## HAProxy Native Golang Client

**HAProxy Native Client** is a client that exposes methods for reading and changing HAProxy configuration files, and executing commands and parsing the output of the HAProxy Runtime API (via unix socket, AKA stats socket in HAProxy). It uses the [models](http://github.com/haproxytech/models) package to shift data around.

## Dependencies

### Internal dependencies

The native client depends on the [config-parser](http://github.com/haproxytech/config-parser).

### External dependencies

- [purell](https://github.com/PuerkitoBio/purell)
- [urlesc](https://github.com/PuerkitoBio/urlesc)
- [govalidator](https://github.com/asaskevich/govalidator)
- [mgo](https://github.com/globalsign/mgo/bson)
- [go-openapi analysis](https://github.com/go-openapi/analysis)
- [go-openapi errors](https://github.com/go-openapi/errors)
- [go-openapi jsonpointer](https://github.com/go-openapi/jsonpointer)
- [go-openapi jsonreference](https://github.com/go-openapi/jsonreference)
- [go-openapi loads](https://github.com/go-openapi/loads)
- [go-openapi runtime](https://github.com/go-openapi/runtime)
- [go-openapi spec](https://github.com/go-openapi/spec)
- [go-openapi strfmt](https://github.com/go-openapi/strfmt)
- [go-openapi swag](https://github.com/go-openapi/swag)
- [go-openapi validate](https://github.com/go-openapi/validate)
- [easyjson](https://github.com/mailru/easyjson/buffer)
- [mapstructure](https://github.com/mitchellh/mapstructure)

## Usage Example

```
// Initialize HAProxy native client
confClient := &configuration.Client{}
confParams := configuration.ClientParams{
    ConfigurationFile: "/etc/haproxy/haproxy.cfg",
    Haproxy:           "/usr/sbin/haproxy",
    UseValidation:     true,
    TransactionDir:    "/tmp/haproxy",
}
err := confClient.Init(confParams)
if err != nil {
    fmt.Println("Error setting up configuration client, using default one")
    confClient, err = configuration.DefaultClient()
    if err != nil {
        fmt.Println("Error setting up default configuration client, exiting...")
        api.ServerShutdown()
    }
}

runtimeClient := &runtime_api.Client{}
globalConf, err := confClient.GetGlobalConfiguration("")
if err == nil {
    socketList := make([]string, 0, 1)
    runtimeAPIs := globalConf.Data.RuntimeApis

    if len(runtimeAPIs) != 0 {
        for _, r := range runtimeAPIs {
            socketList = append(socketList, *r.Address)
        }
        if err := runtimeClient.Init(socketList, "", 0); err != nil {
            fmt.Println("Error setting up runtime client, not using one")
            return nil
		}
    } else {
		fmt.Println("Runtime API not configured, not using it")
		runtimeClient = nil
	}
} else {
    fmt.Println("Cannot read runtime API configuration, not using it")
    runtimeClient = nil
}

client := &client_native.HAProxyClient{}
client.Init(confClient, runtimeClient)

bcks, err := h.Client.Configuration.GetBackends(t)
if err != nil {
    fmt.Println(err.Error())
}
//...

backendsJSON, err := bcks.MarshallBinary()

if err != nil {
    fmt.Println(err.Error())
}

fmt.Println(string(backendsJSON))
//...

```


