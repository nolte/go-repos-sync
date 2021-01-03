package syncbuilder

import (
	"github.com/nolte/go-repos-sync/pkg/reposcrtl/repository"
)

// SelectProcess interface for colleting Import elements from Different Backends.
type SelectProcess interface {
	GetElements() ([]repository.Element, error)
}
