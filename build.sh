#!/bin/bash

version="0.1.0"

for os in linux darwin; do
    for arch in amd64 386; do
        echo "building $os-$arch"
        file="screen_share_remote_go-$version-$os-$arch"
        env GOOS=$os GOARCH=$arch packr build -o build/$version/$file main.go
        cd build/$version
        tar -czf $file.tar.gz $file
        rm $file
        cd ../..
    done
done

for os in windows; do
    for arch in amd64 386; do
        echo "building $os-$arch"
        file="screen_share_remote_go-$version-$os-$arch"
        env GOOS=$os GOARCH=$arch packr build -o build/$version/$file.exe main.go
        cd build/$version
        zip -q $file.zip $file.exe
        rm $file.exe
        cd ../..
    done
done

cd build/$version
sha256sum * > sha256sum.txt

echo "done"
