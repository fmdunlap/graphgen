package graph

import (
	"fmt"
	"reflect"
)

type Edge struct {
	Source NodeId
	Target NodeId
	Attrs  EdgeAttributes
}

type EdgeAttributes struct {
	Label  string
	Weight float64
}

func (attrs EdgeAttributes) String() string {
	if attrs.isEmpty() {
		return ""
	}
	seenFirst := false

	attrsString := "["
	if attrs.Label != "" {
		attrsString += fmt.Sprintf("label: %v", attrs.Label)
		seenFirst = true
	}

	if attrs.Weight != 0.0 {
		if seenFirst {
			attrsString += ", "
		}
		attrsString += fmt.Sprintf("weight: %v", attrs.Weight)
		seenFirst = true
	}

	attrsString += "]"
	return attrsString
}

func (attrs EdgeAttributes) isEmpty() bool {
	return reflect.DeepEqual(attrs, EdgeAttributes{})
}

func (e *Edge) HasAttributes() bool {
	return !e.Attrs.isEmpty()
}

func (e *Edge) definitionString() string {
	edgeString := fmt.Sprintf("%v %v %v", e.Source, "->", e.Target)
	if e.HasAttributes() {
		edgeString += fmt.Sprintf(" %v", e.Attrs.String())
	}
	return edgeString
}
