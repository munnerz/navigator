package validation

import (
	"fmt"

	apimachineryvalidation "k8s.io/apimachinery/pkg/api/validation"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/jetstack/navigator/pkg/apis/navigator"
)

func ValidateImageSpec(img *navigator.ImageSpec, fldPath *field.Path) field.ErrorList {
	el := field.ErrorList{}
	if img.Tag == "" {
		el = append(el, field.Required(fldPath.Child("tag"), ""))
	}
	if img.Repository == "" {
		el = append(el, field.Required(fldPath.Child("repository"), ""))
	}
	return el
}

var supportedESClusterRoles = []string{
	string(navigator.ElasticsearchRoleData),
	string(navigator.ElasticsearchRoleIngest),
	string(navigator.ElasticsearchRoleMaster),
}

func ValidateElasticsearchClusterRole(r navigator.ElasticsearchClusterRole, fldPath *field.Path) field.ErrorList {
	el := field.ErrorList{}
	switch r {
	case navigator.ElasticsearchRoleData:
	case navigator.ElasticsearchRoleIngest:
	case navigator.ElasticsearchRoleMaster:
	default:
		el = append(el, field.NotSupported(fldPath, r, supportedESClusterRoles))
	}
	return el
}

func ValidatElasticsearchClusterNodePool(np *navigator.ElasticsearchClusterNodePool, fldPath *field.Path) field.ErrorList {
	el := ValidateDNS1123Subdomain(np.Name, fldPath.Child("name"))
	el = ValidateElasticsearchPersistence(&np.Persistence, fldPath.Child("persistence"))
	rolesPath := fldPath.Child("roles")
	for i, r := range np.Roles {
		idxPath := rolesPath.Index(i)
		el = append(el, ValidateElasticsearchClusterRole(r, idxPath)...)
	}
	if np.Replicas < 0 {
		el = append(el, field.Invalid(fldPath.Child("replicas"), np.Replicas, "must be greater than zero"))
	}
	// TODO: call k8s.io/kubernetes/pkg/apis/core/validation.ValidateResourceRequirements on np.Resources
	// this will require vendoring kubernetes/kubernetes and switching to use the corev1 ResourceRequirements
	// struct
	return el
}

func ValidateElasticsearchPersistence(cfg *navigator.ElasticsearchClusterPersistenceConfig, fldPath *field.Path) field.ErrorList {
	el := field.ErrorList{}
	if cfg.Enabled && len(cfg.Size) == 0 {
		el = append(el, field.Required(fldPath.Child("size"), ""))
	}
	// TODO: validate size quantity
	return el
}

func ValidateElasticsearchClusterSpec(spec *navigator.ElasticsearchClusterSpec, fldPath *field.Path) field.ErrorList {
	allErrs := ValidateImageSpec(&spec.Image.ImageSpec, fldPath.Child("image"))
	allErrs = append(allErrs, ValidateImageSpec(&spec.Pilot.ImageSpec, fldPath.Child("pilot"))...)
	npPath := fldPath.Child("nodePools")
	allNames := sets.String{}
	for i, np := range spec.NodePools {
		idxPath := npPath.Index(i)
		if allNames.Has(np.Name) {
			allErrs = append(allErrs, field.Duplicate(idxPath.Child("name"), np.Name))
		} else {
			allNames.Insert(np.Name)
		}
		allErrs = append(allErrs, ValidatElasticsearchClusterNodePool(&np, idxPath)...)
	}

	numMasters := countMasterReplicas(spec.NodePools)
	quorom := calculateQuorom(numMasters)
	if numMasters == 0 {
		allErrs = append(allErrs, field.Invalid(npPath, numMasters, "must be at least one master node"))
	} else if spec.MinimumMasters < quorom {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("minimumMasters"), spec.MinimumMasters, fmt.Sprintf("must be a minimum of %d to avoid a split brain scenario", quorom)))
	}

	return allErrs
}

func calculateQuorom(num int64) int64 {
	if num == 1 {
		return 1
	}
	return (num / 2) + 1
}

func countMasterReplicas(pools []navigator.ElasticsearchClusterNodePool) int64 {
	masters := int64(0)
	for _, pool := range pools {
		if hasRole(pool.Roles, navigator.ElasticsearchRoleMaster) {
			masters += pool.Replicas
		}
	}
	return masters
}

func hasRole(set []navigator.ElasticsearchClusterRole, role navigator.ElasticsearchClusterRole) bool {
	for _, s := range set {
		if s == role {
			return true
		}
	}
	return false
}

func ValidateElasticsearchCluster(esc *navigator.ElasticsearchCluster) field.ErrorList {
	allErrs := ValidateObjectMeta(&esc.ObjectMeta, true, apimachineryvalidation.NameIsDNSSubdomain, field.NewPath("metadata"))
	allErrs = append(allErrs, ValidateElasticsearchClusterSpec(&esc.Spec, field.NewPath("spec"))...)
	return allErrs
}
