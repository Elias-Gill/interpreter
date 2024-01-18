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
	ELSE     = "ELSE"
	FOR      = "FOR"
	RETURN   = "RETURN"
	DATATYPE = "DATATYPE" // a datatype declaration token

	// primitive data types
	NUMBER = "NUMBER"
	STRING = "STRING"
	TRUE   = "TRUE"
	FALSE  = "FALSE"

	// especial characters
	COLON     = "COLON"
	SEMICOLON = "SEMICOLON"
	EOF       = "EOF"
	LINEBREAK = "LINEBREAK"
	ILLEGAL   = "ILLEGAL"

	// operators
	PLUS     = "PLUS"
	MINUS    = "MINUS"
	ASTERISC = "ASTERISC"
	BANG     = "BANG"
	COMMA    = "COMMA"
	ASIGN    = "ASIGN"  // =
	EQUALS   = "EQUALS" // ==
	NOTEQUAL = "NOTEQUAL"
	SLASH    = "STROKE"

	// brackets and parenteses
	LBRAC = "LBRAC"
	RBRAC = "RBRAC"
	LPAR  = "LPAR"
	RPAR  = "RPAR"
	LT    = "LT"
	GT    = "GT"
)

var keywords = map[string]TokenType{
	"func":    FUNCTION,
	"var":     VAR,
	"si":      IF,
	"sino":    ELSE,
	"desde":   FOR,
	"retorna": RETURN,

	// datatype keywords
	"entero": DATATYPE,
	"cadena": DATATYPE,
	"true":   TRUE,
	"false":  FALSE,
}

func ResolveType(ident string) TokenType {
	if tType, ok := keywords[ident]; ok {
		return tType
	}

	return IDENT
}
