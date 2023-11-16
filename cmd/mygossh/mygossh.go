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
	"github.com/andesli/mygossh/help"
	"github.com/andesli/mygossh/logs"
	"github.com/andesli/mygossh/machine"
	"github.com/andesli/mygossh/run"
	"github.com/andesli/mygossh/tools"
	"path/filepath"
	"strings"
	"sync"
)

// github.com/andesli/mygossh version
const (
	AppVersion = "mygossh 0.7"
)

var (

	//common options
	port     = flag.String("P", "22", "ssh端口")
	host     = flag.String("h", "", "ssh ip")
	user     = flag.String("u", "root", "ssh用户名")
	psw      = flag.String("p", "", "ssh密码，支持本地密码文件，默认文件路径：~/.mygossh/passwd")
	prunType = flag.String("t", "cmd", "运行模式(cmd|push|pull)：cmd 远程执行命令，默认值；push 推送文件到远程； pull拉取远程文件到本地")

	//batch running options
	ipFile = flag.String("i", "ip.txt", "批量执行是指定ip文件，有关文件格式见下文")
	cons   = flag.Int("c", 30, "批量并发执行的并发度，默认值是30，如果指定为1，mygossh是串行执行")

	//safe options
	encFlag = flag.Bool("e", false, "如果密码传递的是密文，使用-e标记。-e适用于通过-p传递的密码和-i 指定的文件中存放的密码字段")
	force   = flag.Bool("f", false, "如果遇到危险命令mygossh默认是不执行，可以通过-f强制执行，-f 是bool型参数，不指定默认为false")
	psafe   = flag.Bool("s", false, "-s是个bool型，只有到-c被指定为1时才有效，用来控制串行执行报错后是否立即退出")
	//pkey      = flag.String("key", "", "aes key for password decrypt and encryption")
	blackList = []string{"rm", "mkfs", "mkfs.ext3", "make.ext2", "make.ext4", "make2fs", "shutdown", "reboot", "init", "dd"}

	//log options
	plogLevel = flag.String("l", "info", "log level (debug|info|warn|error")
	plogPath  = flag.String("logpath", "./mygossh_log/", "logfile path")
	log       = logs.NewLogger()
	logFile   = "mygossh.log"

	pversion = flag.Bool("version", false, "mygossh version")

	//Timeout
	ptimeout = flag.Int("timeout", 10, "ssh timeout setting")
)

// main
func main() {

	usage := func() {
		fmt.Println(help.Help)
	}

	flag.Parse()

	//version
	if *pversion {
		fmt.Println(AppVersion)
		return
	}

	//if *pkey != "" {
	//	enc.SetKey([]byte(*pkey))
	//}

	if flag.NArg() < 1 {
		usage()
		return
	}

	if *prunType == "" || flag.Arg(0) == "" {
		usage()
		return
	}

	if err := initLog(); err != nil {
		fmt.Printf("init log error:%s\n", err)
		return
	}

	//异步日志，需要最后刷新和关闭
	defer func() {
		log.Flush()
		log.Close()
	}()

	log.Debug("parse flag ok , init log setting ok.")

	switch *prunType {
	//run command on remote server
	case "cmd":
		if flag.NArg() != 1 {
			usage()
			return
		}

		cmd := flag.Arg(0)

		if flag := tools.CheckSafe(cmd, blackList); !flag && *force == false {
			fmt.Printf("Dangerous command in %s", cmd)
			fmt.Printf("You can use the `-f` option to force to excute")
			log.Error("Dangerous command in %s", cmd)
			return
		}

		puser := run.NewUser(*user, *port, *psw, *force, *encFlag)
		log.Info("mygossh -t=cmd  cmd=[%s]", cmd)

		if *host != "" {
			log.Info("[servers]=%s", *host)
			run.SingleRun(*host, cmd, puser, *force, *ptimeout)

		} else {
			cr := make(chan machine.Result)
			ccons := make(chan struct{}, *cons)
			wg := &sync.WaitGroup{}
			run.ServersRun(cmd, puser, wg, cr, *ipFile, ccons, *psafe, *ptimeout)
			wg.Wait()
		}

	//push file or dir  to remote server
	case "scp", "push":

		if flag.NArg() != 2 {
			usage()
			return
		}

		src := flag.Arg(0)
		dst := flag.Arg(1)
		log.Info("mygossh -t=push local-file=%s, remote-path=%s", src, dst)

		puser := run.NewUser(*user, *port, *psw, *force, *encFlag)
		if *host != "" {
			log.Info("[servers]=%s", *host)
			run.SinglePush(*host, src, dst, puser, *force, *ptimeout)
		} else {
			cr := make(chan machine.Result, 20)
			ccons := make(chan struct{}, *cons)
			wg := &sync.WaitGroup{}
			run.ServersPush(src, dst, puser, *ipFile, wg, ccons, cr, *ptimeout)
			wg.Wait()
		}

	//pull file from remote server
	case "pull":
		if flag.NArg() != 2 {
			usage()
			return
		}

		//本地目录
		src := flag.Arg(1)
		//远程文件
		dst := flag.Arg(0)
		log.Info("mygossh -t=pull remote-file=%s  local-path=%s", dst, src)

		puser := run.NewUser(*user, *port, *psw, *force, *encFlag)
		if *host != "" {
			log.Info("[servers]=%s", *host)
			run.SinglePull(*host, puser, src, dst, *force)
		} else {
			run.ServersPull(src, dst, puser, *ipFile, *force)
		}

	default:
		usage()
	}
}

// setting log
func initLog() error {
	switch *plogLevel {
	case "debug":
		log.SetLevel(logs.LevelDebug)
	case "error":
		log.SetLevel(logs.LevelError)
	case "info":
		log.SetLevel(logs.LevelInfo)
	case "warn":
		log.SetLevel(logs.LevelWarn)
	default:
		log.SetLevel(logs.LevelInfo)
	}

	logpath := *plogPath
	err := tools.MakePath(logpath)
	if err != nil {
		return err
	}

	logname := filepath.Join(logpath, logFile)
	logstring := `{"filename":"` + logname + `"}`

	//此处主要是处理windows下文件路径问题,不做转义，日志模块会报如下错误
	//logs.BeeLogger.SetLogger: invalid character 'g' in string escape code
	logstring = strings.Replace(logstring, `\`, `\\`, -1)

	err = log.SetLogger("file", logstring)
	if err != nil {
		return err
	}
	//开启日志异步提升性能
	log.Async()
	return nil
}
