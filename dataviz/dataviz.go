package dataviz

//  Eric's FIRE Calculator
//  Financial Independenc Retire Early is a dream of many
//  Let's try to find out if its possible

import (
	"fmt"
	"html/template"
	"image/color"
	"image/png"

	"log"
	"net/http" // Import the net/http package for web server

	// For checking current working directory if needed
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"

	"strconv"

	"fire_calculator/compute"
)

// Define the path to your HTML template
const htmlTemplatePath = "templates/index.html" // Relative to your project root
const staticFilesDir = "static"                 // Define the directory where your static files are

// RootHandler serves the HTML form.
func RootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Parse the HTML template
	tmpl, err := template.ParseFiles(htmlTemplatePath)
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, "Internal Server Error: Could not load form", http.StatusInternalServerError)
		return
	}

	// Execute the template, sending it to the client
	err = tmpl.Execute(w, nil) // nil because we're not passing any dynamic data to the form yet
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal Server Error: Could not render form", http.StatusInternalServerError)
	}
	log.Println("RootHandler: Served index.html successfully.") // DEBUG
}

// generatePlotData is a helper function to create the plotter.XYs from your slices
// assumes the x axis array is an int, y axis array is a float
func generatePlotData(xData []int, yData []float64, currentAge int) (plotter.XYs, error) {
	if len(xData) != len(yData) {
		return nil, fmt.Errorf("xData and yData slices must have the same length")
	}
	if len(xData) == 0 {
		return nil, fmt.Errorf("input data slices are empty")
	}

	pts := make(plotter.XYs, len(xData))
	for i := range xData {
		pts[i].X = float64(xData[i])/12.0 + float64(currentAge)
		pts[i].Y = yData[i]
		//fmt.Println("%d, %f,%f", i, pts[i].X, pts[i].Y)
	}
	return pts, nil
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
	// Check if the request method is POST
	//if r.Method != http.MethodPost && r.Method != http.MethodGet { // Allow GET for initial image load or direct URL
	//		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	//	return
	//}

	// Parse the form data (for POST requests) or query parameters (for GET requests)
	// r.ParseForm() must be called before accessing r.Form, r.PostForm, or r.FormValue
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form data: "+err.Error(), http.StatusBadRequest)
		return
	}
	log.Println("PlotHandler: Method is GET, proceeding to parse parameters.") // DEBUG

	// Helper function to get a float64 parameter from form/query, with error handling
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

	log.Println("PlotHandler: All parameters parsed successfully. Generating plot.") // DEBUG

	// cmputational logic
	principal, contributions, months, err := compute.SimpleGrowth(initialCapital, monthlyContribution, annualGrowthRate, contributionYears, currentAge, drawDownAge, monthlyDrawAmount, expectedDeathAge)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error calculating growth: %v", err), http.StatusInternalServerError)
		log.Printf("Error calculating growth: %v", err)
		return
	}

	// generate plot data
	principalPoints, err := generatePlotData(months, principal, currentAge)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error preparing cumulative data: %v", err), http.StatusInternalServerError)
		log.Printf("Error preparing cumulative data: %v", err)
		return
	}

	contributionsPoints, err := generatePlotData(months, contributions, currentAge)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error preparing principal data: %v", err), http.StatusInternalServerError)
		log.Printf("Error preparing principal data: %v", err)
		return
	}

	//plot
	p := plot.New()

	p.Title.Text = fmt.Sprintf("Investment Growth Over %d Years (%.2f%% Annual Rate)", totalYears, annualGrowthRate*100)
	p.X.Label.Text = "Age (years)"
	p.Y.Label.Text = "Captial Amount ($)"

	linePrincipal, err := plotter.NewLine(principalPoints)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating principal line: %v", err), http.StatusInternalServerError)
		log.Printf("Error creating principal line: %v", err)
		return
	}
	linePrincipal.Color = color.RGBA{B: 255, A: 255}
	p.Add(linePrincipal)
	p.Legend.Add("Principal", linePrincipal)

	lineContributions, err := plotter.NewLine(contributionsPoints)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating contribution line: %v", err), http.StatusInternalServerError)
		log.Printf("Error creating contribution line: %v", err)
		return
	}
	lineContributions.Color = color.RGBA{R: 255, A: 255}
	p.Add(lineContributions)
	p.Legend.Add("Total Contributions", lineContributions)

	p.Add(plotter.NewGrid())
	p.Legend.Top = true
	p.Legend.Left = true

	// serving plot to HTTP
	w.Header().Set("Content-Type", "image/png") // Tell the browser it's a PNG image

	// Create a PNG drawing canvas that writes directly to the ResponseWriter
	// This canvas satisfies the draw.Canvas interface that p.Draw expects.
	c := vgimg.New(10*vg.Inch, 6*vg.Inch)

	// Draw the plot onto the canvas
	p.Draw(draw.New(c)) // Use draw.New(c) to wrap the vgimg.Png canvas

	// Set the HTTP header to specify the content type.
	w.Header().Set("Content-Type", "image/png")

	// The vgimg.Canvas has a method Image() which returns an image.Image.
	// Alternatively, you can access its public field `img` which is an *image.RGBA (implements image.Image).
	// We then use the standard library's png.Encode to write to the ResponseWriter.
	if err := png.Encode(w, c.Image()); err != nil {
		http.Error(w, "Failed to encode plot to PNG", http.StatusInternalServerError)
		log.Printf("Failed to encode plot to PNG: %v", err)
	}
	log.Println("Plot generated and served via AJAX.")
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
