package handler

import (
	"errors"
)

type (
	trieNode struct {
		char     string
		itemSet  []*Item
		isEnding bool
		children map[rune]*trieNode
	}

	Trie struct {
		root *trieNode
	}
)

func NewTrie() *Trie {
	trieNode := NewTrieNode("/", nil)
	return &Trie{trieNode}
}

func NewTrieNode(char string, item *Item) *trieNode {
	return &trieNode{
		char:     char,
		itemSet:  []*Item{},
		isEnding: false,
		children: make(map[rune]*trieNode),
	}
}

func (t *Trie) Insert(word string, item *Item) error {
	node := t.root
	for _, code := range word {
		value, ok := node.children[code]
		if !ok {
			value = NewTrieNode(string(code), nil)
			node.children[code] = value
		}
		node = value
	}

	for _, value := range node.itemSet {
		if value.method == item.method {
			panic("logic already exists")
		}
	}

	node.itemSet = append(node.itemSet, item)
	node.isEnding = true
	return nil
}

func (t *Trie) Find(method Method, word string) ([]*Item, error) {
	node := t.root

	var itemList []*Item
	for _, code := range word {
		value, ok := node.children[code]
		if !ok {
			return nil, errors.New("path is not unRegistered")
		}
		if value.char == "/" {
			for _, value := range node.itemSet {
				if value.method == PreUse || value.method == PostUse {
					itemList = append(itemList, value)
				}
			}
		}
		node = value
	}
	/*if node.itemSet == nil {
		return nil, errors.New("logic is not unRegistered")
	}*/

	var key int
	var item *Item
	for index, value := range node.itemSet {
		if value.method == method {
			key = index
			item = value
			break
		}
		if index == len(node.itemSet)-1 {
			return nil, errors.New("logic is not unRegistered")
		}
	}
	for i := 0; i < key; i++ {
		value := node.itemSet[i]
		if value.method == PreUse || value.method == PostUse {
			itemList = append(itemList, value)
		}
	}
	itemList = append(itemList, item)

	return itemList, nil
}
