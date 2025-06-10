package main

//  Eric's FIRE Calculator
//  Financial Independenc Retire Early is a dream of many
//  Let's try to find out if its possible

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http" // Import the net/http package for web server
	"os"
	"path/filepath"
	"strconv"
)

// Define the path to your HTML template
const htmlTemplatePath = "templates/index.html" // Relative to your project root
const staticFilesDir = "static"                 // Define the directory where your static files are

var err error
var fireIn *FinancialICs // declares input struct
var fireOut *FinancialResults
var fireTimeSeries *TimeSeriesICs

// RootHandler serves the HTML form.
func RootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// --- Get CSS file modification time ---
	cssPath := filepath.Join(staticFilesDir, "style.css")
	fileInfo, err := os.Stat(cssPath)
	var cssVersion string
	if err != nil {
		log.Printf("Error getting CSS file info for %s: %v. Using default version.", cssPath, err)
		cssVersion = "1" // Fallback version if file not found/error
	} else {
		cssVersion = fmt.Sprintf("%d", fileInfo.ModTime().Unix()) // Use Unix timestamp as version
	}

	// --- Get JavaScript file modification time (NEW) ---
	jsPath := filepath.Join(staticFilesDir, "main.js") // Adjust filename if different
	jsFileInfo, err := os.Stat(jsPath)
	var jsVersion string
	if err != nil {
		log.Printf("Error getting JS file info for %s: %v. Using default version.", jsPath, err)
		jsVersion = "1" // Fallback version
	} else {
		jsVersion = fmt.Sprintf("%d", jsFileInfo.ModTime().Unix()) // Use Unix timestamp
	}

	// Data to pass to the template
	css_js_version := struct {
		CSSVersion string
		JSVersion  string
	}{
		CSSVersion: cssVersion,
		JSVersion:  jsVersion,
	}

	// Parse the HTML template
	tmpl, err := template.ParseFiles(htmlTemplatePath)
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, "Internal Server Error: Could not load form", http.StatusInternalServerError)
		return
	}

	// Execute the template, sending it to the client
	err = tmpl.Execute(w, css_js_version) // nil because we're not passing any dynamic data to the form yet
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal Server Error: Could not render form", http.StatusInternalServerError)
	}
	log.Println("RootHandler: Served index.html successfully.") // DEBUG
}

// plotHandler generates the plot and serves it as a PNG image via HTTP
func PlotHandler(w http.ResponseWriter, r *http.Request) {

	log.Printf("PlotHandler: Received request for path: %s, method: %s", r.URL.Path, r.Method) // DEBUG

	// We expect GET requests from the JavaScript fetch API
	if r.Method != http.MethodGet {
		log.Printf("PlotHandler: Method mismatch. Expected GET, got %s", r.Method) // DEBUG
		http.Error(w, "Method not allowed. Use GET for plot generation.", http.StatusMethodNotAllowed)
		return
	}

	// Parse the form data (for POST requests) or query parameters (for GET requests)
	// r.ParseForm() must be called before accessing r.Form, r.PostForm, or r.FormValue
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form data: "+err.Error(), http.StatusBadRequest)
		return
	}
	log.Println("PlotHandler: Method is GET, proceeding to parse parameters.") // DEBUG

	getFloatParam := func(paramName string) (float64, error) {
		//	s := r.FormValue(paramName) // r.FormValue works for both GET query and POST form data
		s := r.URL.Query().Get(paramName) // r.FormValue works for both GET query and POST form data
		if s == "" {
			return 0, fmt.Errorf("missing parameter: %s", paramName)
		}
		val, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid value for %s: %w", paramName, err)
		}
		log.Printf("PlotHandler: Parsed %s: %f", paramName, val) // DEBUG
		return val, nil
	}

	// Helper function to get an int parameter from form/query, with error handling
	getIntParam := func(paramName string) (int, error) {
		//		s := r.FormValue(paramName)
		s := r.URL.Query().Get(paramName) // r.FormValue works for both GET query and POST form dataa
		if s == "" {
			return 0, fmt.Errorf("missing parameter: %s", paramName)
		}
		val, err := strconv.Atoi(s)
		if err != nil {
			return 0, fmt.Errorf("invalid value for %s: %w", paramName, err)
		}
		log.Printf("PlotHandler: Parsed %s: %d", paramName, val) // DEBUG
		return val, nil
	}

	// initialize fire pointers
	fireIn = &FinancialICs{}
	fireTimeSeries = &TimeSeriesICs{}
	fireOut = &FinancialResults{}

	// Get parameters for GUI
	// For the moment these are single values
	fireIn.InitialCapital, err = getFloatParam("initialCapital")
	if err != nil {
		log.Printf("PlotHandler Error: %v", err) // DEBUG
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fireIn.MonthlyContribution, err = getFloatParam("monthlyContribution")
	if err != nil {
		log.Printf("PlotHandler Error: %v", err) // DEBUG
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fireIn.AnnualGrowthRate, err = getFloatParam("annualGrowthRate")
	if err != nil {
		log.Printf("PlotHandler Error: %v", err) // DEBUG
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fireIn.AnnualGrowthRate = fireIn.AnnualGrowthRate / 100. // convert from percent to fraction
	fireIn.MonthlyPension, err = getFloatParam("monthlyPension")
	if err != nil {
		log.Printf("PlotHandler Error: %v", err) // DEBUG
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fireIn.ExpectedPensionAge, err = getIntParam("expectedPensionAge")
	if err != nil {
		log.Printf("PlotHandler Error: %v", err) // DEBUG
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fireIn.ContributionYears, err = getIntParam("contributionYears")
	if err != nil {
		log.Printf("PlotHandler Error: %v", err) // DEBUG
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fireIn.CurrentAge, err = getIntParam("currentAge")
	if err != nil {
		log.Printf("PlotHandler Error: %v", err) // DEBUG
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fireIn.DrawDownAge, err = getIntParam("drawDownAge")
	if err != nil {
		log.Printf("PlotHandler Error: %v", err) // DEBUG
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fireIn.MonthlyDrawAmount, err = getFloatParam("monthlyDrawAmount")
	if err != nil {
		log.Printf("PlotHandler Error: %v", err) // DEBUG
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fireIn.ExpectedDeathAge, err = getIntParam("expectedDeathAge")
	if err != nil {
		log.Printf("PlotHandler Error: %v", err) // DEBUG
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Println("PlotHandler: All parameters parsed successfully. Generating plot.") // DEBUG

	// computational logic
	fireTimeSeries, err = PrepareFinanceTimeSeriesICs(fireIn)
	fireOut, err = FinancialProjection(fireTimeSeries)

	if err != nil {
		http.Error(w, fmt.Sprintf("Error calculating growth: %v", err), http.StatusInternalServerError)
		log.Printf("Error calculating growth: %v", err)
		return
	}

	// set JSON data struct
	dataForJSON := DataForJSON{
		MonthIndex:        fireOut.Months,
		YearsIndex:        fireOut.Years,
		PrincipalData:     fireOut.Principal,
		ContributionsData: fireOut.Contributions,
		TakeHomeData:      fireOut.TakeHome,
		Title:             fmt.Sprintf("Retirement Projection (Initial: $%.2f, Growth: %.2f%%)", fireIn.InitialCapital, fireIn.AnnualGrowthRate),
		XLabel:            "Age",
		YLabel:            "Portfolio Value ($)",
	}

	// --- Serve JSON data instead of the plot image ---
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(dataForJSON); err != nil {
		log.Printf("Error encoding JSON: %v", err)
		http.Error(w, "Internal Server Error: Could not send data", http.StatusInternalServerError)
		return
	}
	log.Println("Monthly data sent as JSON.")

}

// StartPlottingServer registers all handlers, including for static files.
func StartPlottingServer(port string) {
	// 1. Serve static files from the /static/ URL path
	// http.Dir(staticFilesDir) creates a file system rooted at 'static'
	// http.StripPrefix("/static/", ...) removes the /static/ prefix from the request path
	// when looking for files within the 'static' directory.
	fs := http.FileServer(http.Dir(staticFilesDir))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	log.Printf("Dataviz: Serving static files from /%s/", staticFilesDir) // DEBUG

	// 2. Serve the HTML form at the root URL
	http.HandleFunc("/", RootHandler)

	// 3. Handle form submissions and plot generation at the /plot URL
	http.HandleFunc("/plot", PlotHandler)

	log.Printf("Dataviz: Serving Fire Calculator at http://localhost%s/", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func main() {

	// Print current working directory for debugging if needed
	if dir, err := os.Getwd(); err == nil {
		fmt.Printf("Current working directory: %s\n", dir)
	}

	//	// Register the handler for the "/plot" URL path
	//	http.HandleFunc("/plot", dataviz.PlotHandler)

	port := ":8080"
	fmt.Printf("Server starting on http://localhost%s/plot\n", port)
	fmt.Println("Open your browser to this URL and refresh after code changes.")
	log.Println("Starting Fire Calculator web server...")
	//	dataviz.StartPlottingServer(port)
	StartPlottingServer(port)

	// Start the HTTP server
	// log.Fatal will print the error and exit if ListenAndServe fails
	//	log.Fatal(http.ListenAndServe(port, nil))
}
