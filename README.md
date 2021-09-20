# Fa (failure analysis) 🔩

## Events

The starting point of failure analysis is a list of events. An Event is represented by the following struct:

```go
type Event struct {
    Time float64
	Type string // one of "ttf", "t+"
}
```

Two types of events are supported:

- Failure (type is `"ttf"`)
- Right-censored (type is `"t+"`)

## Empirical methods

Once you have a list of events, you can apply empirical methods to sample the survival function.

Four types of empirical methods are supported:

- Direct Method (DM)
- Improved Direct Method (IDM)
- Median Rank Method (MRM)
- Product Limit Estimator (PLE)

For each of the supported methods, the package exposes a function that receives as input a list of events and computes a list of points that sample the survival function.

For example, to apply the PLE method, one can use the following function:

```go
func EstimatePLE(events []Event) []Point {
    // ...
}
```

## Distribution fitting

Based on the samples of the survival function, it is possible to find a statistical distribution that best represents the underlying trend.

Two statistical distributions are supported:

- Exponential
- Weibull

For each of the supported distribution, the package exposes a function that receives as input a list of points sampling the survival function and returns an object representing the desired distribution.

For example, to find the Weibull distribution that best fits a list of samples, one can use the following function:

```go
func FitWeibull(points []Point) (w Weibull, r float64, err error) {
    // ...
}
```

`FitWeibull` returns an object of type `Weibull`:

```go
type Weibull struct {
	Beta  float64 // shape parameter
	Theta float64 // scale parameter
}
```

It also returns `r`, the index of fit representing the accuracy of the approximation.

Distributions can be sampled with the methods `Prob` and `Cumul` to compute the value of the density and cumulative probability functions respectively (e.g. to draw a graph).
