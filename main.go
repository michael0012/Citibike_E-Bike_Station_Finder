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
	lat := float64(40.647389)
	long := float64(-74.00091)
	const MAX_CITIBIKE_STATIONS = 100
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
	flag.Float64Var(&lat, "lat", lat, "a float 64")
	flag.Float64Var(&long, "long", long, "a float 64")
	flag.Parse()
	citibike.PrintCitiBikeStationsWithElectric(numClassicBikes, lat, long, top)

}
