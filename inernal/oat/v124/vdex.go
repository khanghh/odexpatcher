package v124

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io/ioutil"
	"odexpatcher/inernal/utils"
)

var VdexMagic = [4]byte{0x76, 0x64, 0x65, 0x78}
var VdexVersion06 = [4]byte{0x30, 0x30, 0x36, 0x00}
var VdexHeaderSize = 24

type VdexFile struct {
	*VdexHeader
	DexCheckSums []uint32
	buff         []byte
}

type VdexHeader struct {
	VdexVersion        [4]byte
	DexFileCount       uint32
	DexSize            uint32
	VerifierDepsSize   uint32
	QuickeningInfoSize uint32
}

func (header *VdexHeader) PrintInfo() {
	fmt.Println("======================= VDEX Header ======================")
	fmt.Printf("VdexVersion: %s\n", string(header.VdexVersion[:]))
	fmt.Printf("DexFileCount: %d\n", header.DexFileCount)
	fmt.Printf("DexSize: %d\n", header.DexSize)
	fmt.Printf("VerifierDepsSize: %d\n", header.VerifierDepsSize)
	fmt.Printf("QuickeningInfoSize: %d\n", header.QuickeningInfoSize)
}

func (vdex *VdexFile) PrintInfo() {
	vdex.VdexHeader.PrintInfo()
	fmt.Println("Dex checksums:")
	for idx, checksum := range vdex.DexCheckSums {
		fmt.Printf("  Dex %d: %x\n", idx, checksum)
	}
}

func (vdex *VdexFile) SaveToFile(filePath string) error {
	return ioutil.WriteFile(filePath, vdex.buff, 0644)
}

func (vdex *VdexFile) PatchVDexChecksums(newchecksums []uint32) error {
	writer := utils.NewBinaryWriter(vdex.buff[VdexHeaderSize:])
	for idx := uint32(0); idx < vdex.DexFileCount; idx++ {
		binary.Write(writer, binary.LittleEndian, newchecksums[idx])
	}
	return nil
}

func parseDexChecksums(vdexData []byte, header *VdexHeader) ([]uint32, error) {
	buff := bytes.NewBuffer(vdexData)
	buff.Next(VdexHeaderSize)
	checksums := make([]uint32, header.DexFileCount)
	for idx := uint32(0); idx < header.DexFileCount; idx++ {
		var checksum uint32
		err := binary.Read(buff, binary.LittleEndian, &checksum)
		if err != nil {
			return checksums, nil
		}
		checksums[idx] = checksum
	}
	return checksums, nil
}

func ParseVdex(buff []byte) (*VdexFile, error) {
	vdexHeader, err := parseVdexHeader(buff)
	if err != nil {
		return nil, err
	}
	checksums, err := parseDexChecksums(buff, vdexHeader)
	if err != nil {
		return nil, err
	}

	return &VdexFile{
		VdexHeader:   vdexHeader,
		DexCheckSums: checksums,
		buff:         buff,
	}, nil
}

func parseVdexHeader(odexData []byte) (*VdexHeader, error) {
	buff := bytes.NewBuffer(odexData)
	var magic [4]byte
	err := binary.Read(buff, binary.LittleEndian, &magic)
	if err != nil {
		return nil, err
	}
	if magic != VdexMagic {
		return nil, errors.New("wrong magic header")
	}
	var vdexHeader VdexHeader
	err = binary.Read(buff, binary.LittleEndian, &vdexHeader)
	if err != nil {
		return nil, err
	}
	return &vdexHeader, nil
}
