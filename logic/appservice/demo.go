package appservice

import (
	"context"

	"github.com/go-study-lab/go-mall/api/reply"
	"github.com/go-study-lab/go-mall/api/request"
	"github.com/go-study-lab/go-mall/common/errcode"
	"github.com/go-study-lab/go-mall/common/logger"
	"github.com/go-study-lab/go-mall/common/util"
	"github.com/go-study-lab/go-mall/dal/cache"
	"github.com/go-study-lab/go-mall/logic/do"
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

func (das *DemoAppSvc) CreateDemoOrder(orderRequest *request.DemoOrderCreate) (*reply.DemoOrder, error) {
	demoOrderDo := new(do.DemoOrder)
	err := util.CopyProperties(demoOrderDo, orderRequest)
	if err != nil {
		errcode.Wrap("请求转换demoOrderDo失败", err)
		return nil, err
	}
	demoOrderDo, err = das.demoDomainSvc.CreateDemoOrder(demoOrderDo)
	if err != nil {
		return nil, err
	}
	// 设置缓存和读取, 测试功能用，无实际意义
	cache.SetDemoOrder(das.ctx, demoOrderDo)
	cacheData, _ := cache.GetDemoOrder(das.ctx, demoOrderDo.OrderNo)
	logger.Info(das.ctx, "redis data", "data", cacheData)

	replyDemoOrder := new(reply.DemoOrder)
	err = util.CopyProperties(replyDemoOrder, demoOrderDo)
	if err != nil {
		errcode.Wrap("demoOrderDo转换成replyDemoOrder失败", err)
		return nil, err
	}
	return replyDemoOrder, err
}
