package main

import (
    "fmt"
    "time"
    "strings"
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "github.com/jackc/pgx/v4"
    log "github.com/sirupsen/logrus"
)

var (

)

func main() {

    ConfigureService()
    ConnectPersistence()

    router := gin.New()

    router.GET("/go-timesheets/health", healthCheckHandler)
    router.GET("/go-timesheets/active", getActivePeriodHandler)
    router.GET("/go-timesheets/data", getUserDataHandler)
    router.GET("/go-timesheets/data/:start/:end", getUserTimeRangeDataHandler)

    router.GET("/go-timesheets/analyse", getUserAnalysisHandler)
    router.GET("/go-timesheets/analyse/:start/:end", getUserTimeRangeAnalysisHandler)

    router.POST("/go-timesheets/work_period", createWorkPeriodHandler)
    router.POST("/go-timesheets/break_period/:periodId", createBreakPeriodHandler)

    router.PATCH("/go-timesheets/work_period/:periodId", endWorkPeriodHandler)
    router.PATCH("/go-timesheets/break_period/:breakId", endBreakPeriodHandler)

    router.Run(fmt.Sprintf(":%d", ListenPort))
}


func getUser(ctx *gin.Context) string {
    return ctx.Request.Header.Get("X-Authenticated-Userid")
}

func healthCheckHandler(ctx *gin.Context) {
    log.Debug("received request for health check route")
    StandardHTTP.Success(ctx)
}

func getActivePeriodHandler(ctx *gin.Context) {
    user := getUser(ctx)
    log.Debug(fmt.Sprintf("received request to get active peroid for user %s", user))
    period, err := persistence.getActivePeriod(user)
    if err != nil {
        switch err {
        case pgx.ErrNoRows:
            StandardHTTP.NotFound(ctx)
            return
        default:
            StandardHTTP.InternalServerError(ctx)
            return
        }
    }
    ctx.JSON(200, gin.H{"success": true, "http_code": 200, "payload": period})
}

func getUserDataHandler(ctx *gin.Context) {
    user := getUser(ctx)
    log.Debug(fmt.Sprintf("received request to get user data for user %s", user))
    // get user data from postgres database
    data, err := persistence.getUserData(user)
    if err != nil {
        log.Error(fmt.Errorf("unable to retrieve data for user %s: %v", user, err))
        StandardHTTP.InternalServerError(ctx)
        return
    }
    ctx.JSON(200, gin.H{"success": true, "http_code": 200, "data": data})
}

func getUserTimeRangeDataHandler(ctx *gin.Context) {
    user := getUser(ctx)
    start, end, err := parseTimestamps(ctx.Param("start"), ctx.Param("end"))
    if err != nil {
        log.Error(fmt.Errorf("unable to parse timestamps: %v", err))
        StandardHTTP.InvalidRequest(ctx)
        return
    }

    log.Debug(fmt.Sprintf("received request to get user data for user %s", user))
    // get user data from postgres database
    data, err := persistence.getUserDataOverRange(user, start, end.Add(time.Hour * 24))
    if err != nil {
        log.Error(fmt.Errorf("unable to retrieve data for user %s: %v", user, err))
        StandardHTTP.InternalServerError(ctx)
        return
    }
    // group values by day if specified in query parameters
    groupValues := ctx.DefaultQuery("group", "false")
    if strings.ToLower(groupValues) == "true" {
        log.Debug(fmt.Sprintf("grouping periods by day"))
        ctx.JSON(200, gin.H{"success": true, "http_code": 200, "data": groupPeriodsByDay(data.WorkPeriods, start, end.Add(time.Hour * 24))})
    } else {
        ctx.JSON(200, gin.H{"success": true, "http_code": 200, "data": data.WorkPeriods})
    }
}

func getUserAnalysisHandler(ctx *gin.Context) {
    user := getUser(ctx)
    log.Debug(fmt.Sprintf("received analysis request for user %s", user))
    results, err := analyzeUserTasks(user)
    if err != nil {
        log.Error(fmt.Errorf("unable to analyse user tasks: %v", err))
        StandardHTTP.InternalServerError(ctx)
        return
    }
    ctx.JSON(200, gin.H{"success": true, "http_code": 200, "payload": results})
}

func getUserTimeRangeAnalysisHandler(ctx *gin.Context) {
    user := getUser(ctx)
    start, end, err := parseTimestamps(ctx.Param("start"), ctx.Param("end"))
    if err != nil {
        log.Error(fmt.Errorf("unable to parse timestamps: %v", err))
        StandardHTTP.InvalidRequest(ctx)
        return
    }

    log.Debug(fmt.Sprintf("received time range analysis request for user %s", user))
    results, err := analyseRangedUserTasks(user, start, end)
    if err != nil {
        log.Error(fmt.Errorf("unable to analyse user tasks: %v", err))
        StandardHTTP.InternalServerError(ctx)
        return
    }
    ctx.JSON(200, gin.H{"success": true, "http_code": 200, "payload": results})
}

func createWorkPeriodHandler(ctx *gin.Context) {
    user := getUser(ctx)
    log.Debug(fmt.Sprintf("received request to create new work period for user %s", user))
    // create new work period in database
    period, err := persistence.createWorkPeriod(user)
    if err != nil {
        log.Error(fmt.Errorf("unable to create new work period for user %s: %v", user, err))
        StandardHTTP.InternalServerError(ctx)
        return
    }
    ctx.JSON(200, gin.H{"success": true, "http_code": 200, "payload": period})
}

func createBreakPeriodHandler(ctx *gin.Context) {
    user := getUser(ctx)
    periodId, err := uuid.Parse(ctx.Param("periodId"))
    if err != nil {
        log.Error(fmt.Sprintf("received invalid period ID"))
        StandardHTTP.InvalidRequest(ctx)
        return
    }

    log.Debug(fmt.Sprintf("received request to create new bread period for user %s", user))
    // create new work period in database
    id, err := persistence.createBreakPeriod(periodId)
    if err != nil {
        log.Error(fmt.Errorf("unable to create new break period for user %s: %v", user, err))
        StandardHTTP.InternalServerError(ctx)
        return
    }
    ctx.JSON(200, gin.H{"success": true, "http_code": 200, "id": id})
}

func endWorkPeriodHandler(ctx *gin.Context) {
    periodId, err := uuid.Parse(ctx.Param("periodId"))
    if err != nil {
        log.Error(fmt.Sprintf("received invalid period ID"))
        StandardHTTP.InvalidRequest(ctx)
        return
    }
    // check it work period exist in database
    exists, err := isValidWorkPeriod(periodId)
    if err != nil {
        log.Error(fmt.Errorf("unable to get work period %s: %v", periodId, err))
        StandardHTTP.InternalServerError(ctx)
        return
    }
    if !exists {
        log.Error(fmt.Errorf("invalid work period %s", periodId))
        StandardHTTP.NotFound(ctx)
        return
    }

    log.Debug(fmt.Sprintf("received request to end work period %s", periodId))
    err = persistence.closeWorkPeriod(periodId)
    if err != nil {
        log.Error(fmt.Errorf("unable to close work period %s", periodId))
        StandardHTTP.InternalServerError(ctx)
        return
    }
    ctx.JSON(200, gin.H{"success": true, "http_code": 200, "message": fmt.Sprintf("successfully closed work period %s", periodId)})
}

func endBreakPeriodHandler(ctx *gin.Context) {
    breakId, err := uuid.Parse(ctx.Param("breakId"))
    if err != nil {
        log.Error(fmt.Sprintf("received invalid break ID"))
        StandardHTTP.InvalidRequest(ctx)
        return
    }
    // check if break period exist in database
    exists, err := isValidBreakPeriod(breakId)
    if err != nil {
        log.Error(fmt.Errorf("unable to get work period %s: %v", breakId, err))
        StandardHTTP.InternalServerError(ctx)
        return
    }
    if !exists {
        log.Error(fmt.Errorf("invalid work period %s", breakId))
        StandardHTTP.NotFound(ctx)
        return
    }

    log.Debug(fmt.Sprintf("received request to end break period %s", breakId))
    err = persistence.closeBreakPeriod(breakId)
    if err != nil {
        log.Error(fmt.Errorf("unable to close work period %s", breakId))
        StandardHTTP.InternalServerError(ctx)
        return
    }
    ctx.JSON(200, gin.H{"success": true, "http_code": 200, "message": fmt.Sprintf("successfully closed work period %s", breakId)})
}

