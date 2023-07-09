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

func (rp *RoutePlanner) Plan() (path *Path) {
	queue := []*node{
		{
			id:       constants.ORIGIN,
			goodsNum: 0,
			op:       constants.HOLD,
			father:   nil,
			data:     rp.data,
		},
	}
	pathTable := newPathTable(rp.graph)
	cnt := 2000
	for len(queue) > 0 {
		n := queue[0]
		queue = queue[1:]
		for _, pt := range rp.nextPoints(n) {
			_ = pathTable.get(n.id, pt.id)
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
				goodsNum: goods,
				op:       constants.Op(pt.pointType),
				father:   n,
				data:     &newData,
			}
			if newData.RequiredGoods == 0 {
				cnt--
				if path != nil {
					newPath := GetPathFromNode(pathTable, newNode)
					if path.Len() > newPath.Len() {
						path = newPath
					}
					if cnt == 0 {
						return
					}
				} else {
					path = GetPathFromNode(pathTable, newNode)
				}
				continue
			}
			queue = append(queue, newNode)
		}
	}
	return
}
