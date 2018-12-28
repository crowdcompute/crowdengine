package rpc

// API describes the rpc procedure behaviour
type API struct {
	Namespace string
	Version   string
	Service   interface{}
	Public    bool
}

// Args represent the arguments
type Args struct {
	S string
}

// Result represents the rpc result
type Result struct {
	Output string
}
