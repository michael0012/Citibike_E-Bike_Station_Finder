package main
import (
	"strconv"
	"os"
	"flag"
	"CitiBikeData/backend/citibike"
	"github.com/joho/godotenv"
)



func main(){
	var numClassicBikes, top int
	var lat, long float64
	const MAX_CITIBIKE_STATIONS = 100
	var err1 error = nil
	var err2 error = nil
	err := godotenv.Load()
	if err == nil {
		latString, err1 := strconv.ParseFloat(os.Getenv("LAT"), 64)
		if err1 == nil {
			lat=latString
		}
		longString, err2 := strconv.ParseFloat(os.Getenv("LONG"), 64)
		if err2 == nil {
			long=longString	
		}
	}
	flag.IntVar(&numClassicBikes, "bikes", 0, "an int")
	flag.IntVar(&top, "top", MAX_CITIBIKE_STATIONS, "an int")
	if err != nil || err1 != nil{
		flag.Float64Var(&lat, "lat", 40.647389, "a float 64")
	}
	if err != nil || err2 != nil{
		flag.Float64Var(&long, "long", -74.000917, "a float 64")
	}
	flag.Parse()
	citibike.PrintCitiBikeStationsWithElectric(numClassicBikes, lat, long, top)

}
