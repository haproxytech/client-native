# ![HAProxy](../../assets/images/haproxy-weblogo-210x49.png "HAProxy")

## Generator for: `models/server_params_prepare_for_runtime.go`

This genetator generates a file that will contains the needed functions to prepare a ServerParams to use in combination with `add server` command on the runtime socket.

This file will contain 3 functions:

```
func (p *ServerParams) prepareForRuntimeDoNotSendDisabledFields()
func (p *ServerParams) prepareForRuntimeDoNotSendEnabledFields()
func (p *ServerParams) PrepareFieldsForRuntimeAddServer()
```

They are used for example in Ingress Controller the following way:
```
params.PrepareFieldsForRuntimeAddServer()
serverParams := configuration.SerializeServerParams(defaultServer.ServerParams)
res := cp_params.ServerOptionsString(serverParams)
```

### func (p *ServerParams) prepareForRuntimeDoNotSendDisabledFields()
For example for `Check` that has the values [enabled disabled].

- if the value is `enabled` we must send `add server check`
- if the value is `disabled` we must not send `add server no-check` as `no-check` is not allowed on a dynamic server

`no-check` is the default value.

The purpose is to set `Check` to "" when the value was `disabled` so the commands sent are:
   - `add server check` if value is `enabled`
   - `add server` if value is `disabled`


### func (p *ServerParams) prepareForRuntimeDoNotSendEnabledFields()
It's just the opposite.

For example for `NoSslv3`

- if `enabled` we must send `no-sslv3`
- if `disabled` we must not sent an option


### func (p *ServerParams) PrepareFieldsForRuntimeAddServer()`
is just calling both `PrepareForRuntimeDoNotSendDisabledFields` and  `PrepareForRuntimeDoNotSendEnabledFields`


## WHAT TO DO

Just fill in `server_params_prepare_for_runtime.go` the map:
-  `ServerParamsPrepareForRuntimeMap`

for each field that has an  `// Enum: [enabled disabled]"`
with the correct function to use

## Ensure that the map is always filled in a field is added into ServerParams: CI check

The generator checks all fields in ServerParams that have   `// Enum: [enabled disabled]"` have an entry in the `ServerParamsPrepareForRuntimeMap`.

If a new field is added and not declared in `ServerParamsPrepareForRuntimeMap`, the generator fails with an error and this will make the CI fail.
