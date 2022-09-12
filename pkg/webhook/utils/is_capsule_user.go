package utils

import (
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	"github.com/clastix/capsule/pkg/utils"
)

func IsCapsuleUser(req admission.Request, userGroups []string) bool {
	groupList := utils.NewUserGroupList(req.UserInfo.Groups)
	// if the user is a ServiceAccount belonging to an exluded namespace, definitely, it's not a Capsule user
	// and we can skip the check in case of Capsule user group assigned to system:authenticated
	// (ref: https://github.com/clastix/capsule/issues/234)
	excludedNamespaces := []string{"kube-system", "flux-system", "capsule-system"}
	for _, ns := range excludedNamespaces {
		if groupList.Find(fmt.Sprintf("system:serviceaccounts:%s", ns)) {
			return false
		}
	}

	for _, group := range userGroups {
		if groupList.Find(group) {
			return true
		}
	}

	return false
}
