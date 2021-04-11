package flat

type Building struct {
	Address string
	FlatCount int16
}

type Flat struct {
	Number uint16
	Building Building
}

// Housekeeping app
//
// type Building struct {
//	Address string
//	FlatCount int16
//	EntranceCount int8
// }
//
// type Entrance struct {
//	FlatCount int16
//	FloorCount int8
//	FlatOnFloorCount int8
// }
//
// type Flat struct {
//	Number uint16
//	Entrance Entrance
//	Building Building
// }
