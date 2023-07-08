package bfs

import (
	"github.com/imxyy1soope1/go-rapidrescue/pkg/constants"
)

type node struct {
	id        int
	direction constants.Direction
	father    *node
}
