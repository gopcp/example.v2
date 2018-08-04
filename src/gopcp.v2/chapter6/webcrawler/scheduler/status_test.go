package scheduler

import (
	"sync"
	"testing"
)

func TestCheckStatus(t *testing.T) {
	var currentStatus, wantedStatus Status
	var currentStatusList, wantedStatusList []Status
	// 1.处于正在初始化、正在启动和正在停止状态时，不能有任何的状态改变。
	currentStatusList = []Status{
		SCHED_STATUS_INITIALIZING,
		SCHED_STATUS_STARTING,
		SCHED_STATUS_STOPPING,
	}
	wantedStatus = SCHED_STATUS_INITIALIZING
	for _, currentStatus := range currentStatusList {
		if err := checkStatus(currentStatus, wantedStatus, nil); err == nil {
			t.Fatalf("It still can check status with incorrect current status %q!",
				GetStatusDescription(currentStatus))
		}
	}
	// 2. 想要的状态只能是正在初始化、正在启动和正在停止状态中的一个。
	currentStatus = SCHED_STATUS_UNINITIALIZED
	wantedStatusList = []Status{
		SCHED_STATUS_UNINITIALIZED,
		SCHED_STATUS_INITIALIZED,
		SCHED_STATUS_STARTED,
		SCHED_STATUS_STOPPED,
	}
	for _, wantedStatus := range wantedStatusList {
		if err := checkStatus(currentStatus, wantedStatus, nil); err == nil {
			t.Fatalf("It still can check status with incorrect wanted status %q!",
				GetStatusDescription(wantedStatus))
		}
	}
	// 3. 处于未初始化状态时，不能变为正在启动状态和正在停止状态。
	currentStatus = SCHED_STATUS_UNINITIALIZED
	wantedStatusList = []Status{
		SCHED_STATUS_STARTING,
		SCHED_STATUS_STOPPING,
	}
	for _, wantedStatus := range wantedStatusList {
		if err := checkStatus(currentStatus, wantedStatus, nil); err == nil {
			t.Fatalf("It still can check status with current status %q wanted status %q!",
				GetStatusDescription(currentStatus), GetStatusDescription(wantedStatus))
		}
	}
	wantedStatus = SCHED_STATUS_INITIALIZING
	if err := checkStatus(currentStatus, wantedStatus, nil); err != nil {
		t.Fatalf("An error occurs when checking status: %s (currentStatus: %q, wantedStatus: %q)!",
			err, GetStatusDescription(currentStatus), GetStatusDescription(wantedStatus))
	}
	// 4. 处于已启动状态时，不能变为正在初始化和正在启动状态。
	currentStatus = SCHED_STATUS_STARTED
	wantedStatusList = []Status{
		SCHED_STATUS_INITIALIZING,
		SCHED_STATUS_STARTING,
	}
	for _, wantedStatus := range wantedStatusList {
		if err := checkStatus(currentStatus, wantedStatus, nil); err == nil {
			t.Fatalf("It still can check status with current status %q wanted status %q!",
				GetStatusDescription(currentStatus), GetStatusDescription(wantedStatus))
		}
	}
	wantedStatus = SCHED_STATUS_STOPPING
	if err := checkStatus(currentStatus, wantedStatus, nil); err != nil {
		t.Fatalf("An error occurs when checking status: %s (currentStatus: %q, wantedStatus: %q)!",
			err, GetStatusDescription(currentStatus), GetStatusDescription(wantedStatus))
	}
	// 5. 只要未处于已启动状态就不能变为正在停止状态。
	currentStatusList = []Status{
		SCHED_STATUS_UNINITIALIZED,
		SCHED_STATUS_INITIALIZING,
		SCHED_STATUS_INITIALIZED,
		SCHED_STATUS_STARTING,
		SCHED_STATUS_STOPPING,
		SCHED_STATUS_STOPPED,
	}
	wantedStatus = SCHED_STATUS_STOPPING
	for _, currentStatus := range currentStatusList {
		if err := checkStatus(currentStatus, wantedStatus, nil); err == nil {
			t.Fatalf("It still can check status with current status %q wanted status %q!",
				GetStatusDescription(currentStatus), GetStatusDescription(wantedStatus))
		}
	}
	currentStatus = SCHED_STATUS_STARTED
	if err := checkStatus(currentStatus, wantedStatus, nil); err != nil {
		t.Fatalf("An error occurs when checking status: %s (currentStatus: %q, wantedStatus: %q)!",
			err, GetStatusDescription(currentStatus), GetStatusDescription(wantedStatus))
	}
}

func TestCheckStatusInParallel(t *testing.T) {
	number := 1000
	var lock sync.Mutex
	t.Run("Check status in parallel(1)", func(t *testing.T) {
		for i := 0; i < number; i++ {
			currentStatus := SCHED_STATUS_UNINITIALIZED
			wantedStatus := SCHED_STATUS_INITIALIZING
			if err := checkStatus(currentStatus, wantedStatus, &lock); err != nil {
				t.Fatalf("An error occurs when checking status: %s (currentStatus: %q, wantedStatus: %q)!",
					err, GetStatusDescription(currentStatus), GetStatusDescription(wantedStatus))
			}
		}
	})
	t.Run("Check status in parallel(2)", func(t *testing.T) {
		for i := 0; i < number; i++ {
			currentStatus := SCHED_STATUS_INITIALIZED
			wantedStatusList := []Status{
				SCHED_STATUS_INITIALIZING,
				SCHED_STATUS_STARTING,
			}
			for _, wantedStatus := range wantedStatusList {
				if err := checkStatus(currentStatus, wantedStatus, nil); err != nil {
					t.Fatalf("An error occurs when checking status: %s (currentStatus: %q, wantedStatus: %q)!",
						err, GetStatusDescription(currentStatus), GetStatusDescription(wantedStatus))
				}
			}
		}
	})
	t.Run("Check status in parallel(3)", func(t *testing.T) {
		for i := 0; i < number; i++ {
			currentStatus := SCHED_STATUS_STARTED
			wantedStatus := SCHED_STATUS_STOPPING
			if err := checkStatus(currentStatus, wantedStatus, nil); err != nil {
				t.Fatalf("An error occurs when checking status: %s (currentStatus: %q, wantedStatus: %q)!",
					err, GetStatusDescription(currentStatus), GetStatusDescription(wantedStatus))
			}
		}
	})
	t.Run("Check status in parallel(4)", func(t *testing.T) {
		for i := 0; i < number; i++ {
			currentStatus := SCHED_STATUS_STOPPED
			wantedStatusList := []Status{
				SCHED_STATUS_INITIALIZING,
				SCHED_STATUS_STARTING,
			}
			for _, wantedStatus := range wantedStatusList {
				if err := checkStatus(currentStatus, wantedStatus, nil); err != nil {
					t.Fatalf("An error occurs when checking status: %s (currentStatus: %q, wantedStatus: %q)!",
						err, GetStatusDescription(currentStatus), GetStatusDescription(wantedStatus))
				}
			}
		}
	})
}

func TestGetStatusDescription(t *testing.T) {
	statusMap := map[Status]string{
		SCHED_STATUS_UNINITIALIZED: "uninitialized",
		SCHED_STATUS_INITIALIZING:  "initializing",
		SCHED_STATUS_INITIALIZED:   "initialized",
		SCHED_STATUS_STARTING:      "starting",
		SCHED_STATUS_STARTED:       "started",
		SCHED_STATUS_STOPPING:      "stopping",
		SCHED_STATUS_STOPPED:       "stopped",
		Status(7):                  "unknown",
	}
	for status, expectedDesc := range statusMap {
		desc := GetStatusDescription(status)
		if desc != expectedDesc {
			t.Fatalf("Inconsistent description for status %d: expected: %s, actual: %s",
				status, expectedDesc, desc)
		}
	}
}
