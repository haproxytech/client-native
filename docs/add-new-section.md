# Adding support for a new configuration Section

In this tutorial, we will add support for the `mailers` section to
`client-native`, as defined in [the HAProxy documentation][mailersdoc].

## Prerequisite

A parser for the section must be written in [config-parser], or at least,
the section name must be declared so that it can be parsed automatically.
This is outside the scope of this tutorial though.

## 1. Write the OpenAPI data specification

The first step is to create new data types for our new section, and the data
it contains.

Create a new YAML file [specification/models/configuration/mailers.yaml
](../specification/models/configuration/mailers.yaml), and define the
properties and contents of the section.

>We are using [go-swagger] which is limited to OpenAPI 2.0 for now.

In this case, 2 data types have been defined: `mailers_section` and
`mailer_entry`. We often use the suffixes *_section* and *_entry* to
differentiate 2 types that share the same name and to prevent confusion.

The `mailers` section has an optional property `timeout` which contains
a [duration]. Those are always defined in OpenAPI as integers. And since
0 (zero) is a valid value, add the attribute `x-nullable: true`.

## 2. Reference your new types in the main specification

Edit [haproxy-spec.yaml] under the `definitions` section. Create a definition
for each type and also definitions for arrays of each type. For example:

```yaml
mailers_section:
  $ref: "models/configuration/mailers.yaml#/mailers_section"
mailers_sections:
  title: Mailers Sections
  description: HAProxy mailers_section array
  type: array
  items:
    $ref: "#/definitions/mailers_section"
```

Same for `mailer_entry` and `mailer_entries`.

Then, in the `tags` section, reference your types like so:

```yaml
- name: Mailers
- name: MailerEntry
```

## 3. Validate your changes and generate models

Run `make models` to validate the syntax of the YAML files
and generate Go model files for each type.

The new files will be placed in the [models](../models/) folder.

## 4. Write the OpenAPI path specification

Create a new YAML file [specification/paths/configuration/mailers.yaml
](../specification/paths/configuration/mailers.yaml) and define the REST API
used to get/create/update/delete your new types.

You must define different path definitions to access either 1 object or a list
of objects. For example: `mailers_section_one` and `mailers_sections`.

Use the methods `get`, `post`, `put` and `delete`, to respectively
retrieve, create, update and delete entries.

Note that reference to data types are using the form

```yaml
$ref: "#/definitions/mailers_section"
```

which means we are using references from [haproxy-spec.yaml]
instead of using the specific file for this data type.

## 5. Reference your paths in the main specification

Open [haproxy-spec.yaml] again and define a path for each REST API defined
previously. Add them to the `paths` section. For instance:

```yaml
/services/haproxy/configuration/mailers_section:
  $ref: "paths/configuration/mailers.yaml#/mailers_section"
/services/haproxy/configuration/mailers_section/{name}:
  $ref: "paths/configuration/mailers.yaml#/mailers_section_one"
```

And same for `mailer_entry`.

## 6. Validate your changes

Once again, run `make models` to validate your YAML files
and update the models.

## 7. Create the configuration handlers

You must now write the code responsible for modifying HAProxy's configuration
when the API is called. For each data type, this is done in the `configuration`
folder. In this example, the following files must be created:

 - [mailers_section.go](../configuration/mailers_section.go)
 - [mailers_section_test.go](../configuration/mailers_section_test.go)
 - [mailer_entry.go](../configuration/mailer_entry.go)
 - [mailer_entry_test.go](../configuration/mailer_entry_test.go)

For each data type, create an interface with methods for each API endpoint.
You send or receive data using the types generated previously,
like `models.MailersSections`.

To read or modify HAProxy's configuration, you must use one of the Parser
types, like `parser.Mailers`. These parsers come from [config-parser].

For each type, it is required to write a *Parse* and *Serialize* function.
in this case:

```go
func ParseMailersSection(p parser.Parser, ms *models.MailersSection) error
func SerializeMailersSection(p parser.Parser, data *models.MailersSection) error
func ParseMailerEntry(m types.Mailer) *models.MailerEntry
func SerializeMailerEntry(me models.MailerEntry) types.Mailer
```

## 8. Reference your new interfaces

Open [configuration/interface.go](../configuration/interface.go) and add the
interfaces you created in the previous step to the `Configuration` interface.

Be thorough with the testing.

  *Good luck!*

[mailersdoc]: https://docs.haproxy.org/2.6/configuration.html#3.6
[config-parser]: https://github.com/haproxytech/config-parser
[duration]: https://docs.haproxy.org/2.6/configuration.html#2.5
[go-swagger]: https://goswagger.io/
[haproxy-spec.yaml]: ../specification/haproxy-spec.yaml
