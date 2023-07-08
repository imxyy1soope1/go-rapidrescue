package routing

import (
	"github.com/imxyy1soope1/go-rapidrescue/pkg/constants"
)

type point struct {
	id        int
	pointType constants.PointType
}

type Data struct {
	Carrying       int
	RequiredGoods  int
	MaterialPoints map[int]int
	Quarters       map[int]int
}

type node struct {
	id       int
	g        int
	goodsNum int
	op       constants.Op
	father   *node
	data     *Data
}
