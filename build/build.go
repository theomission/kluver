// Copyright 2015 Christian Gärtner. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package build

import (
	"fmt"
	"github.com/achtern/kluver/lexer"
	"strings"
)

type Shader struct {
	version  string
	vertex   Tokens
	fragment Tokens
	global   Tokens
	compiled GLSL
}

type GLSL struct {
	vertex   string
	fragment string
}

type varDef struct {
	typ  string
	name string
}

type Tokens []lexer.Token

const providePlaceholder = "___PROVIDE___REPLACE___HERE___\n"

func (s *Shader) String() string {
	return fmt.Sprintf("Shader(%s, vertex=%q, fragment=%q, global=%q)", s.version, s.vertex, s.fragment, s.global)
}

func Build(tokenStream <-chan lexer.Token) string {
	global := make(Tokens, 0)
	vertex := make(Tokens, 0)
	fragment := make(Tokens, 0)

	shader := Shader{}

	// phase 0 : global
	// phase 1 : vertex
	// phase 2 : fragment
	phase := 0
	for token := range tokenStream {
		switch token.Typ {
		case lexer.TokenVertex:
			phase = 1
			continue
		case lexer.TokenFragment:
			phase = 2
			continue
		case lexer.TokenEnd:
			phase = 0
			continue
		case lexer.TokenVoid:
			continue
		}

		if token.Typ == lexer.TokenVersionNumber {
			shader.version = token.Val
			continue
		}

		if token.Typ == lexer.TokenAction {
			// we do need the action tokens after lexing
			continue
		}

		switch phase {
		case 0:
			global = append(global, token)
		case 1:
			vertex = append(vertex, token)
		case 2:
			fragment = append(fragment, token)
		default:
			panic("unknow phase")
		}
	}

	shader.global = global
	shader.vertex = vertex
	shader.fragment = fragment

	s, _ := shader.buildVertex()
	return s
}

func (shader *Shader) buildHead() string {
	return "#version " + shader.version + "\n"
}

func (shader *Shader) buildVertex() (string, error) {

	var sb StringBuffer
	sb.append(shader.buildHead())

	providePlaceholderInserted := false
	provides := make([]Tokens, 0)

	for i := 0; i < len(shader.vertex); i++ {
		token := shader.vertex[i]
		switch token.Typ {
		case lexer.TokenRequire:
			if !providePlaceholderInserted {
				// insert provides before the first uniform
				sb.append(providePlaceholder)
				providePlaceholderInserted = true
			}
			sb.append(generateRequire(token, shader.vertex[i+1], shader.vertex[i+2]))
			i += 2
		case lexer.TokenProvide:
			provides = append(provides, Tokens{shader.vertex[i+1], shader.vertex[i+2]})
			sb.append(generateProvideSetting(shader.vertex[i+1], shader.vertex[i+2], shader.vertex[i+3], shader.vertex[i+4]))
			i += 4
		default:
			sb.append(token.Val)
		}
	}

	compiled := sb.String()

	compiled = strings.Replace(compiled, providePlaceholder, generateProvideDecBlock(provides), -1)

	return compiled, nil
}

func (shader *Shader) buildFragment() (string, error) {

	var sb StringBuffer
	sb.append(shader.buildHead())

	for _, token := range shader.fragment {
		sb.append(token.Val)
	}

	return sb.String(), nil
}
