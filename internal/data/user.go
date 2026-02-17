package data

import (
	"context"
	"time"

	"github.com/aarondl/opt/omit"
	"github.com/google/uuid"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/dm"
	"github.com/stephenafamo/bob/dialect/psql/sm"

	"github.com/tencat-dev/go-base/internal/biz"
	"github.com/tencat-dev/go-base/internal/infra/persistence/postgres/bob/models"

	"github.com/go-kratos/kratos/v2/log"
)

type userRepo struct {
	data *Data
	log  *log.Helper
}

// NewUserRepo .
func NewUserRepo(data *Data, logger *log.Helper) biz.UserRepo {
	return &userRepo{
		data: data,
		log:  logger,
	}
}

func (r *userRepo) Save(ctx context.Context, u *biz.User) (*biz.User, error) {
	setter := &models.UserSetter{
		Name:         omit.From(u.Name),
		Email:        omit.From(u.Email),
		PasswordHash: omit.From(u.PasswordHash),
	}

	insertedUser, err := models.Users.Insert(setter).One(ctx, r.data.db)
	if err != nil {
		return nil, err
	}

	// Update the user object with the database-assigned values
	u.ID = insertedUser.ID
	u.CreatedAt = insertedUser.CreatedAt
	u.UpdatedAt = insertedUser.UpdatedAt

	return u, nil
}

func (r *userRepo) Update(ctx context.Context, u *biz.User) (*biz.User, error) {
	setter := &models.UserSetter{
		Name:      omit.From(u.Name),
		Email:     omit.From(u.Email),
		UpdatedAt: omit.From(time.Now().UTC()),
	}

	_, err := models.Users.Update(
		models.UpdateWhere.Users.ID.EQ(u.ID),
		setter.UpdateMod(),
	).Exec(ctx, r.data.db)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (r *userRepo) FindByID(ctx context.Context, id uuid.UUID) (*biz.User, error) {
	user, err := models.FindUser(ctx, r.data.db, id)
	if err != nil {
		return nil, err
	}

	return &biz.User{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (r *userRepo) ListAll(ctx context.Context) ([]*biz.User, error) {
	userslice, err := models.Users.Query(sm.Limit(10)).All(ctx, r.data.db)
	if err != nil {
		return nil, err
	}

	var users []*biz.User
	for _, user := range userslice {
		users = append(users, &biz.User{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		})
	}

	return users, nil
}

func (r *userRepo) DeleteByID(ctx context.Context, id uuid.UUID) error {
	_, err := models.Users.Delete(
		dm.Where(models.Users.Columns.ID.EQ(psql.Arg(id))),
	).Exec(ctx, r.data.db)
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepo) ExistByID(ctx context.Context, id uuid.UUID) (bool, error) {
	exist, err := models.Users.Query(
		sm.Where(models.Users.Columns.ID.EQ(psql.Arg(id))),
	).Exists(ctx, r.data.db)
	if err != nil {
		return exist, err
	}

	return exist, nil
}
