package json

import (
	"github.com/imxyy1soope1/go-rapidrescue/pkg/constants"
)

type AtomPath struct {
	Dest               int
	LeftTurningPoints  []int
	RightTurningPoints []int
	Op                 constants.Op
	GoodsNum           int
}

type Result struct {
	Path   []AtomPath `json:"path"`
	Map    string     `json:"map"`
	String string     `json:"string"`
}
