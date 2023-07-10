package json

type AtomPath struct {
	Dest               int   `json:"dest"`
	LeftTurningPoints  []int `json:"left_turning_points"`
	RightTurningPoints []int `json:"right_turning_points"`
	GoodsNum           int   `json:"goods_num"`
}

type Result struct {
	Path   []AtomPath `json:"path"`
	Map    string     `json:"map"`
	String string     `json:"string"`
}
