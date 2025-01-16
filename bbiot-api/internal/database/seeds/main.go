package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/GDGVIT/bbiot-backend/internal/constants"
	"github.com/GDGVIT/bbiot-backend/internal/database"
)

func main() {
	dbpool, err := database.InitialiseDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to initialise DB: %v\n", err)
		os.Exit(1)
	}

	db := database.New(dbpool)
	// _, err := os.Executable()
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "Error retreiving path: %v\n", err)
	// }

	centering_sheets := readCsvFile("./internal/database/seeds/centering_sheet.csv")
	for _, c := range centering_sheets {
		_, err := db.AddProductTagList(context.Background(), database.AddProductTagListParams{
			ProductType: constants.CENTERING_SHEET,
			Epchex:      c[0],
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error while entering centering_sheet: %v\n", err)
			continue
		}
	}

	props := readCsvFile("./internal/database/seeds/prop.csv")
	for _, p := range props {
		_, err := db.AddProductTagList(context.Background(), database.AddProductTagListParams{
			ProductType: constants.PROP,
			Epchex:      p[0],
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error while entering prop: %v\n", err)
			continue
		}
	}
	spans := readCsvFile("./internal/database/seeds/span.csv")
	for _, s := range spans {
		_, err := db.AddProductTagList(context.Background(), database.AddProductTagListParams{
			ProductType: constants.SPAN,
			Epchex:      s[0],
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error while entering span: %v\n", err)
			continue
		}
	}

	

}

func readCsvFile(filePath string) [][]string {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
	}

	return records
}
