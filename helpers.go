package main

import (
    "github.com/jackc/pgx/v4"
    "github.com/google/uuid"
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