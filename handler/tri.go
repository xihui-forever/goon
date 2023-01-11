package handler

import (
	"errors"
	"github.com/darabuchi/log"
)

type (
	trieNode struct {
		char      string
		logic     *Item
		methodSet []Method
		isEnding  bool
		children  map[rune]*trieNode
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
		char:      char,
		logic:     item,
		methodSet: []Method{},
		isEnding:  false,
		children:  make(map[rune]*trieNode),
	}
}

func (t *Trie) Find(method Method, word string) (*Item, error) {
	node := t.root

	for _, code := range word {
		value, ok := node.children[code]
		if !ok {
			return nil, errors.New("path is not unRegistered")
		}
		node = value
	}

	if method == "" {
		return node.logic, nil
	}

	for _, value := range node.methodSet {
		if value == method {
			log.Info(string(value))
			return node.logic, nil
		}
	}
	log.Errorf("method is illegal")
	return nil, errors.New("method is illegal")
}

func (t *Trie) Insert(method Method, word string, item *Item) error {
	node := t.root
	for _, code := range word {
		value, ok := node.children[code]
		if !ok {
			value = NewTrieNode(string(code), nil)
			node.children[code] = value
		}
		node = value
	}

	if node.logic != nil {
		return errors.New("logic already exist")
	}
	if node.methodSet != nil {
		for _, value := range node.methodSet {
			if value == method {
				return errors.New("method of this logic already exist")
			}
		}
	}

	node.logic = item
	if method != "" {
		node.methodSet = append(node.methodSet, method)
	}
	log.Info(node.methodSet)
	node.isEnding = true
	return nil
}
