package common

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"strconv"
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
func (pkg *Package) ReadWithHeader(reader io.Reader) (err error) {
	sizeBuf := make([]byte, 0, 4+4+4)
	n, err := io.ReadAtLeast(reader, sizeBuf, len(sizeBuf))

	if n < 0 || err != nil {
		return errors.New("read error...")
	}

	len := toInt(sizeBuf[:4])
	headerLen := toInt(sizeBuf[4:8])
	bodyLen := toInt(sizeBuf[8:12])

	total := make([]byte, headerLen + bodyLen)
	n2, err := io.ReadAtLeast(reader, total, headerLen + bodyLen)

	header := total[0:headerLen]
	body := total[headerLen:bodyLen]

	copy([]byte(strconv.Itoa(len))[0:4], pkg.len[:])

	copy([]byte(strconv.Itoa(len))[0:4], pkg.len[:])
	copy([]byte(strconv.Itoa(headerLen))[0:4], pkg.headerLen[:])
	copy([]byte(strconv.Itoa(headerLen))[0:4], pkg.bodyLen[:])

	pkg.header = header
	pkg.body = body

	if n2 < 0 || err != nil {
		return errors.New("read error...")
	}
	return nil
}

//readWithoutHeader
func (pkg *Package) ReadWithoutHeader(reader io.Reader) (err error) {
	buf := make([]byte, 0, 100*1024)
	n, err := reader.Read(buf)

	pkg.body = buf

	pkg.header = make([]byte, 0)
	copy([]byte(strconv.Itoa(0))[0:4], pkg.headerLen[:])
	copy([]byte(strconv.Itoa(n))[0:4], pkg.bodyLen[:])
	copy([]byte(strconv.Itoa(n+0))[0:4], pkg.len[:])

	if err != nil {
		return err
	} else {
		return nil
	}
}

//writeWithHeader
func (pkg *Package) WriteWithHeader(writer io.Writer) (err error) {
	total := make([]byte, 0, 4+4+4+toInt(pkg.len[0:4]))

	copy(total, pkg.len[:])
	copy(total[4:], pkg.headerLen[:])
	copy(total[4+4:], pkg.bodyLen[:])
	copy(total[4+4+4:], pkg.header[:])
	copy(total[4+4+4+len(pkg.header):], pkg.body[:])
	n, err := writer.Write(total)
	if n < 0 || err != nil {
		return errors.New("write error...")
	} else {
		return nil
	}
}

//writeWithoutHeader
func (pkg *Package) WriteWithoutHeader(writer io.Writer) (err error) {
	n, err := writer.Write(pkg.body)
	if n < 0 || err != nil {
		return errors.New("write error...")
	} else {
		return nil
	}
}

//readfully

//value of
func (pkg *Package) ValueOf(header []byte, body []byte) {
	headerLen := len(header)
	bodyLen := len(body)

	len := headerLen + bodyLen

	byteBuf := bytes.NewBuffer(make([]byte, 4))
	binary.Write(byteBuf, binary.LittleEndian, len)
	copy(pkg.len[:], byteBuf.Bytes()[0:4])

	byteBuf = bytes.NewBuffer(make([]byte, 4))
	binary.Write(byteBuf, binary.LittleEndian, headerLen)
	copy(pkg.header[:], byteBuf.Bytes()[0:4])

	byteBuf = bytes.NewBuffer(make([]byte, 4))
	binary.Write(byteBuf, binary.LittleEndian, bodyLen)
	copy(pkg.body[:], byteBuf.Bytes()[0:4])

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

	return totalBytes
}

func toInt(bytes []byte) (int) {
	val := string(bytes)
	intValue, _ := strconv.Atoi(val)
	return intValue
}
