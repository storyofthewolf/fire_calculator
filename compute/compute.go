package compute

//  Eric's FIRE Calculator
//  Financial Independenc Retire Early is a dream of many
//

import (
	"fmt"
)

// Computational Routines
// func SimpleGrowth

// CompoundingGrowth function (copy it here if it's not in a separate package)
// Currently this assumes an initial capital amount
// monthly compounding interest at a fixed growth rate
// monthly contributions to investment/retirement accounts
// number of years contributing to retiremett accouns
// current age

func SimpleGrowth(
	initialCapital float64, // intial amount of money
	monthlyContribution float64, // fixed monthly contribution to retirement/investment accounts
	annualGrowthRate float64, // fixed annual growth rate assumption
	contributionYears int, // number of years expected to contribute to accounts (before retirement)
	currentAge int, // your current ages
	expectedDeathAge int, // when you expect to die
) (
	totalPrincipal []float64,
	totalContributions []float64,
	monthsElapsed []int,
	err error,
) {
	if initialCapital < 0 || monthlyContribution < 0 || annualGrowthRate < 0 {
		return nil, nil, nil, fmt.Errorf("initial principal, monthly contribution, and annual growth rate cannot be negative")
	}
	if currentAge < 0 {
		return nil, nil, nil, fmt.Errorf("age must be greater than zero")
	}

	monthlyGrowthRate := annualGrowthRate / 12.0
	totalMonths := (expectedDeathAge - currentAge) * 12
	contributionMonths := contributionYears * 12

	// slices for running counts of principal accumulation and contributions
	totalPrincipal = make([]float64, 0, totalMonths)
	totalContributions = make([]float64, 0, totalMonths)
	monthsElapsed = make([]int, 0, totalMonths)

	// at t=0 the current principal is the intial capital on had
	currentPrincipal := initialCapital
	currentContributions := initialCapital

	for month := 1; month <= totalMonths; month++ {
		// if still during contributing period add monthly contribution to current principal tally
		// and add month contribution to current contributions tally
		if month <= contributionMonths {
			currentPrincipal += monthlyContribution
			currentContributions += monthlyContribution
		}
		// accrue 1 month of interest
		currentPrincipal *= (1 + monthlyGrowthRate)
		// incrrement slices
		totalPrincipal = append(totalPrincipal, currentPrincipal)
		totalContributions = append(totalContributions, currentContributions)
		monthsElapsed = append(monthsElapsed, month)
	}

	return totalPrincipal, totalContributions, monthsElapsed, nil
}
