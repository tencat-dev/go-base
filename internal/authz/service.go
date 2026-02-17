package authz

import (
	"fmt"
	"sync/atomic"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"

	authzv1 "github.com/tencat-dev/go-base/api/authz/v1"
)

type AuthzRegistry struct {
	data atomic.Value
}

func NewAuthzRegistry() *AuthzRegistry {
	r := &AuthzRegistry{}
	r.data.Store(make(map[string]*authzv1.PermissionOption))
	r.loadFromProto()
	return r
}

func (r *AuthzRegistry) loadFromProto() {
	newMap := make(map[string]*authzv1.PermissionOption)

	protoregistry.GlobalFiles.RangeFiles(func(fd protoreflect.FileDescriptor) bool {
		services := fd.Services()
		for i := 0; i < services.Len(); i++ {
			svc := services.Get(i)
			methods := svc.Methods()

			for j := 0; j < methods.Len(); j++ {
				method := methods.Get(j)
				opts := method.Options()
				if opts == nil {
					continue
				}

				ext := proto.GetExtension(opts, authzv1.E_Permission)
				perm, ok := ext.(*authzv1.PermissionOption)
				if !ok || perm == nil {
					continue
				}

				fullMethod := fmt.Sprintf("/%s/%s",
					string(svc.FullName()),
					string(method.Name()),
				)

				newMap[fullMethod] = perm
			}
		}
		return true
	})

	// 🔥 atomic swap toàn bộ map
	r.data.Store(newMap)
}

func (r *AuthzRegistry) Get(op string) (*authzv1.PermissionOption, bool) {
	m := r.data.Load().(map[string]*authzv1.PermissionOption)
	perm, ok := m[op]
	return perm, ok
}
