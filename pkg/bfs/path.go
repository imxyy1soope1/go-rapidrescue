package bfs

import (
	"fmt"
	"strings"

	"github.com/imxyy1soope1/go-rapidrescue/pkg/constants"
)

type Path struct {
	path               []int
	LeftTurningPoints  []int
	RightTurningPoints []int
	Dest               int
	Graph              *Graph
}

func (p *Path) Len() int {
	return len(p.path) + len(p.LeftTurningPoints) + len(p.RightTurningPoints)
}

func (p *Path) String() string {
	builder := strings.Builder{}
	for _, v := range p.path {
		builder.WriteString(fmt.Sprintf("0x%04X ", v))
	}
	return builder.String()[:len(builder.String())-1]
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
	path := &Path{path: p, Graph: g, Dest: node.id}
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
		directions[i] = p.Graph.getDirection(p.path[i], p.path[i+1])
	}
	p.LeftTurningPoints = make([]int, 0)
	p.RightTurningPoints = make([]int, 0)
	for i := 1; i < len(directions); i++ {
		d := getTurningDirection(directions[i-1], directions[i])
		if d == constants.LEFT {
			p.LeftTurningPoints = append(p.LeftTurningPoints, p.path[i])
		} else if d == constants.RIGHT {
			p.RightTurningPoints = append(p.RightTurningPoints, p.path[i])
		}
	}
}
