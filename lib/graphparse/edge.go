package graphparse

type Edge struct {
    Source NodeId
    Target NodeId
    Attrs  EdgeAttributes
}

type EdgeAttributes struct {
    Label  string
    Weight float64
}
