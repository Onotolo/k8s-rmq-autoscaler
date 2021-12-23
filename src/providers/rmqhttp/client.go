package rmqhttp

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type rmqHTTPClient struct {
	*http.Client
	config Config
}

func (client rmqHTTPClient) getQueueInfo(queue string, vhost string) (*QueueInfo, error) {
	reqUrl := fmt.Sprintf("%s/api/queues/%s/%s", client.config.Url, vhost, queue)
	req, err := http.NewRequest("GET", reqUrl, nil)
	req.SetBasicAuth(client.config.User, client.config.Password)
	if err != nil {
		return nil, err
	}
	response, err := client.Do(req)
	if err != nil || response.StatusCode != http.StatusOK {
		return nil, err
	}
	var info QueueInfo
	if err := json.NewDecoder(response.Body).Decode(&info); err != nil {
		return nil, err
	}
	return &info, nil
}
