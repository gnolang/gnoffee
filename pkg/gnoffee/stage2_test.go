package gnoffee

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"testing"
)

func TestStage2(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantOutput string
		wantErr    bool
	}{
		{
			name: "Basic Test",
			input: `
				package test

				//gnoffee:export baz as Helloer

				type Helloer interface {
				    Hello() string
				}

				type foo struct{}

				func (f *foo) Hello() string {
				    return "Hello from foo!"
				}

				func (f *foo) Bye() { }

				var baz = foo{}

				var _ Helloer = &foo{}
			`,
			wantOutput: `
				package test

				// This function was generated by gnoffee due to the export directive.
				func Hello() string {
					return baz.Hello()
				}
			`,
			wantErr: false,
		},
		{
			name: "Invalid Export Syntax",
			input: `
				package test

				var foo struct{}
				//gnoffee:export foo MyInterface3
				type MyInterface3 interface {
					Baz()
				}
			`,
			wantErr: true,
		},
		{
			name: "Already Annotated With gnoffee Comment",
			input: `
				package test

				var foo = struct{}

				//gnoffee:export foo as MyInterface4
				type MyInterface4 interface {
					Qux()
				}
			`,
			wantOutput: `
				package test

				// This function was generated by gnoffee due to the export directive.
				func Qux() {
					foo.Qux()
				}
			`,
			wantErr: false,
		},
		{
			name: "No Export Directive",
			input: `
				package test

				type SimpleInterface interface {
					Moo()
				}
			`,
			wantOutput: `
				package test
			`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "", tt.input, parser.ParseComments)
			if err != nil {
				t.Fatalf("Failed to parse input: %v", err)
			}

			files := map[string]*ast.File{
				"test.go": file,
			}

			generatedFile, err := Stage2(files)
			switch {
			case err == nil && tt.wantErr:
				t.Fatalf("Expected an error")
			case err != nil && !tt.wantErr:
				t.Fatalf("Error during Stage2 generation: %v", err)
			case err != nil && tt.wantErr:
				return
			case err == nil && !tt.wantErr:
				// noop
			}

			var buf bytes.Buffer
			if err := format.Node(&buf, fset, generatedFile); err != nil {
				t.Fatalf("Failed to format generated output: %v", err)
			}

			generatedCode := normalizeGoCode(buf.String())
			expected := normalizeGoCode(tt.wantOutput)
			if generatedCode != expected {
				t.Errorf("Transformed code does not match expected output.\nExpected:\n\n%v\n\nGot:\n\n%v", expected, generatedCode)
			}
		})
	}
}
