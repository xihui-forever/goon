package goon

import (
	"errors"
	"strings"
)

type (
	trieNode struct {
		char     string
		itemSet  []*Logic
		isEnding bool
		children map[rune]*trieNode
	}

	Trie struct {
		root *trieNode
	}
)

func NewTrie() *Trie {
	trieNode := NewTrieNode("/")
	return &Trie{trieNode}
}

func NewTrieNode(char string) *trieNode {
	return &trieNode{
		char:     char,
		itemSet:  []*Logic{},
		isEnding: false,
		children: make(map[rune]*trieNode),
	}
}

func (t *Trie) Insert(word string, item *Logic) error {
	node := t.root
	word = strings.TrimSuffix(word, "/") + "/"
	for _, code := range word {
		value, ok := node.children[code]
		if !ok {
			value = NewTrieNode(string(code))
			node.children[code] = value
		}
		node = value
	}

	if item.Method() != MethodPreUse && item.Method() != MethodPostUse && item.Method() != MethodUse {
		for _, value := range node.itemSet {
			if value.Method() == item.Method() {
				panic("logic already exists")
			}
		}
	}

	node.itemSet = append(node.itemSet, item)
	node.isEnding = true
	return nil
}

func (t *Trie) Find(method Method, word string) ([]*Logic, error) {
	node := t.root
	word = strings.TrimSuffix(word, "/") + "/"
	var itemList []*Logic
	for index, code := range word {
		value, ok := node.children[code]
		if !ok {
			return nil, errors.New("path is not unRegistered")
		}
		if value.char == "/" {
			for _, value := range value.itemSet {
				switch value.Method() {
				case MethodPreUse, MethodPostUse, MethodUse:
					itemList = append(itemList, value)
				case method:
					if index == len(word)-1 {
						itemList = append(itemList, value)
						goto END
					}
				}
			}
		}
		node = value
	}

END:
	if len(itemList) == 0 {
		return nil, errors.New("path is not unRegistered")
	}
	/*if itemList[len(itemList)-1].Method() != method {
		return nil, errors.New("method is not unRegistered")
	}*/

	return itemList, nil
}
