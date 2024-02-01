package ch_2

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"

	//"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/vg"
)

func Execute_2_5() {

	samp_avg_thousand_run_average_reward := make([]float64, steps)
	const_thousand_run_average_reward := make([]float64, steps)

	samp_avg_thousand_run_average_optimality := make([]float64, steps)
	const_thousand_run_average_optimality := make([]float64, steps)

	// run for 5000 iterations for averaging
	for j := 0; j < runs; j++ {
		// seed
		rand.New(rand.NewSource(int64(j)))

		// initialise test bed
		testBed := make([]float64, arms)

		// init sample average bandit
		var avgBandit Bandit = &sampAvgBandit{

			bandit: bandit{
				k:       arms,
				epsilon: epsilon,
				n:       make([]int, 10),
				q:       make([]float64, 10),

				stepType:  sampleAvgStep,
				testBed:   testBed,
				FieldType: "sampleAverage",
			},
		}

		// init exponential-recency weighted bandit
		var constBandit Bandit = &constStepBandit{

			bandit: bandit{
				k:       arms,
				epsilon: epsilon,
				n:       make([]int, 10),
				q:       make([]float64, 10),

				stepType:  constStep,
				testBed:   testBed,
				FieldType: "constant",
			},
			stepSize: 0.1,
		}

		avgBanditOptimality := make([]float64, steps)
		constBanditOptimality := make([]float64, steps)

		avgBanditAvgReward := make([]float64, steps)
		constBanditAvgReward := make([]float64, steps)

		var avgOptimality float64 = 0
		var avgReward float64 = 0

		for i := 0; i < steps; i++ {

			avgBandit.step()

			// step bandit 1
			reward, optimal := avgBandit.calcReward()

			// update metrics of bandit 1
			if i > 0 {
				avgReward = avgBanditAvgReward[i-1]
				avgBanditAvgReward[i] = avgReward + (reward-avgReward)/float64(i)

				avgOptimality = avgBanditOptimality[i-1]
				if optimal {
					avgBanditOptimality[i] = avgOptimality + (1-avgOptimality)/float64(i)
				} else {
					avgBanditOptimality[i] = avgOptimality + (0-avgOptimality)/float64(i)
				}
			} else {
				avgBanditAvgReward[i] = reward
				if optimal {
					avgBanditOptimality[i] = 1
				} else {
					avgBanditOptimality[i] = 0
				}
			}

			// update thousand run average for bandit 1
			samp_avg_thousand_run_average_reward[i] = samp_avg_thousand_run_average_reward[i] + (avgBanditAvgReward[i]-samp_avg_thousand_run_average_reward[i])/float64(j+1)

			samp_avg_thousand_run_average_optimality[i] = samp_avg_thousand_run_average_optimality[i] + (avgBanditOptimality[i]-samp_avg_thousand_run_average_optimality[i])/float64(j+1)

			// step bandit 2
			reward, optimal = constBandit.calcReward()

			// update metrics of bandit 2
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

			// update thousand run average for bandit 2
			const_thousand_run_average_reward[i] = const_thousand_run_average_reward[i] + (constBanditAvgReward[i]-const_thousand_run_average_reward[i])/float64(j+1)

			const_thousand_run_average_optimality[i] = const_thousand_run_average_optimality[i] + (constBanditOptimality[i]-const_thousand_run_average_optimality[i])/float64(j+1)
		}
	}
	// print results
	fmt.Println("Average reward of constant step: ", const_thousand_run_average_reward[steps-1],
		" and optimality: ", const_thousand_run_average_optimality[steps-1],
		"\nAverage reward of sample average step: ", samp_avg_thousand_run_average_reward[steps-1],
		" and optimality: ", samp_avg_thousand_run_average_optimality[steps-1])

	// plot reward and optimality curves
	// Create a new plot, panicking on error.
	p := plot.New()

	// Set titles and labels
	p.Title.Text = "Reward Curves"
	p.X.Label.Text = "Time Step"
	p.Y.Label.Text = "Average Reward"

	// Action step values
	actionSteps := make([]float64, steps)
	for i := 0; i < steps; i++ {
		actionSteps[i] = float64(i + 1)
	}

	// Create line plotter objects and add them to the plot
	err := plotData(p, actionSteps, const_thousand_run_average_reward, "Exponential-Recency Weighted", color.RGBA{R: 255, A: 255})
	if err != nil {
		log.Panic(err)
	}

	err = plotData(p, actionSteps, samp_avg_thousand_run_average_reward, "Sample Averaged", color.RGBA{G: 255, A: 255})
	if err != nil {
		log.Panic(err)
	}

	// plot reward and optimality curves
	// Create a new plot, panicking on error.
	p2 := plot.New()

	// Set titles and labels
	p2.Title.Text = "Optimal Actions"
	p2.X.Label.Text = "Time Step"
	p2.Y.Label.Text = "Fraction of Actions Optimally Taken"

	// Create line plotter objects and add them to the plot
	err = plotData(p2, actionSteps, const_thousand_run_average_optimality, "Exponential-Recency Weighted", color.RGBA{R: 255, A: 255})
	if err != nil {
		log.Panic(err)
	}

	err = plotData(p2, actionSteps, samp_avg_thousand_run_average_optimality, "Sample Averaged", color.RGBA{G: 255, A: 255})
	if err != nil {
		log.Panic(err)
	}

	// Save the plot to a PNG file.
	if err := p.Save(8*vg.Inch, 8*vg.Inch, "rewards.png"); err != nil {
		log.Panic(err)
	}

	// Save the plot to a PNG file.
	if err := p2.Save(8*vg.Inch, 8*vg.Inch, "optimality.png"); err != nil {
		log.Panic(err)
	}
}
