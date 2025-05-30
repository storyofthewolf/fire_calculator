package dataviz

//  Eric's FIRE Calculator
//  Financial Independenc Retire Early is a dream of many
//  Let's try to find out if its possible

import (
	"encoding/json"
	"fire_calculator/compute"
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

// define struct of financial parameters
/*type financeMain struct {
	months []int,
	years []float64,
	prinicipal []float64,
	contributions []float64,
	monthlyDrawDown []float64,
	monthlyPension []float64
	annualGrowthRate []float64
}
// define struct of taxR rates
type taxRates struct {
	capitalGains float64,
	incomeBrackets []float64,
	federalTaxes []float64
} */

// definte struct for JSON
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

	// Get parameters
	initialCapital, err := getFloatParam("initialCapital")
	if err != nil {
		log.Printf("PlotHandler Error: %v", err) // DEBUG
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	monthlyContribution, err := getFloatParam("monthlyContribution")
	if err != nil {
		log.Printf("PlotHandler Error: %v", err) // DEBUG
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	annualGrowthRate, err := getFloatParam("annualGrowthRate")
	if err != nil {
		log.Printf("PlotHandler Error: %v", err) // DEBUG
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	annualGrowthRate = annualGrowthRate / 100. // convert from percent to fraction

	contributionYears, err := getIntParam("contributionYears")
	if err != nil {
		log.Printf("PlotHandler Error: %v", err) // DEBUG
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	currentAge, err := getIntParam("currentAge")
	if err != nil {
		log.Printf("PlotHandler Error: %v", err) // DEBUG
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	drawDownAge, err := getIntParam("drawDownAge")
	if err != nil {
		log.Printf("PlotHandler Error: %v", err) // DEBUG
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	monthlyDrawAmount, err := getFloatParam("monthlyDrawAmount")
	if err != nil {
		log.Printf("PlotHandler Error: %v", err) // DEBUG
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	expectedDeathAge, err := getIntParam("expectedDeathAge")
	if err != nil {
		log.Printf("PlotHandler Error: %v", err) // DEBUG
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	monthlyPension, err := getFloatParam("monthlyPension")
	if err != nil {
		log.Printf("PlotHandler Error: %v", err) // DEBUG
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	expectedPensionAge, err := getIntParam("expectedPensionAge")
	if err != nil {
		log.Printf("PlotHandler Error: %v", err) // DEBUG
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Println("PlotHandler: All parameters parsed successfully. Generating plot.") // DEBUG

	// call computational logic as needed
	principal, contributions, months, years, takeHome, err := compute.SimpleGrowth(initialCapital,
		monthlyContribution, annualGrowthRate,
		contributionYears, currentAge, drawDownAge,
		monthlyDrawAmount, expectedDeathAge,
		monthlyPension, expectedPensionAge)

	if err != nil {
		http.Error(w, fmt.Sprintf("Error calculating growth: %v", err), http.StatusInternalServerError)
		log.Printf("Error calculating growth: %v", err)
		return
	}

	dataForJSON := DataForJSON{
		MonthIndex:        months,
		YearsIndex:        years,
		PrincipalData:     principal,
		ContributionsData: contributions,
		TakeHomeData:      takeHome,
		Title:             fmt.Sprintf("Retirement Projection (Initial: $%.2f, Growth: %.2f%%)", initialCapital, annualGrowthRate),
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
	log.Println("Take Home", takeHome)
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
