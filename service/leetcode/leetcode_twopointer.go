package leetcode

import "sort"

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

// lc:11
func twoPointersDemo_maxArea(height []int) int {
	start, end := 0, len(height)-1
	var tmp, maxArea int
	for start < end {
		if tmp = min(height[start], height[end]) * (end - start); tmp > maxArea {
			maxArea = tmp
		}
		if height[start] > height[end] {
			end -= 1
		} else {
			start += 1
		}
	}
	return maxArea
}

// lc:15
func twoPointersDemo_threeSum(nums []int) [][]int {
	sort.Ints(nums)
	length := len(nums)
	var result [][]int
	for first := 0; first < length-2; first++ {
		if nums[first] > 0 {
			continue
		}
		if first > 0 && nums[first] == nums[first-1] {
			continue
		}
		for second := first + 1; second < length-1; second++ {
			firstAndSecond := nums[first] + nums[second]
			if firstAndSecond > 0 {
				continue
			}
			if second > first+1 && nums[second] == nums[second-1] {
				continue
			}
			third := length - 1
			for third > second {
				if tmp := nums[third] + firstAndSecond; tmp == 0 {
					result = append(result, []int{nums[first], nums[second], nums[third]})
					break
				} else if tmp < 0 {
					break
				}
				third -= 1
			}
		}
	}
	return result
}
