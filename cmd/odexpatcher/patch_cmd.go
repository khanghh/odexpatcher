package main

import (
	"io/ioutil"
	"log"
	v124 "odexpatcher/inernal/oat/v124"
	"path/filepath"
	"strconv"
	"strings"
)

type PatchCmd struct {
	OdexFile      string `long:"odex" description:"odex file" default:"base.odex"`
	VdexFile      string `long:"vdex" description:"vdex file" default:"base.vdex"`
	ChecksumsFile string `long:"crc32" description:"crc32 checksums to patch"`
	Output        string `long:"out" description:"Output directory"`
}

//dex2oat --dex-file=$dexFile
// --dex-location=$dexLocation
// --oat-file=$oatFile
// --instruction-set=${Art.ISA}
// --instruction-set-variant=${Art.ISA_VARIANT}
// --instruction-set-features=${Art.ISA_FEATURES}

func readChecksums(filename string) []uint32 {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalln("Could not read checksum file:", err)
	}
	checksums := make([]uint32, 0)
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		val, err := strconv.ParseInt(line, 16, 64)
		if err != nil {
			log.Fatalln("Could not parse checksum file:", err)
		}
		checksums = append(checksums, uint32(val))
	}
	return checksums
}

func (cmd *PatchCmd) Execute(args []string) error {
	odexData, err := ioutil.ReadFile(cmd.OdexFile)
	if err != nil {
		log.Fatalln("Could not read odex file:", err)
	}
	odex, err := v124.ParseOdex(odexData)
	if err != nil {
		log.Fatalln("Could not parse odex file:", err)
	}
	vdexData, err := ioutil.ReadFile(cmd.VdexFile)
	if err != nil {
		log.Fatalln("Could not read vdex file:", err)
	}
	vdex, err := v124.ParseVdex(vdexData)
	if err != nil {
		log.Fatalln("Could not parse vdex file:", err)
	}

	checksums := readChecksums(cmd.ChecksumsFile)
	odex.PatchOatDexChecksums(checksums)
	vdex.PatchVDexChecksums(checksums)

	odexOut := filepath.Join(cmd.Output, filepath.Base(cmd.OdexFile))
	vdexOut := filepath.Join(cmd.Output, filepath.Base(cmd.VdexFile))

	odex.SaveToFile(odexOut)
	vdex.SaveToFile(vdexOut)

	return nil
}
