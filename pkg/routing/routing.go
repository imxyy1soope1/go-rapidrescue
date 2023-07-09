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
	nodeChan := make(chan *node, 1)
	pathChan := make(chan *Path, 100)
	workerDoneChan := make(chan struct{})
	pathT := newPathTable(rp.graph)
	workerCount := 0
	cnt := 0
	nodeChan <- &node{
		id:       constants.ORIGIN,
		goodsNum: 0,
		op:       constants.HOLD,
		father:   nil,
		data:     rp.data,
	}

Loop:
	for {
		select {
		case newPath := <-pathChan:
			cnt++
			if path == nil || newPath.Len() < path.Len() {
				path = newPath
			}
			if cnt >= 10000 {
				break Loop
			}
		case <-workerDoneChan:
			workerCount--
			if workerCount == 0 {
				break Loop
			}
		case n := <-nodeChan:
			if workerCount > constants.MAX_WORKER_COUNT {
				continue
			}
			workerCount++
			go rp.newWorker(n, nodeChan, pathChan, workerDoneChan, pathT)
		}
	}
	return
}

func (rp *RoutePlanner) newWorker(n *node, nodeChan chan *node, pathChan chan *Path, workerDoneChan chan struct{}, pathT *pathTable) {
	for _, pt := range rp.nextPoints(n) {
		if pth := pathT.get(n.id, pt.id); pth == nil {
			continue
		}
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
			pathChan <- GetPathFromNode(pathT, newNode)
			continue
		}
		nodeChan <- newNode
	}

	workerDoneChan <- struct{}{}
}
