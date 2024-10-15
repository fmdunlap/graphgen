package graph

import (
	"fmt"
	"reflect"
)

type NodeId string

type Node struct {
	ID    NodeId
	Attrs NodeAttributes
}

type NodeAttributes struct {
	Label string
}

func (attrs NodeAttributes) String() string {
	return fmt.Sprintf(" [label: %v]", attrs.Label)
}

func (n Node) HasAttributes() bool {
	return !reflect.DeepEqual(n.Attrs, NodeAttributes{})
}

func (n Node) definitionString() string {
	nodeString := fmt.Sprintf("%v", n.ID)
	if n.HasAttributes() {
		nodeString += n.Attrs.String()
	}
	return nodeString
}
