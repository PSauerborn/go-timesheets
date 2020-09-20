package main

import (
    "fmt"
    "time"
    "errors"
    "github.com/jackc/pgx/v4"
    "github.com/google/uuid"
    log "github.com/sirupsen/logrus"
)

var (

)


func isValidWorkPeriod(periodId uuid.UUID) (bool, error) {
    _, err := persistence.getWorkPeriod(periodId)
    if err != nil {
        switch err {
        case pgx.ErrNoRows:
            return false, nil
        default:
            return false, err
        }
    }
    return true, nil
}

func isValidBreakPeriod(breakId uuid.UUID) (bool, error) {
    _, err := persistence.getBreakPeriod(breakId)
    if err != nil {
        switch err {
        case pgx.ErrNoRows:
            return false, nil
        default:
            return false, err
        }
    }
    return true, nil
}

func userOwnsPeriod(uid string, periodId uuid.UUID) bool {
    return true
}

func parseTimestamps(start, end string) (time.Time, time.Time, error) {
    layout := "2006-01-02"
    //parse start time
    startTime, err := time.Parse(layout, start)
    if err != nil {
        log.Error(fmt.Errorf("unable to parse start timestamp '%s': %v", start, err))
        return time.Now(), time.Now(), err
    }
    // parse end time
    endTime, err := time.Parse(layout, end)
    if err != nil {
        log.Error(fmt.Errorf("unable to parse end timestamp '%s': %v", end, err))
        return time.Now(), time.Now(), err
    }
    // check if start time is larger than end time
    if startTime.After(endTime) {
        return startTime, endTime, errors.New("start time cannot be larger than end time")
    }
    return startTime, endTime, nil
}


func aggregatePeriods(periods []WorkPeriod, date time.Time) []WorkPeriod {
    aggregate := []WorkPeriod{}
    for _, period := range(periods) {
        if (period.CreatedAt.After(date) && period.CreatedAt.Before(date.Add(time.Hour * 24))) {
            aggregate = append(aggregate, period)
        }
    }
    return aggregate
}

func groupPeriodsByDay(periods []WorkPeriod, start, end time.Time) map[string][]WorkPeriod {
    aggregatedPeriods := map[string][]WorkPeriod{}
    date := start
    for date.Before(end) {
        aggregatedPeriods[date.Format("2006-01-02")] = aggregatePeriods(periods, date)
        date = date.Add(time.Hour * 24)
    }
    return aggregatedPeriods
}