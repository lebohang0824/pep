package semantic

import (
	"fmt"
	"strings"

	"pep/internal/ast"
)

func Analyze(app *ast.App) []error {
	var errs []error

	actionNames := make(map[string]bool)
	for _, a := range app.Actions {
		if actionNames[a.Name] {
			errs = append(errs, fmt.Errorf("duplicate action name '%s'", a.Name))
		}
		actionNames[a.Name] = true
	}

	featureNames := make(map[string]bool)
	for _, f := range app.Features {
		if featureNames[f.Name] {
			errs = append(errs, fmt.Errorf("duplicate feature name '%s'", f.Name))
		}
		featureNames[f.Name] = true
	}

	for _, f := range app.Features {
		paramNames := make(map[string]*ast.Param)
		for _, p := range f.Params {
			paramNames[p.Name] = p
		}

		for _, r := range f.Rules {
			paramName := extractParamFromCondition(r.Condition)
			if paramName != "" {
				if _, ok := paramNames[paramName]; !ok {
					errs = append(errs, fmt.Errorf("rule '%s' references unknown param '%s'", r.Name, paramName))
				}
			}
		}

		for _, e := range f.Events {
			if !actionNames[e.Trigger] {
				errs = append(errs, fmt.Errorf("action '%s' referenced in trigger but not defined", e.Trigger))
			}
		}

		for _, p := range f.Params {
			if p.Type == "enum" && p.Default != nil {
				defStr, ok := p.Default.(string)
				if !ok {
					errs = append(errs, fmt.Errorf("param '%s' has non-string default for enum type", p.Name))
					continue
				}
				found := false
				for _, v := range p.EnumVals {
					if v == defStr {
						found = true
						break
					}
				}
				if !found {
					errs = append(errs, fmt.Errorf("param '%s' default '%s' not in enum values %v", p.Name, defStr, p.EnumVals))
				}
			}

			if p.Required && p.Default != nil {
				errs = append(errs, fmt.Errorf("param '%s' is required but has a default value", p.Name))
			}
		}
	}

	return errs
}

func extractParamFromCondition(condition string) string {
	if strings.HasPrefix(condition, "empty(") && strings.HasSuffix(condition, ")") {
		inner := condition[6 : len(condition)-1]
		return inner
	}
	return ""
}
