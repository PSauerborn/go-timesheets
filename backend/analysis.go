package main

import (
    "fmt"
    "time"
    log "github.com/sirupsen/logrus"
)

var (

)

// function used to analyze list of breaks. both the total number of
// breaks as well as the total number of break hours are returned
func analyseBreaks(breaks []BreakPeriod) BreakPeriodAnalysisResults {
    breakHours := 0.0
    // iterate over breaks and increment total break time
    for _, period := range(breaks) {
        if period.FinishedAt != nil {
            breakHours += period.TotalHours()
        }
    }
    return BreakPeriodAnalysisResults{BreakCount: len(breaks), TotalHours: breakHours}
}

// function used to analyze a list of periods (including breaks)
// all periods are iterated over and the breaks within each period
// are also aggregated and analyzed
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
        if (len(period.Breaks) > 0) {
            breakAnalysis := analyseBreaks(period.Breaks)
            results.TotalWorkHours -= breakAnalysis.TotalHours
            results.TotalBreaks += breakAnalysis.BreakCount
        }
    }
    // evaluate net work hours from total work hours and breaks
    results.NetWorkHours = results.TotalWorkHours - results.TotalBreakHours
    return results
}

// function used to analyse all user tasks. note that all history tasks
// are analysed and returned in the response
func analyzeUserTasks(uid string) (AnalysisResults, error) {
    log.Info(fmt.Sprintf("performaning analysis for user %s", uid))
    results, err := persistence.getUserData(uid)
    if err != nil {
        log.Error(fmt.Errorf("unable to get user data: %v", err))
        return AnalysisResults{}, err
    }
    return analysePeriods(results.WorkPeriods), nil
}

// function used to analyse users tasks over a period of time
func analyseRangedUserTasks(uid string, start, end time.Time) (AnalysisResults, error) {
    log.Info(fmt.Sprintf("performaning analysis for user %s over range %s - %s", uid, start, end))
    results, err := persistence.getUserDataOverRange(uid, start, end)
    if err != nil {
        log.Error(fmt.Errorf("unable to get user data: %v", err))
        return AnalysisResults{}, err
    }
    return analysePeriods(results.WorkPeriods), nil
}

// function used to aggregate break periods by date. all periods that
// that are created during the same day as the given date are returned
func aggregatePeriods(periods []WorkPeriod, date time.Time) []WorkPeriod {
    aggregate := []WorkPeriod{}
    for _, period := range(periods) {
        // if period is on same day as given date, add to array
        if (period.CreatedAt.After(date) && period.CreatedAt.Before(date.Add(time.Hour * 24))) {
            aggregate = append(aggregate, period)
        }
    }
    return aggregate
}

// function used to group work periods into daily buckets. values are returned as
// a map of {<date>: [ periods... ]}
func groupPeriodsByDay(periods []WorkPeriod, start, end time.Time) map[string][]WorkPeriod {
    aggregatedPeriods := map[string][]WorkPeriod{}
    date := start
    for date.Before(end) {
        aggregatedPeriods[date.Format("2006-01-02")] = aggregatePeriods(periods, date)
        date = date.Add(time.Hour * 24)
    }
    return aggregatedPeriods
}

// ###########################################################
// # Define functions used to bucket and analyse bucketed data
// ###########################################################

// function used to analyse bucketed data
// func analyseBuckets(buckets map[time.Time]WorkPeriod) {
//     for bucket, periods := range(buckets) {
//         log.Debug(fmt.Sprintf("processing bucket %+v with %d periods", bucket, len(periods)))
//     }
// }

// function used to determine if period falls within a given bucket
func fallsInBucket(period WorkPeriod, date time.Time, bucket int) bool {
    return period.CreatedAt.After(date) && period.CreatedAt.Before(date.Add(time.Minute * time.Duration(bucket)))
}

// function used to bucket periods around a given date
func bucket(periods []WorkPeriod, date time.Time, bucket int) []WorkPeriod {
    aggregated := []WorkPeriod{}
    for _, period := range(periods) {
        if fallsInBucket(period, date, bucket) {
            aggregated = append(aggregated, period)
        }
    }
    return aggregated
}

// function used to bucket periods by a bucket size (given in minutes)
func bucketPeriods(periods []WorkPeriod, start, end time.Time, bucketSize int) map[time.Time][]WorkPeriod {
    log.Debug(fmt.Sprintf("bucketing periods over date range %s - %s with bucket %d", start, end, bucketSize))
    bucketed := map[time.Time][]WorkPeriod{}
    for start.Before(end) {
        bucketed[start] = bucket(periods, start, bucketSize)
        start = start.Add(time.Minute * time.Duration(bucketSize))
    }
    return bucketed
}