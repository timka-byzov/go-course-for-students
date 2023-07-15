package tagcloud

import "sort"

// TagCloud aggregates statistics about used tags
type TagCloud struct {
	TagStatMap map[string]int
}

// TagStat represents statistics regarding single tag
type TagStat struct {
	Tag             string
	OccurrenceCount int
}

// New should create a valid TagCloud instance
// TODO: You decide whether this function should return a pointer or a value
func New() TagCloud {
	// TODO: Implement this
	return TagCloud{TagStatMap: make(map[string]int)}
}

// AddTag should add a tag to the cloud if it wasn't present and increase tag occurrence count
// thread-safety is not needed
// TODO: You decide whether receiver should be a pointer or a value
func (tc *TagCloud) AddTag(tag string) {
	tc.TagStatMap[tag] += 1
}

// TopN should return top N most frequent tags ordered in descending order by occurrence count
// if there are multiple tags with the same occurrence count then the order is defined by implementation
// if n is greater that TagCloud size then all elements should be returned
// thread-safety is not needed
// there are no restrictions on time complexity
// TODO: You decide whether receiver should be a pointer or a value
func (tc *TagCloud) TopN(n int) []TagStat {
	// TODO: Implement this

	rn := n
	if n > len(tc.TagStatMap) {
		rn = len(tc.TagStatMap)
	}

	res := make([]TagStat, rn)
	pairs := make([]struct {
		Key   string
		Value int
	}, len(tc.TagStatMap))

	i := 0
	for key, value := range tc.TagStatMap {
		pairs[i] = struct {
			Key   string
			Value int
		}{key, value}
		i++
	}

	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].Value > pairs[j].Value
	})

	for j := 0; j < rn; j++ {
		res[j] = TagStat{Tag: pairs[j].Key, OccurrenceCount: pairs[j].Value}
	}

	return res

}
