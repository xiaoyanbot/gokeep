package main
 
import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
)
 
func main() {
	argc := len(os.Args)
	if argc < 2 {
		fmt.Println("Usage:", os.Args[0], " pragram args...")
		return
	}
	workdir, err := os.Getwd()
	if err != nil {
		fmt.Println("运行错误", err.Error())
	}
 
	name := os.Args[1]
	args := os.Args[1:]
	{
		if filepath.Base(name) == name {
			if lp, err := exec.LookPath(name); err != nil {
				log.Println("找不到待执行程序", err.Error())
				return
			} else {
				name = lp
			}
		}
	}
 
	log.Println("程序工作路径:", workdir)
	var cmdline string = strings.Join(args," ")
	log.Println("开始运行:", cmdline)
	var cmd *exec.Cmd
	//创建监听退出chan
	c := make(chan os.Signal)
	//监听指定信号 ctrl+c kill
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		for s := range c {
			switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				log.Println("信号退出", s)
				err = cmd.Process.Kill()
				if err != nil {
					log.Println("清理程序资源:", err.Error())
				}
				os.Exit(0)
			}
		}
	}()
 
	for {
		cmd = &exec.Cmd{
			Path:   name,
			Args:   args,
			Dir:    workdir,
			Stdin:  os.Stdin,
			Stdout: os.Stdout,
			Stderr: os.Stderr,
		}
		// log.Println(cmd.Args)
		err = cmd.Run()
		if err != nil {
			log.Println("程序运行错误:", err.Error())
			err = cmd.Process.Kill()
			if err != nil {
				log.Println("清理程序资源:", err.Error())
			}
			log.Println("开始重启程序")
			continue
		}
		exitcode := cmd.ProcessState.ExitCode()
		if exitcode != 0 {
			log.Println("程序错误退出:", exitcode)
			log.Println("开始重启程序")
			continue
		}
		break
	}
	log.Println("正常退出程序")
}
