package routing

import (
	"math"

	"github.com/imxyy1soope1/go-rapidrescue/pkg/bfs"
	"github.com/imxyy1soope1/go-rapidrescue/pkg/constants"
)

type RoutePlanner struct {
	graph *bfs.Graph
	data  *Data
}

func NewRoutePlanner(graph *bfs.Graph, data *Data) *RoutePlanner {
	return &RoutePlanner{graph: graph, data: data}
}

func (rp *RoutePlanner) nextPoints(node *node) []point {
	ret := make([]point, 0)
	for id, g := range node.data.MaterialPoints {
		if g > 0 && node.data.Carrying != constants.MAX_HOLDING_GOODS && id != node.id {
			ret = append(ret, point{id, constants.MATERIAL})
		}
	}
	for id, g := range node.data.Quarters {
		if g > 0 && node.data.Carrying > 0 && id != node.id {
			ret = append(ret, point{id, constants.QUARTER})
		}
	}
	return ret
}

func (rp *RoutePlanner) Plan() Path {
	pq := bfs.NewPQ()
	n := &node{
		id:       constants.ORIGIN,
		g:        0,
		goodsNum: 0,
		op:       constants.HOLD,
		father:   nil,
		data:     rp.data,
	}
	pq.Push(&bfs.Item{
		Value:    n,
		Priority: 0,
		Index:    0,
	})
	pathTable := newPathTable(rp.graph)
	currIndex := 1
	for pq.Len() > 0 {
		n := bfs.Pop(pq).Value.(*node)
		currIndex--
		for _, pt := range rp.nextPoints(n) {
			path := pathTable.get(n.id, pt.id)
			var goods int
			if pt.pointType == constants.MATERIAL {
				goods = int(math.Min(math.Min(
					float64(constants.MAX_HOLDING_GOODS-n.data.Carrying),
					float64(n.data.MaterialPoints[pt.id]),
				), float64(n.data.RequiredGoods-n.data.Carrying)))
			} else {
				goods = int(math.Min(math.Min(
					float64(n.data.Carrying),
					float64(n.data.Quarters[pt.id]),
				), float64(n.data.RequiredGoods)))
			}
			if goods <= 0 {
				continue
			}
			newData := *n.data
			newData.MaterialPoints = make(map[int]int)
			newData.Quarters = make(map[int]int)
			for id, g := range n.data.MaterialPoints {
				newData.MaterialPoints[id] = g
			}
			for id, g := range n.data.Quarters {
				newData.Quarters[id] = g
			}
			if pt.pointType == constants.MATERIAL {
				newData.MaterialPoints[pt.id] -= goods
				newData.Carrying += goods
			} else {
				newData.Quarters[pt.id] -= goods
				newData.Carrying -= goods
				newData.RequiredGoods -= goods
			}
			newNode := &node{
				id:       pt.id,
				g:        path.Len()*18 - goods,
				goodsNum: goods,
				op:       constants.Op(pt.pointType),
				father:   n,
				data:     &newData,
			}
			if newData.RequiredGoods == 0 {
				return GetPathFromNode(pathTable, newNode)
			}
			pq.Push(&bfs.Item{
				Value:    newNode,
				Priority: 0,
				Index:    currIndex,
			})
			currIndex++
		}
	}
	return nil
}
