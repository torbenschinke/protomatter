package vulkan

import (
	"fmt"
	"runtime"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/torbenschinke/protomatter/logging"
	vk "github.com/vulkan-go/vulkan"
)

type Application struct {
	logger         logging.Logger
	width, height  int
	window         *glfw.Window
	name           string
	instance       vk.Instance
	physicalDevice vk.PhysicalDevice
}

func NewApplication(logger logging.Logger, name string, width, height int) *Application {
	runtime.LockOSThread()
	a := &Application{
		logger: logger,
		width:  width,
		height: height,
		name:   name,
	}
	a.initWindow()
	a.initVulkan()

	return a
}

// initWindow allocates and prepares the Vulkan window using glfw without requiring
// an actual OpenGL context. The window field is valid after calling.
// Panics if no window can be created.
func (a *Application) initWindow() {
	// initialize the GLFW library
	must(glfw.Init())
	vk.SetGetInstanceProcAddr(glfw.GetVulkanGetInstanceProcAddress()) //??
	vk.Init()

	// tell GLFW that we don't need OpenGL
	glfw.WindowHint(glfw.ClientAPI, glfw.NoAPI)

	// TODO: fix me later, needs special care
	glfw.WindowHint(glfw.Resizable, glfw.False)

	// first nil param the actual display to open the window, leave default,
	// last nil is for OpenGL only
	window, err := glfw.CreateWindow(a.width, a.height, a.name, nil, nil)
	must(err)

	a.window = window

	a.logger.Println("glfw init done")
}

func (a *Application) initVulkan() {
	a.createInstance()
	// TODO debug layer
	a.pickPhysicalDevice()
}

func (a *Application) createInstance() {

	// collect information about our application.
	// this may be used by the driver to optimize things
	// for our application
	appInfo := &vk.ApplicationInfo{
		SType:              vk.StructureTypeApplicationInfo,
		PApplicationName:   a.name,
		ApplicationVersion: vk.MakeVersion(1, 0, 0),
		PEngineName:        "protomatter",
		ApiVersion:         vk.ApiVersion10, // TODO switch to 1.1 when raspberry 4 supports it
		EngineVersion:      vk.MakeVersion(1, 0, 0),
	}

	// configure the creation parameters to configure global extensions and
	// validation layers
	extensions := glfw.GetCurrentContext().GetRequiredInstanceExtensions()
	a.logger.Println("vulkan extensions:")
	for _, e := range extensions {
		a.logger.Println(e)
	}

	instInfo := &vk.InstanceCreateInfo{
		SType:                   vk.StructureTypeInstanceCreateInfo,
		PApplicationInfo:        appInfo,
		EnabledExtensionCount:   uint32(len(extensions)),
		PpEnabledExtensionNames: extensions,
	}

	// TODO see https://github.com/vulkan-go/vulkan/issues/56
	// technically correct: allocating the target pointer space in the go-heap
	// may remove that any time, but the stackframe is at least alive
	// until this method returns.
	var v vk.Instance
	must(vk.Error(vk.CreateInstance(instInfo, nil, &v)))
	a.instance = v

	a.logger.Println("vk instance created")

	// inspect the extension properties
	var instExtCount uint32
	must(vk.Error(vk.EnumerateInstanceExtensionProperties("", &instExtCount, nil)))
	a.logger.Println("extension properties", instExtCount)
	instExtProps := make([]vk.ExtensionProperties, instExtCount)
	must(vk.Error(vk.EnumerateInstanceExtensionProperties("", &instExtCount, instExtProps)))
	for _, ext := range instExtProps {
		ext.Deref()
		a.logger.Println(vk.ToString((ext.ExtensionName[:])))
	}

}

// picks a fitting device. We actually may render to multiple, but our engine just picks one.
func (a *Application) pickPhysicalDevice() {
	var deviceCount uint32
	must(vk.Error(vk.EnumeratePhysicalDevices(a.instance, &deviceCount, nil)))

	if deviceCount == 0 {
		panic("no vulkan gpu found")
	}

	physDevs := make([]vk.PhysicalDevice, deviceCount)
	must(vk.Error(vk.EnumeratePhysicalDevices(a.instance, &deviceCount, physDevs)))

	for _, dev := range physDevs {
		if a.isDeviceSuitable(dev) {
			a.physicalDevice = dev
			break
		}
	}
}

// isDeviceSuitable should pick the most capable and performant device OR
// perhaps taking a user-based configuration into account.
func (a *Application) isDeviceSuitable(device vk.PhysicalDevice) bool {
	var deviceProperties vk.PhysicalDeviceProperties
	vk.GetPhysicalDeviceProperties(device, &deviceProperties)
	deviceProperties.Deref()
	deviceProperties.Limits.Deref()

	var deviceFeatures vk.PhysicalDeviceFeatures
	vk.GetPhysicalDeviceFeatures(device, &deviceFeatures)
	deviceFeatures.Deref()

	a.logger.Println(vk.ToString(deviceProperties.DeviceName[:]))
	a.logger.Println("max img dim 2d", deviceProperties.Limits.MaxImageDimension2D)
	a.logger.Println(fmt.Sprintf("%+v", deviceFeatures))

	return true
}

// Run must be called from the outside to block and process events.
func (a *Application) Run() {
	a.logger.Println("entering main loop")
	for !a.window.ShouldClose() {
		glfw.PollEvents()
	}

	a.logger.Println("main loop done")
}

// Close releases any glfw resources.
func (a *Application) Close() error {
	vk.DestroyInstance(a.instance, nil)
	a.window.Destroy()
	glfw.Terminate()

	a.logger.Println("glfw closed")
	return nil
}
