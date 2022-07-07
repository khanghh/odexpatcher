package v124

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

var DexMagic = [4]byte{0x64, 0x65, 0x78, 0x0a}

type DexFile struct {
	*DexHeader
	buf []byte
}

type DexHeader struct {
	DexVersion      [4]byte
	Alder32Checksum uint32
	Signature       [20]uint8
	FileSize        uint32
	HeaderSize      uint32
	EndianTag       uint32
	LinkSize        uint32
	LinkOffset      uint32
	MapOffset       uint32
	StringIdsSize   uint32
	StringIdsOffset uint32
	TypeIdsSize     uint32
	TypeIdsOffset   uint32
	ProtoIdsSize    uint32
	ProtoIdsOffset  uint32
	FieldIdsSize    uint32
	FieldIdsOffset  uint32
	MethodIdsSize   uint32
	MethodIdsOffset uint32
	ClassIdsSize    uint32
	ClassIdsOffset  uint32
	DataSize        uint32
	DataOffset      uint32
}

func (dex *DexFile) PrintInfo() {
	fmt.Println("================ DEX file ========================")
	fmt.Printf("DexVersion: %s\n", string(dex.DexVersion[:]))
	fmt.Printf("Alder32Checksum: %x\n", dex.Alder32Checksum)
	fmt.Printf("Signature: %x\n", dex.Signature)
	fmt.Printf("FileSize: %d\n", dex.FileSize)
	fmt.Printf("EndianTag: %x\n", dex.EndianTag)
	fmt.Printf("StringIdsSize: %d\n", dex.StringIdsSize)
	fmt.Printf("TypeIdsSize: %d\n", dex.TypeIdsSize)
	fmt.Printf("FieldIdsSize: %d\n", dex.FieldIdsSize)
	fmt.Printf("MethodIdsSize: %d\n", dex.MethodIdsSize)
	fmt.Printf("ClassIdsSize: %d\n", dex.ClassIdsSize)
}

func parseDexHeader(dexData []byte) (*DexHeader, error) {
	buff := bytes.NewBuffer(dexData)
	var magic [4]byte
	err := binary.Read(buff, binary.LittleEndian, &magic)
	if err != nil {
		return nil, err
	}
	if magic != DexMagic {
		return nil, errors.New("not a dex file")
	}
	var dexHeader DexHeader
	err = binary.Read(buff, binary.LittleEndian, &dexHeader)
	if err != nil {
		return nil, err
	}
	return &dexHeader, nil
}

func ParseDex(dexData []byte) (*DexFile, error) {
	dexHeader, err := parseDexHeader(dexData)
	if err != nil {
		return nil, err
	}
	return &DexFile{
		DexHeader: dexHeader,
		buf:       dexData,
	}, nil
}
