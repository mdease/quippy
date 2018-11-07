package db

import (
	"encoding/csv"
	"bufio"
	"math/rand"
	"os"
	"time"
)


const path string = "./res/prod.csv"

// Columns: question, num blanks, weight
var db []*Row
var totalWeight int64

// Load the database of questions into memory
func Load() {
	// Open the file
	file, err := os.Open(path)

	if err != nil {
		panic("Could not find database")
	}

	// Create a csv reader
	reader := csv.NewReader(file)

	reader.LazyQuotes = true

	// Read the whole table
	db_temp, err := reader.ReadAll()

	if err != nil {
		panic("Could not load database")
	}

	// Close the file
	err = file.Close()

	if err != nil {
		panic("Could not close file")
	}

	// Init db
	db = make([]*Row, len(db_temp))

	// Create a Row object for each row
	for i, row := range db_temp {
		newRow := NewRow(i, row)
		db[i] = newRow
		totalWeight += newRow.Weight
	}

	// Seed the RNG
	rand.Seed(time.Now().UTC().UnixNano())
}

// Save the database of questions
func Save() {
	// Open the file
	file, err := os.OpenFile(path, os.O_WRONLY, 0777)

	if err != nil {
		panic("Could not find database")
	}

	// Construct a string to write to the file
	data := ""

	// Read each row as csv strings
	for _, row := range db {
		data += row.AsCsvRow()
	}

	// Create a file writer
	writer := bufio.NewWriter(file)

	// Write line-by-line :/
	_, err = writer.WriteString(data)

	if err != nil {
		panic("Could not write to database")
	}

	err = writer.Flush()

	if err != nil {
		panic("Could not save database")
	}

	// Close the file
	err = file.Close()

	if err != nil {
		panic("Could not close file")
	}
}

// Randomly select a question based on weights
func Sample() *Row {
	r := int64(rand.Intn(int(totalWeight)))

	for _, row := range db {
		weight := row.Weight

		// Fields might have negative weight
		if weight > 0 {
			r -= weight
		}

		if r < 0 {
			return row
		}
	}

	panic("Fatal error")
}

// Increase a question's weight (max 5)
func Upvote(index int) {
	row := db[index]

	if row.Weight < 5 {
		row.Weight++
		totalWeight++
	}
}

// Decrease a question's weight
func Downvote(index int) {
	row := db[index]

	row.Weight--
	totalWeight--
}
