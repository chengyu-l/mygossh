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

package local_file

import (
	"bufio"
	"fmt"
	"github.com/andesli/mygossh/auth"
	"github.com/andesli/mygossh/enc"
	"github.com/andesli/mygossh/logs"
	"io"
	"os"
	"path"
	"strings"
)

var (
	log = logs.NewLogger()
)

type LocalConfigDriver struct {
	userPwdMap map[string]string
}

func init() {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		log.Warn("rget user home dir failed")
		return
	}

	// ~/.mygossh/passwd
	localFile := path.Join(userHomeDir, ".mygossh", "passwd")

	fd, err := os.OpenFile(localFile, os.O_RDONLY, 0644)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}

		log.Warn("read passwd file(%s) failed", localFile)
		return
	}

	defer fd.Close()

	var userPwdMap = make(map[string]string, 0)
	reader := bufio.NewReader(fd)
	for {
		data, readErr := reader.ReadString('\n')
		if readErr == io.EOF && len(data) == 0 {
			break
		}

		if len(data) > 0 && data[len(data)-1] == '\n' {
			data = data[:len(data)-1]
		}

		index := strings.Index(data, " ")
		if index == 0 {
			continue
		}

		user := strings.TrimSpace(data[:index])
		encPass := strings.TrimSpace(data[index:])
		userPwdMap[user] = encPass
	}

	db := &LocalConfigDriver{
		userPwdMap: userPwdMap,
	}

	auth.Register("local_file", db)
}

func (dv LocalConfigDriver) GetPassword(host, user string) (string, error) {
	pass, ok := dv.userPwdMap[user]
	if ok {
		skey := enc.GetKey()
		text, decErr := enc.AesDecEncode(pass, skey[:16])
		if decErr != nil {
			log.Warn("decode password failed, user:%s", user)
			return "", decErr
		}

		return string(text), nil
	}

	return "", fmt.Errorf("no password")
}
