package cut_test

import (
	"fmt"

	"github.com/gloo-foo/testable"

	command "github.com/gloo-foo/cmd-cut"
)

func ExampleCut_openEndedRange() {
	// echo "abcdef" | cut -b2-   (byte 2 through end)
	output, _ := testable.Test(
		command.Cut(command.CutBytes("2-")),
		"abcdef\n",
	)
	fmt.Print(output)
	// Output:
	// bcdef
}

func ExampleCut_fromStartRange() {
	// echo "abcdef" | cut -c-3   (characters 1 through 3)
	output, _ := testable.Test(
		command.Cut(command.CutChars("-3")),
		"abcdef\n",
	)
	fmt.Print(output)
	// Output:
	// abc
}
