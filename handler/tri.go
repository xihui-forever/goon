package handler

import (
	"errors"
)

type (
	trieNode struct {
		char     string
		mapper   map[Method]*Item
		useSet   map[Method][]*Item
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
		mapper:   make(map[Method]*Item),
		useSet:   make(map[Method][]*Item),
		isEnding: false,
		children: make(map[rune]*trieNode),
	}
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

	if method != PreUse && method != PostUse {
		if node.mapper[method] != nil {
			panic("logic already exists")
		} else {
			node.mapper[method] = item
		}
	}

	node.useSet[method] = append(node.useSet[method], item)

	node.isEnding = true
	return nil
}

func (t *Trie) Find(method Method, word string) ([]*trieNode, error) {
	node := t.root

	var nodeList []*trieNode
	//var itemList []*Item
	for _, code := range word {
		value, ok := node.children[code]
		if !ok {
			return nil, errors.New("path is not unRegistered")
		}
		if value.char == "/" {
			nodeList = append(nodeList, node)
		}
		node = value
	}

	if node.mapper[method] == nil {
		return nil, errors.New("logic is not unRegistered")
	}
	nodeList = append(nodeList, node)

	/*for _, value := range nodeList {
		if value.mapper[PreUse] != nil {
			itemList = append(itemList, )
		}
	}

	itemList = append(itemList, node.mapper[method])

	for i := len(nodeList) - 1; i >= 0; i-- {
		if nodeList[i].mapper[PostUse] != nil {
			itemList = append(itemList, nodeList[i].mapper[PostUse])
		}
	}*/

	return nodeList, nil
}
