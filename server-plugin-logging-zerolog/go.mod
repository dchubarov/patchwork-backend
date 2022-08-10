module twowls.org/patchwork/plugin/logging/zerolog

go 1.19

require (
	github.com/rs/zerolog v1.27.0
	twowls.org/patchwork/commons v0.0.0-00010101000000-000000000000
)

require (
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8 // indirect
)

replace twowls.org/patchwork/commons => ../commons
