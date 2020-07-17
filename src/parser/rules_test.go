package parser

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

const rulesA = `
rules_version = '2'; 
service cloud.firestore {
	match /foo/{bar} {
	}
}
`

func GetFileContent(s string) string {
	result, err := ioutil.ReadFile(s)
	if err != nil {
		panic("failed to read file")
	}
	return string(result)

}

func TestRules(t *testing.T) {
	tests := []struct {
		name string
		input string
		expected string
	}{
		{
			name:     "basic",
			input:    rulesA,
			expected: "",
		},
		{
			name: "file",
			input: GetFileContent("../../testdata/firestore.rules"),
			expected: "",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tokens := New(test.input)
			rules, err := ParseRules(tokens)
			assert.Nil(t, err)
			assert.Equal(t, test.expected, rules.String())
		})
	}
}

