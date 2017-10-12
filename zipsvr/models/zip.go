package models

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

// Keep it capitalized in order for JSON files to encode them
type Zip struct {
	Code  string
	City  string
	State string
}

// declaring a slice of a pointer
type ZipSlice []*Zip

type ZipIndex map[string]ZipSlice

// parameter = file name
// return = zipslice and error (u can return multiple things)
func LoadZips(fileName string) (ZipSlice, error) {

	// whats with the two variables??
	f, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}

	// the reader will read one line at the time, instead of opening up the entire file
	reader := csv.NewReader(f)
	_, err = reader.Read() // the underscore is an ignorer variable.
	if err != nil {
		return nil, fmt.Errorf("error reading header row: %v", err)
	}

	// the make() will preallocating space to be x length. very efficient and much faster
	zips := make(ZipSlice, 0, 43000)
	for {
		// io.eof is the error that signifies end of iteration
		fields, err := reader.Read() // this will split the row into strings
		if err == io.EOF {
			// the indicator that it is the end of the file
			return zips, nil
		}
		if err != nil {
			return nil, fmt.Errorf("error reading record: %v", err)
		}

		// the & takes and creates an address
		z := &Zip{
			Code:  fields[0],
			City:  fields[3],
			State: fields[6],
		}

		zips = append(zips, z)
	}
}
