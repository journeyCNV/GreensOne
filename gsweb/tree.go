package gsweb

import (
	"errors"
	"strings"
)

type TrieTree struct {
	root *node
}

type node struct {
	isLast   bool                //该节点是否为最终节点，是否能完成一个独立的uri
	segment  string              //节点字符串
	handlers []ControllerHandler //这个节点包含的控制器 + 中间件
	children []*node             //这个节点下的子节点
}

func NewTrieTree() *TrieTree {
	root := NewNode()
	return &TrieTree{root}
}

func NewNode() *node {
	return &node{
		isLast:   false,
		segment:  "",
		children: []*node{},
	}
}

func (t *TrieTree) FindHandler(uri string) []ControllerHandler {
	matchNode := t.root.matchNode(uri)
	if matchNode == nil {
		return nil
	}
	return matchNode.handlers
}

func (t *TrieTree) AddRouter(uri string, handlers []ControllerHandler) error {
	n := t.root
	if n.matchNode(uri) != nil {
		return errors.New("route exist:" + uri)
	}

	segments := strings.Split(uri, "/")
	for i, seg := range segments {
		if !isWildSegment(seg) {
			seg = strings.ToLower(seg)
		}
		isLast := i == len(segments)-1

		// 找到匹配的节点
		var satisNode *node
		children := n.filterChildNodes(seg)
		if len(children) > 0 {
			for _, child := range children {
				for child.segment == seg {
					satisNode = child
					break
				}
			}
		}
		if satisNode == nil {
			currNode := NewNode()
			currNode.segment = seg
			if isLast {
				currNode.isLast = true
				currNode.handlers = handlers
			}
			n.children = append(n.children, currNode)
			satisNode = currNode
		}

		n = satisNode
	}

	return nil

}

// 判断是否是通配符
func isWildSegment(segment string) bool {
	return strings.HasPrefix(segment, ":")
}

func (n *node) filterChildNodes(segment string) []*node {
	if len(n.children) == 0 {
		return nil
	}

	// 如果是通配符，所有子节点都满足
	if isWildSegment(segment) {
		return n.children
	}

	nodes := make([]*node, 0, len(n.children))
	for _, child := range n.children {
		if isWildSegment(child.segment) {
			nodes = append(nodes, child)
		} else if child.segment == segment {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

// 判断路由是否已经出现在树中
func (n *node) matchNode(uri string) *node {
	segments := strings.SplitN(uri, "/", 2)
	segment := segments[0]
	if !isWildSegment(segment) {
		segment = strings.ToLower(segment)
	}

	children := n.filterChildNodes(segment)
	if children == nil || len(children) == 0 { //如果没有子节点
		return nil
	}

	// 如果是最后一段uri
	if len(segments) == 1 {
		for _, child := range children {
			if child.isLast { // 如果子节点是最后一个节点
				return child
			}
		}
		return nil
	}

	// 如果没走完就继续走
	for _, child := range children {
		match := child.matchNode(segments[1])
		if match != nil {
			return match
		}
	}
	return nil
}
