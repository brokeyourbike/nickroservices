package data

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
)

func TestNewExchangeRates(t *testing.T) {
	r, err := NewExchangeRates(hclog.Default())
	assert.NoError(t, err)

	fmt.Printf("%#v", r.rates)
}
