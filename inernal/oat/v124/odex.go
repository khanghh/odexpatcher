package v124

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"odexpatcher/inernal/utils"
)

var OatMagic = [4]byte{0x6f, 0x61, 0x74, 0x0a}
var OatVersion124 = [4]byte{0x31, 0x32, 0x34, 0x00}

const HeaderOffset = 4096

type OdexFile struct {
	*OdexHeader
	OatDexFiles     []*OatDex
	OatDexCheckSums []checksumRecord
	buff            []byte
}

type checksumRecord struct {
	Offset   uint32
	CheckSum uint32
}

type OdexHeader struct {
	OatVersion                       [4]byte
	AdlerChecksum                    uint32
	InstructionSet                   uint32
	InstructionSetFeaturesBitmap     uint32
	DexFileCount                     uint32
	ExecutableOffset                 uint32
	I2IBridgeOffset                  uint32
	I2CBridgeOffset                  uint32
	JNIdlsymLookupOffset             uint32
	QuickGenericJNITrampolineOffset  uint32
	QuickImtConflictTrampolineOffset uint32
	QuickResolutionTrampolineOffset  uint32
	QuickToInterpreterBridgeOffset   uint32
	ImagePatchDelta                  uint32
	ImageFileLocationOatChecksum     uint32
	ImageFileLocationOatDataBegin    uint32
	KeyValueStoreSize                uint32
}

func (header *OdexHeader) PrintInfo() {
	fmt.Println("======================= ODEX Header ======================")
	fmt.Printf("OatVersion: %s\n", string(header.OatVersion[:]))
	fmt.Printf("AdlerChecksum: %x\n", header.AdlerChecksum)
	fmt.Printf("InstructionSet: %d\n", header.InstructionSet)
	fmt.Printf("InstructionSetFeaturesBitmap: %d\n", header.InstructionSetFeaturesBitmap)
	fmt.Printf("DexFileCount: %d\n", header.DexFileCount)
	fmt.Printf("ExecutableOffset: %d\n", header.ExecutableOffset)
}

type OatDex struct {
	LocationSize       uint32
	LocationData       []uint8
	Checksum           uint32
	FileOffset         uint32
	ClassOffsetsOffset uint32
	LookupTableOffset  uint32
}

func (odex *OdexFile) PrintInfo() {
	odex.OdexHeader.PrintInfo()
	fmt.Println("OatDex files:")
	for idx := uint32(0); idx < odex.DexFileCount; idx++ {
		oatDex := odex.OatDexFiles[idx]
		fmt.Printf("  OatDex Num: %d\n", idx)
		fmt.Printf("    LocationSize: %d\n", oatDex.LocationSize)
		fmt.Printf("    LocationData: %x\n", oatDex.LocationData)
		fmt.Printf("    Checksum: %x\n", oatDex.Checksum)
		fmt.Printf("    FileOffset: %x\n", oatDex.FileOffset)
		fmt.Printf("    ClassOffsetsOffset: %x\n", oatDex.ClassOffsetsOffset)
		fmt.Printf("    LookupTableOffset: %x\n", oatDex.LookupTableOffset)
	}
}

func (odex *OdexFile) SaveToFile(filePath string) error {
	return ioutil.WriteFile(filePath, odex.buff, 0644)
}

func ParseOdex(buff []byte) (*OdexFile, error) {
	odexHeader, err := parseOdexHeader(buff)
	if err != nil {
		return nil, err
	}
	oatDexFiles, err := parseOatDexList(buff, odexHeader)
	if err != nil {
		return nil, err
	}
	return &OdexFile{
		OdexHeader:  odexHeader,
		OatDexFiles: oatDexFiles,
		buff:        buff,
	}, nil
}

func (odex *OdexFile) PatchOatDexChecksums(newChecksums []uint32) error {
	offset := HeaderOffset + 18*4 + odex.KeyValueStoreSize
	writer := utils.NewBinaryWriter(odex.buff[offset:])
	for idx := uint32(0); idx < odex.DexFileCount; idx++ {
		oatDexFile := odex.OatDexFiles[idx]
		writer.Seek(4+int64(oatDexFile.LocationSize), io.SeekCurrent)
		binary.Write(writer, binary.LittleEndian, newChecksums[idx])
		writer.Seek(4*3, io.SeekCurrent)
	}
	return nil
}

func parseOatDexList(odexData []byte, odexHeader *OdexHeader) ([]*OatDex, error) {
	oatDexFilesOffset := HeaderOffset + 18*4 + odexHeader.KeyValueStoreSize
	buff := bytes.NewBuffer(odexData[oatDexFilesOffset:])
	oatDexFiles := make([]*OatDex, odexHeader.DexFileCount)
	for idx := uint32(0); idx < odexHeader.DexFileCount; idx++ {
		var locationSize uint32
		binary.Read(buff, binary.LittleEndian, &locationSize)
		locationData := make([]byte, locationSize)
		binary.Read(buff, binary.LittleEndian, &locationData)
		var checksum uint32
		binary.Read(buff, binary.LittleEndian, &checksum)
		var fileOffset uint32
		binary.Read(buff, binary.LittleEndian, &fileOffset)
		var classOffsetsOffset uint32
		binary.Read(buff, binary.LittleEndian, &classOffsetsOffset)
		var lookupTableOffset uint32
		binary.Read(buff, binary.LittleEndian, &lookupTableOffset)
		oatDex := &OatDex{
			LocationSize:       locationSize,
			LocationData:       locationData,
			Checksum:           checksum,
			FileOffset:         fileOffset,
			ClassOffsetsOffset: classOffsetsOffset,
			LookupTableOffset:  lookupTableOffset,
		}
		oatDexFiles[idx] = oatDex
	}
	return oatDexFiles, nil
}

func parseOdexHeader(odexData []byte) (*OdexHeader, error) {
	buff := bytes.NewBuffer(odexData)
	buff.Next(HeaderOffset)
	var magic [4]byte
	err := binary.Read(buff, binary.LittleEndian, &magic)
	if err != nil {
		return nil, err
	}
	if magic != OatMagic {
		return nil, errors.New("wrong magic header")
	}
	var odexHeader OdexHeader
	err = binary.Read(buff, binary.LittleEndian, &odexHeader)
	if err != nil {
		return nil, err
	}
	return &odexHeader, nil
}
