package fa

import (
	"errors"
	"math"
	"sort"
)

type Weibull struct {
	Beta  float64 // shape parameter
	Theta float64 // scale parameter
}

// logProb computes the natural logarithm of the value of the probability
// density function at x. -Inf is returned if x is less than zero.
//
// Special cases occur when x == 0, and the result depends on the shape
// parameter as follows:
//  If 0 < K < 1, logProb returns +Inf.
//  If K == 1, logProb returns 0.
//  If K > 1, logProb returns -Inf.
func (w Weibull) logProb(x float64) float64 {
	if x < 0 {
		return math.Inf(-1)
	}
	if x == 0 && w.Beta == 1 {
		return 0
	}
	return math.Log(w.Beta) - math.Log(w.Theta) + (w.Beta-1)*(math.Log(x)-math.Log(w.Theta)) - math.Pow(x/w.Theta, w.Beta)
}

// Prob computes the value of the probability density function at x.
func (w Weibull) Prob(x float64) float64 {
	if x < 0 {
		return 0
	}
	return math.Exp(w.logProb(x))
}

// Cumul computes the value of the cumulative density function at x.
func (w Weibull) Cumul(x float64) float64 {
	if x < 0 {
		return 0
	}
	return -math.Expm1(-math.Pow(x/w.Theta, w.Beta))
}

// logSurvival returns the log of the survival function (complementary CDF) at x.
func (w Weibull) logSurvival(x float64) float64 {
	if x < 0 {
		return 0
	}
	return -math.Pow(x/w.Theta, w.Beta)
}

// Survival returns the survival function (complementary CDF) at x.
func (w Weibull) Survival(x float64) float64 {
	return math.Exp(w.logSurvival(x))
}

// FitWeibull fits a Weibull distribution to the survival function sampled by the given points.
func FitWeibull(points []Point) (w Weibull, r float64, err error) {
	sort.Slice(points, func(i, j int) bool {
		return points[i].X < points[j].X
	})

	origPoints := points
	for _, p := range origPoints {
		if p.X == 0 {
			points = points[1:]
		} else if p.X < 0 {
			return Weibull{}, 0, errors.New("cannot fit Weibull on negative points")
		}
	}

	trans := make([]Point, len(points))
	avg := Point{}
	for i := 0; i < len(points); i++ {
		trans[i] = Point{
			X: math.Log(points[i].X),
			Y: math.Log(math.Log(1 / points[i].Y)),
		}
		avg.X += trans[i].X
		avg.Y += trans[i].Y
	}
	avg.X /= float64(len(points))
	avg.Y /= float64(len(points))

	diffProdSum := 0.0
	xSqDiffSum, ySqDiffSum := 0.0, 0.0
	for i := 0; i < len(points); i++ {
		xDiff, yDiff := trans[i].X-avg.X, trans[i].Y-avg.Y
		diffProdSum += xDiff * yDiff
		xSqDiffSum += math.Pow(xDiff, 2)
		ySqDiffSum += math.Pow(yDiff, 2)
	}

	// Weibull parameters
	w.Beta = diffProdSum / xSqDiffSum
	w.Theta = math.Exp(avg.X - (avg.Y / w.Beta))

	// index of fit
	sqrt1 := math.Sqrt(xSqDiffSum / float64(len(points)))
	sqrt2 := math.Sqrt(ySqDiffSum / float64(len(points)))
	r = (diffProdSum / float64(len(points))) / (sqrt1 * sqrt2)

	return w, r, nil
}
