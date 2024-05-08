package evaluator

import (
	"fmt"
	"testing"

	"github.com/sl2.0/objects"
	"github.com/sl2.0/parser"
)

// --- Testing utils ---

func testBool(t *testing.T, evaluated objects.Object, expected bool) {
	if evaluated.Type() != objects.BOOL_OBJ {
		t.Errorf("Expected 'Object Boolean' type. Got %s", evaluated.Inspect())
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

func testInteger(t *testing.T, evaluated objects.Object, expected int64) {
	if evaluated.Type() != objects.INTEGER_OBJ {
		t.Errorf("Expected 'Object integer' type. Got %s", evaluated.Inspect())
		return
	}

	res, ok := evaluated.(*objects.Integer)
	if !ok {
		t.Errorf("Cannot parse to 'Object integer'")
		return
	}

	if res.Value != expected {
		t.Errorf("Expected '%d'. Got %d", expected, res.Value)
	}
}

func testString(t *testing.T, evaluated objects.Object, expected string) {
	if evaluated.Type() != objects.STRING_OBJ {
		t.Errorf("Expected 'Object String' type. Got %s", evaluated.Inspect())
		return
	}

	res, ok := evaluated.(*objects.String)
	if !ok {
		t.Errorf("Cannot parse to 'Object integer'")
		return
	}

	if res.Value != expected {
		t.Errorf("Expected '%s'. Got %s", expected, res.Value)
	}
}

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

	ev := NewFromProgram(p)
	evaluated := ev.EvalProgram(objects.NewStorage())

	if evaluated == nil {
		t.Errorf("Evaluator returned a nil value")
		return nil
	}

	return evaluated
}
