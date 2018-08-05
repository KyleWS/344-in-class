package testing

import (
	"testing"
)

func TestReverse(t *testing.T) {
	cases := []struct {
		name           string
		input          string
		expectedOutput string
	}{
		{
			name:           "empty string",
			input:          "",
			expectedOutput: "",
		},
		{
			name:           "single character",
			input:          "a",
			expectedOutput: "a",
		},
		{
			name:           "two characters",
			input:          "ab",
			expectedOutput: "ba",
		},
		{
			name:           "stressed",
			input:          "stressed",
			expectedOutput: "desserts",
		},
		{
			name:           "high unicode",
			input:          "Hello, 世界",
			expectedOutput: "界世 ,olleH",
		},
	}

	for _, c := range cases {
		if output := Reverse(c.input); output != c.expectedOutput {
			t.Errorf("%s: got %s but expected %s", c.name, output, c.expectedOutput)
		}
	}
}

func TestGetGreeting(t *testing.T) {
	cases := []struct {
		name           string
		input          string
		expectedOutput string
	}{
		{
			name:           "empty string",
			input:          "",
			expectedOutput: "Hello, World!",
		},
		{
			name:           "non-empty string",
			input:          "Gorgeous",
			expectedOutput: "Hello, Gorgeous!",
		},
	}

	for _, c := range cases {
		if output := GetGreeting(c.input); output != c.expectedOutput {
			t.Errorf("%s: got %s but expected %s", c.name, output, c.expectedOutput)
		}
	}
}

func TestParseSize(t *testing.T) {
	cases := []struct {
		name           string
		input          string
		expectedOutput *Size
	}{
		{
			name:           "empty string",
			input:          "",
			expectedOutput: &Size{},
		},
		{
			name:           "valid",
			input:          "10x20",
			expectedOutput: &Size{10, 20},
		},
		{
			name:           "invalid height",
			input:          "10xfoo",
			expectedOutput: &Size{10, 0},
		},
	}

	for _, c := range cases {
		if output := ParseSize(c.input); output.Height != c.expectedOutput.Height || output.Width != c.expectedOutput.Width {
			t.Errorf("%s: got %v but expected %v", c.name, output, c.expectedOutput)
		}
	}
}

func TestLateDaysConsume(t *testing.T) {
	ld := NewLateDays()
	//Consume() should return 3 then 2 then 1 then 0
	//then continue to return 0 after that
	for i := 3; i > -10; i-- {
		expectedOutput := i
		if expectedOutput < 0 {
			expectedOutput = 0
		}
		if output := ld.Consume("test"); output != expectedOutput {
			t.Errorf("iteration %d: got %d but expected %d", i, output, expectedOutput)
		}
	}
}
