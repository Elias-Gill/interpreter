package evaluator

import (
	"fmt"
	"strings"
	"testing"

	"github.com/sl2.0/objects"
	"github.com/sl2.0/parser"
)

func TestIntegerEvaluation(t *testing.T) {
	evaluated := parseAndEval(t, "1123;")

	if evaluated == nil {
		return
	}

	testInteger(t, evaluated, 1123)
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

func TestMinusOperator(t *testing.T) {
	evaluated := parseAndEval(t, "-12")

	if evaluated == nil {
		return
	}

	testInteger(t, evaluated, -12)
}

func TestInfixArithmetic(t *testing.T) {
	testCases := []struct {
		tcase    string
		expected int
	}{
		{tcase: "-12 + 24 - -12", expected: 24},
		{tcase: "-12 - 12 * -2 ", expected: 12},
		{tcase: "(-12 + 24) * 2 ", expected: 24},
		{tcase: "-(11 + 1) * 2 ", expected: -24},
	}

	for _, tc := range testCases {
		evaluated := parseAndEval(t, tc.tcase)

		if evaluated == nil {
			continue
		}

		testInteger(t, evaluated, int64(tc.expected))
	}
}

func TestInfixComparison(t *testing.T) {
	testCases := []struct {
		tcase    string
		expected bool
	}{
		{tcase: "-12 + 24 < -12", expected: false},
		{tcase: "-12 + 24 > -12", expected: true},
		{tcase: "(-12 + 24) == 12 ", expected: true},
		{tcase: "-(11 + 1) != 2 ", expected: true},
	}

	for _, tc := range testCases {
		evaluated := parseAndEval(t, tc.tcase)

		if evaluated == nil {
			continue
		}

		testBool(t, evaluated, tc.expected)
	}
}

func TestIfEvaluation(t *testing.T) {
	testCases := []struct {
		tcase    string
		expected bool
	}{
		{tcase: "si(-12 + 24 < -12){true}sino{false}", expected: false},
		{tcase: "si(true){true}sino{false}", expected: true},
		{tcase: "si(false){true}sino{false}", expected: false},
	}

	for _, tc := range testCases {
		evaluated := parseAndEval(t, tc.tcase)

		if evaluated == nil {
			continue
		}

		testBool(t, evaluated, tc.expected)
	}
}

func TestReturnStatement(t *testing.T) {
	testCases := []struct {
		tcase    string
		expected int64
	}{
		{tcase: "2*8;retorna 2; 2*2", expected: 2},
		{tcase: "si(true){retorna 123}; true", expected: 123},
	}

	for _, tc := range testCases {
		evaluated := parseAndEval(t, tc.tcase)

		if evaluated == nil {
			continue
		}

		testInteger(t, evaluated, tc.expected)
	}
}

func TestEvalError(t *testing.T) {
	testCases := []struct {
		tcase    string
		expected string
	}{
		{tcase: "2*true;", expected: "Expected right value of '*' to be an integer."},
		{tcase: "true*2;", expected: "Expected right value to be a boolean."},
		{tcase: "si(true*2){2}", expected: 
            "Expected boolean expression for 'if' condition.\n"+
                "\tExpected right value to be a boolean." + 
                "\n\tGot: 2",
        },
	}

	for _, tc := range testCases {
		evaluated := parseAndEval(t, tc.tcase)

		if evaluated == nil {
			continue
		}

		if !strings.HasPrefix(evaluated.Inspect(), tc.expected) {
			t.Errorf(
				"Bad message:\nExpected: \n%s\nActual: \n%s",
				tc.expected,
				evaluated.Inspect(),
			)
		}
	}
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

	ev := NewFromProgram(p)
	evaluated := ev.EvalProgram()

	if evaluated == nil {
		t.Errorf("Evaluator returned a nil value")
		return nil
	}

	return evaluated
}

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

func testInteger(t *testing.T, evaluated objects.Object, expected int64) bool {
	if evaluated.Type() != objects.INTEGER_OBJ {
		t.Errorf("Expected 'Object integer' type. Got %s", evaluated.Inspect())
		return false
	}

	res, ok := evaluated.(*objects.Integer)
	if !ok {
		t.Errorf("Cannot parse to 'Object integer'")
		return false
	}

	if res.Value != expected {
		t.Errorf("Expected '%d'. Got %d", expected, res.Value)
	}

	return true
}
