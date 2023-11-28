//go:build !ignore_autogenerated

/*
Copyright 2023.

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

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PodDestroyer) DeepCopyInto(out *PodDestroyer) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PodDestroyer.
func (in *PodDestroyer) DeepCopy() *PodDestroyer {
	if in == nil {
		return nil
	}
	out := new(PodDestroyer)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *PodDestroyer) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PodDestroyerList) DeepCopyInto(out *PodDestroyerList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]PodDestroyer, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PodDestroyerList.
func (in *PodDestroyerList) DeepCopy() *PodDestroyerList {
	if in == nil {
		return nil
	}
	out := new(PodDestroyerList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *PodDestroyerList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PodDestroyerSpec) DeepCopyInto(out *PodDestroyerSpec) {
	*out = *in
	in.Selector.DeepCopyInto(&out.Selector)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PodDestroyerSpec.
func (in *PodDestroyerSpec) DeepCopy() *PodDestroyerSpec {
	if in == nil {
		return nil
	}
	out := new(PodDestroyerSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PodDestroyerStatus) DeepCopyInto(out *PodDestroyerStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PodDestroyerStatus.
func (in *PodDestroyerStatus) DeepCopy() *PodDestroyerStatus {
	if in == nil {
		return nil
	}
	out := new(PodDestroyerStatus)
	in.DeepCopyInto(out)
	return out
}