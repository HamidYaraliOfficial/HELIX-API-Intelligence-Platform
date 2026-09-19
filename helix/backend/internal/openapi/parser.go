// Package openapi implements the HELIX Contract Engine for OpenAPI documents.
// It normalizes both Swagger 2.0 and OpenAPI 3.x (in either JSON or YAML)
// into a single internal Contract representation so the rest of the
// platform (Catalog, Discovery, Diff Engine, ...) only ever deals with one
// shape of data.
package openapi

import (
	"fmt"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// Contract is HELIX's normalized, version-agnostic representation of an
// API specification.
type Contract struct {
	Title       string
	Description string
	Version     string
	BaseURL     string
	AuthType    string
	Endpoints   []ContractEndpoint
}

// ContractEndpoint is a normalized operation (method + path).
type ContractEndpoint struct {
	Path        string
	Method      string
	OperationID string
	Summary     string
	Description string
	AuthType    string
	Deprecated  bool
	Parameters  []ContractParameter
	Responses   []ContractResponse
}

// ContractParameter is a normalized request parameter.
type ContractParameter struct {
	Name        string
	In          string
	Type        string
	Required    bool
	Description string
}

// ContractResponse is a normalized response definition.
type ContractResponse struct {
	StatusCode  string
	Description string
	ContentType string
}

var httpMethods = []string{"get", "put", "post", "delete", "options", "head", "patch", "trace"}

// Parse accepts raw bytes of an OpenAPI/Swagger document in either JSON or
// YAML (JSON is valid YAML, so a single code path handles both) and returns
// a normalized Contract, or an error describing why parsing failed.
func Parse(raw []byte) (*Contract, error) {
	var doc map[string]interface{}
	if err := yaml.Unmarshal(raw, &doc); err != nil {
		return nil, fmt.Errorf("document is not valid JSON or YAML: %w", err)
	}
	if doc == nil {
		return nil, fmt.Errorf("document is empty")
	}

	if _, isV2 := doc["swagger"]; isV2 {
		return parseV2(doc)
	}
	if _, isV3 := doc["openapi"]; isV3 {
		return parseV3(doc)
	}

	return nil, fmt.Errorf("unrecognized document: missing 'swagger' or 'openapi' root key")
}

// ---- Swagger 2.0 -----------------------------------------------------

func parseV2(doc map[string]interface{}) (*Contract, error) {
	info := asMap(doc["info"])
	c := &Contract{
		Title:       asString(info["title"]),
		Description: asString(info["description"]),
		Version:     asString(info["version"]),
	}

	scheme := "https"
	if schemes := asStringSlice(doc["schemes"]); len(schemes) > 0 {
		scheme = schemes[0]
	}
	host := asString(doc["host"])
	basePath := asString(doc["basePath"])
	if host != "" {
		c.BaseURL = strings.TrimRight(fmt.Sprintf("%s://%s%s", scheme, host, basePath), "/")
	}

	c.AuthType = detectAuthTypeV2(asMap(doc["securityDefinitions"]))

	paths := asMap(doc["paths"])
	for _, path := range sortedKeys(paths) {
		operations := asMap(paths[path])
		for _, method := range httpMethods {
			opRaw, ok := operations[method]
			if !ok {
				continue
			}
			op := asMap(opRaw)
			ep := ContractEndpoint{
				Path:        path,
				Method:      strings.ToUpper(method),
				OperationID: asString(op["operationId"]),
				Summary:     asString(op["summary"]),
				Description: asString(op["description"]),
				Deprecated:  asBool(op["deprecated"]),
				AuthType:    c.AuthType,
			}
			for _, p := range asSlice(op["parameters"]) {
				pm := asMap(p)
				ep.Parameters = append(ep.Parameters, ContractParameter{
					Name:        asString(pm["name"]),
					In:          asString(pm["in"]),
					Type:        asString(pm["type"]),
					Required:    asBool(pm["required"]),
					Description: asString(pm["description"]),
				})
			}
			responses := asMap(op["responses"])
			for _, code := range sortedKeys(responses) {
				rm := asMap(responses[code])
				ep.Responses = append(ep.Responses, ContractResponse{
					StatusCode:  code,
					Description: asString(rm["description"]),
					ContentType: "application/json",
				})
			}
			c.Endpoints = append(c.Endpoints, ep)
		}
	}

	return c, nil
}

func detectAuthTypeV2(defs map[string]interface{}) string {
	for _, raw := range defs {
		d := asMap(raw)
		switch asString(d["type"]) {
		case "oauth2":
			return "oauth2"
		case "basic":
			return "basic"
		case "apiKey":
			return "apiKey"
		}
	}
	return "none"
}

// ---- OpenAPI 3.x -------------------------------------------------------

func parseV3(doc map[string]interface{}) (*Contract, error) {
	info := asMap(doc["info"])
	c := &Contract{
		Title:       asString(info["title"]),
		Description: asString(info["description"]),
		Version:     asString(info["version"]),
	}

	if servers := asSlice(doc["servers"]); len(servers) > 0 {
		c.BaseURL = asString(asMap(servers[0])["url"])
	}

	components := asMap(doc["components"])
	c.AuthType = detectAuthTypeV3(asMap(components["securitySchemes"]))

	paths := asMap(doc["paths"])
	for _, path := range sortedKeys(paths) {
		operations := asMap(paths[path])
		for _, method := range httpMethods {
			opRaw, ok := operations[method]
			if !ok {
				continue
			}
			op := asMap(opRaw)
			ep := ContractEndpoint{
				Path:        path,
				Method:      strings.ToUpper(method),
				OperationID: asString(op["operationId"]),
				Summary:     asString(op["summary"]),
				Description: asString(op["description"]),
				Deprecated:  asBool(op["deprecated"]),
				AuthType:    c.AuthType,
			}
			for _, p := range asSlice(op["parameters"]) {
				pm := asMap(p)
				schema := asMap(pm["schema"])
				ep.Parameters = append(ep.Parameters, ContractParameter{
					Name:        asString(pm["name"]),
					In:          asString(pm["in"]),
					Type:        asString(schema["type"]),
					Required:    asBool(pm["required"]),
					Description: asString(pm["description"]),
				})
			}
			if reqBody := asMap(op["requestBody"]); len(reqBody) > 0 {
				content := asMap(reqBody["content"])
				for contentType := range content {
					ep.Parameters = append(ep.Parameters, ContractParameter{
						Name:        "body",
						In:          "body",
						Type:        contentType,
						Required:    asBool(reqBody["required"]),
						Description: asString(reqBody["description"]),
					})
				}
			}
			responses := asMap(op["responses"])
			for _, code := range sortedKeys(responses) {
				rm := asMap(responses[code])
				contentType := ""
				content := asMap(rm["content"])
				for ct := range content {
					contentType = ct
					break
				}
				ep.Responses = append(ep.Responses, ContractResponse{
					StatusCode:  code,
					Description: asString(rm["description"]),
					ContentType: contentType,
				})
			}
			c.Endpoints = append(c.Endpoints, ep)
		}
	}

	return c, nil
}

func detectAuthTypeV3(schemes map[string]interface{}) string {
	for _, raw := range schemes {
		s := asMap(raw)
		switch asString(s["type"]) {
		case "oauth2":
			return "oauth2"
		case "openIdConnect":
			return "oidc"
		case "http":
			if asString(s["scheme"]) == "bearer" {
				return "bearer"
			}
			return "basic"
		case "apiKey":
			return "apiKey"
		case "mutualTLS":
			return "mtls"
		}
	}
	return "none"
}

// ---- generic YAML/JSON helpers -----------------------------------------

func asMap(v interface{}) map[string]interface{} {
	if v == nil {
		return map[string]interface{}{}
	}
	switch m := v.(type) {
	case map[string]interface{}:
		return m
	case map[interface{}]interface{}:
		out := make(map[string]interface{}, len(m))
		for k, val := range m {
			out[fmt.Sprintf("%v", k)] = val
		}
		return out
	default:
		return map[string]interface{}{}
	}
}

func asSlice(v interface{}) []interface{} {
	if s, ok := v.([]interface{}); ok {
		return s
	}
	return nil
}

func asStringSlice(v interface{}) []string {
	var out []string
	for _, item := range asSlice(v) {
		out = append(out, asString(item))
	}
	return out
}

func asString(v interface{}) string {
	if v == nil {
		return ""
	}
	return fmt.Sprintf("%v", v)
}

func asBool(v interface{}) bool {
	b, _ := v.(bool)
	return b
}

func sortedKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
