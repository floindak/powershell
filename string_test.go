// Copyright (c) 2017 Gorillalabs. All rights reserved.

package powershell

import "testing"

func TestRandomStrings(t *testing.T) {
	r1 := randomString(8)
	r2 := randomString(8)

	if r1 == r2 {
		t.Error("Failed to create random strings: The two generated strings are identical.")
	} else if len(r1) != 16 {
		t.Errorf("Expected the random string to contain 16 characters, but got %d.", len(r1))
	}
}

func TestQuotingArguments(t *testing.T) {
	testcases := [][]string{
		{"", "''"},
		{"test", "'test'"},
		{"two words", "'two words'"},
		{"quo\"ted", "'quo\"ted'"},
		{"quo'ted", "'quo\"ted'"},
		{"quo\\'ted", "'quo\\\"ted'"},
		{"quo\"t'ed", "'quo\"t\"ed'"},
		{"es\\caped", "'es\\caped'"},
		{"es`caped", "'es`caped'"},
		{"es\\`caped", "'es\\`caped'"},
	}

	for i, testcase := range testcases {
		quoted := quoteArg(testcase[0])

		if quoted != testcase[1] {
			t.Errorf("test %02d failed: input '%s', expected %s, actual %s", i+1, testcase[0], testcase[1], quoted)
		}
	}
}
