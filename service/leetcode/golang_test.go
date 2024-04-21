package leetcode

import "testing"

func TestPlay(t *testing.T) {
	t.Run("hello", func(t *testing.T) {
		root := BuildListNode([]int{1, 2, 3, 4, 5})
		rev := reverseLinkedList(root)
		t.Log(rev)
	})
}

func TestList(t *testing.T) {
	list := []*Node{
		{
			Id:       0,
			ParentId: -1, // 当作nil
			Name:     "xx0",
		}, {
			Id:       1,
			ParentId: 0,
			Name:     "xx1",
		}, {
			Id:       2,
			ParentId: 0,
			Name:     "xx2",
		}, {
			Id:       3,
			ParentId: 1,
			Name:     "xx3",
		}, {
			Id:       4,
			ParentId: 2,
			Name:     "xx2",
		},
	}
	root := list2Tree(list)
	t.Log(root)
}
