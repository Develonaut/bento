package itamae

import (
	"context"
	"fmt"
	"time"

	"github.com/Develonaut/bento/pkg/neta"
)

// executeNode executes a single node (handles all node types).
func (i *Itamae) executeNode(
	ctx context.Context,
	def *neta.Definition,
	execCtx *executionContext,
	result *Result,
) error {
	// Check context cancellation
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// Handle different node types
	switch def.Type {
	case "group":
		return i.executeGroup(ctx, def, execCtx, result)
	case "loop":
		return i.executeLoop(ctx, def, execCtx, result)
	case "parallel":
		return i.executeParallel(ctx, def, execCtx, result)
	default:
		return i.executeSingle(ctx, def, execCtx, result)
	}
}

// executeSingle executes a single (non-group) neta.
func (i *Itamae) executeSingle(
	ctx context.Context,
	def *neta.Definition,
	execCtx *executionContext,
	result *Result,
) error {
	// Progress callback: starting
	i.notifyProgress(def.ID, "starting")

	// Send messenger event: node started
	if i.messenger != nil {
		i.messenger.SendNodeStarted(def.ID, def.Name, def.Type)
	}

	if i.logger != nil {
		msg := msgNetaStarted()
		// Use debug level for individual neta execution to reduce log verbosity
		// This keeps logs clean while still providing detail with --verbose
		i.logger.Debug(msg.format(),
			"neta_id", def.ID,
			"neta_type", def.Type)
	}

	// Get neta implementation from pantry
	netaImpl, err := i.pantry.GetNew(def.Type)
	if err != nil {
		return newNodeError(def.ID, def.Type, "get neta", err)
	}

	// Prepare parameters with execution context
	// Resolve templates in parameters using current execution context
	params := make(map[string]interface{})
	for k, v := range def.Parameters {
		params[k] = execCtx.resolveValue(v)
	}

	// Add execution context for template resolution (for neta that need it)
	params["_context"] = execCtx.toMap()

	// Add streaming callback for shell-command neta (Phase 8.5)
	// This enables real-time output from long-running processes like Blender
	params["_onOutput"] = func(line string) {
		if i.logger != nil {
			i.logger.Stream(line)
		}
	}

	// Track execution time for messenger
	start := time.Now()

	// Execute neta
	output, err := netaImpl.Execute(ctx, params)

	duration := time.Since(start)

	// Send messenger event: node completed (with duration and error if any)
	if i.messenger != nil {
		i.messenger.SendNodeCompleted(def.ID, duration, err)
	}

	if err != nil {
		return newNodeError(def.ID, def.Type, "execute", err)
	}

	// Store output
	execCtx.set(def.ID, output)
	result.NodeOutputs[def.ID] = output
	result.NodesExecuted++

	// Progress callback: completed
	i.notifyProgress(def.ID, "completed")

	if i.logger != nil {
		msg := msgNetaCompleted()
		// Use debug level for individual neta completion to reduce log verbosity
		i.logger.Debug(msg.format(),
			"neta_id", def.ID,
			"neta_type", def.Type)
	}

	return nil
}

// executeGroup executes a group neta (container with child nodes).
func (i *Itamae) executeGroup(
	ctx context.Context,
	def *neta.Definition,
	execCtx *executionContext,
	result *Result,
) error {
	i.notifyProgress(def.ID, "starting")

	// Send messenger event: node started
	if i.messenger != nil {
		i.messenger.SendNodeStarted(def.ID, def.Name, def.Type)
	}

	if i.logger != nil {
		msg := msgGroupStarted()
		i.logger.Info(msg.format(),
			"group_id", def.ID,
			"group_name", def.Name,
			"child_count", len(def.Nodes))
	}

	// Handle empty group
	if len(def.Nodes) == 0 {
		i.notifyProgress(def.ID, "completed")
		// Send messenger event: node completed
		if i.messenger != nil {
			i.messenger.SendNodeCompleted(def.ID, 0, nil)
		}
		return nil
	}

	// Track execution time for messenger
	start := time.Now()

	// Build execution graph
	g, err := buildGraph(def)
	if err != nil {
		if i.messenger != nil {
			i.messenger.SendNodeCompleted(def.ID, time.Since(start), err)
		}
		return newNodeError(def.ID, "group", "build graph", err)
	}

	// Check for cycles
	if g.hasCycle() {
		err := newNodeError(def.ID, "group", "validate",
			fmt.Errorf("circular dependency detected"))
		if i.messenger != nil {
			i.messenger.SendNodeCompleted(def.ID, time.Since(start), err)
		}
		return err
	}

	// Execute nodes in topological order
	if err := i.executeGraph(ctx, g, execCtx, result); err != nil {
		if i.messenger != nil {
			i.messenger.SendNodeCompleted(def.ID, time.Since(start), err)
		}
		return err
	}

	duration := time.Since(start)

	i.notifyProgress(def.ID, "completed")
	// Note: Group execution is tracked by child nodes, not the group itself

	// Send messenger event: node completed
	if i.messenger != nil {
		i.messenger.SendNodeCompleted(def.ID, duration, nil)
	}

	if i.logger != nil {
		msg := msgGroupCompleted()
		i.logger.Info(msg.format(),
			"group_id", def.ID,
			"group_name", def.Name,
			"child_count", len(def.Nodes))
	}

	return nil
}

// executeGraph executes all nodes in a graph in topological order.
func (i *Itamae) executeGraph(
	ctx context.Context,
	g *graph,
	execCtx *executionContext,
	result *Result,
) error {
	executed := make(map[string]bool)
	queue := g.getStartNodes()

	for len(queue) > 0 {
		// Check context cancellation
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Get next node to execute
		node := queue[0]
		queue = queue[1:]

		// Skip if already executed
		if executed[node.ID] {
			continue
		}

		// Execute node
		if err := i.executeNode(ctx, node, execCtx, result); err != nil {
			return err
		}

		executed[node.ID] = true
		g.markExecuted(node.ID)

		// Add ready children to queue
		for _, target := range g.getTargets(node.ID) {
			if g.isReady(target.ID) && !executed[target.ID] {
				queue = append(queue, target)
			}
		}
	}

	return nil
}

// notifyProgress calls the progress callback if set.
func (i *Itamae) notifyProgress(nodeID string, status string) {
	if i.onProgress != nil {
		i.onProgress(nodeID, status)
	}
}
