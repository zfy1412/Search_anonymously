package send

type User struct {
	Id       string
	Name     string
	Password string
}
type Node struct {
	X float64 `json:"lng"`
	Y float64 `json:"lat"`
}
type Date struct {
	Stime string `json:"stime"`
	Etime string `json:"etime"`
	Node  []Node `json:"path"`
}
