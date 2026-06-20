package generator

import (
	"encoding/json"

	"pep/internal/ast"
)

type SpecOutput struct {
	AppName  string            `json:"appName"`
	Meta     map[string]string `json:"meta"`
	Actions  []ActionOutput    `json:"actions"`
	Features []FeatureOutput   `json:"features"`
}

type ActionOutput struct {
	Name   string        `json:"name"`
	Meta   map[string]string `json:"meta"`
	Params []ParamOutput `json:"params"`
}

type FeatureOutput struct {
	Name      string            `json:"name"`
	Meta      map[string]string `json:"meta"`
	Functions []FunctionOutput  `json:"functions"`
	Params    []ParamOutput     `json:"params"`
	Rules     []RuleOutput      `json:"rules"`
	Events    []EventOutput     `json:"events"`
}

type ParamOutput struct {
	Name     string      `json:"name"`
	Type     string      `json:"type"`
	Enum     []string    `json:"enum,omitempty"`
	Required bool        `json:"required"`
	Default  interface{} `json:"default,omitempty"`
}

type RuleOutput struct {
	Name      string `json:"name"`
	Condition string `json:"condition"`
}

type FunctionOutput struct {
	Name   string   `json:"name"`
	Params []string `json:"params"`
	Body   string   `json:"body"`
}

type EventOutput struct {
	Condition string `json:"condition"`
	Trigger   string `json:"trigger"`
}

func Generate(app *ast.App, pretty bool) ([]byte, error) {
	output := SpecOutput{
		AppName:  app.Name,
		Meta:     app.Meta,
		Actions:  make([]ActionOutput, 0),
		Features: make([]FeatureOutput, 0),
	}

	for _, a := range app.Actions {
		ao := ActionOutput{
			Name:   a.Name,
			Meta:   a.Meta,
			Params: convertParams(a.Params),
		}
		if ao.Meta == nil {
			ao.Meta = make(map[string]string)
		}
		if ao.Params == nil {
			ao.Params = make([]ParamOutput, 0)
		}
		output.Actions = append(output.Actions, ao)
	}

	for _, f := range app.Features {
		fo := FeatureOutput{
			Name:      f.Name,
			Meta:      f.Meta,
			Functions: convertFunctions(f.Functions),
			Params:    convertParams(f.Params),
			Rules:     convertRules(f.Rules),
			Events:    convertEvents(f.Events),
		}
		if fo.Meta == nil {
			fo.Meta = make(map[string]string)
		}
		if fo.Functions == nil {
			fo.Functions = make([]FunctionOutput, 0)
		}
		if fo.Params == nil {
			fo.Params = make([]ParamOutput, 0)
		}
		if fo.Rules == nil {
			fo.Rules = make([]RuleOutput, 0)
		}
		if fo.Events == nil {
			fo.Events = make([]EventOutput, 0)
		}
		output.Features = append(output.Features, fo)
	}

	if pretty {
		return json.MarshalIndent(output, "", "  ")
	}
	return json.Marshal(output)
}

func convertParams(params []*ast.Param) []ParamOutput {
	var out []ParamOutput
	for _, p := range params {
		po := ParamOutput{
			Name:     p.Name,
			Type:     p.Type,
			Required: p.Required,
		}
		if p.Type == "enum" {
			po.Enum = p.EnumVals
		}
		if p.Default != nil {
			po.Default = p.Default
		}
		out = append(out, po)
	}
	return out
}

func convertRules(rules []*ast.Rule) []RuleOutput {
	var out []RuleOutput
	for _, r := range rules {
		out = append(out, RuleOutput{
			Name:      r.Name,
			Condition: r.Condition,
		})
	}
	return out
}

func convertFunctions(functions []*ast.Function) []FunctionOutput {
	var out []FunctionOutput
	for _, f := range functions {
		out = append(out, FunctionOutput{
			Name:   f.Name,
			Params: f.Params,
			Body:   f.Body,
		})
	}
	return out
}

func convertEvents(events []*ast.Event) []EventOutput {
	var out []EventOutput
	for _, e := range events {
		out = append(out, EventOutput{
			Condition: e.Condition,
			Trigger:   e.Trigger,
		})
	}
	return out
}
