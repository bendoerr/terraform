// Package backend provides interfaces that the CLI uses to interact with
// Terraform. A backend provides the abstraction that allows the same CLI
// to simultaneously support both local and remote operations for seamlessly
// using Terraform in a team environment.
package backend

import (
	"github.com/hashicorp/terraform/config/module"
	"github.com/hashicorp/terraform/state"
	"github.com/hashicorp/terraform/terraform"
)

// Backend is the minimal interface that must be implemented to enable Terraform.
type Backend interface {
	// Ask for input and configure the backend. Similar to
	// terraform.ResourceProvider.
	//Input(*terraform.ResourceConfig) (*terraform.ResourceConfig, error)
	Validate(*terraform.ResourceConfig) ([]string, []error)
	Configure(*terraform.ResourceConfig) error

	// State returns the current state for this environment. This state may
	// not be loaded locally: the proper APIs should be called on state.State
	// to load the state.
	State() (state.State, error)
}

// Enhanced implements additional behavior on top of a normal backend.
//
// Enhanced backends allow customizing the behavior of Terraform operations.
// This allows Terraform to potentially run operations remotely, load
// configurations from external sources, etc.
type Enhanced interface {
	Backend

	// Operation performs a Terraform operation such as refresh, plan, apply.
	// It is up to the implementation to determine what "performing" means.
	// This should block until the operation is complete.
	Operation(*Operation) error
}

// An operation represents an operation for Terraform to execute.
type Operation struct {
	Type OperationType // Enum representing the available ops

	// Configuration related to the operation:

	Id          string       // For local, this is a plan file
	Destroy     bool         // Destroy mode
	Module      *module.Tree // Can be nil
	Parallelism int          // Max in-flight parallel for terraform.Context
	Targets     []string
	Variables   map[string]interface{}

	// Meta

	UIIn     terraform.UIInput  // Ui to read
	UIOut    terraform.UIOutput // Ui to write
	CancelCh <-chan struct{}    // Graceful exit by closing
}
