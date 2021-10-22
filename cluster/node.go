package cluster

import (
	"fmt"

	"github.com/aceaura/libra/cluster/component"
)

type Node struct {
	opts       []nodeOpt
	components []component.Component
}

func NewNode(opts ...nodeOpt) *Node {
	return &Node{opts: opts}
}

func (n *Node) Boot() error {
	for _, s := range n.components {
		fmt.Println(s)
	}
	return nil
}

type nodeOpt func(*Node)
type nodeOption struct{}

var NodeOption nodeOption

func (nodeOption) WithComponent(c component.Component) nodeOpt {
	return func(s *Node) {
		s.components = append(s.components, c)
	}
}

func (n *Node) WithComponent(c component.Component) *Node {
	n.opts = append(n.opts, NodeOption.WithComponent(c))
	return n
}
