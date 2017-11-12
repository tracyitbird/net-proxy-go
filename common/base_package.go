package common

import (
	"io"
	"errors"
)

type BasePackage struct {
	len [4]byte
	headerLen [4]byte
	bodyLen [4]byte
	header []byte
	body []byte
}

//readWithHeader
func (pkg *BasePackage) ReadWithHeader (reader io.Reader) (err error){
	sizeBuf := make([]byte, 0, 4 + 4 + 4)
	n, err := io.ReadAtLeast(reader, sizeBuf, len(sizeBuf))

	if n < 0 || err != nil {
		return errors.New("read error...")
	}

	len := int(n)
	headerLen := int(sizeBuf[4:8])
	bodyLen := int(sizeBuf[8:12])

	total := make([]byte, 0, len)
	n2, err := io.ReadAtLeast(reader, total, len)


	header := total[0:headerLen]
	body := total[headerLen:bodyLen]

	pkg.len = [4]byte(len)
	pkg.headerLen = [4]byte(headerLen)
	pkg.bodyLen = [4]byte(bodyLen)
	pkg.header = header
	pkg.body = body

	if n2 < 0 || err != nil {
		return errors.New("read error...")
	}
	return nil
}

//readWithoutHeader
func (pkg *BasePackage) ReadWithoutHeader (reader io.Reader) (err error) {
	buf := make([]byte, 0, 100 * 1024)
	n, err := reader.Read(buf)


	pkg.body = buf
	pkg.bodyLen = [4]byte(n)

	pkg.header = make([]byte, 0, 0)
	pkg.headerLen = [4]byte(0)

	pkg.len = [4]byte(n + 0)
	if err != nil{
		return err
	} else {
		return nil
	}
}

//writeWithHeader
func (pkg *BasePackage) WriteWithHeader (writer io.Writer) (err error){
	total := make([]byte, 0, 4 + 4 + 4 + int(pkg.len))
	copy(total, []byte(pkg.len))
	copy(total[4:], []byte(pkg.headerLen))
	copy(total[4+4:], []byte(pkg.bodyLen))
	copy(total[4+4+4:], pkg.header)
	copy(total[4+4+4+len(pkg.header):], pkg.body)
	n, err := writer.Write(total)
	if n < 0 || err != nil {
		return errors.New("write error...")
	} else {
		return nil
	}
}

//writeWithoutHeader
func (pkg *BasePackage) WriteWithoutHeader (writer io.Writer) (err error){
	n, err := writer.Write(pkg.body)
	if n < 0 || err != nil {
		return errors.New("write error...")
	} else {
		return nil
	}
}

//readfully