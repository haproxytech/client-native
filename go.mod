module github.com/haproxytech/client-native

go 1.12

require (
	github.com/go-openapi/errors v0.19.0
	github.com/go-openapi/strfmt v0.19.0
	github.com/google/uuid v1.1.1
	github.com/haproxytech/config-parser v1.0.3
	github.com/haproxytech/models v1.1.1-0.20190808132148-c7bd6054bdcc
	github.com/mitchellh/mapstructure v1.1.2
	github.com/pkg/errors v0.8.1
	github.com/stretchr/testify v1.3.0
)

replace github.com/haproxytech/models => ../models
