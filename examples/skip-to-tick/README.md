# Skip to Tick Example

This example demonstrates how to use the `SkipToTick` feature to start parsing a demo from a specific tick, rather than from the beginning.

## Overview

When you only need to analyze a specific portion of a demo (for example, the last few rounds), you can use `SkipToTick` to improve performance by skipping all events before that tick. The parser will still process the demo data to maintain correct game state, but events won't be dispatched to your handlers until the target tick is reached.

## Usage

```bash
go run skip_to_tick.go -demo /path/to/demo.dem
```

## How It Works

The example configures the parser with a target tick:

```go
cfg := demoinfocs.ParserConfig{
    MsgQueueBufferSize: 100000,
    SkipToTick:         10000,  // Start from tick 10000
}

p := demoinfocs.NewParserWithConfig(f, cfg)
```

With this configuration:
- All events before tick 10000 will be suppressed (not dispatched to handlers)
- Game state is still updated correctly throughout the demo
- Once tick 10000 is reached, events are dispatched normally
- This provides a performance benefit when you don't need early-game data

## Benefits

- **Performance**: Skip processing events you don't need
- **Targeted Analysis**: Focus on specific rounds or time periods
- **Memory Efficiency**: Event handlers aren't invoked for skipped ticks
- **Correct State**: Game state (player positions, scores, etc.) is still accurate at the target tick

## Use Cases

- Analyzing specific rounds (e.g., round 20 onwards)
- Processing only overtime rounds
- Debugging issues that occur late in a match
- Creating highlights from specific time periods
- Reducing processing time for long demos when only late-game data is needed
