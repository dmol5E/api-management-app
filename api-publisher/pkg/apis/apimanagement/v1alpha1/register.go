package v1alpha1

import (
	"github.com/dmol5e/api-management-app/api-publisher/pkg/apis/apimanagement"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	// Kind is normally the CamelCased singular type. The resource manifest uses this.
	Kind string = "APIConfig"
	// GroupVersion is the version.
	GroupVersion string = "v1alpha1"
	// Plural is the plural name used in /apis/<group>/<version>/<plural>
	Plural string = "apiconfigs"
	// Singular is used as an alias on kubectl for display.
	Singular string = "apiconfig"
	// CRDName is the CRD name for RouteConfig.
	CRDName string = Plural + "." + apimanagement.GroupName
	// ShortName is the short alias for the CRD.
	ShortName string = "apicfg"
)

var (
	// SchemeGroupVersion is the primary key to register CRD in k8s
	SchemeGroupVersion = schema.GroupVersion{Group: apimanagement.GroupName, Version: GroupVersion}
	SchemeBuilder      = runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme        = SchemeBuilder.AddToScheme
)

// Resource takes an unqualified resource and returns a Group qualified GroupResource
func Resource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&APIConfig{},
		&APIConfigList{},
	)
	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}
