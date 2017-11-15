package common

type PackageHandler interface {
	Handle(pkg *Package) (newPkg Package)
}

type EncryptHandler struct {

}

type DecryptHandler struct {

}

type CompressHandler struct {

}

type DecompressHandler struct {
	
}
