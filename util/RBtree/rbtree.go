package rbtree

import (
	keystruct "basic/util/KeyStruct"
	"sync"
)

const (
	RED   = 0
	BLACK = 1
)

type RBTreeNode struct {
	Left   *RBTreeNode
	Right  *RBTreeNode
	Parent *RBTreeNode
	Color  uint8

	key keystruct.KeyStruct
	val interface{}
}

/* 5 Attributes of RBTree
1.the node is either red or black
2.the root node is black
3.every leaves node is black(and they have nil values)
4.if a node is red, its children are black
5.every route from node to its successors poccesses same counts of black nodes

*/
type RBTree struct {
	root         *RBTreeNode
	_NIL         *RBTreeNode //to avoid hanging pointer, but not necessary
	elementCount int
	mtx          sync.RWMutex
}

func (rbt *RBTree) leftRotate(x *RBTreeNode) {
	if x.Right == rbt._NIL {
		return
	}
	y := x.Right
	x.Right = y.Left
	if y.Left != rbt._NIL {
		y.Left.Parent = x
	}
	y.Parent = x.Parent

	if x.Parent == rbt._NIL {
		rbt.root = y
	} else if x == x.Parent.Left {
		x.Parent.Left = y
	} else {
		x.Parent.Right = y
	}

	y.Left = x
	x.Parent = y
}

func (rbt *RBTree) rightRotate(x *RBTreeNode) {
	if x.Left == rbt._NIL {
		return
	}
	y := x.Left
	x.Left = y.Right
	if y.Right != rbt._NIL {
		y.Right.Parent = x
	}
	y.Parent = x.Parent

	if x.Parent == rbt._NIL {
		rbt.root = y
	} else if x == x.Parent.Left {
		x.Parent.Left = y
	} else {
		x.Parent.Right = y
	}

	y.Right = x
	x.Parent = y
}

func (rbt *RBTree) insert(i *RBTreeNode) *RBTreeNode {
	//step 1: insert the node like normal binary search tree
	//set its color red
	//left is small and right is big
	cur := rbt.root
	pos := rbt._NIL

	for cur != rbt._NIL {
		pos = cur
		if cur.key.CompareBiggerThan(i.key) {
			cur = cur.Left
		} else if i.key.CompareBiggerThan(cur.key) {
			cur = cur.Right
		} else {
			//maybe change the value here
			return cur
		}
	}

	i.Parent = pos
	if pos == rbt._NIL {
		//there is not root
		rbt.root = i
	} else if pos.key.CompareBiggerThan(i.key) {
		pos.Left = i
	} else {
		pos.Right = i
	}
	rbt.elementCount++
	//step 2 : maintains the balance tree
	rbt.fixAfterInsert(i)
	return i
}

func (rbt *RBTree) fixAfterInsert(i *RBTreeNode) {
	//new node z is a red node
	//if z's parent is black, it can insert directly & won't change the consistent
	//from bottom to top, maintains the consistent of RBTree
	for i.Parent.Color == RED {
		if i.Parent == i.Parent.Parent.Left {
			//y is an uncle node
			uncle := i.Parent.Parent.Right
			if uncle.Color == RED {
				i.Parent.Color = BLACK
				uncle.Color = BLACK
				i.Parent.Parent.Color = RED
				//and it will continue, if z is the root, it should be black
				//if new z's parent is black, we do nothing
				i = i.Parent.Parent
			} else {
				if i == i.Parent.Right {
					i = i.Parent
					rbt.leftRotate(i)
				}
				i.Parent.Color = BLACK
				i.Parent.Parent.Color = RED
				rbt.rightRotate(i.Parent.Parent)
			}
		} else {
			uncle := i.Parent.Parent.Left
			if uncle.Color == RED {
				i.Parent.Color = BLACK
				uncle.Color = BLACK
				i.Parent.Parent.Color = RED
				i = i.Parent.Parent
			} else {
				if i == i.Parent.Left {
					i = i.Parent
					rbt.rightRotate(i)
				}
				i.Parent.Color = BLACK
				i.Parent.Parent.Color = RED
				rbt.leftRotate(i.Parent.Parent)
			}
		}
	}
	rbt.root.Color = BLACK
}

func (rbt *RBTree) findSuccessor(cur *RBTreeNode) *RBTreeNode {
	if cur == rbt._NIL {
		return cur
	}

	if cur.Right != rbt._NIL {
		tmp := cur.Right
		for tmp.Left != rbt._NIL {
			tmp = tmp.Left
		}
		return tmp
	}

	tmp := cur.Parent
	for tmp != rbt._NIL && cur == tmp.Right {
		cur = tmp
	}
	return tmp
}

func (rbt *RBTree) delete(key keystruct.KeyStruct) *RBTreeNode {
	toDelete := rbt.search(key)
	if toDelete == rbt._NIL {
		return toDelete
	}
	//a copy of node to delete
	ret := &RBTreeNode{
		Left:   rbt._NIL,
		Right:  rbt._NIL,
		Parent: rbt._NIL,
		Color:  toDelete.Color,
		key:    toDelete.key,
		val:    toDelete.val,
	}

	return ret
}

//used for inner search and can be package to outside
func (rbt *RBTree) search(key keystruct.KeyStruct) *RBTreeNode {
	cur := rbt.root

	for cur != rbt._NIL {
		if cur.key.CompareBiggerThan(key) {
			cur = cur.Right
		} else if key.CompareBiggerThan(cur.key) {
			cur = cur.Left
		} else {
			break
		}
	}

	return cur
}

func New() (rbt RBTree) {
	initNode := &RBTreeNode{nil, nil, nil, BLACK, keystruct.DefaultKey{}, nil}
	rbt = RBTree{
		root:         initNode,
		_NIL:         initNode,
		elementCount: 0,
		mtx:          sync.RWMutex{},
	}
	return
}
