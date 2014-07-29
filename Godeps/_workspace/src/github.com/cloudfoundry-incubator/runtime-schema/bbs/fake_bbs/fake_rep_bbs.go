// This file was generated by counterfeiter
package fake_bbs

import (
	"github.com/cloudfoundry-incubator/runtime-schema/bbs"
	"github.com/cloudfoundry-incubator/runtime-schema/bbs/services_bbs"

	"github.com/cloudfoundry-incubator/runtime-schema/models"

	"sync"
	"time"
)

type FakeRepBBS struct {
	MaintainExecutorPresenceStub        func(heartbeatInterval time.Duration, executorPresence models.ExecutorPresence) (services_bbs.Presence, <-chan bool, error)
	maintainExecutorPresenceMutex       sync.RWMutex
	maintainExecutorPresenceArgsForCall []struct {
		heartbeatInterval time.Duration
		executorPresence  models.ExecutorPresence
	}
	maintainExecutorPresenceReturns struct {
		result1 services_bbs.Presence
		result2 <-chan bool
		result3 error
	}
	WatchForDesiredTaskStub        func() (<-chan models.Task, chan<- bool, <-chan error)
	watchForDesiredTaskMutex       sync.RWMutex
	watchForDesiredTaskArgsForCall []struct{}
	watchForDesiredTaskReturns struct {
		result1 <-chan models.Task
		result2 chan<- bool
		result3 <-chan error
	}
	ClaimTaskStub        func(taskGuid string, executorID string) error
	claimTaskMutex       sync.RWMutex
	claimTaskArgsForCall []struct {
		taskGuid   string
		executorID string
	}
	claimTaskReturns struct {
		result1 error
	}
	StartTaskStub        func(taskGuid string, executorID string, containerHandle string) error
	startTaskMutex       sync.RWMutex
	startTaskArgsForCall []struct {
		taskGuid        string
		executorID      string
		containerHandle string
	}
	startTaskReturns struct {
		result1 error
	}
	CompleteTaskStub        func(taskGuid string, failed bool, failureReason string, result string) error
	completeTaskMutex       sync.RWMutex
	completeTaskArgsForCall []struct {
		taskGuid      string
		failed        bool
		failureReason string
		result        string
	}
	completeTaskReturns struct {
		result1 error
	}
	ReportActualLRPAsStartingStub        func(lrp models.ActualLRP, executorID string) error
	reportActualLRPAsStartingMutex       sync.RWMutex
	reportActualLRPAsStartingArgsForCall []struct {
		lrp        models.ActualLRP
		executorID string
	}
	reportActualLRPAsStartingReturns struct {
		result1 error
	}
	ReportActualLRPAsRunningStub        func(lrp models.ActualLRP, executorId string) error
	reportActualLRPAsRunningMutex       sync.RWMutex
	reportActualLRPAsRunningArgsForCall []struct {
		lrp        models.ActualLRP
		executorId string
	}
	reportActualLRPAsRunningReturns struct {
		result1 error
	}
	RemoveActualLRPStub        func(lrp models.ActualLRP) error
	removeActualLRPMutex       sync.RWMutex
	removeActualLRPArgsForCall []struct {
		lrp models.ActualLRP
	}
	removeActualLRPReturns struct {
		result1 error
	}
	RemoveActualLRPForIndexStub        func(processGuid string, index int, instanceGuid string) error
	removeActualLRPForIndexMutex       sync.RWMutex
	removeActualLRPForIndexArgsForCall []struct {
		processGuid  string
		index        int
		instanceGuid string
	}
	removeActualLRPForIndexReturns struct {
		result1 error
	}
	WatchForStopLRPInstanceStub        func() (<-chan models.StopLRPInstance, chan<- bool, <-chan error)
	watchForStopLRPInstanceMutex       sync.RWMutex
	watchForStopLRPInstanceArgsForCall []struct{}
	watchForStopLRPInstanceReturns struct {
		result1 <-chan models.StopLRPInstance
		result2 chan<- bool
		result3 <-chan error
	}
	ResolveStopLRPInstanceStub        func(stopInstance models.StopLRPInstance) error
	resolveStopLRPInstanceMutex       sync.RWMutex
	resolveStopLRPInstanceArgsForCall []struct {
		stopInstance models.StopLRPInstance
	}
	resolveStopLRPInstanceReturns struct {
		result1 error
	}
}

func (fake *FakeRepBBS) MaintainExecutorPresence(heartbeatInterval time.Duration, executorPresence models.ExecutorPresence) (services_bbs.Presence, <-chan bool, error) {
	fake.maintainExecutorPresenceMutex.Lock()
	defer fake.maintainExecutorPresenceMutex.Unlock()
	fake.maintainExecutorPresenceArgsForCall = append(fake.maintainExecutorPresenceArgsForCall, struct {
		heartbeatInterval time.Duration
		executorPresence  models.ExecutorPresence
	}{heartbeatInterval, executorPresence})
	if fake.MaintainExecutorPresenceStub != nil {
		return fake.MaintainExecutorPresenceStub(heartbeatInterval, executorPresence)
	} else {
		return fake.maintainExecutorPresenceReturns.result1, fake.maintainExecutorPresenceReturns.result2, fake.maintainExecutorPresenceReturns.result3
	}
}

func (fake *FakeRepBBS) MaintainExecutorPresenceCallCount() int {
	fake.maintainExecutorPresenceMutex.RLock()
	defer fake.maintainExecutorPresenceMutex.RUnlock()
	return len(fake.maintainExecutorPresenceArgsForCall)
}

func (fake *FakeRepBBS) MaintainExecutorPresenceArgsForCall(i int) (time.Duration, models.ExecutorPresence) {
	fake.maintainExecutorPresenceMutex.RLock()
	defer fake.maintainExecutorPresenceMutex.RUnlock()
	return fake.maintainExecutorPresenceArgsForCall[i].heartbeatInterval, fake.maintainExecutorPresenceArgsForCall[i].executorPresence
}

func (fake *FakeRepBBS) MaintainExecutorPresenceReturns(result1 services_bbs.Presence, result2 <-chan bool, result3 error) {
	fake.MaintainExecutorPresenceStub = nil
	fake.maintainExecutorPresenceReturns = struct {
		result1 services_bbs.Presence
		result2 <-chan bool
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeRepBBS) WatchForDesiredTask() (<-chan models.Task, chan<- bool, <-chan error) {
	fake.watchForDesiredTaskMutex.Lock()
	defer fake.watchForDesiredTaskMutex.Unlock()
	fake.watchForDesiredTaskArgsForCall = append(fake.watchForDesiredTaskArgsForCall, struct{}{})
	if fake.WatchForDesiredTaskStub != nil {
		return fake.WatchForDesiredTaskStub()
	} else {
		return fake.watchForDesiredTaskReturns.result1, fake.watchForDesiredTaskReturns.result2, fake.watchForDesiredTaskReturns.result3
	}
}

func (fake *FakeRepBBS) WatchForDesiredTaskCallCount() int {
	fake.watchForDesiredTaskMutex.RLock()
	defer fake.watchForDesiredTaskMutex.RUnlock()
	return len(fake.watchForDesiredTaskArgsForCall)
}

func (fake *FakeRepBBS) WatchForDesiredTaskReturns(result1 <-chan models.Task, result2 chan<- bool, result3 <-chan error) {
	fake.WatchForDesiredTaskStub = nil
	fake.watchForDesiredTaskReturns = struct {
		result1 <-chan models.Task
		result2 chan<- bool
		result3 <-chan error
	}{result1, result2, result3}
}

func (fake *FakeRepBBS) ClaimTask(taskGuid string, executorID string) error {
	fake.claimTaskMutex.Lock()
	defer fake.claimTaskMutex.Unlock()
	fake.claimTaskArgsForCall = append(fake.claimTaskArgsForCall, struct {
		taskGuid   string
		executorID string
	}{taskGuid, executorID})
	if fake.ClaimTaskStub != nil {
		return fake.ClaimTaskStub(taskGuid, executorID)
	} else {
		return fake.claimTaskReturns.result1
	}
}

func (fake *FakeRepBBS) ClaimTaskCallCount() int {
	fake.claimTaskMutex.RLock()
	defer fake.claimTaskMutex.RUnlock()
	return len(fake.claimTaskArgsForCall)
}

func (fake *FakeRepBBS) ClaimTaskArgsForCall(i int) (string, string) {
	fake.claimTaskMutex.RLock()
	defer fake.claimTaskMutex.RUnlock()
	return fake.claimTaskArgsForCall[i].taskGuid, fake.claimTaskArgsForCall[i].executorID
}

func (fake *FakeRepBBS) ClaimTaskReturns(result1 error) {
	fake.ClaimTaskStub = nil
	fake.claimTaskReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeRepBBS) StartTask(taskGuid string, executorID string, containerHandle string) error {
	fake.startTaskMutex.Lock()
	defer fake.startTaskMutex.Unlock()
	fake.startTaskArgsForCall = append(fake.startTaskArgsForCall, struct {
		taskGuid        string
		executorID      string
		containerHandle string
	}{taskGuid, executorID, containerHandle})
	if fake.StartTaskStub != nil {
		return fake.StartTaskStub(taskGuid, executorID, containerHandle)
	} else {
		return fake.startTaskReturns.result1
	}
}

func (fake *FakeRepBBS) StartTaskCallCount() int {
	fake.startTaskMutex.RLock()
	defer fake.startTaskMutex.RUnlock()
	return len(fake.startTaskArgsForCall)
}

func (fake *FakeRepBBS) StartTaskArgsForCall(i int) (string, string, string) {
	fake.startTaskMutex.RLock()
	defer fake.startTaskMutex.RUnlock()
	return fake.startTaskArgsForCall[i].taskGuid, fake.startTaskArgsForCall[i].executorID, fake.startTaskArgsForCall[i].containerHandle
}

func (fake *FakeRepBBS) StartTaskReturns(result1 error) {
	fake.StartTaskStub = nil
	fake.startTaskReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeRepBBS) CompleteTask(taskGuid string, failed bool, failureReason string, result string) error {
	fake.completeTaskMutex.Lock()
	defer fake.completeTaskMutex.Unlock()
	fake.completeTaskArgsForCall = append(fake.completeTaskArgsForCall, struct {
		taskGuid      string
		failed        bool
		failureReason string
		result        string
	}{taskGuid, failed, failureReason, result})
	if fake.CompleteTaskStub != nil {
		return fake.CompleteTaskStub(taskGuid, failed, failureReason, result)
	} else {
		return fake.completeTaskReturns.result1
	}
}

func (fake *FakeRepBBS) CompleteTaskCallCount() int {
	fake.completeTaskMutex.RLock()
	defer fake.completeTaskMutex.RUnlock()
	return len(fake.completeTaskArgsForCall)
}

func (fake *FakeRepBBS) CompleteTaskArgsForCall(i int) (string, bool, string, string) {
	fake.completeTaskMutex.RLock()
	defer fake.completeTaskMutex.RUnlock()
	return fake.completeTaskArgsForCall[i].taskGuid, fake.completeTaskArgsForCall[i].failed, fake.completeTaskArgsForCall[i].failureReason, fake.completeTaskArgsForCall[i].result
}

func (fake *FakeRepBBS) CompleteTaskReturns(result1 error) {
	fake.CompleteTaskStub = nil
	fake.completeTaskReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeRepBBS) ReportActualLRPAsStarting(lrp models.ActualLRP, executorID string) error {
	fake.reportActualLRPAsStartingMutex.Lock()
	defer fake.reportActualLRPAsStartingMutex.Unlock()
	fake.reportActualLRPAsStartingArgsForCall = append(fake.reportActualLRPAsStartingArgsForCall, struct {
		lrp        models.ActualLRP
		executorID string
	}{lrp, executorID})
	if fake.ReportActualLRPAsStartingStub != nil {
		return fake.ReportActualLRPAsStartingStub(lrp, executorID)
	} else {
		return fake.reportActualLRPAsStartingReturns.result1
	}
}

func (fake *FakeRepBBS) ReportActualLRPAsStartingCallCount() int {
	fake.reportActualLRPAsStartingMutex.RLock()
	defer fake.reportActualLRPAsStartingMutex.RUnlock()
	return len(fake.reportActualLRPAsStartingArgsForCall)
}

func (fake *FakeRepBBS) ReportActualLRPAsStartingArgsForCall(i int) (models.ActualLRP, string) {
	fake.reportActualLRPAsStartingMutex.RLock()
	defer fake.reportActualLRPAsStartingMutex.RUnlock()
	return fake.reportActualLRPAsStartingArgsForCall[i].lrp, fake.reportActualLRPAsStartingArgsForCall[i].executorID
}

func (fake *FakeRepBBS) ReportActualLRPAsStartingReturns(result1 error) {
	fake.ReportActualLRPAsStartingStub = nil
	fake.reportActualLRPAsStartingReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeRepBBS) ReportActualLRPAsRunning(lrp models.ActualLRP, executorId string) error {
	fake.reportActualLRPAsRunningMutex.Lock()
	defer fake.reportActualLRPAsRunningMutex.Unlock()
	fake.reportActualLRPAsRunningArgsForCall = append(fake.reportActualLRPAsRunningArgsForCall, struct {
		lrp        models.ActualLRP
		executorId string
	}{lrp, executorId})
	if fake.ReportActualLRPAsRunningStub != nil {
		return fake.ReportActualLRPAsRunningStub(lrp, executorId)
	} else {
		return fake.reportActualLRPAsRunningReturns.result1
	}
}

func (fake *FakeRepBBS) ReportActualLRPAsRunningCallCount() int {
	fake.reportActualLRPAsRunningMutex.RLock()
	defer fake.reportActualLRPAsRunningMutex.RUnlock()
	return len(fake.reportActualLRPAsRunningArgsForCall)
}

func (fake *FakeRepBBS) ReportActualLRPAsRunningArgsForCall(i int) (models.ActualLRP, string) {
	fake.reportActualLRPAsRunningMutex.RLock()
	defer fake.reportActualLRPAsRunningMutex.RUnlock()
	return fake.reportActualLRPAsRunningArgsForCall[i].lrp, fake.reportActualLRPAsRunningArgsForCall[i].executorId
}

func (fake *FakeRepBBS) ReportActualLRPAsRunningReturns(result1 error) {
	fake.ReportActualLRPAsRunningStub = nil
	fake.reportActualLRPAsRunningReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeRepBBS) RemoveActualLRP(lrp models.ActualLRP) error {
	fake.removeActualLRPMutex.Lock()
	defer fake.removeActualLRPMutex.Unlock()
	fake.removeActualLRPArgsForCall = append(fake.removeActualLRPArgsForCall, struct {
		lrp models.ActualLRP
	}{lrp})
	if fake.RemoveActualLRPStub != nil {
		return fake.RemoveActualLRPStub(lrp)
	} else {
		return fake.removeActualLRPReturns.result1
	}
}

func (fake *FakeRepBBS) RemoveActualLRPCallCount() int {
	fake.removeActualLRPMutex.RLock()
	defer fake.removeActualLRPMutex.RUnlock()
	return len(fake.removeActualLRPArgsForCall)
}

func (fake *FakeRepBBS) RemoveActualLRPArgsForCall(i int) models.ActualLRP {
	fake.removeActualLRPMutex.RLock()
	defer fake.removeActualLRPMutex.RUnlock()
	return fake.removeActualLRPArgsForCall[i].lrp
}

func (fake *FakeRepBBS) RemoveActualLRPReturns(result1 error) {
	fake.RemoveActualLRPStub = nil
	fake.removeActualLRPReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeRepBBS) RemoveActualLRPForIndex(processGuid string, index int, instanceGuid string) error {
	fake.removeActualLRPForIndexMutex.Lock()
	defer fake.removeActualLRPForIndexMutex.Unlock()
	fake.removeActualLRPForIndexArgsForCall = append(fake.removeActualLRPForIndexArgsForCall, struct {
		processGuid  string
		index        int
		instanceGuid string
	}{processGuid, index, instanceGuid})
	if fake.RemoveActualLRPForIndexStub != nil {
		return fake.RemoveActualLRPForIndexStub(processGuid, index, instanceGuid)
	} else {
		return fake.removeActualLRPForIndexReturns.result1
	}
}

func (fake *FakeRepBBS) RemoveActualLRPForIndexCallCount() int {
	fake.removeActualLRPForIndexMutex.RLock()
	defer fake.removeActualLRPForIndexMutex.RUnlock()
	return len(fake.removeActualLRPForIndexArgsForCall)
}

func (fake *FakeRepBBS) RemoveActualLRPForIndexArgsForCall(i int) (string, int, string) {
	fake.removeActualLRPForIndexMutex.RLock()
	defer fake.removeActualLRPForIndexMutex.RUnlock()
	return fake.removeActualLRPForIndexArgsForCall[i].processGuid, fake.removeActualLRPForIndexArgsForCall[i].index, fake.removeActualLRPForIndexArgsForCall[i].instanceGuid
}

func (fake *FakeRepBBS) RemoveActualLRPForIndexReturns(result1 error) {
	fake.RemoveActualLRPForIndexStub = nil
	fake.removeActualLRPForIndexReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeRepBBS) WatchForStopLRPInstance() (<-chan models.StopLRPInstance, chan<- bool, <-chan error) {
	fake.watchForStopLRPInstanceMutex.Lock()
	defer fake.watchForStopLRPInstanceMutex.Unlock()
	fake.watchForStopLRPInstanceArgsForCall = append(fake.watchForStopLRPInstanceArgsForCall, struct{}{})
	if fake.WatchForStopLRPInstanceStub != nil {
		return fake.WatchForStopLRPInstanceStub()
	} else {
		return fake.watchForStopLRPInstanceReturns.result1, fake.watchForStopLRPInstanceReturns.result2, fake.watchForStopLRPInstanceReturns.result3
	}
}

func (fake *FakeRepBBS) WatchForStopLRPInstanceCallCount() int {
	fake.watchForStopLRPInstanceMutex.RLock()
	defer fake.watchForStopLRPInstanceMutex.RUnlock()
	return len(fake.watchForStopLRPInstanceArgsForCall)
}

func (fake *FakeRepBBS) WatchForStopLRPInstanceReturns(result1 <-chan models.StopLRPInstance, result2 chan<- bool, result3 <-chan error) {
	fake.WatchForStopLRPInstanceStub = nil
	fake.watchForStopLRPInstanceReturns = struct {
		result1 <-chan models.StopLRPInstance
		result2 chan<- bool
		result3 <-chan error
	}{result1, result2, result3}
}

func (fake *FakeRepBBS) ResolveStopLRPInstance(stopInstance models.StopLRPInstance) error {
	fake.resolveStopLRPInstanceMutex.Lock()
	defer fake.resolveStopLRPInstanceMutex.Unlock()
	fake.resolveStopLRPInstanceArgsForCall = append(fake.resolveStopLRPInstanceArgsForCall, struct {
		stopInstance models.StopLRPInstance
	}{stopInstance})
	if fake.ResolveStopLRPInstanceStub != nil {
		return fake.ResolveStopLRPInstanceStub(stopInstance)
	} else {
		return fake.resolveStopLRPInstanceReturns.result1
	}
}

func (fake *FakeRepBBS) ResolveStopLRPInstanceCallCount() int {
	fake.resolveStopLRPInstanceMutex.RLock()
	defer fake.resolveStopLRPInstanceMutex.RUnlock()
	return len(fake.resolveStopLRPInstanceArgsForCall)
}

func (fake *FakeRepBBS) ResolveStopLRPInstanceArgsForCall(i int) models.StopLRPInstance {
	fake.resolveStopLRPInstanceMutex.RLock()
	defer fake.resolveStopLRPInstanceMutex.RUnlock()
	return fake.resolveStopLRPInstanceArgsForCall[i].stopInstance
}

func (fake *FakeRepBBS) ResolveStopLRPInstanceReturns(result1 error) {
	fake.ResolveStopLRPInstanceStub = nil
	fake.resolveStopLRPInstanceReturns = struct {
		result1 error
	}{result1}
}

var _ bbs.RepBBS = new(FakeRepBBS)
