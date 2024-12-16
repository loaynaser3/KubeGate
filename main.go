package main

import (
	"github.com/loaynaser3/KubeGate/cmd"
	"github.com/loaynaser3/KubeGate/pkg/logging"
)

func main() {
	logging.Initialize()
	cmd.Execute()
}
