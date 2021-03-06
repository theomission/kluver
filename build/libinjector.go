// Copyright 2015 Christian Gärtner. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package build

import (
	"fmt"
	"github.com/achtern/kluver/lexer"
)

func injectLibVertex(shader *Shader, libIndex map[int][]string) {
outer:
	for {
		for i := 0; i < len(shader.vertex); i++ {
			if shader.vertex[i].Typ != lexer.TokenYield {
				if i == len(shader.vertex)-1 {
					break outer
				}
				continue
			}
			// contains everything before the yield token
			everythingBefore := append(Tokens(nil), shader.vertex[:i]...)

			// contains everything after the yield AND endstatement token
			everythingAfter := append(Tokens(nil), shader.vertex[i+2:]...)

			shader.vertex = everythingBefore

			include := false
			for _, lib := range shader.libs {
				for _, libToken := range lib.vertex {
					switch libToken.Typ {
					case lexer.TokenExport:
						include = true
						continue
					case lexer.TokenExportEnd:
						include = false
					}
					if include {
						shader.vertex = append(shader.vertex, libToken)
					}
				}
			}

			shader.vertex = append(shader.vertex, everythingAfter...)
			// break loop
			// this way we start at the beginning again
			break
		}
	}
}

func injectLibFragment(shader *Shader, libIndex map[int][]string) {
	libGetterIdentifier := make([]string, 0)
	// _L_ib_G_etter_I_dentifier
	usedLGIs := make([]string, 0)

	newFragment := make(Tokens, 0)

	// libIndex -> template -> its tokens
	templates := make(map[int]map[string]Tokens)
	// libIndex -> supply -> its tokens
	supplies := make(map[int]map[string]Tokens)
	// libIndex -> supply -> its parent
	suppliesParent := make(map[int]map[string]string)
	addToTemplate := ""
	addToSupplies := ""
	suppliesParentName := ""

	include := false
	skipNext := 0
	for libIndex, lib := range shader.libs {
		// setup maps
		templates[libIndex] = make(map[string]Tokens)
		supplies[libIndex] = make(map[string]Tokens)
		suppliesParent[libIndex] = make(map[string]string)

		for libTokenIndex, libToken := range lib.fragment {
			if skipNext > 0 {
				skipNext -= 1
				continue
			}
			switch libToken.Typ {
			case lexer.TokenExport:
				include = true
				continue
			case lexer.TokenExportEnd:
				include = false
			case lexer.TokenGet:
				hash := GetHash(libIndex, libTokenIndex)
				libToken.Val = fmt.Sprintf("vec4 get%s", hash)
				libGetterIdentifier = append(libGetterIdentifier, hash)
				if addToSupplies == "" && addToTemplate == "" {
					// we have a "normal" lib
					// add lgi to used list
					usedLGIs = append(usedLGIs, hash)
				}
			case lexer.TokenTemplate:
				addToTemplate = lib.fragment[libTokenIndex+1].Val
				skipNext += 2 // ignore the pointer && name
				continue // ignore the template token
			case lexer.TokenTemplateEnd:
				addToTemplate = ""
			case lexer.TokenSupply:
				addToSupplies = lib.fragment[libTokenIndex+1].Val
				skipNext = 2 // skip name of the supply and the pointer
				suppliesParentName = ""

				if lib.fragment[libTokenIndex+2].Typ == lexer.TokenColon {
					// this supply extends a template
					suppliesParentName = lib.fragment[libTokenIndex+3].Val
					skipNext += 2 // skip the parent and the pointer
				}
				continue // ignore the supply token
			case lexer.TokenSupplyEnd:
				addToSupplies = ""
			}
			if include {
				newFragment = append(newFragment, libToken)
			}
			if addToTemplate != "" {
				templates[libIndex][addToTemplate] = append(templates[libIndex][addToTemplate], libToken)
			}
			if addToSupplies != "" {
				supplies[libIndex][addToSupplies] = append(supplies[libIndex][addToSupplies], libToken)
				suppliesParent[libIndex][addToSupplies] = suppliesParentName
			}
		}
	}

	includedTemplateBlocks := make([]string, 0)

	// included requsted supplies
	for i, supplyNamesReq := range libIndex {
		for i2, supply := range supplies {
			// supply is of type "map[string]Tokens"
			if i == i2 {
				for supplyName, tokens := range supply {
					for _, supplyNameReq := range supplyNamesReq {
						if supplyName == supplyNameReq {
							// add tokens of parent
							parentName := suppliesParent[i][supplyName]
							// but only if it has not been included yet
							if !ContainsString(parentName, includedTemplateBlocks) {
								tokensOfParent := templates[i][parentName]
								for _, token := range tokensOfParent {
									newFragment = append(newFragment, token)
								}
								includedTemplateBlocks = append(includedTemplateBlocks, parentName)
							}
							// add tokens of supply
							for _, token := range tokens {
								newFragment = append(newFragment, token)
								if token.Typ == lexer.TokenGet {
									// remove the first 8 characters (vec4 get)
									usedLGIs = append(usedLGIs, token.Val[8:])
								}
							}
						}
					}
				}
			}
		}
	}

	shader.fragment = append(newFragment, shader.fragment...)

	// for every yield token, we have to call the @get functions of all libs
	for i := 0; i < len(shader.fragment); i++ {
		if shader.fragment[i].Typ == lexer.TokenYield {
			var sb StringBuffer
			for _, hash := range usedLGIs {

				sb.Append(fmt.Sprintf(
					"\t%s = get%s(%s);\n",    // fn call
					shader.fragment[i+1].Val, // actionVar
					hash, // libGetterIdentifier
					shader.fragment[i+1].Val)) // actionVar
			}
			shader.fragment[i].Val = sb.String()
		}
	}
}
