package main

/* install following dependencies*/

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	//"github.com/paulmach/go.geojson"

	geo "github.com/paulmach/go.geo"
)

const ApiVersion = "v1"
const Profile = "driving"
const RouteService = "route"
const TripService = "trip"

type OSRMClient struct {
	baseUrl string
}
type RouteOptions struct {
	Locations    []geo.Pointer
	Overview     string
	Steps        bool
	Alternatives bool
}

func NewOSRMClient(baseUrl string) *OSRMClient {
	return &OSRMClient{baseUrl: baseUrl}
}

type Point [2]float64
type Pointer interface {
	Point() *Point
}
type RouteResponse struct {
	Routes []struct {
		Geometry string  `json:"geometry"`
		Distance float64 `json:"distance"`
		Duration float64 `json:"duration"`
	} `json:"routes"`
}

func (opts RouteOptions) buildRequest(baseUrl, service string) (*http.Request, error) {
	if service != RouteService && service != TripService {
		return nil, fmt.Errorf("Service should be one of %s, %s", RouteService, TripService)
	}
	if len(opts.Locations) < 2 {
		return nil, fmt.Errorf("you need to have atleast 2 locations")
	}
	var locs []string
	for _, loc := range opts.Locations {
		locs = append(locs, fmt.Sprintf("%f,%f", loc.Point().Lng(), loc.Point().Lat()))
	}
	locations := strings.Join(locs, ";")
	reqUrl := fmt.Sprintf("%s/%s/%s/%s/%s", baseUrl, service, ApiVersion, Profile, locations)
	req, err := http.NewRequest(http.MethodGet, reqUrl, nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Add("overview", opts.Overview)
	q.Add("steps", strconv.FormatBool(opts.Steps))
	if service == RouteService {
		q.Add("alternatives", strconv.FormatBool(opts.Alternatives))
	}
	if service == TripService {
		q.Add("source", "first")
		q.Add("destination", "last")
		q.Add("roundtrip", "false")
	}
	req.URL.RawQuery = q.Encode()
	//fmt.Println(req)
	return req, nil
}

func (client *OSRMClient) GetRoute(opts RouteOptions) (*RouteResponse, error) {
	req, err := opts.buildRequest(client.baseUrl, RouteService)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var routeResp RouteResponse
	json.Unmarshal(body, &routeResp)
	if err != nil {
		return nil, err
	}
	return &routeResp, nil
}
func osrmURL() string {
	osrmURL := os.Getenv("OSRM_URL")
	if osrmURL == "" {
		//osrmURL = "http://ec2-54-255-208-196.ap-southeast-1.compute.amazonaws.com:5000"
		osrmURL = "http://localhost:5000"
	}
	return osrmURL
}
func getRoute(full bool, points []geo.Pointer) (*RouteResponse, error) {
	osrmapi := NewOSRMClient(osrmURL())
	overview := "simplified"
	if full {
		overview = "full"
	}
	options := RouteOptions{
		Alternatives: false,
		Overview:     overview,
		Steps:        false,
	}
	for _, pt := range points {
		options.Locations = append(options.Locations, pt.Point())
	}
	return osrmapi.GetRoute(options)
}
