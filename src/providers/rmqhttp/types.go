package rmqhttp

import (
	"github.com/medal-labs/k8s-rmq-autoscaler/base/provider"
)

type Config struct {
	Name     provider.Name
	Url      string
	User     string
	Password string
}

type AppConfig struct {
	QueueName string `k8s-annotation:"queue"`
	Vhost     string `k8s-annotation:"vhost"`
}

type QueueInfo struct {
	Consumers       int    `json:"consumers"`
	IdleSince       string `json:"idle_since"`
	Messages        int    `json:"messages"`
	MessagesDetails struct {
		Rate float64 `json:"rate"`
	} `json:"messages_details"`
	MessagesReady        int `json:"messages_ready"`
	MessagesReadyDetails struct {
		Rate float64 `json:"rate"`
	} `json:"messages_ready_details"`
	MessagesUnacknowledged        int `json:"messages_unacknowledged"`
	MessagesUnacknowledgedDetails struct {
		Rate float64 `json:"rate"`
	} `json:"messages_unacknowledged_details"`
	Name  string `json:"name"`
	Node  string `json:"node"`
	State string `json:"state"`
	Vhost string `json:"vhost"`
}
