package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/lorciv/fa"
)

var (
	methodFlag  = flag.String("m", "idm", "Empirical method to be used to sample the survival function")
	verboseFlag = flag.Bool("v", false, "Show the samples of the survival function")
)

func main() {
	log.SetPrefix("fa: ")
	log.SetFlags(0)

	flag.Parse()

	var events []fa.Event

	scan := bufio.NewScanner(os.Stdin)
	i := 1
	for scan.Scan() {
		e, err := parseEvent(scan.Text())
		if err != nil {
			log.Fatalf("line %d: %v", i, err)
		}
		events = append(events, e)
		i++
	}
	if err := scan.Err(); err != nil {
		log.Fatal(err)
	}

	var points []fa.Point

	switch *methodFlag {
	case "dm":
		points = fa.EstimateDM(events)
	case "idm":
		points = fa.EstimateIDM(events)
	case "mrm":
		points = fa.EstimateMRM(events)
	case "ple":
		points = fa.EstimatePLE(events)
	default:
		log.Fatalf("invalid method %q", *methodFlag)
	}

	weib, r, err := fa.FitWeibull(points)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Weibull:\tbeta/theta/r = %f/%f/%f\n", weib.Beta, weib.Theta, r)

	exp, r, err := fa.FitExponential(points)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Exponential:\trate/r = %f/%f\n", exp.Rate, r)

	if *verboseFlag {
		tw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.AlignRight)
		fmt.Fprintln(tw, "i\tt\tR(t)\tWeibull\tExp\t")
		fmt.Fprintln(tw, "-\t-\t----\t-------\t---\t")
		for i, p := range points {
			fmt.Fprintf(tw, "%d\t%.1f\t%.3f\t%.3f\t%.3f\t\n", i, p.X, p.Y, weib.Survival(p.X), exp.Survival(p.X))
		}
		tw.Flush()
	}
}

func parseEvent(s string) (fa.Event, error) {
	split := strings.Split(s, ",")
	if len(split) != 2 {
		return fa.Event{}, errors.New("could not parse event: wrong number of fields")
	}

	time, err := strconv.ParseFloat(split[0], 64)
	if err != nil {
		return fa.Event{}, fmt.Errorf("could not parse event: %v", err)
	}

	if split[1] != "ttf" && split[1] != "t+" {
		return fa.Event{}, errors.New("could not parse event: invalid type")
	}

	return fa.Event{
		Time: time,
		Type: split[1],
	}, nil
}
