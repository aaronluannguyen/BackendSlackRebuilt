package indexes

import "sync"

//TODO: implement a trie data structure that stores
//keys of type string and values of type int64

type Trie struct {
	root *TrieNode
	mx sync.RWMutex
	entries int
}

type TrieNode struct {
	key rune
	values int64set
	links trienodelinks
}

//NewTrie constructs a new Trie.
func NewTrie() *Trie {
	return &Trie{}
}

//Len returns the number of entries in the trie.
func (t *Trie) Len() int {
	return t.entries
}

//Add adds a key and value to the trie.
func (t *Trie) Add(key string, value int64) {
	currNode := t.root
	for _, c := range key {
		currNode = AddHelper(currNode, c)
	}
	// check if int64 already has value. If not, add.
}

//Find finds `max` values matching `prefix`. If the trie
//is entirely empty, or the prefix is empty, or max == 0,
//or the prefix is not found, this returns a nil slice.
func (t *Trie) Find(prefix string, max int) []int64 {
	panic("implement this function according to the comments above")
}

//Remove removes a key/value pair from the trie
//and trims branches with no values.
func (t *Trie) Remove(key string, value int64) {
	panic("implement this function according to the comments above")
}

func AddHelper(node *TrieNode, ch rune) *TrieNode {
	nextNode, exists := node.links[ch]
	if exists {
		return nextNode
	}
	node.links.add(ch)
	return node.links[ch]
}