package command

import (
	"bytes"
	"strconv"
	"strings"

	gloo "github.com/gloo-foo/framework"
	"github.com/gloo-foo/framework/patterns"
)

// position is a 1-based index into a line's bytes, runes, or fields.
type position int

// interval is an inclusive 1-based range of positions. An unbounded upper end
// (a trailing "N-" spec) is modelled by hi == unbounded, which contains every
// position from lo onward.
type interval struct {
	lo position
	hi position
}

// unbounded marks an interval whose upper end runs to the end of the line.
const unbounded position = 0

// contains reports whether pos falls inside the interval.
func (iv interval) contains(pos position) bool {
	if pos < iv.lo {
		return false
	}
	return iv.hi == unbounded || pos <= iv.hi
}

// selection is a parsed position spec together with the complement flag.
type selection struct {
	intervals      []interval
	isComplemented bool
}

// selected reports whether the 1-based pos should be emitted, applying the
// complement flag (XOR): a member position is emitted unless complemented.
func (s selection) selected(pos position) bool {
	in := false
	for _, iv := range s.intervals {
		if iv.contains(pos) {
			in = true
			break
		}
	}
	return in != s.isComplemented
}

// Cut returns a Command that selects fields, bytes, or characters from each
// input line, equivalent to the Unix cut(1) utility.
//
// Modes (mutually exclusive, resolved in this order):
//   - CutBytes: select bytes by position (-b)
//   - CutChars: select characters (runes) by position (-c)
//   - CutFields (with CutDelimiter): select fields by position (-f/-d)
//
// CutComplement inverts the selection for any mode. Specs accept comma-separated
// 1-based positions and ranges, including open-ended ranges: "1,3-5", "2-", "-3".
func Cut(opts ...any) gloo.Command[[]byte, []byte] {
	cfg := parseConfig(opts)
	switch {
	case cfg.bytesSpec != "":
		return positionCommand(parseSpec(cfg.bytesSpec, cfg.complementEnabled), selectBytes)
	case cfg.charsSpec != "":
		return positionCommand(parseSpec(cfg.charsSpec, cfg.complementEnabled), selectChars)
	default:
		return fieldsCommand(cfg)
	}
}

// positionSpec is a comma-separated list of 1-based positions and ranges,
// e.g. "1-3,5,7-": the argument grammar shared by the -b and -c specs.
type positionSpec string

// cutConfig holds parsed option values for the Cut command.
type cutConfig struct {
	delimiter         CutDelimiter
	bytesSpec         positionSpec
	charsSpec         positionSpec
	fields            []position
	complementEnabled cutComplementFlag
}

// parseConfig folds the variadic options into a cutConfig.
func parseConfig(opts []any) cutConfig {
	var cfg cutConfig
	for _, o := range opts {
		cfg = cfg.with(o)
	}
	return cfg
}

// with folds a single option value into the config, returning the updated
// copy. Values of any other type are ignored.
func (c cutConfig) with(o any) cutConfig {
	switch v := o.(type) {
	case CutDelimiter:
		c.delimiter = v
	case CutFieldsOpt:
		c.fields = toPositions(v)
	case CutBytes:
		c.bytesSpec = positionSpec(v)
	case CutChars:
		c.charsSpec = positionSpec(v)
	case cutComplementFlag:
		c.complementEnabled = v
	}
	return c
}

// toPositions widens the 1-based field indices into positions.
func toPositions(fields CutFieldsOpt) []position {
	out := make([]position, len(fields))
	for i, f := range fields {
		out[i] = position(f)
	}
	return out
}

// parseSpec parses a spec like "1-3,5,7-" into a selection. Unparseable or
// empty parts are skipped; a spec that yields no intervals is left empty, which
// (without complement) emits nothing — matching the silent-skip behaviour of
// the original byte/char modes.
func parseSpec(spec positionSpec, isComplement cutComplementFlag) selection {
	var ivs []interval
	for _, part := range strings.Split(string(spec), ",") {
		if iv, ok := parseInterval(specPart(strings.TrimSpace(part))); ok {
			ivs = append(ivs, iv)
		}
	}
	return selection{intervals: ivs, isComplemented: bool(isComplement)}
}

// specPart is one comma-free element of a positionSpec: "N", "N-M", "N-", or
// "-M", or one of the two sides of such a range.
type specPart string

// parseInterval parses one comma-free spec part: "N", "N-M", "N-", or "-M".
func parseInterval(part specPart) (interval, bool) {
	dash := strings.Index(string(part), "-")
	if dash < 0 {
		return singleInterval(part)
	}
	return rangeInterval(part[:dash], part[dash+1:])
}

// singleInterval parses a bare "N" into the interval [N, N].
func singleInterval(s specPart) (interval, bool) {
	n, ok := parsePosition(s)
	if !ok {
		return interval{}, false
	}
	return interval{lo: n, hi: n}, true
}

// rangeInterval parses the two sides of a "lo-hi" spec, where either side may be
// empty: "N-" runs to the end and "-M" starts at position 1.
func rangeInterval(loStr, hiStr specPart) (interval, bool) {
	lo, hi := position(1), unbounded
	if loStr != "" {
		parsed, ok := parsePosition(loStr)
		if !ok {
			return interval{}, false
		}
		lo = parsed
	}
	if hiStr != "" {
		parsed, ok := parsePosition(hiStr)
		if !ok {
			return interval{}, false
		}
		hi = parsed
	}
	return interval{lo: lo, hi: hi}, true
}

// parsePosition parses a positive 1-based position.
func parsePosition(s specPart) (position, bool) {
	n, err := strconv.Atoi(string(s))
	if err != nil {
		return 0, false
	}
	return position(n), true
}

// itemSelector renders the selected items of a line given a selection.
type itemSelector func([]byte, selection) []byte

// positionCommand builds a Map command that applies sel via fn to each line.
func positionCommand(sel selection, fn itemSelector) gloo.Command[[]byte, []byte] {
	return patterns.Map(func(line []byte) ([]byte, error) {
		return fn(line, sel), nil
	})
}

// selectBytes keeps the selected bytes of line, in input order.
func selectBytes(line []byte, sel selection) []byte {
	out := make([]byte, 0, len(line))
	for i, b := range line {
		if sel.selected(position(i + 1)) {
			out = append(out, b)
		}
	}
	return out
}

// selectChars keeps the selected runes of line, in input order.
func selectChars(line []byte, sel selection) []byte {
	runes := []rune(string(line))
	out := make([]rune, 0, len(runes))
	for i, r := range runes {
		if sel.selected(position(i + 1)) {
			out = append(out, r)
		}
	}
	return []byte(string(out))
}

// fieldsRequested reports whether the -f flag named any fields; without a
// request the field mode passes every line through unchanged, like GNU cut.
type fieldsRequested bool

// fieldsCommand builds the field-selection (-f/-d) mode command.
func fieldsCommand(cfg cutConfig) gloo.Command[[]byte, []byte] {
	delim := []byte(delimiterOrTab(cfg.delimiter))
	sel := fieldSelection(cfg.fields, cfg.complementEnabled)
	hasFields := fieldsRequested(len(cfg.fields) > 0)
	return patterns.Map(func(line []byte) ([]byte, error) {
		return cutFields(line, delim, sel, hasFields), nil
	})
}

// delimiterOrTab defaults an empty delimiter to a tab, like GNU cut.
func delimiterOrTab(d CutDelimiter) string {
	if d == "" {
		return "\t"
	}
	return string(d)
}

// fieldSelection turns the requested 1-based field indices into a selection of
// single-position intervals.
func fieldSelection(fields []position, isComplement cutComplementFlag) selection {
	ivs := make([]interval, len(fields))
	for i, f := range fields {
		ivs[i] = interval{lo: f, hi: f}
	}
	return selection{intervals: ivs, isComplemented: bool(isComplement)}
}

// cutFields selects fields from one line. With no fields requested, or when the
// line contains no delimiter, the line passes through unchanged. Selected fields
// are emitted in input order (cut semantics), joined by the delimiter.
func cutFields(line, delim []byte, sel selection, hasFields fieldsRequested) []byte {
	if !bool(hasFields) {
		return line
	}
	parts := bytes.Split(line, delim)
	if len(parts) == 1 {
		return line
	}
	selected := make([][]byte, 0, len(parts))
	for i, part := range parts {
		if sel.selected(position(i + 1)) {
			selected = append(selected, part)
		}
	}
	return bytes.Join(selected, delim)
}
