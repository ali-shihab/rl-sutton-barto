package ch_2

import (
	"image/color"
	"math"
	"math/rand"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
)

// calc reward for given action
func (b *bandit) calcReward() (float64, bool) {

	var reward float64
	var rewardIndex int
	var optimal bool

	r := rand.Float64()

	if thrshld := 1 - b.epsilon; r < thrshld {
		_, index := max(b.q)
		mu := b.testBed[index]
		reward = rand.NormFloat64() + mu
		rewardIndex = index
	} else {
		r2 := rand.Intn(10)

		mu := b.testBed[r2]
		reward = rand.NormFloat64() + mu
		rewardIndex = r2
	}

	if b.FieldType == "sampleAverage" {
		b.updateQ(reward, rewardIndex)
	} else {
		b.constUpdateQ(reward, rewardIndex)
	}

	// check if optimal action was chosen
	if _, optIndex := max(b.testBed); rewardIndex == optIndex {
		optimal = true
	} else {
		optimal = false
	}
	return reward, optimal
}

func (b *bandit) updateQ(reward float64, index int) {

	b.n[index] += 1
	n := float64(b.n[index])
	b.q[index] += (reward - b.q[index]) / n
}

func (b *bandit) constUpdateQ(reward float64, index int) {

	b.n[index] += 1
	b.q[index] += 0.1 * (reward - b.q[index])
}

// random walk of test bed
func (b *bandit) step() {
	for i := 0; i < len(b.testBed); i++ {
		mu := rand.NormFloat64() * 0.01
		b.testBed[i] = b.testBed[i] + mu
	}
}

// return test bed reward distributiions
func (b *bandit) GetTestBed() []float64 {
	return b.testBed
}

// find max of slice
func max(slice []float64) (float64, int) {
	maxVal := math.Inf(-1)
	var index int
	for i, v := range slice {
		if v > maxVal {
			maxVal = v
			index = i
		}
	}
	return maxVal, index
}

// plotData creates a line plotter and adds it to the given plot
func plotData(p *plot.Plot, x, y []float64, label string, lineColor color.Color) error {
	pts := make(plotter.XYs, len(x))
	for i := range x {
		pts[i].X = x[i]
		pts[i].Y = y[i]
	}

	line, err := plotter.NewLine(pts)
	if err != nil {
		return err
	}

	line.Color = lineColor

	p.Add(line)
	p.Legend.Add(label, line)

	return nil
}
