package leetcode

import (
	"container/list"
	"log"
	"math"
	"sort"
)

// 时间复杂度
// 空间复杂度

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a >= b {
		return a
	}
	return b
}

func twoSum(nums []int, target int) []int {
	// use hash to store value
	hashTable := make(map[int]int)
	for i, x := range nums {
		if p, ok := hashTable[target-x]; ok {
			return []int{p, i}
		}
		hashTable[x] = i
	}
	return nil
}

// lc:283
func twoPointersDemo_MoveZeroes(nums []int) {
	// Given an integer array nums, move all 0's to the end of it while maintaining the relative order of the non-zero elements.
	// Example: Input: nums = [0,1,0,3,12], Output: [1,3,12,0,0]
	left, right, n := 0, 0, len(nums)
	for right < n {
		if nums[right] != 0 {
			nums[left], nums[right] = nums[right], nums[left]
			left++
		}
		right++
	}
}

// lc:3
func slideWindowDemo_LengthOfLongestSubstring(s string) int {
	start, end, n, maxLength := 0, 0, len(s), 0
	if n <= 1 {
		return n
	}
	// character last index
	dict := make(map[rune]int)
	asRune := []rune(s)
	for end < n {
		log.Println("before", start, end)
		if index, exist := dict[asRune[end]]; exist {
			log.Println("exist, old index", index)
			maxLength = max(maxLength, end-start)
			if index >= start {
				start = index + 1
			}
		}
		dict[asRune[end]] = end
		end += 1
		log.Println("after", start, end)
	}
	maxLength = max(maxLength, end-start)
	return maxLength
}

// lc:560
func subArrayDemo_subArraySum(nums []int, k int) int {
	n := len(nums)
	if n == 0 {
		return 0
	}
	match := 0
	for begin := 0; begin < n; begin++ {
		var tmp int
		for end := begin; end < n; end++ {
			tmp += nums[end]
			if tmp == k {
				match += 1
				continue
			}
		}
	}
	return match
}

// lc:560
func subArrayDemo_subArraySumV2(nums []int, k int) int {
	// nums[0...i] 前i个数的和
	// nums[i...j] = nums[end] - nums[start] = k -> sum[end] - k = sum[start]
	vals := make(map[int]int)
	vals[0] = 1
	var sum int
	var match int
	for i := 0; i < len(nums); i++ {
		sum += nums[i]
		if count, exist := vals[sum-k]; exist {
			match += count
		}
		vals[sum] += 1
	}
	return match
}

// lc:53, array
func normalArrayDemo_maxSubArray(nums []int) int {
	// Given an integer array nums, find the subarray with the largest sum, and return its sum.
	// 动态规划
	n := len(nums)
	if n == 0 {
		return 0
	}
	bestValue := make([]int, n)
	bestValue[0] = nums[0]
	maxValue := nums[0]
	if n == 1 {
		return maxValue
	}
	for index := 1; index < n; index++ {
		bestValue[index] = max(nums[index], bestValue[index-1]+nums[index])
		maxValue = max(maxValue, bestValue[index])
	}
	return maxValue
}

// lc: 73, matrix
func matrixDemo_setZeroes(matrix [][]int) {
	yLen := len(matrix)
	if yLen == 0 {
		return
	}
	xLen := len(matrix[0])
	var points [][2]int
	for i := 0; i < yLen; i++ {
		for j := 0; j < xLen; j++ {
			if matrix[i][j] == 0 {
				points = append(points, [2]int{i, j})
			}
		}
	}
	spread := func(x, y int) {
		for yCursor := 0; yCursor < yLen; yCursor++ {
			matrix[yCursor][x] = 0
		}
		for xCursor := 0; xCursor < xLen; xCursor++ {
			matrix[y][xCursor] = 0
		}
	}
	for _, point := range points {
		spread(point[0], point[1])
	}
}

type ListNode struct {
	Val  int
	Next *ListNode
}

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

func BuildListNode(nums []int) *ListNode {
	root := &ListNode{}
	cur := root
	for i, v := range nums {
		if i == 0 {
			cur.Val = v
		} else {
			cur.Next = &ListNode{
				Val: v,
			}
			cur = cur.Next
		}
	}
	return root
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

// lc: 234
func linkedListDemo_isPalindrome(head *ListNode) bool {
	// Given the head of a singly linked list, return true if it is a palindrome or false otherwise.
	return false
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

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

func BuildBinaryTree(nums []int) *TreeNode {
	// level-1 [0]
	// level-2 [1 | 2] 1+2
	// level-3 [3 4 | 5 6] 1+2+4
	// level-4 [7 8 9 10 | 11 12 13 14]
	var build func(parent int) *TreeNode
	build = func(loc int) *TreeNode {
		if loc >= len(nums) {
			return nil
		}
		if nums[loc] == -1 {
			return nil
		}
		node := &TreeNode{
			Val: nums[loc],
		}
		node.Left = build(loc*2 + 1)
		node.Right = build(loc*2 + 2)
		return node
	}
	return build(0)
}

func BinaryBfsRead(root *TreeNode) []int {
	var result []int
	fifo := list.New()
	fifo.PushBack(root)
	for {
		v := fifo.Front()
		if v == nil {
			break
		}
		fifo.Remove(v)
		n := v.Value.(*TreeNode)
		result = append(result, n.Val)
		if n.Left != nil {
			fifo.PushBack(n.Left)
		} else if n.Right != nil {
			fifo.PushBack(&TreeNode{
				Val: -1,
			})
		}
		if n.Right != nil {
			fifo.PushBack(n.Right)
		} else if n.Left != nil {
			fifo.PushBack(&TreeNode{
				Val: -1,
			})
		}
	}
	return result
}

// lc: 94
func binaryTreeDemo_inorderTraversal(root *TreeNode) []int {
	if root == nil {
		return nil
	}
	var result []int
	result = append(result, binaryTreeDemo_inorderTraversal(root.Left)...)
	result = append(result, root.Val)
	result = append(result, binaryTreeDemo_inorderTraversal(root.Right)...)
	return result
}

// lc: 104
func binaryTreeDemo_maxDepth(root *TreeNode) int {
	// Given the root of a binary tree, return its maximum depth.
	var dfs func(node *TreeNode) int
	dfs = func(node *TreeNode) int {
		if node == nil {
			return 0
		}
		return 1 + max(dfs(node.Left), dfs(node.Right))
	}
	return dfs(root)
}

// lc: 226
func binaryTreeDemo_invertTree(root *TreeNode) *TreeNode {
	// Given the root of a binary tree, invert the tree, and return its root.
	var invert func(node *TreeNode)
	invert = func(node *TreeNode) {
		if node == nil {
			return
		}
		invert(node.Left)
		invert(node.Right)
		node.Left, node.Right = node.Right, node.Left
	}
	invert(root)
	return root
}

// lc: 101
func binaryTreeDemo_isSymmetric(root *TreeNode) bool {
	// Given the root of a binary tree, check whether it is a mirror of itself (i.e., symmetric around its center).
	//       1
	//     2      2
	//  3  4    4   5
	// 6  7 8 9  9 8  7 6

	var checkEqual func(a, b *TreeNode) bool
	checkEqual = func(a, b *TreeNode) bool {
		if a != nil && b != nil {
			return a.Val == b.Val && checkEqual(a.Left, b.Right) && checkEqual(a.Right, b.Left)
		}
		if a == nil && b == nil {
			return true
		}
		return false
	}
	return checkEqual(root.Left, root.Right)
}

func binaryTreeDemo_sortedArrayToBST(nums []int) *TreeNode {
	var buildTree func(n []int) *TreeNode
	buildTree = func(n []int) *TreeNode {
		length := len(n)
		if length == 0 {
			return nil
		}
		if length == 1 {
			return &TreeNode{
				Val: n[0],
			}
		}
		mid := length / 2
		node := &TreeNode{
			Val: n[mid],
		}
		node.Left = buildTree(n[:mid])
		node.Right = buildTree(n[mid+1:])
		return node
	}
	return buildTree(nums)
}

// lc: 543
func binaryTreeDemo_diameterOfBinaryTree(root *TreeNode) int {
	var maxV int
	var diameter func(*TreeNode) (max int)
	diameter = func(tn *TreeNode) int {
		if tn == nil {
			return 0
		}
		leftRoute := diameter(tn.Left)
		rightRoute := diameter(tn.Right)
		maxV = max(leftRoute+rightRoute, maxV)
		return max(leftRoute, rightRoute) + 1
	}
	diameter(root)
	return maxV
}

// lc: 200, grid
func gridDemo_numIslands(grid [][]byte) int {
	return -1
}

func binarySearchDemo_searchInsert(nums []int, target int) int {
	left, right := 0, len(nums)-1
	pos := (right - left) / 2
	for left <= right {
		if nums[pos] == target {
			return pos
		}
		if target < nums[pos] {
			right = right - 1
		} else {
			left = left + 1
		}
		pos = left + (right-left)/2
	}
	return pos
}

func stackDemo_isValidParentheses(s string) bool {
	var ls []string
	for _, r := range s {
		c := string(r)
		if c == "(" || c == "[" || c == "{" {
			ls = append(ls, c)
		} else if c == ")" {
			end := len(ls) - 1
			if end < 0 {
				return false
			}
			if ls[end] == "(" {
				ls = ls[:end]
			} else {
				return false
			}
		} else if c == "]" {
			end := len(ls) - 1
			if end < 0 {
				return false
			}
			if ls[end] == "[" {
				ls = ls[:end]
			} else {
				return false
			}
		} else if c == "}" {
			end := len(ls) - 1
			if end < 0 {
				return false
			}
			if ls[end] == "{" {
				ls = ls[:end]
			} else {
				return false
			}
		} else {
			return false
		}
	}
	return len(ls) == 0
}

func greedDemo_maxProfit(prices []int) int {
	cost, profit := math.MaxInt, 0
	for _, price := range prices {
		cost = min(cost, price)
		profit = max(profit, price-cost)
	}
	return profit
}

func dynamicProgrammingDemo_climbStairs(n int) int {
	// 转移方程 f(x) = f(x-1) + f(x-2)
	dp := make([]int, n+1)
	dp[0] = 1
	dp[1] = 1
	for i := 2; i <= n; i++ {
		dp[i] = dp[i-1] + dp[i-2]
	}
	return dp[n]
}

func dynamicProgrammingDemo_knapsack(weights []int, prices []int, total int) int {
	// x 轴: 背包容量
	// y 轴: 可选择的商品
	// 转移方程 f(x) = 已用背包最大价值 + 剩余商品最大价值
	dp := make([][]int, 0, len(weights)+1)
	for i := 0; i <= len(weights); i++ {
		dp = append(dp, make([]int, total+1))
	}
	for w := weights[0]; w <= total; w++ {
		dp[1][w] = prices[0]
	}
	for item := 2; item <= len(weights); item++ {
		for w := 1; w <= total; w++ {
			if weights[item] <= w {
				// take current, drop others
				price1 := dp[item][total-weights[item]] + prices[item]
				price2 := dp[item-1][w]
				dp[item][w] = max(price1, price2)
			} else {
				dp[item][w] = dp[item-1][w]
			}
		}
		log.Println(dp[item])
	}
	return 0
}

func dynamicProgrammingDemo_generatePascal(numRows int) [][]int {
	gen := func(nums []int) []int {
		res := make([]int, 0, len(nums)+1)
		res = append(res, nums[0])
		for i := 1; i < len(nums); i++ {
			res = append(res, nums[i]+nums[i-1])
		}
		res = append(res, nums[len(nums)-1])
		return res
	}
	rows := make([][]int, 0, numRows)
	rows = append(rows, []int{1})
	for i := 1; i < numRows; i++ {
		rows = append(rows, gen(rows[i-1]))
	}
	return rows
}

func miscDemo_singleNumber(nums []int) int {
	dict := make(map[int]struct{}, len(nums))
	for _, v := range nums {
		if _, exist := dict[v]; exist {
			delete(dict, v)
			continue
		}
		dict[v] = struct{}{}
	}
	if len(dict) == 1 {
		for v := range dict {
			return v
		}
	}
	// 异或: 不相同的为1
	var v int
	for _, num := range nums {
		v = v ^ num
	}
	return v
}

// lc: 169
func miscDemo_majorityElement(nums []int) int {
	// 排序后
	sort.Ints(nums)
	return nums[len(nums)/2]
}
