package schedule

import (
	"errors"
	"time"
)

type ScheduleOption struct {
}

type Schedule struct {
	Options         ScheduleOption
	schedules       map[string]func(time.Time, map[string]string, interface{}) (*time.Time, error)
	scheduleOptions map[string]interface{}
}

func (s *Schedule) CheckSupportedSchedule(name string) bool {
	return s.schedules[name] != nil
}

func (s *Schedule) AddSchedule(name string, f func(time.Time, map[string]string, interface{}) (*time.Time, error), options interface{}) {
	if s.schedules == nil {
		s.schedules = make(map[string]func(time.Time, map[string]string, interface{}) (*time.Time, error))
	}
	if s.scheduleOptions == nil {
		s.scheduleOptions = make(map[string]interface{})
	}
	s.schedules[name] = f
	s.scheduleOptions[name] = options
}

func (s *Schedule) Schedule(pivotTime time.Time, name string, params map[string]string) (*time.Time, error) {
	if s.schedules[name] == nil {
		return nil, errors.New("schedule not found")
	}
	return s.schedules[name](pivotTime, params, s.scheduleOptions[name])
}

var Scheduler = Schedule{}
