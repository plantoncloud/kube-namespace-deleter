package kubenamespacedeleter

import (
	"fmt"
	_delete "github.com/plantoncloud/kube-namespace-deleter/cmd/kubenamespacedeleter/delete"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var debug bool

var rootCmd = &cobra.Command{
	Use:   "kube-namespace-deleter",
	Short: "CLI to Delete Kubernetes Namespace stuck in Terminating State",
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "set log level to debug")
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.DisableSuggestions = true
	cobra.OnInitialize(func() {
		if debug {
			log.SetLevel(log.DebugLevel)
			log.Debug("running in debug mode")
		}
	})
	rootCmd.AddCommand(_delete.Delete)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
