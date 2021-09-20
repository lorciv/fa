package fa

import (
	"errors"
	"math"
	"sort"
)

type Exponential struct {
	Rate float64
}

// logProb computes the natural logarithm of the value of the probability density function at x.
func (e Exponential) logProb(x float64) float64 {
	if x < 0 {
		return math.Inf(-1)
	}
	return math.Log(e.Rate) - e.Rate*x
}

// Prob computes the value of the probability density function at x.
func (e Exponential) Prob(x float64) float64 {
	return math.Exp(e.logProb(x))
}

// Cumul computes the value of the cumulative density function at x.
func (e Exponential) Cumul(x float64) float64 {
	if x < 0 {
		return 0
	}
	return -math.Expm1(-e.Rate * x)
}

// Survival returns the survival function (complementary CDF) at x.
func (e Exponential) Survival(x float64) float64 {
	if x < 0 {
		return 1
	}
	return math.Exp(-e.Rate * x)
}

// FitExponential fits a Weibull distribution to the survival function sampled by the given points.
func FitExponential(points []Point) (e Exponential, r float64, err error) {
	sort.Slice(points, func(i, j int) bool {
		return points[i].X < points[j].X
	})

	origPoints := points
	for _, p := range origPoints {
		if p.X == 0 {
			points = points[1:]
		} else if p.X < 0 {
			return Exponential{}, 0, errors.New("cannot fit Exponential on negative points")
		}
	}

	trans := make([]Point, len(points))
	avg := Point{}
	prodSum, xSquareSum := 0.0, 0.0
	for i := 0; i < len(points); i++ {
		x, y := points[i].X, math.Log(1/points[i].Y)
		trans[i] = Point{
			X: x,
			Y: y,
		}
		avg.X += x
		avg.Y += y
		prodSum += x * y
		xSquareSum += math.Pow(x, 2)
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

	// exponential dist parameter
	e.Rate = prodSum / xSquareSum

	// index of fit
	sqrt1 := math.Sqrt(xSqDiffSum / float64(len(points)))
	sqrt2 := math.Sqrt(ySqDiffSum / float64(len(points)))
	r = (diffProdSum / float64(len(points))) / (sqrt1 * sqrt2)

	return e, r, nil
}
