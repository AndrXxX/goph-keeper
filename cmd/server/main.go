package main

import (
	"log"

	"github.com/AndrXxX/goph-keeper/pkg/buildformatter"
	"github.com/AndrXxX/goph-keeper/pkg/logger"
)

var buildVersion string
var buildDate string
var buildCommit string

func main() {
	if iErr := logger.Initialize("Info"); iErr != nil {
		log.Fatal(iErr)
	}
	buildFormatter := buildformatter.BuildFormatter{
		Labels: []string{"Build version", "Build date", "Build commit"},
		Values: []string{buildVersion, buildDate, buildCommit},
	}
	for _, bInfo := range buildFormatter.Format() {
		logger.Log.Info(bInfo)
	}
}
