package indexes

type trienodelinks map[rune]*TrieNode


func (t trienodelinks) add(value rune) bool {
	_, exists := t[value]
	t[value] = &TrieNode{}
	return !exists
}

func (t trienodelinks) remove(value rune) bool {
	_, exists := t[value]
	delete(t, value)
	return exists
}

func (t trienodelinks) has(value rune) bool {
	_, exists := t[value]
	return exists
}

func (t trienodelinks) all() []rune {
	result := make([]rune, 0 , len(t))
	for k := range t {
		result = append(result, k)
	}
	return result
}