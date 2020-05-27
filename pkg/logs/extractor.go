package logs

import (
	"bytes"
	"context"
	"fmt"
	"io"

	pipelinev1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

// K8sExtractor is an implementation of the Extractor interface, getting the log
// output using a Kubernetes client.
type K8sExtractor struct {
	clientset kubernetes.Interface
	streamer  logStreamer
}

type logStreamer func(ns, name string, c kubernetes.Interface) (io.ReadCloser, error)

// New creates and returns a new K8sExtractor.
func New(c kubernetes.Interface) *K8sExtractor {
	return &K8sExtractor{clientset: c, streamer: streamLogsForPod}
}

// PipelineRun fetches the logs for each Pod associated with a PipelineRun.
func (k *K8sExtractor) PipelineRun(ctx context.Context, pr *pipelinev1.PipelineRun) (map[string][]byte, error) {
	prLogData := map[string][]byte{}
	for _, tr := range pr.Status.TaskRuns {
		logs, err := logsForPod(ctx, pr.ObjectMeta.Namespace, tr.Status.PodName, k.clientset, k.streamer)
		if err != nil {
			return nil, err
		}
		prLogData[tr.Status.PodName] = logs
	}
	return prLogData, nil
}

func logsForPod(ctx context.Context, ns, name string, c kubernetes.Interface, streamer logStreamer) ([]byte, error) {
	podLogs, err := streamer(ns, name, c)
	if err != nil {
		return nil, fmt.Errorf("error in opening stream: %w", err)
	}
	defer podLogs.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)
	if err != nil {
		return nil, fmt.Errorf("error in copy logs from pod to buffer: %w", err)
	}
	return buf.Bytes(), nil
}

// This is not currently testable using the fake client, as GetLogs() returns an
// empty Request.
func streamLogsForPod(ns, name string, c kubernetes.Interface) (io.ReadCloser, error) {
	podLogOpts := corev1.PodLogOptions{}
	req := c.CoreV1().Pods(ns).GetLogs(name, &podLogOpts)
	return req.Stream()
}
