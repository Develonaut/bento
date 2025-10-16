package screens

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"bento/pkg/neta"
	"bento/pkg/omise/styles"
)

// buildPathForNode constructs node path for tracking
func buildPathForNode(basePath string, index int) string {
	if basePath == "" {
		return fmt.Sprintf("%d", index)
	}
	return fmt.Sprintf("%s.%d", basePath, index)
}

// parseDepth calculates nesting level from path
func parseDepth(path string) int {
	if path == "" {
		return 0
	}
	return strings.Count(path, ".") + 1
}

// flattenDefinition converts tree to flat list with paths
func flattenDefinition(def neta.Definition, basePath string) []NodeState {
	if !def.IsGroup() {
		return flattenSingleNode(def, basePath)
	}
	return flattenGroupNodes(def, basePath)
}

// flattenSingleNode creates state for a single node bento
func flattenSingleNode(def neta.Definition, basePath string) []NodeState {
	return []NodeState{{
		path:     basePath,
		name:     def.Name,
		nodeType: def.Type,
		status:   NodePending,
		depth:    0,
	}}
}

// flattenGroupNodes flattens all nodes in a group recursively
func flattenGroupNodes(def neta.Definition, basePath string) []NodeState {
	states := []NodeState{}

	for idx, child := range def.Nodes {
		path := buildPathForNode(basePath, idx)
		states = append(states, createNodeState(child, path))

		if child.IsGroup() {
			childStates := flattenDefinition(child, path)
			states = append(states, childStates...)
		}
	}

	return states
}

// createNodeState builds a NodeState from definition and path
func createNodeState(def neta.Definition, path string) NodeState {
	return NodeState{
		path:     path,
		name:     def.Name,
		nodeType: def.Type,
		status:   NodePending,
		depth:    parseDepth(path),
	}
}

// handleInitMsg initializes node states from definition
func (e Executor) handleInitMsg(msg ExecutionInitMsg) (Executor, tea.Cmd) {
	// Flatten definition to get all nodes
	e.nodeStates = flattenDefinition(msg.Definition, "")
	return e, nil
}

// handleNodeStarted updates node to running state
func (e Executor) handleNodeStarted(msg NodeStartedMsg) (Executor, tea.Cmd) {
	for i := range e.nodeStates {
		if e.nodeStates[i].path == msg.Path {
			e.nodeStates[i].status = NodeRunning
			e.nodeStates[i].startTime = time.Now()
			break
		}
	}
	return e, nil
}

// handleNodeCompleted updates node to completed/failed state
func (e Executor) handleNodeCompleted(msg NodeCompletedMsg) (Executor, tea.Cmd) {
	for i := range e.nodeStates {
		if e.nodeStates[i].path == msg.Path {
			e.nodeStates[i].duration = msg.Duration
			if msg.Error != nil {
				e.nodeStates[i].status = NodeFailed
			} else {
				e.nodeStates[i].status = NodeCompleted
			}
			break
		}
	}

	// Update progress based on completion
	e.progressPercent = e.calculateProgressFromNodes()

	return e, nil
}

// calculateProgressFromNodes returns completion percentage based on nodes
func (e Executor) calculateProgressFromNodes() float64 {
	if len(e.nodeStates) == 0 {
		return 0.0
	}

	completed := 0
	for _, node := range e.nodeStates {
		if node.status == NodeCompleted || node.status == NodeFailed {
			completed++
		}
	}

	return float64(completed) / float64(len(e.nodeStates))
}

// formatNodeLine renders single node with status icon and timing
func (e Executor) formatNodeLine(node NodeState) string {
	indent := strings.Repeat("  ", node.depth)
	icon := e.getNodeIcon(node.status)
	line := e.buildNodeLine(indent, icon, node)
	return e.styleNodeLine(line, node.status)
}

// buildNodeLine constructs the node line with icon and duration
func (e Executor) buildNodeLine(indent, icon string, node NodeState) string {
	line := fmt.Sprintf("%s%s %s", indent, icon, node.name)

	if node.status == NodeCompleted || node.status == NodeFailed {
		durationStr := node.duration.Round(time.Millisecond).String()
		line = fmt.Sprintf("%s (%s)", line, durationStr)
	}

	return line
}

// styleNodeLine applies styling based on node status
func (e Executor) styleNodeLine(line string, status NodeStatus) string {
	switch status {
	case NodeCompleted:
		return styles.SuccessStyle.Render(line)
	case NodeFailed:
		return styles.ErrorStyle.Render(line)
	case NodeRunning:
		return line
	default:
		return styles.Subtle.Render(line)
	}
}

// getNodeIcon returns emoji/character for node status
func (e Executor) getNodeIcon(status NodeStatus) string {
	switch status {
	case NodeRunning:
		return e.spinner.View() // Animated spinner
	case NodeCompleted:
		return emojiSuccess // ✓
	case NodeFailed:
		return emojiFailure // ✗
	default:
		return "•" // Pending
	}
}
