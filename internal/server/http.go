package server

import (
	"encoding/json"
	"reflect"

	"github.com/bytedance/sonic"
	"github.com/go-kratos/kratos/contrib/middleware/validate/v2"
	"github.com/go-kratos/kratos/v2/encoding"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/gorilla/handlers"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	authv1 "github.com/tencat-dev/go-base/api/auth/v1"
	authzv1 "github.com/tencat-dev/go-base/api/authz/v1"
	userv1 "github.com/tencat-dev/go-base/api/user/v1"
	"github.com/tencat-dev/go-base/internal/authz"
	"github.com/tencat-dev/go-base/internal/conf"
)

type HttpServer transport.Server

// NewHTTPServer new an HTTP server.
func NewHTTPServer(
	c *conf.HTTPServer,
	userService userv1.UserServiceServer,
	authService authv1.AuthServiceServer,
	authzService authzv1.AuthzServiceServer,
	logger log.Logger,
	authzMiddleware authz.AuthzMiddleware,
) HttpServer {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			logging.Server(logger),
			validate.ProtoValidate(),
			middleware.Middleware(authzMiddleware),
		),
		http.Filter(handlers.CORS(
			handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
			handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}),
			handlers.AllowedOrigins([]string{"*"}),
		)),
	}
	if c.Network != "" {
		opts = append(opts, http.Network(c.Network))
	}
	if c.Addr != "" {
		opts = append(opts, http.Address(c.Addr))
	}
	if c.Timeout != nil {
		opts = append(opts, http.Timeout(c.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	userv1.RegisterUserServiceHTTPServer(srv, userService)
	authv1.RegisterAuthServiceHTTPServer(srv, authService)
	authzv1.RegisterAuthzServiceHTTPServer(srv, authzService)
	return srv
}

// Name is the name registered for the json codec.
const Name = "json"

var (
	// MarshalOptions is a configurable JSON format marshaller.
	MarshalOptions = protojson.MarshalOptions{
		UseProtoNames:   true,
		EmitUnpopulated: true,
	}
	// UnmarshalOptions is a configurable JSON format parser.
	UnmarshalOptions = protojson.UnmarshalOptions{
		DiscardUnknown: true,
	}
)

func init() {
	encoding.RegisterCodec(codec{})
}

// codec is a Codec implementation with json.
type codec struct{}

func (codec) Marshal(v interface{}) ([]byte, error) {
	switch m := v.(type) {
	case json.Marshaler:
		return m.MarshalJSON()
	case proto.Message:
		return MarshalOptions.Marshal(m)
	default:
		return sonic.Marshal(m)
	}
}

func (codec) Unmarshal(data []byte, v interface{}) error {
	switch m := v.(type) {
	case json.Unmarshaler:
		return m.UnmarshalJSON(data)
	case proto.Message:
		return UnmarshalOptions.Unmarshal(data, m)
	default:
		rv := reflect.ValueOf(v)
		for rv := rv; rv.Kind() == reflect.Ptr; {
			if rv.IsNil() {
				rv.Set(reflect.New(rv.Type().Elem()))
			}
			rv = rv.Elem()
		}
		if m, ok := reflect.Indirect(rv).Interface().(proto.Message); ok {
			return UnmarshalOptions.Unmarshal(data, m)
		}
		return sonic.Unmarshal(data, m)
	}
}

func (codec) Name() string {
	return Name
}
