package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"pep/internal/lexer"
	"pep/internal/parser"
	"pep/internal/semantic"
	"pep/internal/generator"
)

var (
	outputFile string
	pretty     bool
)

func step(msg string, args ...interface{}) {
	fmt.Printf("  \u2022 %s ...\n", fmt.Sprintf(msg, args...))
}

func ok(msg string, args ...interface{}) {
	fmt.Printf("    \u2713 %s\n", fmt.Sprintf(msg, args...))
}

var rootCmd = &cobra.Command{
	Use:   "pep [input-file]",
	Short: "Pep spec language compiler - compile .pep files to JSON.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		inputPath := args[0]

		fmt.Printf("\n  pep compiler v%s\n", "0.1.0")
		fmt.Println("  \u2500" + strings.Repeat("\u2500", 38))
		fmt.Println()

		step("Reading %s", inputPath)
		data, err := os.ReadFile(inputPath)
		if err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}
		ok("%d bytes read", len(data))

		step("Parsing")
		l := lexer.NewLexer(string(data))
		p := parser.NewParser(l)
		app, parseErrs := p.Parse()
		if len(parseErrs) > 0 {
			for _, e := range parseErrs {
				fmt.Fprintf(os.Stderr, "\n  \u2717 Parse error: %s\n", e)
			}
			return fmt.Errorf("parsing failed")
		}
		ok("%d actions, %d features", len(app.Actions), len(app.Features))

		step("Semantic analysis")
		semErrs := semantic.Analyze(app)
		if len(semErrs) > 0 {
			for _, e := range semErrs {
				fmt.Fprintf(os.Stderr, "\n  \u2717 Semantic error: %s\n", e)
			}
			return fmt.Errorf("semantic analysis failed")
		}
		ok("all checks passed")

		step("Generating JSON")
		jsonData, err := generator.Generate(app, pretty)
		if err != nil {
			return fmt.Errorf("generation failed: %w", err)
		}
		ok("%d bytes", len(jsonData))

		outPath := outputFile
		if outPath == "" {
			ext := filepath.Ext(inputPath)
			outPath = inputPath[:len(inputPath)-len(ext)] + ".json"
		}

		err = os.WriteFile(outPath, jsonData, 0644)
		if err != nil {
			return fmt.Errorf("failed to write output: %w", err)
		}

		fmt.Println()
		fmt.Println("  \u2500" + strings.Repeat("\u2500", 38))
		fmt.Printf("  output \u2192 %s\n", outPath)
		fmt.Println()

		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file path")
	rootCmd.Flags().BoolVarP(&pretty, "pretty", "p", false, "Pretty print JSON")
}
