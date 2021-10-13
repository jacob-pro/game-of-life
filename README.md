# Game of Life

Bristol COMS20001_2019 Concurrent Computing (Yr 2), 
Coursework 1: Game of Life \
Pair project: [@jacob-pro](https://github.com/jacob-pro) [@michaelwagstaff](https://github.com/michaelwagstaff)

Confirmed mark: 82

Various concurrent implementations of Conway's Game of Life. 

The primary stage of the assignment was to implement the game using Go channels in a halo-exchange system, 
the secondary stage allowed us to choose our own method with the goal of making it as fast as possible. We
took an unconventional but very fast approach - writing the core logic in Rust and linking it as a static library!

## Building

To run in Goland/Jetbrains add a `go build` build configuration in the `src` directory. 
Make sure to enable `run.processes.with.pty` in the registry so that termbox will work

- Install Rust https://www.rust-lang.org/tools/install

To use cgo on Windows requires the Rust library be compiled with the GCC/GNU toolchain. 
- `rustup toolchain install stable-gnu`
- In an MINGW64 terminal `pacman -S mingw-w64-x86_64-gcc`
- Install GNU make `pacman -S make`

Build: `cd src && make`

Running tests:
- `go test -args -i "halo"`
- `go test -bench . -args -i "halo"`
