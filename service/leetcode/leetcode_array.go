package leetcode

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
func arrayDemo_maxSubArray(nums []int) int {
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

// lc:56
func arrayDemo_merge(intervals [][]int) [][]int {
	merge := func(a, b []int) []int {
		if a[0] <= b[0] && a[1] >= b[1] {
			return []int{a[0], a[1]}
		}
		if a[0] >= b[0] && a[1] <= b[1] {
			return []int{b[0], b[1]}
		}
		if a[1] >= b[0] && a[0] <= b[0] {
			return []int{a[0], b[1]}
		}
		if a[0] <= b[1] && a[0] >= b[0] {
			return []int{b[0], a[1]}
		}
		return nil
	}
	var result [][]int
	prevArray := intervals[0]
	for index := 1; index < len(intervals); index++ {
		mer := merge(prevArray, intervals[index])
		if mer == nil {
			result = append(result, prevArray)
			prevArray = intervals[index]
		} else {
			prevArray = mer
		}
	}
	result = append(result, prevArray)
	return result
}

// lc:189
func arrayDemo_rotate(nums []int, k int) {
	length := len(nums)
	start := length - k%length
	result := make([]int, 0, len(nums))
	result = append(result, nums[start:]...)
	result = append(result, nums[:start]...)
	nums = nums[:0]
	nums = append(nums, result...)
}

// lc:238
func arrayDemo_productExceptSelf(nums []int) []int {
	length := len(nums)
	leftResult := make([]int, length)
	rightResult := make([]int, length)
	sumLeft, sumRight := 1, 1
	for i := 0; i < length; i++ {
		if i == 0 {
			leftResult[0] = sumLeft
			rightResult[length-1-i] = sumRight
		} else {
			sumLeft = sumLeft * nums[i-1]
			leftResult[i] = sumLeft

			sumRight = sumRight * nums[length-i]
			rightResult[length-1-i] = sumRight
		}
	}
	result := make([]int, 0, length)
	for i := 0; i < length; i++ {
		result = append(result, leftResult[i]*rightResult[i])
	}
	return result
}
