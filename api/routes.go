package api

import (
	"camera-control/stream"
	"encoding/json"
	"net/http"
)

type StartRequest struct {
	CameraID string `json:"camera_id"`
	RTSP     string `json:"rtsp_url"`
	N        int    `json:"every_nth"`
}

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/start", startStreamHandler)
	mux.HandleFunc("/stop", stopStreamHandler)
}

func startStreamHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var req StartRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	go stream.StartStream(req.CameraID, req.RTSP, req.N)
	w.WriteHeader(http.StatusAccepted)
}

func stopStreamHandler(w http.ResponseWriter, r *http.Request) {
	cameraID := r.URL.Query().Get("camera_id")
	if cameraID == "" {
		http.Error(w, "Missing camera_id", http.StatusBadRequest)
		return
	}

	stream.StopStream(cameraID)
	w.WriteHeader(http.StatusOK)
}
