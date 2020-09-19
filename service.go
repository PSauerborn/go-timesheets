package main

import (
    "fmt"
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    log "github.com/sirupsen/logrus"
)

var (

)

func main() {

    ConfigureService()
    ConnectPersistence()

    router := gin.New()

    router.GET("/go-timesheets/health", healthCheckHandler)
    router.GET("/go-timesheets/data", getUserDataHandler)

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

func createWorkPeriodHandler(ctx *gin.Context) {
    user := getUser(ctx)
    log.Debug(fmt.Sprintf("received request to create new work period for user %s", user))
    // create new work period in database
    id, err := persistence.createWorkPeriod(user)
    if err != nil {
        log.Error(fmt.Errorf("unable to create new work period for user %s: %v", user, err))
        StandardHTTP.InternalServerError(ctx)
        return
    }
    ctx.JSON(200, gin.H{"success": true, "http_code": 200, "id": id})
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

