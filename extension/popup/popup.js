// DOM Elements
const detectedTab = document.getElementById("tab-detected");
const activeTab = document.getElementById("tab-active");
const detectedStreamsContent = document.getElementById("detected-streams");
const activeStreamsContent = document.getElementById("active-streams");
const detectedList = document.getElementById("detected-list");
const activeList = document.getElementById("active-list");
const refreshDetectedBtn = document.getElementById("refresh-detected");
const clearDetectedBtn = document.getElementById("clear-detected");
const refreshActiveBtn = document.getElementById("refresh-active");
const streamForm = document.getElementById("stream-form");
const formTitle = document.getElementById("form-title");
const streamAction = document.getElementById("stream-action");
const streamName = document.getElementById("stream-name");
const streamOrigin = document.getElementById("stream-origin");
const streamUrl = document.getElementById("stream-url");
const cancelFormBtn = document.getElementById("cancel-form");
const submitFormBtn = document.getElementById("submit-form");
const streamDetailsForm = document.getElementById("stream-details");
const statusMessage = document.getElementById("status-message");

// Server config elements
const serverConfig = document.getElementById("server-config");
const serverUrl = document.getElementById("server-url");
const testConnectionBtn = document.getElementById("test-connection");
const saveConfigBtn = document.getElementById("save-config");
const connectionStatus = document.getElementById("connection-status");
const contentContainer = document.getElementById("content-container");
const serverSettingsBtn = document.getElementById("server-settings");

// Templates
const detectedStreamTemplate = document.getElementById("detected-stream-template");
const activeStreamTemplate = document.getElementById("active-stream-template");

// API URL - will be updated from server config
let API_SERVER = "http://localhost:8080";

// Current form state
let currentFormData = null;

// Load server configuration
function loadServerConfig() {
  browser.runtime.sendMessage({ action: "getServerConfig" })
    .then(response => {
      if (response && response.serverConfig) {
        serverUrl.value = response.serverConfig.url;
        API_SERVER = response.serverConfig.url;
        
        if (response.serverConfig.isConnected) {
          showContent();
        } else {
          showServerConfig();
        }
      } else {
        showServerConfig();
      }
    })
    .catch(error => {
      console.error("Error loading server config:", error);
      showServerConfig();
    });
}

// Test connection to server
function testConnection() {
  const url = serverUrl.value.trim();
  
  if (!url) {
    setConnectionStatus("Please enter a server URL", true);
    return;
  }
  
  // Validate URL format
  if (!url.match(/^https?:\/\/.+/)) {
    setConnectionStatus("Invalid URL format. Use http://hostname:port", true);
    return;
  }
  
  setConnectionStatus("Testing connection...");
  
  // Update server URL in background
  browser.runtime.sendMessage({ 
    action: "updateServerConfig",
    url: url
  })
  .then(() => {
    // Update local API_SERVER variable
    API_SERVER = url;
    // Test connection
    return browser.runtime.sendMessage({ action: "testServerConnection" });
  })
  .then(result => {
    if (result.success) {
      let message = "Connection successful";
      if (result.status && result.status.streams_count !== undefined) {
        message += ` - ${result.status.streams_count} active streams`;
      }
      setConnectionStatus(message, false);
      // Wait a moment before showing content
      setTimeout(showContent, 1000);
    } else {
      setConnectionStatus(`Connection failed: ${result.error}`, true);
    }
  })
  .catch(error => {
    setConnectionStatus(`Error: ${error.message}`, true);
  });
}

// Set connection status message
function setConnectionStatus(message, isError = false) {
  connectionStatus.textContent = message;
  connectionStatus.className = "connection-status";
  if (isError) {
    connectionStatus.classList.add("error");
  } else if (message.includes("successful")) {
    connectionStatus.classList.add("success");
  }
}

// Show server configuration
function showServerConfig() {
  serverConfig.classList.remove("hidden");
  contentContainer.classList.add("hidden");
}

// Show main content
function showContent() {
  serverConfig.classList.add("hidden");
  contentContainer.classList.remove("hidden");
  loadDetectedStreams();
}

// Tab switching
detectedTab.addEventListener("click", function() {
  setActiveTab("detected");
});

activeTab.addEventListener("click", function() {
  setActiveTab("active");
});

function setActiveTab(tabName) {
  if (tabName === "detected") {
    detectedTab.classList.add("active");
    activeTab.classList.remove("active");
    detectedStreamsContent.classList.add("active");
    activeStreamsContent.classList.remove("active");
    loadDetectedStreams();
  } else {
    activeTab.classList.add("active");
    detectedTab.classList.remove("active");
    activeStreamsContent.classList.add("active");
    detectedStreamsContent.classList.remove("active");
    loadActiveStreams();
  }
}

// Load detected streams from background script
function loadDetectedStreams() {
  detectedList.innerHTML = '<div class="loading">Loading detected streams...</div>';
  browser.runtime.sendMessage({ action: "getDetectedStreams" })
    .then(response => {
      displayDetectedStreams(response.streams);
    })
    .catch(error => {
      detectedList.innerHTML = `<div class="empty-message">Error: ${error.message}</div>`;
    });
}

// Load active streams from server
function loadActiveStreams() {
  activeList.innerHTML = '<div class="loading">Loading active streams...</div>';
  browser.runtime.sendMessage({ action: "getServerStreams" })
    .then(streams => {
      displayActiveStreams(streams);
    })
    .catch(error => {
      activeList.innerHTML = `<div class="empty-message">Error connecting to server: ${error.message}</div>`;
    });
}

// Display detected streams
function displayDetectedStreams(streams) {
  detectedList.innerHTML = '';
  
  if (!streams || streams.length === 0) {
    detectedList.innerHTML = '<div class="empty-message">No M3U8 streams detected. Browse to a page with HLS content.</div>';
    return;
  }
  
  streams.forEach(stream => {
    const streamElement = detectedStreamTemplate.content.cloneNode(true);
    
    // Set stream info
    const urlElement = streamElement.querySelector(".stream-url");
    urlElement.textContent = stream.url;
    urlElement.title = stream.url;
    
    const originText = stream.origin || "No origin detected";
    streamElement.querySelector(".stream-origin").textContent = `Origin: ${originText}`;
    
    // Add button event listener
    const addButton = streamElement.querySelector(".add-stream-btn");
    addButton.addEventListener("click", () => {
      showAddForm(stream);
    });
    
    detectedList.appendChild(streamElement);
  });
}

// Display active streams
function displayActiveStreams(streams) {
  activeList.innerHTML = '';
  
  if (!streams || Object.keys(streams).length === 0) {
    activeList.innerHTML = '<div class="empty-message">No active streams. Add a stream from the Detected tab.</div>';
    return;
  }
  
  for (const [name, stream] of Object.entries(streams)) {
    const streamElement = activeStreamTemplate.content.cloneNode(true);
    
    // Set stream info
    streamElement.querySelector(".stream-name").textContent = name;
    
    const urlElement = streamElement.querySelector(".stream-url");
    urlElement.textContent = stream.url;
    urlElement.title = stream.url;
    
    streamElement.querySelector(".stream-origin").textContent = stream.origin;
    
    const statusElement = streamElement.querySelector(".stream-status");
    statusElement.textContent = `Status: ${stream.status}`;
    if (stream.status === "active") {
      statusElement.classList.add("active");
    } else if (stream.status === "error" || stream.status === "failed") {
      statusElement.classList.add("error");
    }
    
    // Setup buttons
    const updateButton = streamElement.querySelector(".update-stream-btn");
    updateButton.addEventListener("click", () => {
      showUpdateForm(name, stream);
    });
    
    const removeButton = streamElement.querySelector(".remove-stream-btn");
    removeButton.addEventListener("click", () => {
      removeStream(name);
    });
    
    const viewLink = streamElement.querySelector(".view-stream-btn");
    viewLink.href = `${API_SERVER}/stream/${name}/index.m3u8`;
    
    activeList.appendChild(streamElement);
  }
}

// Show form to add a new stream
function showAddForm(stream = null) {
  streamAction.value = "add";
  formTitle.textContent = "Add Stream";
  submitFormBtn.textContent = "Add Stream";
  
  // Clear or pre-fill form
  streamName.value = "";
  streamOrigin.value = stream ? stream.origin : "";
  streamUrl.value = stream ? stream.url : "";
  
  // Store current form data
  currentFormData = stream;
  
  // Show form
  streamForm.classList.remove("hidden");
  streamName.focus();
}

// Show form to update a stream
function showUpdateForm(name, stream) {
  streamAction.value = "update";
  formTitle.textContent = "Update Stream";
  submitFormBtn.textContent = "Update Stream";
  
  // Pre-fill form
  streamName.value = name;
  streamName.disabled = true; // Can't change name when updating
  
  // Make sure we have the correct origin and URL values
  streamOrigin.value = stream.origin || "";
  streamUrl.value = stream.url || "";
  
  // If URL or origin is missing or "unknown", try to use detected streams as options
  if (!stream.url || stream.url === "unknown") {
    // Add dropdown for detected streams
    loadDetectedStreamsForSelection();
  }
  
  // Store current form data
  currentFormData = {
    name: name,
    ...stream
  };
  
  // Show form
  streamForm.classList.remove("hidden");
  streamOrigin.focus();
}

// Load detected streams to use for selection
function loadDetectedStreamsForSelection() {
  // Create a dropdown element if it doesn't exist
  let selectContainer = document.getElementById("detected-streams-select-container");
  if (!selectContainer) {
    // Create container
    selectContainer = document.createElement("div");
    selectContainer.id = "detected-streams-select-container";
    selectContainer.className = "form-group";
    
    // Add label
    const label = document.createElement("label");
    label.textContent = "Select from detected streams:";
    selectContainer.appendChild(label);
    
    // Create select element
    const select = document.createElement("select");
    select.id = "detected-streams-select";
    selectContainer.appendChild(select);
    
    // Add to form before the URL field
    const urlFormGroup = streamUrl.parentElement;
    streamDetailsForm.insertBefore(selectContainer, urlFormGroup);
    
    // Add change event listener
    select.addEventListener("change", () => {
      const selectedValue = select.value;
      if (selectedValue) {
        const [selectedUrl, selectedOrigin] = selectedValue.split("||");
        streamUrl.value = selectedUrl;
        streamOrigin.value = selectedOrigin || "";
      }
    });
  }
  
  // Get the select element
  const select = document.getElementById("detected-streams-select");
  select.innerHTML = '<option value="">-- Select a stream --</option>';
  
  // Get detected streams and populate the dropdown
  browser.runtime.sendMessage({ action: "getDetectedStreams" })
    .then(response => {
      if (response.streams && response.streams.length > 0) {
        response.streams.forEach(stream => {
          const option = document.createElement("option");
          option.value = `${stream.url}||${stream.origin || ""}`;
          
          // Simple display with just the URL
          option.textContent = stream.url.substring(0, 50) + (stream.url.length > 50 ? "..." : "");
          select.appendChild(option);
        });
      } else {
        // No detected streams
        const option = document.createElement("option");
        option.disabled = true;
        option.textContent = "No streams detected";
        select.appendChild(option);
      }
    });
}

// Hide the form
function hideForm() {
  streamForm.classList.add("hidden");
  streamName.disabled = false;
  currentFormData = null;
  
  // Remove the detected streams select container if it exists
  const selectContainer = document.getElementById("detected-streams-select-container");
  if (selectContainer) {
    selectContainer.remove();
  }
}

// Add a new stream
function addStream(name, origin, url) {
  setStatus("Adding stream...");
  
  console.log(`Sending to server - Name: ${name}, Origin: ${origin}, URL: ${url}`);
  
  browser.runtime.sendMessage({
    action: "addStream",
    name: name,
    origin: origin,
    url: url
  })
  .then(response => {
    if (response.error) {
      setStatus(`Error: ${response.error}`, true);
    } else {
      setStatus(`Stream '${name}' added successfully`);
      hideForm();
      
      // Switch to active tab
      setActiveTab("active");
    }
  })
  .catch(error => {
    setStatus(`Error: ${error.message}`, true);
  });
}

// Update a stream
function updateStream(name, origin, url) {
  setStatus("Updating stream...");
  
  browser.runtime.sendMessage({
    action: "updateStream",
    name: name,
    origin: origin,
    url: url
  })
  .then(response => {
    if (response.error) {
      setStatus(`Error: ${response.error}`, true);
    } else {
      setStatus(`Stream '${name}' updated successfully`);
      hideForm();
      loadActiveStreams();
    }
  })
  .catch(error => {
    setStatus(`Error: ${error.message}`, true);
  });
}

// Remove a stream
function removeStream(name) {
  if (!confirm(`Are you sure you want to remove the stream '${name}'?`)) {
    return;
  }
  
  setStatus("Removing stream...");
  
  browser.runtime.sendMessage({
    action: "removeStream",
    name: name
  })
  .then(response => {
    if (response.error) {
      setStatus(`Error: ${response.error}`, true);
    } else {
      setStatus(`Stream '${name}' removed successfully`);
      
      // Force a hard refresh after a delay
      setTimeout(() => {
        // First, clear any existing streams in the UI
        activeList.innerHTML = '<div class="loading">Refreshing list...</div>';
        
        // Then fetch the new list from server
        setTimeout(loadActiveStreams, 500);
      }, 500);
    }
  })
  .catch(error => {
    setStatus(`Error: ${error.message}`, true);
  });
}

// Set status message
function setStatus(message, isError = false) {
  statusMessage.textContent = message;
  statusMessage.style.color = isError ? "#f44336" : "#666";
  
  // Clear after 5 seconds
  setTimeout(() => {
    statusMessage.textContent = "";
  }, 5000);
}

// Event Listeners
refreshDetectedBtn.addEventListener("click", loadDetectedStreams);
clearDetectedBtn.addEventListener("click", () => {
  if (confirm("Clear all detected streams?")) {
    browser.runtime.sendMessage({ action: "clearDetectedStreams" })
      .then(() => {
        loadDetectedStreams();
        setStatus("Detected streams cleared");
      });
  }
});

refreshActiveBtn.addEventListener("click", loadActiveStreams);

cancelFormBtn.addEventListener("click", hideForm);

streamDetailsForm.addEventListener("submit", event => {
  event.preventDefault();
  
  const name = streamName.value.trim();
  const origin = streamOrigin.value.trim();
  const url = streamUrl.value.trim();
  
  if (!name || !origin || !url) {
    setStatus("All fields are required", true);
    return;
  }
  
  if (streamAction.value === "add") {
    addStream(name, origin, url);
  } else {
    updateStream(name, origin, url);
  }
});

// Initialize
document.addEventListener("DOMContentLoaded", () => {
  loadServerConfig();
  
  // Event listeners for server config
  testConnectionBtn.addEventListener("click", testConnection);
  saveConfigBtn.addEventListener("click", () => {
    testConnection();
  });
  
  serverSettingsBtn.addEventListener("click", showServerConfig);
});