package firebase

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/api/option"
)

var Module = fx.Module("gmc.firebase", fx.Provide(
	newFirebaseApp))

func newFirebaseApp(cfg *viper.Viper, logger *zap.Logger) *firebase.App {
	opt := option.WithCredentialsFile(cfg.GetString("app.firebase.credentialsFile"))
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		logger.Fatal("error initializing app", zap.Error(err))
		return nil
	}
	return app
}
