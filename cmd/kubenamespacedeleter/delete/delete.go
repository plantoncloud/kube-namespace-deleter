package _delete

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/plantoncloud/kube-namespace-deleter/pkg/kubernetes/namespace/deleter"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var Delete = &cobra.Command{
	Use:   "delete",
	Short: "Deletes Kubernetes Namespace stuck in Terminating State",
	Run:   deleteHandler,
}

func deleteHandler(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		log.Fatalf("namespace name must be provided")
	}

	namespaceName := args[0]

	executor := deleter.RealDeleter{}
	if err := deleter.Delete(namespaceName, executor); err != nil {
		log.Fatalf("%#v", errors.Wrap(err, "failed to delete namespace"))
	}

	fmt.Printf("namespace %s deleted successfully\n", namespaceName)
}
