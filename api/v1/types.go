package v1

type LocalObjectReference struct {
	// Name of the referent.
	Name string `json:"name"`
}

// CABundleReference identifies a Secret or ConfigMap containing a PEM-encoded
// CA bundle. Kind defaults to Secret and Key defaults to ca.crt for backwards
// compatibility with the original caBundleRef format.
type CABundleReference struct {
	// Name of the referent.
	Name string `json:"name"`
	// Kind of the referent. Supported values are Secret and ConfigMap.
	// +optional
	// +kubebuilder:validation:Enum=Secret;ConfigMap
	Kind string `json:"kind,omitempty"`
	// Key in the referent containing the PEM-encoded CA bundle.
	// +optional
	Key string `json:"key,omitempty"`
	// Keys in the referent containing PEM-encoded CA certificates. Values are
	// concatenated in the listed order. If set, Keys takes precedence over Key.
	// +optional
	// +kubebuilder:validation:MinItems=1
	Keys []string `json:"keys,omitempty"`
}
