// Store for detected streams
let detectedStreams = [];

// Server configuration
let serverConfig = {
  url: "http://localhost:8080",  // Default
  isConnected: false
};

// Filter specifically for .m3u8 requests and capture their headers
browser.webRequest.onBeforeSendHeaders.addListener(
  captureM3U8Headers,
  { urls: ["*://*/*m3u8*"] },
  ["requestHeaders"]
);

// Function to capture M3U8 request headers
function captureM3U8Headers(details) {
  // Only track XHR/media requests (exclude extension's own requests)
  if (details.tabId !== -1 && (details.type === "xmlhttprequest" || details.type === "media")) {
    const url = details.url;
    
    // Only process actual m3u8 URLs
    if (url.includes(".m3u8")) {
      console.log("M3U8 Request detected:", url);
      console.log("Headers:", details.requestHeaders);
      
      // Extract the Origin header
      let origin = "";
      for (const header of details.requestHeaders) {
        if (header.name.toLowerCase() === "origin") {
          origin = header.value;
          console.log("Found Origin header:", origin);
          break;
        }
      }
      
      // Find if we already have this URL in our detected streams
      const existingIndex = detectedStreams.findIndex(stream => stream.url === url);
      
      if (existingIndex >= 0) {
        // Update the existing stream with the origin
        if (origin) {
          detectedStreams[existingIndex].origin = origin;
          saveDetectedStreams();
        }
      } else {
        // Create a new stream entry
        const newStream = {
          url: url,
          origin: origin,
          detected: new Date().toISOString(),
          id: Date.now().toString()
        };
        
        console.log("Adding new stream:", newStream);
        
        // Add to detected streams
        detectedStreams.push(newStream);
        saveDetectedStreams();
        updateBadge();
      }
    }
  }
  
  return { requestHeaders: details.requestHeaders };
}

// Save streams to storage
function saveDetectedStreams() {
  browser.storage.local.set({ streams: detectedStreams });
}

// Load streams from storage
function loadDetectedStreams() {
  browser.storage.local.get("streams").then(result => {
    if (result.streams) {
      detectedStreams = result.streams;
      updateBadge();
    }
  });
}

// Save server config to storage
function saveServerConfig() {
  browser.storage.local.set({ serverConfig: serverConfig });
}

// Load server config from storage
function loadServerConfig() {
  browser.storage.local.get("serverConfig").then(result => {
    if (result.serverConfig) {
      serverConfig = result.serverConfig;
      // Test connection on load
      testServerConnection();
    }
  });
}

// Test connection to server
function testServerConnection() {
  return fetch(`${serverConfig.url}/api/status`, {
    headers: {
      'Cache-Control': 'no-cache'
    }
  })
  .then(response => {
    if (!response.ok) {
      throw new Error(`Server returned ${response.status}: ${response.statusText}`);
    }
    return response.json();
  })
  .then(data => {
    console.log("Server status:", data);
    serverConfig.isConnected = true;
    saveServerConfig();
    return { success: true, status: data };
  })
  .catch(error => {
    console.error("Server connection test failed:", error);
    serverConfig.isConnected = false;
    saveServerConfig();
    return { success: false, error: error.message };
  });
}

// Update badge with count of streams
function updateBadge() {
  browser.browserAction.setBadgeText({
    text: detectedStreams.length > 0 ? detectedStreams.length.toString() : ""
  });
  browser.browserAction.setBadgeBackgroundColor({ color: "#4688F1" });
}

// API functions
function addStreamToServer(name, origin, url) {
  return fetch(`${serverConfig.url}/api/stream/start`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json"
    },
    body: JSON.stringify({ name, origin, url })
  })
  .then(response => response.json());
}

function updateStreamOnServer(name, origin, url) {
  // Stop the existing stream first
  return fetch(`${serverConfig.url}/api/stream/stop/${name}`, {
    method: "POST"
  })
  .then(response => {
    if (!response.ok) {
      throw new Error(`Failed to stop stream: ${response.status} ${response.statusText}`);
    }
    return response.json();
  })
  .then(() => {
    // Small delay to ensure cleanup is complete
    return new Promise(resolve => setTimeout(resolve, 500));
  })
  .then(() => {
    // Then start a new one with the same name
    return addStreamToServer(name, origin, url);
  });
}

function removeStreamFromServer(name) {
  return fetch(`${serverConfig.url}/api/stream/stop/${name}`, {
    method: "POST",
    headers: {
      'Cache-Control': 'no-cache'
    }
  })
  .then(response => response.json());
}

function getStreamsFromServer() {
  // Add a cache-busting parameter to ensure we get fresh data
  return fetch(`${serverConfig.url}/api/streams?nocache=${Date.now()}`, {
    headers: {
      'Cache-Control': 'no-cache',
      'Pragma': 'no-cache'
    }
  })
  .then(response => {
    if (!response.ok) {
      throw new Error(`Server returned ${response.status}: ${response.statusText}`);
    }
    return response.json();
  });
}

// Initialize
loadDetectedStreams();
loadServerConfig();

// Listen for messages from popup
browser.runtime.onMessage.addListener((message, sender, sendResponse) => {
  if (message.action === "getDetectedStreams") {
    sendResponse({ streams: detectedStreams });
  }
  else if (message.action === "clearDetectedStreams") {
    detectedStreams = [];
    saveDetectedStreams();
    updateBadge();
    sendResponse({ success: true });
  }
  else if (message.action === "addStream") {
    addStreamToServer(message.name, message.origin, message.url)
      .then(result => sendResponse(result))
      .catch(error => sendResponse({ error: error.toString() }));
    return true; // Required for async sendResponse
  }
  else if (message.action === "updateStream") {
    updateStreamOnServer(message.name, message.origin, message.url)
      .then(result => sendResponse(result))
      .catch(error => sendResponse({ error: error.toString() }));
    return true; // Required for async sendResponse
  }
  else if (message.action === "removeStream") {
    removeStreamFromServer(message.name)
      .then(result => sendResponse(result))
      .catch(error => sendResponse({ error: error.toString() }));
    return true; // Required for async sendResponse
  }
  else if (message.action === "getServerStreams") {
    getStreamsFromServer()
      .then(result => sendResponse(result))
      .catch(error => sendResponse({ error: error.toString() }));
    return true; // Required for async sendResponse
  }
  else if (message.action === "getServerConfig") {
    sendResponse({ serverConfig });
  }
  else if (message.action === "updateServerConfig") {
    serverConfig.url = message.url;
    saveServerConfig();
    sendResponse({ success: true });
  }
  else if (message.action === "testServerConnection") {
    testServerConnection()
      .then(result => sendResponse(result))
      .catch(error => sendResponse({ error: error.toString() }));
    return true; // Required for async sendResponse
  }
});