// +build !ignore_autogenerated

/*
Copyright 2019 The sensu-operator Authors

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

// Code generated by deepcopy-gen. DO NOT EDIT.

package v1beta1

import (
	v1 "k8s.io/api/core/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ABSBackupSource) DeepCopyInto(out *ABSBackupSource) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ABSBackupSource.
func (in *ABSBackupSource) DeepCopy() *ABSBackupSource {
	if in == nil {
		return nil
	}
	out := new(ABSBackupSource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ABSRestoreSource) DeepCopyInto(out *ABSRestoreSource) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ABSRestoreSource.
func (in *ABSRestoreSource) DeepCopy() *ABSRestoreSource {
	if in == nil {
		return nil
	}
	out := new(ABSRestoreSource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BackupPolicy) DeepCopyInto(out *BackupPolicy) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BackupPolicy.
func (in *BackupPolicy) DeepCopy() *BackupPolicy {
	if in == nil {
		return nil
	}
	out := new(BackupPolicy)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BackupSource) DeepCopyInto(out *BackupSource) {
	*out = *in
	if in.S3 != nil {
		in, out := &in.S3, &out.S3
		*out = new(S3BackupSource)
		**out = **in
	}
	if in.ABS != nil {
		in, out := &in.ABS, &out.ABS
		*out = new(ABSBackupSource)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BackupSource.
func (in *BackupSource) DeepCopy() *BackupSource {
	if in == nil {
		return nil
	}
	out := new(BackupSource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BackupSpec) DeepCopyInto(out *BackupSpec) {
	*out = *in
	if in.EtcdEndpoints != nil {
		in, out := &in.EtcdEndpoints, &out.EtcdEndpoints
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.BackupPolicy != nil {
		in, out := &in.BackupPolicy, &out.BackupPolicy
		*out = new(BackupPolicy)
		**out = **in
	}
	in.BackupSource.DeepCopyInto(&out.BackupSource)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BackupSpec.
func (in *BackupSpec) DeepCopy() *BackupSpec {
	if in == nil {
		return nil
	}
	out := new(BackupSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BackupStatus) DeepCopyInto(out *BackupStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BackupStatus.
func (in *BackupStatus) DeepCopy() *BackupStatus {
	if in == nil {
		return nil
	}
	out := new(BackupStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterCondition) DeepCopyInto(out *ClusterCondition) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterCondition.
func (in *ClusterCondition) DeepCopy() *ClusterCondition {
	if in == nil {
		return nil
	}
	out := new(ClusterCondition)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterSpec) DeepCopyInto(out *ClusterSpec) {
	*out = *in
	if in.Pod != nil {
		in, out := &in.Pod, &out.Pod
		*out = new(PodPolicy)
		(*in).DeepCopyInto(*out)
	}
	if in.TLS != nil {
		in, out := &in.TLS, &out.TLS
		*out = new(TLSPolicy)
		(*in).DeepCopyInto(*out)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterSpec.
func (in *ClusterSpec) DeepCopy() *ClusterSpec {
	if in == nil {
		return nil
	}
	out := new(ClusterSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterStatus) DeepCopyInto(out *ClusterStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]ClusterCondition, len(*in))
		copy(*out, *in)
	}
	in.Members.DeepCopyInto(&out.Members)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterStatus.
func (in *ClusterStatus) DeepCopy() *ClusterStatus {
	if in == nil {
		return nil
	}
	out := new(ClusterStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HandlerSocket) DeepCopyInto(out *HandlerSocket) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HandlerSocket.
func (in *HandlerSocket) DeepCopy() *HandlerSocket {
	if in == nil {
		return nil
	}
	out := new(HandlerSocket)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HookList) DeepCopyInto(out *HookList) {
	*out = *in
	if in.Hooks != nil {
		in, out := &in.Hooks, &out.Hooks
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HookList.
func (in *HookList) DeepCopy() *HookList {
	if in == nil {
		return nil
	}
	out := new(HookList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MemberSecret) DeepCopyInto(out *MemberSecret) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MemberSecret.
func (in *MemberSecret) DeepCopy() *MemberSecret {
	if in == nil {
		return nil
	}
	out := new(MemberSecret)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MembersStatus) DeepCopyInto(out *MembersStatus) {
	*out = *in
	if in.Ready != nil {
		in, out := &in.Ready, &out.Ready
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Unready != nil {
		in, out := &in.Unready, &out.Unready
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MembersStatus.
func (in *MembersStatus) DeepCopy() *MembersStatus {
	if in == nil {
		return nil
	}
	out := new(MembersStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ObjectMeta) DeepCopyInto(out *ObjectMeta) {
	*out = *in
	if in.Labels != nil {
		in, out := &in.Labels, &out.Labels
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Annotations != nil {
		in, out := &in.Annotations, &out.Annotations
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ObjectMeta.
func (in *ObjectMeta) DeepCopy() *ObjectMeta {
	if in == nil {
		return nil
	}
	out := new(ObjectMeta)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PodPolicy) DeepCopyInto(out *PodPolicy) {
	*out = *in
	if in.Labels != nil {
		in, out := &in.Labels, &out.Labels
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.NodeSelector != nil {
		in, out := &in.NodeSelector, &out.NodeSelector
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Affinity != nil {
		in, out := &in.Affinity, &out.Affinity
		*out = new(v1.Affinity)
		(*in).DeepCopyInto(*out)
	}
	in.Resources.DeepCopyInto(&out.Resources)
	if in.Tolerations != nil {
		in, out := &in.Tolerations, &out.Tolerations
		*out = make([]v1.Toleration, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.SensuEnv != nil {
		in, out := &in.SensuEnv, &out.SensuEnv
		*out = make([]v1.EnvVar, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.PersistentVolumeClaimSpec != nil {
		in, out := &in.PersistentVolumeClaimSpec, &out.PersistentVolumeClaimSpec
		*out = new(v1.PersistentVolumeClaimSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.Annotations != nil {
		in, out := &in.Annotations, &out.Annotations
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.SecurityContext != nil {
		in, out := &in.SecurityContext, &out.SecurityContext
		*out = new(v1.PodSecurityContext)
		(*in).DeepCopyInto(*out)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PodPolicy.
func (in *PodPolicy) DeepCopy() *PodPolicy {
	if in == nil {
		return nil
	}
	out := new(PodPolicy)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ProxyRequests) DeepCopyInto(out *ProxyRequests) {
	*out = *in
	if in.EntityAttributes != nil {
		in, out := &in.EntityAttributes, &out.EntityAttributes
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ProxyRequests.
func (in *ProxyRequests) DeepCopy() *ProxyRequests {
	if in == nil {
		return nil
	}
	out := new(ProxyRequests)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RestoreSource) DeepCopyInto(out *RestoreSource) {
	*out = *in
	if in.S3 != nil {
		in, out := &in.S3, &out.S3
		*out = new(S3RestoreSource)
		**out = **in
	}
	if in.ABS != nil {
		in, out := &in.ABS, &out.ABS
		*out = new(ABSRestoreSource)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RestoreSource.
func (in *RestoreSource) DeepCopy() *RestoreSource {
	if in == nil {
		return nil
	}
	out := new(RestoreSource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RestoreSpec) DeepCopyInto(out *RestoreSpec) {
	*out = *in
	in.RestoreSource.DeepCopyInto(&out.RestoreSource)
	out.SensuCluster = in.SensuCluster
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RestoreSpec.
func (in *RestoreSpec) DeepCopy() *RestoreSpec {
	if in == nil {
		return nil
	}
	out := new(RestoreSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RestoreStatus) DeepCopyInto(out *RestoreStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RestoreStatus.
func (in *RestoreStatus) DeepCopy() *RestoreStatus {
	if in == nil {
		return nil
	}
	out := new(RestoreStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *S3BackupSource) DeepCopyInto(out *S3BackupSource) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new S3BackupSource.
func (in *S3BackupSource) DeepCopy() *S3BackupSource {
	if in == nil {
		return nil
	}
	out := new(S3BackupSource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *S3RestoreSource) DeepCopyInto(out *S3RestoreSource) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new S3RestoreSource.
func (in *S3RestoreSource) DeepCopy() *S3RestoreSource {
	if in == nil {
		return nil
	}
	out := new(S3RestoreSource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SensuAsset) DeepCopyInto(out *SensuAsset) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SensuAsset.
func (in *SensuAsset) DeepCopy() *SensuAsset {
	if in == nil {
		return nil
	}
	out := new(SensuAsset)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *SensuAsset) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SensuAssetList) DeepCopyInto(out *SensuAssetList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]SensuAsset, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SensuAssetList.
func (in *SensuAssetList) DeepCopy() *SensuAssetList {
	if in == nil {
		return nil
	}
	out := new(SensuAssetList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *SensuAssetList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SensuAssetSpec) DeepCopyInto(out *SensuAssetSpec) {
	*out = *in
	if in.Metadata != nil {
		in, out := &in.Metadata, &out.Metadata
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Filters != nil {
		in, out := &in.Filters, &out.Filters
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	in.SensuMetadata.DeepCopyInto(&out.SensuMetadata)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SensuAssetSpec.
func (in *SensuAssetSpec) DeepCopy() *SensuAssetSpec {
	if in == nil {
		return nil
	}
	out := new(SensuAssetSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SensuAssetStatus) DeepCopyInto(out *SensuAssetStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SensuAssetStatus.
func (in *SensuAssetStatus) DeepCopy() *SensuAssetStatus {
	if in == nil {
		return nil
	}
	out := new(SensuAssetStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SensuBackup) DeepCopyInto(out *SensuBackup) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SensuBackup.
func (in *SensuBackup) DeepCopy() *SensuBackup {
	if in == nil {
		return nil
	}
	out := new(SensuBackup)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *SensuBackup) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SensuBackupList) DeepCopyInto(out *SensuBackupList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]SensuBackup, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SensuBackupList.
func (in *SensuBackupList) DeepCopy() *SensuBackupList {
	if in == nil {
		return nil
	}
	out := new(SensuBackupList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *SensuBackupList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SensuCheckConfig) DeepCopyInto(out *SensuCheckConfig) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SensuCheckConfig.
func (in *SensuCheckConfig) DeepCopy() *SensuCheckConfig {
	if in == nil {
		return nil
	}
	out := new(SensuCheckConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *SensuCheckConfig) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SensuCheckConfigList) DeepCopyInto(out *SensuCheckConfigList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]SensuCheckConfig, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SensuCheckConfigList.
func (in *SensuCheckConfigList) DeepCopy() *SensuCheckConfigList {
	if in == nil {
		return nil
	}
	out := new(SensuCheckConfigList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *SensuCheckConfigList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SensuCheckConfigSpec) DeepCopyInto(out *SensuCheckConfigSpec) {
	*out = *in
	if in.Handlers != nil {
		in, out := &in.Handlers, &out.Handlers
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.RuntimeAssets != nil {
		in, out := &in.RuntimeAssets, &out.RuntimeAssets
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Subscriptions != nil {
		in, out := &in.Subscriptions, &out.Subscriptions
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.ExtendedAttributes != nil {
		in, out := &in.ExtendedAttributes, &out.ExtendedAttributes
		*out = make([]byte, len(*in))
		copy(*out, *in)
	}
	if in.CheckHooks != nil {
		in, out := &in.CheckHooks, &out.CheckHooks
		*out = make([]HookList, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Subdue != nil {
		in, out := &in.Subdue, &out.Subdue
		*out = new(TimeWindowWhen)
		(*in).DeepCopyInto(*out)
	}
	if in.ProxyRequests != nil {
		in, out := &in.ProxyRequests, &out.ProxyRequests
		*out = new(ProxyRequests)
		(*in).DeepCopyInto(*out)
	}
	if in.OutputMetricHandlers != nil {
		in, out := &in.OutputMetricHandlers, &out.OutputMetricHandlers
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.EnvVars != nil {
		in, out := &in.EnvVars, &out.EnvVars
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	in.SensuMetadata.DeepCopyInto(&out.SensuMetadata)
	in.Validation.DeepCopyInto(&out.Validation)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SensuCheckConfigSpec.
func (in *SensuCheckConfigSpec) DeepCopy() *SensuCheckConfigSpec {
	if in == nil {
		return nil
	}
	out := new(SensuCheckConfigSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SensuCheckConfigStatus) DeepCopyInto(out *SensuCheckConfigStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SensuCheckConfigStatus.
func (in *SensuCheckConfigStatus) DeepCopy() *SensuCheckConfigStatus {
	if in == nil {
		return nil
	}
	out := new(SensuCheckConfigStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SensuCluster) DeepCopyInto(out *SensuCluster) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SensuCluster.
func (in *SensuCluster) DeepCopy() *SensuCluster {
	if in == nil {
		return nil
	}
	out := new(SensuCluster)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *SensuCluster) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SensuClusterList) DeepCopyInto(out *SensuClusterList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]SensuCluster, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SensuClusterList.
func (in *SensuClusterList) DeepCopy() *SensuClusterList {
	if in == nil {
		return nil
	}
	out := new(SensuClusterList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *SensuClusterList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SensuClusterRef) DeepCopyInto(out *SensuClusterRef) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SensuClusterRef.
func (in *SensuClusterRef) DeepCopy() *SensuClusterRef {
	if in == nil {
		return nil
	}
	out := new(SensuClusterRef)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SensuHandler) DeepCopyInto(out *SensuHandler) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SensuHandler.
func (in *SensuHandler) DeepCopy() *SensuHandler {
	if in == nil {
		return nil
	}
	out := new(SensuHandler)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *SensuHandler) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SensuHandlerList) DeepCopyInto(out *SensuHandlerList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]SensuHandler, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SensuHandlerList.
func (in *SensuHandlerList) DeepCopy() *SensuHandlerList {
	if in == nil {
		return nil
	}
	out := new(SensuHandlerList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *SensuHandlerList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SensuHandlerSpec) DeepCopyInto(out *SensuHandlerSpec) {
	*out = *in
	out.Socket = in.Socket
	if in.Handlers != nil {
		in, out := &in.Handlers, &out.Handlers
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Filters != nil {
		in, out := &in.Filters, &out.Filters
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.EnvVars != nil {
		in, out := &in.EnvVars, &out.EnvVars
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.RuntimeAssets != nil {
		in, out := &in.RuntimeAssets, &out.RuntimeAssets
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	in.SensuMetadata.DeepCopyInto(&out.SensuMetadata)
	in.Validation.DeepCopyInto(&out.Validation)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SensuHandlerSpec.
func (in *SensuHandlerSpec) DeepCopy() *SensuHandlerSpec {
	if in == nil {
		return nil
	}
	out := new(SensuHandlerSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SensuHandlerStatus) DeepCopyInto(out *SensuHandlerStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SensuHandlerStatus.
func (in *SensuHandlerStatus) DeepCopy() *SensuHandlerStatus {
	if in == nil {
		return nil
	}
	out := new(SensuHandlerStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SensuRestore) DeepCopyInto(out *SensuRestore) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SensuRestore.
func (in *SensuRestore) DeepCopy() *SensuRestore {
	if in == nil {
		return nil
	}
	out := new(SensuRestore)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *SensuRestore) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SensuRestoreList) DeepCopyInto(out *SensuRestoreList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]SensuRestore, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SensuRestoreList.
func (in *SensuRestoreList) DeepCopy() *SensuRestoreList {
	if in == nil {
		return nil
	}
	out := new(SensuRestoreList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *SensuRestoreList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StaticTLS) DeepCopyInto(out *StaticTLS) {
	*out = *in
	if in.Member != nil {
		in, out := &in.Member, &out.Member
		*out = new(MemberSecret)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StaticTLS.
func (in *StaticTLS) DeepCopy() *StaticTLS {
	if in == nil {
		return nil
	}
	out := new(StaticTLS)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TLSPolicy) DeepCopyInto(out *TLSPolicy) {
	*out = *in
	if in.Static != nil {
		in, out := &in.Static, &out.Static
		*out = new(StaticTLS)
		(*in).DeepCopyInto(*out)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TLSPolicy.
func (in *TLSPolicy) DeepCopy() *TLSPolicy {
	if in == nil {
		return nil
	}
	out := new(TLSPolicy)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TimeWindowDays) DeepCopyInto(out *TimeWindowDays) {
	*out = *in
	if in.All != nil {
		in, out := &in.All, &out.All
		*out = make([]*TimeWindowTimeRange, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(TimeWindowTimeRange)
				**out = **in
			}
		}
	}
	if in.Sunday != nil {
		in, out := &in.Sunday, &out.Sunday
		*out = make([]*TimeWindowTimeRange, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(TimeWindowTimeRange)
				**out = **in
			}
		}
	}
	if in.Monday != nil {
		in, out := &in.Monday, &out.Monday
		*out = make([]*TimeWindowTimeRange, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(TimeWindowTimeRange)
				**out = **in
			}
		}
	}
	if in.Tuesday != nil {
		in, out := &in.Tuesday, &out.Tuesday
		*out = make([]*TimeWindowTimeRange, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(TimeWindowTimeRange)
				**out = **in
			}
		}
	}
	if in.Wednesday != nil {
		in, out := &in.Wednesday, &out.Wednesday
		*out = make([]*TimeWindowTimeRange, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(TimeWindowTimeRange)
				**out = **in
			}
		}
	}
	if in.Thursday != nil {
		in, out := &in.Thursday, &out.Thursday
		*out = make([]*TimeWindowTimeRange, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(TimeWindowTimeRange)
				**out = **in
			}
		}
	}
	if in.Friday != nil {
		in, out := &in.Friday, &out.Friday
		*out = make([]*TimeWindowTimeRange, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(TimeWindowTimeRange)
				**out = **in
			}
		}
	}
	if in.Saturday != nil {
		in, out := &in.Saturday, &out.Saturday
		*out = make([]*TimeWindowTimeRange, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(TimeWindowTimeRange)
				**out = **in
			}
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TimeWindowDays.
func (in *TimeWindowDays) DeepCopy() *TimeWindowDays {
	if in == nil {
		return nil
	}
	out := new(TimeWindowDays)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TimeWindowTimeRange) DeepCopyInto(out *TimeWindowTimeRange) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TimeWindowTimeRange.
func (in *TimeWindowTimeRange) DeepCopy() *TimeWindowTimeRange {
	if in == nil {
		return nil
	}
	out := new(TimeWindowTimeRange)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TimeWindowWhen) DeepCopyInto(out *TimeWindowWhen) {
	*out = *in
	in.Days.DeepCopyInto(&out.Days)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TimeWindowWhen.
func (in *TimeWindowWhen) DeepCopy() *TimeWindowWhen {
	if in == nil {
		return nil
	}
	out := new(TimeWindowWhen)
	in.DeepCopyInto(out)
	return out
}
