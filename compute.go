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

// CompoundingGrowth function (copy it here if it's not in a separate package)
// Currently this assumes an initial capital amount
// monthly compounding interest at a fixed growth rate
// monthly contributions to investment/retirement accounts
// number of years contributing to retiremett accouns
// current age

func SimpleGrowth(input *FinancialICs) (*FinancialResults, error) {

	if input.InitialCapital < 0 || input.MonthlyContribution < 0 || input.AnnualGrowthRate < 0 {
		return nil, fmt.Errorf("initial principal, monthly contribution, and annual growth rate cannot be negative")
	}
	if input.CurrentAge < 0 {
		return nil, fmt.Errorf("age must be greater than zero")
	}

	// define local variables
	monthlyGrowthRate := input.AnnualGrowthRate / 12.0
	totalMonths := (input.ExpectedDeathAge-input.CurrentAge)*12 + 1
	contributionMonths := input.ContributionYears * 12
	drawDownStart := (input.DrawDownAge - input.CurrentAge) * 12
	pensionStart := (input.ExpectedPensionAge - input.CurrentAge) * 12

	// slices for running counts of principal accumulation and contributions
	totalPrincipal := make([]float64, 0, totalMonths)
	totalContributions := make([]float64, 0, totalMonths)
	takeHome := make([]float64, 0, totalMonths)
	monthsElapsed := make([]int, 0, totalMonths)
	yearsElapsed := make([]float64, 0, totalMonths)

	// at t=0 the current principal is the intial capital on had
	currentPrincipal := input.InitialCapital
	currentContributions := input.InitialCapital
	currentTakeHome := 0.0

	for month := 0; month <= totalMonths; month++ {
		// if still during contributing period add monthly contribution to current principal tally
		// and add month contribution to current contributions tally

		// incrrement slices
		totalPrincipal = append(totalPrincipal, currentPrincipal)
		totalContributions = append(totalContributions, currentContributions)
		monthsElapsed = append(monthsElapsed, month)
		yearsElapsed = append(yearsElapsed, float64(month)/12.)
		takeHome = append(takeHome, currentTakeHome)

		if month <= contributionMonths {
			currentPrincipal += input.MonthlyContribution
			currentContributions += input.MonthlyContribution
		}
		if month >= drawDownStart {
			currentPrincipal -= input.MonthlyDrawAmount
			currentTakeHome = input.MonthlyDrawAmount
		}
		if month >= drawDownStart {
			currentPrincipal -= input.MonthlyDrawAmount
			currentTakeHome = input.MonthlyDrawAmount
		}
		if month >= pensionStart {
			currentTakeHome += input.MonthlyPension
		}
		if currentPrincipal <= 0.0 {
			currentPrincipal = 0.0
		}
		// accrue 1 month of interest
		currentPrincipal *= (1 + monthlyGrowthRate)

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
