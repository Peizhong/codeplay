package register

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnRegister(t *testing.T) {
	err := UnRegister("news/content-10.128.64.23:2379")
	assert.NoError(t, err)
}
