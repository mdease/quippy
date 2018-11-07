package db

import (
	"strconv"
	"strings"
)


// Object representing a row in the question database
type Row struct {
	Index int
	Question string
	NumBlanks int64
	Weight int64
	Answers []string
	RespondentIDs []string
}

// Create a new Row object based on the csv string data
func NewRow(i int, row []string) *Row {
	question := row[0]

	// Parse the numBlanks field
	numBlanksStr := row[1]
	numBlanks, err := strconv.ParseInt(numBlanksStr, 10, 64)

	if err != nil {
		panic("Error parsing database")
	}

	// Parse the weight field
	weightStr := row[2]
	weight, err := strconv.ParseInt(weightStr, 10, 64)

	if err != nil {
		panic("Error parsing database")
	}

	return &Row { Index: i, Question: question, NumBlanks: numBlanks, Weight: weight }
}

// Return the given Row object as its csv representation
func (r *Row) AsCsvRow() string {
	// Escape quotes within the question
	question := r.Question
	question = strings.Replace(question, "\"", "\"\"", -1)
	question = "\"" + question + "\""

	// Cast the numBlanks field
	numBlanks := strconv.FormatInt(r.NumBlanks, 10)

	// Cast the weight field
	weight := strconv.FormatInt(r.Weight, 10)

	return question + "," + numBlanks + "," + weight + "\n"
}
