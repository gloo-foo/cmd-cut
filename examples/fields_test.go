package cut_test

import (
	"fmt"

	command "github.com/gloo-foo/cmd-cut"
	"github.com/gloo-foo/testable"
)

func ExampleCut_fields() {
	// echo "one:two:three:four" | cut -d: -f2
	output, _ := testable.Test(
		command.Cut(command.CutDelimiter(":"), command.CutFields(2)),
		"one:two:three:four\n",
	)
	fmt.Print(output)
	// Output:
	// two
}
