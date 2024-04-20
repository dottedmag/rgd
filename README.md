# rgd (readable goroutines dump)

This tool refines Go goroutines dumps to improve their readability:
```
goroutine 97 gp=0x14000584540 m=nil [select]:
runtime.gopark(0x14000185af0?, 0x2?, 0x8?, 0x9?, 0x14000185a94?)
	/opt/homebrew/Cellar/go/1.22.2/libexec/src/runtime/proc.go:402 +0xc8 fp=0x14000185940 sp=0x14000185920 pc=0x10029c0d8
github.com/dottedmag/retry.Do({0x100a6fb68, 0x140004a8b40}, {0x5f5e100, 0x3ff8000000000000, 0x3fc0000000000000, 0x0, 0x6fc23ac00, 0x0, 0x0, 0x0, ...}, ...)
	/Users/dottedmag/go/pkg/mod/github.com/dottedmag/retry@v0.0.0-20240406161445-a48d7c84477a/retry.go:205 +0x43c fp=0x14000185b50 sp=0x14000185a50 pc=0x100828f5c
main.main()
	testmain.go:51 +0x16c fp=0x14000047f40 sp=0x14000047eb0 pc=0x100853a0c
testing.tRunner(0x140004b8680, 0x140006b6e40)
	/opt/homebrew/Cellar/go/1.22.2/libexec/src/testing/testing.go:1689 +0xec fp=0x14000185fb0 sp=0x14000185f60 pc=0x10036257c
created by testing.(*T).Run in goroutine 95
	/opt/homebrew/Cellar/go/1.22.2/libexec/src/testing/testing.go:1742 +0x318
```
becomes
```
97   [select] gp=0x14000584540 m=nil:
  -runtime.gopark(0x14000185af0?, 0x2?, 0x8?, 0x9?, 0x14000185a94?)
    -GO 1.22.2-/runtime/proc.go:402 +0xc8 fp=0x14000185940 sp=0x14000185920 pc=0x10029c0d8
  github.com/dottedmag/retry.Do({0x100a6fb68, 0x140004a8b40}, {0x5f5e100, 0x3ff8000000000000, 0x3fc0000000000000, 0x0, 0x6fc23ac00, 0x0, 0x0, 0x0, ...}, ...)
    +CACHE+/github.com/dottedmag/retry@v0.0.0-20240406161445-a48d7c84477a/retry.go:205 +0x43c fp=0x14000185b50 sp=0x14000185a50 pc=0x100828f5c
  main.main()
    testmain.go:51 +0x16c fp=0x14000047f40 sp=0x14000047eb0 pc=0x100853a0c
  -testing.tRunner(0x140004b8680, 0x140006b6e40)
    -GO 1.22.2-/testing/testing.go:1689 +0xec fp=0x14000185fb0 sp=0x14000185f60 pc=0x10036257c
  .created by testing.(*T).Run in goroutine 95
    -GO 1.22.2-/testing/testing.go:1742 +0x318
```

## Usage

```
rgd (<input file>|-)
```
