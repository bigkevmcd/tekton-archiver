package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func makeArchiveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "archive",
		Short: "archive Tekton PipelineRuns and TaskRuns",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	return cmd
}

// stopper returns a channel that remains open until an interrupt is received.
func stopper() chan struct{} {
	stop := make(chan struct{})
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		logrus.Warn("Interrupt received, attempting clean shutdown...")
		close(stop)
		<-c
		logrus.Error("Second interrupt received, force exiting...")
		os.Exit(1)
	}()
	return stop
}
