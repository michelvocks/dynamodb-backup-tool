package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func main() {
	// Input parameters
	dataFolder := flag.String("data", "data", "path to the data folder")
	tableName := flag.String("table", "", "table name to export")
	mode := flag.String("mode", "export", "export or import that should be executed")
	flag.Parse()

	switch *mode {
	case "export":
		// Export data
		exportData(*tableName, *dataFolder)
	case "import":
		// Import data
		importData(*tableName, *dataFolder)
	default:
		panic(errors.New("mode is unknown"))
	}
}

func importData(tableName, dataFolder string) {
	sess, err := session.NewSessionWithOptions(session.Options{
		// enable shared config support.
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		panic(err)
	}

	svc := dynamodb.New(sess)

	// Get all files
	files, err := ioutil.ReadDir(dataFolder)
	if err != nil {
		panic(err)
	}

	// Iterate files
	for _, f := range files {
		// load file
		raw, err := ioutil.ReadFile(path.Join(dataFolder, f.Name()))
		if err != nil {
			panic(err)
		}

		// Transform raw data into stuct
		var items []map[string]*dynamodb.AttributeValue
		err = json.Unmarshal(raw, &items)
		if err != nil {
			panic(err)
		}

		// Iterate items
		ri := map[string][]*dynamodb.WriteRequest{}
		for _, v := range items {
			pq := &dynamodb.PutRequest{
				Item: v,
			}

			ri[tableName] = append(ri[tableName], &dynamodb.WriteRequest{
				PutRequest: pq,
			})
		}

		// Insert items into dynamodb
		input := &dynamodb.BatchWriteItemInput{
			RequestItems: ri,
		}
		result, err := svc.BatchWriteItem(input)
		if err != nil {
			panic(err)
		}
		fmt.Println(result)
	}
}

func exportData(tableName, dataFolder string) {
	sess, err := session.NewSessionWithOptions(session.Options{
		// enable shared config support.
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		panic(err)
	}

	svc := dynamodb.New(sess)
	input := &dynamodb.ScanInput{
		TableName: &tableName,
	}

	// Create output folder
	err = os.MkdirAll(dataFolder, os.ModePerm)
	if err != nil {
		panic(err)
	}

	pagenum := 0
	err = svc.ScanPages(input, func(page *dynamodb.ScanOutput, lastPage bool) bool {
		pagenum++
		fmt.Printf("Current page: %d\n", pagenum)

		// Open file
		f, err := os.Create(path.Join(dataFolder, fmt.Sprintf("%d.json", pagenum)))
		if err != nil {
			panic(err)
		}
		defer f.Close()

		// Marshal
		out, err := json.Marshal(page.Items)
		if err != nil {
			panic(err)
		}

		_, err = f.Write(out)
		if err != nil {
			panic(err)
		}
		f.Sync()
		return true
	})
	if err != nil {
		panic(err)
	}
}
