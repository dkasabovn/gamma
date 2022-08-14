package user

import "gamma/app/system/stability"

func init() {
	stability.LoadDependencies(stability.UserSvc())
}
