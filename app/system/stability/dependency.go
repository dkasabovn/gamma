package stability

import (
	"gamma/app/services/user"
)

type Loader func()

func UserSvc() Loader {
	return func() {
		user.GetUserService()
	}
}

func LoadDependencies(opts ...Loader) {
	for _, loader := range opts {
		loader()
	}
}
