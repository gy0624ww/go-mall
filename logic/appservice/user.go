package appservice

import (
	"context"

	"github.com/go-study-lab/go-mall/api/reply"
	"github.com/go-study-lab/go-mall/common/logger"
	"github.com/go-study-lab/go-mall/common/util"
	"github.com/go-study-lab/go-mall/logic/domainservice"
)

type UserAppSvc struct {
	ctx           context.Context
	userDomainSvc *domainservice.UserDomainSvc
}

func NewUserAppSvc(ctx context.Context) *UserAppSvc {
	return &UserAppSvc{
		ctx:           ctx,
		userDomainSvc: domainservice.NewUserDomainSvc(ctx),
	}
}
func (us *UserAppSvc) GenToken() (*reply.TokenReply, error) {
	token, err := us.userDomainSvc.GenAuthToken(12345678, "h5", "")
	if err != nil {
		return nil, err
	}
	logger.Info(us.ctx, "generate token success", "tokenData", token)
	tokenReply := new(reply.TokenReply)
	util.CopyProperties(tokenReply, token)
	return tokenReply, err
}
