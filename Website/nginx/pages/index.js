document.addEventListener("DOMContentLoaded", () => {
    const stats = [
      { id: "stat1Value", name: "Avg Site Latency" },
      { id: "stat2Value", name: "Projects" },
      { id: "stat3Value", name: "Time Saved" },
      { id: "stat4Value", name: "Success Rate" },
      { id: "stat5Value", name: "Clients" },
      { id: "stat6Value", name: "Downloads" },
      { id: "stat7Value", name: "Experience" },
      { id: "stat8Value", name: "Uptime" },
      { id: "stat9Value", name: "Awards" },
      { id: "stat10Value", name: "Lines of Code" },
    ];
  
    const fetchStats = async () => {
      try {
        // Replace with your actual API endpoint
        const response = await fetch("https://your-api-endpoint.com/stats");
        
        if (!response.ok) {
          throw new Error(`HTTP error! Status: ${response.status}`);
        }
  
        const data = await response.json();
  
        stats.forEach((stat) => {
          const valueElement = document.getElementById(stat.id);
          const nameElement = document.getElementById(stat.id.replace("Value", "Name"));
  
          if (valueElement && data[stat.id]) {
            valueElement.textContent = data[stat.id]; // Update value
          }
  
          if (nameElement && data[stat.id + "Name"]) {
            nameElement.textContent = data[stat.id + "Name"]; // Update name
          }
        });
      } catch (error) {
        console.error("Error fetching stats:", error);
        // If an error occurs, original values are retained
      }
    };
  
    // Fetch stats immediately and every 30 seconds
    fetchStats();
    setInterval(fetchStats, 30000);
  });
  