package main

//  Eric's FIRE Calculator
//  Financial Independenc Retire Early is a dream of many
//

import (
	"fmt"
)

// Computational Routines
// func SimpleGrowth
// func applyTaxes

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

	fmt.Println(totalMonths, contributionMonths, drawDownStart)

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
