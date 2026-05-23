package main

import "testing"

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "hello world    ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "   Go is great!  ",
			expected: []string{"go", "is", "great!"},
		},
		{
			input:    "   Multiple   spaces   here  ",
			expected: []string{"multiple", "spaces", "here"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("Expected %d words, got %d", len(c.expected), len(actual))
			return
		}
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("Expected word '%s', got '%s'", expectedWord, word)
				return
			}
		}
	}
}
