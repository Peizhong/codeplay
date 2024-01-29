package mongodb_exporter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListService(t *testing.T) {
	ListService()
}

func TestRegister(t *testing.T) {
	err := RegisterService("10.131.249.25", 28029, "xj-db-wangpz", "xj", "xj-db", "xj-db")
	assert.NoError(t, err)
}

func TestDeRegister(t *testing.T) {
	err := DeRegister("xj-db-wangpz-10.131.249.25:28028")
	assert.NoError(t, err)
}
