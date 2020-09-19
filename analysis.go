package main

import (
	"fmt"
	"time"
	log "github.com/sirupsen/logrus"
)

var (

)

func analyseBreaks(breaks []BreakPeriod) BreakPeriodAnalysisResults {
	breakHours := 0.0
	for _, period := range(breaks) {
		if period.FinishedAt != nil {
			breakHours += period.TotalHours()
		}
	}
	return BreakPeriodAnalysisResults{BreakCount: len(breaks), TotalHours: breakHours}
}

func analysePeriods(periods []WorkPeriod) AnalysisResults {
	results := AnalysisResults{}
	// iterate over work periods and perform analysis
	for _, period := range(periods) {
		results.TotalPeriods += 1
		// if period has finish date, evaluate work period
		if period.FinishedAt != nil {
			results.TotalWorkHours += period.TotalHours()
		}
		// perform analysis on breaks and add total to results
		breakAnalysis := analyseBreaks(period.Breaks)
		results.TotalWorkHours -= breakAnalysis.TotalHours
		results.TotalBreaks += breakAnalysis.BreakCount
	}
	// evaluate net work hours from total work hours and breaks
	results.NetWorkHours = results.TotalWorkHours - results.TotalBreakHours
	return results
}

func analyzeUserTasks(uid string) (AnalysisResults, error) {
	log.Info(fmt.Sprintf("performaning analysis for user %s", uid))
	results, err := persistence.getUserData(uid)
	if err != nil {
		log.Error(fmt.Errorf("unable to get user data: %v", err))
		return AnalysisResults{}, err
	}
	return analysePeriods(results.WorkPeriods), nil
}

func analyseRangedUserTasks(uid string, start, end time.Time) (AnalysisResults, error) {
	log.Info(fmt.Sprintf("performaning analysis for user %s over range %s - %s", uid, start, end))
	results, err := persistence.getUserDataOverRange(uid, start, end)
	if err != nil {
		log.Error(fmt.Errorf("unable to get user data: %v", err))
		return AnalysisResults{}, err
	}
	return analysePeriods(results.WorkPeriods), nil
}