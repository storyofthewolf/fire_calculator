package main

//  Eric's FIRE Calculator
//  Financial Independenc Retire Early is a dream of many
//  Let's try to find out if its possible

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
