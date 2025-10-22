package itamae

import (
	"fmt"
	"hash/fnv"
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

// statusWordsRunning contains fun status words for running nodes.
// Source: .claude/STATUS_WORDS.md (synchronized with pkg/miso/sushi.go)
var statusWordsRunning = []string{
	"Tasting",
	"Sampling",
	"Trying",
	"Enjoying",
	"Devouring",
	"Nibbling",
	"Savoring",
	"Testing",
}

// statusWordsCompleted contains fun status words for completed nodes.
// Source: .claude/STATUS_WORDS.md (synchronized with pkg/miso/sushi.go)
var statusWordsCompleted = []string{
	"Savored",
	"Devoured",
	"Enjoyed",
	"Relished",
	"Finished",
	"Consumed",
	"Completed",
	"Perfected",
}

// randomSushi returns a random sushi emoji from the approved list.
func randomSushi() string {
	return sushiEmojis[rand.Intn(len(sushiEmojis))]
}

// randomErrorEmoji returns a random error emoji from the approved list.
func randomErrorEmoji() string {
	return errorEmojis[rand.Intn(len(errorEmojis))]
}

// getStatusWord returns a fun varied status word based on node name.
// Uses deterministic hash to ensure same node gets same word.
func getStatusWord(name string, isRunning bool) string {
	h := fnv.New32a()
	h.Write([]byte(name))
	hash := h.Sum32()

	if isRunning {
		return statusWordsRunning[hash%uint32(len(statusWordsRunning))]
	}
	return statusWordsCompleted[hash%uint32(len(statusWordsCompleted))]
}

// msgBentoStarted creates a message for bento execution start.
// Format matches CLI output: "🍱 Running Bento: [name]"
func msgBentoStarted(name string) logMessage {
	return logMessage{
		emoji: "🍱",
		text:  "Running Bento: " + name,
	}
}

// msgBentoCompleted creates a message for bento execution completion.
// Format matches CLI output: "🍥 Delicious! Bento executed successfully in [duration]"
func msgBentoCompleted(duration string) logMessage {
	return logMessage{
		emoji: randomSushi(),
		text:  "Delicious! Bento executed successfully in " + duration,
	}
}

// msgBentoFailed creates a message for bento execution failure.
// Format: "❌ Failed! Bento execution failed in [duration]"
func msgBentoFailed(duration string) logMessage {
	return logMessage{
		emoji: randomErrorEmoji(),
		text:  "Failed! Bento execution failed in " + duration,
	}
}

// msgNetaStarted creates a message for neta execution start.
func msgNetaStarted() logMessage {
	return logMessage{
		emoji: "", // No emoji for individual neta logs
		text:  "Executing neta",
	}
}

// msgGroupStarted creates a message for group execution start.
// Format: "│ ┌─ Tasting NETA:group name"
func msgGroupStarted(depth int, name string) logMessage {
	indent := getIndent(depth)
	statusWord := getStatusWord(name, true)
	return logMessage{
		emoji: "",
		text:  indent + "  ┌─ " + statusWord + " NETA:group " + name,
	}
}

// msgGroupCompleted creates a message for group execution completion.
// Format: "│ └─ Finished NETA:group name (2ms)"
func msgGroupCompleted(depth int, name, duration string) logMessage {
	indent := getIndent(depth)
	statusWord := getStatusWord(name, false)
	return logMessage{
		emoji: "",
		text:  indent + "  └─ " + statusWord + " NETA:group " + name + " (" + duration + ")",
	}
}

// msgLoopStarted creates a message for loop execution start.
// Format: "│  │  ┌─ Sampling NETA:loop name"
func msgLoopStarted(depth int, name string) logMessage {
	indent := getIndent(depth)
	statusWord := getStatusWord(name, true)
	return logMessage{
		emoji: "",
		text:  indent + "  ┌─ " + statusWord + " NETA:loop " + name,
	}
}

// msgLoopCompleted creates a message for loop execution completion.
// Format: "│  │  └─ Perfected NETA:loop name (2ms, 75%)"
func msgLoopCompleted(depth int, name, duration string, progressPct int) logMessage {
	indent := getIndent(depth)
	statusWord := getStatusWord(name, false)
	return logMessage{
		emoji: "",
		text:  fmt.Sprintf("%s  └─ %s NETA:loop %s (%s, %d%%)", indent, statusWord, name, duration, progressPct),
	}
}

// msgChildNodeStarted creates a message for child node execution start.
// Depth indicates nesting level: 0=root, 1=in group, 2=in loop, etc.
// Format: "│  │  ┌─ Tasting NETA:type name"
func msgChildNodeStarted(depth int, nodeType, name string) logMessage {
	indent := getIndent(depth)
	statusWord := getStatusWord(name, true)
	return logMessage{
		emoji: "",
		text:  indent + "  ┌─ " + statusWord + " NETA:" + nodeType + " " + name,
	}
}

// msgChildNodeCompleted creates a message for child node execution completion.
// Format: "│  │  └─ Devoured NETA:type name (2ms, 10%)"
func msgChildNodeCompleted(depth int, nodeType, name, duration string, progressPct int) logMessage {
	indent := getIndent(depth)
	statusWord := getStatusWord(name, false)
	return logMessage{
		emoji: "",
		text:  fmt.Sprintf("%s  └─ %s NETA:%s %s (%s, %d%%)", indent, statusWord, nodeType, name, duration, progressPct),
	}
}

// getIndent returns the indentation string based on depth.
// Preserves parent dividers at each nesting level:
// Depth 0 = "│"
// Depth 1 = "│  │"
// Depth 2 = "│  │  │"
func getIndent(depth int) string {
	if depth == 0 {
		return "│"
	}
	base := "│"
	for i := 0; i < depth; i++ {
		base += "  │"
	}
	return base
}
