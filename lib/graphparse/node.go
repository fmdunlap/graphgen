package graphparse

type NodeId string

type Node struct {
    ID    NodeId
    Attrs NodeAttributes
}

type NodeAttributes struct {
    Label string
}
