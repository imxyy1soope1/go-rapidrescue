package json

type Config struct {
	Map             [][]int     `json:"map"`
	BrokenRoads     [][2]int    `json:"broken_roads"`
	NonTuringPoints []int       `json:"non_turing_points"`
	RequiredGoods   int         `json:"required_goods"`
	MaterialPoints  map[int]int `json:"material_points"`
	Quarters        map[int]int `json:"quarters"`
}
