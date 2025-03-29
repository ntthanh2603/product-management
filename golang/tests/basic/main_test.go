package basic

import "testing"

func TestAddOne(t *testing.T) {
	var (
		input  = 1
		output = 2
	)
	actual := AddOne(input)
	if actual != output {
		t.Errorf("AddOne(%d) = %d, output %d", input, actual, output)
	}
}
