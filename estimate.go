package fa

import (
	"fmt"
	"sort"
)

type Event struct {
	Time float64
	Type string // one of "ttf", "t+"
}

type Point struct {
	X, Y float64
}

func (p Point) String() string {
	return fmt.Sprintf("(%f, %f)", p.X, p.Y)
}

// EstimateDM computes an estimate of the survival function according to the Direct Method.
func EstimateDM(events []Event) []Point {
	return estimateCompleteData(events, func(i, n int) float64 {
		return 1 - (float64(i) / float64(n))
	})
}

// EstimateDM computes an estimate of the survival function according to the Improved Direct Method.
func EstimateIDM(events []Event) []Point {
	return estimateCompleteData(events, func(i, n int) float64 {
		return 1 - (float64(i) / float64(n+1))
	})
}

// EstimateDM computes an estimate of the survival function according to the Median Rank Method.
func EstimateMRM(events []Event) []Point {
	return estimateCompleteData(events, func(i, n int) float64 {
		return 1 - ((float64(i) - 0.3) / (float64(n) + 0.4))
	})
}

func estimateCompleteData(events []Event, f func(n, i int) float64) []Point {
	filteredEvents := []Event{
		{Time: 0},
	}
	for _, e := range events {
		if e.Type == "ttf" {
			filteredEvents = append(filteredEvents, e)
		}
	}
	sort.Slice(filteredEvents, func(i, j int) bool {
		return filteredEvents[i].Time < filteredEvents[j].Time
	})

	var points []Point
	n := len(filteredEvents) - 1
	for i, e := range filteredEvents {
		if i == 0 {
			points = append(points, Point{X: 0, Y: 1})
			continue
		}
		points = append(points, Point{
			X: e.Time,
			Y: f(i, n),
		})
	}
	return points
}

// EstimateDM computes an estimate of the survival function according to the Product Limit Estimator.
func EstimatePLE(events []Event) []Point {
	filteredEvents := []Event{
		{Time: 0},
	}
	filteredEvents = append(filteredEvents, events...)
	sort.Slice(filteredEvents, func(i, j int) bool {
		return filteredEvents[i].Time < filteredEvents[j].Time
	})

	var values []float64
	n := len(events)
	for i, e := range filteredEvents {
		if i == 0 {
			values = append(values, 1)
			continue
		}
		corr := 1.0
		if e.Type == "ttf" {
			corr = float64(n+1-i) / float64(n+2-i)
		}
		values = append(values, values[i-1]*corr)
	}

	var points []Point
	for i, e := range filteredEvents {
		if i == 0 || e.Type == "ttf" {
			points = append(points, Point{
				X: e.Time,
				Y: values[i],
			})
		}
	}
	return points
}
