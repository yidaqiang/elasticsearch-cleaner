package types

import "time"

type Indices []Index

func (idxs Indices) Len() int {
	return len(idxs)
}
func (idxs Indices) Swap(i, j int) {
	idxs[i], idxs[j] = idxs[j], idxs[i]
}
func (idxs Indices) Less(i, j int) bool {
	return idxs[j].CreateTime.Before(idxs[i].CreateTime)
}

type Index struct {
	Name       string
	CreateTime time.Time
}
