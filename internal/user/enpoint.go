package user

import (
	"context"

	"github.com/Jhon-2801/courses-meta/meta"
	responses "github.com/Jhon-2801/lib-response/response"
)

type (
	Controller func(ctx context.Context, request interface{}) (response interface{}, err error)

	Endpoints struct {
		Create Controller
		Get    Controller
		GetAll Controller
		Update Controller
		Delete Controller
	}

	CreateReq struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		Phone     string `json:"phone"`
	}

	GetReq struct {
		ID string
	}

	GetAllReq struct {
		FirstName string
		LastName  string
		Limit     int
		Page      int
	}

	UpdateReq struct {
		ID        string
		FirstName *string `json:"first_name"`
		LastName  *string `json:"last_name"`
		Email     *string `json:"email"`
		Phone     *string `json:"phone"`
	}
	DeleteReq struct {
		ID string
	}
	Response struct {
		Status int         `json:"status"`
		Data   interface{} `json:"data,omitempty"`
		Err    string      `json:"error,omitempty"`
		Meta   *meta.Meta  `json:"meta,omitempty"`
	}

	Config struct {
		LimPageDef string
	}
)

func MakeEndpoints(s Service, config Config) Endpoints {

	return Endpoints{
		Create: makeCreateEndpoint(s),
		Get:    makeGetEndpoint(s),
		GetAll: makeGetAllEndpoint(s, config),
		Update: makeUpdateEndpoint(s),
		Delete: makeDeleteEndpoint(s),
	}
}
func makeCreateEndpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(CreateReq)

		if req.FirstName == "" {
			return nil, responses.BadRequest("first name is requerid")
		}
		if req.LastName == "" {
			return nil, responses.BadRequest("last name is requerid")
		}

		user, err := s.Create(ctx, req.FirstName, req.LastName, req.Email, req.Phone)
		if err != nil {
			return nil, responses.InternalServerError(err.Error())
		}
		return responses.Created("succes", user, nil), nil
	}
}

func makeGetEndpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {

		req := request.(GetReq)

		user, err := s.Get(ctx, req.ID)

		if err != nil {
			return nil, responses.NotFound(err.Error())
		}

		return responses.OK("success", user, nil), nil
	}
}

func makeGetAllEndpoint(s Service, config Config) Controller {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {

		req := request.(GetAllReq)

		filters := Filters{
			FirstName: req.FirstName,
			LastName:  req.LastName,
		}
		count, err := s.Count(ctx, filters)
		if err != nil {
			return nil, responses.InternalServerError(err.Error())
		}
		meta, err := meta.New(req.Page, req.Limit, count, config.LimPageDef)
		if err != nil {
			return nil, responses.InternalServerError(err.Error())
		}
		users, err := s.GetAll(ctx, filters, meta.Offset(), meta.Limit())

		if err != nil {
			return nil, responses.InternalServerError(err.Error())
		}
		return responses.OK("success", users, meta), nil
	}
}

func makeUpdateEndpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(UpdateReq)

		if req.FirstName == nil || *req.FirstName == "" {
			return nil, responses.BadRequest("first name is required")
		}
		if req.LastName == nil || *req.LastName == "" {
			return nil, responses.BadRequest("last name is required")
		}

		err = s.Update(ctx, req.ID, req.FirstName, req.LastName, req.Email, req.Phone)
		if err != nil {
			return nil, responses.InternalServerError(err.Error())
		}
		return responses.OK("success", nil, nil), nil
	}
}

func makeDeleteEndpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(DeleteReq)

		if err := s.Delete(ctx, req.ID); err != nil {
			return nil, responses.InternalServerError(err.Error())
		}
		return responses.OK("success", nil, nil), nil
	}
}
