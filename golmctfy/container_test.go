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

func TestContainerPause(t *testing.T) {
	err := testNormalCases("lmctfy_container_pause", func() error {
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
		return c.Pause()
	}, "lmctfy_new_container_api", "lmctfy_container_api_get_container")
	if err != nil {
		t.Error(err)
	}
}

func TestContainerResume(t *testing.T) {
	err := testNormalCases("lmctfy_container_resume", func() error {
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
		return c.Resume()
	}, "lmctfy_new_container_api", "lmctfy_container_api_get_container")
	if err != nil {
		t.Error(err)
	}
}

func TestContainerKillAll(t *testing.T) {
	err := testNormalCases("lmctfy_container_killall", func() error {
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
		return c.KillAll()
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

func TestContainerListThreads(t *testing.T) {
	err := testNormalCases("lmctfy_container_list_threads", func() error {
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
		threads, err := c.ListThreads(CONTAINER_LIST_POLICY_SELF)
		if err != nil {
			return err
		}
		// XXX predefined
		if len(threads) != 2 {
			t.Errorf("Should be 2 threads, but received %v", len(threads))
			return nil
		}
		if threads[0] != 1 || threads[1] != 2 {
			t.Errorf("Some thread ID is wrong: %+v", threads)
		}
		return nil
	}, "lmctfy_new_container_api", "lmctfy_container_api_get_container")
	if err != nil {
		t.Error(err)
	}
}

func TestContainerListProcesses(t *testing.T) {
	err := testNormalCases("lmctfy_container_list_processes", func() error {
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
		processes, err := c.ListProcesses(CONTAINER_LIST_POLICY_SELF)
		if err != nil {
			return err
		}
		// XXX predefined
		if len(processes) != 3 {
			t.Errorf("Should be 3 processes, but received %v", len(processes))
			return nil
		}
		if processes[2] != 1 || processes[1] != 2 || processes[0] != 3 {
			t.Errorf("Some process ID is wrong: %+v", processes)
		}
		return nil
	}, "lmctfy_new_container_api", "lmctfy_container_api_get_container")
	if err != nil {
		t.Error(err)
	}
}

func TestContainerListSubcontainers(t *testing.T) {
	err := testNormalCases("lmctfy_container_list_subcontainers", func() error {
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
		subs, err := c.ListSubcontainers(CONTAINER_LIST_POLICY_SELF)
		if err != nil {
			return err
		}
		defer func() {
			for _, c := range subs {
				c.Close()
			}
		}()
		// XXX predefined
		if len(subs) != 2 {
			t.Errorf("Should be 2 subcontainers, but received %v", len(subs))
			return nil
		}
		if subs[0].Name() != "/a" {
			t.Errorf("First container should be /a. But got %v", subs[0].Name())
		}
		if subs[1].Name() != "/b" {
			t.Errorf("Second container should be /b. But got %v", subs[1].Name())
		}
		return nil
	}, "lmctfy_new_container_api", "lmctfy_container_api_get_container")
	if err != nil {
		t.Error(err)
	}
}

func TestContainerStats(t *testing.T) {
	err := testNormalCases("lmctfy_container_stats_raw", func() error {
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
		stats, err := c.Stats(CONTAINER_STATS_TYPE_SUMMARY)
		if err != nil {
			return err
		}
		if stats == nil {
			t.Errorf("stats should not be nil")
		}
		return nil
	}, "lmctfy_new_container_api", "lmctfy_container_api_get_container")
	if err != nil {
		t.Error(err)
	}
}
