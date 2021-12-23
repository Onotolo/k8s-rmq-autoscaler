package provider

import (
	"github.com/medal-labs/k8s-rmq-autoscaler/base/parameter"
	"github.com/medal-labs/k8s-rmq-autoscaler/base/scalable"
	"k8s.io/klog"
	"sync"
)

type AppContext struct {
	*baseContext
	Parameters []parameter.Name
}

type ResultAppContext struct {
	*baseContext
}

type baseContext struct {
	App          scalable.App
	ProviderName Name
	mx           sync.Mutex
	isDone       bool
	results      chan Result
	done         chan struct{}
}

func (ctx AppContext) PutResult(parameters ProvidedParameters) {
	if ctx.IsCanceled() {
		return
	}
	select {
	case <-ctx.done:
	case ctx.results <- Result{
		Parameters: parameters,
	}:
	}
}

func (ctx AppContext) Error(err error) {
	if ctx.IsCanceled() {
		return
	}
	select {
	case <-ctx.done:
	case ctx.results <- Result{
		Error: err,
	}:
	}
}

func (ctx AppContext) IsCanceled() bool {
	return ctx.isClosed()
}

func (ctx AppContext) Finish() {
	if klog.V(3) {
		klog.Infof("Finishing app context for '%s' app and '%s' provider", ctx.App.Name, ctx.ProviderName)
	}
	ctx.close()
}

func (ctx ResultAppContext) GetNextResult() (Result, bool) {
	select {
	case <-ctx.done:
		return Result{}, false
	case res := <-ctx.results:
		return res, true
	}
}

func (ctx ResultAppContext) Cancel() {
	if klog.V(3) {
		klog.Infof("Cancelling app context for '%s' app and '%s' provider", ctx.App.Name, ctx.ProviderName)
	}
	ctx.close()
}

func (ctx ResultAppContext) IsCanceled() bool {
	return ctx.isClosed()
}

func (s *baseContext) close() {
	s.mx.Lock()
	defer s.mx.Unlock()
	if !s.isDone {
		close(s.done)
		s.isDone = true
	}
}

func (s *baseContext) isClosed() bool {
	s.mx.Lock()
	defer s.mx.Unlock()
	return s.isDone
}

func newConnectedContexts(app scalable.App, cfg Config, paramsNames []parameter.Name) (AppContext, ResultAppContext) {
	appResults := make(chan Result)
	appCancel := make(chan struct{})

	s := baseContext{
		App:          app,
		ProviderName: cfg.Name,
		results:      appResults,
		done:         appCancel,
	}
	ctx := AppContext{
		baseContext: &s,
		Parameters:  paramsNames,
	}
	resultCtx := ResultAppContext{
		baseContext: &s,
	}
	return ctx, resultCtx
}
