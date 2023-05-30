package httpserver

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"github.com/lcnssantos/online-activity/internal/app"
	"github.com/lcnssantos/online-activity/internal/infra/concurrency"
)

type Controller struct {
	appService app.AppService
}

func NewController(appService app.AppService) Controller {
	return Controller{appService: appService}
}

func (t *Controller) GetActivity(ctx *gin.Context) {
	activity, err := t.appService.GetActivity(ctx)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, activity)
}

func (t *Controller) GetBrazilActivity(ctx *gin.Context) {
	brazilActivity, err := t.appService.GetBrazilActivity(ctx)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, brazilActivity)
}

// func (t *Controller) GetGeoActivity(ctx *gin.Context) {
// 	geoActivity, err := t.appService.GetGeoActivity(ctx)

// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, err)
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, geoActivity)
// }

func (t *Controller) Get24hHistory(ctx *gin.Context) {
	activities, err := t.appService.GetHistoryByMinutes(ctx, 24*60)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, activities)
}

func (t *Controller) GetHistoryByMinutes(ctx *gin.Context) {
	minutes, err := strconv.Atoi(ctx.Param("minutes"))

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	activities, err := t.appService.GetHistoryByMinutes(ctx, int64(minutes))

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, activities)
}

func (t *Controller) GetBrazil24hHistory(ctx *gin.Context) {
	activities, err := t.appService.GetBrazilHistoryByMinutes(ctx, 24*60)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, activities)
}

func (t *Controller) GetBrazilHistoryByMinutes(ctx *gin.Context) {
	minutes, err := strconv.Atoi(ctx.Param("minutes"))

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	activities, err := t.appService.GetBrazilHistoryByMinutes(ctx, int64(minutes))

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, activities)
}

// func (t *Controller) GetGeo24hHistory(ctx *gin.Context) {
// 	activities, err := t.appService.GetGeoHistoryByMinutes(ctx, 24*60)

// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, err)
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, activities)
// }

// func (t *Controller) GetGeoHistoryByMinutes(ctx *gin.Context) {
// 	minutes, err := strconv.Atoi(ctx.Param("minutes"))

// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, err)
// 		return
// 	}

// 	activities, err := t.appService.GetGeoHistoryByMinutes(ctx, int64(minutes))

// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, err)
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, activities)
// }

func (t *Controller) saveCurrent(ctx *gin.Context) {
	asyncTasks := concurrency.ExecuteConcurrentTasks(
		concurrency.TaskInput{
			Task: func() (interface{}, error) {
				err := t.appService.SaveActivity(ctx)
				return nil, err
			},
			Tag: "save-activity",
		},
		concurrency.TaskInput{
			Task: func() (interface{}, error) {
				err := t.appService.SaveBrazilActivity(ctx)
				return nil, err
			},
			Tag: "save-brazil-activity",
		},
		concurrency.TaskInput{
			Task: func() (interface{}, error) {
				err := t.appService.SaveGeoActivity(ctx)
				return nil, err
			},
			Tag: "save-geo-activity",
		},
	)

	stdTask := asyncTasks[0]
	brTask := asyncTasks[1]
	geoTask := asyncTasks[2]

	if brTask.Err != nil {
		log.Error().Err(brTask.Err).Msg("BR_TASK")
	}

	if geoTask.Err != nil {
		log.Error().Err(geoTask.Err).Msg("GEO_TASK")
	}

	if stdTask.Err != nil {
		log.Error().Err(stdTask.Err).Msg("STD_TASK")
	}
}
