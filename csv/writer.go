package csv

import (
	"bytes"
	"encoding/csv"
)

func WriteToCSV(header []string, content [][]string) (csvFile []byte) {
	buffer := &bytes.Buffer{}
	w := csv.NewWriter(buffer)
	err := w.Write(header)
	if err != nil {
		return nil
	}
	for _, val := range content {
		if err = w.Write(val); err != nil {
			return nil
		}
	}
	w.Flush()
	return buffer.Bytes()
}
