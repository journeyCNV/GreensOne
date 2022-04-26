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
	parent   *node
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
		parent:   &node{},
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
	// 先判断树里有没有这个uri
	if n.matchNode(uri) != nil {
		return errors.New("route exist:" + uri)
	}

	segments := strings.Split(uri, "/")
	for i, seg := range segments {
		if !isWildSegment(seg) {
			seg = strings.ToLower(seg)
		}
		isLast := i == len(segments)-1

		// 先找找有没有匹配的节点
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

		// 如果没有匹配的节点，进行新建
		if satisNode == nil {
			currNode := NewNode()
			currNode.segment = seg
			currNode.parent = n
			// 如果已经遍历到uri的最后一段了，存储一下
			if isLast {
				currNode.isLast = true
				currNode.handlers = handlers
			}
			n.children = append(n.children, currNode)
			satisNode = currNode // 将新建当前节点赋值给(满足条件的)记录节点
		}

		// 更换上级节点，往下走
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

// 将uri解析为params
func (n *node) parseParamsFormEndNode(uri string) map[string]string {
	currNode := n
	paramMap := map[string]string{}
	segments := strings.Split(uri, "/")
	segL := len(segments) - 1
	for i := segL; i >= 0; i-- {
		if currNode.segment == "" {
			break
		}
		if isWildSegment(currNode.segment) {
			paramMap[currNode.segment[1:]] = segments[i] //去掉通配符的:
		}
		currNode = currNode.parent
	}
	return paramMap
}
