package main

import (
	"github.com/ehsundar/dopamine/cmd"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetReportCaller(true)

	cmd.Execute()
}
