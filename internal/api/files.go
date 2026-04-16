// Package api provides file handling endpoints.
package api

import (
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/firasmosbehi/coddy/internal/session"
)

// UploadFile handles file uploads to a session's sandbox.
func (h *Handlers) UploadFile(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Path[len("/sessions/"):]
	sessionID = strings.TrimSuffix(sessionID, "/upload")
	sessionID = strings.TrimSuffix(sessionID, "/")

	if sessionID == "" {
		respondError(w, http.StatusBadRequest, "Session ID required")
		return
	}

	// Parse multipart form
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10 MB max
		respondError(w, http.StatusBadRequest, "Failed to parse form: "+err.Error())
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		respondError(w, http.StatusBadRequest, "Failed to get file: "+err.Error())
		return
	}
	defer file.Close()

	// Get destination path (optional, defaults to filename)
	dstPath := r.FormValue("path")
	if dstPath == "" {
		dstPath = header.Filename
	}

	// Validate path (prevent directory traversal)
	dstPath = filepath.Clean(dstPath)
	if strings.Contains(dstPath, "..") {
		respondError(w, http.StatusBadRequest, "Invalid path")
		return
	}

	// Get session with sandbox
	ctx := r.Context()
	sessWithSandbox, err := h.sessions.Store().GetWithSandbox(ctx, sessionID)
	if err != nil {
		if err == session.ErrSessionNotFound {
			respondError(w, http.StatusNotFound, "Session not found")
			return
		}
		if err == session.ErrSessionExpired {
			respondError(w, http.StatusGone, "Session expired")
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Create temp file
	tmpFile, err := io.ReadAll(file)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to read file: "+err.Error())
		return
	}

	// Write to temp location
	tmpPath := "/tmp/coddy_upload_" + header.Filename
	if err := writeTempFile(tmpPath, tmpFile); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to write temp file: "+err.Error())
		return
	}

	// Upload to sandbox
	if err := sessWithSandbox.Sandbox.UploadFile(ctx, tmpPath, dstPath); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to upload to sandbox: "+err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"message":  "File uploaded successfully",
		"filename": header.Filename,
		"path":     dstPath,
		"size":     len(tmpFile),
	})
}

// DownloadFile handles file downloads from a session's sandbox.
func (h *Handlers) DownloadFile(w http.ResponseWriter, r *http.Request) {
	// Extract session ID and file path from URL
	// Format: /sessions/:id/files/:path
	path := r.URL.Path[len("/sessions/"):]
	parts := strings.SplitN(path, "/files/", 2)

	if len(parts) != 2 {
		respondError(w, http.StatusBadRequest, "Invalid path format")
		return
	}

	sessionID := parts[0]
	filePath := parts[1]

	if sessionID == "" || filePath == "" {
		respondError(w, http.StatusBadRequest, "Session ID and file path required")
		return
	}

	// Validate path
	filePath = filepath.Clean(filePath)
	if strings.Contains(filePath, "..") {
		respondError(w, http.StatusBadRequest, "Invalid path")
		return
	}

	// Get session with sandbox
	ctx := r.Context()
	sessWithSandbox, err := h.sessions.Store().GetWithSandbox(ctx, sessionID)
	if err != nil {
		if err == session.ErrSessionNotFound {
			respondError(w, http.StatusNotFound, "Session not found")
			return
		}
		if err == session.ErrSessionExpired {
			respondError(w, http.StatusGone, "Session expired")
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Download file from sandbox
	data, err := sessWithSandbox.Sandbox.DownloadFile(ctx, filePath)
	if err != nil {
		respondError(w, http.StatusNotFound, "File not found: "+err.Error())
		return
	}

	// Set headers
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filepath.Base(filePath)))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))

	w.Write(data)
}

// ListFiles lists files in a session's sandbox.
func (h *Handlers) ListFiles(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Path[len("/sessions/"):]
	sessionID = strings.TrimSuffix(sessionID, "/files")
	sessionID = strings.TrimSuffix(sessionID, "/")

	if sessionID == "" {
		respondError(w, http.StatusBadRequest, "Session ID required")
		return
	}

	// Get path parameter (optional)
	path := r.URL.Query().Get("path")
	if path == "" {
		path = "."
	}

	// Get session with sandbox
	ctx := r.Context()
	sessWithSandbox, err := h.sessions.Store().GetWithSandbox(ctx, sessionID)
	if err != nil {
		if err == session.ErrSessionNotFound {
			respondError(w, http.StatusNotFound, "Session not found")
			return
		}
		if err == session.ErrSessionExpired {
			respondError(w, http.StatusGone, "Session expired")
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// List files
	files, err := sessWithSandbox.Sandbox.ListFiles(ctx, path)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to list files: "+err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"path":  path,
		"files": files,
		"count": len(files),
	})
}

func writeTempFile(path string, data []byte) error {
	// This is a helper to write temp files for upload
	// In production, use proper temp file handling
	return nil // Placeholder
}
