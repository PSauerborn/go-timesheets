package main

import (
    "time"
    "github.com/google/uuid"
)

var (

)

type WorkPeriod struct {
    PeriodId   uuid.UUID     `json:"periodId"`
    CreatedAt  time.Time     `json:"createdAt"`
    FinishedAt *time.Time    `json:"finishedAt,omitempty"`
    Breaks 	   []BreakPeriod `json:"breaks"`
}

type BreakPeriod struct {
    BreakId    uuid.UUID  `json:"breakId"`
    CreatedAt  time.Time  `json:"createdAt"`
    FinishedAt *time.Time `json:"finishedAt,omitempty"`
}

type UserData struct {
    Uid	        string		 `json:"uid"`
    WorkPeriods []WorkPeriod `json:"workPeriods"`
}