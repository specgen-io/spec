# spec

Spec is a API specification format. It currently can define models and HTTP-based enpoints.

## Table of Contents
  - [Spec Format](#spec-format)
    - [Format Overview](#format-overview)
    - [Spec Structure](#spec-structure)
        - [Meta Information](#meta-information)
        - [Operations](#operations)
        - [Models](#models)
        - [Short Form](#short-form)
    - [Types](#types)
        - [Primitive Types](#primitive-types)
        - [Nullable Types](#nullable-types)
        - [Structured Types](#structured-types)
    - [Model](#model)
        - [Object Model](#object-model)
        - [Model Field](#model-field)
        - [Enum Model](#enum-model)
    - [Operation](#operation)
        - [Endpoint](#endpoint)
        - [Header Parameters](#header-parameters)
        - [Query Parameters](#query-parameters)
        - [Request Body](#request-body)
        - [Response](#response)
  - [Code Generation](#code-generation)
    - [OpenAPI](#openapi)
    - [Scala Jackson Models](#scala-jackson-models)
    - [Scala Play Service](#scala-play-service)
    - [Scala Https Client](#scala-https-client)


## Spec Format

### Format Overview

Spec format is based on YAML. Each spec file is YAML file.
In many aspects spec format resembles to [OpenAPI](https://github.com/OAI/OpenAPI-Specification/blob/master/versions/3.0.0.md). Main purpose of spec is to provide more compact (then OpenAPI) way of defining API and also limit some capabilities of OpenAPI.

### Spec Structure

Here's an example of simplest spec:
```yaml
idl_version: 0               # meta information
service_name: example
version: '1'

operations:                  # HTTP operations
  sample:                    # group of operations
    get_sample:
      endpoint: GET /sample
      response:
        ok: Sample

models:                      # models specification
  Sample:                    # name of the model
    field1: string
    field2: int
```

Spec YAML file consists of following sections:
* Meta information
* Operations
* Models

#### Meta Information

Meta information is presented in form for keys at the top level of the YAML file. In the example above `idl_version` and `service_name` are meta information fields.

Here's the list of supported meta information fields:

| Name         | Description                                                                   |
| ------------ | ----------------------------------------------------------------------------- |
| idl_version  | Version of spec format                                                        |
| service_name | Name of the service, should be in [kebab-case](http://wiki.c2.com/?KebabCase) |
| version      | Version of the specification                                                  |

#### Operations

Operations are defining HTTP requests in `operations` section of spec file. Operations are grouped, groups are used in code generation and normally should bundle together related operations. The operation group name should be in [snake_case](https://en.wikipedia.org/wiki/Snake_case).

#### Models

Models section allows to define custom user types, including dictionaries, enums, etc. Described models can be used in operations where needed just by name of the model.

#### Short form

Most entities are defined in YAML as dictinary. Some entities also support short form definition. Short form allows to define APIs in very compact way. While long form allows to specify all details of the API. Short form is basically a single field from entity full definition hence the dictionary is not needed to define the entity but only one field value is enough. Such short form field is marked accordingly in the documentation below, if it exists.

Compare long and short forms of defining field as an example:
```yaml
field: string               # <- short form

field:                      # <- long form
  type: string
  description: some field   # allows to specify more information 
```

### Types

JSON supports very limited number of types: string, number, boolean, object, array, null. Specifying JSON type is often not enough when it comes HTTP API definition and to what is allowed/prohibited as query/header/body value. For example, if the endpoint expects date and time in ISO 8601 format then in JSON it's just string, though API user supposed to pass a string only in a specific format.

The null is a problem across all JSON types. All fields can be null in JSON. Though usually API is very sensitive to null values and does not allow nulls everywhere.

Spec has it's own list of supported types to close gaps mentioned above and to provide declarative way of defining types. This section describes these types.

#### Primitive Types

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

#### Nullable Types

By default all types can't have `null` value. The `?` modifier after type defines nullable type. For example `string` can not be `null` though `string?` can have null value.

#### Structured Types

| Spec type | JSON type | Notes                                          |
| --------- | --------- | ---------------------------------------------- |
| array<_>  | array     | array items of the same type _                 |
| map<_>    | object    | object with property values of the same type _ |

Structured types are similar to generic types: `array<string>` represents array of strings.

### Model

#### Object Model

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

As table above shows object model could be defined in short form with fields only:

```yaml
Model:
  field1: string
  field2: int
```

#### Model Field

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


#### Enum Model

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

### Operation

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
  body:
    description: sample that will be created
    type: Sample
  response:
    ok: Sample
    forbidden: empty
```

Here's information about operation definition fields:

| Field       | Required             | Details                                                  |
| ----------- | -------------------- | -------------------------------------------------------- |
| description | no                   | description of the field, used for documentation         |
| endpoint    | yes                  | HTTP endpoint of operation                               |
| header      | no                   | dictionary of HTTP header parameters and types           |
| query       | no                   | dictionary of HTTP query parameters and types            |
| body        | yes for POST and PUT | HTTP body of the request                                 |
| response    | yes                  | dictionary of supported HTTP responses and response body |

#### Endpoint

Endpoint determines what url operation is processing. Endpoint has following format: `METHOD url`. Method might be one of follwing: `GET`, `POST`, `PUT`, `DELETE`. Url is just a string always starting from `/`. Url can contain parameters in following format: `{param:type}`. Here's an example for enpoint with url parameter:
```
GET /sample/{id:uuid}
```

#### Header Parameters

The `header` field of [operation](#operation) allows to describe HTTP header parameters. The `header` is a dictionary where key is the name of the header parameter and value is the definition of the parameter. Here are fields of parameter definition:

| Field       | Required | Short form | Details                                              |
| ----------- | -------- | ---------- | ---------------------------------------------------- |
| description | no       |            | description of the parameter, used for documentation |
| type        | yes      | yes        | type of the parameter                                |

Here's an example of two header parameters definition: `Authorization` (`string`) and `X-Request-Id` (`uuid`):

```yaml
header:
  Authorization:
    type: string
    description: authorization token
  X-Request-Id:
    type: uuid
    description: original request id passed
```

Header parameters names should be in Pascal-Kebab case.

Header parameters could be declared in a short form:

```yaml
header:
  Authorization: string
  X-Request-Id: string
```

#### Query Parameters

The `query` field of [operation](#operation) allows to describe HTTP query parameters. The `query` is a dictionary where key is the name of the query parameter and value is the definition of the parameter. Here are fields of parameter definition:

| Field       | Required | Short form | Details                                              |
| ----------- | -------- | ---------- | ---------------------------------------------------- |
| description | no       |            | description of the parameter, used for documentation |
| type        | yes      | yes        | type of the parameter                                |

Here's an example of two query parameters definition: `page_size` and `page_number`, both of them are `int`:

```yaml
query:
  page_size:
    type: int
    description: size of the page
  page_number:
    type: int
    description: number of requested page
```

Query parameters names should be in snake_case case.

Query parameters could be declared in a short form:

```yaml
query:
  page_size: int
  page_number: int
```

#### Request Body

The `body` field of [operation](#operation) defines body of HTTP request. Here are fields of body definition:

| Field       | Required | Short form | Details                                              |
| ----------- | -------- | ---------- | ---------------------------------------------------- |
| description | no       |            | description of body                                  |
| type        | yes      | yes        | type that represents body                            |

Here's an example of body definition represented by custom type `Sample`:

```yaml
body:
  description: sample that will be created
  type: Sample
```

Body could be declared in a short form:
```yaml
body: Sample
```

#### Response

The `response` field of [operation](#operation) defines all possible responses of HTTP request. The `response` is a dictionary where key is the name of the response and value is the definition of the response. The name of responses should be in a text form according to [RFC 7231](https://tools.ietf.org/html/rfc7231) but in snake_case. So, `OK` will be `ok`, `Unauthorized` - `unauthorized`, `Method Not Allowed` - `method_not_allowed`. 

Here are fields of response definition:

| Field       | Required | Short form | Details                                              |
| ----------- | -------- | ---------- | ---------------------------------------------------- |
| description | no       |            | description of response                              |
| type        | yes      | yes        | type that represents response body                   |

Here's an example of response definition with 2 responses: `ok` and `forbidden`, the `ok` response has a body of type `Sample`:

```yaml
response:
  ok:
    description: sample is created
    type: Sample
  forbidden:
    description: user is forbidden to create sample
    type: empty
```

## Code Generation

### OpenAPI

### Scala Jackson Models

### Scala Play Service

### Scala Https Client
