package data

import (
	"context"

	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/sm"

	"github.com/tencat-dev/go-base/internal/biz"
	"github.com/tencat-dev/go-base/internal/infra/persistence/postgres/bob/models"

	"github.com/go-kratos/kratos/v2/log"
)

type authRepo struct {
	data *Data
	log  *log.Helper
}

// NewAuthRepo .
func NewAuthRepo(data *Data, logger *log.Helper) biz.AuthRepo {
	return &authRepo{
		data: data,
		log:  logger,
	}
}

func (r *authRepo) FindByEmail(ctx context.Context, email string) (*biz.Auth, error) {
	auth, err := models.Users.Query(
		sm.Where(models.Users.Columns.Email.EQ(psql.Arg(email))),
	).One(ctx, r.data.db)
	if err != nil {
		return nil, err
	}

	return &biz.Auth{
		ID:           auth.ID,
		Name:         auth.Name,
		Email:        auth.Email,
		PasswordHash: auth.PasswordHash,
	}, nil
}
