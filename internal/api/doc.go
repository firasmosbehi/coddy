// Package api provides the HTTP API for Coddy.
//
// API Endpoints:
//
// Health & Status:
//   GET /health                    - Health check
//   GET /stats                     - Server statistics
//
// Sessions:
//   POST /sessions                 - Create new session
//   GET  /sessions                 - List all sessions
//   GET  /sessions/:id             - Get session details
//   DEL  /sessions/:id             - Delete session
//
// File Operations:
//   POST /sessions/:id/upload      - Upload file (multipart/form-data)
//   GET  /sessions/:id/files       - List files in session
//   GET  /sessions/:id/files/:path - Download file
//
// Real-time Communication:
//   WS   /ws/sessions/:id          - WebSocket for streaming chat
//
// Example Usage:
//
// Create a session:
//   curl -X POST http://localhost:8000/sessions
//
// Upload a file:
//   curl -X POST -F "file=@data.csv" http://localhost:8000/sessions/:id/upload
//
// List files:
//   curl http://localhost:8000/sessions/:id/files
//
// Download a file:
//   curl http://localhost:8000/sessions/:id/files/data.csv -o data.csv
//
// WebSocket connection (JavaScript):
//   const ws = new WebSocket('ws://localhost:8000/ws/sessions/:id');
//   ws.onopen = () => {
//     ws.send(JSON.stringify({
//       type: "chat",
//       payload: JSON.stringify({message: "Hello!"})
//     }));
//   };
//   ws.onmessage = (event) => {
//     const msg = JSON.parse(event.data);
//     console.log(msg.content);
//   };
//
package api
