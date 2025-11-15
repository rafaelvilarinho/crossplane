package revision

import (
	"github.com/crossplane/crossplane-runtime/v2/pkg/logging"
	v1 "github.com/crossplane/crossplane/v2/apis/pkg/v1"
	v2 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// ResourceFilter filter resources of package based on selectors
type ResourceFilter struct {
	excluded []v1.ResourceSelector
	log      logging.Logger
}

// NewResourceFilter Creates a new filter of resources
func NewResourceFilter(excluded []v1.ResourceSelector, log logging.Logger) *ResourceFilter {
	return &ResourceFilter{
		excluded: excluded,
		log:      log,
	}
}

// Filter an object list of runtime.Object based on selectors
func (f *ResourceFilter) Filter(objects []runtime.Object) []runtime.Object {
	var filtered []runtime.Object

	for _, obj := range objects {
		if f.shouldExclude(obj) {
			continue
		}

		filtered = append(filtered, obj)
	}

	return filtered
}

// shouldExclude define if an object must be excluded
func (f *ResourceFilter) shouldExclude(obj runtime.Object) bool {
	u, ok := obj.(*v2.CustomResourceDefinition)
	f.log.Debug("-----------> shouldExclude 1", "u", u, "ok", ok, "obj", obj)
	if !ok {
		return false
	}

	kind := u.GroupVersionKind().Kind
	name := u.GetName()
	group := u.GroupVersionKind().Group
	f.log.Debug("-----------> shouldExclude 2", "kind", kind, "name", name, "group", group)
	f.log.Debug("-----------> shouldExclude 3", "excluded", f.excluded)

	if len(f.excluded) == 0 {
		return false
	}

	return f.matchesAny(f.excluded, kind, name, group)

}

// matchesAny verifies if resource corresponds to some selector
func (f *ResourceFilter) matchesAny(selectors []v1.ResourceSelector, kind, name, group string) bool {
	for _, selector := range selectors {
		if f.matches(selector, kind, name, group) {
			return false
		}
	}

	return true
}

// matches verifies if a resource corresponds to one specific selector
func (f *ResourceFilter) matches(selector v1.ResourceSelector, kind, name, group string) bool {
	f.log.Debug("-------> TEST", "selectorName", selector.Name, "kind", kind, "name", name, "group", group)

	if selector.Group != "" && selector.Group == group {
		return true
	}

	if selector.Name != "" && selector.Name != name {
		return true
	}

	return false
}

// GetFilterStats return filtering statistics for logs
func (f *ResourceFilter) GetFilterStats(original, filtered []runtime.Object) map[string]int {
	return map[string]int{
		"original": len(original),
		"filtered": len(filtered),
		"removed":  len(original) - len(filtered),
	}
}
