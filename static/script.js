// script.js

document.addEventListener('DOMContentLoaded', function() {
    const monitoringDataElem = document.getElementById('monitoring-data');
  
    // Function to handle SSE events
    function handleSSE(event) {
      const eventData = JSON.parse(event.data);
      // Update the UI with the received monitoring data
      const timestamp = new Date(eventData.timestamp).toLocaleString();
      const html = `
        <p>Timestamp: ${timestamp}</p>
        <p>CPU Usage: ${eventData.cpu}%</p>
        <p>Memory Usage: ${eventData.memory}%</p>
      `;
      monitoringDataElem.innerHTML = html;
    }
  
    // Connect to SSE endpoint
    const eventSource = new EventSource('/events');
    
    // Add event listener to handle incoming SSE events
    eventSource.addEventListener('message', handleSSE);
  
    // Handle errors
    eventSource.onerror = function(event) {
      console.error('SSE error:', event);
      // Attempt to reconnect
      eventSource.close();
    };
  });
  