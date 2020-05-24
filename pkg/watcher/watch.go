package watcher

import (
	"fmt"

	pipelinev1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	pipelineclientset "github.com/tektoncd/pipeline/pkg/client/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labelsv1 "k8s.io/apimachinery/pkg/labels"

	"github.com/bigkevmcd/tekton-archiver/pkg/logs"
)

type logger interface {
	Infof(template string, args ...interface{})
	Infow(msg string, keysAndValues ...interface{})
	Errorf(template string, args ...interface{})
	Errorw(msg string, keysAndValues ...interface{})
}

func WatchPipelineRuns(stop <-chan struct{}, e logs.Extractor, tektonClient pipelineclientset.Interface, ns string, l logger) {
	l.Infow("starting to watch for PipelineRuns", "ns", ns)
	api := tektonClient.TektonV1beta1().PipelineRuns(ns)
	listOptions := metav1.ListOptions{
		LabelSelector: labelsv1.Set(map[string]string{"archive.tekton.dev": "true"}).AsSelector().String(),
	}
	watcher, err := api.Watch(listOptions)
	if err != nil {
		l.Errorf("failed to watch PipelineRuns: %s", err)
		return
	}
	ch := watcher.ResultChan()
	for {
		select {
		case <-stop:
			return
		case v := <-ch:
			pr := v.Object.(*pipelinev1.PipelineRun)
			err := handlePipelineRun(tektonClient, pr, l)
			if err != nil {
				l.Infow(fmt.Sprintf("error handling PipelineRun: %s", err), "name", pr.ObjectMeta.Name)
			}
		}
	}
}

func handlePipelineRun(tektonClient pipelineclientset.Interface, pr *pipelinev1.PipelineRun, l logger) error {
	newState := runState(pr)
	l.Infof("Received a PipelineRun %#v %s", pr.Status, newState)
	return nil
}
