package types

import "fmt"

type TargetRegistry int

const (
	TargetRegistryUnknown = iota
	TargetRegistryAws
	TargetRegistryGcp
)

func (p TargetRegistry) String() string {
	return [...]string{"unknown", "aws", "gcp"}[p]
}

func ParseTargetRegistry(p string) (TargetRegistry, error) {
	switch p {
	case TargetRegistry(TargetRegistryAws).String():
		return TargetRegistryAws, nil
	case TargetRegistry(TargetRegistryGcp).String():
		return TargetRegistryGcp, nil
	}
	return TargetRegistryUnknown, fmt.Errorf("unknown target registry string: '%s', defaulting to unknown", p)
}

type ImageSwapPolicy int

const (
	ImageSwapPolicyAlways = iota
	ImageSwapPolicyExists
)

func (p ImageSwapPolicy) String() string {
	return [...]string{"always", "exists"}[p]
}

func ParseImageSwapPolicy(p string) (ImageSwapPolicy, error) {
	switch p {
	case ImageSwapPolicy(ImageSwapPolicyAlways).String():
		return ImageSwapPolicyAlways, nil
	case ImageSwapPolicy(ImageSwapPolicyExists).String():
		return ImageSwapPolicyExists, nil
	}
	return ImageSwapPolicyExists, fmt.Errorf("unknown image swap policy string: '%s', defaulting to exists", p)
}

type ImageCopyPolicy int

const (
	ImageCopyPolicyDelayed = iota
	ImageCopyPolicyImmediate
	ImageCopyPolicyForce
)

func (p ImageCopyPolicy) String() string {
	return [...]string{"delayed", "immediate", "force"}[p]
}

func ParseImageCopyPolicy(p string) (ImageCopyPolicy, error) {
	switch p {
	case ImageCopyPolicy(ImageCopyPolicyDelayed).String():
		return ImageCopyPolicyDelayed, nil
	case ImageCopyPolicy(ImageCopyPolicyImmediate).String():
		return ImageCopyPolicyImmediate, nil
	case ImageCopyPolicy(ImageCopyPolicyForce).String():
		return ImageCopyPolicyForce, nil
	}
	return ImageCopyPolicyDelayed, fmt.Errorf("unknown image copy policy string: '%s', defaulting to delayed", p)
}
