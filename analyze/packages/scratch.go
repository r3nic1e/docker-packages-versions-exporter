package packages

type ScratchPackageGetter struct{}

func init() {
	registerPackageGetter(&ScratchPackageGetter{}, "scratch")
}

func (*ScratchPackageGetter) GetPackages() (pkgs []string) {
	return
}
