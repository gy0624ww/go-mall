package appservice

import (
	"context"

	"github.com/go-study-lab/go-mall/logic/domainservice"
)

type DemoAppSvc struct {
	ctx           context.Context
	demoDomainSvc *domainservice.DemoDomainSvc
}

func NewDemoAppSvc(ctx context.Context) *DemoAppSvc {
	return &DemoAppSvc{
		ctx:           ctx,
		demoDomainSvc: domainservice.NewDemoDomainSvc(ctx),
	}
}

//func (das *DemoAppSvc)DoSomething() {
//	demo, err := das.demoDomainSvc.GetDemoEntity(id)
//	if err != nil {
//		logger.New(das.ctx).Error("DemoAppSvc DoSomething err", err)
//		return err
//	}
//	......
//}

func (das *DemoAppSvc) GetDemoIdentities() ([]int64, error) {
	demos, err := das.demoDomainSvc.GetDemos()
	if err != nil {
		return nil, err
	}
	identities := make([]int64, 0, len(demos))

	for _, demo := range demos {
		identities = append(identities, demo.Id)
	}
	return identities, nil
}
