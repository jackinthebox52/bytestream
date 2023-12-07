# Project Name

## Description

This project is currently under development and serves as a proof of concept for maintaining and rebroadcasting multiple HLS streams on a network. It includes a web server that serves both web assets (HTML, CSS, JS) and a REST API, as well as a Go backend.

## Features

 - Web Server: The web server serves both web assets and a REST API. It uses the Gin Web Framework to handle HTTP requests and responses.
 - Go Backend: The Go backend is responsible for handling HLS streams. It uses the ffmpeg library to ingest streams and rebroadcast them.
 - Stream Viewing: The web application allows users to view the HLS streams. It uses the hls.js library to play the streams in the browser.

## Installation

This section will contain instructions on how to install and run the project.

## Usage

The project runs a web server on port 8081. You can access the web application by navigating to http://localhost:8081 in your web browser. The main page displays a list of available streams. You can view a stream by clicking on its name, which will redirect you to the player page.

To add a new stream, you can use the POST /bstreams endpoint with the stream URL and referrer as parameters.

## Development

Contributions to the project are welcome. To contribute, fork the repository, make your changes, and submit a pull request. Please ensure that your code follows the Go coding standards and includes appropriate unit tests.

## Roadmap

Future improvements for the project include adding authentication to the API, improving the web interface, and optimizing the stream handling for better performance.

## License

This project is licensed under the MIT License. Please see the LICENSE file for more details.