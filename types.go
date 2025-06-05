package main

// define structure types used through out the program

// ------- JSON data struct ------
type DataForJSON struct {
	MonthIndex        []int     `json:"months"`
	YearsIndex        []float64 `json:"years"`
	PrincipalData     []float64 `json:"principal"`
	ContributionsData []float64 `json:"contributions"`
	TakeHomeData      []float64 `json:"takeHome"`
	Title             string    `json:"title"`
	XLabel            string    `json:"xLabel"`
	YLabel            string    `json:"yLabel"`
}

// GrowthInput groups all parameters needed to run the growth simulation.
type FinancialICs struct {
	InitialCapital      float64 // Initial amount of money
	MonthlyContribution float64 // Fixed monthly contribution to retirement/investment accounts
	AnnualGrowthRate    float64 // Fixed annual growth rate assumption (e.g., 0.07 for 7%)
	ContributionYears   int     // Number of years expected to contribute to accounts (before retirement)
	CurrentAge          int     // Your current age
	DrawDownAge         int     // Age at which you start drawing down funds
	MonthlyDrawAmount   float64 // Fixed monthly draw amount during retirement
	ExpectedDeathAge    int     // When you expect to die
	MonthlyPension      float64 // Fixed monthly draw amount during retirement
	ExpectedPensionAge  int     // age when starting to take Pension or Social Security
}

// GrowthOutput holds all the results from the growth simulation.
type FinancialResults struct {
	Principal     []float64 // Principal balance over time (e.g., end-of-month balances)
	Contributions []float64 // Cumulative contributions over time
	TakeHome      []float64 // monthly take home amount
	Months        []int     // Month numbers (e.g., 1, 2, 3...)
	Years         []float64 // Years elapsed from start
}
