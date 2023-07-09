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
	END Op = iota
	HOLD
	RELEASE
)

func (op Op) String() string {
	switch op {
	case END:
		return "结束"
	case HOLD:
		return "装载"
	case RELEASE:
		return "卸下"
	default:
		return "未知"
	}
}

const (
	ORIGIN = iota + 0x0A01
	DEST
)

var MAX_WORKER_COUNT = runtime.NumCPU()
