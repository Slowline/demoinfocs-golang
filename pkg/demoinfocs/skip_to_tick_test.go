package demoinfocs_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	demoinfocs "github.com/markus-wa/demoinfocs-golang/v5/pkg/demoinfocs"
	events "github.com/markus-wa/demoinfocs-golang/v5/pkg/demoinfocs/events"
)

func TestSkipToTick(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test due to -short flag")
	}

	t.Run("parse from beginning", func(t *testing.T) {
		f, err := os.Open(s2DemPath)
		assert.NoError(t, err)
		defer f.Close()

		p := demoinfocs.NewParser(f)
		defer p.Close()

		var killCount int
		var minTick int = 999999999
		var maxTick int = 0

		p.RegisterEventHandler(func(e events.Kill) {
			killCount++
			tick := p.GameState().IngameTick()
			if tick < minTick {
				minTick = tick
			}
			if tick > maxTick {
				maxTick = tick
			}
		})

		err = p.ParseToEnd()
		assert.NoError(t, err)
		assert.Greater(t, killCount, 0, "Should have recorded kills")
		t.Logf("Full parse: %d kills from tick %d to %d", killCount, minTick, maxTick)
	})

	t.Run("skip to tick 50000", func(t *testing.T) {
		f, err := os.Open(s2DemPath)
		assert.NoError(t, err)
		defer f.Close()

		config := demoinfocs.ParserConfig{
			MsgQueueBufferSize: 1000,
			SkipToTick:         50000,
		}

		p := demoinfocs.NewParserWithConfig(f, config)
		defer p.Close()

		var killCount int
		var minTick int = 999999999
		var maxTick int = 0

		p.RegisterEventHandler(func(e events.Kill) {
			killCount++
			tick := p.GameState().IngameTick()
			if tick < minTick {
				minTick = tick
			}
			if tick > maxTick {
				maxTick = tick
			}
		})

		err = p.ParseToEnd()
		assert.NoError(t, err)

		if killCount > 0 {
			assert.GreaterOrEqual(t, minTick, 50000, "First kill should be at or after tick 50000")
			t.Logf("Skipped parse: %d kills from tick %d to %d", killCount, minTick, maxTick)
		} else {
			t.Log("No kills after tick 50000")
		}
	})

	t.Run("skip to tick 10000", func(t *testing.T) {
		f, err := os.Open(s2DemPath)
		assert.NoError(t, err)
		defer f.Close()

		config := demoinfocs.ParserConfig{
			MsgQueueBufferSize: 1000,
			SkipToTick:         10000,
		}

		p := demoinfocs.NewParserWithConfig(f, config)
		defer p.Close()

		var eventCount int
		var minTick int = 999999999

		p.RegisterEventHandler(func(e events.Kill) {
			eventCount++
			tick := p.GameState().IngameTick()
			if tick < minTick {
				minTick = tick
			}
		})

		p.RegisterEventHandler(func(e events.RoundStart) {
			eventCount++
			tick := p.GameState().IngameTick()
			if tick < minTick {
				minTick = tick
			}
		})

		err = p.ParseToEnd()
		assert.NoError(t, err)

		if eventCount > 0 {
			assert.GreaterOrEqual(t, minTick, 10000, "First event should be at or after tick 10000")
			t.Logf("Skipped parse: %d events starting from tick %d", eventCount, minTick)
		} else {
			t.Log("No events after tick 10000")
		}
	})
}
