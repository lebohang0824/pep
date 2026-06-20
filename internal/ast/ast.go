package ast

// App is the root node
type App struct {
	Name     string
	Meta     map[string]string
	Features []*Feature
	Actions  []*Action
}

// Action represents an "action" block
type Action struct {
	Name   string
	Meta   map[string]string
	Params []*Param
}

// Function represents a "function" block inside a feature
type Function struct {
	Name   string
	Params []string
	Body   string
}

// Feature represents a "feature" block
type Feature struct {
	Name      string
	Meta      map[string]string
	Functions []*Function
	Params    []*Param
	Rules     []*Rule
	Events    []*Event
}

// Param represents a "param" block
type Param struct {
	Name     string
	Type     string
	EnumVals []string
	Required bool
	Default  interface{}
}

// Rule represents a "rule" block
type Rule struct {
	Name      string
	Condition string
}

// Event represents an event trigger inside an "events" block
type Event struct {
	Condition string
	Trigger   string
}
