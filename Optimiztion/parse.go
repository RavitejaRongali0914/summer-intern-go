package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

func parse1() {
	csvfile, err := os.Open("a.csv")
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}

	// Parse the file
	r := csv.NewReader(csvfile)
	//var trips []trip

	// Iterate through the records
	for {
		record, err := r.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err)
		}

		startshift := time.Date(2019, 9, 21, 0, 00, 00, 0, time.UTC)
		r1, _ := strconv.Atoi(record[0])
		r2, _ := strconv.Atoi(record[1])
		layout := "2006-01-02 15:04"
		start, _ := time.Parse(layout, record[2]+" "+record[3])
		end, _ := time.Parse(layout, record[2]+" "+record[4])
		startlat, _ := strconv.ParseFloat(record[5], 64)
		startlang, _ := strconv.ParseFloat(record[6], 64)
		endlat, _ := strconv.ParseFloat(record[7], 64)
		endlang, _ := strconv.ParseFloat(record[8], 64)
		diff := start.Sub(startshift).Minutes()
		shift := int(diff) / int(30)
		_, ok := tripsbyshift[shift]
		if !ok {
			tripsbyshift[shift] = make([]trip, 0)
		}
		tripsbyshift[shift] = append(tripsbyshift[shift], trip{r1, r2, start, end, Location{startlat, startlang}, Location{endlat, endlang}, false, float64(0), 0})

	}
}
func csvtrips() {
	file, err := os.Create("csvtrips.csv")
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()
	var w [][]string
	for _, data := range tripsbyshift {
		if data != nil {
			for _, data1 := range data {
				s := data1.starttime.String()
				e := data1.endtime.String()
				sl := fmt.Sprintf("%f", data1.startlocation)
				el := fmt.Sprintf("%f", data1.endlocation)
				d := fmt.Sprintf("%f", data1.tripduration)
				w = append(w, []string{strconv.Itoa(data1.tripid), strconv.Itoa(data1.capacity), s, e, sl, el, strconv.FormatBool(data1.hasallocated), d + "m", "cab " + strconv.Itoa(data1.driver)})
			}
		}

	}
	for _, dd := range w {
		err = writer.Write(dd)
		if err != nil {
			log.Fatalln("Couldn't open the csv file", err)
		}
	}
}
func csvdrivers() {

	file, err := os.Create("csvdrivers.csv")
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()
	var w [][]string
	for _, data := range globalD {
		s := data.last_trip_end.String()
		//e := data1.endtime.String()
		sl := fmt.Sprintf("%f", data.last_location)
		//el := fmt.Sprintf("%f", data1.endlocation)
		//gg := strings.Trim(strings.Replace(fmt.Sprint(data.trips), " ", ",", -1), "[]")
		var gg string
		var totaldis float64
		for i, kk := range data.trips {
			totaldis = totaldis + idmap[kk].tripduration
			gg = gg + strconv.Itoa(kk)
			if i+1 < len(data.trips) {
				gg = gg + ","
			}
		}
		d := fmt.Sprintf("%f", totaldis)
		//fmt.Println(i, gg)
		w = append(w, []string{strconv.Itoa(data.Id), s, sl, gg, d})

	}
	for _, dd := range w {
		err = writer.Write(dd)
		if err != nil {
			log.Fatalln("Couldn't open the csv file", err)
		}
	}

}

// func write(){

// 	file, err := os.Create("result.csv")
// 	if err != nil {
// 		log.Fatalln("Couldn't open the csv file", err)
// 	}
// 	defer file.Close()
// 	writer := csv.NewWriter(file)
// 	defer writer.Flush()
// 	for _,data := range tr

// }*/
