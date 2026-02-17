package authz

import (
	"context"

	"github.com/casbin/casbin/v3"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/go-kratos/kratos/v2/transport"
	jwtv5 "github.com/golang-jwt/jwt/v5"

	"github.com/tencat-dev/go-base/internal/conf"
)

type AuthzMiddleware middleware.Middleware

func NewAuthzMiddleware(
	jwtConf *conf.JWT,
	e casbin.IEnforcer,
	r *AuthzRegistry,
) AuthzMiddleware {
	jwtMiddleware := jwt.Server(
		func(*jwtv5.Token) (any, error) {
			return []byte(jwtConf.Secret), nil
		},
		jwt.WithSigningMethod(jwtv5.SigningMethodHS256),
	)

	return func(next middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (any, error) {
			tr, ok := transport.FromServerContext(ctx)
			if !ok {
				return nil, errors.Unauthorized("NO_TRANSPORT", "no transport")
			}

			fullMethod := tr.Operation()

			// 🔥 Fast path: public API
			perm, exists := r.Get(fullMethod)
			if !exists {
				return next(ctx, req)
			}

			// 🔥 Protected API → verify JWT first
			handlerWithJWT := jwtMiddleware(func(ctx context.Context, req any) (any, error) {

				token, ok := jwt.FromContext(ctx)
				if !ok {
					return nil, errors.Unauthorized("NO_USER", "no user")
				}

				sub, err := token.GetSubject()
				if err != nil {
					return nil, errors.Unauthorized("INVALID_TOKEN", err.Error())
				}

				allowed, err := e.Enforce(sub, perm.Object, perm.Action)
				if err != nil {
					return nil, err
				}
				if !allowed {
					return nil, errors.Forbidden("ACCESS_DENIED", "permission denied")
				}

				return next(ctx, req)
			})

			return handlerWithJWT(ctx, req)
		}
	}
}
