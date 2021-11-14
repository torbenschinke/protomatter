# protomatter
an experimental vulkan 3d engine for linux (raspberry 4).
Implementation based on https://vulkan-tutorial.com/.

## develop on raspberry 4 (Raspberry OS _bullseye_)

```bash
sudo apt install mesa-vulkan-drivers
sudo apt install libvulkan-dev libvulkan1 vulkan-tools
```

## grabbing proper wrappers

```bash
go get -u -tags=vulkan github.com/go-gl/glfw/v3.3/glfw
go get github.com/vulkan-go/vulkan
```