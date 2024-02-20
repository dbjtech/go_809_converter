package eventer

var LocateTypes = struct {
	CELL int
	GPS  int
	WIFI int
}{
	CELL: 0,
	GPS:  1,
	WIFI: 2,
}

var TermPositionTypes = struct {
	UNKNOWN int
	GPS     int
	CELL    int
}{
	UNKNOWN: 0,
	GPS:     1,
	CELL:    2,
}

var LcFixStatus = struct {
	FAILED  int
	SUCCESS int
}{
	FAILED:  0,
	SUCCESS: 1,
}

var MoveStatus = struct {
	RUN  int
	STOP int
}{
	RUN:  1,
	STOP: 2,
}

var TermStatus = struct {
	RUN  int
	STOP int
}{
	RUN:  2,
	STOP: 1,
}

var LocationType = struct {
	RUN  string
	STOP string
}{
	RUN:  "1",
	STOP: "2",
}

var CarStatus = struct {
	RUN  int
	STOP int
}{
	RUN:  2,
	STOP: 1,
}

var PushAcountFilter = struct {
	HOST_ACOUNT []int
	MONITOR_IN  []int
	WIRELESS    []int
	WIRED       []int
}{
	HOST_ACOUNT: []int{1},
	MONITOR_IN:  []int{1, 2, 3, 4},
	WIRELESS:    []int{1, 4},
	WIRED:       []int{1, 3},
}
