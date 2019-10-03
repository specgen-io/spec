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
| ------------ | ---------------------- |
| idl_version  | Version of spec format |
| service_name | Name of the service, should be in [kebab-case](http://wiki.c2.com/?KebabCase)|
| package_name | Name of the package, used for code generation |

Operations are describing HTTP-based operations in `operations` section of spec file. Operations are grouped, groups are used in code generation and normally should bundle together related operations. The operation group name should be in [snake_case](https://en.wikipedia.org/wiki/Snake_case).

Models section allows to describe custom user types, including dictionaries, enums, etc. Described models can be used in operations where needed just by name of the model.

## Types

JSON supports very limited number of types: string, number, boolean, object, array, null. Specifying JSON type is often not enough when it comes to describing HTTP API and what is allowed/prohibited as query/header/body value. For example, if the endpoint expects date and time in ISO 8601 format then in JSON it's just string, though API user supposed to pass a string only in a specific format.

The null is a problem across all JSON types. All fields can be null in JSON. Though usually API is very sensitive to null values and does not allow nulls everywhere.

Spec has it's own list of supported types to close gaps mentioned above and to provide declarative way of describing types. This section describes these types.

**Primitive types**

| Spec type         | JSON type | Notes                                             |
| ----------------- | --------- | ------------------------------------------------- |
| byte              | number    | -128 to 127                                       |
| short <br> int16  | number    | -32768 to 32767                                   |
| int <br> int32    | number    | -2147483648 to 2147483647                         |
| long <br> int64   | number    | -9223372036854775808 to 9223372036854775807       |
| float             | number    | 32 bit IEEE 754 single-precision float            |
| double            | number    | 64 bit IEEE 754 double-precision float            |
| decimal           | number    | arbitrary-precision signed decimal                |
| bool <br> boolean | boolean   |                                                   |
| char              | string    | single symbol string                              |
| string <br> str   | string    |                                                   |
| uuid              | string    | lower case hex symbols with hyphens as 8-4-4-4-12 |
| date              | string    | ISO 8601 yyyy-mm-dd                               |
| datetime          | string    | ISO 8601 yyyy-mm-ddThh:mm:ss.ffffff               |
| time              | string    | ISO 8601 hh:mm:ss.ffffff                          |
| json              | object    | any JSON                                          |
| empty             | N/A       | represents nothing, similar to unit is some langs |

**Nullable types**

By default all types can't have `null` value. The `?` modifier after type describes nullable type. For example `string` can't be `null` though `string?` can have null value.

**Structured types**

| Spec type | JSON type | Notes                                          |
| --------- | --------- | ---------------------------------------------- |
| array<_>  | array     | array items of the same type _                 |
| map<_>    | object    | object with property values of the same type _ |

Structured types are similar to generic types: `array<string>` represents array of strings.

## Model

**Object model**

Here's an example of object model definition:
```yaml
Model:
  description: the model
  fields:
    field1: string
    field2: int
```
Here's information about object model fields:

| Field       | Required | Short form | Details                                          |
| ----------- | -------- | ---------- | ------------------------------------------------ |
| description | no       |            | description of the model, used for documentation |
| fields      | yes      | yes        | dictionary of fields, keys are names of fields   |

As table above shows object model could be described in short form with fields only:

```yaml
Model:
  field1: string
  field2: int
```

**Model field**

Here's an example of field definition:
```yaml
field1:
  description: some field
  type: string
```

Here's information about field definition fields:

| Field       | Required | Short form | Details                                          |
| ----------- | -------- | ---------- | ------------------------------------------------ |
| description | no       |            | description of the field, used for documentation |
| type        | yes      | yes        | type of the field                                |


**Enum model**

Enum is represented in JSON as a string with limited set of possible values.

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

Here's an example of operation:
```yaml
create_sample:
  description: creates sample
  endpoint: POST /sample
  header:
    Authorization: string
  query:
    sample_id: uuid
    user_id: int?
  body: Sample
  response:
    ok: Sample
    forbidden: empty
```

Here's information about operation definition fields:

| Field       | Required             | Details                                                  |
| ----------- | -------------------- | -------------------------------------------------------- |
| description | no                   | description of the field, used for documentation         |
| endpoint    | yes                  | HTTP endpoint of operation                               |
| header      | no                   | dictionary of HTTP headers parameters and types          |
| query       | no                   | dictionary of HTTP query parameters and types            |
| body        | yes for POST and PUT | HTTP body of the request                                 |
| response    | yes                  | dictionary of supported HTTP responses and response body |

**Endpoint**

Endpoint determines what url operation is processing. Endpoint has following format: `METHOD url`. Method might be one of follwing: `GET`, `POST`, `PUT`, `DELETE`. Url is just a string always starting from `/`. Url can contain parameters in following format: `{param:type}`. Here's an example for enpoint with url parameter:
```
GET /sample/{id:uuid}
```