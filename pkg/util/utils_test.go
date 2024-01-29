package util

import (
	"errors"
	"fmt"
	"testing"
)

func TestErrors(t *testing.T) {
	oe := errors.New("graceful exit")
	e2 := fmt.Errorf("%w err", ErrGracefulExit)
	t.Log(errors.Is(e2, ErrGracefulExit), errors.Is(e2, oe))
}
