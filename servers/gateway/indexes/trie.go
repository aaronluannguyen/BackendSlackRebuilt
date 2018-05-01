package indexes

import (
	"sync"
	"sort"
)

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
	prevNode *TrieNode
}

//NewTrie constructs a new Trie.
func NewTrie() *Trie {
	return &Trie{
		root: NewTrieNode(),
	}
}

//Len returns the number of entries in the trie.
func (t *Trie) Len() int {
	return t.entries
}

//Add adds a key and value to the trie.
func (t *Trie) Add(key string, value int64) {
	t.mx.Lock()
	defer t.mx.Unlock()

	currNode := t.root
	for _, c := range key {
		currNode = AddHelper(currNode, c)
	}
	added := currNode.values.add(value)
	if added {
		t.entries++
	}
}

//Find finds `max` values matching `prefix`. If the trie
//is entirely empty, or the prefix is empty, or max == 0,
//or the prefix is not found, this returns a nil slice.
func (t *Trie) Find(prefix string, max int) []int64 {
	t.mx.RLock()
	defer t.mx.RUnlock()

	currNode := t.root
	for _, c := range prefix {
		tempNode := GetPrefixLastNodeHelper(currNode, c)
		if tempNode != nil {
			currNode = tempNode
		} else {
			currNode = nil
			break
		}
	}
	var valSet []int64
	if currNode != nil {
		FindValuesMatchingPrefix(currNode, max, &valSet)
	}
	return valSet
}

//Remove removes a key/value pair from the trie
//and trims branches with no values.
func (t *Trie) Remove(key string, value int64) {
	t.mx.Lock()
	defer t.mx.Unlock()

	currNode := t.root
	for _, c := range key {
		currNode = RemoveHelper(currNode, c)
		if currNode == nil {
			break
		}
		if removed := currNode.values.remove(value); removed {
			t.entries--
		}
		TrimTrie(currNode)
	}
}

func NewTrieNode() *TrieNode {
	return &TrieNode{
		values: make(int64set),
		links: make(trienodelinks),
	}
}

//AddHelper serves as a helper function that takes in a node and
//current ch and checks if that node's links contain a node that has
//that rune as its key. If there isn't an existing node for that ch, then
//a new one will be added
func AddHelper(node *TrieNode, ch rune) *TrieNode {
	if node.links.has(ch) {
		return node.links[ch]
	}
	node.links.add(ch)
	node.links[ch] = NewTrieNode()
	node.links[ch].prevNode = node
	return node.links[ch]
}

func RemoveHelper(node *TrieNode, ch rune) *TrieNode {
	nextNode, exists := node.links[ch]
	if exists {
		return nextNode
	}
	return nil
}

func TrimTrie(node *TrieNode) {
	currNode := node
	for len(currNode.values) == 0 && len(currNode.links) == 0 {
		value := currNode.key
		delete(currNode.prevNode.links, value)
		currNode = currNode.prevNode
	}
}

func GetPrefixLastNodeHelper(node *TrieNode, ch rune) *TrieNode {
	if node.links.has(ch) {
		return node.links[ch]
	}
	return nil
}

func FindValuesMatchingPrefix(node *TrieNode, max int, valSet *[]int64) {
	sortedVals := node.values.all()
	for _, v := range sortedVals {
		if len(*valSet) == max {
			return
		}
		*valSet = append(*valSet, v)
	}
	sortedKeys := GetSortedRuneKeys(node.links)
	for _, k := range sortedKeys {
		FindValuesMatchingPrefix(node.links[k], max, valSet)
	}
	return
}

func GetSortedRuneKeys(linkMap trienodelinks) []rune {
	sortedKeys := make([]rune, 0, len(linkMap))
	for k := range linkMap {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Slice(sortedKeys, func(i, j int) bool {return sortedKeys[i] < sortedKeys[j]})
	return sortedKeys
}