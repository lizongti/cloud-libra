package cluster

import (
	"fmt"

	"github.com/aceaura/libra/cluster/component"
)

type Node struct {
	components []component.Component
}

func NewNode(opts ...nodeOpt) *Node {
	n := &Node{}

	for _, opt := range opts {
		opt(n)
	}

	return n
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
	return func(n *Node) {
		n.WithComponent(c)
	}
}

func (n *Node) WithComponent(c component.Component) *Node {
	n.components = append(n.components, c)
	return n
}
