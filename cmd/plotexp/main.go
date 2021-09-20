package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/lorciv/fa"
)

var (
	rate = flag.Float64("rate", 1, "Rate of the exponential distribution")
	from = flag.Float64("from", 0, "Lower bound for x")
	to   = flag.Float64("to", 10, "Upper bound for x")
	step = flag.Float64("step", 1, "Step to advange for x")
)

func main() {
	log.SetPrefix("")
	log.SetFlags(0)

	flag.Parse()

	exp := fa.Exponential{
		Rate: *rate,
	}

	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.AlignRight)
	fmt.Fprintln(tw, "x\tf(x)\tF(x)\t")
	fmt.Fprintln(tw, "-\t----\t----\t")
	for i := *from; i < *to; i += *step {
		fmt.Fprintf(tw, "%.3f\t%.3f\t%.3f\t\n", i, exp.Prob(i), exp.Cumul(i))
	}
	tw.Flush()
}
