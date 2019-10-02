# spec

Spec is a API specification format. It currently can describe models and HTTP-based enpoints.

## Format Overview

Spec format is based on YAML. Each spec file is YAML file.
In many aspects spec format resembles to [OpenAPI](https://github.com/OAI/OpenAPI-Specification/blob/master/versions/3.0.0.md). Main purpose of spec is to provide more compact (then OpenAPI) way of describing API and also limit some capabilities of OpenAPI.

Most entities are described in YAML with dictinary. Som entities also support short form. Short form allows to describe APIs in very compact way. While long form allows to specify all details of the API. Short form is basically a single field from entity full description hence the dictionary is not needed to define the entity. You can see [field definition](Put link here) for example. Such short form field is marked accordingly in the documentation below, if it exists.

## Spec Structure

Here's an example of simplest spec:
```yaml
# meta information
idl_version: 0
service_name: example
package_name: com.company.example

operations:                  # HTTP operations
  sample:                    # group of operations
    get_sample:              # name of the operation
      endpoint: GET /sample  # endpoint
      response:              # enpoint responses
        ok: Sample           # HTTP status name: response model

models:                      # models specification
  Sample:                    # name of the model
    field: string            # field name: type
```

Spec YAML file consists of following sections:
* Meta information
* Operations
* Models

Meta information is presented in form for keys at the top level of the YAML file. In the example above `idl_version` and `service_name` are meta information fields.

Here's the list of supported meta information fields:

| Name         | Description            |
|--------------|------------------------|
| idl_version  | Version of spec format |
| service_name | Name of the service, should be in [kebab-case](http://wiki.c2.com/?KebabCase)|
| package_name | Name of the package, used for code generation |

Operations are describing HTTP-based operations in `operations` section of spec file. Operations are grouped, groups are used in code generation and normally should bundle together related operations. The operation group name should be in [snake_case](https://en.wikipedia.org/wiki/Snake_case).

Models section allows to describe custom user types, including dictionaries, enums, etc. Described models can be used in operations where needed just by name of the model.

## Types

This section describes types supported in spec format.

*Primitive types*

| Spec type         | JSON type | Notes |
|-------------------|-----------|-------|
| byte              |||
| short <br> int16  |||
| int <br> int32    |||
| long <br> int64   |||
| float             |||
| double            |||
| decimal           |||
| bool <br> boolean |||
| char              |||
| string <br> str   |||
| uuid              |||
| date              |||
| datetime          |||
| time              |||
| json              |||

*Nullable types*

By default all types can't have `null` value. The `?` modifier after type describes nullable type. For example `string` can't be `null` though `string?` can have null value.

*Structured types*

| Spec type | JSON type | Notes |
| --------- | --------- | ----- |
| array<_>  | Array     ||
| map<_>    | Object    ||

## Model

*Object model*

Here's an example of object model:
```yaml
Model:
  description: the model
  fields:
    field1: string
    field2: int
```

*Enum model*

Here's an example of enum model:
```yaml
Model:
  description: the model
  enum:
  - first
  - second
  - third
```

## Operation

Operation description here