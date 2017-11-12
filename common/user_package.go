package common

type UserPackage struct {
	basePkg BasePackage
}

//get address and port, this info only use in first connected package
func (userPkg *UserPackage) getAddressAndPort() (addr string, port int) {
	header := userPkg.basePkg.header
	addrLen := int(header[0:4])
	addr = string(header[4:addrLen])
	port = int(header[4 + addrLen :])
	return addr, port
}