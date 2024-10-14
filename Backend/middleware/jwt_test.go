package middleware

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadKeys(t *testing.T) {
	_, _, _, err := LoadKeys()
	assert.NoError(t, err)
}


