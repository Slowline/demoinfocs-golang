package main

import (
	"fmt"
	"log"
	"os"

	ex "github.com/markus-wa/demoinfocs-golang/v5/examples"
	demoinfocs "github.com/markus-wa/demoinfocs-golang/v5/pkg/demoinfocs"
	events "github.com/markus-wa/demoinfocs-golang/v5/pkg/demoinfocs/events"
)

// Run like this: go run skip_to_tick.go -demo /path/to/demo.dem
//
// This example demonstrates how to use SkipToTick to start parsing from a specific tick.
// Events before the target tick will not be dispatched to handlers, improving performance
// when you only need to analyze a specific portion of a demo.
func main() {
	f, err := os.Open(ex.DemoPathFromArgs())
	if err != nil {
		log.Panic("failed to open demo file: ", err)
	}
	defer f.Close()

	// Parse the demo from tick 10000 onwards
	// All events before tick 10000 will be skipped
	cfg := demoinfocs.ParserConfig{
		MsgQueueBufferSize: 100000,
		SkipToTick:         10000, // Start processing from tick 10000
	}

	p := demoinfocs.NewParserWithConfig(f, cfg)
	defer p.Close()

	var eventCount int
	var firstTick int = -1

	// Register handler for kills
	p.RegisterEventHandler(func(e events.Kill) {
		if firstTick == -1 {
			firstTick = p.GameState().IngameTick()
			fmt.Printf("First event received at tick %d (target was %d)\n", firstTick, cfg.SkipToTick)
		}
		eventCount++
		fmt.Printf("[Tick %d] %s killed %s with %s\n",
			p.GameState().IngameTick(),
			e.Killer.Name,
			e.Victim.Name,
			e.Weapon.String())
	})

	// Register handler for round starts
	p.RegisterEventHandler(func(e events.RoundStart) {
		if firstTick == -1 {
			firstTick = p.GameState().IngameTick()
			fmt.Printf("First event received at tick %d (target was %d)\n", firstTick, cfg.SkipToTick)
		}
		eventCount++
		fmt.Printf("[Tick %d] Round started\n", p.GameState().IngameTick())
	})

	// Parse to end
	err = p.ParseToEnd()
	if err != nil {
		log.Panic("failed to parse demo: ", err)
	}

	fmt.Printf("\nTotal events processed: %d\n", eventCount)
	if firstTick != -1 {
		fmt.Printf("First event was at tick %d (skipped %d ticks)\n", firstTick, firstTick)
	} else {
		fmt.Println("No events were found after the skip tick")
	}
}
