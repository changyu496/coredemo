package framework

import (
	"errors"
	"strings"
)

type Tree struct {
	root *node // 根节点
}

type node struct {
	isLast   bool                // 代表这个节点是否可以成为最终路由规则，该节点是否能成为一个独立的uri，是否自身就是一个终极节点
	segment  string              // uri中的字符串，代表这个节点表示的路由中某个段的字符串
	handlers []ControllerHandler // 代表这个节点包含的控制器，用于最终加载调用
	childs   []*node             // 代表这个节点下的子节点
}

func newNode() *node {
	return &node{
		isLast:  false,
		segment: "",
		childs:  []*node{},
	}
}

func NewTree() *Tree {
	return &Tree{root: newNode()}
}

// 判断一个segment是否是通用的segment，以:开头
func isWildSegment(segment string) bool {
	return strings.HasPrefix(segment, ":")
}

// 过滤下一层满足segment规则的子节点
func (n *node) filterChildNodes(segment string) []*node {
	if len(n.childs) == 0 {
		return nil
	}
	//  如果segment是通配符，则所有下一层子节点都满足需求
	if isWildSegment(segment) {
		return n.childs
	}
	nodes := make([]*node, 0, len(n.childs))
	// 过滤所有下一层子节点
	for _, cnode := range n.childs {
		if isWildSegment(cnode.segment) {
			nodes = append(nodes, cnode)
		} else if segment == cnode.segment {
			nodes = append(nodes, cnode)
		}
	}
	return nodes
}

// 判断路由是否已经在节点的所有子节点树中存在了
func (n *node) matchNode(url string) *node {
	// 使用分隔符将uri切割为两个部分
	segments := strings.SplitN(url, "/", 2)
	segment := segments[0]
	if !isWildSegment(segment) {
		segment = strings.ToUpper(segment)
	}
	// 匹配符合的下一层子节点
	cnodes := n.filterChildNodes(segment)
	// 如果档期子节点没有一个符合，那么说明这个uri之前一定不存在，直接返回nil
	if cnodes == nil || len(cnodes) == 0 {
		return nil
	}
	// 如果只有一个segment，则是最后一个标记
	if len(segments) == 1 {
		// 如果segment已经是最后一个节点，判断这些cnode是否有isLast标志
		for _, tn := range cnodes {
			if tn.isLast {
				return tn
			}
		}
		// 都不是最后一个节点
		return nil
	}
	//  如果有多个segment，则需要依次递归每个子节点继续进行查找
	for _, tn := range cnodes {
		tnMatch := tn.matchNode(segments[1])
		if tnMatch != nil {
			return tnMatch
		}
	}
	return nil
}

// AddRoute 增加路由节点
func (tree *Tree) AddRoute(uri string, handlers []ControllerHandler) error {
	n := tree.root
	// 确认是否有路由冲突
	if n.matchNode(uri) != nil {
		return errors.New("route exists:" + uri)
	}
	segments := strings.Split(uri, "/")
	for index, segment := range segments {
		// 最终进入Node segment的字段
		if !isWildSegment(segment) {
			segment = strings.ToUpper(segment)
		}
		isLast := index == len(segments)-1

		var objNode *node //标记是否有合适的子节点
		childNodes := n.filterChildNodes(segment)
		// 如果有匹配的子节点
		if len(childNodes) > 0 {
			for _, cnode := range childNodes {
				if cnode.segment == segment {
					objNode = cnode
					break
				}
			}
		}
		// 如果没有找到，则创建一个新的node节点
		if objNode == nil {
			cnode := newNode()
			cnode.segment = segment
			if isLast {
				cnode.isLast = true
				cnode.handlers = handlers
			}
			n.childs = append(n.childs, cnode)
			objNode = cnode
		}
		n = objNode
	}
	return nil
}

func (tree *Tree) FindHandler(url string) []ControllerHandler {
	matchNode := tree.root.matchNode(url)
	if matchNode != nil {
		return matchNode.handlers
	}
	return nil
}
