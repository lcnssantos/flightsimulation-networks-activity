package app

import (
	"context"
	"time"

	"github.com/lcnssantos/online-activity/internal/infra/concurrency"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/lcnssantos/online-activity/internal/domain"
)

type AppService struct {
	ivaoFacade   domain.Facade
	posconFacade domain.Facade
	vatsimFacade domain.Facade
	repository   domain.Repository
}

func NewAppService(
	ivaoFacade domain.Facade,
	posconFacade domain.Facade,
	vatsimFacade domain.Facade,
	repo domain.Repository,
) AppService {
	return AppService{
		ivaoFacade:   ivaoFacade,
		posconFacade: posconFacade,
		vatsimFacade: vatsimFacade,
		repository:   repo,
	}
}

func (a *AppService) loadActivity(ctx context.Context) (ivao *domain.Activity, vatsim *domain.Activity, poscon *domain.Activity, err error) {
	asyncTasks := concurrency.ExecuteConcurrentTasks(
		concurrency.TaskInput{
			Task: func() (interface{}, error) {
				return a.ivaoFacade.GetActivity(ctx)
			},
			Tag: "ivao-activity",
		}, concurrency.TaskInput{
			Task: func() (interface{}, error) {
				return a.posconFacade.GetActivity(ctx)
			},
			Tag: "poscon-activity",
		}, concurrency.TaskInput{
			Task: func() (interface{}, error) {
				return a.vatsimFacade.GetActivity(ctx)
			},
			Tag: "vatsim-activity",
		})

	ivaoTask := asyncTasks[0]
	posconTask := asyncTasks[1]
	vatsimTask := asyncTasks[2]

	if ivaoTask.Err != nil {
		ivao = &domain.Activity{Pilot: 0, ATC: 0}
	} else {
		ivao = ivaoTask.Result.(*domain.Activity)
	}

	if posconTask.Err != nil {
		poscon = &domain.Activity{Pilot: 0, ATC: 0}
	} else {
		poscon = posconTask.Result.(*domain.Activity)
	}

	if vatsimTask.Err != nil {
		vatsim = &domain.Activity{Pilot: 0, ATC: 0}
	} else {
		vatsim = vatsimTask.Result.(*domain.Activity)
	}

	return
}

func (a *AppService) GetActivity(ctx context.Context) (*domain.NetworkActivity, error) {
	ivao, vatsim, poscon, err := a.loadActivity(ctx)

	if err != nil {
		return nil, err
	}

	return &domain.NetworkActivity{
		Date:   time.Now(),
		IVAO:   *ivao,
		VATSIM: *vatsim,
		POSCON: *poscon,
	}, nil
}

func (a *AppService) SaveActivity(ctx context.Context) error {
	ivao, vatsim, poscon, err := a.loadActivity(ctx)

	if err != nil {
		return err
	}

	return a.repository.SaveActivity(ctx, domain.NetworkActivity{
		ID:     primitive.NewObjectID(),
		Date:   time.Now(),
		IVAO:   *ivao,
		VATSIM: *vatsim,
		POSCON: *poscon,
	})
}

func (a *AppService) GetHistoryByMinutes(ctx context.Context, minutes int64) ([]domain.NetworkActivity, error) {
	return a.repository.GetActivityByMinutes(ctx, minutes)
}

func (a *AppService) loadBrazilActivity(ctx context.Context) (ivao *domain.Activity, vatsim *domain.Activity, poscon *domain.Activity, err error) {
	asyncTasks := concurrency.ExecuteConcurrentTasks(
		concurrency.TaskInput{
			Task: func() (interface{}, error) {
				return a.ivaoFacade.GetBrazilActivity(ctx)
			},
			Tag: "ivao-brazil-activity",
		}, concurrency.TaskInput{Task: func() (interface{}, error) {
			return a.posconFacade.GetBrazilActivity(ctx)
		}, Tag: "poscon-brazil-activity"}, concurrency.TaskInput{
			Task: func() (interface{}, error) {
				return a.vatsimFacade.GetBrazilActivity(ctx)
			},
			Tag: "vatsim-brazil-activity",
		})

	ivaoTask := asyncTasks[0]
	posconTask := asyncTasks[1]
	vatsimTask := asyncTasks[2]

	if ivaoTask.Err != nil {
		ivao = &domain.Activity{Pilot: 0, ATC: 0}
	} else {
		ivao = ivaoTask.Result.(*domain.Activity)
	}

	if posconTask.Err != nil {
		poscon = &domain.Activity{Pilot: 0, ATC: 0}
	} else {
		poscon = posconTask.Result.(*domain.Activity)
	}

	if vatsimTask.Err != nil {
		vatsim = &domain.Activity{Pilot: 0, ATC: 0}
	} else {
		vatsim = vatsimTask.Result.(*domain.Activity)
	}

	return
}

func (a *AppService) GetBrazilActivity(ctx context.Context) (*domain.NetworkActivity, error) {
	ivao, vatsim, poscon, err := a.loadBrazilActivity(ctx)

	if err != nil {
		return nil, err
	}

	return &domain.NetworkActivity{
		Date:   time.Now(),
		IVAO:   *ivao,
		VATSIM: *vatsim,
		POSCON: *poscon,
	}, nil
}

func (a *AppService) SaveBrazilActivity(ctx context.Context) error {
	ivao, vatsim, poscon, err := a.loadBrazilActivity(ctx)

	if err != nil {
		return err
	}

	return a.repository.SaveBrazilActivity(ctx, domain.NetworkActivity{
		ID:     primitive.NewObjectID(),
		Date:   time.Now(),
		IVAO:   *ivao,
		VATSIM: *vatsim,
		POSCON: *poscon,
	})
}

func (a *AppService) GetBrazilHistoryByMinutes(ctx context.Context, minutes int64) ([]domain.NetworkActivity, error) {
	return a.repository.GetBrazilActivityByMinutes(ctx, minutes)
}

func (a *AppService) loadGeoActivity(ctx context.Context) (ivao *domain.GeoActivity, vatsim *domain.GeoActivity, poscon *domain.GeoActivity, err error) {
	asyncTasks := concurrency.ExecuteConcurrentTasks(
		concurrency.TaskInput{
			Task: func() (interface{}, error) {
				return a.ivaoFacade.GetGeoActivity(ctx)
			},
			Tag: "ivao-geo-activity",
		}, concurrency.TaskInput{
			Task: func() (interface{}, error) {
				return a.posconFacade.GetGeoActivity(ctx)
			},
			Tag: "poscon-geo-activity",
		}, concurrency.TaskInput{
			Task: func() (interface{}, error) {
				return a.vatsimFacade.GetGeoActivity(ctx)
			},
			Tag: "vatsim-geo-activity",
		})

	ivaoTask := asyncTasks[0]
	posconTask := asyncTasks[1]
	vatsimTask := asyncTasks[2]

	if ivaoTask.Err != nil {
		ivao = &domain.GeoActivity{}
	} else {
		ivao = ivaoTask.Result.(*domain.GeoActivity)
	}

	if posconTask.Err != nil {
		poscon = &domain.GeoActivity{}
	} else {
		poscon = posconTask.Result.(*domain.GeoActivity)
	}

	if vatsimTask.Err != nil {
		vatsim = &domain.GeoActivity{}
	} else {
		vatsim = vatsimTask.Result.(*domain.GeoActivity)
	}

	return
}

func (a *AppService) GetGeoActivity(ctx context.Context) (*domain.GeoNetworkActivity, error) {
	ivao, vatsim, poscon, err := a.loadGeoActivity(ctx)

	if err != nil {
		return nil, err
	}

	return &domain.GeoNetworkActivity{
		Date:   time.Now(),
		IVAO:   *ivao,
		VATSIM: *vatsim,
		POSCON: *poscon,
	}, nil
}

func (a *AppService) SaveGeoActivity(ctx context.Context) error {
	ivao, vatsim, poscon, err := a.loadGeoActivity(ctx)

	if err != nil {
		return err
	}

	return a.repository.SaveGeoActivity(ctx, domain.GeoNetworkActivity{
		ID:     primitive.NewObjectID(),
		Date:   time.Now(),
		IVAO:   *ivao,
		VATSIM: *vatsim,
		POSCON: *poscon,
	})
}

func (a *AppService) GetGeoHistoryByMinutes(ctx context.Context, minutes int64) ([]domain.GeoNetworkActivity, error) {
	return a.repository.GetGeoActivityByMinutes(ctx, minutes)
}
