package extension

const PortAnnouncePrefix = "IAS_EXTENSION_PORT="

type ExecuteRequest struct {
	Action string                 `json:"action"`
	Params map[string]interface{} `json:"params"`
}

type ExecuteResponse struct {
	Success bool                   `json:"success"`
	Data    map[string]interface{} `json:"data,omitempty"`
	Error   string                 `json:"error,omitempty"`
}
