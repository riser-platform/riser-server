/*
KNative is not yet consumable using Go Modules (https://github.com/knative/serving/issues/3368).
Since it's desirable to have a statically typed object model for KNative resources, some parts of this code
were pulled from https://github.com/knative/serving@e8fe5efc640acc6da99c545fc826abb2c997a207 and
falls under the following Copyright license.


Copyright 2019 The Knative Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package resources

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knative.dev/pkg/apis"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Service acts as a top-level container that manages a Route and Configuration
// which implement a network service. Service exists to provide a singular
// abstraction which can be access controlled, reasoned about, and which
// encapsulates software lifecycle decisions such as rollout policy and
// team resource ownership. Service acts only as an orchestrator of the
// underlying Routes and Configurations (much as a kubernetes Deployment
// orchestrates ReplicaSets), and its usage is optional but recommended.
//
// The Service's controller will track the statuses of its owned Configuration
// and Route, reflecting their statuses and conditions as its own.
//
// See also: https://knative.dev/serving/blob/master/docs/spec/overview.md#service
type Service struct {
	metav1.ObjectMeta `json:"metadata"`
	metav1.TypeMeta   `json:",inline"`
	ServiceSpec       `json:"spec"`
}

// ServiceSpec represents the configuration for the Service object.
// A Service's specification is the union of the specifications for a Route
// and Configuration.  The Service restricts what can be expressed in these
// fields, e.g. the Route must reference the provided Configuration;
// however, these limitations also enable friendlier defaulting,
// e.g. Route never needs a Configuration name, and may be defaulted to
// the appropriate "run latest" spec.
type ServiceSpec struct {
	// ServiceSpec inlines an unrestricted ConfigurationSpec.
	ConfigurationSpec `json:",inline"`

	// ServiceSpec inlines RouteSpec and restricts/defaults its fields
	// via webhook.  In particular, this spec can only reference this
	// Service's configuration and revisions (which also influences
	// defaults).
	RouteSpec `json:",inline"`
}

// ConfigurationSpec holds the desired state of the Configuration (from the client).
type ConfigurationSpec struct {
	// Template holds the latest specification for the Revision to be stamped out.
	// +optional
	Template RevisionTemplateSpec `json:"template"`
}

// RevisionTemplateSpec describes the data a revision should have when created from a template.
// Based on: https://github.com/kubernetes/api/blob/e771f807/core/v1/types.go#L3179-L3190
type RevisionTemplateSpec struct {
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// +optional
	Spec RevisionSpec `json:"spec,omitempty"`
}

// RevisionSpec holds the desired state of the Revision (from the client).
type RevisionSpec struct {
	corev1.PodSpec `json:",inline"`

	// ContainerConcurrency specifies the maximum allowed in-flight (concurrent)
	// requests per container of the Revision.  Defaults to `0` which means
	// concurrency to the application is not limited, and the system decides the
	// target concurrency for the autoscaler.
	// +optional
	ContainerConcurrency *int64 `json:"containerConcurrency,omitempty"`

	// TimeoutSeconds holds the max duration the instance is allowed for
	// responding to a request.  If unspecified, a system default will
	// be provided.
	// +optional
	TimeoutSeconds *int64 `json:"timeoutSeconds,omitempty"`
}

// RouteSpec holds the desired state of the Route (from the client).
type RouteSpec struct {
	// Traffic specifies how to distribute traffic over a collection of
	// revisions and configurations.
	// +optional
	Traffic []TrafficTarget `json:"traffic,omitempty"`
}

// TrafficTarget holds a single entry of the routing table for a Route.
type TrafficTarget struct {
	// Tag is optionally used to expose a dedicated url for referencing
	// this target exclusively.
	// +optional
	Tag string `json:"tag,omitempty"`

	// RevisionName of a specific revision to which to send this portion of
	// traffic.  This is mutually exclusive with ConfigurationName.
	// +optional
	RevisionName string `json:"revisionName,omitempty"`

	// ConfigurationName of a configuration to whose latest revision we will send
	// this portion of traffic. When the "status.latestReadyRevisionName" of the
	// referenced configuration changes, we will automatically migrate traffic
	// from the prior "latest ready" revision to the new one.  This field is never
	// set in Route's status, only its spec.  This is mutually exclusive with
	// RevisionName.
	// +optional
	ConfigurationName string `json:"configurationName,omitempty"`

	// LatestRevision may be optionally provided to indicate that the latest
	// ready Revision of the Configuration should be used for this traffic
	// target.  When provided LatestRevision must be true if RevisionName is
	// empty; it must be false when RevisionName is non-empty.
	// +optional
	LatestRevision *bool `json:"latestRevision,omitempty"`

	// Percent indicates that percentage based routing should be used and
	// the value indicates the percent of traffic that is be routed to this
	// Revision or Configuration. `0` (zero) mean no traffic, `100` means all
	// traffic.
	// When percentage based routing is being used the follow rules apply:
	// - the sum of all percent values must equal 100
	// - when not specified, the implied value for `percent` is zero for
	//   that particular Revision or Configuration
	// +optional
	Percent *int64 `json:"percent,omitempty"`

	// URL displays the URL for accessing named traffic targets. URL is displayed in
	// status, and is disallowed on spec. URL must contain a scheme (e.g. http://) and
	// a hostname, but may not contain anything else (e.g. basic auth, url path, etc.)
	// +optional
	URL *apis.URL `json:"url,omitempty"`
}
