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

	// data types
	INTEGER = "INTEGER"
	STRING  = "STRING"
	BOOL    = "BOOL"

	// especial characters
	COLON     = "COLON"
	SEMICOLON = "SEMICOLON"
	EOF       = "EOF"
	ILLEGAL    = "ILLEGAL"

	// operators
	PLUS  = "PLUS"
	MINUS = "MINUS"
	ASIGN = "ASIGN"

	// brackets and parenteses
	LBRAC = "LBRAC"
	RBRAC = "RBRAC"
	LPAR  = "LPAR"
	RPAR  = "RPAR"
)

var keywords = map[string]TokenType{
	"func":   FUNCTION,
	"var":    VAR,
	"entero": INTEGER,
}

func TokenizeIdent(ident string) TokenType {
	if ty, ok := keywords[ident]; ok {
		return ty
	}

	return IDENT
}
