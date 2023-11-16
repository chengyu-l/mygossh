// Copyright 2018 github.com/andesli/mygossh Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// Author: andes
// Email: email.tata@qq.com

package main

import (
	"flag"
	"fmt"
	"github.com/andesli/mygossh/enc"
	"os"
)

var (
	encry = flag.Bool("e", false, "encryption(加密)")
	decry = flag.Bool("d", false, "decrypt(解密)")
	//pkey  = flag.String("k", "", "aes加密密钥")
)

func usage() {
	flag.PrintDefaults()
	os.Exit(1)
}

func main() {

	flag.Parse()

	//if *pkey != "" {
	//	enc.SetKey([]byte(*pkey))
	//}

	if flag.NArg() < 1 {
		usage()
	}
	opstr := flag.Arg(0)

	rzkey := enc.GetKey()

	if *encry {
		text, err := fenc(opstr, rzkey)
		if err != nil {
			fmt.Printf("Error:%s\n", err)
			return
		}
		fmt.Println(text)
	} else if *decry {
		text, err := fdec(opstr, rzkey)
		if err != nil {
			fmt.Printf("Error:%s\n", err)
			return
		}
		fmt.Println(string(text))

	} else {
		usage()
	}
}

func fdec(pass string, key []byte) ([]byte, error) {
	skey := key[:16]
	return enc.AesDecEncode(pass, skey)
}
func fenc(pass string, key []byte) (string, error) {
	skey := key[:16]
	return enc.AesEncEncode([]byte(pass), skey)
}
