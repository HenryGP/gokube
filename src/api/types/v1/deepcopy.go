package v1

import "k8s.io/apimachinery/pkg/runtime"

// DeepCopyInto copies all properties of this object into another object of the
// same type that is provided as a pointer.
func (in *MongoDB) DeepCopyInto(out *MongoDB) {
	out.TypeMeta = in.TypeMeta
	out.ObjectMeta = in.ObjectMeta
	out.Spec = MongoSpec{
		Credentials: in.Spec.Credentials,
		Project:     in.Spec.Project,
		Version:     in.Spec.Version,
		Type:        in.Spec.Type,
	}
}

// DeepCopyObject returns a generically typed copy of an object
func (in *MongoDB) DeepCopyObject() runtime.Object {
	out := MongoDB{}
	in.DeepCopyInto(&out)
	return &out
}

// DeepCopyObject returns a generically typed copy of an object
func (in *MongoDBList) DeepCopyObject() runtime.Object {
	out := MongoDBList{}
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta

	if in.Items != nil {
		out.Items = make([]MongoDB, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&out.Items[i])
		}
	}
	return &out
}
