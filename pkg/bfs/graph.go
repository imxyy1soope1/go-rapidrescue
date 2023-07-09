package bfs

import (
	"fmt"
	"math"
	"strings"

	"github.com/imxyy1soope1/go-rapidrescue/pkg/constants"
)

type road = [2]int

type Graph struct {
	graph           [][]int
	links           map[int][]int
	brokenRoads     []road
	nonTuringPoints []int
}

func (g *Graph) String() string {
	builder := strings.Builder{}
	for _, line := range g.graph {
		for _, point := range line {
			if point == int(constants.ROAD) {
				builder.WriteString("-0x001")
			} else {
				builder.WriteString(fmt.Sprintf("0x%04X ", point))
			}
		}
		builder.WriteString("\n")
	}
	return builder.String()
}

func NewGraph(graph [][]int, brokenRoads []road, nonTuringPoints []int) *Graph {
	g := &Graph{}
	g.graph = graph
	g.links = make(map[int][]int)
	g.brokenRoads = brokenRoads
	g.nonTuringPoints = nonTuringPoints
	g.getLinks()
	return g
}

func (g *Graph) getLinks() {
	for x, line := range g.graph {
		for y, point := range line {
			if point <= int(constants.BLOCK) {
				continue
			}
			g.links[point] = make([]int, 0)
			direct := [3]int{-1, 0, 1}

			for _, i := range direct {
				for _, j := range direct {
					if math.Abs(float64(i)) == math.Abs(float64(j)) {
						continue
					}
					newX := x + i
					newY := y + j
					for g.pointIsValid(newX, newY) {
						if g.graph[newX][newY] > 0 {
							g.links[point] = append(g.links[point], g.graph[newX][newY])
							break
						}
						newX += i
						newY += j
					}
				}
			}
		}
	}
}

func (g *Graph) isBroken(rd road) bool {
	for _, p := range g.brokenRoads {
		if (rd[0] == p[0] && rd[1] == p[1]) || (rd[0] == p[1] && rd[1] == p[0]) {
			return true
		}
	}
	return false
}

func (g *Graph) pointIsValid(x, y int) bool {
	return 0 <= x && x < len(g.graph) && 0 <= y && y < len(g.graph[0]) && constants.PointType(g.graph[x][y]) != constants.BLOCK
}

func (g *Graph) neighbers(point int) []int {
	ret := make([]int, 0, len(g.links[point]))
	for _, p := range g.links[point] {
		if !g.isBroken(road{point, p}) {
			ret = append(ret, p)
		}
	}
	return ret
}

func (g *Graph) getDirection(first, second int) constants.Direction {
	x1, y1 := -1, -1
	x2, y2 := -1, -1
	for x, line := range g.graph {
		for y, point := range line {
			if point == first {
				x1, y1 = x, y
			} else if point == second {
				x2, y2 = x, y
			}
		}
	}
	if x1 == x2 {
		if y1 < y2 {
			return constants.E
		} else {
			return constants.W
		}
	} else {
		if x1 < x2 {
			return constants.S
		} else {
			return constants.N
		}
	}
}

func (g *Graph) nonTuring(id int) bool {
	for _, p := range g.nonTuringPoints {
		if p == id {
			return true
		}
	}
	return false
}

func (g *Graph) Bfs(origin, dest int) (path *Path) {
	queue := []*node{
		{
			id:        origin,
			direction: constants.UNDEF_DIRECTION,
			father:    nil,
		},
	}
	for _, r := range g.brokenRoads {
		if r[0] == origin || r[1] == origin || r[0] == dest || r[1] == dest {
			return
		}
	}
	cnt := 4
	for len(queue) > 0 {
		n := queue[0]
		queue = queue[1:]
		for _, neigh := range g.neighbers(n.id) {
			if n.father != nil && neigh == n.father.id {
				continue
			}
			direction := g.getDirection(n.id, neigh)
			nd := &node{
				id:        neigh,
				direction: direction,
				father:    n,
			}
			if neigh == dest {
				cnt--
				if path != nil {
					newpath := GetPathFromNode(g, nd)
					if path.Len() > newpath.Len() {
						path = newpath
					}
					if cnt == 0 {
						return
					}
				} else {
					path = GetPathFromNode(g, nd)
				}
				continue
			}
			queue = append(queue, nd)
		}
	}
	return
}
