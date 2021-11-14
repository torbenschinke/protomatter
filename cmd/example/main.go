package main

import (
	"github.com/torbenschinke/protomatter/logging"
	"github.com/torbenschinke/protomatter/vulkan"
)

func main() {
	app := vulkan.NewApplication(logging.NewLogger(), "protomatter engine example", 800, 600)
	app.Run()
	defer app.Close()
}
