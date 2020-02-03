package main

import (
	"fmt"

	geo "github.com/paulmach/go.geo"
	"github.com/paulmach/go.geo/quadtree"
)

// type Pointer interface {
// 	// CenterPoint is kind of a weird name, but it's meant to not overlap
// 	// with any stuct attributes.
// 	CenterPoint() *geo.Point // lng/lat, or other, point
// }

func ClusterizeProximity(qTree *quadtree.Quadtree, packets []packet, vehicles []vehicle) []Cluster {
	var clusters []Cluster
	//	cpyEmployees := append([]Employee(nil), employees...)
	//	cpyEmployees = By(distanceToSite).Sort(cpyEmployees, false)
	//	clusterized := make(map[Employee]bool)
	cpypackets := append([]packet(nil), packets...)
	cpypackets = By(distanceToSite).Sort(cpypackets, false)
	clusterized := make(map[packet]bool)

	//fmt.Println(cpypackets)

	for _, p := range cpypackets {
		if clusterized[p] {
			continue
		}
		//fmt.Println(23)

		v := GetNextLargest(vehicles)

		var weight float64
		var volume float64
		//fmt.Println(23)
		//		pp := geo.NewPoint(p.longitude, p.latitude)
		kPointers := qTree.FindKNearestMatchingGeo(p.Point(), 50, func(r geo.Pointer) bool {
			var q packet
			for _, k := range packets {
				if r.Point()[0] == k.longitude && r.Point()[1] == k.latitude {
					q = k
				}
			}
			_, ok := clusterized[q]
			return !ok
		}, 10000)

		//fmt.Println("crossed")

		var tmppacks []packet

		for _, p := range kPointers {
			var emp packet
			var j int
			for i, k := range packets {
				if p.Point()[0] == k.longitude && p.Point()[1] == k.latitude {
					emp = k
					j = i
				}
			}
			if weight <= v.weight && volume <= v.volume && emp.weight < v.weight && emp.volume < v.volume {

				weight += emp.weight
				volume += emp.volume
				clusterized[emp] = true
				packets[j].done = true
				emp.done = true
				tmppacks = append(tmppacks, emp)
				//fmt.Println(emp.id)
			}
		}
		//fmt.Println(v)
		smallerV := GetVehicleToFit(weight, volume, vehicles, v)
		//fmt.Println(weight, volume, smallerV, v)
		//fmt.Println(smallerV)
		perct := (weight / smallerV.weight) * 100
		var cluster Cluster
		cluster = Cluster{tmppacks, smallerV.index, weight, volume, perct}
		clusters = append(clusters, cluster)

		/// for the employees obatined from FindKNearestMatchingGeo we update clusterized[employee] == true
		// and create a tempemps which contain the employees obatined from FindKNearestMatchingGeo.

	}
	for _, pac := range packets {
		if !pac.done {
			fmt.Println("contact manager", pac)
			notclusteredpac = append(notclusteredpac, pac)
		}
	}

	///building the clusters from this values.

	return clusters
}
