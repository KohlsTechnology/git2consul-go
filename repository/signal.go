package repository

type Signal struct {
	Type   string
	Branch string
}

// TODO: Possibly use enums for the type
// type signal uint
//
// const (
// 	CLONE signal = iota
// 	UPDATE
// )
//
// var signals = []string{
// 	"CLONE",
// 	"UPDATE",
// }
//
// func (s signal) String() string {
// 	return signals[s]
// }

func (r *Repository) GetSignal() Signal {
	return <-r.signal
}
