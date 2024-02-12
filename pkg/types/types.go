package types

import "fmt"

type Registry int

const (
	RegistryUnknown = iota
	RegistryAWS
	RegistryGCP
	RegistryGeneric
)

func (p Registry) String() string {
	return [...]string{"unknown", "aws", "gcp", "generic"}[p]
}

func ParseRegistry(p string) (Registry, error) {
	switch p {
	case Registry(RegistryAWS).String():
		return RegistryAWS, nil
	case Registry(RegistryGCP).String():
		return RegistryGCP, nil
	case Registry(RegistryGeneric).String():
		return RegistryGeneric, nil
	}
	return RegistryUnknown, fmt.Errorf("unknown target registry string: '%s', defaulting to unknown", p)
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
	ImageCopyPolicyNone
)

func (p ImageCopyPolicy) String() string {
	return [...]string{"delayed", "immediate", "force", "none"}[p]
}

func ParseImageCopyPolicy(p string) (ImageCopyPolicy, error) {
	switch p {
	case ImageCopyPolicy(ImageCopyPolicyDelayed).String():
		return ImageCopyPolicyDelayed, nil
	case ImageCopyPolicy(ImageCopyPolicyImmediate).String():
		return ImageCopyPolicyImmediate, nil
	case ImageCopyPolicy(ImageCopyPolicyForce).String():
		return ImageCopyPolicyForce, nil
	case ImageCopyPolicy(ImageCopyPolicyNone).String():
		return ImageCopyPolicyNone, nil
	}
	return ImageCopyPolicyDelayed, fmt.Errorf("unknown image copy policy string: '%s', defaulting to delayed", p)
}
