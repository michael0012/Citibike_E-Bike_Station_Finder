package main
import (
	"strconv"
	"os"
	"flag"
	"CitiBikeData/backend/citibike"
	"github.com/joho/godotenv"
)



func main(){
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	var numClassicBikes, top int
	var lat, long float64
	const MAX_CITIBIKE_STATIONS = 100
	latString, err1 := strconv.ParseFloat(os.Getenv("LAT"), 64)
	if err1 == nil {
		lat=latString
	}
	longString, err2 := strconv.ParseFloat(os.Getenv("LONG"), 64)
	if err2 == nil {
		long=longString	
	}
	flag.IntVar(&numClassicBikes, "bikes", 0, "an int")
	flag.IntVar(&top, "top", MAX_CITIBIKE_STATIONS, "an int")
	if err1 != nil{
		flag.Float64Var(&lat, "lat", 40.646389, "a float 64")
	}
	if err2 != nil{
		flag.Float64Var(&long, "long", -74.001918, "a float 64")
	}
	flag.Parse()
	citibike.PrintCitiBikeStationsWithElectric(numClassicBikes, lat, long, top)

}
