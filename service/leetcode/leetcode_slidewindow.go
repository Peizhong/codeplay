package leetcode

import "log"

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

// lc:438
func slideWindowDemo_findAnagrams(s string, p string) []int {
	expectDict := make(map[byte]int)
	for _, c := range p {
		expectDict[byte(c)] += 1
	}
	actualDict := make(map[byte]int)
	var result []int
	length := len(p)
	for start := 0; start < len(s)-length+1; start++ {
		if start == 0 {
			for index := 0; index < length; index++ {
				actualDict[s[start+index]] += 1
			}
		} else {
			actualDict[byte(s[start-1])] -= 1
			actualDict[byte(s[start+length])] += 1
		}
		allMatch := true
		for expectV, expectCount := range expectDict {
			if actualCount := actualDict[expectV]; actualCount != expectCount {
				allMatch = false
				break
			}
		}
		if allMatch {
			result = append(result, start)
		}
	}
	return result
}
