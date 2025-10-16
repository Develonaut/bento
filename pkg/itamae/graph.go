package itamae

import (
	"context"
	"fmt"
	"time"

	"bento/pkg/neta"
)

// executeGraph executes nodes in topological order using ExecutionGraphStore
func (i *Itamae) executeGraph(ctx context.Context, def neta.Definition) (neta.Result, error) {
	graph, nodeMap, err := i.prepareGraphExecution(def)
	if err != nil {
		return neta.Result{}, err
	}

	outputs, err := i.executeGraphLayers(ctx, graph, nodeMap, def.Edges)
	if err != nil {
		return neta.Result{}, err
	}

	return i.finalizeGraphExecution(outputs, graph)
}

// prepareGraphExecution analyzes graph and initializes execution state
func (i *Itamae) prepareGraphExecution(def neta.Definition) (*neta.ExecutionGraph, map[string]neta.Definition, error) {
	graph, err := neta.AnalyzeExecutionGraph(def)
	if err != nil {
		return nil, nil, fmt.Errorf("graph analysis failed: %w", err)
	}

	i.store.InitializeGraph(
		graph.Nodes,
		graph.Edges,
		graph.ExecutionOrder,
		graph.CriticalPath,
		graph.TotalWeight,
		graph.MaxParallelism,
	)

	nodeMap := buildNodeMap(def.Nodes)
	return graph, nodeMap, nil
}

// buildNodeMap creates lookup map from node ID to Definition
func buildNodeMap(nodes []neta.Definition) map[string]neta.Definition {
	nodeMap := make(map[string]neta.Definition)
	for _, node := range nodes {
		nodeMap[node.ID] = node
	}
	return nodeMap
}

// executeGraphLayers executes all layers in topological order
func (i *Itamae) executeGraphLayers(ctx context.Context, graph *neta.ExecutionGraph, nodeMap map[string]neta.Definition, edges []neta.NodeEdge) (map[string]neta.Result, error) {
	outputs := make(map[string]neta.Result)

	for _, layer := range graph.ExecutionOrder {
		if err := i.executeGraphLayer(ctx, layer, nodeMap, edges, outputs); err != nil {
			return outputs, err
		}
	}

	return outputs, nil
}

// executeGraphLayer executes all nodes in a single layer
func (i *Itamae) executeGraphLayer(ctx context.Context, layer []string, nodeMap map[string]neta.Definition, edges []neta.NodeEdge, outputs map[string]neta.Result) error {
	for _, nodeID := range layer {
		result, err := i.executeGraphNode(ctx, nodeID, nodeMap, edges, outputs)
		if err != nil {
			return err
		}
		outputs[nodeID] = result
	}
	return nil
}

// executeGraphNode executes a single node within graph context
func (i *Itamae) executeGraphNode(ctx context.Context, nodeID string, nodeMap map[string]neta.Definition, edges []neta.NodeEdge, outputs map[string]neta.Result) (neta.Result, error) {
	nodeDef, ok := nodeMap[nodeID]
	if !ok {
		return neta.Result{}, fmt.Errorf("node definition not found for ID: %s", nodeID)
	}

	i.startGraphNode(nodeID, nodeDef)

	result, duration, err := i.executeTimedNode(ctx, nodeDef, outputs, edges)
	if err != nil {
		return i.handleGraphNodeError(nodeID, duration, err)
	}

	i.completeGraphNode(nodeID, duration)
	return result, nil
}

// startGraphNode initializes node execution state
func (i *Itamae) startGraphNode(nodeID string, nodeDef neta.Definition) {
	i.store.SetNodeState(nodeID, neta.NodeStateExecuting, "")
	i.notifyNodeStarted(nodeID, nodeDef.Name, nodeDef.Type)
}

// executeTimedNode executes node and returns result with duration
func (i *Itamae) executeTimedNode(ctx context.Context, nodeDef neta.Definition, outputs map[string]neta.Result, edges []neta.NodeEdge) (neta.Result, time.Duration, error) {
	start := time.Now()
	result, err := i.executeNodeWithDataFlow(ctx, nodeDef, outputs, edges)
	duration := time.Since(start)
	return result, duration, err
}

// handleGraphNodeError handles node execution failure
func (i *Itamae) handleGraphNodeError(nodeID string, duration time.Duration, err error) (neta.Result, error) {
	i.store.SetNodeState(nodeID, neta.NodeStateError, err.Error())
	i.notifyNodeCompleted(nodeID, duration, err)
	i.store.CompleteExecution()
	return neta.Result{}, err
}

// completeGraphNode marks node as successfully completed
func (i *Itamae) completeGraphNode(nodeID string, duration time.Duration) {
	i.store.SetNodeState(nodeID, neta.NodeStateCompleted, "")
	i.notifyNodeCompleted(nodeID, duration, nil)
}

// executeNodeWithDataFlow executes a node with data flow from previous nodes
func (i *Itamae) executeNodeWithDataFlow(ctx context.Context, def neta.Definition, outputs map[string]neta.Result, edges []neta.NodeEdge) (neta.Result, error) {
	if def.IsGroup() {
		return i.executeGroup(ctx, def, "")
	}

	exec, err := i.pantry.Get(def.Type)
	if err != nil {
		return neta.Result{}, fmt.Errorf("node type not found: %s: %w", def.Type, err)
	}

	params := i.prepareNodeParams(def, outputs, edges)
	return exec.Execute(ctx, params)
}

// prepareNodeParams builds parameters with data flow from edges
func (i *Itamae) prepareNodeParams(def neta.Definition, outputs map[string]neta.Result, edges []neta.NodeEdge) map[string]interface{} {
	params := copyParams(def.Parameters)
	injectInputFromEdges(params, def.ID, outputs, edges)
	return params
}

// copyParams creates a copy of parameters map
func copyParams(original map[string]interface{}) map[string]interface{} {
	params := make(map[string]interface{})
	for k, v := range original {
		params[k] = v
	}
	return params
}

// injectInputFromEdges finds edges targeting this node and injects input
func injectInputFromEdges(params map[string]interface{}, nodeID string, outputs map[string]neta.Result, edges []neta.NodeEdge) {
	if _, hasInput := params["input"]; hasInput {
		return
	}

	for _, edge := range edges {
		if edge.Target == nodeID {
			if sourceResult, ok := outputs[edge.Source]; ok {
				params["input"] = sourceResult.Output
			}
		}
	}
}

// finalizeGraphExecution returns the result of the last executed node
func (i *Itamae) finalizeGraphExecution(outputs map[string]neta.Result, graph *neta.ExecutionGraph) (neta.Result, error) {
	i.store.CompleteExecution()

	if len(graph.ExecutionOrder) == 0 {
		return neta.Result{Output: nil}, nil
	}

	lastLayer := graph.ExecutionOrder[len(graph.ExecutionOrder)-1]
	if len(lastLayer) == 0 {
		return neta.Result{Output: nil}, nil
	}

	lastNodeID := lastLayer[0]
	if result, ok := outputs[lastNodeID]; ok {
		return result, nil
	}

	return neta.Result{Output: nil}, nil
}
