// Copyright (c) 2017 Gorillalabs. All rights reserved.

package powershell

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
)

func randomString(bytes int) string {
	c := bytes
	b := make([]byte, c)

	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}

	return hex.EncodeToString(b)
}

func createBoundary() string {
	return fmt.Sprintf("#powershell_boundary_%s#", randomString(12))
}

func quoteArg(s string) string {
	return "'" + strings.Replace(s, "'", "\"", -1) + "'"
}