package tests

// All timeouts are in seconds
const (
	CONTAINER_NAME       = "container_name"
	PORT                 = "port"
	DEFAULT_TIMEOUT      = "default_timeout"
	FIRST_CONTAINER_WAIT = "first_container_wait"
)

var (
	StepOptions map[string]interface{} = map[string]interface{}{
		CONTAINER_NAME:       "servicecontrol-server_control",
		PORT:                 10080,
		DEFAULT_TIMEOUT:      int64(10),
		FIRST_CONTAINER_WAIT: int64(180),
	}
)
