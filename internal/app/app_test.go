package app

import (
	"bytes"
	"database/sql"
	"net/http"
	"net/http/cookiejar"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	appcfg "github.com/tsatke/verbose-broccoli/internal/app/config"
	"golang.org/x/net/nettest"
)

func TestAppSuite(t *testing.T) {
	gin.SetMode(gin.TestMode)
	suite.Run(t, new(AppSuite))
}

type AppSuite struct {
	suite.Suite

	app     *App
	cookies *cookiejar.Jar
}

func (suite *AppSuite) SetupTest() {
	suite.cookies, _ = cookiejar.New(nil)

	lis, err := nettest.NewLocalListener("tcp")
	suite.NoError(err)

	log := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).
		With().
		Timestamp().
		Logger()

	var opts []Option
	opts = append(opts, WithLogger(log))

	pgHost := os.Getenv("PG_HOST")
	if pgHost != "" {
		suite.T().Logf("using database at %v", pgHost)

		vp := viper.New()
		vp.Set(appcfg.PGEndpoint, pgHost)
		vp.Set(appcfg.PGPort, 5432)
		vp.Set(appcfg.PGDatabase, "postgres")
		vp.Set(appcfg.PGUsername, "postgres")
		vp.Set(appcfg.PGPassword, "postgres")

		dbProvider, err := NewPostgresDatabaseProvider(log, appcfg.Config{Viper: vp}, false)
		suite.NoError(err)
		suite.Require().NoError(dbProvider.tx(func(tx *sql.Tx) error {
			_, err := tx.Exec(`
DELETE FROM au_document_acls;
DELETE FROM au_document_headers;
`)
			return err
		}))

		opts = append(opts, WithDocumentRepo(NewPostgresDocumentRepo(dbProvider)))
	}

	suite.app = New(lis, opts...)
	suite.IsType(&MemObjectStorage{}, suite.app.objects)
	suite.IsType(&MemAuthService{}, suite.app.auth)

	go func() {
		if err := suite.app.Run(); err != nil {
			panic(err)
		}
	}()
}

func (suite *AppSuite) TearDownTest() {
	if suite.app != nil {
		suite.NoError(suite.app.Close())
	}
}

func (suite *AppSuite) login() string {
	user := uuid.New().String()
	pass := uuid.New().String()

	suite.createUser(user, pass)

	suite.
		Request("POST", "/auth/login").
		BodyJSON(M{
			"username": user,
			"password": pass,
		}).
		ExpectJSON(http.StatusOK, M{
			"success": true,
		})

	return user
}

func (suite *AppSuite) logout() {
	suite.
		Request("GET", "/auth/logout").
		ExpectJSON(http.StatusOK, M{
			"success": true,
		})
}

func (suite *AppSuite) createUser(user, pass string) {
	suite.app.auth.(*MemAuthService).data[user] = pass
}

func (suite *AppSuite) createContent(id string, content []byte) {
	suite.NoError(suite.app.objects.(*MemObjectStorage).Create(DocID(id), bytes.NewReader(content)))
}
