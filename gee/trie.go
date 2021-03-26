package gee

import (
	"strings"
)

type Node struct {
	pattern string
	part    string
	child   []*Node
	isWild  bool
}

//获取第一个匹配成功的子节点
func (n *Node) matchChild(part string) *Node {
	for _, child := range n.child {
		if child.part == part || child.isWild {
			return child
		}
	}

	return nil
}

//获取全部匹配的子节点
func (n *Node) matchChildren(part string) []*Node {
	nodes := make([]*Node, 0)
	for _, child := range n.child {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}

	return nodes
}

//新增子节点
func (n *Node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		n.pattern = pattern
		return
	}

	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		child = &Node{
			part:   part,
			isWild: part[0] == ':' || part[0] == '*',
		}
		n.child = append(n.child, child)
	}

	child.insert(pattern, parts, height+1)
}

//查找
func (n *Node) search(parts []string, height int) *Node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}

	part := parts[height]
	children := n.matchChildren(part)
	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}

	}

	return nil
}
