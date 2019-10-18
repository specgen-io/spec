package spec

import (
	"gopkg.in/yaml.v2"
	"regexp"
	"strings"
)

type Operation struct {
	Endpoint     string       `yaml:"endpoint"`
	Description  *string      `yaml:"description"`
	Body         *Body        `yaml:"body"`
	HeaderParams HeaderParams `yaml:"header"`
	QueryParams  QueryParams  `yaml:"query"`
	Responses    Responses    `yaml:"response"`
	Method       string
	Url          string
	UrlParams    UrlParams
}

func (self *Operation) Init() {
	method, url, params := ParseEndpoint(self.Endpoint)
	self.Method = method
	self.Url = url
	self.UrlParams = params
}

type NamedOperation struct {
	Name Name
	Operation
}

type Operations []NamedOperation

type UrlParam struct {
	Name Name
	Type Type
}

func ParseEndpoint(endpoint string) (string, string, UrlParams) {
	endpointParts := strings.SplitN(endpoint, " ", 2)
	method := endpointParts[0]
	url := endpointParts[1]
	re := regexp.MustCompile(`\{[a-z][a-z0-9]*([a-z][a-z0-9]*)*:[a-z0-9_<>\\?]*\}`)
	matches := re.FindAllStringIndex(url, -1)
	params := UrlParams{}
	cleanUrl := url
	for _, match := range matches {
		start := match[0]
		end := match[1]
		originalParamStr := url[start:end]
		paramStr := originalParamStr
		paramStr = strings.Replace(paramStr, "{", "", -1)
		paramStr = strings.Replace(paramStr, "}", "", -1)
		paramParts := strings.Split(paramStr, ":")
		paramName := strings.TrimSpace(paramParts[0])
		paramType := strings.TrimSpace(paramParts[1])
		param := NamedParam{Name: Name{paramName}, Param: *NewParam(ParseType(paramType), nil)}
		params = append(params, param)

		cleanUrl = strings.Replace(cleanUrl, originalParamStr, UrlParamStr(paramName), 1)
	}
	return method, cleanUrl, params
}

func UrlParamStr(paramName string) string {
	return "{" + paramName + "}"
}

func (value *Operations) UnmarshalYAML(unmarshal func(interface{}) error) error {
	data := make(map[string]Operation)
	err := unmarshal(&data)
	if err != nil {
		return err
	}

	names := make(yaml.MapSlice, 0)
	err = unmarshal(&names)
	if err != nil {
		return err
	}

	array := make([]NamedOperation, len(names))
	for index, item := range names {
		key := item.Key.(string)
		name := Name{key}
		err := name.Check(SnakeCase)
		if err != nil {
			return err
		}
		operation := NamedOperation{Name: name, Operation: data[key]}
		operation.Init()
		array[index] = operation
	}

	*value = array
	return nil
}

type Api struct {
	Name       Name
	Operations Operations
}

type Apis []Api

func (value *Apis) UnmarshalYAML(unmarshal func(interface{}) error) error {
	data := make(map[string]Operations)
	err := unmarshal(&data)
	if err != nil {
		return err
	}

	names := make(yaml.MapSlice, 0)
	err = unmarshal(&names)
	if err != nil {
		return err
	}

	array := make([]Api, len(names))
	for index, item := range names {
		key := item.Key.(string)
		name := Name{key}
		err := name.Check(SnakeCase)
		if err != nil {
			return err
		}
		dataItem := data[key]
		array[index] = Api{Name: name, Operations: dataItem}
	}

	*value = array
	return nil
}
