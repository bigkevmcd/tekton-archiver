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

type logStreamer func(ns, name string, c kubernetes.Interface) (io.ReadCloser, error)

var defaultLogStreamer logStreamer = streamLogsForPod

func PipelineRunLogs(ctx context.Context, pr *pipelinev1.PipelineRun, clientset kubernetes.Interface) (map[string][]byte, error) {
	prLogData := map[string][]byte{}
	for _, tr := range pr.Status.TaskRuns {
		logs, err := logsForPod(ctx, pr.ObjectMeta.Namespace, tr.Status.PodName, clientset)
		if err != nil {
			return nil, err
		}
		prLogData[tr.Status.PodName] = logs
	}
	return prLogData, nil
}

func logsForPod(ctx context.Context, ns, name string, c kubernetes.Interface) ([]byte, error) {
	podLogs, err := defaultLogStreamer(ns, name, c)
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
