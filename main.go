package main

import (
	"flag"

	"github.com/pdiorio/santashuffle/internal/notify"
	"github.com/pdiorio/santashuffle/internal/selection"
)

var (
	participantPtr *string
	gmailPtr       *string
	seedVal        *int64
	dryrun         *bool
)

func init() {
	participantPtr = flag.String("participants", "conf/participants.yaml.example", "path to a yaml file with participant metadata")
	gmailPtr = flag.String("config", "conf/gmail.yaml.example", "path to a yaml file with email configuration")
	seedVal = flag.Int64("seed", 0, "initial seed value; 0 indicates initialize randomness with the current execution time")
	dryrun = flag.Bool("dryrun", false, "execute in dryrun mode (examine secret results & do not send messages)")
}

func main() {
	flag.Parse()

	matches := selection.RunSelection(*participantPtr, *seedVal, *dryrun)
	notify.NotifyPariticpants(matches, *gmailPtr, *dryrun)
}
