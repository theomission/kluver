// Copyright 2015 Christian Gärtner. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package lexer

import (
	"fmt"
)

type TokenType int

type Token struct {
	Typ TokenType
	Pos int
	Val string
}

func (i Token) String() string {
	switch i.Typ {
	case TokenEOF:
		return "TokenEOF"
	case TokenError:
		return i.Val
	}

	if len(i.Val) > 10 {
		return fmt.Sprintf("%s: %.10q[...]", i.Typ, i.Val)
	}
	return fmt.Sprintf("%s: %q", i.Typ, i.Val)
}

const (
	TokenError TokenType = iota

	TokenEOF
	TokenVoid
	TokenVersion
	TokenVersionNumber
	TokenEndStatement
	TokenImport
	TokenImportPath
	TokenVertex
	TokenEnd
	TokenFragment
	TokenGLSL
	TokenYield
	TokenActionVar
	TokenProvide
	TokenRequire
	TokenRequest
	TokenAction
	TokenTypeDef
	TokenNameDec
	TokenAssign
	TokenGLSLAction
	TokenWrite
	TokenWriteOpenBracket
	TokenWriteCloseBracket
	TokenWriteSlot
)

func (i TokenType) String() string {
	switch i {
	case TokenVoid:
		return "TokenVoid"
	case TokenVersion:
		return "TokenVersion"
	case TokenVersionNumber:
		return "TokenVersionNumber"
	case TokenEndStatement:
		return "TokenEndStatement"
	case TokenImport:
		return "TokenImport"
	case TokenImportPath:
		return "TokenImport"
	case TokenVertex:
		return "TokenVertex"
	case TokenEnd:
		return "TokenEnd"
	case TokenFragment:
		return "TokenFragment"
	case TokenGLSL:
		return "TokenGLSL"
	case TokenAction:
		return "TokenAction"
	case TokenRequire:
		return "TokenRequire"
	case TokenProvide:
		return "TokenProvide"
	case TokenRequest:
		return "TokenRequest"
	case TokenYield:
		return "TokenYield"
	case TokenActionVar:
		return "TokenActionVar"
	case TokenTypeDef:
		return "TokenTypeDef"
	case TokenNameDec:
		return "TokenNameDec"
	case TokenAssign:
		return "TokenAssign"
	case TokenGLSLAction:
		return "TokenGLSLAction"
	case TokenWrite:
		return "TokenWrite"
	case TokenWriteOpenBracket:
		return "TokenWriteOpenBracket"
	case TokenWriteCloseBracket:
		return "TokenWriteCloseBracket"
	case TokenWriteSlot:
		return "TokenWriteSlot"
	default:
		return "unknown"
	}
}
