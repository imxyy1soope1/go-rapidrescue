package routing

import (
	"fmt"
	"strings"

	"github.com/imxyy1soope1/go-rapidrescue/pkg/bfs"
	"github.com/imxyy1soope1/go-rapidrescue/pkg/constants"
	"github.com/imxyy1soope1/go-rapidrescue/pkg/json"
)

type p struct {
	path     *bfs.Path
	op       constants.Op
	goodsNum int
}

type Path []p

func (pth *Path) ToResult() json.Result {
	res := json.Result{}
	for _, v := range *pth {
		res.Path = append(res.Path, json.AtomPath{
			Dest:               v.path.Dest,
			LeftTurningPoints:  v.path.LeftTurningPoints,
			RightTurningPoints: v.path.RightTurningPoints,
			Op:                 v.op,
			GoodsNum:           v.goodsNum,
		})
	}
	pt := []p(*pth)
	res.Map = fmt.Sprint(pt[0].path.Graph)
	res.String = fmt.Sprint(pth)
	return res
}

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
		if v.op == constants.END {
			builder.WriteString("去往0x0A02，结束")
		} else {
			builder.WriteString(fmt.Sprintf("去往0x%04X, %s%d件物资\n", v.path.Dest, v.op, v.goodsNum))
		}
	}
	return builder.String()
}

func GetPathFromNode(pt *pathTable, n *node) *Path {
	path := Path{
		p{
			path:     pt.get(n.id, constants.DEST),
			op:       constants.END,
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
