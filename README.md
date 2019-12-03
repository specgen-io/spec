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
        - [Array and Dictionary](#array-and-dictionary)
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
    - [Scala Play Service](#scala-play-service)
    - [Scala Sttp Client](#scala-sttp-client)
    - [Scala Models](#scala-models)

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
| :----------- | :---------------------------------------------------------------------------- |
| idl_version  | Version of spec format                                                        |
| service_name | Name of the service, should be in [kebab-case](http://wiki.c2.com/?KebabCase) |
| version      | Version of the specification                                                  |

#### Operations

Operations are defining HTTP requests in `operations` section of spec file. Operations are grouped, groups are used in code generation and normally should bundle together related operations. The operation group name should be in [snake_case](https://en.wikipedia.org/wiki/Snake_case).

#### Models

Models section allows to define custom user types, including dictionaries, enums, etc. Described models can be used in operations where needed just by name of the model.

#### Short form

Most entities are defined in YAML as dictionary. Some entities also support short form definition. Short form allows to define APIs in very compact way. While long form allows to specify all details of the API entities. Short form is better readable by humans.

Compare long and short forms of defining field as an example:
```yaml
# short form
field: string = the value     # some field

# long form
field:
  type: string
  default: the value
  description: some field 
```

Entities supporting short form are documented accordingly. Descriptions of most entities could be provided in a comment on the same line where key for entity is defined. In the example above `some field` is a description for the `field` specified via comment on the same line as the field key `field`.

### Types

JSON supports very limited number of types: string, number, boolean, object, array, null. Specifying JSON type is often not enough when it comes HTTP API definition and to what is allowed/prohibited as query/header/body value. For example, if the endpoint expects date and time in ISO 8601 format then in JSON it's just a string, though API user supposed to pass a string only in a specific format.

The null is a problem across all JSON types. All fields can be null in JSON. Though usually API is very sensitive to null values and does not allow nulls everywhere.

Spec has it's own list of supported types to close gaps mentioned above and to provide declarative way of defining types. This section describes these types.

#### Primitive Types

| Spec type         | JSON type | Notes                                             |
| :---------------- | :-------- | :------------------------------------------------ |
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
| json              | object    | any unstructured JSON                             |
| empty             | N/A       | represents nothing, similar to unit is some langs |

#### Nullable Types

By default all types can't have `null` value. The `?` modifier after type defines nullable type. For example `string` can not be `null` though `string?` can have null value.

#### Array and Dictionary

Following modifiers allow to specify data structures:

| Modifier   | JSON type | Notes                                       |
| :--------- | :-------- | :-------------------------------------------|
| the_type[] | array     | array of items of the_type                  |
| the_type{} | object    | dictionary with property values of the_type |

For example `string[]` represents array of strings. The type `int{}` represents JSON object where all properties values are integers.

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

| Field       | Required | Details                                          |
| :---------- | :------- | :----------------------------------------------- |
| fields      | yes      | dictionary of fields, keys are names of fields   |
| description | no       | description of the model, used for documentation |

As table above shows object model could be defined in short form with fields only:

```yaml
Model:  # the model
  field1: string
  field2: int
```

#### Model Field

Here's an example of field definition:
```yaml
field1:
  type: string
  default: some default value
  description: some field
```

Here's information about field definition fields:

| Field       | Required | Details                                          |
| :---------- | :------- | :----------------------------------------------- |
| type        | yes      | type of the field                                |
| default     | no       | default value for the field                      |
| description | no       | description of the field, used for documentation |

When field is defined in a short form with type only, the default value might be defined through `=` separator. Here's an example of short form equivalent to long field definition from above

```yaml
field1: string = some default value  # some field
```

#### Enum Model

Enum is represented in JSON as a string with limited set of possible values.

Here's an example of enum model:
```yaml
Model:
  description: the model
  enum:
    first:
      description: First option
    second:
      description: Second option
    third:
      description: Third option
```

Here's information about enum definition fields:

| Field       | Required | Details                                                                              |
| :---------- | :------- | :----------------------------------------------------------------------------------- |
| enum        | yes      | either list of strings or dictionary with keys and values with enum item information |
| description | no       | description of the model, used for documentation                                     |

Here's an equivalent example for short form of enum definition:

```yaml
Model:  # the model
  enum:
  - first   # First option
  - second  # Second option
  - third   # Third option
```

Enum item definition supports following fields.

| Field       | Required | Details                                                                              |
| :---------- | :------- | :----------------------------------------------------------------------------------- |
| description | no       | description of the item, used for documentation                                      |

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
| endpoint    | yes                  | HTTP endpoint of operation                               |
| body        | yes for POST and PUT | HTTP body of the request                                 |
| response    | yes                  | dictionary of supported HTTP responses and response body |
| header      | no                   | dictionary of HTTP header parameters and types           |
| query       | no                   | dictionary of HTTP query parameters and types            |
| description | no                   | description of the field, used for documentation         |

#### Endpoint

Endpoint determines what url operation is processing. Endpoint has following format: `METHOD url`. Method might be one of follwing: `GET`, `POST`, `PUT`, `DELETE`. Url is just a string always starting from `/`. Url can contain parameters in following format: `{param:type}`. Here's an example for enpoint with url parameter:
```
GET /sample/{id:uuid}
```

#### Header Parameters

The `header` field of [operation](#operation) allows to describe HTTP header parameters. The `header` is a dictionary where key is the name of the header parameter and value is the definition of the parameter. Here are fields of parameter definition:

| Field       | Required | Details                                              |
| :---------- | :------- | :--------------------------------------------------- |
| type        | yes      | type of the parameter                                |
| default     | no       | default value for the parameter                      |
| description | no       | description of the parameter, used for documentation |

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

Header parameters names should be in Pascal-Kebab case because this is how they should appear in real HTTP request.

Header parameters could be declared in a short form:

```yaml
header:
  Authorization: string  # authorization token
  X-Request-Id: string   # original request id passed
```

When short form is used the default value for parameter could be specified through `=` separator:

```yaml
header:
  X-Request-Id: string = some default id   # original request id passed
```


#### Query Parameters

The `query` field of [operation](#operation) allows to describe HTTP query parameters. The `query` is a dictionary where key is the name of the query parameter and value is the definition of the parameter. Here are fields of parameter definition:

| Field       | Required | Details                                              |
| :---------- | :------- | :--------------------------------------------------- |
| type        | yes      | type of the parameter                                |
| default     | no       | default value for the parameter                      |
| description | no       | description of the parameter, used for documentation |

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
  page_size: int     # size of the page
  page_number: int   # number of requested page
```

When short form is used the default value for parameter could be specified through `=` separator:

```yaml
query:
  page_size: int = 100   # size of the page
  page_number: int = 0   # number of requested page
```

#### Request Body

The `body` field of [operation](#operation) defines body of HTTP request. Here are fields of body definition:

| Field       | Required | Details                                              |
| :---------- | :------- | :--------------------------------------------------- |
| type        | yes      | type that represents body                            |
| description | no       | description of body                                  |

Here's an example of body definition represented by custom type `Sample`:

```yaml
body:
  type: Sample
  description: sample that will be created
```

Body could be declared in a short form:
```yaml
body: Sample  # sample that will be created
```

#### Response

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

The `response` field of [operation](#operation) defines all possible responses of HTTP request. The `response` is a dictionary where key is the name of the response and value is the definition of the response. The name of responses should be in a text form according to [RFC 7231](https://tools.ietf.org/html/rfc7231) but in snake_case. So, `OK` will be `ok`, `Unauthorized` - `unauthorized`, `Method Not Allowed` - `method_not_allowed`. 

Here are fields of response definition:

| Field       | Required | Details                                              |
| :---------- | :------- | :--------------------------------------------------- |
| type        | yes      | type that represents response body                   |
| description | no       | description of response                              |

Response could be declared in a short form, here's short form equivalent of responses above:

```yaml
response:
  ok: Sample        # sample is created
  forbidden: empty  # user is forbidden to create sample
```

## Code Generation

Code generation from spec is implemented in [specgen](https://github.com/ModaOperandi/specgen) command line tool. Sections below describes how to get generated code in various scenarios.

### OpenAPI

Spec might be translated into [OpenAPI](https://swagger.io/docs/specification/about/) (aka Swagger) almost without losses where OpenAPI supports features that are used in spec.
Following command line generates OpenAPI from spec:
```shell script
$ specgen openapi --spec-file spec.yaml --out-file swagger.yaml
```

### Scala

Spec scala code generation is packaged in form of SBT plugins. Add `sbt-spec` plugin library into your project using following line in your `plugins.sbt`:
```
addSbtPlugin("spec" % "sbt-spec" % "<version>")
```
Instead of `<version>` use latest available from [Bintray](https://bintray.com/sbt/sbt-plugin-releases?filterByPkgName=sbt-spec).
The `sbt-spec` library contains several plugins each serving specific purpose and wrapping specific `specgen` tool command.

Some of generated code relies on base libraries that are hosted in JCenter repository. Therefore includion of that repository is required in your `build.sbt`:
```
useJCenter := true
```

Under the hood all spec SBT plugins are running [specgen](https://github.com/ModaOperandi/specgen) command line tool.

#### Scala Play Service

Spec code generation has full support for Scala Play applications. Models, controllers, services traits and services scaffolded implementations are generated by `SpecPlay` plugin.

To enable server side code generation from spec enable `SpecPlay` plugin in your `build.sbt`:
```scala
enablePlugins(SpecPlay)
```
Following dependencies are defined in `SpecPlay` and should be added into `libraryDependencies`:
```
libraryDependencies ++= Seq(
  specPlay,
  swaggerUI
)
```

Here are `SpecPlay` settings that allow customization of code generation:

| Setting          | Default                              | Description                                                                                |
| :--------------- | :----------------------------------- | :----------------------------------------------------------------------------------------- |
| specFile         | `spec.yaml`                          | Path to spec file relative to project folder                                               |
| specSwagger      | `public`                             | Path to folder for generated swagger file                                                  |
| specGeneratePath | `target/scala/src_managed/main/spec` | Path to generated code is placed                                                           |
| specServicesPath | `app/services`                       | Path to scaffolded services files; these services are scaffolded only if they do not exist |
| specRoutesPath   | `conf`                               | Path to folder for generated Play routes file                                              |

Default settings are aligned with the standard [Play application layout](https://www.playframework.com/documentation/2.7.x/Anatomy).

The generated swagger (OpenAPI) spec is hosted at `/docs` route of the Play application.

#### Scala Sttp Client

Spec code generation can generate HTTP client in Scala based on spec. The generated client is using [sttp](https://github.com/softwaremill/sttp) client under the hood with provided ability to plugin backend of your taste.

To enable client code generation from spec enable `SpecSttp` plugin in your `build.sbt`:
```scala
enablePlugins(SpecSttp)
```
Following dependencies are defined in `SpecSttp` and should be added into `libraryDependencies`:
```
libraryDependencies ++= Seq(
  specSttp
)
```

Here are `SpecSttp` settings that allow customization of code generation:

| Setting          | Default                              | Description                                  |
| :--------------- | :----------------------------------- | :------------------------------------------- |
| specFile         | `spec.yaml`                          | Path to spec file relative to project folder |
| specGeneratePath | `target/scala/src_managed/main/spec` | Path to generated code is placed             |

#### Scala Models

Spec code generation can generate only (de)serializable to JSON models without any application specifics.

To enable models generation from spec enable `SpecModels` plugin in your `build.sbt`:
```scala
enablePlugins(SpecModels)
```
Following dependencies are defined in `SpecModels` and should be added into `libraryDependencies`:
```
libraryDependencies ++= Seq(
  specCirce
)
```

Here are `SpecModels` settings that allow customization of code generation:

| Setting          | Default                              | Description                                  |
| :--------------- | :----------------------------------- | :------------------------------------------- |
| specFile         | `spec.yaml`                          | Path to spec file relative to project folder |
| specGeneratePath | `target/scala/src_managed/main/spec` | Path to generated code is placed             |
