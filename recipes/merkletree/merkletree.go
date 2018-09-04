package merkletree

import "crypto/sha256"

type Node struct {
	left, right *Node
	data        []byte
}

// 从下往上生长
// 给到左子节点,
func NewNode(left, right *Node, data []byte) *Node {
	node := &Node{}

	if left == nil && right == nil {
		hash := sha256.Sum256(data)
		node.data = hash[:]
	} else {
		preHashes := append(left.data, right.data...)
		hash := sha256.Sum256(preHashes)
		node.data = hash[:]
	}

	node.left = left
	node.right = right
	return node
}

type Tree struct {
	head *Node
}

func NewTree(data [][]byte) *Tree {
	// 保证节点是偶数,拿最后一个补
	if len(data)%2 != 0 {
		data = append(data, data[len(data)-1])
	}

	var nodes []*Node
	for i := 0; i < len(data); i++ {
		nodes = append(nodes, NewNode(nil, nil, data[i]))
	}

}
