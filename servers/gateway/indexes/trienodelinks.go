package indexes

type trienodelinks map[rune]*TrieNode


func (t trienodelinks) add(value rune) bool {
	_, exists := t[value]
	t[value] = NewTrieNode()
	return !exists
}

func (t trienodelinks) has(value rune) bool {
	_, exists := t[value]
	return exists
}