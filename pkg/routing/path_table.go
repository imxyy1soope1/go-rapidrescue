package routing

import (
	"sync"

	"github.com/imxyy1soope1/go-rapidrescue/pkg/bfs"
)

type pathTable struct {
	m     sync.Map
	graph *bfs.Graph
}

func newPathTable(graph *bfs.Graph) *pathTable {
	return &pathTable{m: sync.Map{}, graph: graph}
}

func (pt *pathTable) get(origin, dest int) *bfs.Path {
	if path, ok := pt.m.Load(origin + dest*0x10000); ok {
		return path.(*bfs.Path)
	} else {
		path := pt.graph.Bfs(origin, dest)
		pt.m.Store(origin+dest*0x10000, path)
		return path
	}
}
