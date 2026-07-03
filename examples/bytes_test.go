package cut_test

import (
	"fmt"

	"github.com/gloo-foo/testable"

	command "github.com/gloo-foo/cmd-cut"
)

func ExampleCut_bytes() {
	// echo "abcdefgh" | cut -b1,3,5
	output, _ := testable.Test(
		command.Cut(command.CutBytes("1,3,5")),
		"abcdefgh\n",
	)
	fmt.Print(output)
	// Output:
	// ace
}
