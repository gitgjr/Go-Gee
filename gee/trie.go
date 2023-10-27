package gee

import "strings"

type Node struct {
	pattern  string // 待匹配路由，例如 /p/:lang
	part     string // 路由中的一部分，例如 :lang
	children []*Node
	isWild   bool // 是否精确匹配，part 含有 : 或 * 时为true
}

func (n *Node) matchChild(part string) *Node {
	for _, child := range n.children {
		if child.part == part || child.isWild == true {
			return child
		}
	}
	return nil
}

func (n *Node) matchChildren(part string) []*Node {
	nodes := make([]*Node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild == true {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

func (n *Node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		n.pattern = pattern
		return
	}
	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		child = &Node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1) //TODO get migure out
}

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
