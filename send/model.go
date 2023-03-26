package send

type User struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
}
type Node struct {
	X float64 `json:"lng"`
	Y float64 `json:"lat"`
}
type Trajectory struct {
	Id  int     `json:"id"`
	Tid int     `json:"tid"`
	X   float64 `json:"x"`
	Y   float64 `json:"y"`
}
type Date struct {
	Stime string `json:"stime"`
	Etime string `json:"etime"`
	Node  []Node `json:"path"`
}
