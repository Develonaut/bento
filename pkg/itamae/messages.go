package itamae

import (
	"math/rand"
)

// logMessage provides standardized log messages with optional emojis.
// This ensures all log output uses approved emojis from .claude/EMOJIS.md
// and fun terminology from .claude/STATUS_WORDS.md.
type logMessage struct {
	emoji string
	text  string
}

// format returns the formatted log message with emoji if present.
func (m logMessage) format() string {
	if m.emoji != "" {
		return m.emoji + " " + m.text
	}
	return m.text
}

// sushiEmojis contains approved sushi-themed emojis for logs.
// Source: .claude/EMOJIS.md (synchronized with pkg/miso/sushi.go)
var sushiEmojis = []string{
	"🍣", "🍙", "🥢", "🍥", "🍱", "🍜", "🍡", "🍢",
	"🦐", "🦑", "🐟", "🍤", "🥟", "🥡", "🍶", "🍵", "🥠", "🧋",
}

// errorEmojis contains approved error emojis for failed operations.
var errorEmojis = []string{
	"👹", "👺", "💀", "☠️", "💥", "🔥", "⚠️", "❌", "🚫", "🤢",
}

// randomSushi returns a random sushi emoji from the approved list.
func randomSushi() string {
	return sushiEmojis[rand.Intn(len(sushiEmojis))]
}

// randomErrorEmoji returns a random error emoji from the approved list.
func randomErrorEmoji() string {
	return errorEmojis[rand.Intn(len(errorEmojis))]
}

// msgBentoStarted creates a message for bento execution start.
func msgBentoStarted() logMessage {
	return logMessage{
		emoji: randomSushi(),
		text:  "Starting bento execution",
	}
}

// msgBentoCompleted creates a message for bento execution completion.
func msgBentoCompleted() logMessage {
	return logMessage{
		emoji: randomSushi(),
		text:  "Bento execution completed",
	}
}

// msgBentoFailed creates a message for bento execution failure.
func msgBentoFailed() logMessage {
	return logMessage{
		emoji: randomErrorEmoji(),
		text:  "Bento execution failed",
	}
}

// msgNetaStarted creates a message for neta execution start.
func msgNetaStarted() logMessage {
	return logMessage{
		emoji: "", // No emoji for individual neta logs
		text:  "Executing neta",
	}
}

// msgNetaCompleted creates a message for neta execution completion.
func msgNetaCompleted() logMessage {
	return logMessage{
		emoji: "", // No emoji for individual neta logs
		text:  "Neta completed",
	}
}

// msgGroupStarted creates a message for group execution start.
func msgGroupStarted() logMessage {
	return logMessage{
		emoji: "", // No emoji for group logs
		text:  "│ ┌─ GROUP START",
	}
}

// msgGroupCompleted creates a message for group execution completion.
func msgGroupCompleted() logMessage {
	return logMessage{
		emoji: "", // No emoji for group logs
		text:  "│ └─ GROUP END",
	}
}

// msgLoopStarted creates a message for loop execution start.
func msgLoopStarted(mode string) logMessage {
	return logMessage{
		emoji: "", // No emoji for loop logs
		text:  "│   ┌─ LOOP START (" + mode + ")",
	}
}

// msgLoopCompleted creates a message for loop execution completion.
func msgLoopCompleted(mode string) logMessage {
	return logMessage{
		emoji: "", // No emoji for loop logs
		text:  "│   └─ LOOP END (" + mode + ")",
	}
}
