package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/spf13/cobra"
	pipelineclientset "github.com/tektoncd/pipeline/pkg/client/clientset/versioned"
	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"knative.dev/pkg/signals"

	"github.com/bigkevmcd/tekton-archiver/pkg/archiver/s3"
	"github.com/bigkevmcd/tekton-archiver/pkg/logs"
	"github.com/bigkevmcd/tekton-archiver/pkg/watcher"
)

func makeArchiveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "archive",
		Short: "archive Tekton PipelineRuns",
		RunE: func(cmd *cobra.Command, args []string) error {
			clusterConfig, err := rest.InClusterConfig()
			if err != nil {
				return fmt.Errorf("failed to create a cluster config: %s", err)
			}

			tektonClient, err := pipelineclientset.NewForConfig(clusterConfig)
			if err != nil {
				return fmt.Errorf("failed to create the tekton client: %v", err)
			}

			coreClient, err := kubernetes.NewForConfig(clusterConfig)
			if err != nil {
				return fmt.Errorf("failed to create the core client: %v", err)
			}
			extractor := logs.New(coreClient)
			logger, _ := zap.NewProduction()
			defer func() {
				err := logger.Sync() // flushes buffer, if any
				if err != nil {
					log.Println(err)
				}
			}()
			sugar := logger.Sugar()
			stopCh := signals.SetupSignalHandler()
			watcher.WatchPipelineRuns(stopCh, extractor, createUploader(os.Getenv("BUCKET_NAME")), tektonClient, "default", sugar)
			<-stopCh
			return nil
		},
	}
	return cmd
}

func createUploader(bucket string) *s3.S3Archiver {
	sess := session.Must(session.NewSession())
	return s3.New(s3.NewS3Uploader(bucket, s3manager.NewUploader(sess)))

}
