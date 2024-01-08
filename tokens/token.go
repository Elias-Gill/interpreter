package tokens

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

// token types
const (
	// keywords
	VAR      = "VAR"
	FUNCTION = "FUNCTION"
	IDENT    = "IDENT"
	IF       = "IF"
	FOR      = "FOR"
	RETURN   = "RETURN"
	DATATYPE = "DATATYPE" // a datatype declaration token

	// primitive data types
	NUMBER = "NUMBER"
	STRING = "STRING"
	BOOL   = "BOOL"

	// especial characters
	COLON     = "COLON"
	SEMICOLON = "SEMICOLON"
	EOF       = "EOF"
	ILLEGAL   = "ILLEGAL"

	// operators
	PLUS    = "PLUS"
	MINUS   = "MINUS"
	ASIGN   = "ASIGN"
	COMPARE = "COMPARE"

	// brackets and parenteses
	LBRAC = "LBRAC"
	RBRAC = "RBRAC"
	LPAR  = "LPAR"
	RPAR  = "RPAR"
)

var keywords = map[string]TokenType{
	"func":   FUNCTION,
	"var":    VAR,
	"true":   BOOL,
	"false":  BOOL,
	"si":     IF,
	"desde":  FOR,
	"return": RETURN,

	// datatypes keywords
	"entero": DATATYPE,
	"cadena": DATATYPE,
}

func TokenizeIdent(ident string) TokenType {
	if ty, ok := keywords[ident]; ok {
		return ty
	}

	return IDENT
}
