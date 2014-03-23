package golmctfy

import (
	. "containers_lmctfy"
	"fmt"
	"testing"
	"time"
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

func TestNotification(t *testing.T) {
	err := testNormalCases("lmctfy_container_register_notification_raw", func() error {
		api, err := NewContainerApi()
		defer api.Close()
		if err != nil {
			return fmt.Errorf("This should not fail: %v", err)
		}
		containerName := "/container"
		c, err := api.Get(containerName)
		defer c.Close()
		if err != nil {
			return fmt.Errorf("Get() should not fail: %v", err)
		}
		ch := make(chan *Event)
		var spec EventSpec
		notifId, err := c.RegisterNotification(&spec, ch)
		if err != nil {
			return err
		}
		testNotify := func(nid uint64, evtErr error) {
			notifyContainer(c, evtErr)
			select {
			case evt := <-ch:
				if evt.Container != c {
					t.Errorf("should return the same container")
				}
				if evt.NotifId != nid {
					t.Errorf("returned notif id is %v; should be %v", evt.NotifId, nid)
				}
				if evtErr == nil {
					if evt.Error != nil {
						t.Errorf("received notification has message %v; should be nil", evt.Error)
					}
				} else {
					if evt.Error.Error() != evtErr.Error() {
						t.Errorf("received notification has message %v; should be %v", evt.Error, evtErr)
					}
				}
			case <-time.After(3 * time.Second):
				t.Errorf("Timeout")
			}
		}
		testNotify(notifId, nil)
		testNotify(notifId, fmt.Errorf("some event"))
		return expectCall("lmctfy_container_unregister_notification", 0, "", func() error {
			return c.UnregisterNotification(notifId)
		})
		return nil
	}, "lmctfy_new_container_api", "lmctfy_container_api_get_container")
	if err != nil {
		t.Error(err)
	}
}
