package scalable

import (
	"github.com/medal-labs/k8s-rmq-autoscaler/common"
	"time"
)

// App struct used to store information about a scalable instance
type App struct {
	Ref           interface{}
	Annotations   *map[string]string
	Key           string
	Name          string
	ReadyReplicas int
	Replicas      int
	UpdatedDate   time.Time
}

type AppId = string

func (app *App) ParseAnnotations(v interface{}, prefixes ...string) error {
	return common.ParseK8sAnnotations(*app.Annotations, v, prefixes...)
}
