package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/dottedmag/must"
)

func check(t *testing.T, src, expected string) {
	t.Run(src, func(t *testing.T) {
		rdr := strings.NewReader(src + "\n")
		wr := bytes.NewBuffer(nil)

		must.OK(formatDump("/home/test", wr, rdr))

		assert.Equal(t, expected, strings.TrimSuffix(wr.String(), "\n"))
	})
}

func TestFormat(t *testing.T) {
	// Goroutine headers are reformatted
	check(t,
		"goroutine 97 gp=0x14000584540 m=nil [select]:",
		"97   [select] gp=0x14000584540 m=nil:")
	// Go stdlib source locations are indented and shortened
	check(t,
		"\t/opt/homebrew/Cellar/go/1.22.2/libexec/src/runtime/proc.go:1 +0x2 fp=0x3 sp=0x4 pc=0x5",
		"    -GO 1.22.2-/runtime/proc.go:1 +0x2 fp=0x3 sp=0x4 pc=0x5")

	// Go modules cache locations are indented and shortened
	check(t,
		"\t/home/test/go/pkg/mod/foo.com/bar/baz@v1.2.3/retry.go:6 +0x7 fp=0x8 sp=0x9 pc=0xa",
		"    +CACHE+/foo.com/bar/baz@v1.2.3/retry.go:6 +0x7 fp=0x8 sp=0x9 pc=0xa")

	// source locations in $HOME are indented and shortened
	check(t,
		"\t/home/test/proj/file.go:9 +0xa fp=0xb sp=0xc pc=0xd",
		"    <HOME>/proj/file.go:9 +0xa fp=0xb sp=0xc pc=0xd")

	// all other source locations are indented
	check(t,
		"\t/opt/homebrew/Cellar/nogo/func.go:1 +0x2 fp=0x3 sp=0x4 pc=0x5",
		"    /opt/homebrew/Cellar/nogo/func.go:1 +0x2 fp=0x3 sp=0x4 pc=0x5")
	check(t,
		"\t/tmp/a.go:1 +0x2 fp=0x3 sp=0x4 pc=0x5",
		"    /tmp/a.go:1 +0x2 fp=0x3 sp=0x4 pc=0x5")
	check(t,
		"\ta.go:1 +0x2 fp=0x3 sp=0x4 pc=0x5",
		"    a.go:1 +0x2 fp=0x3 sp=0x4 pc=0x5")

	// stdlib function names are indented and have a prefix
	check(t,
		"database/sql.OpenDB.gowrap1()",
		"  -database/sql.OpenDB.gowrap1()")
	// other function names are indented
	check(t,
		"foo.bar/baz.Uuu({})",
		"  foo.bar/baz.Uuu({})")
	check(t,
		"main.main()",
		"  main.main()")
	// Creation records are indented and have a prefix
	check(t,
		"created by testing.(*T).Run in goroutine 66",
		"  .created by testing.(*T).Run in goroutine 66")

	// signal names, PC and register values are untouched
	check(t,
		"SIGQUIT: quit",
		"SIGQUIT: quit")
	check(t,
		"PC=0x1 m=0 sigcode=0",
		"PC=0x1 m=0 sigcode=0")
	check(t,
		"r0\t0x0",
		"r0\t0x0")
	check(t,
		"fault\t0x1",
		"fault\t0x1")
}
