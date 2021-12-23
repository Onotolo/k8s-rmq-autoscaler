package loop

import (
	"context"
	"errors"
	"github.com/medal-labs/k8s-rmq-autoscaler/base/executor"
	"github.com/medal-labs/k8s-rmq-autoscaler/base/scalable"
	"github.com/medal-labs/k8s-rmq-autoscaler/base/strategy"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog"
	"sync"
	"time"
)

type AutoscalerLoop struct {
	add      chan *v1.Deployment
	delete   chan *v1.Deployment
	apps     map[string]scalable.App
	client   *kubernetes.Clientset
	recorder record.EventRecorder
}

type Config struct {
	ExecutorCfg     executor.Config
	InCluster       bool
	Namespaces      string
	LoopTickSeconds int
}

func Launch(ctx context.Context, cfg Config) error {

	l := AutoscalerLoop{
		apps:   make(map[string]scalable.App),
		delete: make(chan *v1.Deployment),
		add:    make(chan *v1.Deployment),
	}
	var err error

	l.client, err = discover(ctx, &l, cfg.InCluster, cfg.Namespaces)
	if err != nil {
		return err
	}

	l.recorder, err = eventRecorder(l.client)
	if err != nil {
		return err
	}

	go func() {

		loopTick := time.NewTicker(time.Duration(cfg.LoopTickSeconds) * time.Second)
		defer loopTick.Stop()
		for {
			select {
			case deployment := <-l.add:
				if err := l.addDeployment(deployment); err != nil {
					klog.Error(err)
					continue
				}
			case deployment := <-l.delete:
				key, _ := cache.MetaNamespaceKeyFunc(deployment)
				klog.Infof("%s: deleting app", key)
				delete(l.apps, key)

			case <-loopTick.C:

				startTime := time.Now()
				apps := make([]scalable.App, len(l.apps))

				appIndex := 0
				for _, app := range l.apps {
					apps[appIndex] = app
					appIndex += 1
				}

				results, errs := executor.Launch(cfg.ExecutorCfg, apps)

				wg := sync.WaitGroup{}
				wg.Add(2)

				go func() {
					defer wg.Done()
					for err := range errs {
						l.handleError(err)
					}
				}()
				go func() {
					defer wg.Done()
					for result := range results {
						l.applyScalingResult(ctx, result, l.recorder)
					}
				}()
				wg.Wait()
				if klog.V(2) {
					klog.Infof("Finished scaling round in %s", time.Now().Sub(startTime))
				}

			case <-ctx.Done():
				klog.Info("Finishing due to context cancellation")
				return
			}
		}
	}()
	return nil
}

func (l AutoscalerLoop) addDeployment(deployment *v1.Deployment) error {
	key, _ := cache.MetaNamespaceKeyFunc(deployment)

	app, err := createApp(deployment, key)

	if err != nil {
		klog.Error(err)
		return err
	}
	if _, ok := l.apps[key]; ok {
		// Already exist
		klog.Infof("%s: updating app", key)
	} else {
		klog.Infof("%s: new app", key)
	}
	l.apps[key] = *app

	return nil
}

func createApp(deployment *v1.Deployment, key string) (*scalable.App, error) {

	if _, ok := deployment.ObjectMeta.Annotations[AnnotationPrefix+Enable]; !ok {
		return nil, errors.New(key + " not concerned by autoscaling, skipping")
	}

	return &scalable.App{
		Ref:           deployment,
		Key:           key,
		Name:          deployment.Name,
		Replicas:      int(*deployment.Spec.Replicas),
		ReadyReplicas: int(deployment.Status.ReadyReplicas),
		UpdatedDate:   time.Now(),
		Annotations:   &deployment.Annotations,
	}, nil
}

func (l *AutoscalerLoop) handleError(err executor.Error) {
	klog.Errorf("Got error during strategies execution: %s", err)
	baseErr, ok := err.(executor.BaseError)
	if !ok {
		return
	}
	depl, ok := baseErr.App.Ref.(*v1.Deployment)
	l.recorder.Eventf(depl, corev1.EventTypeWarning, "ASWarning", "error during scaling", err)
}

func (l *AutoscalerLoop) applyScalingResult(ctx context.Context, result strategy.Result, recorder record.EventRecorder) {
	app := result.App

	ref, ok := app.Ref.(*v1.Deployment)
	if !ok {
		klog.Errorf("%s app ref contains value of unexpected type: expected *v1.Deployment, got %t", app.Key, app.Ref)
		return
	}
	if result.Skip {
		klog.Infof("%s scaling will be skipped", app.Key)
		return
	}
	if int(*ref.Spec.Replicas) == result.RequiredReplicas {
		klog.Infof("%s scaling will be skipped: requested replicas number hasn't changed", app.Key)
		recorder.Eventf(ref, corev1.EventTypeNormal, "ASSkip", "Scaling will be skipped: requested replicas number hasn't changed")
		return
	}
	klog.Infof("%s Will be updated from %d replicas to %d", app.Key, app.Replicas, result.RequiredReplicas)
	newReplicas := int32(result.RequiredReplicas)

	increment := result.RequiredReplicas - app.Replicas

	if increment > 0 {
		recorder.Eventf(ref, corev1.EventTypeNormal, "ASScaleUp", "Scaling up to %d replicas", newReplicas)
	} else if increment < 0 {
		recorder.Eventf(ref, corev1.EventTypeNormal, "ASScaleDown", "Scaling down to %d replicas", newReplicas)
	}

	ref.Spec.Replicas = &newReplicas
	newRef, err := l.client.AppsV1().Deployments(ref.Namespace).Update(ctx, ref, metav1.UpdateOptions{})

	if err != nil {
		klog.Errorf("Error during deployment (%s) update, retry later (%s)", app.Key, err)
	} else {
		app.Ref = newRef
	}
}

// eventRecorder returns an EventRecorder type that can be
// used to post Events to different object's lifecycles.
func eventRecorder(
	kubeClient *kubernetes.Clientset) (record.EventRecorder, error) {
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(klog.Infof)
	eventBroadcaster.StartRecordingToSink(
		&typedcorev1.EventSinkImpl{
			Interface: kubeClient.CoreV1().Events("")})
	recorder := eventBroadcaster.NewRecorder(
		scheme.Scheme,
		corev1.EventSource{Component: "autoscaler.rabbitmq"})
	return recorder, nil
}
