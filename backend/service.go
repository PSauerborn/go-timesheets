package main

import (
    "fmt"
    "time"
    "strconv"
    "strings"
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "github.com/jackc/pgx/v4"
    log "github.com/sirupsen/logrus"
    jaeger "github.com/PSauerborn/jaeger-negroni"
)

var (

)

type UserIDMetric struct {}

func(metric UserIDMetric) MetricName() string {
    return "uid"
}

func(metric UserIDMetric) EvaluateMetric(ctx *gin.Context) interface{} {
    user := getUser(ctx)
    log.Debug(fmt.Sprintf("setting jaeger spans with user ID %s", user))
    return user
}

func main() {
    // configure environment variables and connect persistence layer to database
    ConfigureService()
    ConnectPersistence()

    // create new jaeger config and add uid metric
    config := jaeger.Config("jaeger-agent", "go-timesheets-api", 6831)
    config.PreRequestMetrics = append(config.PreRequestMetrics, UserIDMetric{})

    tracer := jaeger.NewTracer(config)
    defer tracer.Close()

    router := gin.New()
    router.Use(jaeger.JaegerNegroni(config))

    // create handlers for user data routes
    router.GET("/go-timesheets/health", healthCheckHandler)
    router.GET("/go-timesheets/active", getActivePeriodHandler)
    router.GET("/go-timesheets/data", getUserDataHandler)
    router.GET("/go-timesheets/data/:start/:end", getUserTimeRangeDataHandler)
    router.GET("/go-timesheets/bucket_analysis/:start/:end", getUserBucketAnalysisHandler)

    // create handlers for user data analysis routes
    router.GET("/go-timesheets/analyse", getUserAnalysisHandler)
    router.GET("/go-timesheets/analyse/:start/:end", getUserTimeRangeAnalysisHandler)
    // create handlers to create work and break periods
    router.POST("/go-timesheets/work_period", createWorkPeriodHandler)
    router.POST("/go-timesheets/break_period/:periodId", createBreakPeriodHandler)
    // create handlers to end work and break periods
    router.PATCH("/go-timesheets/work_period/:periodId", endWorkPeriodHandler)
    router.PATCH("/go-timesheets/break_period/:breakId", endBreakPeriodHandler)

    router.Run(fmt.Sprintf(":%d", ListenPort))
}

// function used to retrieve authenticated user ID from header
func getUser(ctx *gin.Context) string {
    return ctx.Request.Header.Get("X-Authenticated-Userid")
}

// handler function used for basic health checks
func healthCheckHandler(ctx *gin.Context) {
    log.Debug("received request for health check route")
    StandardHTTP.Success(ctx)
}

// function used to retrieve current active period from database
// note that a 404 response is returned if no active period exists
func getActivePeriodHandler(ctx *gin.Context) {
    // retrieve current user and get active period
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

// function used to retrieve user data from database
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

// function used to retrieve user data from database for specific time range
func getUserTimeRangeDataHandler(ctx *gin.Context) {
    user := getUser(ctx)
    // get start and end time from url and parse into time.Time objects
    start, end, err := parseTimestamps(ctx.Param("start"), ctx.Param("end"), "2006-01-02")
    if err != nil {
        log.Error(fmt.Errorf("unable to parse timestamps: %v", err))
        StandardHTTP.InvalidRequestWithMessage(ctx, "invalid timestamp(s)")
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

// function used to return aggregated results for user data
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

// function used to return aggregated results for user over a specific time range
func getUserTimeRangeAnalysisHandler(ctx *gin.Context) {
    user := getUser(ctx)
    // get start and end time from url and parse into time.Time objects
    start, end, err := parseTimestamps(ctx.Param("start"), ctx.Param("end"), "2006-01-02")
    if err != nil {
        log.Error(fmt.Errorf("unable to parse timestamps: %v", err))
        StandardHTTP.InvalidRequestWithMessage(ctx, "invalid timestamp(s)")
        return
    }
    // analyse users tasks over time range
    log.Debug(fmt.Sprintf("received time range analysis request for user %s", user))
    results, err := analyseRangedUserTasks(user, start, end)
    if err != nil {
        log.Error(fmt.Errorf("unable to analyse user tasks: %v", err))
        StandardHTTP.InternalServerError(ctx)
        return
    }
    ctx.JSON(200, gin.H{"success": true, "http_code": 200, "payload": results})
}

func getUserBucketAnalysisHandler(ctx *gin.Context) {
    user := getUser(ctx)
    // get start and end time from url and parse into time.Time objects
    start, end, err := parseTimestamps(ctx.Param("start"), ctx.Param("end"), "2006-01-02T15:04")
    if err != nil {
        log.Error(fmt.Errorf("unable to parse timestamps: %v", err))
        StandardHTTP.InvalidRequestWithMessage(ctx, "invalid timestamp(s)")
        return
    }
    log.Debug(fmt.Sprintf("received bucket analysis request for user %s", user))
    // retrieve bucket size from query string and parse to integer
    bucketSizeString := ctx.DefaultQuery("bucket_size", "1440")
    bucketSize, err := strconv.Atoi(bucketSizeString)
    if err != nil {
        log.Error(fmt.Errorf("received invalid bucket size: %v", err))
        StandardHTTP.InvalidRequestWithMessage(ctx, "invalid bucket interval")
        return
    }
    // execute bucket analysis and return results
    results, err := executeBucketAnalysis(user, start, end, bucketSize)
    if err != nil {
        log.Error(fmt.Errorf("unable to execute bucket analysis: %v", err))
        StandardHTTP.InternalServerError(ctx)
        return
    }
    ctx.JSON(200, gin.H{"success": true, "http_code": 200, "payload": results})
}

// function used to create a new work period in the database
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

// function used to create new break period in database
func createBreakPeriodHandler(ctx *gin.Context) {
    user := getUser(ctx)
    // retrieve and parse period id from URL
    periodId, err := uuid.Parse(ctx.Param("periodId"))
    if err != nil {
        log.Error(fmt.Sprintf("received invalid period ID"))
        StandardHTTP.InvalidRequestWithMessage(ctx, "invalid period id")
        return
    }

    log.Debug(fmt.Sprintf("received request to create new bread period for user %s", user))
    // create new work period in database
    payload, err := persistence.createBreakPeriod(periodId)
    if err != nil {
        log.Error(fmt.Errorf("unable to create new break period for user %s: %v", user, err))
        StandardHTTP.InternalServerError(ctx)
        return
    }
    ctx.JSON(200, gin.H{"success": true, "http_code": 200, "payload": payload})
}

// function used to end a specific work period
func endWorkPeriodHandler(ctx *gin.Context) {
    periodId, err := uuid.Parse(ctx.Param("periodId"))
    if err != nil {
        log.Error(fmt.Sprintf("received invalid period ID"))
        StandardHTTP.InvalidRequestWithMessage(ctx, "invalid period id")
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

// function used to end a particular break period
func endBreakPeriodHandler(ctx *gin.Context) {
    breakId, err := uuid.Parse(ctx.Param("breakId"))
    if err != nil {
        log.Error(fmt.Sprintf("received invalid break ID"))
        StandardHTTP.InvalidRequestWithMessage(ctx, "invalid break id")
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

