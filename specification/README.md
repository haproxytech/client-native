# ![HAProxy](../assets/images/haproxy-weblogo-210x49.png "HAProxy")

## HAProxy Data Plane API Specification

This is the [OpenAPI 2.0 (fka Swagger)](https://github.com/OAI/OpenAPI-Specification/blob/master/versions/2.0.md) specification for the [HAProxy Data Plane API project](https://github.com/haproxytech/dataplaneapi)

## Contributing

When contributing, change files located in paths/ and models/ directories and the haproxy-spec.yaml, and then build the resulting one-file spec `build/haproxy_spec.yaml`.

```bash
go run cmd/specification/*.go -file specification/haproxy-spec.yaml > specification/build/haproxy_spec.yaml
```

this is already integrated into [make models](../Makefile) command

For commit messages and general style please follow the haproxy project's [CONTRIBUTING guide](https://github.com/haproxy/haproxy/blob/master/CONTRIBUTING) and use that where applicable.


##  Note: some source yaml files in paths/ and models/ are go templates.*
The templates are applied in the
```bash
go run cmd/specification/*.go -file specification/haproxy-spec.yaml > specification/build/haproxy_spec.yaml
```
so when it's applied when you run the [make models](../Makefile) command

The files to update for the go templates are in `configuration/parents`:
- [constants.go](../configuration/parents/constants.go)
  - <xx>ChildType: list of objects that have a parent like 'server'
  - CnParentType are the Config Parser names of the parents.
- [parents.go](../configuration/parents/parents.go)
  - Contains for each child type the list of parents

For example in [haproxy-spec.yaml](../specification/haproxy-spec.yaml):
```
  {{ range parents "server" -}}
  /services/haproxy/configuration/{{ .PathParentType }}/{parent_name}/servers:
    $ref: "paths/configuration/server.yaml#/servers"
  {{ end -}}
```

`parents` is
```
 func Parents(childType string) []Parent {
```
defined in `parents.go` and is used in [server.yaml](../specification/paths/configuration/server.yaml)
```
servers:
  get:
    summary: Return an array of servers
    description: Returns an array of all servers that are configured in specified backend.
    operationId: getAllServer{{ .ParentType }}
...
```
