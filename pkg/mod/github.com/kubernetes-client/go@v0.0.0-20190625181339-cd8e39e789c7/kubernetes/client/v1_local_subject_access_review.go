/*
 * Kubernetes
 *
 * No description provided (generated by Swagger Codegen https://github.com/swagger-api/swagger-codegen)
 *
 * API version: v1.10.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package client

// LocalSubjectAccessReview checks whether or not a user or group can perform an action in a given namespace. Having a namespace scoped resource makes it much easier to grant namespace scoped policy that includes permissions checking.
type V1LocalSubjectAccessReview struct {

	// APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#resources
	ApiVersion string `json:"apiVersion,omitempty"`

	// Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds
	Kind string `json:"kind,omitempty"`

	Metadata *V1ObjectMeta `json:"metadata,omitempty"`

	// Spec holds information about the request being evaluated.  spec.namespace must be equal to the namespace you made the request against.  If empty, it is defaulted.
	Spec *V1SubjectAccessReviewSpec `json:"spec"`

	// Status is filled in by the server and indicates whether the request is allowed or not
	Status *V1SubjectAccessReviewStatus `json:"status,omitempty"`
}
