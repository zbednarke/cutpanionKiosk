function createWeightChart(chartData) {
    // chartData is expected to be an array of 7 numbers
    const ctx = document.getElementById('weightChartCanvas').getContext('2d');

    new Chart(ctx, {
        type: 'line',
        data: {
            labels: ['Day -6', 'Day -5', 'Day -4', 'Day -3', 'Day -2', 'Day -1', 'Today'],
            datasets: [{
                label: 'Weight (lbs)',
                data: chartData,
                borderColor: 'rgba(54, 162, 235, 1)',  // a nice blue
                backgroundColor: 'rgba(54, 162, 235, 0.2)',
                fill: true,
                tension: 0.1  // makes the line slightly curved
            }]
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            scales: {
                y: {
                    beginAtZero: false
                }
            },
            plugins: {
                legend: {
                    display: true
                },
                title: {
                    display: true,
                    text: 'Weight Chart'
                }
            }
        }
    });
}

// NEW: Create the Deficit Chart
function createDeficitChart(chartData) {
    const container = document.getElementById('deficit-chart');
    container.innerHTML = '<canvas id="deficitChartCanvas"></canvas>';
    const ctx = document.getElementById('deficitChartCanvas').getContext('2d');
    new Chart(ctx, {
        type: 'line',
        data: {
            labels: ['Day -6', 'Day -5', 'Day -4', 'Day -3', 'Day -2', 'Day -1', 'Today'],
            datasets: [{
                label: 'Deficit (kcal)',
                data: chartData,
                borderColor: 'rgba(255, 159, 64, 1)',
                backgroundColor: 'rgba(255, 159, 64, 0.2)',
                fill: true,
                tension: 0.1
            }]
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            scales: { y: { beginAtZero: true } },
            plugins: {
                legend: { display: true },
                title: { display: true, text: 'Deficit Chart' }
            }
        }
    });
}

// NEW: Create the Protein Chart
function createProteinChart(chartData) {
    const container = document.getElementById('protein-chart');
    container.innerHTML = '<canvas id="proteinChartCanvas"></canvas>';
    const ctx = document.getElementById('proteinChartCanvas').getContext('2d');
    new Chart(ctx, {
        type: 'line',
        data: {
            labels: ['Day -6', 'Day -5', 'Day -4', 'Day -3', 'Day -2', 'Day -1', 'Today'],
            datasets: [{
                label: 'Protein (g)',
                data: chartData,
                borderColor: 'rgba(75, 192, 192, 1)',
                backgroundColor: 'rgba(75, 192, 192, 0.2)',
                fill: true,
                tension: 0.1
            }]
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            scales: { y: { beginAtZero: true } },
            plugins: {
                legend: { display: true },
                title: { display: true, text: 'Protein Chart' }
            }
        }
    });
}

// NEW: Create the Calories Chart
function createCaloriesChart(chartData) {
    const container = document.getElementById('calories-chart');
    container.innerHTML = '<canvas id="caloriesChartCanvas"></canvas>';
    const ctx = document.getElementById('caloriesChartCanvas').getContext('2d');
    new Chart(ctx, {
        type: 'line',
        data: {
            labels: ['Day -6', 'Day -5', 'Day -4', 'Day -3', 'Day -2', 'Day -1', 'Today'],
            datasets: [{
                label: 'Calories',
                data: chartData,
                borderColor: 'rgba(153, 102, 255, 1)',
                backgroundColor: 'rgba(153, 102, 255, 0.2)',
                fill: true,
                tension: 0.1
            }]
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            scales: { y: { beginAtZero: true } },
            plugins: {
                legend: { display: true },
                title: { display: true, text: 'Calories Chart' }
            }
        }
    });
}

function UIUpdate() {
    fetch("/api/data")
        .then(response => response.json())
        .then(data => {
            // Quote
            quoteEl.textContent = data.quote
                ? `"${data.quote}"`
                : "No quote available";

            // Workout
            workoutEl.textContent = data.workout_today
                ? data.workout_today.toUpperCase()
                : "No workout info";

            if (data.weight_chart_data && Array.isArray(data.weight_chart_data)) {
                createWeightChart(data.weight_chart_data);
            }
            if (data.deficit_chart_data && Array.isArray(data.deficit_chart_data)) {
                createDeficitChart(data.deficit_chart_data);
            }
            if (data.protein_chart_data && Array.isArray(data.protein_chart_data)) {
                createProteinChart(data.protein_chart_data);
            }
            if (data.calories_chart_data && Array.isArray(data.calories_chart_data)) {
                createCaloriesChart(data.calories_chart_data);
            }
        })
        .catch(err => {
            console.error("Failed to fetch data:", err);
            quoteEl.textContent = "Error fetching data";
            workoutEl.textContent = "Error fetching data";
            scheduleEl.textContent = "Error fetching data";
        });
}


document.addEventListener("DOMContentLoaded", () => {
    const timeEl = document.getElementById("time");
    const quoteEl = document.getElementById("quote");
    const workoutEl = document.getElementById("workout");
    const scheduleEl = document.getElementById("schedule");

    // Update the time every second
    function updateTime() {
        const now = new Date();
        // e.g., "1:45 PM" (no seconds)
        timeEl.textContent = now.toLocaleTimeString([], {
            hour: '2-digit',
            minute: '2-digit'
        });
    }
    updateTime();
    setInterval(updateTime, 1000);
    UIUpdate()
    setInterval(UIUpdate, 60000)

    // Fetch data from /api/data and display

});
