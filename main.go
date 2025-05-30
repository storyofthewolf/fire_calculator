package main

//  Eric's FIRE Calculator
//  Financial Independenc Retire Early is a dream of many
//  Let's try to find out if its possible

//I would like to add
//2. add social-security and pension benefit consideration
//3. add tax consideration, combined with draw-down and SS/P benefits to give take home
//4. add non-constant contributions, draw down
//5. perhaps add a stochastic component to the annual growth rate, since the market doesnt always go up

import (
	"fire_calculator/dataviz"
	"fmt"
	"log"
	"os" // For checking current working directory if needed
)

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
	dataviz.StartPlottingServer(port)

	// Start the HTTP server
	// log.Fatal will print the error and exit if ListenAndServe fails
	//	log.Fatal(http.ListenAndServe(port, nil))
}
