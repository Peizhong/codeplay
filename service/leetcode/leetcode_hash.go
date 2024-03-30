package leetcode

// lc:1
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
