package evaluator

import (
	"fmt"
	"testing"

	"github.com/sl2.0/objects"
	"github.com/sl2.0/parser"
)

func TestIntegerEvaluation(t *testing.T) {
	evaluated := parseAndEval(t, "1123;")

	if evaluated == nil {
		return
	}

	if evaluated.Type() != objects.INTEGER_OBJ {
		t.Fatalf("Expected 'Object integer' type. Got %s", evaluated.Type())
	}

	value, ok := evaluated.(*objects.Integer)
	if !ok {
		t.Fatalf("Cannot parse to 'Object integer'")
	}

	if value.Value != 1123 {
		t.Fatalf("Expected '1123'. Got %d", value.Value)
	}
}

func TestBooleanEvaluation(t *testing.T) {
	evaluated := parseAndEval(t, "true")

	if evaluated == nil {
		return
	}

	testBool(t, evaluated, true)
}

func TestBangOperator(t *testing.T) {
	evaluated := parseAndEval(t, "!true")

	if evaluated == nil {
		return
	}

	testBool(t, evaluated, false)
}

// --- Testing utils ---

func parseAndEval(t *testing.T, input string) objects.Object {
	const colorMagenta = "\033[35m"
	const colorNone = "\033[0m"

	parser := parser.NewParser(input)
	p := parser.ParseProgram()

	// check errors
	if len(parser.Errors()) != 0 {
		s := ""

		for _, value := range parser.Errors() {
			s += fmt.Sprintf("\n%s", value)
		}

		t.Errorf("Found %s%d%s parsing errors: %s",
			colorMagenta, len(parser.Errors()), colorNone, s)
		return nil
	}

	if p == nil {
		t.Errorf("ParseProgram() returned nil")
		return nil
	}

	if len(p.Statements) == 0 {
		t.Errorf("Empty AST")
		return nil
	}

	evaluated := Eval(p)
	if evaluated == nil {
		t.Errorf("Evaluator returned a nil value")
		return nil
	}

	return evaluated
}

func testBool(t *testing.T, evaluated objects.Object, expected bool) {
	if evaluated.Type() != objects.BOOL_OBJ {
		t.Errorf("Expected 'Object Boolean' type. Got %s", evaluated.Type())
		return
	}

	value, ok := evaluated.(*objects.Boolean)
	if !ok {
		t.Errorf("Cannot parse to 'Object integer'")
		return
	}

	if value.Value != expected {
		t.Errorf("Expected '%v'. Got %v", expected, value.Value)
	}
}
