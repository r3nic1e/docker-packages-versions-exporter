package packages

import "fmt"

var packageGetters map[string]PackageGetter

type PackageGetter interface {
	GetPackages() []string
}

type UnknownOSError struct{}

func (UnknownOSError) Error() string {
	return fmt.Sprintf("Unknown OS")
}

func NewPackageGetter(os string) (PackageGetter, error) {
	if packageGetter, ok := packageGetters[os]; ok {
		return packageGetter, nil
	} else {
		return nil, UnknownOSError{}
	}
}

func registerPackageGetter(packageGetter PackageGetter, os string) {
	if packageGetters == nil {
		packageGetters = make(map[string]PackageGetter)
	}
	packageGetters[os] = packageGetter
}
