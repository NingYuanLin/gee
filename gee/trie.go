package gee

import "strings"

type node struct {
	pattern  string  // 路由的全称，只有在根节点，才!="" 例如/p/:lang
	part     string  // 当前节点代表的路由部分，例如:lang
	children []*node // 子节点
	isWild   bool    // 是否精确匹配 当part带有:或者*时，为true
}

// 第一个匹配成功的节点，用于插入
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// 所有匹配成功的节点，用于查找
func (n *node) matchChildren(part string) []*node {
	children := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			children = append(children, child)
		}
	}
	return children
}

// 插入child node
func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		n.pattern = pattern
		return
	}
	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		child = &node{
			part:   part,
			isWild: part[0] == ':' || part[0] == '*',
		}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1)
}

func (n *node) search(parts []string, height int) *node {
	// 实际上 如果n.part = "*xxx"，那一定没有children
	if height == len(parts) || strings.HasPrefix(n.part, "*") {
		// 这是为了防止以下情况：
		// 路由为 /a/b/c 但是实际的请求为/a/b
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
