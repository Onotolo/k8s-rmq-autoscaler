package executor

import (
	"fmt"
	"github.com/medal-labs/k8s-rmq-autoscaler/base/provider"
	"github.com/medal-labs/k8s-rmq-autoscaler/base/scalable"
	"github.com/medal-labs/k8s-rmq-autoscaler/base/strategy"
)

type Error error

type BaseError struct {
	App      scalable.App
	Strategy strategy.Config
	Err      error
}

type ProviderError struct {
	BaseError
	ProviderName provider.Name
}

func (e BaseError) Error() string {
	return fmt.Sprintf("Failed to get scaling result for %s: %s", e.App.Name, e.Err)
}

func (e ProviderError) Error() string {
	return fmt.Sprintf("Failed to get scaling result for %s: %s", e.App.Name, e.Err)
}
