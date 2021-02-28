package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +resource:path=apimanagement
type RouteConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              RouteConfigSpec `json:"spec"`
}

// +k8s:deepcopy-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +resource:path=apimanagement
type RouteConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []RouteConfig `json:"items"`
}

// RouteConfigSpec is configuration of routes for API-management
// +k8s:deepcopy-gen=true
type RouteConfigSpec struct {
	Gateway string  `json:"gateway"`
	Routes  []Route `json:"routes"`
}

type Route struct {
	Destination Destination `json:"destination"`
	Match       Match       `json:"match"`
	PathRewrite string      `json:"pathRewrite"`
}

type Match struct {
	Path string `json:"path"`
}

type Destination struct {
	Address Address `json:"address"`
}

type Address struct {
	Host string `json:"host"`
	Port uint32 `json:"port"`
}
