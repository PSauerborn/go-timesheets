package main

import (
    "fmt"
    "time"
    "context"
    "github.com/jackc/pgx/v4"
    "github.com/jackc/pgx/v4/pgxpool"
    "github.com/google/uuid"
    log "github.com/sirupsen/logrus"
)

var persistence *Persistence

type Persistence struct {
    conn *pgxpool.Pool
}

// function used to connect postgres connection
func ConnectPersistence() {
    log.Info(fmt.Sprintf("attempting postgres connection with connection string %s", PostgresConnection))
    db, err := pgxpool.Connect(context.Background(), PostgresConnection)
    if err != nil {
        log.Fatal(fmt.Errorf("unable to connect to postgres server: %v", err))
    }
    log.Info("successfully connected to postgres")
    // connect persistence and assign to persistence var
    persistence = &Persistence{db}
}

func(db Persistence) createWorkPeriod(uid string) (uuid.UUID, error) {
    log.Debug(fmt.Sprintf("creating new work period for user %s", uid))
    periodId := uuid.New()
    _, err := db.conn.Exec(context.Background(), "INSERT INTO work_periods(period_id, uid) VALUES($1,$2)", periodId, uid)
    if err != nil {
        log.Error(fmt.Errorf("unable to create new work period: %v", err))
        return periodId, err
    }
    log.Info(fmt.Sprintf("successfully created new work period with ID %s", periodId))
    return periodId, nil
}

func(db Persistence) createBreakPeriod(periodId uuid.UUID) (uuid.UUID, error) {
    log.Debug(fmt.Sprintf("creating new break period for work period %s", periodId))
    breakId := uuid.New()
    _, err := db.conn.Exec(context.Background(), "INSERT INTO break_periods(break_id, period_id) VALUES($1,$2)", breakId, periodId)
    if err != nil {
        log.Error(fmt.Errorf("unable to create new work period: %v", err))
        return breakId, err
    }
    log.Info(fmt.Sprintf("successfully created new break period %s", breakId))
    return breakId, nil
}

func(db Persistence) getUserData(uid string) (UserData, error) {
    log.Debug(fmt.Sprintf("fetching data for user %s", uid))

    rows, err := db.conn.Query(context.Background(), "SELECT period_id FROM work_periods WHERE uid=$1", uid)
    if err != nil {
        log.Error(fmt.Errorf("unable to retrieve work periods for user %s: %v", uid, err))
        switch err {
        case pgx.ErrNoRows:
            return UserData{ Uid: uid, WorkPeriods: []WorkPeriod{}}, nil
        default:
            return UserData{}, err
        }
    }

    periods := []WorkPeriod{}
    // iterate over period ID's and retrieve full period
    for rows.Next() {
        var periodId uuid.UUID
        err := rows.Scan(&periodId)
        if err != nil {
            log.Error(fmt.Errorf("unable to process work period: %v", err))
            continue
        }

        // retrieve period and all breaks from database
        period, err := db.getWorkPeriod(periodId)
        if err != nil {
            log.Error(fmt.Errorf("unable to retrieve work period %s", periodId))
        } else {
            periods = append(periods, period)
        }
    }
    return UserData{Uid: uid, WorkPeriods: periods}, nil
}

func(db Persistence) getUserDataOverRange(uid string, start, end time.Time) (UserData, error) {
    log.Debug(fmt.Sprintf("fetching data for user %s", uid))

    rows, err := db.conn.Query(context.Background(), "SELECT period_id FROM work_periods WHERE uid=$1 AND created_at > $2 AND created_at < $3", uid, start, end)
    if err != nil {
        log.Error(fmt.Errorf("unable to retrieve work periods for user %s: %v", uid, err))
        switch err {
        case pgx.ErrNoRows:
            return UserData{ Uid: uid, WorkPeriods: []WorkPeriod{}}, nil
        default:
            return UserData{}, err
        }
    }

    periods := []WorkPeriod{}
    // iterate over period ID's and retrieve full period
    for rows.Next() {
        var periodId uuid.UUID
        err := rows.Scan(&periodId)
        if err != nil {
            log.Error(fmt.Errorf("unable to process work period: %v", err))
            continue
        }

        // retrieve period and all breaks from database
        period, err := db.getWorkPeriod(periodId)
        if err != nil {
            log.Error(fmt.Errorf("unable to retrieve work period %s", periodId))
        } else {
            periods = append(periods, period)
        }
    }
    return UserData{Uid: uid, WorkPeriods: periods}, nil
}


func(db Persistence) getBreakPeriod(breakId uuid.UUID) (BreakPeriod, error) {
    log.Debug(fmt.Sprintf("retrieving break period %s", breakId))

    var (periodId uuid.UUID; createdAt time.Time; finishedAt *time.Time)
    breakPeriod := db.conn.QueryRow(context.Background(), "SELECT period_id,created_at,finished_at FROM break_periods WHERE break_id=$1", breakId)
    err := breakPeriod.Scan(&periodId, &createdAt, &finishedAt)
    if err != nil {
        log.Error(fmt.Errorf("unable to retrieve break period %s: %v", breakId, err))
        return BreakPeriod{}, err
    }
    return BreakPeriod{BreakId: breakId, CreatedAt: createdAt, FinishedAt: finishedAt}, nil
}

func(db Persistence) getBreakPeriods(periodId uuid.UUID) ([]BreakPeriod, error) {
    log.Debug(fmt.Sprintf("retrieving break periods for work period %s", periodId))
    breaks := []BreakPeriod{}

    rows, err := db.conn.Query(context.Background(), "SELECT break_id FROM break_periods WHERE period_id=$1", periodId)
    if err != nil {
        log.Error(fmt.Errorf("unable to retrieve break periods for period ID %s: %v", periodId, err))
        switch err {
        case pgx.ErrNoRows:
            return breaks, nil
        default:
            return breaks, err
        }
    }

    for rows.Next() {
        var breakId uuid.UUID
        err := rows.Scan(&breakId)
        if err != nil {
            log.Error(fmt.Errorf("unable to parse breakId: %v", err))
            continue
        }
        breakPeriod, err := db.getBreakPeriod(breakId)
        if err != nil {
            log.Error(fmt.Errorf("unable to retrieve break %s: %v", breakId, err))
        } else {
            breaks = append(breaks, breakPeriod)
        }
    }
    return breaks, nil
}

func(db Persistence) getWorkPeriod(periodId uuid.UUID) (WorkPeriod, error) {
    log.Debug(fmt.Sprintf("retrieving work period %s", periodId))

    var (createdAt time.Time; finishedAt *time.Time)

    period := db.conn.QueryRow(context.Background(), "SELECT created_at,finished_at FROM work_periods WHERE period_id=$1", periodId)
    err := period.Scan(&createdAt, &finishedAt)
    if err != nil {
        log.Error(fmt.Errorf("unable to retrieve work period %s: %v", periodId, err))
        return WorkPeriod{}, err
    }
    // get break periods for the work period from database
    breaks, err := db.getBreakPeriods(periodId)
    if err != nil {
        log.Error(fmt.Errorf("unable to get break periods for work period %s: %v", periodId, err))
        return WorkPeriod{}, nil
    }
    return WorkPeriod{PeriodId: periodId, CreatedAt: createdAt, FinishedAt: finishedAt, Breaks: breaks}, nil
}

func(db Persistence) closeWorkPeriod(periodId uuid.UUID) error {
    log.Debug(fmt.Sprintf("closing work period %s", periodId))
    _, err := db.conn.Exec(context.Background(), "UPDATE work_periods SET finished_at=$1 WHERE period_id=$2", time.Now(), periodId)
    if err != nil {
        log.Error(fmt.Errorf("unable to close work period %s: %v", periodId, err))
        return err
    }
    log.Info(fmt.Sprintf("successfully updated work period %s", periodId))
    return nil
}

func(db Persistence) closeBreakPeriod(breakId uuid.UUID) error {
    log.Debug(fmt.Sprintf("closing work break %s", breakId))
    _, err := db.conn.Exec(context.Background(), "UPDATE break_periods SET finished_at=$1 WHERE break_id=$2", time.Now(), breakId)
    if err != nil {
        log.Error(fmt.Errorf("unable to close work period %s: %v", breakId, err))
        return err
    }
    log.Info(fmt.Sprintf("successfully updated work period %s", breakId))
    return nil
}