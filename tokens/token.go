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
	NUMBER   = "NUMBER"
	IF       = "IF"
	FOR      = "FOR"
	RETURN   = "RETURN"

	// data types
	INTEGER = "INTEGER"
	STRING  = "STRING"
	BOOL    = "BOOL"

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
	"entero": NUMBER,
	"true":   BOOL,
	"false":  BOOL,
	"si":     IF,
	"desde":  FOR,
	"cadena": STRING,
    "return": RETURN,
}

func TokenizeIdent(ident string) TokenType {
	if ty, ok := keywords[ident]; ok {
		return ty
	}

	return IDENT
}
