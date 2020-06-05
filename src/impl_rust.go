package main

// #cgo CFLAGS: -I../rust/include
// #cgo windows LDFLAGS: ../rust/target/release/gol_rust.lib -lws2_32 -luserenv
// #cgo linux darwin LDFLAGS: ../rust/target/release/libgol_rust.a -ldl
// #include <gol.h>
import "C"

type rust struct {
	gol    *C.GameOfLife
	height int
	width  int
}

// Stage 5 custom high performance implementation
func initRust(world world, threads int) implementation {

	gol := C.gol_init((*C.uchar)(&world.matrix[0]), C.int32_t(world.height), C.int32_t(world.width), C.int32_t(threads))

	return &rust{
		gol:    gol,
		height: world.height,
		width:  world.width,
	}
}

func (r *rust) nextTurn() {
	C.gol_next_turn(r.gol)
}

func (r *rust) getWorld() world {
	// Load the world into a slice
	b := make([]byte, r.width*r.height)
	C.gol_get_world(r.gol, (*C.uchar)(&b[0]))

	return world{
		width:  r.width,
		height: r.height,
		matrix: b,
	}
}

func (r *rust) close() {
	C.gol_destroy(r.gol)
}
