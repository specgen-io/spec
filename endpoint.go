package spec

import (
	"gopkg.in/yaml.v3"
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
	method, url, params, err := parseEndpoint(node.Value, node)
	if err != nil {
		return yamlError(node, err.Error())
	}
	*value = Endpoint{Method: method, Url: url, UrlParams: params}
	return nil
}

func parseEndpoint(endpoint string, node *yaml.Node) (string, string, UrlParams, error) {
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

		typ, err := parseType(paramType)
		if err != nil {
			return "", "", nil, err
		}

		param := &NamedParam{
			Name: Name{Source: paramName, Location: node},
			DefinitionDefault: DefinitionDefault{
				definitionDefault: definitionDefault{
					Type: TypeLocated{
						Definition: *typ,
						Location:   node,
					},
				},
				Location: node,
			},
		}

		params = append(params, *param)

		cleanUrl = strings.Replace(cleanUrl, originalParamStr, UrlParamStr(paramName), 1)
	}
	return method, cleanUrl, params, nil
}

func UrlParamStr(paramName string) string {
	return "{" + paramName + "}"
}
