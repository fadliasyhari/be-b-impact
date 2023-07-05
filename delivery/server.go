package delivery

import (
	"fmt"
	"net/http"

	"be-b-impact.com/csr/config"
	"be-b-impact.com/csr/delivery/api/middleware"
	"be-b-impact.com/csr/delivery/controller"
	"be-b-impact.com/csr/manager"
	"be-b-impact.com/csr/model"
	"be-b-impact.com/csr/usecase"
	"be-b-impact.com/csr/utils/authenticator"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

type Server struct {
	ucManager    manager.UseCaseManager
	engine       *gin.Engine
	host         string
	log          *logrus.Logger
	authUseCase  usecase.AuthUseCase
	tokenService authenticator.AccessToken
}

func (s *Server) initController() {
	s.engine.Use(middleware.LogRequestMiddleware(s.log))
	tokenMdw := middleware.NewTokenValidator(s.tokenService)
	controller.NewUsersController(s.engine, s.ucManager.UsersUseCase(), tokenMdw)
	controller.NewAuthController(s.engine, s.ucManager.UsersUseCase(), s.authUseCase, tokenMdw)
	controller.NewCategoryController(s.engine, s.ucManager.CategoryUseCase(), tokenMdw)
	controller.NewTagController(s.engine, s.ucManager.TagUseCase(), tokenMdw)
	controller.NewContentController(s.engine, s.ucManager.ContentUseCase(), tokenMdw)
	controller.NewProposalController(s.engine, s.ucManager.ProposalUseCase(), s.ucManager.ProposalDetailUseCase(), s.ucManager.FileUseCase(), s.ucManager.ProgressUseCase(), s.ucManager.ProposalProgressUseCase(), s.ucManager.UsersUseCase(), tokenMdw)
	controller.NewProgressController(s.engine, s.ucManager.ProgressUseCase(), tokenMdw)
	controller.NewEventController(s.engine, s.ucManager.EventUseCase(), s.ucManager.EventParticipantUseCase(), tokenMdw)
}

func (s *Server) Run() {
	s.initController()
	err := s.engine.Run(s.host)
	if err != nil {
		panic(err.Error())
	}
}

func NewServer() *Server {
	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	infra, err := manager.NewInfraManager(cfg)
	if err != nil {
		panic(err)
	}
	repo := manager.NewRepositoryManager(infra)
	uc := manager.NewUseCaseManager(repo)

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Replace with your desired origins
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	opt, err := redis.ParseURL(cfg.RedisConfig.Url)
	if err != nil {
		panic(err)
	}
	client := redis.NewClient(opt)
	tokenService := authenticator.NewTokenService(cfg.TokenConfig, client)
	authUc := usecase.NewAuthUseCase(tokenService, repo.UsersRepo(), client)
	r.GET("/migrate", func(c *gin.Context) {
		infra.Conn().DB()
		err := infra.Migrate(
			&model.User{},
			&model.Profile{},
			&model.Category{},
			&model.Image{},
			&model.Content{},
			&model.Tag{},
			&model.TagsContent{},
			&model.ContentDetail{},
			&model.Proposal{},
			&model.ProposalDetail{},
			&model.File{},
			&model.EventImage{},
			&model.Event{},
			&model.EventParticipant{},
			&model.Progress{},
			&model.ProposalProgress{},
		)
		if err != nil {
			panic(err)
		}
		c.JSON(http.StatusOK, gin.H{"message": "Migration completed"})
	})

	return &Server{
		ucManager:    uc,
		engine:       r,
		host:         fmt.Sprintf("%s:%s", cfg.ApiHost, cfg.ApiPort),
		log:          infra.Log(),
		authUseCase:  authUc,
		tokenService: tokenService,
	}
}
