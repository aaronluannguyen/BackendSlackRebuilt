package indexes

import (
	"testing"
	"reflect"
)

//TODO: implement automated tests for your trie data structure

func CreateInitialTrie() *Trie {
	trie := NewTrie()
	trie.Add("he", 1)
	trie.Add("hello", 2)
	trie.Add("hey", 3)
	trie.Add("go", 4)
	trie.Add("good", 5)
	trie.Add("goose", 6)
	trie.Add("世界", 7)

	return trie
}

func TestTrie_Add(t *testing.T) {
	cases := [] struct {
		key string
		val int64
		expectedValue int64
		expectedLength int
	} {
		{
			"he",
			1,
			1,
			1,
		},
		{
			"hello",
			2,
			2,
			2,
		},
		{
			"hey",
			3,
			3,
			3,
		},
		{
			"go",
			4,
			4,
			4,
		},
		{
			"good",
			5,
			5,
			5,
		},
		{
			"goose",
			6,
			6,
			6,
		},
	}
	trie := NewTrie()
	for _, c := range cases {
		trie.Add(c.key, c.val)
		result := trie.Find(c.key, 1)
		if result[0] != c.expectedValue {
			t.Errorf("expected %d for %s but got %d", c.expectedValue, c.key, result[0])
		}
		if trie.Len() != c.expectedLength {
			t.Errorf("incorrect trie length: got %d but expected %d", trie.Len(), c.expectedLength)
		}
	}
}

func TestTrie_Find(t *testing.T) {
	cases := [] struct {
		name string
		prefix string
		expectedSlice []int64
		max int
	} {
		{
			"Find H prefix without hitting max",
			"h",
			[]int64{1, 2, 3},
			6,
		},
		{
			"Find H prefix with max of 2",
			"h",
			[]int64{1, 2},
			2,
		},
		{
			"Test empty returned slice due to max of 0",
			"h",
			nil,
			0,
		},
		{
			"Test empty returned slice due to non-existent prefix",
			"helooo",
			nil,
			5,
		},
		{
			"Find he prefix alphabetical test",
			"he",
			[]int64{1, 2},
			2,
		},
		{
			"Find unicode prefix alphabetical test",
			"世",
			[]int64{7},
			1,
		},
	}

	trie := CreateInitialTrie()
	for _, c := range cases {
		result := trie.Find(c.prefix, c.max)
		if !reflect.DeepEqual(result, c.expectedSlice) {
			t.Errorf("case %s: incorrect slice returned. Got %v but expected %v", c.name, result, c.expectedSlice)
		}
	}
}

func TestTrie_Remove(t *testing.T) {
	cases := [] struct {
		name string
		word string
		val int64
		expectedSlice []int64
		max int
	} {
		{
			"Remove hello",
			"hello",
			2,
			nil,
			1,
		},
		{
			"Remove non-existent word",
			"zello",
			2,
			nil,
			1,
		},

	}

	trie := CreateInitialTrie()
	for _, c := range cases {
		trie.Remove(c.word, c.val)
		result := trie.Find(c.word, c.max)
		if !reflect.DeepEqual(result, c.expectedSlice) {
			t.Errorf("case %s: incorrect slice returned. Got %v but expected %v", c.name, result, c.expectedSlice)
		}
	}

	trie = CreateInitialTrie()
	trie.Remove("hello", 2)
	result := trie.Find("hel", 2)
	if result != nil {
		t.Errorf("trimming did not happen properly")
	}
}