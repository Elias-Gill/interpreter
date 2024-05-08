package evaluator

import (
	"strings"
	"testing"
)

func TestIntegerEvaluation(t *testing.T) {
	evaluated := parseAndEval(t, "1123;")

	if evaluated == nil {
		return
	}

	testInteger(t, evaluated, 1123)
}

func TestStringEvaluation(t *testing.T) {
	evaluated := parseAndEval(t, `"personal"`)

	if evaluated == nil {
		return
	}

	testString(t, evaluated, "personal")
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

func TestInfixStrings(t *testing.T) {
	testCases := []struct {
		tcase    string
		expected string
	}{
		{tcase: `"Hola" + "chau"`, expected: "Holachau"},
		{tcase: `"Hola " + "personal "`, expected: "Hola personal "},
	}

	for _, tc := range testCases {
		evaluated := parseAndEval(t, tc.tcase)

		if evaluated == nil {
			continue
		}

		testString(t, evaluated, tc.expected)
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
		{tcase: `"Hola" == "chau"`, expected: false},
		{tcase: `"Hola" == "Hola"`, expected: true},
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
		{tcase: "si(true*2){2}", expected: "Expected boolean expression for 'if' condition.\n" +
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

func TestEvalIdentifier(t *testing.T) {
	testCases := []struct {
		tcase    string
		expected int64
	}{
		{tcase: "var nuevo=2; nuevo", expected: 2},
		{tcase: "var a=3; a*2;", expected: 6},
		{tcase: "var a=3; var b=4; a+b;", expected: 7},
	}

	for _, tt := range testCases {
		testInteger(t, parseAndEval(t, tt.tcase), tt.expected)
	}
}

func TestFunctionCalls(t *testing.T) {
	testCases := []struct {
		tcase    string
		expected int64
	}{
		{
			tcase: `func algo() {
                retorna 2;
            }
            algo();`,
			expected: 2},

		{
			tcase: `func algo(a, b) {
            retorna a * b;
            }
            algo(2, 8);`,
			expected: 16},
	}

	for _, tt := range testCases {
		p := parseAndEval(t, tt.tcase)
		if p == nil {
			continue
		}
		testInteger(t, p, tt.expected)
	}
}
