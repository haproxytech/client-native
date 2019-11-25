# ![HAProxy](assets/images/haproxy-weblogo-210x49.png "HAProxy")

## HAProxy Native Golang Client

**HAProxy Native Client** is a client that exposes methods for reading and changing HAProxy configuration files, and executing commands and parsing the output of the HAProxy Runtime API (via unix socket, AKA stats socket in HAProxy). It uses the [models](http://github.com/haproxytech/models) package to shift data around.

## Usage Example

```go
// Initialize HAProxy native client
confClient := &configuration.Client{}
confParams := configuration.ClientParams{
    ConfigurationFile:      "/etc/haproxy/haproxy.cfg",
    Haproxy:                "/usr/sbin/haproxy",
    UseValidation:          true,
    PersistentTransactions: true,
    TransactionDir:         "/tmp/haproxy",
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

## Generating interfaces

if you change arguments of Client functions, you need to regenerate interface

```bash
ifacemaker -f './configuration/*.go'  -p client_native -i IConfigurationClient -s Client -c "This file is generated, don't edit manually, see README.md for details." > client_interface.go
ifacemaker -f 'runtime/*.go' -p client_native -i IRuntimeClient -s Client -c "This file is generated, don't edit manually, see README.md for details." > runtimeclient_interface.go
```

You will have to make some manual editing of generated files, because this process is not perfect and generated code won't compile.

## Contributing

For commit messages and general style please follow the haproxy project's [CONTRIBUTING guide](https://github.com/haproxy/haproxy/blob/master/CONTRIBUTING) and use that where applicable.
