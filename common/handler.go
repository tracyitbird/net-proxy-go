package common

type Handler interface {
	handle(pkg *BasePackage) (newPkg BasePackage)
}