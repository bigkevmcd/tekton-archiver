package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	pipelineclientset "github.com/tektoncd/pipeline/pkg/client/clientset/versioned"
	"go.uber.org/zap"
	"k8s.io/client-go/rest"
	"knative.dev/pkg/signals"

	"github.com/bigkevmcd/tekton-archiver/pkg/watcher"
)

func makeArchiveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "archive",
		Short: "archive Tekton PipelineRuns and TaskRuns",
		RunE: func(cmd *cobra.Command, args []string) error {
			clusterConfig, err := rest.InClusterConfig()
			if err != nil {
				return fmt.Errorf("failed to create a cluster config: %s", err)
			}

			tektonClient, err := pipelineclientset.NewForConfig(clusterConfig)
			if err != nil {
				return fmt.Errorf("failed to create the tekton client: %v", err)
			}

			logger, _ := zap.NewProduction()
			defer func() {
				err := logger.Sync() // flushes buffer, if any
				if err != nil {
					log.Println(err)
				}
			}()
			sugar := logger.Sugar()
			stopCh := signals.SetupSignalHandler()
			watcher.WatchPipelineRuns(stopCh, tektonClient, "default", sugar)
			<-stopCh
			return nil
		},
	}
	return cmd
}
