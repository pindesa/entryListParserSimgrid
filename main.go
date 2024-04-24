package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"unicode/utf16"
)

type JsonData struct {
	Entries        []Entry `json:"entries"`
	ForceEntryList int     `json:"forceEntryList"`
}

type Entry struct {
	Driver                       []Driver `json:"drivers"`
	RaceNumber                   int      `json:"raceNumber"`
	ForcedCarModel               int      `json:"forcedCarModel"`
	OverrideDriverInfo           int      `json:"overrideDriverInfo"`
	DefaultGridPos               int      `json:"defaultGridPosition"`
	BallastKg                    int      `json:"ballastKg"`
	Restrictor                   int      `json:"restrictor"`
	CustomCar                    string   `json:"customCar"`
	OverrideCarModelForCustomCar int      `json:"overrideCarModelForCustomCar"`
	IsServerAdmin                int      `json:"isServerAdmin"`
}

type Driver struct {
	FirstName      string `json:"firstName"`
	LastName       string `json:"lastName"`
	ShortName      string `json:"shortName"`
	Nationality    int    `json:"nationality"`
	DriverCategory int    `json:"driverCategory"`
	PlayerID       string `json:"playerID"`
}

func writeFileUtf16(name string, outputJson []byte) error {
	var bytes [2]byte
	const BOM = '\ufffe' //LE. for BE '\ufeff'

	file, err := os.Create(name)
	if err != nil {
		fmt.Printf("Can't open file. %v", err)
		return err
	}
	defer file.Close()

	bytes[0] = BOM >> 8
	bytes[1] = BOM & 255

	file.Write(bytes[0:])
	runes := utf16.Encode([]rune(string(outputJson)))
	for _, r := range runes {
		bytes[1] = byte(r >> 8)
		bytes[0] = byte(r & 255)
		file.Write(bytes[0:])
	}
	return nil
}

func splitDrivers(inputFile, outputFilePath string) error {
	// Read the input JSON file
	data, err := os.ReadFile(inputFile)
	if err != nil {
		return err
	}

	// Unmarshal the JSON into a map with string keys and json.RawMessage values
	// var jsonData map[string]json.RawMessage
	// if err := json.Unmarshal(data, &jsonData); err != nil {
	// 	return err
	// }

	var jsonStruct JsonData
	if err := json.Unmarshal(data, &jsonStruct); err != nil {
		return err
	}
	// Extract the "entries" key
	// entriesRaw, ok := jsonData["entries"]
	// if !ok {
	// 	return fmt.Errorf("invalid JSON format: 'entries' key not found")
	// }

	// // Unmarshal the "entries" value into a slice of Entry
	// var entries []Entry
	// if err := json.Unmarshal(entriesRaw, &entries); err != nil {
	// 	return err
	// }

	// Create an array to store individual driver entries
	var driversJson JsonData
	// Iterate over each entry and create a separate entry for each driver

	for _, entry := range jsonStruct.Entries {
		//fmt.Println(entry)
		for _, driver := range entry.Driver {
			//fmt.Println(driver)
			if driver.FirstName == "" {
				continue
			}
			driverEntry := Entry{
				Driver:                       []Driver{driver},
				RaceNumber:                   -1,
				ForcedCarModel:               entry.ForcedCarModel,
				OverrideDriverInfo:           entry.OverrideDriverInfo,
				DefaultGridPos:               entry.DefaultGridPos,
				BallastKg:                    entry.BallastKg,
				Restrictor:                   entry.Restrictor,
				CustomCar:                    entry.CustomCar,
				OverrideCarModelForCustomCar: 0,
				IsServerAdmin:                entry.IsServerAdmin,
			}
			//fmt.Println(driverEntry)
			driversJson.Entries = append(driversJson.Entries, driverEntry) //[i] = driverEntry

		}
	}
	driversJson.ForceEntryList = 1
	// Marshal the driver entries back to JSON with an indentation of 2 spaces
	driverEntriesJSON, err := json.MarshalIndent(driversJson, "", "  ")
	if err != nil {
		return err
	}
	// var bytes[2]byte
	// cons BOM = '\ufffe'
	// bytes[0] = BOM >> 8
	// bytes[1] = BOM & 255
	// runeJson := utf16.DecodeRune(driverEntriesJSON, 'U+FFFE')

	// Write the JSON to the output file

	//this part creates a file
	// than it writes bytes into to be able to encode it as utf16 le BOM

	if err := writeFileUtf16(outputFilePath, driverEntriesJSON); err != nil {
		return err
	}

	fmt.Printf("Saved driver entries to %s\n", outputFilePath)

	return nil
}

func main() {
	// Example usage:
	inputFile := flag.String("if", "entrylist_team.json", "input json file with full path")
	outputFile := flag.String("of", "entrylist_solo.json", "output json file with full path")
	flag.Parse()
	fmt.Println(*inputFile)
	fmt.Println(*outputFile)
	var inputFilePath string = *inputFile
	outputFilePath := *outputFile
	// inputFilePath := "/home/pindesa/golang/goProjects/jsonParseParmic/entrylist_team.json"
	// outputFilePath := "/home/pindesa/golang/goProjects/jsonParseParmic/test_team.json"

	if err := splitDrivers(inputFilePath, outputFilePath); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
