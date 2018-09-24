package merkletree

import "crypto/sha256"

type Node struct {
	left, right *Node
	sig         []byte
}

func NewLeafNode(sig []byte) *Node {
	leaf := &Node{}
	leaf.sig = sig
	return leaf
}

func NewParentNode(left, right *Node) *Node {
	parent := &Node{}

	// 必须有一个左子节点
	if left == nil {
		return nil
	}

	if right == nil {
		parent.sig = left.sig
	} else {
		combineSigs := append(left.sig, right.sig...)
		hash := sha256.Sum256(combineSigs)
		parent.sig = hash[:]
	}
	return parent
}

// 生成最底层的叶子节点,并挂载在最底层的父节点上,返回最底层的父节点数组
func ButtomLevelNodes(sigs [][]byte) []*Node {
	parents := make([]*Node, 0)

	for i := 0; i < len(sigs)-1; i += 2 {
		leftNode := NewLeafNode(sigs[i])
		rightNode := NewLeafNode(sigs[i+1])
		parentNode := NewParentNode(leftNode, rightNode)

		parents = append(parents, parentNode)
	}

	if len(sigs)/2 != 0 {
		leftNode := NewLeafNode(sigs[len(sigs)-1])
		parentNode := NewParentNode(leftNode, nil)
		parents = append(parents, parentNode)
	}

	return parents
}

// 根据指定层的非叶子节点数组,生成父节点数组
func InternalLevelNodes(nodes []*Node) []*Node {
	parents := make([]*Node, 0)

	for i := 0; i < len(nodes)-1; i += 2 {
		parentNode := NewParentNode(nodes[i], nodes[i+1])
		parents = append(parents, parentNode)
	}

	if len(nodes)/2 != 0 {
		parentNode := NewParentNode(nodes[len(nodes)-1], nil)
		parents = append(parents, parentNode)
	}

	return parents
}

type Tree struct {
	head *Node

	// 总计节点数量
	nodesNum int

	// 深度
	depth int
}

func NewTree(sigs [][]byte) *Tree {
	if len(sigs) <= 1 {
		return nil
	}

	tree := &Tree{}
	parents := ButtomLevelNodes(sigs)
	tree.depth = 1
	tree.nodesNum += len(parents)

	for ; len(parents) > 1; {
		parents = InternalLevelNodes(parents)

		tree.depth++
		tree.nodesNum += len(parents)
	}

	if len(parents) != 0 {
		panic("no root node")
	}

	tree.head = parents[0]
	return tree
}
