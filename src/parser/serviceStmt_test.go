package parser

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const input1 = `
service cloud.firestore {
	match /databases/{database}/documents {
		allow read, write: if false;
	}
}
`
const input2 = `
service cloud.firestore {
	match /databases/{database}/documents {
		function fish() {
			let fowl = fly();
			return fowl;
		}
		allow read, write: if fish();
	}
}
`
func TestParseServiceStmt(t *testing.T) {
	tests := []struct {
		name string
		input string
		expected string
	}{
		{
			name: "empty",
			input: "service cloud.firestore {}",
			expected: "service cloud.firestore {  }",
		},
		{
			name: "match",
			input: input1,
			expected: "service cloud.firestore { match /databases/{database}/documents {allow read, write: if false;} }",
		},
		{
			name: "match2",
			input: input2,
			expected: "service cloud.firestore { match /databases/{database}/documents {function fish () { let fowl = fly(); return fowl; }  allow read, write: if fish();} }",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tokens := New(test.input)
			ss, err := ParseServiceStmt(tokens)
			assert.Nil(t, err)
			assert.Equal(t, test.expected, ss.String())

		})
	}
}

