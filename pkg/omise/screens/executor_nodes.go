package screens

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"bento/pkg/neta"
	"bento/pkg/omise/components"
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

	// Update sequence display with initial steps
	e.sequence = e.sequence.SetSteps(e.convertNodesToSteps())

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

	// Update sequence display
	e.sequence = e.sequence.SetSteps(e.convertNodesToSteps())

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

	// Update sequence display
	e.sequence = e.sequence.SetSteps(e.convertNodesToSteps())

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

// convertNodesToSteps converts NodeStates to Sequence Steps
func (e Executor) convertNodesToSteps() []components.Step {
	steps := make([]components.Step, len(e.nodeStates))
	for i, node := range e.nodeStates {
		steps[i] = components.Step{
			Name:     node.name,
			Type:     node.nodeType,
			Status:   convertNodeStatusToStepStatus(node.status),
			Duration: node.duration,
			Depth:    node.depth,
		}
	}
	return steps
}

// convertNodeStatusToStepStatus converts NodeStatus to StepStatus
func convertNodeStatusToStepStatus(status NodeStatus) components.StepStatus {
	switch status {
	case NodePending:
		return components.StepPending
	case NodeRunning:
		return components.StepRunning
	case NodeCompleted:
		return components.StepCompleted
	case NodeFailed:
		return components.StepFailed
	default:
		return components.StepPending
	}
}
