package main

import (
	"fmt"
	"sort"

	//"github.com/MOOVE-Network/clustering_core/clustering"
	//"github.com/MOOVE-Network/clustering_core/clustering"
	//"github.com/MOOVE-Network/clustering_core/clustering"
	//"github.com/MOOVE-Network/clustering_core/clustering"
	geo "github.com/paulmach/go.geo"
	"github.com/paulmach/go.geo/quadtree"
)

// type Location struct {
// 	Lat float64 `json:"lat"`
// 	Lng float64 `json:"lng"`
// }

type Cluster struct {
	//id             int
	sites []packet
	//hasVehicle       bool
	Vehicletype int
	crweight    float64
	crvolume    float64
	efficiency  float64
}

type vehicle struct {
	index  int
	weight float64
	volume float64
}

type packet struct {
	id             int     //`json:"id"`
	weight         float64 //`json:"weight"`
	volume         float64 //`json:"volume"`
	latitude       float64 //`json:"latitude"`
	longitude      float64 //`json:"longitute"`
	DistanceToSite float64 //`json:"distanceToSite"`
	done           bool
}

var tripsbystate = make(map[string][]packet)

//var Vehicles []vehicle
type By func(e1, e2 *packet) bool

type packetSorter struct {
	packets []packet
	by      func(e1, e2 *packet) bool
}

func distanceToSite(e1, e2 *packet) bool {
	return e1.DistanceToSite < e2.DistanceToSite
}

func (es *packetSorter) Len() int {
	return len(es.packets)
}

func (es *packetSorter) Less(i, j int) bool {
	return es.by(&es.packets[i], &es.packets[j])
}
func (es *packetSorter) Swap(i, j int) {
	es.packets[i], es.packets[j] = es.packets[j], es.packets[i]
}

func ImportPacketPoints(packets []packet) *quadtree.Quadtree {
	var pointers []geo.Pointer
	for _, e := range packets {
		pointers = append(pointers, geo.NewPoint(e.longitude, e.latitude))
	}
	qTree := quadtree.NewFromPointers(pointers)
	return qTree
}

func (by By) Sort(pp []packet, reverse bool) []packet {
	es := &packetSorter{
		packets: pp,
		by:      by,
	}
	if reverse {
		sort.Sort(sort.Reverse(es))
	} else {
		sort.Sort(es)
	}
	return es.packets
}

func GetVehicleToFit(weight float64, volume float64, v []vehicle, maxv vehicle) vehicle {
	//	var capSizes []int
	//var we, vo float64
	//fmt.Println("entered get", weight, volume)
	var minv, ans vehicle
	minv = maxv
	//minv.weight = weight
	//minv.volume = volume
	for _, k := range v {
		if k.weight >= weight && k.volume >= volume {
			if k.weight <= minv.weight && k.volume <= minv.volume {
				minv = k
				ans = k
				//	fmt.Println(ans, k)
			}
		}
	}

	return ans
}

func (e packet) Point() *geo.Point {
	return geo.NewPointFromLatLng(e.latitude, e.longitude)
}

func NewPointFromLatLng(lat, lng float64) *Point {
	return &Point{lng, lat}
}

func GetNextLargest(v []vehicle) vehicle {

	//var we, vo float64
	k := v[len(v)-1]
	return k
}
func (e *packet) ReInitDistance(from *geo.Point) {
	pp := geo.NewPoint(e.longitude, e.latitude)
	var k []geo.Pointer
	k = append(k, pp)
	k = append(k, from)
	routeResp, _ := getRoute(false, k)
	var newDistance float64
	newDistance = routeResp.Routes[0].Distance
	e.DistanceToSite = newDistance
}

var v1 = vehicle{1, 200, 90}
var v2 = vehicle{2, 450, 67}
var v3 = vehicle{3, 780, 89}
var Vehicles = []vehicle{v1, v2, v3}
var mclusters []Cluster
var notclusteredpac []packet

func main() {

	parse()
	done := make(chan bool)
	for _, z := range tripsbystate {
		go func(f []packet) {
			var newstate []packet
			for _, e := range f {
				e.ReInitDistance(geo.NewPoint(72.867851, 19.155148))
				//fmt.Println(e.DistanceToSite)

				newstate = append(newstate, e)
			}
			//fmt.Println(newstate)
			qtree := ImportPacketPoints(newstate)
			c := ClusterizeProximity(qtree, newstate, Vehicles)
			for _, cc := range c {

				if cc.sites != nil {
					mclusters = append(mclusters, cc)
					//fmt.Println("---------------------")
					fmt.Println(cc)
				}
			}

			//fmt.Println("------------------------------------------------------------------------")
			done <- true
		}(z)
	}
	for i := 0; i < len(tripsbystate); i++ {
		<-done
	}
	csvdrivers()
	//var clusters []Cluster
	//for _, f := range newstate {

}
