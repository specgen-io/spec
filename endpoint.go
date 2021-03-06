package spec

import (
	"errors"
	"gopkg.in/vsapronov/yaml.v3"
	"regexp"
	"strings"
)

type Endpoint struct {
	Method    string
	Url       string
	UrlParams UrlParams
}

func (value *Endpoint) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.ScalarNode {
		return yamlError(node, "operation endpoint should be string")
	}
	endpoint, err := parseEndpoint(node.Value, node)
	if err != nil {
		return yamlError(node, err.Error())
	}
	*value = *endpoint
	return nil
}

func parseEndpoint(endpoint string, node *yaml.Node) (*Endpoint, error) {
	spaces_count := strings.Count(endpoint, " ")
	if spaces_count != 1 {
		return nil, errors.New("endpoint should be in format 'METHOD url'")
	}
	endpointParts := strings.SplitN(endpoint, " ", 2)
	method := endpointParts[0]
	err := HttpMethod.Check(method)
	if err != nil {
		return nil, err
	}
	url := endpointParts[1]
	re := regexp.MustCompile(`\{[a-z][a-z0-9_]*:[a-z0-9_<>\\?]*\}`)
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

		typ, err := parseType(paramType)
		if err != nil {
			return nil, err
		}

		param := &NamedParam{
			Name: Name{Source: paramName, Location: node},
			DefinitionDefault: DefinitionDefault{
				Type:     Type{*typ, node},
				Location: node,
			},
		}

		params = append(params, *param)

		cleanUrl = strings.Replace(cleanUrl, originalParamStr, UrlParamStr(paramName), 1)
	}
	return &Endpoint{Method: method, Url: cleanUrl, UrlParams: params}, nil
}

func UrlParamStr(paramName string) string {
	return "{" + paramName + "}"
}
