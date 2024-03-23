package service

import (
	"context"

	"golang_socmed/internal/entity"
	helpers "golang_socmed/internal/helper"
	"golang_socmed/internal/model"
	"golang_socmed/internal/model/converter"
	"golang_socmed/internal/repository"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type PostService struct {
	Repository *repository.PostRepository
	Validate   *validator.Validate
	Log        *logrus.Logger
}

func NewPostService(r *repository.PostRepository, validate *validator.Validate, log *logrus.Logger) *PostService {
	return &PostService{Repository: r, Validate: validate, Log: log}
}

func (s *PostService) List(ctx context.Context, filter *model.PostFilter) ([]model.PostResponse, error) {

	if err := helpers.ValidationError(s.Validate, filter); err != nil {
		s.Log.WithError(err).Error("failed to validate request query params")
		return nil, &fiber.Error{
			Code:    400,
			Message: err.Error(),
		}
	}

	posts, err := s.Repository.List(filter)
	if err != nil {
		s.Log.WithError(err).Error("failed get post lists")
		return nil, err
	}

	newPosts := make([]model.PostResponse, len(posts))
	for i, post := range posts {
		newPosts[i] = *converter.PostConverter(&post)
	}

	return newPosts, nil
}

func (s *PostService) Create(ctx context.Context, request *model.PostRequest, userId string) error {
	// if err := s.Validate.Struct(request); err != nil {
	if err := helpers.ValidationError(s.Validate, request); err != nil {
		s.Log.Error("failed to validate request body")
		return err
	}

	newRequest := &entity.Post{
		ID:         uuid.New().String(),
		PostInHtml: request.PostInHtml,
		Tags:       request.Tags,
		UserId:     userId,
	}

	err := s.Repository.Create(newRequest)
	if err != nil {
		//s.Log.Error("failed to insert new data")
		return err
	}

	return nil
}
