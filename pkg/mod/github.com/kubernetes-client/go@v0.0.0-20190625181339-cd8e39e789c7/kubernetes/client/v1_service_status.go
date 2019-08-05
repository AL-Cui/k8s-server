/*
 * Kubernetes
 *
 * No description provided (generated by Swagger Codegen https://github.com/swagger-api/swagger-codegen)
 *
 * API version: v1.10.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package client

// ServiceStatus represents the current status of a service.
type V1ServiceStatus struct {

	// LoadBalancer contains the current status of the load-balancer, if one is present.
	LoadBalancer *V1LoadBalancerStatus `json:"loadBalancer,omitempty"`
}
