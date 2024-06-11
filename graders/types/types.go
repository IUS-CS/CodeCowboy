package types

import "time"

type GraderReturn struct {
	Passed   float64
	Failed   float64
	Coverage float64
	TimeLate time.Duration
}
