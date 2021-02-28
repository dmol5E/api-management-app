package v1alpha1

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/dmol5e/api-management-app/api-publisher/pkg/apis/apimanagement"
	apiextensionv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiextension "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/wait"
)

// CreateCRD creates CustomResourceDefinition for RouteConfig
func CreateCRD(ctx context.Context, clientSet apiextension.Interface) (*apiextensionv1.CustomResourceDefinition, error) {
	crd := &apiextensionv1.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: CRDName,
		},
		Spec: apiextensionv1.CustomResourceDefinitionSpec{
			Group: apimanagement.GroupName,
			Versions: []apiextensionv1.CustomResourceDefinitionVersion{
				{
					Name:    SchemeGroupVersion.Version,
					Served:  true,
					Storage: true,
				},
			},
			Scope: apiextensionv1.NamespaceScoped,
			Names: apiextensionv1.CustomResourceDefinitionNames{
				Plural:     Plural,
				Singular:   Singular,
				Kind:       reflect.TypeOf(RouteConfig{}).Name(),
				ShortNames: []string{ShortName},
			},
		},
	}

	crd, err := clientSet.ApiextensionsV1().CustomResourceDefinitions().Create(ctx, crd, metav1.CreateOptions{})
	if err != nil {
		if apierrors.IsAlreadyExists(err) {
			return crd, nil
		}
		return nil, err
	}

	// Wait for CRD creation.
	err = wait.Poll(5*time.Second, 60*time.Second, func() (bool, error) {
		crd, err = clientSet.ApiextensionsV1().CustomResourceDefinitions().Get(ctx, CRDName, metav1.GetOptions{})
		if err != nil {
			fmt.Printf("Fail to wait for CRD Jinghzhu creation: %+v\n", err)

			return false, err
		}
		for _, cond := range crd.Status.Conditions {
			switch cond.Type {
			case apiextensionv1.Established:
				if cond.Status == apiextensionv1.ConditionTrue {
					return true, err
				}
			case apiextensionv1.NamesAccepted:
				if cond.Status == apiextensionv1.ConditionFalse {
					fmt.Printf("Name conflict while wait for CRD Jinghzhu creation: %s, %+v\n", cond.Reason, err)
				}
			}
		}

		return false, err
	})

	// If there is an error, delete the object to keep it clean.
	if err != nil {
		fmt.Println("Try to cleanup")
		deleteErr := clientSet.ApiextensionsV1beta1().CustomResourceDefinitions().Delete(ctx, CRDName, metav1.DeleteOptions{})
		if deleteErr != nil {
			fmt.Printf("Fail to delete CRD Jinghzhu: %+v\n", deleteErr)

			return nil, errors.NewAggregate([]error{err, deleteErr})
		}

		return nil, err
	}

	return crd, nil
}
