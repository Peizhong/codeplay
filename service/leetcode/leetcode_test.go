package leetcode

import (
	"testing"

	"github.com/peizhong/codeplay/pkg/logger"
)

func TestMain(m *testing.M) {
	logger.InitLogger()
	defer logger.Flush()

	m.Run()
}
