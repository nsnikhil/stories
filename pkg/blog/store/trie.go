package store

import (
	"fmt"
	"regexp"
)

const (
	startOfCapitalLetters = 65
	endOfCapitalLetters   = 90
	startOfSmallLetters   = 97
	endOfSmallLetters     = 122
	capitalA              = 'A'
	smallA                = 'a'
	invalidIndex          = -1
)

type Trie interface {
	insert(s, id string) []error
	getIDs(query string) (map[string]bool, []error)
}

type trieNode struct {
	links     [26]*trieNode
	endOfWord bool
	ids       []string
}

func newTrieNode() *trieNode {
	return &trieNode{}
}

func (t *trieNode) setEndOfWord() {
	t.endOfWord = true
}

func (t *trieNode) isEndOfWord() bool {
	return t.endOfWord
}

func (t *trieNode) getIDs() []string {
	return t.ids
}

func (t *trieNode) addID(id string) {
	t.ids = append(t.ids, id)
}

type characterTrie struct {
	root *trieNode
}

func removeNonCharacter(words []string) []string {
	var result []string

	reg := regexp.MustCompile(`[^a-zA-Z]+`)

	for _, word := range words {
		if len(word) == 0 {
			continue
		}

		//ignore word if it contains something other than characters
		if ok := reg.MatchString(word); ok {
			continue
		}

		result = append(result, word)
	}
	return result
}

func splitQuery(query string) []string {
	//split by punctuations or new line
	return removeNonCharacter(regexp.MustCompile(`[.,!?:;\s]`).Split(query, -1))
}

func (ct *characterTrie) insert(s, id string) []error {
	words := splitQuery(s)

	resErr := make([]error, 0)

	for _, word := range words {
		if len(word) == 0 {
			continue
		}

		if err := insert(ct, word, id); err != nil {
			resErr = append(resErr, err)
		}
	}

	return resErr
}

func (ct *characterTrie) getIDs(query string) (map[string]bool, []error) {
	words := splitQuery(query)

	res := make(map[string]bool)
	resErr := make([]error, 0)

	for _, word := range words {
		if len(word) == 0 {
			continue
		}

		if err := search(ct, word, res); err != nil {
			resErr = append(resErr, err)
		}
	}

	return res, resErr
}

func insert(ct *characterTrie, word, id string) error {
	curr := ct.root

	for _, char := range word {
		idx, err := index(char)
		if err != nil {
			return err
		}

		if curr.links[idx] == nil {
			curr.links[idx] = newTrieNode()
		}

		curr = curr.links[idx]
	}

	curr.setEndOfWord()
	curr.addID(id)

	return nil
}

func search(ct *characterTrie, word string, res map[string]bool) error {
	curr := ct.root

	for _, char := range word {
		idx, err := index(char)
		if err != nil {
			return err
		}

		if curr.links[idx] == nil {
			return fmt.Errorf("%s not found", word)
		}

		curr = curr.links[idx]
	}

	if !curr.isEndOfWord() {
		return fmt.Errorf("%s not found", word)
	}

	for _, id := range curr.getIDs() {
		if !res[id] {
			res[id] = true
		}
	}

	return nil
}

func index(char int32) (int32, error) {
	if char >= startOfCapitalLetters && char <= endOfCapitalLetters {
		return char - capitalA, nil
	}

	if char >= startOfSmallLetters && char <= endOfSmallLetters {
		return char - smallA, nil
	}

	return invalidIndex, fmt.Errorf("%c is not a valid character", char)
}

func NewCharacterTrie() Trie {
	return &characterTrie{
		root: newTrieNode(),
	}
}
