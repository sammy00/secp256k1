// Copyright 2015 Conformal Systems LLC. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file is ignored during the regular build due to the following build tag.
// It is called by go generate and used to automatically generate pre-computed
// tables used to accelerate operations.
// +build ignore

package main

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"fmt"
	"log"
	"os"

	btcec "github.com/sammyne/secp256k1"
)

func main() {
	fi, err := os.Create("secp256k1.go")
	if err != nil {
		log.Fatal(err)
	}
	defer fi.Close()

	// Compress the serialized byte points.
	serialized := btcec.S256().SerializedBytePoints()
	var compressed bytes.Buffer
	w := zlib.NewWriter(&compressed)
	if _, err := w.Write(serialized); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	w.Close()

	// Encode the compressed byte points with base64.
	encoded := make([]byte, base64.StdEncoding.EncodedLen(compressed.Len()))
	base64.StdEncoding.Encode(encoded, compressed.Bytes())

	fmt.Fprintln(fi, "// Copyright (c) 2015 Conformal Systems LLC.")
	fmt.Fprintln(fi, "// Use of this source code is governed by an ISC")
	fmt.Fprintln(fi, "// license that can be found in the LICENSE file.")
	fmt.Fprintln(fi)
	fmt.Fprintln(fi, "package secp256k1")
	fmt.Fprintln(fi)
	fmt.Fprintln(fi, "// Auto-generated file (see genprecomps.go)")
	fmt.Fprintln(fi, "// DO NOT EDIT")
	fmt.Fprintln(fi)
	fmt.Fprintf(fi, "var secp256k1BytePoints = []byte(%q)\n", encoded)
}
