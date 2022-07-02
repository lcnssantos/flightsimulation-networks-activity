package app

import (
	"context"
	"time"

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
) *AppService {
	return &AppService{
		ivaoFacade:   ivaoFacade,
		posconFacade: posconFacade,
		vatsimFacade: vatsimFacade,
		repository:   repo,
	}
}

func (a *AppService) loadActivity(ctx context.Context) (ivao *domain.Activity, vatsim *domain.Activity, poscon *domain.Activity, err error) {
	ivao, err = a.ivaoFacade.GetActivity(ctx)

	if err != nil {
		return
	}

	poscon, err = a.posconFacade.GetActivity(ctx)

	if err != nil {
		return
	}

	vatsim, err = a.vatsimFacade.GetActivity(ctx)

	if err != nil {
		return
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
	ivao, err = a.ivaoFacade.GetBrazilActivity(ctx)

	if err != nil {
		return
	}

	poscon, err = a.posconFacade.GetBrazilActivity(ctx)

	if err != nil {
		return
	}

	vatsim, err = a.vatsimFacade.GetBrazilActivity(ctx)

	if err != nil {
		return
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
	ivao, err = a.ivaoFacade.GetGeoActivity(ctx)

	if err != nil {
		return
	}

	poscon, err = a.posconFacade.GetGeoActivity(ctx)

	if err != nil {
		return
	}

	vatsim, err = a.vatsimFacade.GetGeoActivity(ctx)

	if err != nil {
		return
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
		Date:   time.Now(),
		IVAO:   *ivao,
		VATSIM: *vatsim,
		POSCON: *poscon,
	})
}

func (a *AppService) GetGeoHistoryByMinutes(ctx context.Context, minutes int64) ([]domain.GeoNetworkActivity, error) {
	return a.repository.GetGeoActivityByMinutes(ctx, minutes)
}
