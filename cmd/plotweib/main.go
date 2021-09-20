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
	shape = flag.Float64("shape", 1, "Shape parameter of the Weibull distribution")
	scale = flag.Float64("scale", 1, "Scale parameter of the Weibull distribution")
	from  = flag.Float64("from", 0, "Lower bound for x")
	to    = flag.Float64("to", 10, "Upper bound for x")
	step  = flag.Float64("step", 1, "Step to advange for x")
)

func main() {
	log.SetPrefix("")
	log.SetFlags(0)

	flag.Parse()

	weib := fa.Weibull{
		Beta:  *shape,
		Theta: *scale,
	}

	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.AlignRight)
	fmt.Fprintln(tw, "x\tf(x)\tF(x)\t")
	fmt.Fprintln(tw, "-\t----\t----\t")
	for i := *from; i < *to; i += *step {
		fmt.Fprintf(tw, "%.3f\t%.3f\t%.3f\t\n", i, weib.Prob(i), weib.Cumul(i))
	}
	tw.Flush()
}
