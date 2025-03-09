# app.py - Main Flask Application
from flask import Flask, request, jsonify, send_from_directory, make_response
import os
import subprocess
import threading
import time
import json
import re

app = Flask(__name__)

# Add no-cache headers to all responses
@app.after_request
def add_header(response):
    response.headers['Cache-Control'] = 'no-store, no-cache, must-revalidate, max-age=0'
    response.headers['Pragma'] = 'no-cache'
    response.headers['Expires'] = '0'
    return response

# Constants
STREAMS_DIR = '/var/www/streams'
STREAMS_FILE = '/var/www/streams/streams.json'
active_streams = {}  # Dictionary to track running stream processes

# Create streams directory if it doesn't exist
try:
    os.makedirs(STREAMS_DIR, exist_ok=True)
    print(f"Stream directory exists at: {STREAMS_DIR}")
    print(f"Contents of directory: {os.listdir(STREAMS_DIR)}")
except Exception as e:
    print(f"Error with streams directory: {str(e)}")

# Validate input strings
def validate_input(text, input_type):
    if input_type == 'name':
        # Allow alphanumeric, dashes, and underscores
        return bool(re.match(r'^[a-zA-Z0-9_-]+$', text))
    elif input_type == 'origin':
        # Basic URL validation
        return bool(re.match(r'^https?://[\w.-]+(?:\.[\w.-]+)+[\w\-._~:/?#[\]@!$&\'()*+,;=]*$', text))
    return False

# Stream handler function
def handle_stream(name, origin, url):
    stream_dir = os.path.join(STREAMS_DIR, name)
    os.makedirs(stream_dir, exist_ok=True)
    
    # Update stream status
    active_streams[name]['status'] = 'active'
    
    try:
        # Prepare ffmpeg command
        cmd = [
            'ffmpeg', 
            '-headers', f'Referer: {origin}/\r\nOrigin: {origin}\r\n',
            '-i', url,
            '-map', '0:1',  # Assuming video is stream 1
            '-map', '0:2',  # Assuming audio is stream 2
            '-c', 'copy',
            '-f', 'segment',
            '-segment_time', '15',
            '-segment_format', 'mpegts',
            '-segment_list', f'{stream_dir}/index.m3u8',
            '-segment_list_flags', '+live',
            f'{stream_dir}/segment_%d.ts'
        ]
        
        # Start ffmpeg process
        process = subprocess.Popen(cmd)
        
        # Store process in active_streams
        active_streams[name]['process'] = process
        process.wait()
        
    except Exception as e:
        active_streams[name]['status'] = 'error'
        active_streams[name]['error'] = str(e)

# API Endpoints
@app.route('/api/stream/start', methods=['POST'])
def start_stream():
    data = request.json
    
    # Print the received data to stdout
    print("=== RECEIVED STREAM START REQUEST ===")
    print(f"Request data: {data}")
    print("=== END REQUEST DATA ===")
    
    # Validate input
    if not data or 'name' not in data or 'origin' not in data or 'url' not in data:
        return jsonify({'error': 'Missing required parameters'}), 400
    
    name = data['name']
    origin = data['origin']
    url = data['url']
    
    print(f"Stream name: {name}")
    print(f"Origin: {origin}")
    print(f"URL: {url}")
    
    # Validate name and origin
    if not validate_input(name, 'name'):
        return jsonify({'error': 'Invalid name format'}), 400
    
    if not validate_input(origin, 'origin'):
        return jsonify({'error': 'Invalid origin format'}), 400
    
    # Check if stream with this name already exists
    if name in active_streams and active_streams[name].get('status') == 'active':
        return jsonify({'error': f'Stream with name {name} already exists'}), 409
    
    # Create stream entry
    active_streams[name] = {
        'name': name,
        'origin': origin,
        'url': url,
        'status': 'initializing',
        'created_at': time.time()
    }
    
    # Start stream in a separate thread
    stream_thread = threading.Thread(
        target=handle_stream,
        args=(name, origin, url)
    )
    stream_thread.daemon = True
    stream_thread.start()
    
    return jsonify({
        'status': 'success',
        'message': f'Stream {name} started',
        'stream_url': f'/stream/{name}/index.m3u8'
    }), 201

@app.route('/api/stream/stop/<name>', methods=['POST'])
def stop_stream(name):
    if name not in active_streams:
        return jsonify({'error': 'Stream not found'}), 404
    
    try:
        # If there's a process, terminate it
        if 'process' in active_streams[name] and active_streams[name]['process']:
            try:
                active_streams[name]['process'].terminate()
            except Exception as e:
                print(f"Error terminating process: {str(e)}")
        
        # Clean up stream files
        stream_dir = os.path.join(STREAMS_DIR, name)
        if os.path.exists(stream_dir):
            try:
                # Remove all files in the directory
                for file in os.listdir(stream_dir):
                    file_path = os.path.join(stream_dir, file)
                    if os.path.isfile(file_path):
                        os.unlink(file_path)
                
                # Remove the directory itself
                os.rmdir(stream_dir)
                print(f"Removed stream directory: {stream_dir}")
            except Exception as e:
                print(f"Error removing stream files: {str(e)}")
        
        # This is the critical part - make sure we remove the stream from the dictionary
        if name in active_streams:
            del active_streams[name]
            print(f"Removed stream '{name}' from active streams")
        
        return jsonify({'status': 'success', 'message': f'Stream {name} stopped and removed'}), 200
    except Exception as e:
        print(f"Error in stop_stream: {str(e)}")
        return jsonify({'error': f'Failed to stop stream: {str(e)}'}), 500

@app.route('/api/streams', methods=['GET'])
def list_streams():
    # First, clean up any streams that might have processes that are no longer running
    streams_to_remove = []
    for name, data in active_streams.items():
        if 'process' in data and data['process']:
            # Check if process is still running
            if data['process'].poll() is not None:  # Process has terminated
                streams_to_remove.append(name)
    
    # Remove dead streams
    for name in streams_to_remove:
        if name in active_streams:
            del active_streams[name]
    
    # Now return the clean list
    streams_info = {}
    for name, data in active_streams.items():
        # Copy only shareable data (no process objects)
        streams_info[name] = {
            'name': data['name'],
            'origin': data['origin'],
            'status': data['status'],
            'created_at': data.get('created_at'),
            'stream_url': f'/stream/{name}/index.m3u8'
        }
    return jsonify(streams_info)

# Serve stream files
@app.route('/stream/<name>/<path:filename>')
def serve_stream_file(name, filename):
    stream_dir = os.path.join(STREAMS_DIR, name)
    return send_from_directory(stream_dir, filename)

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=8080)