package vulkan

import (
	"fmt"

	vk "github.com/vulkan-go/vulkan"
)

// must asserts that err is nil
func must(err error) {
	if err != nil {
		panic(err)
	}
}

// mustVk asserts that the res is VK_SUCCESS.
func mustVk(res vk.Result) {
	if res != vk.Success {
		panic(fmt.Errorf("VkResult is %d", res))
	}
}
