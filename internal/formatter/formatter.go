package formatter

import (
	"encoding/json"
	"fmt"
)

// Pretty cleanly formats a given data structure for easily reading.
func Pretty(v interface{}) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(b))
}
