package cut_test

import (
	"fmt"

	"github.com/gloo-foo/testable"

	command "github.com/gloo-foo/cmd-cut"
)

func ExampleCut_characters() {
	// echo "Hello World" | cut -c1-5
	output, _ := testable.Test(
		command.Cut(command.CutChars("1-5")),
		"Hello World\n",
	)
	fmt.Print(output)
	// Output:
	// Hello
}
