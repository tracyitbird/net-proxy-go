package common

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"log"
	"net"
)

type Package struct {
	len       [4]byte
	headerLen [4]byte
	bodyLen   [4]byte
	header    []byte
	body      []byte
}

func NewPackage() *Package {
	return &Package{}
}

func (pkg *Package) GetHeader() []byte {
	return pkg.header
}

func (pkg *Package) GetBody() []byte {
	return pkg.body
}

//readWithHeader
func (pkg *Package) ReadWithHeader(reader net.Conn) (err error) {
	sizeBuf := make([]byte, 4+4+4)
	n, err := io.ReadAtLeast(reader, sizeBuf, 12)

	if n < 0 || err != nil {
		return errors.New("read error...")
	}

	log.Printf("read with header read size buf len = %v", n)
	log.Printf("size buf = %v", sizeBuf)

	len := BytesToInt(sizeBuf[:4])
	headerLen := BytesToInt(sizeBuf[4:8])
	bodyLen := BytesToInt(sizeBuf[8:12])

	total := make([]byte, headerLen + bodyLen)
	n2, err := io.ReadAtLeast(reader, total, headerLen + bodyLen)

	header := total[0:headerLen]
	body := total[headerLen:bodyLen]

	log.Printf("read with header, len = %v, headerLen = %v, bodyLen = %v", len, headerLen, bodyLen)

	copy(pkg.len[:], sizeBuf[:4])
	copy(pkg.headerLen[:], sizeBuf[4:8])
	copy(pkg.bodyLen[:], sizeBuf[8:12])
	pkg.header = header
	pkg.body = body

	if n2 < 0 || err != nil {
		return errors.New("read error...")
	}
	return nil
}

//readWithoutHeader
func (pkg *Package) ReadWithoutHeader(reader io.Reader) (err error) {
	buf := make([]byte, 0, 100 * 1024)
	n, err := reader.Read(buf)

	pkg.body = buf
	pkg.header = make([]byte, 0)

	len := 0 + n + 12
	headerLen := 0
	bodyLen := n

	copy(pkg.len[:], IntToBytes(len)[:])
	copy(pkg.headerLen[:], IntToBytes(headerLen)[:])
	copy(pkg.bodyLen[:], IntToBytes(bodyLen)[:])

	if err != nil {
		return err
	} else {
		return nil
	}
}

//writeWithHeader
//func (pkg *Package) WriteWithHeader(writer io.Writer) (err error) {
//	total := make([]byte, 0, 4+4+4+toInt(pkg.len[0:4]))
//
//	copy(total, pkg.len[:])
//	copy(total[4:], pkg.headerLen[:])
//	copy(total[4+4:], pkg.bodyLen[:])
//	copy(total[4+4+4:], pkg.header[:])
//	copy(total[4+4+4+len(pkg.header):], pkg.body[:])
//	n, err := writer.Write(total)
//	if n < 0 || err != nil {
//		return errors.New("write error...")
//	} else {
//		return nil
//	}
//}
//
////writeWithoutHeader
//func (pkg *Package) WriteWithoutHeader(writer io.Writer) (err error) {
//	n, err := writer.Write(pkg.body)
//	if n < 0 || err != nil {
//		return errors.New("write error...")
//	} else {
//		return nil
//	}
//}

//readfully

//value of
func (pkg *Package) ValueOf(header []byte, body []byte) {
	headerLen := len(header)
	bodyLen := len(body)
	len := headerLen + bodyLen + 12

	copy(pkg.len[:], IntToBytes(len)[:])
	copy(pkg.headerLen[:], IntToBytes(headerLen)[:])
	copy(pkg.bodyLen[:], IntToBytes(bodyLen)[:])

	pkg.header = header
	pkg.body = body
}

//to bytes
func (pkg *Package) ToBytes() []byte {
	totalBytes := make([]byte, 4+4+4+len(pkg.header)+len(pkg.body))

	copy(totalBytes[:4], pkg.len[:])
	copy(totalBytes[4:8], pkg.headerLen[:])
	copy(totalBytes[8:12], pkg.bodyLen[:])
	copy(totalBytes[12:12+len(pkg.header)], pkg.header[:])
	copy(totalBytes[12+len(pkg.header):12+len(pkg.header)+len(pkg.body)], pkg.body[:])

	log.Printf("len = %v, headerLen = %v, bodyLen = %v, header = %v, body = %v", pkg.len, pkg.headerLen, pkg.bodyLen, pkg.header, pkg.body)
	return totalBytes
}

func IntToBytes(n int) []byte {
	tmp := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, &tmp)
	return bytesBuffer.Bytes()
}

func BytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)
	var tmp int32
	binary.Read(bytesBuffer, binary.BigEndian, &tmp)
	return int(tmp)
}
