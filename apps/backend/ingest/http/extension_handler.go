package http

import (
	"encoding/json"
	"fmt"
	ias_extension "ias/automation/extension"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

func GetExtensionList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if os.Getenv("IAS_ENABLE_EXTENSION") != "true" {
		http.Error(w, `{"error":"extensions are disabled (IAS_ENABLE_EXTENSION != true)"}`, http.StatusServiceUnavailable)
		return
	}
	slog.Info("Listing loaded extensions", "process", "extension_handler")
	mgr := ias_extension.GetGlobal()
	if mgr == nil {
		http.Error(w, `{"error":"extension manager not initialized"}`, http.StatusInternalServerError)
		return
	}
	extensions := mgr.List()
	if extensions == nil {
		extensions = []map[string]interface{}{}
	}
	jsonData, err := json.Marshal(extensions)
	if err != nil {
		slog.Error("Failed to marshal extension list", "error", err)
		http.Error(w, `{"error":"failed to marshal extension list"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func ServeExtensionUI(w http.ResponseWriter, r *http.Request) {
	if os.Getenv("IAS_ENABLE_EXTENSION") != "true" {
		http.Error(w, "extensions are disabled", http.StatusServiceUnavailable)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/extensions/")
	slashIdx := strings.Index(path, "/")

	var name string
	var forwardPath string

	if slashIdx == -1 {
		name = path
		forwardPath = "/"
	} else {
		name = path[:slashIdx]
		forwardPath = path[slashIdx:]
	}

	if name == "" {
		http.Error(w, "missing extension name in path", http.StatusBadRequest)
		return
	}

	mgr := ias_extension.GetGlobal()
	if mgr == nil {
		http.Error(w, "extension manager not initialized", http.StatusInternalServerError)
		return
	}

	port, ok := mgr.GetPort(name)
	if !ok {
		http.Error(w, fmt.Sprintf("extension %q not found", name), http.StatusNotFound)
		return
	}

	target, err := url.Parse(fmt.Sprintf("http://localhost:%d", port))
	if err != nil {
		slog.Error("Failed to parse extension target URL", "extension", name, "port", port, "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		slog.Error("Extension UI proxy error", "extension", name, "port", port, "error", err)
		http.Error(w, fmt.Sprintf("extension %q unreachable", name), http.StatusBadGateway)
	}

	r.URL.Path = forwardPath
	slog.Debug("Proxying extension request",
		"extension", name, "port", port, "forward", forwardPath, "process", "extension_handler")
	proxy.ServeHTTP(w, r)
}
