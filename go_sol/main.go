package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"
)

type CityData struct {
	City       string `json:"city"`
	StateShort string `json:"state_short"`
	StateFull  string `json:"state_full"`
}

type StateData struct {
	StateShort string `json:"state_short"`
	StateFull  string `json:"state_full"`
}

type Data interface {
	GetKey(record []string) string
	NewInstance(record []string) interface{}
	GetClassName() string
}

type CityDataProcessor struct{}

func (c *CityDataProcessor) GetKey(record []string) string {
	return strings.Join(record[0:3], ",")
}

func (c *CityDataProcessor) NewInstance(record []string) interface{} {
	return CityData{
		City:       record[0],
		StateShort: record[1],
		StateFull:  record[2],
	}
}

func (c *CityDataProcessor) GetClassName() string {
	return "Cities"
}

type StateDataProcessor struct{}

func (s *StateDataProcessor) GetKey(record []string) string {
	return strings.Join(record[1:3], ",")
}

func (s *StateDataProcessor) NewInstance(record []string) interface{} {
	return StateData{
		StateShort: record[1],
		StateFull:  record[2],
	}
}

func (s *StateDataProcessor) GetClassName() string {
	return "States"
}

func Process(reader *csv.Reader, dataProcessor Data, outputFileName string) {
	count := -1
	dataMap := make(map[string]interface{})

	for {
		record, err := reader.Read()

		if err != nil {
			break
		}

		if len(record) == 0 {
			break
		}

		if count == -1 {
			count++
			continue
		}

		txt := dataProcessor.GetKey(record)

		if _, exists := dataMap[txt]; !exists {
			dataMap[txt] = dataProcessor.NewInstance(record)
			count++
		}
	}

	keys := make([]string, 0, len(dataMap))

	for key := range dataMap {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	arr := make([]interface{}, 0, len(dataMap))

	for _, k := range keys {
		arr = append(arr, dataMap[k])
	}

	// Write results to an output file
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()

	// Marshal the data to JSON
	jsonData, err := json.Marshal(arr)
	if err != nil {
		log.Fatal(err)
	}

	// Write the JSON data to the file
	_, err = outputFile.Write(jsonData)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Total %s: %d\n", strings.ToLower(dataProcessor.GetClassName()), len(dataMap))
}

func main() {
	defer timeTrack(time.Now(), "Process")

	inputFileName := "us_cities_states_counties.csv"
	file, err := os.Open(inputFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = '|'

	Process(reader, &CityDataProcessor{}, "cities.json")

	file.Seek(0, 0)
	reader = csv.NewReader(file)
	reader.Comma = '|'

	Process(reader, &StateDataProcessor{}, "states.json")
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %f sec", name, elapsed.Seconds())
}
