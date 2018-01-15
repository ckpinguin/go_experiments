package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
)

func main() {
	file, err := os.Open("./Post_Adressdaten20171226.csv")
	if err != nil {
		log.Fatalln("Error opening file:", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)

	reader.Comma = ';'
	// reader.ReuseRecord = true
	// reader.FieldsPerRecord = 16
	lineCount := 0

	for {
		lineCount++
		record, err := reader.Read()
		if err == io.EOF {
			break
			// We must ignore all other errors, when we have lines with
			// different numbers of record fields
			// If i holds a T, then t will be the underlying value and ok will be true.
		} else if err, ok := err.(*csv.ParseError); ok {
			if err.Err == csv.ErrFieldCount {
				log.Println(record)
				continue
			} else if err.Err == csv.ErrBareQuote {
				log.Println("Barequote! Line", lineCount, "record:", record)
				continue
			}
		} else if err != nil {
			log.Fatalln("Error reading record:", err, record)
			return
		}

		// log.Println("Record", lineCount, "is", record, "and has", len(record), "fields")
		for i := 0; i < len(record); i++ {
			log.Println(" ", record[i])
		}
		log.Println()
	}
	log.Println("Finished reading, total no. of lines:", lineCount)
}
