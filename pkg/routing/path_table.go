package routing

import (
	"github.com/imxyy1soope1/go-rapidrescue/pkg/bfs"
)

type pathTable struct {
	m     map[int]*bfs.Path
	graph *bfs.Graph
}

func newPathTable(graph *bfs.Graph) *pathTable {
	return &pathTable{m: make(map[int]*bfs.Path), graph: graph}
}

func (pt *pathTable) get(origin, dest int) *bfs.Path {
	if path, ok := pt.m[origin+dest*0x10000]; ok {
		return path
	} else {
		pt.m[origin+dest*0x10000] = pt.graph.Bfs(origin, dest)
	}
	return pt.m[origin+dest*0x10000]
}
