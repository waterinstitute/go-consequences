// +build utf8

package main
//Main code using csv file as an input

import (
	"encoding/csv"
	"fmt"
	"os"
	"flag"
	"strconv"

	"github.com/USACE/go-consequences/consequences"
	"github.com/USACE/go-consequences/hazards"
	"github.com/USACE/go-consequences/resultswriters"
	"github.com/USACE/go-consequences/structures"
	//"github.com/USACE/go-consequences/compute"
	//"github.com/USACE/go-consequences/hazardproviders"
)

//define structure inventory
func main() {
	//define result writer
	var conversion float64
	deepPath := flag.String("depth", "peak_waterlevel.csv", "Depth elevation file generated by a model (adcirc or ras)")
	siPath := flag.String("si", "SI.csv", "Structure Inventory CSV")
	result := flag.String("result", "result.shp", "Shapefile with the result of Go Consequences for structure Inventory and damage")
	unit := flag.String("unit", "meters", "Unit of depth values")
	flag.Parse()
	myEventResults, c := resultswriters.InitShpResultsWriter(*result, "event_results")
	if c != nil {
		panic(c)
	}
	defer myEventResults.Close()

	if *unit == "meters" {
		conversion = 3.28084
		_ = conversion		
	} else {
		conversion = 1.0
		_ = conversion
	}

	otp := structures.JsonOccupancyTypeProvider{}
	otp.InitDefault()
	m := otp.OccupancyTypeMap()
	defaultOcctype := m["RES1-1SNB"]

	filePathStructure, err := os.Open(*siPath)
    if err != nil {
        fmt.Println("Couldn't open structure inventory file")
        panic(err)
    }

	defer filePathStructure.Close()

	structureInventoryLines, err2 := csv.NewReader(filePathStructure).ReadAll()

	if err2 != nil {
		fmt.Println("Couldn't read structure inventory file")
		panic(err2)
	}

	filePathDepth, err3 := os.Open(*deepPath)

	if err3 != nil {
		fmt.Println("Can't find depth file")
		panic(err3)
	}

	defer filePathDepth.Close()

	DepthLines, err4 := csv.NewReader(filePathDepth).ReadAll()

	if err4 != nil {
		fmt.Println("Couldn't read depth file")
		panic(err4)
	}

	for i, line := range structureInventoryLines {
		if i > 0 {
			structure := structures.StructureStochastic{}
			structure.Name = line[0]
			structure.DamCat = line[4]
			structure.CBFips = line[1]
			structure.X, _ = strconv.ParseFloat(line[2], 64)
			structure.Y, _ = strconv.ParseFloat(line[3], 64)
			foundHt, _ := strconv.ParseFloat(line[8], 64)
			structure.FoundHt = consequences.ParameterValue{Value: foundHt}
			structVal, _ := strconv.ParseFloat(line[6], 64)
			structure.StructVal = consequences.ParameterValue{Value: structVal}
			contVal, _ := strconv.ParseFloat(line[7], 64)
			structure.ContVal = consequences.ParameterValue{Value: contVal}
			structure.FoundType = line[9]

			OccTypeName := line[5]

			var occtype = defaultOcctype
			if ot, ok := m[OccTypeName]; ok {
				occtype = ot
			} else {
				occtype = defaultOcctype
				msg := "Using default " + OccTypeName + " not found"
				fmt.Println(msg)
				//return s, errors.New(msg)
			}

			structure.OccType = occtype

			hazard := hazards.DepthEvent{}
			depth, err5 := strconv.ParseFloat(DepthLines[i][3], 64)
			if err5 != nil {
				fmt.Println("Depth Error")
				depth = 0
			}
			hazard.SetDepth(depth*conversion)

			result, err6 := structure.Compute(hazard)

			if err6 != nil {
				fmt.Println("Compute Error")
			}
			if err6 == nil {
				myEventResults.Write(result)
			}

		}

	}

}