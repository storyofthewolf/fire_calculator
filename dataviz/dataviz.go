package dataviz

//  Eric's FIRE Calculator
//  Financial Independenc Retire Early is a dream of many
//  Let's try to find out if its possible

import (
	"fmt"
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

	"fire_calculator/compute"
)

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

	initialInv := 360000.0
	monthlyCont := 3000.0
	annualRate := 0.05
	savingYears := 30
	currentAge := 41

	principal, contributions, months, err := compute.SimpleGrowth(initialInv, monthlyCont, annualRate, savingYears, currentAge)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error calculating growth: %v", err), http.StatusInternalServerError)
		log.Printf("Error calculating growth: %v", err)
		return
	}

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

	p := plot.New()

	p.Title.Text = fmt.Sprintf("Investment Growth Over %d Years (%.2f%% Annual Rate)", savingYears, annualRate*100)
	p.X.Label.Text = "Age (years)"
	p.Y.Label.Text = "Amount ($)"

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
	c := vgimg.New(8*vg.Inch, 6*vg.Inch)

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
}
