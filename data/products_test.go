package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChecksValidation(t *testing.T) {
	p := &Product{Name: "John", Price: 1.00, SKU: "abc-abc-abc"}

	err := p.Validate()

	assert.NoError(t, err)
}
