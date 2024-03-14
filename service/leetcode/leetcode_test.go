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

func TestHello(t *testing.T) {
	type st struct {
		A int
	}
	v1 := &st{
		A: 1,
	}
	v2 := &st{
		A: 1,
	}
	t.Log(v1 == v2)
}

func TestAll(t *testing.T) {
	t.Run("slide window", func(t *testing.T) {
		slideWindowDemo_LengthOfLongestSubstring("abcabcbb")
	})
	t.Run("sub array", func(t *testing.T) {
		subArrayDemo_subArraySumV2([]int{1}, 0)
	})
	t.Run("linked list", func(t *testing.T) {
		intersect := &ListNode{
			Val: 2,
			Next: &ListNode{
				Val: 4,
			},
		}
		linkA := &ListNode{
			Val: 1,
			Next: &ListNode{
				Val: 9,
				Next: &ListNode{
					Val:  1,
					Next: intersect,
				},
			},
		}
		linkB := &ListNode{
			Val:  3,
			Next: intersect,
		}
		i := linkedListDemo_getIntersectionNode(linkA, linkB)
		t.Log(i)
	})
	t.Run("lined list 2", func(t *testing.T) {
		head := BuildListNode([]int{1, 2, 3, 4, 5})
		v := linkedListDemo_reverseList(head)
		t.Log(v)
	})
	t.Run("binary search", func(t *testing.T) {
		v := binarySearchDemo_searchInsert([]int{1, 3, 5, 6}, 5)
		t.Log(v)
	})
	t.Run("build binary", func(t *testing.T) {
		v := BuildBinaryTree([]int{1, 2, 2, -1, 3, -1, 3})
		t.Log(v)
		nums := BinaryBfsRead(v)
		t.Log(nums)
	})
	t.Run("binary diameter", func(t *testing.T) {
		v := BuildBinaryTree([]int{1, 2})
		x := binaryTreeDemo_diameterOfBinaryTree(v)
		t.Log(x)
	})
	t.Run("dynamic program", func(t *testing.T) {
		v := dynamicProgrammingDemo_knapsack([]int{1, 2, 3}, []int{100, 200, 300}, 10)
		t.Log(v)
	})
}
