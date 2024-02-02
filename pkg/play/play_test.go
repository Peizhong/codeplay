package play

import "testing"

func TestCompare(t *testing.T) {
	type s struct {
		X *struct {
			Y int
		}
	}
	a := s{
		X: &struct{ Y int }{
			Y: 1,
		},
	}
	b := s{
		X: &struct{ Y int }{
			Y: 1,
		},
	}
	t.Log(a == b)
	c1 := make(chan int)
	c2 := make(chan int)
	t.Log(c1 == c2)

	var i1 any = make([]string, 1)
	var i2 any = make([]string, 1)
	t.Log(i1 == i2)

	m := make(map[string]struct{})
	t.Log(m)
}
