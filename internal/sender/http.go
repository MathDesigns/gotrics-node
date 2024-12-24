package sender

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func SendMetrics(serverAddr string, metrics interface{}) error {
	data, err := json.Marshal(metrics)
	if err != nil {
		return err
	}

	_, err = http.Post(serverAddr+"/metrics", "application/json", bytes.NewBuffer(data))
	return err
}
