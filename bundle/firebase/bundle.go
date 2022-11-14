package firebase

import (
	"context"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"

	"slice"
)

type Bundle struct {
	ProjectID       string `envconfig:"FIREBASE_PROJECT_ID" required:"True"`
	Authentication  bool
	CredentialsFile string `envconfig:"FIREBASE_CREDENTIALS_FILE"`
	EmulatorHost    string `envconfig:"FIREBASE_AUTH_EMULATOR_HOST"`
}

// Build provide database to di.
func (b *Bundle) Build(builder slice.ContainerBuilder) {
	builder.Provide(b.NewFirebaseApp)

	if b.Authentication {
		builder.Provide(b.NewFirebaseAuth)
	}
}

// Boot implements Bundle interface.
func (b *Bundle) Boot(_ context.Context, interactor slice.Container) (err error) {
	return
}

// Shutdown implements Bundle interface.
func (b *Bundle) Shutdown(_ context.Context, _ slice.Container) (err error) {
	return nil
}

func (b *Bundle) NewFirebaseApp(logger slice.Logger) (*firebase.App, error) {
	defer logger.Infof("firebase", "Register firebase app")

	opts := make([]option.ClientOption, 0)

	if b.CredentialsFile != "" {
		opts = append(opts, option.WithCredentialsFile(b.CredentialsFile))
	}

	app, _ := firebase.NewApp(
		context.Background(),
		&firebase.Config{
			ProjectID: b.ProjectID,
		},
		opts...)
	return app, nil
}

func (b *Bundle) NewFirebaseAuth(logger slice.Logger, app *firebase.App) (*auth.Client, error) {
	defer logger.Infof("firebase", "Register firebase auth")

	return app.Auth(context.Background())
}
