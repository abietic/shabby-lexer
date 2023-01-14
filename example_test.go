// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package shabbylexer_test

import (
	"fmt"

	"github.com/abietic/shabbylexer"
	"github.com/abietic/shabbylexer/token"
)

func ExampleScanner_Scan() {
	// src is the input that we want to tokenize.
	// src := []byte("cos(x) + 1i*sin(x) // Euler")
	src := []byte("cos(x) + 1.30 * sin(x) // Test shabby-lexer")

	// Initialize the scanner.
	var s shabbylexer.Lexer
	fset := token.NewFileSet()                      // positions are relative to fset
	file := fset.AddFile("", fset.Base(), len(src)) // register input "file"
	s.Init(file, src, nil /* no error handler */, shabbylexer.ScanComments)

	// Repeated calls to Scan yield the token sequence found in the input.
	for {
		pos, tok, lit := s.Scan()
		if tok == token.EOF {
			break
		}
		fmt.Printf("%s\t%s\t%q\n", fset.Position(pos), tok, lit)
	}

	// output:
	// 1:1	IDENT	"cos"
	// 1:4	(	""
	// 1:5	IDENT	"x"
	// 1:6	)	""
	// 1:8	+	""
	// 1:10	FLOAT	"1.30"
	// 1:15	*	""
	// 1:17	IDENT	"sin"
	// 1:20	(	""
	// 1:21	IDENT	"x"
	// 1:22	)	""
	// 1:24	;	"\n"
	// 1:24	COMMENT	"// Test shabby-lexer"
}
