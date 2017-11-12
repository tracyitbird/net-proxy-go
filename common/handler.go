package common

type Handler interface {
	handle(pkg *Package) (newPkg Package)
}
