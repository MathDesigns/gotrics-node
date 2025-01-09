package sender

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func SendMetrics(serverAddr string, metrics interface{}) error {
	data, err := json.Marshal(metrics)
	if err != nil {
		return err
	}

	fmt.Printf("Sending the following data: %s\n", string(data))

	_, err = http.Post(serverAddr+"/metrics", "application/json", bytes.NewBuffer(data))
	return err
}
