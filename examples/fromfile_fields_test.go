package cut_test

import (
	"fmt"
	"os"

	command "github.com/gloo-foo/cmd-cut"
	"github.com/gloo-foo/testable"
)

func ExampleCut_fromFile_fields() {
	// cut -d: -f2 testdata/fields.txt
	data, err := os.ReadFile("testdata/fields.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "read testdata: %v\n", err)
		return
	}
	output, _ := testable.Test(
		command.Cut(command.CutDelimiter(":"), command.CutFields(2)),
		string(data),
	)
	fmt.Print(output)
	// Output:
	// two
	// beta
	// second
}
