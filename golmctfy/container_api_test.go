package golmctfy

import (
	. "containers_lmctfy"
	"fmt"
	"testing"
)

func testNormalCases(fn string, f func() error, preamble ...string) error {
	defer resetMockEnv()
	err := expectCall(fn, 0, "", f, preamble...)
	if err != nil {
		return err
	}
	err = assertExpectations()
	if err != nil {
		return err
	}
	err = expectCall(fn, 2, "error message", f, preamble...)
	if err != nil {
		return err
	}
	err = assertExpectations()
	if err != nil {
		return err
	}
	return nil
}

func TestCMock(t *testing.T) {
	defer resetMockEnv()
	err := testNormalCases("some_function_which_should_never_be_called", func() error {
		var spec InitSpec
		return InitMachine(&spec)
	})
	if err == nil {
		t.Error("Should return error")
	}
	err = expectCall("should_be_called", 0, "", func() error {
		return nil
	})
	err = assertExpectations()
	if err == nil {
		t.Error("Should return error")
	}
}

func TestInitMachine(t *testing.T) {
	var spec InitSpec
	err := testNormalCases("lmctfy_init_machine_raw", func() error {
		return InitMachine(&spec)
	})
	if err != nil {
		t.Error(err)
	}
}

func TestNewContainerApi(t *testing.T) {
	err := testNormalCases("lmctfy_new_container_api", func() error {
		api, err := NewContainerApi()
		defer api.Close()
		return err
	})
	if err != nil {
		t.Error(err)
	}
}

func TestCreateContainer(t *testing.T) {
	err := testNormalCases("lmctfy_container_api_create_container_raw", func() error {
		api, err := NewContainerApi()
		defer api.Close()
		if err != nil {
			return fmt.Errorf("This should not fail: %v", err)
		}
		var spec ContainerSpec
		containerName := "/container"
		c, err := api.Create(containerName, &spec)
		if err != nil {
			return err
		}
		defer c.Close()
		cn := c.Name()
		if containerName != cn {
			return fmt.Errorf("Container name should be %v, but received %v", containerName, cn)
		}

		return nil
	}, "lmctfy_new_container_api")
	if err != nil {
		t.Error(err)
	}

}

func TestGetContainer(t *testing.T) {
	err := testNormalCases("lmctfy_container_api_get_container", func() error {
		api, err := NewContainerApi()
		defer api.Close()
		if err != nil {
			return fmt.Errorf("This should not fail: %v", err)
		}
		containerName := "/container"
		c, err := api.Get(containerName)
		if err != nil {
			return err
		}
		defer c.Close()
		cn := c.Name()
		if containerName != cn {
			return fmt.Errorf("Container name should be %v, but received %v", containerName, cn)
		}

		return nil
	}, "lmctfy_new_container_api")
	if err != nil {
		t.Error(err)
	}
}

func TestDestroyContainer(t *testing.T) {
	err := testNormalCases("lmctfy_container_api_destroy_container", func() error {
		api, err := NewContainerApi()
		defer api.Close()
		if err != nil {
			return fmt.Errorf("This should not fail: %v", err)
		}
		containerName := "/container"
		c, err := api.Get(containerName)
		if err != nil {
			return err
		}
		defer c.Close()

		return api.Destroy(c)
	}, "lmctfy_new_container_api", "lmctfy_container_api_get_container")
	if err != nil {
		t.Error(err)
	}
}

func TestDetectContainer(t *testing.T) {
	err := testNormalCases("lmctfy_container_api_detect_container", func() error {
		api, err := NewContainerApi()
		defer api.Close()
		if err != nil {
			return fmt.Errorf("This should not fail: %v", err)
		}
		c, err := api.Detect(0)
		if err != nil {
			return err
		}
		if len(c) == 0 {
			return fmt.Errorf("Should return a container's name")
		}
		return nil
	}, "lmctfy_new_container_api")
	if err != nil {
		t.Error(err)
	}
}
