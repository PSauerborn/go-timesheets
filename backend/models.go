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

func(period WorkPeriod) TotalHours() float64 {
    finished := *(period.FinishedAt)
    return finished.Sub(period.CreatedAt).Hours()
}

type ActiveWorkPeriod struct {
    PeriodId    uuid.UUID          `json:"periodId"`
    CreatedAt   time.Time          `json:"createdAt"`
    ActiveSince float64            `json:"activeSince"`
    ActiveBreak *ActiveBreakPeriod `json:"activeBreak"`
}

type BreakPeriod struct {
    BreakId    uuid.UUID  `json:"breakId"`
    CreatedAt  time.Time  `json:"createdAt"`
    FinishedAt *time.Time `json:"finishedAt,omitempty"`
}

func(period BreakPeriod) TotalHours() float64 {
    finished := *(period.FinishedAt)
    return finished.Sub(period.CreatedAt).Hours()
}

type ActiveBreakPeriod struct {
    BreakId   uuid.UUID `json:"breakId"`
    CreatedAt time.Time `json:"createdAt"`
}

type UserData struct {
    Uid	        string		 `json:"uid"`
    WorkPeriods []WorkPeriod `json:"workPeriods"`
}

type RangedAnalysisResults struct {
    Start   time.Time       `json:"start"`
    End     time.Time       `json:"end"`
    Results AnalysisResults `json:"results"`
}

type AnalysisResults struct {
    TotalPeriods    int     `json:"totalPeriods"`
    TotalBreaks     int     `json:"totalBreaks"`
    TotalWorkHours  float64 `json:"totalWorkHours"`
    TotalBreakHours float64 `json:"totalBreakHours"`
    NetWorkHours    float64 `json:"netWorkHours"`
}

type WorkPeriodAnalysisResults struct {
    PeriodCount int    `json:"totalPeriods"`
    TotalHours float64 `json:"totalHours"`
}

type BreakPeriodAnalysisResults struct {
    BreakCount int     `json:"totalPeriods"`
    TotalHours float64 `json:"totalHours"`
}

type BucketAnalysis struct {
    TotalWorkHours  float64   `json:"totalWorkHours"`
    TotalBreakHours float64   `json:"totalBreakHours"`
    NetWorkHours    float64   `json:"netWorkHours"`
    StartTime       time.Time `json:"startTime"`
    EndTime         time.Time `json:"endTime"`
    TotalPeriods    int       `json:"totalPeriods"`
    TotalBreaks     int       `json:"totalBreaks"`
}