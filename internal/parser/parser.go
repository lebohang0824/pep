package parser

import (
	"fmt"
	"strconv"

	"pep/internal/ast"
	"pep/internal/lexer"
)

type Parser struct {
	l         *lexer.Lexer
	curToken  lexer.Token
	peekToken lexer.Token
	errors    []string
}

func NewParser(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) curTokenIs(tok lexer.TokenType) bool {
	return p.curToken.Type == tok
}

func (p *Parser) peekTokenIs(tok lexer.TokenType) bool {
	return p.peekToken.Type == tok
}

func (p *Parser) expectPeek(tok lexer.TokenType) bool {
	if p.peekTokenIs(tok) {
		p.nextToken()
		return true
	}
	p.appendError("expected %s but got %s (line=%d,col=%d)", tok, p.peekToken.Type, p.peekToken.Line, p.peekToken.Column)
	return false
}

func (p *Parser) expectCurrent(tok lexer.TokenType) bool {
	if p.curTokenIs(tok) {
		p.nextToken()
		return true
	}
	p.appendError("expected %s but got %s (line=%d,col=%d)", tok, p.curToken.Type, p.curToken.Line, p.curToken.Column)
	return false
}

func (p *Parser) appendError(format string, args ...interface{}) {
	p.errors = append(p.errors, fmt.Sprintf(format, args...))
}

func (p *Parser) Parse() (*ast.App, []error) {
	app := p.parseApp()

	// Support top-level actions/features outside the app block,
	// with tolerance for extra end tokens between blocks
	for {
		for p.curTokenIs(lexer.END) {
			p.nextToken()
		}
		if p.curTokenIs(lexer.ACTION) {
			app.Actions = append(app.Actions, p.parseAction())
		} else if p.curTokenIs(lexer.FEATURE) {
			app.Features = append(app.Features, p.parseFeature())
		} else {
			break
		}
	}

	errs := make([]error, len(p.errors))
	for i, e := range p.errors {
		errs[i] = fmt.Errorf("%s", e)
	}
	return app, errs
}

func (p *Parser) parseApp() *ast.App {
	app := &ast.App{
		Meta: make(map[string]string),
	}

	if !p.expectCurrent(lexer.APP) {
		return app
	}

	if !p.curTokenIs(lexer.IDENT) {
		p.appendError("expected app name but got %s (line=%d,col=%d)", p.curToken.Type, p.curToken.Line, p.curToken.Column)
		return app
	}
	app.Name = p.curToken.Literal
	p.nextToken()

	if p.curTokenIs(lexer.META) {
		app.Meta = p.parseMeta()
	}

	for p.curTokenIs(lexer.ACTION) || p.curTokenIs(lexer.FEATURE) {
		switch p.curToken.Type {
		case lexer.ACTION:
			app.Actions = append(app.Actions, p.parseAction())
		case lexer.FEATURE:
			app.Features = append(app.Features, p.parseFeature())
		}
	}

	if !p.curTokenIs(lexer.END) && p.curToken.Type != lexer.EOF {
		p.appendError("expected END or EOF but got %s (line=%d,col=%d)", p.curToken.Type, p.curToken.Line, p.curToken.Column)
	} else {
		p.expectCurrent(lexer.END)
	}

	return app
}

func (p *Parser) parseMeta() map[string]string {
	m := make(map[string]string)

	p.nextToken() // consume META

	for !p.curTokenIs(lexer.END) && p.curToken.Type != lexer.EOF {
		if !p.curTokenIs(lexer.IDENT) {
			p.appendError("expected meta key but got %s (line=%d,col=%d)", p.curToken.Type, p.curToken.Line, p.curToken.Column)
			p.nextToken()
			continue
		}
		key := p.curToken.Literal
		p.nextToken()

		if !p.expectCurrent(lexer.COLON) {
			p.advanceToEnd()
			return m
		}

		if !p.curTokenIs(lexer.STRING) {
			p.appendError("expected meta value but got %s (line=%d,col=%d)", p.curToken.Type, p.curToken.Line, p.curToken.Column)
			p.advanceToEnd()
			return m
		}
		value := p.curToken.Literal
		p.nextToken()

		if _, exists := m[key]; exists {
			p.appendError("duplicate meta key '%s' (line=%d,col=%d)", key, p.curToken.Line, p.curToken.Column)
		}
		m[key] = value
	}

	if !p.curTokenIs(lexer.END) {
		p.appendError("expected END in meta block (line=%d,col=%d)", p.curToken.Line, p.curToken.Column)
	} else {
		p.nextToken()
	}

	return m
}

func (p *Parser) parseAction() *ast.Action {
	a := &ast.Action{
		Meta: make(map[string]string),
	}

	p.nextToken() // consume ACTION

	if !p.curTokenIs(lexer.IDENT) {
		p.appendError("expected action name but got %s (line=%d,col=%d)", p.curToken.Type, p.curToken.Line, p.curToken.Column)
		p.advanceToEnd()
		return a
	}
	a.Name = p.curToken.Literal
	p.nextToken()

	if p.curTokenIs(lexer.META) {
		a.Meta = p.parseMeta()
	}

	if p.curTokenIs(lexer.PARAMS) {
		a.Params = p.parseParams()
	}

	if !p.curTokenIs(lexer.END) && p.curToken.Type != lexer.EOF {
		p.appendError("expected END after action but got %s (line=%d,col=%d)", p.curToken.Type, p.curToken.Line, p.curToken.Column)
		p.advanceToEnd()
	} else {
		p.expectCurrent(lexer.END)
	}

	return a
}

func (p *Parser) parseFunction() *ast.Function {
	fn := &ast.Function{}

	p.nextToken() // consume FUNCTION

	if !p.curTokenIs(lexer.IDENT) {
		p.appendError("expected function name but got %s (line=%d,col=%d)", p.curToken.Type, p.curToken.Line, p.curToken.Column)
		p.advanceToEnd()
		return fn
	}
	fn.Name = p.curToken.Literal
	p.nextToken()

	if !p.expectCurrent(lexer.LPAREN) {
		p.advanceToEnd()
		return fn
	}

	if p.curTokenIs(lexer.IDENT) {
		fn.Params = append(fn.Params, p.curToken.Literal)
		p.nextToken()
	}

	if !p.expectCurrent(lexer.RPAREN) {
		p.advanceToEnd()
		return fn
	}

	// Expect return keyword
	if !p.expectCurrent(lexer.RETURN) {
		p.advanceToEnd()
		return fn
	}

	// Collect body tokens until END
	var bodyParts []string
	for p.curToken.Type != lexer.EOF && !p.curTokenIs(lexer.END) {
		bodyParts = append(bodyParts, p.curToken.Literal)
		p.nextToken()
	}
	fn.Body = ""
	if len(bodyParts) > 0 {
		// Join tokens with spaces to reconstruct the expression
		for i, part := range bodyParts {
			if i > 0 {
				// Add space between tokens unless it's a delimiter
				last := bodyParts[i-1]
				if last != "(" && last != "." && part != ")" && part != "." && part != ":" {
					fn.Body += " "
				}
			}
			fn.Body += part
		}
	}

	// Consume END
	if !p.curTokenIs(lexer.END) {
		p.appendError("expected END after function but got %s (line=%d,col=%d)", p.curToken.Type, p.curToken.Line, p.curToken.Column)
	} else {
		p.nextToken()
	}

	return fn
}

func (p *Parser) parseFeature() *ast.Feature {
	f := &ast.Feature{
		Meta:      make(map[string]string),
		Functions: []*ast.Function{},
		Params:    []*ast.Param{},
		Rules:     []*ast.Rule{},
		Events:    []*ast.Event{},
	}

	p.nextToken() // consume FEATURE

	if !p.curTokenIs(lexer.IDENT) {
		p.appendError("expected feature name but got %s (line=%d,col=%d)", p.curToken.Type, p.curToken.Line, p.curToken.Column)
		p.advanceToEnd()
		return f
	}
	f.Name = p.curToken.Literal
	p.nextToken()

	// Loop to handle blocks in any order. The feature ends when we see
	// FEATURE, ACTION, APP (next top-level item), EOF, or a bare END
	// not followed by a sub-block keyword.
	for p.curToken.Type != lexer.EOF {
		// If the next top-level item starts, the feature ends implicitly
		// (parseRules already consumed any trailing ENDs).
		if p.curTokenIs(lexer.FEATURE) || p.curTokenIs(lexer.ACTION) ||
			p.curTokenIs(lexer.APP) {
			return f
		}

		switch {
		case p.curTokenIs(lexer.META):
			f.Meta = p.parseMeta()
		case p.curTokenIs(lexer.FUNCTION):
			f.Functions = append(f.Functions, p.parseFunction())
		case p.curTokenIs(lexer.PARAMS):
			f.Params = p.parseParams()
		case p.curTokenIs(lexer.RULES):
			f.Rules = p.parseRules()
		case p.curTokenIs(lexer.EVENTS):
			f.Events = p.parseEvents()
		case p.curTokenIs(lexer.END):
			p.nextToken()
			return f
		default:
			p.appendError("unexpected token '%s' in feature (line=%d,col=%d)",
				p.curToken.Literal, p.curToken.Line, p.curToken.Column)
			p.advanceToEnd()
			return f
		}
	}

	return f
}

func (p *Parser) parseParams() []*ast.Param {
	var params []*ast.Param

	p.nextToken() // consume PARAMS

	for !p.curTokenIs(lexer.END) && p.curToken.Type != lexer.EOF {
		if !p.expectCurrent(lexer.PARAM) {
			p.advanceToEnd()
			return params
		}

		param := p.parseParam()
		params = append(params, param)
	}

	if !p.curTokenIs(lexer.END) {
		p.appendError("expected END in params block (line=%d,col=%d)", p.curToken.Line, p.curToken.Column)
	} else {
		p.nextToken()
	}

	return params
}

func (p *Parser) parseParam() *ast.Param {
	param := &ast.Param{}

	if !p.curTokenIs(lexer.IDENT) {
		p.appendError("expected param name but got %s (line=%d,col=%d)", p.curToken.Type, p.curToken.Line, p.curToken.Column)
		p.advanceToEnd()
		return param
	}
	param.Name = p.curToken.Literal
	p.nextToken()

	for !p.curTokenIs(lexer.END) && p.curToken.Type != lexer.EOF {
		switch p.curToken.Literal {
		case "type":
			p.nextToken()
			if !p.expectCurrent(lexer.COLON) {
				p.advanceToEnd()
				return param
			}
			if p.curTokenIs(lexer.ENUM) {
				param.Type = "enum"
				p.nextToken()
				if !p.expectCurrent(lexer.LBRACE) {
					p.advanceToEnd()
					return param
				}
				for !p.curTokenIs(lexer.RBRACE) && p.curToken.Type != lexer.EOF {
					if p.curTokenIs(lexer.STRING) {
						param.EnumVals = append(param.EnumVals, p.curToken.Literal)
						p.nextToken()
					} else if p.curTokenIs(lexer.COMMA) {
						p.nextToken()
					} else {
						p.appendError("expected string or '}' in enum but got %s (line=%d,col=%d)", p.curToken.Type, p.curToken.Line, p.curToken.Column)
						p.advanceToEnd()
						return param
					}
				}
				if !p.expectCurrent(lexer.RBRACE) {
					p.advanceToEnd()
					return param
				}
				if len(param.EnumVals) == 0 {
					p.appendError("enum must have at least one value (line=%d,col=%d)", p.curToken.Line, p.curToken.Column)
				}
			} else if p.curTokenIs(lexer.IDENT) && (p.curToken.Literal == "string" || p.curToken.Literal == "integer") {
				param.Type = p.curToken.Literal
				p.nextToken()
			} else {
				p.appendError("expected type (string, integer, or enum) but got %s (line=%d,col=%d)", p.curToken.Type, p.curToken.Line, p.curToken.Column)
				p.advanceToEnd()
				return param
			}
		case "required":
			p.nextToken()
			if !p.expectCurrent(lexer.COLON) {
				p.advanceToEnd()
				return param
			}
			if p.curTokenIs(lexer.TRUE) {
				param.Required = true
				p.nextToken()
			} else if p.curTokenIs(lexer.FALSE) {
				param.Required = false
				p.nextToken()
			} else {
				p.appendError("expected true or false for required but got %s (line=%d,col=%d)", p.curToken.Type, p.curToken.Line, p.curToken.Column)
				p.advanceToEnd()
				return param
			}
		case "default":
			p.nextToken()
			if !p.expectCurrent(lexer.COLON) {
				p.advanceToEnd()
				return param
			}
			if p.curTokenIs(lexer.STRING) {
				param.Default = p.curToken.Literal
				p.nextToken()
			} else if p.curTokenIs(lexer.INTEGER) {
				val, err := strconv.Atoi(p.curToken.Literal)
				if err != nil {
					p.appendError("invalid integer default '%s' (line=%d,col=%d)", p.curToken.Literal, p.curToken.Line, p.curToken.Column)
				} else {
					param.Default = val
				}
				p.nextToken()
			} else {
				p.appendError("expected string or integer for default but got %s (line=%d,col=%d)", p.curToken.Type, p.curToken.Line, p.curToken.Column)
				p.advanceToEnd()
				return param
			}
		default:
			p.appendError("unexpected token '%s' in param (line=%d,col=%d)", p.curToken.Literal, p.curToken.Line, p.curToken.Column)
			p.nextToken()
		}
	}

	if !p.curTokenIs(lexer.END) {
		p.appendError("expected END after param (line=%d,col=%d)", p.curToken.Line, p.curToken.Column)
	} else {
		p.nextToken()
	}

	return param
}

func (p *Parser) parseRules() []*ast.Rule {
	var rules []*ast.Rule

	p.nextToken() // consume RULES

	for p.curToken.Type != lexer.EOF {
		// Separator END between rules — skip and continue
		if p.curTokenIs(lexer.END) {
			if p.peekTokenIs(lexer.RULE) {
				p.nextToken()
				continue
			}
			break // rules block's END
		}

		if !p.expectCurrent(lexer.RULE) {
			p.advanceToEnd()
			return rules
		}

		if !p.curTokenIs(lexer.IDENT) {
			p.appendError("expected rule name but got %s (line=%d,col=%d)", p.curToken.Type, p.curToken.Line, p.curToken.Column)
			p.advanceToEnd()
			return rules
		}
		name := p.curToken.Literal
		p.nextToken()

		if !p.expectCurrent(lexer.IF) {
			p.advanceToEnd()
			return rules
		}

		if p.curTokenIs(lexer.NOT) {
			// Function-call rule: if not FunctionName(args)
			p.nextToken()

			if !p.curTokenIs(lexer.IDENT) {
				p.appendError("expected function name after not but got %s (line=%d,col=%d)", p.curToken.Type, p.curToken.Line, p.curToken.Column)
				p.advanceToEnd()
				return rules
			}
			funcName := p.curToken.Literal
			p.nextToken()

			if !p.expectCurrent(lexer.LPAREN) {
				p.advanceToEnd()
				return rules
			}

			funcArgs := ""
			if p.curTokenIs(lexer.IDENT) {
				funcArgs = p.curToken.Literal
				p.nextToken()
			}

			if !p.expectCurrent(lexer.RPAREN) {
				p.advanceToEnd()
				return rules
			}

			if !p.expectCurrent(lexer.REJECT) {
				p.advanceToEnd()
				return rules
			}

			if !p.curTokenIs(lexer.END) {
				p.appendError("expected END after rule but got %s (line=%d,col=%d)", p.curToken.Type, p.curToken.Line, p.curToken.Column)
				p.advanceToEnd()
				return rules
			}
			p.nextToken()

			cond := "not " + funcName + "(" + funcArgs + ")"
			rules = append(rules, &ast.Rule{
				Name:      name,
				Condition: cond,
			})
		} else if p.curTokenIs(lexer.EMPTY) {
			// Empty-check rule: if empty(ParamName)
			p.nextToken()

			if !p.expectCurrent(lexer.LPAREN) {
				p.advanceToEnd()
				return rules
			}

			if !p.curTokenIs(lexer.IDENT) {
				p.appendError("expected parameter name in empty() but got %s (line=%d,col=%d)", p.curToken.Type, p.curToken.Line, p.curToken.Column)
				p.advanceToEnd()
				return rules
			}
			paramName := p.curToken.Literal
			p.nextToken()

			if !p.expectCurrent(lexer.RPAREN) {
				p.advanceToEnd()
				return rules
			}

			if !p.expectCurrent(lexer.REJECT) {
				p.advanceToEnd()
				return rules
			}

			if !p.curTokenIs(lexer.END) {
				p.appendError("expected END after rule but got %s (line=%d,col=%d)", p.curToken.Type, p.curToken.Line, p.curToken.Column)
				p.advanceToEnd()
				return rules
			}
			p.nextToken()

			rules = append(rules, &ast.Rule{
				Name:      name,
				Condition: "empty(" + paramName + ")",
			})
		} else {
			p.appendError("expected EMPTY or NOT after if but got %s (line=%d,col=%d)", p.curToken.Type, p.curToken.Line, p.curToken.Column)
			p.advanceToEnd()
			return rules
		}
	}

	// Consume all END tokens at the rules block boundary.
	// Stop if the next token is FEATURE, ACTION, APP, or EOF — those
	// belong to the enclosing block, not to the rules block.
	for p.curTokenIs(lexer.END) {
		if p.peekTokenIs(lexer.FEATURE) || p.peekTokenIs(lexer.ACTION) ||
			p.peekTokenIs(lexer.APP) || p.peekToken.Type == lexer.EOF {
			break
		}
		p.nextToken()
	}

	return rules
}

func (p *Parser) parseEvents() []*ast.Event {
	var events []*ast.Event

	p.nextToken() // consume EVENTS

	for !p.curTokenIs(lexer.END) && p.curToken.Type != lexer.EOF {
		switch {
		case p.curTokenIs(lexer.IF):
			p.nextToken()

			if !p.curTokenIs(lexer.IDENT) {
				p.appendError("expected identifier in condition but got %s (line=%d,col=%d)", p.curToken.Type, p.curToken.Line, p.curToken.Column)
				p.advanceToEnd()
				return events
			}
			left := p.curToken.Literal
			p.nextToken()

			if !p.expectCurrent(lexer.EQ_EQ) {
				p.advanceToEnd()
				return events
			}

			if !p.curTokenIs(lexer.STRING) {
				p.appendError("expected string value in condition but got %s (line=%d,col=%d)", p.curToken.Type, p.curToken.Line, p.curToken.Column)
				p.advanceToEnd()
				return events
			}
			right := p.curToken.Literal
			p.nextToken()

			if !p.expectCurrent(lexer.TRIGGER) {
				p.advanceToEnd()
				return events
			}

			if !p.curTokenIs(lexer.IDENT) {
				p.appendError("expected action name after trigger but got %s (line=%d,col=%d)", p.curToken.Type, p.curToken.Line, p.curToken.Column)
				p.advanceToEnd()
				return events
			}
			actionName := p.curToken.Literal
			p.nextToken()

			if !p.curTokenIs(lexer.END) {
				p.appendError("expected END after event but got %s (line=%d,col=%d)", p.curToken.Type, p.curToken.Line, p.curToken.Column)
				p.advanceToEnd()
				return events
			}
			p.nextToken()

			events = append(events, &ast.Event{
				Condition: left + " == \"" + right + "\"",
				Trigger:   actionName,
			})

		case p.curTokenIs(lexer.TRIGGER):
			p.nextToken()

			if !p.curTokenIs(lexer.IDENT) {
				p.appendError("expected action name after trigger but got %s (line=%d,col=%d)", p.curToken.Type, p.curToken.Line, p.curToken.Column)
				p.advanceToEnd()
				return events
			}
			actionName := p.curToken.Literal
			p.nextToken()

			events = append(events, &ast.Event{
				Trigger: actionName,
			})

		default:
			p.appendError("expected IF or TRIGGER but got %s (line=%d,col=%d)", p.curToken.Type, p.curToken.Line, p.curToken.Column)
			p.advanceToEnd()
			return events
		}
	}

	if !p.curTokenIs(lexer.END) {
		p.appendError("expected END in events block (line=%d,col=%d)", p.curToken.Line, p.curToken.Column)
	} else {
		p.nextToken()
	}

	return events
}

func (p *Parser) advanceToEnd() {
	for p.curToken.Type != lexer.EOF && p.curToken.Type != lexer.END && p.curToken.Type != lexer.ACTION && p.curToken.Type != lexer.FEATURE && p.curToken.Type != lexer.FUNCTION {
		p.nextToken()
	}
	if p.curTokenIs(lexer.END) {
		p.nextToken()
	}
}
