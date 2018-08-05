package models

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

//Zip represents a zip code record.
//The `json:"..."` field tags allow us to change
//the name of the field when it is encoded into JSON
//see https://drstearns.github.io/tutorials/gojson/
type Zip struct {
	Code  string `json:"code,omitempty"`
	City  string `json:"city,omitempty"`
	State string `json:"state,omitempty"`
}

//ZipSlice is a slice of *Zip,
//that is, pointers to Zip struct instances
type ZipSlice []*Zip

//ZipIndex is a map from strings to ZipSlices
type ZipIndex map[string]ZipSlice

//LoadZips loads the zip code records from a CSV file,
//returning a ZipSlice or an error.
func LoadZips(fileName string) (ZipSlice, error) {
	//open the file
	f, err := os.Open(fileName)
	//if that generated an error, return the error to
	//the caller with a bit of context added to the front
	//of the error message. The %v token will be replaced with
	//the error's message.
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}

	//create a new CSV stream reader over the file
	//this will read only one line at a time, as opposed to
	//reading the entire file into memory, which could be huge
	reader := csv.NewReader(f)

	//read and discard the first row, as it contains the field names
	//the `_` variable is used in Go to ignore a return value, as the
	//compiler requires you to use any variables you declare.
	_, err = reader.Read()
	//if there was an error, return it to the caller with some context
	if err != nil {
		return nil, fmt.Errorf("error reading header row: %v", err)
	}

	//make a ZipSlice with zero-length but 43000 capacity.
	//this ensures we do only one memory allocation up-front,
	//which will be much more efficient than letting the slice
	//re-allocate every time we run out of capacity. If you know
	//how many elements you will need up-front, use the make()
	//function to pre-allocate the underlying array.
	zips := make(ZipSlice, 0, 43000)

	//loop until we get an io.EOF error, which indicates the end
	//of the file, or until we get an CSV parsing error
	for {
		//read the next line
		fields, err := reader.Read()
		//if we got an io.EOF error, we reached the end of the file
		if err == io.EOF {
			//return the ZipSlice and no error
			return zips, nil
		}
		//if we got any other sort of error, return that to the caller
		//TODO: give the caller some hints by adding a line number to
		//the context you add to the error string
		if err != nil {
			return nil, fmt.Errorf("error reading record: %v", err)
		}

		//crate and initialize a new Zip instance
		z := &Zip{
			Code:  fields[0],
			City:  fields[3],
			State: fields[6],
		}

		//append that to the our zip slice
		zips = append(zips, z)
	}
}
