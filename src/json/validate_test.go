package json

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	assert := assert.New(t)
	err := Validate(`["k1","value"`)
	assert.Equal(nil, err)
}
