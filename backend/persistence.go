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

// function used to create new work period in postgres datebase
func(db Persistence) createWorkPeriod(uid string) (ActiveWorkPeriod, error) {
    log.Debug(fmt.Sprintf("creating new work period for user %s", uid))
    periodId := uuid.New()
    now := time.Now()
    // create new work period and parse into ActiveWorkPeriod struct
    _, err := db.conn.Exec(context.Background(), "INSERT INTO work_periods(period_id, uid, created_at) VALUES($1,$2,$3)", periodId, uid, now)
    if err != nil {
        log.Error(fmt.Errorf("unable to create new work period: %v", err))
        return ActiveWorkPeriod{}, err
    }
    log.Info(fmt.Sprintf("successfully created new work period with ID %s", periodId))
    return ActiveWorkPeriod{PeriodId: periodId, CreatedAt: now}, nil
}

// function used to create new break period in database
func(db Persistence) createBreakPeriod(periodId uuid.UUID) (uuid.UUID, error) {
    log.Debug(fmt.Sprintf("creating new break period for work period %s", periodId))
    breakId := uuid.New()
    // create new break period and insert into database
    _, err := db.conn.Exec(context.Background(), "INSERT INTO break_periods(break_id, period_id) VALUES($1,$2)", breakId, periodId)
    if err != nil {
        log.Error(fmt.Errorf("unable to create new work period: %v", err))
        return breakId, err
    }
    log.Info(fmt.Sprintf("successfully created new break period %s", breakId))
    return breakId, nil
}

// function used to retrieve user data. all work periods are
// retrieved first, and the list of work periods is then used
// to retrieve the list of break periods, which are all combined
func(db Persistence) getUserData(uid string) (UserData, error) {
    log.Debug(fmt.Sprintf("fetching data for user %s", uid))
    // retrieve all periods from database for user
    rows, err := db.conn.Query(context.Background(), "SELECT period_id FROM work_periods WHERE uid=$1 AND finished_at IS NOT NULL", uid)
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

// function used to retrieve user data within a specific time range.
// note that this is the equivalent of db.getUserData() with the
// additional timestamp constraint
func(db Persistence) getUserDataOverRange(uid string, start, end time.Time) (UserData, error) {
    log.Debug(fmt.Sprintf("fetching data for user %s", uid))
    // retrieve all periods from database what are completed
    rows, err := db.conn.Query(context.Background(), "SELECT period_id FROM work_periods WHERE uid=$1 AND created_at > $2 AND created_at < $3 AND finished_at IS NOT NULL", uid, start, end)
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

// function used to get a specific break period from the database
func(db Persistence) getBreakPeriod(breakId uuid.UUID) (BreakPeriod, error) {
    log.Debug(fmt.Sprintf("retrieving break period %s", breakId))

    var (periodId uuid.UUID; createdAt time.Time; finishedAt *time.Time)
    // execute postgres query to get break from database
    breakPeriod := db.conn.QueryRow(context.Background(), "SELECT period_id,created_at,finished_at FROM break_periods WHERE break_id=$1", breakId)
    err := breakPeriod.Scan(&periodId, &createdAt, &finishedAt)
    if err != nil {
        log.Error(fmt.Errorf("unable to retrieve break period %s: %v", breakId, err))
        return BreakPeriod{}, err
    }
    return BreakPeriod{BreakId: breakId, CreatedAt: createdAt, FinishedAt: finishedAt}, nil
}

// function used to retrieve all break periods associated with a particular
// work period
func(db Persistence) getBreakPeriods(periodId uuid.UUID) ([]BreakPeriod, error) {
    log.Debug(fmt.Sprintf("retrieving break periods for work period %s", periodId))
    breaks := []BreakPeriod{}
    // get all break ID's associated with period ID
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
        // retrieve break details from database using break ID
        breakPeriod, err := db.getBreakPeriod(breakId)
        if err != nil {
            log.Error(fmt.Errorf("unable to retrieve break %s: %v", breakId, err))
        } else {
            breaks = append(breaks, breakPeriod)
        }
    }
    return breaks, nil
}

// function used to retrieve work period from database given a work period ID
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

// function used to close work period given work period ID
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

// function used to close break period given particular break period ID
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

// function used to retrieve current break period from database.
// this is done by getting all work periods, ordering by timstamp
// and selecting the latest entry. Note that only non-completed
// periods are selected
func(db Persistence) getActivePeriod(uid string) (ActiveWorkPeriod, error) {
    log.Debug(fmt.Sprintf("retrieving active work period for user %s", uid))

    var (periodId uuid.UUID; createdAt time.Time)
    period := db.conn.QueryRow(context.Background(), "SELECT period_id, created_at FROM work_periods WHERE uid=$1 AND finished_at IS NULL ORDER BY created_at DESC LIMIT 1", uid)
    err := period.Scan(&periodId, &createdAt)
    if err != nil {
        log.Error(fmt.Errorf("unable to retrieve active user period for user %s", uid))
        return ActiveWorkPeriod{}, err
    }
    // evaluate time that period has been active for given current date and created
    active := time.Now().Sub(createdAt)
    return ActiveWorkPeriod{PeriodId: periodId, CreatedAt: createdAt, ActiveSince: active.Hours()}, nil
}