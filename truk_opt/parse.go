package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	//"time"
)

func parse() {
	csvfile, err := os.Open("errorTo.csv")

	if err != nil {

		log.Fatalln("Couldn't open the csv file", err)
	}

	// Parse the file
	r := csv.NewReader(csvfile)

	//var trips []trip
	i := 0

	//var tripsbystate = make(map[string][]packet)

	// Iterate through the records
	for {

		if i == 0 {
			_, err := r.Read()

			if err == io.EOF {

				break
			}

			if err != nil {
				log.Fatal(err)
			}

		}
		if i > 0 {

			record, err := r.Read()

			if err == io.EOF {

				break
			}

			if err != nil {
				log.Fatal(err)
			}

			id, _ := strconv.Atoi(record[1])

			state := record[8]
			//fmt.Println(layout, start)
			lat, _ := strconv.ParseFloat(record[13], 64)
			lang, _ := strconv.ParseFloat(record[14], 64)

			//r2, _ := strconv.Atoi(record[1])
			weight, _ := strconv.ParseFloat(record[23], 64)
			volume, _ := strconv.ParseFloat(record[25], 64)

			_, ok := tripsbystate[state]
			if !ok {
				tripsbystate[state] = make([]packet, 0)
			}
			tripsbystate[state] = append(tripsbystate[state], packet{id, weight, volume, lat, lang, 0, false})

		}

		i = i + 1

	}

	// for state1, i := range tripsbystate {
	// 	for _, j := range i {
	// 		fmt.Println(state1, " ", j.id)
	// 	}
	// }
}
func csvdrivers() {

	file, err := os.Create("csvtrips.csv")
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()
	var w [][]string
	w = append(w, []string{"tripId", "typeofvehicle", "weight", "volume", "efficiency in %", "Packet Id's"})
	for l, data := range mclusters {
		l = l + 1
		//s := data.last_trip_end.String()
		//e := data1.endtime.String()
		wei := fmt.Sprintf("%f", data.crweight)
		vol := fmt.Sprintf("%f", data.crvolume)
		eff := fmt.Sprintf("%f", data.efficiency)
		//gg := strings.Trim(strings.Replace(fmt.Sprint(data.trips), " ", ",", -1), "[]")
		var gg string
		//var totaldis float64
		for i, kk := range data.sites {
			//totaldis = totaldis + idmap[kk].tripduration
			gg = gg + strconv.Itoa(kk.id)
			if i+1 < len(data.sites) {
				gg = gg + ","
			}
		}
		//d := fmt.Sprintf("%f", totaldis)
		//fmt.Println(i, gg)
		w = append(w, []string{strconv.Itoa(l), strconv.Itoa(data.Vehicletype), wei, vol, eff, gg})

	}
	for j, cc := range notclusteredpac {
		j = j + 1
		v := len(mclusters)
		wei := fmt.Sprintf("%f", cc.weight)
		vol := fmt.Sprintf("%f", cc.volume)
		// 	eff := fmt.Sprintf("%f", data.efficiency)
		w = append(w, []string{strconv.Itoa(j + v), "need large vehicle", wei, vol, "+iff", " "})
	}
	for _, dd := range w {
		err = writer.Write(dd)
		if err != nil {
			log.Fatalln("Couldn't open the csv file", err)
		}
	}

}
