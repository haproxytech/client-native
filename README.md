# ![HAProxy](../assets/images/haproxy-weblogo-210x49.png "HAProxy")

## HAProxy Native Golang Client

**HAProxy Native Client** is a client that exposes methods for reading and changing HAProxy configuration files, and executing commands and parsing the output of the HAProxy Runtime API (via unix socket, AKA stats socket in HAProxy). It uses the [models](http://github.com/haproxytech/models) package to shift data around.

## Dependencies

### Internal dependencies

The native client currently depends on the [config-parser](http://github.com/haproxytech/config-parser) and the [lbctl](http://github.com/HAPEE/lbctl) parser. The tendency is to drop the **lbctl** in favour of the native **config-parser** as the later project is developed.

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
confClient, err = configuration.DefaultClient()
if err != nil {
    fmt.Println("Error setting up default configuration client, exiting...")
}

var nbproc int64
data, err := confClient.GlobalParser.GetGlobalAttr("nbproc")
if err != nil {
    nbproc = int64(1)
} else {
    d := data.(*simple.SimpleNumber)
    nbproc = d.Value
}

statsSocket := ""
data, err = confClient.GlobalParser.GetGlobalAttr("stats socket")
if err == nil {
    statsSockets := data.(*stats.SocketLines)
    statsSocket = statsSockets.SocketLines[0].Path
} else {
    fmt.Println("Error getting stats socket")
    fmt.Println(err.Error())
}

if statsSocket == "" {
    fmt.Println("Stats socket not configured, no runtime client initiated")
    runtimeClient = nil
} else {
    socketList := make([]string, 0, 1)
    if nbproc > 1 {
        for i := int64(0); i < nbproc; i++ {
            socketList = append(socketList, fmt.Sprintf("%v.%v", statsSocket, i))
        }
    } else {
        socketList = append(socketList, statsSocket)
    }
    err := runtimeClient.Init(socketList)
    if err != nil {
        fmt.Println("Error setting up runtime client, not using one")
        runtimeClient = nil
    }
}

client := &HAProxyClient{}
client.Init(confClient, runtimeClient)

rawConfig, err := client.Configuration.GetRawConfiguration()

stats, err := client.Runtime.GetStats()
```


