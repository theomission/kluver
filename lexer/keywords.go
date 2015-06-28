package lexer

const (
	eof           = -1
	endStatement  = ";"
	importLib     = "@import"
	vertex        = "#---VERTEX---#"
	fragment      = "#---FRAGMENT---#"
	end           = "#---END---#"
	action        = "@"
	actionRequire = "require"
	actionProvide = "provide"
	actionYield   = "yield"
	actionRequest = "request"
	actionWrite   = "write"
	actionAssign  = "="

	writeOpenBracket  = "("
	writeCloseBracket = ")"
)
