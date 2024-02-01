package ch_2

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"math/rand"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/vg"
)

const (
	a       = 0.1
	avgwndw = 100000
	s       = 200000
)

func Execute_2_11() {
	paramCurve := EpsilonParameterStudy(a, s, avgwndw, runs)

	// plot parameter curve
	// Create a new plot, panicking on error.
	p := plot.New()

	// Set titles and labels
	p.Title.Text = "Parameter Curve - Epsilon"
	p.X.Label.Text = "Epsilon"
	p.Y.Label.Text = "Average Reward"

	// Action step values
	epsilons := make([]float64, 6)
	for i := -7; i < -1; i++ {
		epsilons[i+7] = math.Pow(2, float64(i))
	}

	// Create line plotter objects and add them to the plot
	err := plotData(p, epsilons, paramCurve, "Exponential-Recency Weighted", color.RGBA{R: 255, A: 255})
	if err != nil {
		log.Panic(err)
	}

	// Save the plot to a PNG file.
	if err := p.Save(8*vg.Inch, 8*vg.Inch, "paramcurves.png"); err != nil {
		log.Panic(err)
	}
}

// parameter curve estimation for epsilon-greedy bandits on epsilon in {2^-7, ..., 2^-2}
func EpsilonParameterStudy(a float64, s int, avgwndw int, runs int) []float64 {
	// set up plots of rewards & optmality curves
	// Create a new plot, panicking on error.
	p := plot.New()

	// Set titles and labels
	p.Title.Text = "Reward Curves"
	p.X.Label.Text = "Time Step"
	p.Y.Label.Text = "Average Reward"

	// plot reward and optimality curves
	// Create a new plot, panicking on error.
	p2 := plot.New()

	// Set titles and labels
	p2.Title.Text = "Optimal Actions"
	p2.X.Label.Text = "Time Step"
	p2.Y.Label.Text = "Fraction of Actions Optimally Taken"

	// Action step values
	actionSteps := make([]float64, s)
	for i := 0; i < s; i++ {
		actionSteps[i] = float64(i + 1)
	}

	paramCurve := make([]float64, 6)

	for i := -7; i < -1; i++ {
		e := math.Pow(2, float64(i))

		rewardCurve, optimalityCurve := Simulate(e, a, s, runs)
		avg := AverageRewardOverWindow(rewardCurve, avgwndw)

		paramCurve[i+7] = avg

		// plot params
		epsilonString := fmt.Sprintf("%f", e)
		label := "Epsilon" + epsilonString
		scale := math.Pow((avg / rewardCurve[s-1]), 2)

		// Create line plotter objects and add them to the plot
		err := plotData(p, actionSteps, rewardCurve, label, color.RGBA{R: uint8(scale * 255), G: 255, A: 255})
		if err != nil {
			log.Panic(err)
		}

		// Create line plotter objects and add them to the plot
		err = plotData(p2, actionSteps, optimalityCurve, label, color.RGBA{R: uint8(scale * 255), G: 255, A: 255})
		if err != nil {
			log.Panic(err)
		}
	}

	// Save the plot to a PNG file.
	if err := p.Save(8*vg.Inch, 8*vg.Inch, "ex_2_11_rewards.png"); err != nil {
		log.Panic(err)
	}

	// Save the plot to a PNG file.
	if err := p2.Save(8*vg.Inch, 8*vg.Inch, "ex_2_11_optimality.png"); err != nil {
		log.Panic(err)
	}

	return paramCurve
}

// simulate s steps of 10-armed bandit on nonstationary testbed for runs iterations for constant step size a
func Simulate(e float64, a float64, s int, runs int) ([]float64, []float64) {
	fmt.Println(e)
	const_thousand_run_average_reward := make([]float64, s)
	const_thousand_run_average_optimality := make([]float64, s)

	// run for x iterations for averaging
	for j := 0; j < runs; j++ {
		// seed
		rand.New(rand.NewSource(int64(j)))

		// initialise test bed
		testBed := make([]float64, arms)

		// init exponential-recency weighted bandit
		var constBandit Bandit = &constStepBandit{

			bandit: bandit{
				k:       arms,
				epsilon: e,
				n:       make([]int, 10),
				q:       make([]float64, 10),

				stepType:  constStep,
				testBed:   testBed,
				FieldType: "constant",
			},
			stepSize: a,
		}

		constBanditOptimality := make([]float64, s)
		constBanditAvgReward := make([]float64, s)

		var avgOptimality float64 = 0
		var avgReward float64 = 0

		for i := 0; i < s; i++ {

			// step bandit
			constBandit.step()

			// calculate reward
			reward, optimal := constBandit.calcReward()

			// update metrics of bandit
			if i > 0 {
				avgReward = constBanditAvgReward[i-1]
				constBanditAvgReward[i] = avgReward + (reward-avgReward)/float64(i)

				avgOptimality = constBanditOptimality[i-1]
				if optimal {
					constBanditOptimality[i] = avgOptimality + (1-avgOptimality)/float64(i)
				} else {
					constBanditOptimality[i] = avgOptimality + (0-avgOptimality)/float64(i)
				}
			} else {
				constBanditAvgReward[i] = reward
				if optimal {
					constBanditOptimality[i] = 1
				} else {
					constBanditOptimality[i] = 0
				}
			}

			// update thousand run average for bandit
			const_thousand_run_average_reward[i] = const_thousand_run_average_reward[i] + (constBanditAvgReward[i]-const_thousand_run_average_reward[i])/float64(j+1)

			const_thousand_run_average_optimality[i] = const_thousand_run_average_optimality[i] + (constBanditOptimality[i]-const_thousand_run_average_optimality[i])/float64(j+1)
		}
	}
	// print results
	fmt.Println("Average reward of constant step: ", const_thousand_run_average_reward[s-1],
		" and optimality: ", const_thousand_run_average_optimality[s-1])

	return const_thousand_run_average_reward, const_thousand_run_average_optimality
}

// calculate the average of the last avgwndw steps of the reward curve
func AverageRewardOverWindow(rc []float64, avggwndw int) float64 {
	sum := 0.0

	for _, val := range rc {
		sum += val
	}
	return sum / float64(avggwndw)
}
