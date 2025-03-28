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
  
    // Fetch data from /api/data and display
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
  
        // Schedule (if the backend provides an array of events)
        // Example data.schedule = [
        //   { startTime: "9:00 AM", endTime: "10:00 AM", title: "Morning Workout" },
        //   { startTime: "1:00 PM", endTime: "2:00 PM", title: "Lunch with Sam" }
        // ]
        if (data.schedule && Array.isArray(data.schedule) && data.schedule.length > 0) {
          data.schedule.forEach(event => {
            const eventDiv = document.createElement("div");
            eventDiv.classList.add("schedule-item");
            eventDiv.textContent = `${event.startTime} - ${event.endTime}: ${event.title}`;
            scheduleEl.appendChild(eventDiv);
          });
        } else {
          scheduleEl.textContent = "No events scheduled";
        }
      })
      .catch(err => {
        console.error("Failed to fetch data:", err);
        quoteEl.textContent = "Error fetching data";
        workoutEl.textContent = "Error fetching data";
        scheduleEl.textContent = "Error fetching data";
      });
  });
  