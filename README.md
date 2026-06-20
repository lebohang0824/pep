# Pep

Pep is a spec language for defining application behaviors, actions, and features. The Pep compiler reads `.pep` files and produces a JSON specification file consumable by AI-powered code generators, documentation tools, and scaffolding pipelines.

## Installation

```bash
git clone <repo-url> && cd pep
go build -o pep .
```

Requires Go 1.21+.

## Usage

```bash
pep file.pep           # compiles to file.json
pep file.pep -o out.json   # custom output path
pep file.pep -p            # pretty-printed JSON
```

The compiler prints a clean progress summary:

```
  pep compiler v0.1.0
  ───────────────────────────────────────

  • Reading file.pep ...
    ✓ 2779 bytes read
  • Parsing ...
    ✓ 1 actions, 5 features
  • Semantic analysis ...
    ✓ all checks passed
  • Generating JSON ...
    ✓ 1780 bytes

  ───────────────────────────────────────
  output → file.json
```

## Pep Language

A `.pep` file describes an application — its metadata, actions, and features with validation rules and event triggers.

### Structure

```
app AppName
    meta
        key: "value"
    end

    action ActionName
        params
            param ParamName
                type: string
                required: true
            end
        end
    end

    feature FeatureName
        meta
            description: "..."
        end

        params
            param ParamName
                type: integer
                required: true
            end
        end

        rules
            rule RuleName
                if empty(ParamName)
                    reject
            end
        end

        events
            if Field == "value"
                trigger ActionName
            end
        end
    end
end
```

### Language Features

| Construct | Description |
|---|---|
| `app` | Root block — declares the application |
| `meta` | Key-value metadata (`key: "value"`) |
| `action` | Declares an action with typed params |
| `feature` | Declares a feature with params, rules, events, and functions |
| `param` | Parameter with `type`, `required`, `default` properties |
| `rule` | Validation rule — supports `empty()` and function calls |
| `events` | Event triggers — conditional (`if cond trigger action end`) or unconditional (`trigger action`) |
| `function` | Reusable validation function with `return` |

### Parameter Types

- `string`
- `integer`
- `enum{"VAL1","VAL2"}`

### Rules

```
if empty(ParamName)          reject  end
if not FuncName(args)        reject  end
```

### Events

```pep
# conditional
if Status == "DONE"
    trigger SendCompletionEmail
end

# unconditional
trigger NotifyAdmin
```

## VS Code Extension

The `vscode-pep/` directory contains a VS Code extension providing syntax highlighting and snippets for `.pep` files.

```bash
cp -r vscode-pep ~/.vscode/extensions/pep-lang.pep-lang
# Reload VS Code
```

## Project Structure

```
pep/
├── main.go                  # Entry point
├── cmd/root.go              # CLI (Cobra)
├── internal/
│   ├── ast/                 # AST node definitions
│   ├── lexer/               # Tokenizer (35+ token types)
│   ├── parser/              # Recursive descent parser
│   ├── semantic/            # Semantic analysis & validation
│   └── generator/           # JSON output generator
├── testdata/
│   ├── valid.pep            # Sample input
│   └── expected.json        # Golden test output
└── vscode-pep/              # VS Code language extension
```

## License

MIT
