package golmctfy

import (
	. "containers_lmctfy"
	"fmt"
	"testing"
)

func TestContainerEnter(t *testing.T) {
	err := testNormalCases("lmctfy_container_enter", func() error {
		api, err := NewContainerApi()
		defer api.Close()
		if err != nil {
			return fmt.Errorf("This should not fail: %v", err)
		}
		containerName := "/container"
		c, err := api.Get(containerName)
		if err != nil {
			return fmt.Errorf("Get() should not fail: %v", err)
		}
		defer c.Close()
		tids := []int{1, 2, 3}
		return c.Enter(tids)
	}, "lmctfy_new_container_api", "lmctfy_container_api_get_container")
	if err != nil {
		t.Error(err)
	}
}

func TestContainerRun(t *testing.T) {
	err := testNormalCases("lmctfy_container_run_raw", func() error {
		api, err := NewContainerApi()
		defer api.Close()
		if err != nil {
			return fmt.Errorf("This should not fail: %v", err)
		}
		containerName := "/container"
		c, err := api.Get(containerName)
		if err != nil {
			return fmt.Errorf("Get() should not fail: %v", err)
		}
		defer c.Close()
		var spec RunSpec
		args := []string{"/bin/echo", "hello"}
		_, err = c.Run(args, &spec)
		return err
	}, "lmctfy_new_container_api", "lmctfy_container_api_get_container")
	if err != nil {
		t.Error(err)
	}
}

func TestContainerUpdate(t *testing.T) {
	err := testNormalCases("lmctfy_container_update_raw", func() error {
		api, err := NewContainerApi()
		defer api.Close()
		if err != nil {
			return fmt.Errorf("This should not fail: %v", err)
		}
		containerName := "/container"
		c, err := api.Get(containerName)
		if err != nil {
			return fmt.Errorf("Get() should not fail: %v", err)
		}
		defer c.Close()
		var spec ContainerSpec
		return c.Update(CONTAINER_UPDATE_POLICY_DIFF, &spec)
	}, "lmctfy_new_container_api", "lmctfy_container_api_get_container")
	if err != nil {
		t.Error(err)
	}
}
