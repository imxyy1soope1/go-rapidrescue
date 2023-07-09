package constants

import "runtime"

type PointType int

const (
	ROAD PointType = iota - 1
	BLOCK
	MATERIAL
	QUARTER
)

func IsBlock(pointType int) bool {
	return PointType(pointType) == BLOCK
}

type Direction int8

const (
	UNDEF_DIRECTION Direction = iota - 1
	N
	E
	S
	W
)

type RotaionDirection uint8

const (
	STRAIGHT RotaionDirection = iota
	LEFT
	RIGHT
)

const MAX_HOLDING_GOODS = 18

type Op uint8

const (
	UNDEF_OP Op = iota
	HOLD
	RELEASE
)

func (op Op) String() string {
	switch op {
	case HOLD:
		return "获取"
	case RELEASE:
		return "释放"
	default:
		return "unknown"
	}
}

const (
	ORIGIN = iota + 0x0A01
	DEST
)

var MAX_WORKER_COUNT = runtime.NumCPU()
