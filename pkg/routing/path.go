package routing

import (
	"fmt"
	"strings"

	"github.com/imxyy1soope1/go-rapidrescue/pkg/bfs"
	"github.com/imxyy1soope1/go-rapidrescue/pkg/constants"
)

type p struct {
	path     *bfs.Path
	op       constants.Op
	goodsNum int
}

type Path []p

func (p *Path) Len() int {
	sum := 0
	for _, v := range *p {
		sum += v.path.Len()
	}
	return sum
}

func (p *Path) String() string {
	builder := strings.Builder{}
	for _, v := range *p {
		builder.WriteString(fmt.Sprintf("%s %s %d\n", v.path, v.op, v.goodsNum))
	}
	return builder.String()
}

func GetPathFromNode(pt *pathTable, n *node) *Path {
	path := Path{
		p{
			path:     pt.get(n.id, constants.DEST),
			op:       constants.UNDEF_OP,
			goodsNum: 0,
		},
	}
	for n.father != nil {
		path = append(path, p{
			path:     pt.get(n.father.id, n.id),
			op:       n.op,
			goodsNum: n.goodsNum,
		})
		n = n.father
	}
	// reverse
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return &path
}
