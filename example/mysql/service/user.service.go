package service

import (
	"context"
	"github.com/happycrud/crud/example/mysql/api"
	"github.com/happycrud/crud/example/mysql/crud"
	"github.com/happycrud/crud/example/mysql/crud/user"
	"math"
	"strings"
	"time"

	"github.com/happycrud/xsql"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// UserServiceImpl UserServiceImpl
type UserServiceImpl struct {
	api.UnimplementedUserServiceServer
	Client *crud.Client
}

type IValidateUser interface {
	ValidateUser(a *api.User) error
}

// CreateUser CreateUser
func (s *UserServiceImpl) CreateUser(ctx context.Context, req *api.User) (*api.User, error) {
	if checker, ok := interface{}(s).(IValidateUser); ok {
		if err := checker.ValidateUser(req); err != nil {
			return nil, err
		}
	}

	a := &user.User{
		Id:    0,
		Name:  req.GetName(),
		Age:   req.GetAge(),
		Ctime: time.Now(),
		Mtime: time.Now(),
	}
	var err error
	_, err = s.Client.User.
		Create().
		SetUser(a).
		Save(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	// query after create and return
	a2, err := s.Client.Master.User.
		Find().
		Where(
			user.IdOp.EQ(a.Id),
		).
		One(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return convertUser(a2), nil
}

// DeleteUser DeleteUser
func (s *UserServiceImpl) DeleteUser(ctx context.Context, req *api.UserId) (*emptypb.Empty, error) {
	_, err := s.Client.User.
		Delete().
		Where(
			user.IdOp.EQ(req.GetId()),
		).
		Exec(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

// Updateuser UpdateUser
func (s *UserServiceImpl) UpdateUser(ctx context.Context, req *api.UpdateUserReq) (*api.User, error) {
	if checker, ok := interface{}(s).(IValidateUser); ok {
		if err := checker.ValidateUser(req.User); err != nil {
			return nil, err
		}
	}
	if len(req.GetUpdateMask()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "empty filter condition")
	}
	update := s.Client.User.Update()
	for _, v := range req.GetUpdateMask() {
		switch v {
		case user.Name:
			update.SetName(req.GetUser().GetName())
		case user.Age:
			update.SetAge(req.GetUser().GetAge())
		case user.Ctime:
			t, err := time.ParseInLocation("2006-01-02 15:04:05", req.GetUser().GetCtime(), time.Local)
			if err != nil {
				return nil, status.Error(codes.InvalidArgument, err.Error())
			}
			update.SetCtime(t)
		case user.Mtime:
			t, err := time.ParseInLocation("2006-01-02 15:04:05", req.GetUser().GetMtime(), time.Local)
			if err != nil {
				return nil, status.Error(codes.InvalidArgument, err.Error())
			}
			update.SetMtime(t)
		}
	}
	_, err := update.
		Where(
			user.IdOp.EQ(req.GetUser().GetId()),
		).
		Save(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	// query after update and return
	a, err := s.Client.Master.User.
		Find().
		Where(
			user.IdOp.EQ(req.GetUser().GetId()),
		).
		One(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return convertUser(a), nil
}

// GetUser GetUser
func (s *UserServiceImpl) GetUser(ctx context.Context, req *api.UserId) (*api.User, error) {
	a, err := s.Client.User.
		Find().
		Where(
			user.IdOp.EQ(req.GetId()),
		).
		One(ctx)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return convertUser(a), nil
}

// ListUsers ListUsers
func (s *UserServiceImpl) ListUsers(ctx context.Context, req *api.ListUsersReq) (*api.ListUsersResp, error) {
	page := req.GetPage()
	size := req.GetPageSize()
	if size <= 0 {
		size = 20
	}
	offset := size * (page - 1)
	if offset < 0 {
		offset = 0
	}
	finder := s.Client.User.
		Find().
		Offset(offset).
		Limit(size)

	if req.GetOrderByField() == api.UserField_User_unknow {
		req.OrderByField = api.UserField_User_id
	}
	odb := strings.TrimPrefix(req.GetOrderByField().String(), "User_")
	if req.GetOrderByDesc() {
		finder.OrderDesc(odb)
	} else {
		finder.OrderAsc(odb)
	}
	counter := s.Client.User.
		Find().
		Count()

	var ps []*xsql.Predicate
	for _, v := range req.GetFilters() {
		p, err := xsql.GenP(strings.TrimPrefix(v.Field.String(), "User_"), v.Op, v.Value)
		if err != nil {
			return nil, err
		}
		ps = append(ps, p)
	}
	if len(ps) > 0 {
		p := xsql.And(ps...)
		finder.WhereP(p)
		counter.WhereP(p)
	}
	list, err := finder.All(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	count, err := counter.Int64(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	pageCount := int32(math.Ceil(float64(count) / float64(size)))

	return &api.ListUsersResp{Users: convertUserList(list), TotalCount: int32(count), PageCount: pageCount}, nil
}

func convertUser(a *user.User) *api.User {
	return &api.User{
		Id:    a.Id,
		Name:  a.Name,
		Age:   a.Age,
		Ctime: a.Ctime.Format("2006-01-02 15:04:05"),
		Mtime: a.Mtime.Format("2006-01-02 15:04:05"),
	}
}

func convertUserList(list []*user.User) []*api.User {
	ret := make([]*api.User, 0, len(list))
	for _, v := range list {
		ret = append(ret, convertUser(v))
	}
	return ret
}
