package deploy

import (
	"context"
	"github.com/jerson/deployer/pkg/deployer"
	"github.com/sirupsen/logrus"
	"time"
)

func delay(ctx context.Context, times int) (context.Context, error) {
	log := ctx.Value(deployer.ContextKeyLog).(*logrus.Entry)

	i := 0
	for {

		log.Info("wait ", times-i)
		time.Sleep(time.Second)

		i++
		if i >= times {
			return ctx, nil
		}
	}
}
