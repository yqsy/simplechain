package merkletree

import (
	"crypto/sha256"
	"reflect"
	"testing"
)

func create2EleTree() *Tree {
	tree := NewTree([][]byte{
		[]byte("52e422506d8238ef3196b41db4c41ee0afd659b6"),
		[]byte("6d0b51991ac3806192f3cb524a5a5d73ebdaacf8"),
	})
	return tree
}

func create3EleTree() *Tree {
	tree := NewTree([][]byte{
		[]byte("52e422506d8238ef3196b41db4c41ee0afd659b6"),
		[]byte("6d0b51991ac3806192f3cb524a5a5d73ebdaacf8"),
		[]byte("461848c8b70e5a57bd94008b2622796ec26db657"),
	})
	return tree
}

func create4EleTree() *Tree {
	tree := NewTree([][]byte{
		[]byte("52e422506d8238ef3196b41db4c41ee0afd659b6"),
		[]byte("6d0b51991ac3806192f3cb524a5a5d73ebdaacf8"),
		[]byte("461848c8b70e5a57bd94008b2622796ec26db657"),
		[]byte("c938037dc70d107b3386a86df7fef17a9983cf53"),
	})
	return tree
}

func create10EleTree() *Tree {
	tree := NewTree([][]byte{
		[]byte("52e422506d8238ef3196b41db4c41ee0afd659b6"),
		[]byte("6d0b51991ac3806192f3cb524a5a5d73ebdaacf8"),
		[]byte("461848c8b70e5a57bd94008b2622796ec26db657"),
		[]byte("c938037dc70d107b3386a86df7fef17a9983cf53"),
		[]byte("d9312928e5702168348fe67ee2a3e3a1b7bc7c93"),
		[]byte("506d93ebff5365d8f5dd9fedd4a063949be831a4"),
		[]byte("e45922755802b52f11599d4746035ecad18c0c46"),
		[]byte("994d89c38e5b9384235696a0efea5b6b93efb270"),
		[]byte("26fe8e189fd5bb3fe56d4d3def6494802cb8cba3"),
		[]byte("3cf4172b27b7b182db0dd68276f08f7c27561c32"),
	})
	return tree
}

func Test2Nodes(t *testing.T) {
	tree := create2EleTree()
	if tree.GetDepth() != 1 || tree.GetNodesNum() != 1 {
		t.Fatal("err")
	}
	tmpsha256 := sha256.Sum256(append([]byte("52e422506d8238ef3196b41db4c41ee0afd659b6"), []byte("6d0b51991ac3806192f3cb524a5a5d73ebdaacf8")...))
	if !reflect.DeepEqual(tree.GetRoot().sig, tmpsha256[:]) {
		t.Fatal("err")
	}
}

func Test3Nodes(t *testing.T) {
	tree := create3EleTree()

	if tree.GetDepth() != 2 || tree.GetNodesNum() != 3 {
		t.Fatal("err")
	}
	tmpsha2561 := sha256.Sum256(append([]byte("52e422506d8238ef3196b41db4c41ee0afd659b6"), []byte("6d0b51991ac3806192f3cb524a5a5d73ebdaacf8")...))
	tmpsha2562 := []byte("461848c8b70e5a57bd94008b2622796ec26db657")
	tmpsha2563 := sha256.Sum256(append(tmpsha2561[:], tmpsha2562...))
	if !reflect.DeepEqual(tree.GetRoot().sig, tmpsha2563[:]) {
		t.Fatal("err")
	}
}

func Test4Nodes(t *testing.T) {
	tree := create4EleTree()

	if tree.GetDepth() != 2 || tree.GetNodesNum() != 3 {
		t.Fatal("err")
	}
	tmpsha2561 := sha256.Sum256(append([]byte("52e422506d8238ef3196b41db4c41ee0afd659b6"), []byte("6d0b51991ac3806192f3cb524a5a5d73ebdaacf8")...))
	tmpsha2562 := sha256.Sum256(append([]byte("461848c8b70e5a57bd94008b2622796ec26db657"), []byte("c938037dc70d107b3386a86df7fef17a9983cf53")...))
	tmpsha2563 := sha256.Sum256(append(tmpsha2561[:], tmpsha2562[:]...))
	if !reflect.DeepEqual(tree.GetRoot().sig, tmpsha2563[:]) {
		t.Fatal("err")
	}
}

func Test10Nodes(t *testing.T) {
	tree := create10EleTree()
	if tree.GetDepth() != 4 || tree.GetNodesNum() != 11 {
		t.Fatal("err")
	}
}
