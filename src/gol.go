package main

import (
	"fmt"
	"sync"
	"time"
)

type distributorState struct {
	currentTurn  int
	currentWorld World
	locker       sync.Locker
}

// distributor divides the work between workers and interacts with other goroutines.
// d.keyChan is nil when launched by the test framework
func distributor(p golParams, d distributorChans, alive chan []cell) {

	// The state lock is used to synchronise the Distributor, Keyboard, and Ticker threads.
	// If no user is connected (e.g. testing/benchmark) then we do not need real synchronisation
	// Because only the Distributor thread will be running
	var lock sync.Locker
	if d.keyChan != nil {
		lock = &sync.Mutex{}
	} else {
		lock = NopLocker{}
	}

	state := distributorState{
		currentTurn:  0,
		currentWorld: LoadWorldFromPgm(p.imageHeight, p.imageWidth, d),
		locker:       lock,
	}

	if d.keyChan != nil {
		go handleKeyboard(&state, d)
		go ticker(&state)
	}

	f := getImplementationFunction(p.implementationName)

	// Calculate the new state of Game of Life after the given number of turns.
	turnLocal := 0
	for turnLocal < p.turns {
		//Clone the World into local memory
		state.locker.Lock()
		worldLocal := state.currentWorld.Clone()
		state.locker.Unlock()

		//Perform the computation
		f(&worldLocal, p.threads)
		turnLocal++

		//Update the state with new world
		state.locker.Lock()
		state.currentWorld = worldLocal
		state.currentTurn = turnLocal
		state.locker.Unlock()
	}

	// Make sure that the Io has finished any output before exiting (there may be a save in progress).
	d.io.command <- ioCheckIdle
	<-d.io.idle

	// Return the coordinates of cells that are still alive.
	state.locker.Lock()
	alive <- state.currentWorld.CalculateAlive()
}

func handleKeyboard(state *distributorState, d distributorChans) {
	for {
		keyPress := <-d.keyChan
		switch keyPress {
		case 's':
			state.locker.Lock()
			state.currentWorld.SaveToPgm(d, state.currentTurn)
			state.locker.Unlock()
		case 'p':
			state.locker.Lock()
			fmt.Printf("Paused. Press p to continue...\n")
			// Wait for p to be pressed again
			// Whilst the mutex remains locked the GoL and Ticker threads will be stuck
			for {
				keyPress := <-d.keyChan
				if keyPress == 'p' {
					break
				}
			}
			fmt.Printf("Continuing...\n")
			state.locker.Unlock()
		case 'q':
			state.locker.Lock()
			state.currentWorld.SaveToPgm(d, state.currentTurn)
			// Make sure that the Io has finished any output before exiting (so that the save completes)
			d.io.command <- ioCheckIdle
			<-d.io.idle
			exit()
		}
	}
}

func ticker(state *distributorState) {
	ticker := time.NewTicker(2 * time.Second)
	for {
		<-ticker.C
		state.locker.Lock()
		world := state.currentWorld.Clone()
		turn := state.currentTurn
		state.locker.Unlock()
		fmt.Printf("On turn: %d there are %d alive\n", turn, len(world.CalculateAlive()))
	}
}

func getImplementationFunction(name string) func(world *World, threads int) {
	if name == "" {
		return ImplementationDefault.function()
	} else {
		impl, err := implementationFromName(name)
		check(err)
		return impl.function()
	}
}
