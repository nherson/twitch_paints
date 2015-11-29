package main

import "testing"

func TestStripInput(t *testing.T) {
	// See https://golang.org/doc/code.html#Testing
	// Make a simple struct for input and output
	mp := MessageProcessor{make(chan PaintCommand)}
	cases := []struct {
		in, out string
	}{
		{"   (24, 25,26, 123456)    ", "(24,25,26,123456)"},
		{"(åß,®,œ,®)", "(,,,)"},
		{"", ""},
	}
	// Iterate over the inputs and outputs and make sure they match up
	for _, c := range cases {
		result := mp.StripInput(c.in)
		if result != c.out {
			t.Errorf("Expected: %q, Got: %q", c.out, result)
		}
	}
}

func TestParse(t *testing.T) {
	mp := MessageProcessor{make(chan PaintCommand)}

	// Happy path testing
	goodCases := []struct {
		in  string
		out PaintCommand
	}{ // Add more cases here as needed
		{"(24,25,123456)", PaintCommand{X: 24, Y: 25, Red: 0x12, Green: 0x34, Blue: 0x56}},
		{"(1,5,abcdef)", PaintCommand{X: 1, Y: 5, Red: 0xab, Green: 0xcd, Blue: 0xef}},
		{"(799,599,090909)", PaintCommand{X: 799, Y: 599, Red: 0x9, Green: 0x9, Blue: 0x9}},
		{"(2,25,0B1D2F)", PaintCommand{X: 2, Y: 25, Red: 0x0b, Green: 0x1d, Blue: 0x2f}},
		{"(0,0,000000)", PaintCommand{X: 0, Y: 0, Red: 0x0, Green: 0x0, Blue: 0x0}},
	}
	for _, c := range goodCases {
		r, e := mp.Parse(c.in)
		// There should be no errors and the result matches the expected
		if e != nil || *r != c.out {
			t.Errorf("Error with input %q, expected %q and got %q", c.in, c.out, r)
		}
	}

	// Sad path testing
	badCases := []string{ // Add more cases here as needed
		"",
		"()",
		"(1,,)",
		"(",
		"23,23,ffffff)",
		"(25,aaaaaa)",
		"(2)",
		"(800,50,123456)", // out of bounds
		"(30,600,123456)", // out of bounds
		"(-56,1,111111)",  // negative X
		"(1,-56,111111)",  // negative Y
	}
	for _, badCase := range badCases {
		r, e := mp.Parse(badCase)
		if r != nil {
			t.Errorf("Bad input to MessageProcessor.Parse returned non-nil: %q", badCase)
		} else if e == nil {
			t.Errorf("Bad input to MessageProcessor.Parse didn't return an error")
		}
	}
}
