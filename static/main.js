// Store a reference to the chart instance globally or in a scope accessible by the event listener
// This allows us to destroy and re-create the chart when new data arrives
let myChart;

document.getElementById('calculatorForm').addEventListener('submit', async function(event) {
    event.preventDefault(); // Prevent default form submission

    const form = event.target;
    const formData = new FormData(form);
    const queryParams = new URLSearchParams(formData).toString();

    const chartCanvas = document.getElementById('financialChart');
    const errorMessageDiv = document.getElementById('errorMessage');

    errorMessageDiv.textContent = ''; // Clear previous error messages
    chartCanvas.style.display = 'block'; // Ensure canvas is visible if previously hidden by an error

    try {
        const response = await fetch(`/plot?${queryParams}`);

        if (!response.ok) {
            const errorText = await response.text();
            throw new Error(errorText || `HTTP error! status: ${response.status}`);
        }

        const data = await response.json(); // Parse the JSON response

        // --- Input Validation for received data ---
        if (!data.months || !Array.isArray(data.months) ||
            !data.years  || !Array.isArray(data.years) ||
            !data.principal || !Array.isArray(data.principal) ||
            !data.contributions || !Array.isArray(data.contributions)) {
            throw new Error("Invalid data format received from server: missing or malformed arrays.");
        }
        if (data.months.length === 0 || data.principal.length === 0 || data.contributions.length === 0) {
            throw new Error("No data points generated. Please check input values.");
        }
        if (data.months.length !== data.principal.length || data.months.length !== data.contributions.length) {
            throw new Error("Data arrays have inconsistent lengths.");
        }

        // Prepare data for Chart.js
        const labels = data.months.map(month => `Month ${month}`); // Format labels for X-axis
//        const labels = data.years.map(year => `Year ${year}`); // Format labels for X-axis
        const principalData = data.principal;
        const contributionsData = data.contributions;

        // If a chart already exists, destroy it before creating a new one
        if (myChart) {
            myChart.destroy();
        }

        // Create the chart using Chart.js
        myChart = new Chart(chartCanvas, { // Assign the new chart to myChart
            type: 'line',
            data: {
                labels: labels,
                datasets: [{
                    label: data.title,
                    data: principalData,
                    borderColor: 'blue',
                    backgroundColor: 'rgba(0, 0, 255, 0.1)', // Optional: light fill
                    fill: false, // Set to true if you want an area chart
                    pointRadius: 1, // Smaller dots for monthly data
                    pointHoverRadius: 6,
                    pointHitRadius: 10,
                },
                {
                    label: 'Total Contributions', // Label for the contributions line
                     data: contributionsData,
                    borderColor: 'green', // Different color for contributions
                    backgroundColor: 'rgba(0, 128, 0, 0.1)',
                    fill: false,
                    pointRadius: 1,
                    pointHoverRadius: 6,
                    pointHitRadius: 10,
                }]
            },
            options: {
                responsive: true, // Make chart responsive to container size
                maintainAspectRatio: false, // Allow aspect ratio to change
                scales: {
                    x: {
                        title: {
                            display: true,
                            text: data.xLabel
                        }
                    },
                    y: {
                        title: {
                            display: true,
                            text: data.yLabel
                        },
                        ticks: {
                            callback: function(value, index, values) {
                                return '$' + value.toLocaleString(); // Format with commas
                             }
                        }
                    }
                },
                plugins: {
                    tooltip: {
                        callbacks: {
                            label: function(context) {
                                let label = context.dataset.label || '';
                                if (label) {
                                    label += ': ';
                                }
                                if (context.parsed.y !== null) {
                                    label += '$' + context.parsed.y.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 }); // Format with commas and 2 decimals
                                }
                                return label;
                            },
                            title: function(context) {
                                 // context.label is typically the X-axis label (age string)
                                const age = parseFloat(context.label);
                                const years = Math.floor(age);
                                const months = Math.round((age - years) * 12); // Calculate months from fractional part
                                return `Age: ${years} years, ${months} months`;
                            }
                        }
                    },
                    legend: {
                        display: false // You only have one dataset, so legend might be redundant
                    }
                }
            }
        });

    } catch (error) {
        console.error("Error fetching or rendering data:", error);
        errorMessageDiv.textContent = `Error: ${error.message}`;
        chartCanvas.style.display = 'none'; // Hide chart on error
         if (myChart) { // Destroy existing chart if error occurs
            myChart.destroy();
        }
    }
});

// Trigger initial plot generation on page load with default values
window.addEventListener('load', () => {
    document.getElementById('calculatorForm').dispatchEvent(new Event('submit'));
 });
   
