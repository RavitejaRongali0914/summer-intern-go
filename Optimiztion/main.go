package main

import (
	"sync"
	"time"
)

type Location struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}
type driverOp struct {
	Id            int
	last_location Location
	//Duration_work time.Duration
	last_trip_end time.Time
	vehicle_size  int
	//available     bool
	trips []int
}
type trip struct {
	tripid        int
	capacity      int
	starttime     time.Time
	endtime       time.Time
	startlocation Location
	endlocation   Location
	hasallocated  bool
	tripduration  float64
	driver        int
}

var D []int
var globalD []driverOp
var mutex = &sync.Mutex{}
var tripsbyshift = make(map[int][]trip)
var startshift1 = time.Date(2019, 9, 21, 0, 00, 00, 0, time.UTC)
var idmap = make(map[int]trip)

func main() {

	parse1()
	for id := 0; id < 48; id++ {
		for index, i := range globalD {
			if i.last_trip_end.Sub(startshift1).Minutes() < (float64(30*id + 10)) {
				D = append(D, index)
			}
		}
		if tripsbyshift[id] != nil {
			findavailabledrivers(tripsbyshift[id])
		}
		D = D[:0]
	}
	// var c int
	// for _, r := range globalD {
	// 	fmt.Println(r)
	// 	c = c + len(r.trips)

	// }
	// for i, _ := range tripsbyshift {
	// 	if tripsbyshift[i] != nil {
	// 		fmt.Println(tripsbyshift[i])
	// 	}
	// }
	csvtrips()
	csvdrivers()
	// for _, tt := range idmap {
	// 	fmt.Println(tt)
	// }
	//fmt.Println("total trips", c)

}
