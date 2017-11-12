package common

type PackageHandler interface {
	Handle(pkg *Package) (newPkg Package)
}
