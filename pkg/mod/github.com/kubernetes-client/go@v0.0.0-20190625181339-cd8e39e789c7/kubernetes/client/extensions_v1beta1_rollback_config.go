/*
 * Kubernetes
 *
 * No description provided (generated by Swagger Codegen https://github.com/swagger-api/swagger-codegen)
 *
 * API version: v1.10.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package client

// DEPRECATED.
type ExtensionsV1beta1RollbackConfig struct {

	// The revision to rollback to. If set to 0, rollback to the last revision.
	Revision int64 `json:"revision,omitempty"`
}
