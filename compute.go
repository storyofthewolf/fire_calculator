package main

//  Eric's FIRE Calculator
//  Financial Independenc Retire Early is a dream of many
//

import (
	"fmt"
	// "time"
	// "math/rand"
)

// Computational Routines
// func SimpleGrowth
// func applyTaxes

// generateRandomGrowthRates generates a slice of random monthly growth rates
// based on a normal distribution.
// mean: the average annual growth rate (e.g., 0.08 for 8%)
// stdDev: the standard deviation of annual growth rate
// numMonths: the number of months to generate growth rates for
/* func generateRandomGrowthRates(mean, stdDev float64, numMonths int) []float64 {
	// Create a new random number generator.
	// Using time.Now().UnixNano() ensures a different sequence each time.
	// For reproducible results, you'd use a fixed seed, e.g., rand.NewSource(42)
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	// Adjust mean and stdDev for monthly rates (assuming annual inputs)
	// Simple approximation: divide by 12 for mean, sqrt(12) for stdDev
	// For more rigorous financial modeling, consider continuous compounding or other methods.
	monthlyMean := mean / 12.0
	monthlyStdDev := stdDev / (12.0 * 0.5) // A common rule of thumb for converting annual to monthly std dev (dividing by sqrt(12))

	growthRates := make([]float64, numMonths)

	for i := 0; i < numMonths; i++ {
		// NormFloat64 returns a normally distributed float64 with mean 0 and stddev 1.
		// We then scale and shift it to our desired mean and stdDev.
		rate := r.NormFloat64()*monthlyStdDev + monthlyMean
		growthRates[i] = rate
	}

	return growthRates
} */

func PrepareFinanceTimeSeriesICs(input *FinancialICs) (*TimeSeriesICs, error) {

	if input.InitialCapital < 0 || input.MonthlyContribution < 0 {
		return nil, fmt.Errorf("initial principal, and monthly contribution cannot be negative")
	}
	if input.CurrentAge < 0 {
		return nil, fmt.Errorf("age must be greater than zero")
	}

	// define local
	totalMonths := (input.ExpectedDeathAge-input.CurrentAge)*12 + 1
	contributionMonths := input.ContributionYears * 12
	drawDownStart := (input.DrawDownAge - input.CurrentAge) * 12
	pensionStart := (input.ExpectedPensionAge - input.CurrentAge) * 12

	// slices for running counts of principal accumulation and contributions
	monthlyContribution := make([]float64, 0, totalMonths)
	monthlyGrowthRate := make([]float64, 0, totalMonths)
	monthlyDrawAmount := make([]float64, 0, totalMonths)
	monthlyPension := make([]float64, 0, totalMonths)

	for month := 0; month <= totalMonths; month++ {
		// if still during contributing period add monthly contribution to current principal tally
		// and add month contribution to current contributions tally
		if month <= contributionMonths {
			monthlyContribution = append(monthlyContribution, input.MonthlyContribution)
		} else {
			monthlyContribution = append(monthlyContribution, 0.0)
		}
		if month >= drawDownStart {
			monthlyDrawAmount = append(monthlyDrawAmount, input.MonthlyDrawAmount)
		} else {
			monthlyDrawAmount = append(monthlyDrawAmount, 0.0)
		}
		if month >= pensionStart {
			monthlyPension = append(monthlyPension, input.MonthlyPension)
		} else {
			monthlyPension = append(monthlyPension, 0.0)
		}
		monthlyGrowthRate = append(monthlyGrowthRate, input.AnnualGrowthRate/12.)
	}

	output := &TimeSeriesICs{
		InitialCapital:      input.InitialCapital,
		MonthlyContribution: monthlyContribution,
		MonthlyGrowthRate:   monthlyGrowthRate,
		MonthlyDrawAmount:   monthlyDrawAmount,
		MonthlyPension:      monthlyPension,
		TotalMonths:         totalMonths,
	}

	return output, nil
}

func FinancialProjection(input *TimeSeriesICs) (*FinancialResults, error) {

	// slices for running counts of principal accumulation and contributions
	totalPrincipal := make([]float64, 0, input.TotalMonths)
	totalContributions := make([]float64, 0, input.TotalMonths)
	takeHome := make([]float64, 0, input.TotalMonths)
	monthsElapsed := make([]int, 0, input.TotalMonths)
	yearsElapsed := make([]float64, 0, input.TotalMonths)

	// at t=0 the current principal is the intial capital on had
	currentPrincipal := input.InitialCapital
	currentContributions := input.InitialCapital
	currentTakeHome := 0.0

	for month := 0; month <= input.TotalMonths; month++ {
		// if still during contributing period add monthly contribution to current principal tally
		// and add month contribution to current contributions tally
		//fmt.Println(totalPrincipal[month], totalContributions[month], totalPrincipal[month])
		// incrrement slices
		totalPrincipal = append(totalPrincipal, currentPrincipal)
		totalContributions = append(totalContributions, currentContributions)
		monthsElapsed = append(monthsElapsed, month)
		yearsElapsed = append(yearsElapsed, float64(month)/12.)
		takeHome = append(takeHome, currentTakeHome)

		currentPrincipal += input.MonthlyContribution[month]
		currentContributions += input.MonthlyContribution[month]
		currentPrincipal -= input.MonthlyDrawAmount[month]
		currentTakeHome = input.MonthlyDrawAmount[month] + input.MonthlyPension[month]
		if currentPrincipal <= 0.0 {
			currentPrincipal = 0.0
		}
		// accrue 1 month of interest
		currentPrincipal *= (1 + input.MonthlyGrowthRate[month])

	}

	//
	output := &FinancialResults{
		Principal:     totalPrincipal,
		Contributions: totalContributions,
		Months:        monthsElapsed,
		Years:         yearsElapsed,
		TakeHome:      takeHome,
	}

	return output, nil
}
