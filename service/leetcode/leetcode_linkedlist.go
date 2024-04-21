package leetcode

// lc: 160
func linkedListDemo_getIntersectionNode(headA, headB *ListNode) *ListNode {
	// Given the heads of two singly linked-lists headA and headB, return the node at which the two lists intersect.
	aCursor := headA
	bCursor := headB
	all := map[*ListNode]struct{}{
		aCursor: struct{}{},
		bCursor: struct{}{},
	}
	for {
		aCursor = aCursor.Next
		bCursor = bCursor.Next
		if aCursor == nil && bCursor == nil {
			return nil
		}
		if aCursor != nil {
			if _, exist := all[aCursor]; exist {
				return aCursor
			}
			all[aCursor] = struct{}{}
		}
		if bCursor != nil {
			if _, exist := all[bCursor]; exist {
				return bCursor
			}
			all[bCursor] = struct{}{}
		}
	}
}

// lc: 160
func linkedListDemo_getIntersectionNodeV2(headA, headB *ListNode) *ListNode {
	// 两个链表同时走，如果a走到底了，转到b上面走。如果b走到底了，转到a上面走。如果链表相交，他们会重合。
	aCursor := headA
	bCursor := headB
	for {
		if aCursor == bCursor {
			return aCursor
		}
		aCursor = aCursor.Next
		bCursor = bCursor.Next
		if aCursor == nil && bCursor == nil {
			return nil
		}
		if aCursor == nil {
			aCursor = headB
		}
		if bCursor == nil {
			bCursor = headA
		}
	}
}

// lc: 206
func linkedListDemo_reverseList(head *ListNode) *ListNode {
	var prev *ListNode
	cur := head
	if cur == nil {
		return nil
	}
	for {
		// 1->2->3->4
		// 1<-2->3->4
		// 1<-2<-3->4
		// 1<-2<-3<-4
		next := cur.Next
		cur.Next = prev
		if next == nil {
			return cur
		}
		prev = cur
		cur = next
	}
}

func reverseLinkedList(head *ListNode) *ListNode {
	writer := &ListNode{
		Val: head.Val,
	}
	reader := head.Next
	for {
		if reader == nil {
			break
		}
		newNode := &ListNode{
			Val:  reader.Val,
			Next: writer,
		}
		writer = newNode
		reader = reader.Next
	}
	return writer
}

// lc: 234
func linkedListDemo_isPalindrome(head *ListNode) bool {
	// Given the head of a singly linked list, return true if it is a palindrome or false otherwise.
	revHead := reverseLinkedList(head)
	cursorA := head
	cursorB := revHead
	for {
		if cursorA.Val != cursorB.Val {
			return false
		}
		cursorA = cursorA.Next
		cursorB = cursorB.Next
		if cursorA == nil {
			break
		}
	}
	return true
}

// lc: 141
func linkedListDemo_hasCycle(head *ListNode) bool {
	if head == nil {
		return false
	}
	slow := head
	fast := head.Next
	for {
		if fast == nil || fast.Next == nil || fast.Next.Next == nil {
			return false
		}
		if fast == slow {
			return true
		}
		fast = fast.Next.Next
		slow = slow.Next
	}
}
