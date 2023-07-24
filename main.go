package main

import (
	"github.com/plantoncloud/kube-namespace-deleter/cmd/kubenamespacedeleter"
	clipanic "github.com/plantoncloud/kube-namespace-deleter/internal/cli/panic"
)

func main() {
	kubenamespacedeleter.Execute()
	finished := new(bool)
	defer clipanic.Handle(finished)
	kubenamespacedeleter.Execute()
	*finished = true
}
