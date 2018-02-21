// Copyright 2018 rootkiwi
//
// screen_share_remote_go is licensed under GNU General Public License 3 or later.
//
// See LICENSE for more details.

// +build ignore

package main

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

const version = "0.1.0"

func printUsage() {
	fmt.Println("Usage building screen_share_remote_go:")
	fmt.Println()
	fmt.Println("No argument will build for running platform, outputting binary in bin/dev")
	fmt.Println("Argument 'release' will build and archive for all platforms, outputting in bin/release")
}

func main() {
	args := os.Args[1:]
	switch len(args) {
	case 0:
		buildDev()
		return
	case 1:
		if args[0] == "release" {
			buildRelease()
			return
		}
	}
	printUsage()
}

func buildDev() {
	fileName := fmt.Sprintf("screen_share_remote_go-%s", version)
	relativeBinPath := "bin/dev/" + fileName
	fmt.Println(relativeBinPath)
	cmd := exec.Command("packr", "build", "-o", relativeBinPath, "cmd/screen_share_remote_go/main.go")
	if err := cmd.Run(); err != nil {
		log.Fatalf("error running build: %v\n", err)
	}
}

type platform struct {
	os     string
	arches []string
	archiver
}

type archiver interface {
	archive(path string)
}

func buildRelease() {
	var tarGz = new(tarGzArchiver)
	var zip = new(zipArchiver)

	platforms := []*platform{
		{"linux", []string{"amd64", "386", "arm", "arm64"}, tarGz},
		{"darwin", []string{"amd64", "386"}, tarGz},
		{"windows", []string{"amd64", "386"}, zip},
	}

	relativeBinDirPath := "bin/release/" + version
	if err := os.RemoveAll(relativeBinDirPath); err != nil {
		log.Fatalf("error removing directory: %s: %v\n", relativeBinDirPath, err)
	}
	fmt.Println(relativeBinDirPath)

	for _, p := range platforms {
		for _, arch := range p.arches {
			fmt.Printf("%s-%s\n", p.os, arch)
			fileName := fmt.Sprintf("screen_share_remote_go-%s-%s-%s", version, p.os, arch)
			relativeBinPath := filepath.Join(relativeBinDirPath, fileName)
			cmd := exec.Command("env", "GOOS="+p.os, "GOARCH="+arch,
				"packr", "build", "-o", relativeBinPath, "cmd/screen_share_remote_go/main.go")
			if err := cmd.Run(); err != nil {
				log.Fatalf("error running build: %s-%s: %v\n", p.os, arch, err)
			}
			p.archive(relativeBinPath)
			if err := os.Remove(relativeBinPath); err != nil {
				log.Fatalf("error removing bin: %s: %v\n", relativeBinPath, err)
			}
		}
	}

	os.Chdir(relativeBinDirPath)
	filesToSum, err := filepath.Glob("*")
	if err != nil {
		log.Fatalln(err)
	}
	cmd := exec.Command("sha256sum", filesToSum...)
	sumFile, err := os.Create("sha256sum.txt")
	if err != nil {
		log.Fatalln(err)
	}
	defer sumFile.Close()
	cmd.Stdout = sumFile
	if err := cmd.Run(); err != nil {
		log.Fatalln(err)
	}
}

type tarGzArchiver struct{}

func (gz *tarGzArchiver) archive(path string) {
	inFile, err := os.Open(path)
	if err != nil {
		log.Fatalln(err)
	}
	defer inFile.Close()
	outFile, err := os.Create(path + ".tar.gz")
	if err != nil {
		log.Fatalln(err)
	}
	defer outFile.Close()

	gw, err := gzip.NewWriterLevel(outFile, gzip.BestCompression)
	if err != nil {
		log.Fatalln(err)
	}
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()

	info, err := inFile.Stat()
	if err != nil {
		log.Fatalln(err)
	}
	_, fileName := filepath.Split(inFile.Name())
	h := &tar.Header{
		Name:    fileName,
		Size:    info.Size(),
		Mode:    int64(info.Mode()),
		ModTime: info.ModTime(),
	}

	err = tw.WriteHeader(h)
	if err != nil {
		log.Fatal(err)
	}
	_, err = io.Copy(tw, inFile)
	if err != nil {
		log.Fatalln(err)
	}
}

type zipArchiver struct{}

func (gz *zipArchiver) archive(path string) {
	inFile, err := os.Open(path)
	if err != nil {
		log.Fatalln(err)
	}
	defer inFile.Close()
	outFile, err := os.Create(path + ".zip")
	if err != nil {
		log.Fatalln(err)
	}
	defer outFile.Close()

	zw := zip.NewWriter(outFile)
	defer zw.Close()

	info, err := inFile.Stat()
	if err != nil {
		log.Fatalln(err)
	}

	h, err := zip.FileInfoHeader(info)
	if err != nil {
		log.Fatalln(err)
	}
	h.Method = zip.Deflate
	w, err := zw.CreateHeader(h)
	if err != nil {
		log.Fatalln(err)
	}
	_, err = io.Copy(w, inFile)
	if err != nil {
		log.Fatalln(err)
	}
}
