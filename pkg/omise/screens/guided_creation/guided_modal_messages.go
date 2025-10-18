package guided_creation

import "bento/pkg/neta"

// GuidedCompleteMsg is sent when the guided flow completes (success, cancel, or error)
type GuidedCompleteMsg struct {
	Success    bool
	Definition *neta.Definition
	Err        error
	Cancelled  bool
}
