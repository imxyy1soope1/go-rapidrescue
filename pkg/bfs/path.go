package bfs

import (
	"github.com/imxyy1soope1/go-rapidrescue/pkg/constants"
)

type Path struct {
	path               []int
	leftTurningPoints  []int
	rightTurningPoints []int
	graph              *Graph
}

func (p *Path) Len() int {
	return len(p.path)*2 + len(p.leftTurningPoints) + len(p.rightTurningPoints)
}

func GetPathFromNode(g *Graph, node *node) *Path {
	p := []int{node.id}
	for node.father != nil {
		p = append(p, node.father.id)
		node = node.father
	}
	// reverse
	for i, j := 0, len(p)-1; i < j; i, j = i+1, j-1 {
		p[i], p[j] = p[j], p[i]
	}
	path := &Path{path: p, graph: g}
	path.getTurningPoints()
	return path
}

func getTurningDirection(first, second constants.Direction) constants.RotaionDirection {
	var (
		flag      = first - second
		direction constants.RotaionDirection
	)
	if flag == 0 {
		return constants.STRAIGHT
	} else if flag >= 0 {
		direction = constants.LEFT
	} else {
		direction = constants.RIGHT
	}
	if flag == -3 {
		direction = constants.LEFT
	} else if flag == 3 {
		direction = constants.RIGHT
	}
	return direction
}

func (p *Path) getTurningPoints() {
	directions := make([]constants.Direction, len(p.path)-1)
	for i := 0; i < len(p.path)-1; i++ {
		directions[i] = p.graph.getDirection(p.path[i], p.path[i+1])
	}
	p.leftTurningPoints = make([]int, 0)
	p.rightTurningPoints = make([]int, 0)
	for i := 1; i < len(directions); i++ {
		d := getTurningDirection(directions[i-1], directions[i])
		if d == constants.LEFT {
			p.leftTurningPoints = append(p.leftTurningPoints, p.path[i])
		} else if d == constants.RIGHT {
			p.rightTurningPoints = append(p.rightTurningPoints, p.path[i])
		}
	}
}
