package spec

import (
	"errors"
	"gopkg.in/yaml.v3"
	"regexp"
	"strings"
)

type Endpoint struct {
	Method    string
	Url       string
	UrlParams UrlParams
}

func NewEndpoint(endpoint string) Endpoint {
	method, url, params := parseEndpoint(endpoint)
	return Endpoint{Method: method, Url: url, UrlParams: params}
}

func (value *Endpoint) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.ScalarNode {
		return errors.New("endpoint should be string")
	}
	*value = NewEndpoint(node.Value)
	return nil
}

func parseEndpoint(endpoint string) (string, string, UrlParams) {
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
		param := NewParam(paramName, ParseType(paramType), nil, nil)
		params = append(params, *param)

		cleanUrl = strings.Replace(cleanUrl, originalParamStr, UrlParamStr(paramName), 1)
	}
	return method, cleanUrl, params
}

func UrlParamStr(paramName string) string {
	return "{" + paramName + "}"
}
