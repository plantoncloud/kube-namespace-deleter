package kubenamespacedeleter

import (
	"fmt"
	"github.com/plantoncloud/kube-namespace-deleter/pkg/kubernetes/namespace/deleter"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var debug bool

var rootCmd = &cobra.Command{
	Use:   "kube-namespace-deleter",
	Short: "Deletes Kubernetes Namespace stuck in Terminating State",
	Run:   deleteHandler,
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "set log level to debug")
	rootCmd.PersistentFlags().StringP("namespace", "n", "", "Name of the Namespace to Delete")
	rootCmd.DisableSuggestions = true
	cobra.OnInitialize(func() {
		if debug {
			log.SetLevel(log.DebugLevel)
			log.Debug("running in debug mode")
		}
	})
}

func deleteHandler(cmd *cobra.Command, args []string) {
	namespaceName, err := cmd.Flags().GetString("namespace")
	if err != nil {
		log.Fatalf("%#v", err)
	}
	if namespaceName == "" {
		err = cmd.Help()
		if err != nil {
			log.Fatalf("%#v", err)
		}
		return
	}
	executor := deleter.RealDeleter{}
	if err = deleter.Delete(namespaceName, executor); err != nil {
		log.Fatalf("%#v", err)
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
