package main
//Main code using raster as an input
import (
	"fmt"
	"flag"
	"github.com/USACE/go-consequences/compute"
	"github.com/USACE/go-consequences/hazardproviders"
	"github.com/USACE/go-consequences/resultswriters"
	"github.com/USACE/go-consequences/structureprovider"
)

func main() {
	rasterPath := flag.String("raster", "Mosaic_Ida_ft_4326.tif", "Raster elevation generated by a model (adcirc or ras)")
	siPath := flag.String("si", "SI.shp", "Structure Inventory shapefile")
	result := flag.String("result", "result.shp", "Shapefile with the result of Go Consequences for structure Inventory and damage")
	flag.Parse()
	fmt.Println("Start")
	//define hazard
	myEventHazard, e := hazardproviders.Init(*rasterPath)
	if e != nil {
		panic(e)
	}
	defer myEventHazard.Close()
	shpStructureProvider, b := structureprovider.InitSHP(*siPath)
	if b != nil {
		panic(b)
	}
	//check for close
	//define result writer
	myEventResults, c := resultswriters.InitShpResultsWriter(*result, "event_results")
	if c != nil {
		panic(c)
	}
	defer myEventResults.Close()
	//run the compute
	compute.StreamAbstract(myEventHazard, shpStructureProvider, myEventResults)
	fmt.Println("End")
}

