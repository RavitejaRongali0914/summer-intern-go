package main

import (
	"fmt"

	geo "github.com/paulmach/go.geo"
)

func findavailabledrivers(m []trip) {
	//var wg sync.WaitGroup
	//wg.Add(len(m))
	//fmt.Println("strating", len(m))
	done := make(chan bool)
	for i, n := range m {
		go func(k1 trip, j int) {
			//k1 := n
			p1 := geo.NewPoint(k1.startlocation.Lng, k1.startlocation.Lat)
			end := geo.NewPoint(k1.endlocation.Lng, k1.endlocation.Lat)
			var x []geo.Pointer
			x = append(x, p1)
			x = append(x, end)
			routeRespe, _ := getRoute(false, x)
			distance1 := routeRespe.Routes[0].Distance

			if D != nil {
				for _, q := range D {
					var p []geo.Pointer
					p = append(p, p1)
					p2 := geo.NewPoint(globalD[q].last_location.Lng, globalD[q].last_location.Lat)
					p = append(p, p2)
					routeResp, _ := getRoute(false, p)
					distance := routeResp.Routes[0].Distance
					if distance < 25000 {
						mutex.Lock()
						if k1.starttime.Sub(globalD[q].last_trip_end).Seconds() > 1200 {
							//fmt.Println("Found old driver", D[q].Id, D[q].trips)
							//fmt.Println(k1.tripid, k.Id, distance)
							//mutex.Lock()
							globalD[q].trips = append(globalD[q].trips, k1.tripid)
							globalD[q].last_location = k1.endlocation
							globalD[q].last_trip_end = k1.endtime
							k1.hasallocated = true
							//fmt.Println(D[q])
							fmt.Println(distance1)
							m[j].tripduration = distance1
							m[j].hasallocated = true
							m[j].driver = globalD[q].Id
							idmap[m[j].tripid] = m[j]

						}
						mutex.Unlock()
						//mutex.Unlock()
						if k1.hasallocated {
							//fmt.Println(k1.tripid)
							break
						}
					}

				}
			}
			if !k1.hasallocated {
				//fmt.Println("assigning new driver for ", k1.tripid)
				var z []int
				z = append(z, k1.tripid)
				newdriver := driverOp{len(globalD) + 1, k1.endlocation, k1.endtime, k1.capacity, z}
				globalD = append(globalD, newdriver)
				k1.hasallocated = true
				mutex.Lock()
				k1.tripduration = distance1
				m[j].hasallocated = true
				m[j].tripduration = distance1
				m[j].driver = newdriver.Id
				k1.driver = newdriver.Id
				idmap[m[j].tripid] = k1
				mutex.Unlock()
				//fmt.Println(D)
				//fmt.Println(e)

			}
			done <- true
		}(n, i)
	}
	for i := 0; i < len(m); i++ {
		<-done
	}
}
