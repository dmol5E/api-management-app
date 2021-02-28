package v1alpha

import (
	"context"
	"reflect"

	apiextensionv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiextension "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	CRDKind     string = "RouteConfig"
	CRDPlural   string = "routeconfigs"
	CRDSingular string = "routeconfig"
	CRDGroup    string = "cloud.apimanagement"
	CRDVersion  string = "v1alpha1"
	FullCRDName string = CRDPlural + "." + CRDGroup
)

func CreateCRD(ctx context.Context, clientset apiextension.Interface) error {
	crd := &apiextensionv1.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{Name: FullCRDName},
		Spec: apiextensionv1.CustomResourceDefinitionSpec{
			Group: CRDGroup,
			Scope: apiextensionv1.NamespaceScoped,
			Names: apiextensionv1.CustomResourceDefinitionNames{
				Plural:   CRDPlural,
				Singular: CRDSingular,
				Kind:     reflect.TypeOf(RouteConfig{}).Name(),
			},
		},
	}

	_, err := clientset.ApiextensionsV1().CustomResourceDefinitions().Create(ctx, crd, metav1.CreateOptions{})
	if err != nil {
		if apierrors.IsAlreadyExists(err) {
			return nil
		}
		return err
	}
	return nil
}
