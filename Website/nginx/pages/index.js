// main.js
document.addEventListener("DOMContentLoaded", () => {
    // This function is called once the HTML is fully loaded.
    fetchAndDisplayStats();
    //refresh stats every 5 seconds+
    setInterval(fetchAndDisplayStats, 5000);

  });
  
  function fetchAndDisplayStats() {
    fetch("/api/stats")  // Nginx will proxy /api/stats to Go backend
      .then(response => response.json())
      .then(data => {
        // Suppose your backend returns JSON like:
        // {
        //   "avgLatencyMs": 123,
        //   "projectsCount": 42,
        //   "timeSaved": 9999
        // }
  
        document.getElementById("avg-latency").innerText = data.avgLatencyMs + " ms";
        document.getElementById("projects").innerText = data.projectsCount;
        document.getElementById("time-saved").innerText = data.timeSaved + "+";
      })
      .catch(err => {
        console.error("Error fetching stats:", err);
      });
  }
  