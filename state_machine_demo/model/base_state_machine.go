package model

type BaseStatus interface {
}

type BaseEvent interface {
}

type StatusEventPair struct {
	status BaseStatus
	event  BaseEvent
}

func (pair StatusEventPair) equals(other StatusEventPair) bool {
	return pair.status == other.status && pair.event == other.event
}

type StateMachine struct {
	statusEventMap map[StatusEventPair]BaseStatus
}

func (sm *StateMachine) accept(sourceStatus BaseStatus, event BaseEvent, targetStatus BaseStatus) {
	pair := StatusEventPair{status: sourceStatus, event: event}
	sm.statusEventMap[pair] = targetStatus
}

func (sm *StateMachine) getTargetStatus(sourceStatus BaseStatus, event BaseEvent) BaseStatus {
	pair := StatusEventPair{status: sourceStatus, event: event}
	baseStatus := sm.statusEventMap[pair]
	return baseStatus
}
