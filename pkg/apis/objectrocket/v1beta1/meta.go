package v1beta1

// ObjectMeta is the sensu object metadata
type ObjectMeta struct {
	// Name must be unique within a namespace. Name is primarily intended for creation
	// idempotence and configuration definition.
	Name string `json:"name,omitempty" yaml: "name,omitempty"`
	// Namespace defines a logical grouping of objects within which each object name must
	// be unique.
	Namespace string `json:"namespace,omitempty" yaml: "namespace,omitempty"`
	// Map of string keys and values that can be used to organize and categorize
	// (scope and select) objects. May also be used in filters and token
	// substitution.
	// TODO: Link to Sensu documentation.
	// More info: http://kubernetes.io/docs/user-guide/labels
	Labels map[string]string `json:"labels,omitempty" yaml: ",labels,omitempty"`
	// Annotations is an unstructured key value map stored with a resource that may be
	// set by external tools to store and retrieve arbitrary metadata. They are not
	// queryable and should be preserved when modifying objects.
	// TODO: Link to Sensu documentation.
	// More info: http://kubernetes.io/docs/user-guide/annotations
	Annotations map[string]string `json:"annotations,omitempty" yaml: "annotations,omitempty" `
}
