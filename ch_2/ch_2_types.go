package ch_2

type stepType int

const (
	// enums to represent bandit step types
	sampleAvgStep stepType = iota
	constStep
	varStep

	// constant to represent k & steps
	arms    = 10
	steps   = 20000
	epsilon = 0.1
	runs    = 500
)

// define interface to represent a bandit
type Bandit interface {
	step()
	calcReward() (float64, bool)
	updateQ(float64, int)
	constUpdateQ(float64, int)
	GetTestBed() []float64
}

// bandit model
type bandit struct {
	k       int
	epsilon float64
	n       []int
	q       []float64

	stepType
	testBed []float64

	FieldType string
}

// sample average bandit
type sampAvgBandit struct {
	bandit
}

// exponential-recency weighted bandit
type constStepBandit struct {
	bandit
	stepSize float64
}
