package main

import (
	"fmt"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"github.com/julienschmidt/httprouter"
	"strings"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type JsonResult struct {
	Results []struct {
		Geometry struct {
			Location struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"location"`
		} `json:"geometry"`
	} `json:"results"`
}

type JsonName struct {
	Id         bson.ObjectId `json:"id" bson:"_id"`
	Name       string        `json:"name"`
	Address    string        `json:"address"`
	City       string        `json:"city"`
	State      string        `json:"state"`
	Zip        string        `json:"zip"`
	Coordinate struct {
		Lat float64 `json:"lat"`
		Lng float64 `json:"lng"`
	} `json:"coordinate"`
}

func GetLocations(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
	con, err := mgo.Dial("mongodb://cmpe273:1234@ds045454.mongolab.com:45454/cmpe273")
	if err != nil {
		panic(err)
	}
	defer con.Close()
	con.SetMode(mgo.Monotonic, true)
	c := con.DB("cmpe273").C("addressBook")
	id := p.ByName("name")
	oid := bson.ObjectIdHex(id)
	var result JsonName
	c.FindId(oid).One(&result)
	if err != nil {
		log.Fatal(err)
	}
	oid = bson.ObjectId(result.Id)
	b2, err := json.Marshal(result)
	if err != nil {
	}
	fmt.Fprintf(rw, string(b2))
}

func PostLocations(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
	var myjson3 JsonName
	s3 := json.NewDecoder(req.Body)
	err := s3.Decode(&myjson3)
	StartQuery := "http://maps.google.com/maps/api/geocode/json?address="
	WhereQuery := myjson3.Address + " " + myjson3.City + " " + myjson3.State
	WhereQuery = strings.Replace(WhereQuery, " ", "+", -1)
	EndQuery := "&sensor=false"
	Url1 := StartQuery + WhereQuery + EndQuery
	res, err := http.Get(Url1)
	if err != nil {
		log.Fatal(err)
	}
	robots, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	var JsonResult1 JsonResult
	err = json.Unmarshal(robots, &JsonResult1)
	if err != nil {
		log.Fatal(err)
	}
	myjson3.Id = bson.NewObjectId()
	myjson3.Coordinate.Lat = JsonResult1.Results[0].Geometry.Location.Lat
	myjson3.Coordinate.Lng = JsonResult1.Results[0].Geometry.Location.Lng
	if err != nil {
	}
	con, err := mgo.Dial("mongodb://cmpe273:1234@ds045454.mongolab.com:45454/cmpe273")
	if err != nil {
		panic(err)
	}
	defer con.Close()
	con.SetMode(mgo.Monotonic, true)
	c := con.DB("cmpe273").C("addressBook")
	err = c.Insert(myjson3)
	if err != nil {
		log.Fatal(err)
	}
	result := JsonName{}
	id := myjson3.Id.Hex()
	oid := bson.ObjectIdHex(id)
	c.FindId(oid).One(&result)
	if err != nil {
		log.Fatal(err)
	}
	oid = bson.ObjectId(result.Id)
	b2, err := json.Marshal(result)
	if err != nil {
	}
	fmt.Fprintf(rw, string(b2))
}

func PutLocations(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
	var myjson3 JsonName
	s3 := json.NewDecoder(req.Body)
	err := s3.Decode(&myjson3)
	con, err := mgo.Dial("mongodb://cmpe273:1234@ds045454.mongolab.com:45454/cmpe273")
	if err != nil {
		panic(err)
	}
	defer con.Close()
	con.SetMode(mgo.Monotonic, true)
	c := con.DB("cmpe273").C("addressBook")
	id := p.ByName("name")
	oid := bson.ObjectIdHex(id)
	var result JsonName
	c.FindId(oid).One(&result)
	if myjson3.Name != "" {
		result.Name = myjson3.Name
	}
	if myjson3.Address != "" {
		result.Address = myjson3.Address
	}
	if myjson3.City != "" {
		result.City = myjson3.City
	}
	if myjson3.State != "" {
		result.State = myjson3.State
	}
	if myjson3.Zip != "" {
		result.Zip = myjson3.Zip
	}
	StartQuery := "http://maps.google.com/maps/api/geocode/json?address="
	WhereQuery := result.Address + " " + result.City + " " + result.State
	WhereQuery = strings.Replace(WhereQuery, " ", "+", -1)
	EndQuery := "&sensor=false"
	Url1 := StartQuery + WhereQuery + EndQuery
	res, err := http.Get(Url1)
	if err != nil {
		log.Fatal(err)
	}
	robots, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	var JsonResult1 JsonResult
	err = json.Unmarshal(robots, &JsonResult1)
	if err != nil {
		log.Fatal(err)
	}
	result.Coordinate.Lat = JsonResult1.Results[0].Geometry.Location.Lat
	result.Coordinate.Lng = JsonResult1.Results[0].Geometry.Location.Lng
	c.UpdateId(oid, result)
	if err != nil {
		log.Fatal(err)
	}
	oid = bson.ObjectId(result.Id)
	b2, err := json.Marshal(result)
	if err != nil {
	}
	fmt.Fprintf(rw, string(b2))
}

func DeleteLocations(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
	fmt.Fprintf(rw, "Delete Id: %s\n", p.ByName("name"))
	con, err := mgo.Dial("mongodb://cmpe273:1234@ds045454.mongolab.com:45454/cmpe273")
	if err != nil {
		panic(err)
	}
	defer con.Close()
	con.SetMode(mgo.Monotonic, true)
	c := con.DB("cmpe273").C("addressBook")
	id := p.ByName("name")
	oid := bson.ObjectIdHex(id)
	c.RemoveId(oid)
	fmt.Fprintf(rw, "Deleted: %s\n", p.ByName("name"))
}

func main() {
	mux := httprouter.New()
	mux.POST("/locations", PostLocations)
	mux.GET("/locations/:name", GetLocations)
	mux.PUT("/locations/:name", PutLocations)
	mux.DELETE("/locations/:name", DeleteLocations)
	server := http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: mux,
	}
	server.ListenAndServe()
}