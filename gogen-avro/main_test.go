package main

import "testing"

func TestCodegenComment(t *testing.T) {
	/*
		tt := []struct {
			in  []string
			out string
		}{
			{[]string{"file_1.avsc", "file_2.avsc"}, `// Code generated by github.com/clear-street/gogen-avro. DO NOT EDIT. */
	/*
	 * SOURCES:
	 *     file_1.avsc
	 *     file_2.avsc
	 */
	/*
	 * SOURCE:
	 *     file_1.avsc
	 */

	/*`},
	}

	/*	for _, tc := range tt {
		assert.Equal(t, tc.out, codegenComment(tc.in))
	}*/
}