package global

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"runtime"
	"syscall"
	"time"

	"github.com/ethtweet/ethtweet/logs"
)

func CheckWindowsMysqld() error {
	if err := PingMysql(); err == nil {
		return nil
	}

	if err := RunWindowsMysqld(); err != nil {
		return fmt.Errorf("start windows mysqld failed:%w", err)
	}

	return PingMysql()
}

func RunWindowsMysqld() error {
	if runtime.GOOS != "windows" {
		msg := "not windows os, should not have been here"
		logs.PrintlnWarning(msg)
		return fmt.Errorf("msg")
	}

	_, err := os.Stat("mysql/bin/mysqld.exe")
	if err != nil {
		msg := "mysqld not found"
		logs.PrintlnInfo(msg)
		return fmt.Errorf(msg)
	}

	logs.PrintlnSuccess("mysqld found")
	_, err = os.Stat("mysql/data/ibdata1")
	//初始化数据库
	if err != nil {
		if err = ExecCmd("cmd.exe", "/c", ".\\mysql\\bin\\mysqld.exe --default-authentication-plugin=mysql_native_password --initialize-insecure --user=root --console"); err == nil {
			return nil
		}
	}

	if err != nil {
		logs.PrintDebugErr("exec cmd.exe with mysql_native_password err:%v", err)
	}

	// here we know mysql isn't running, so we try to start it again.
	// todo:与上面的cmd.exe有啥区别?
	if err = ExecCmd("cmd.exe", "/c", ".\\mysql\\bin\\mysqld.exe --console"); err != nil {
		return fmt.Errorf("exec cmd.exe --console err:%w", err)
	}

	// 休眠2秒给mysql启动
	time.Sleep(2 * time.Second)
	return nil
}

func PingMysql() error {
	// add mysql ping, avoid other proc occupy the port 3306
	conn, err := net.DialTimeout("tcp", "127.0.0.1:3306", 1*time.Second)
	if err == nil && conn != nil {
		conn.Close()
		return nil
	}
	return fmt.Errorf("ping mysql failed:%w", err)
}

func ExecCmd(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	logs.PrintlnSuccess("init mysqld")
	err := cmd.Start()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			status := exitErr.Sys().(syscall.WaitStatus)
			switch {
			case status.Exited():
				logs.PrintDebugErr("Return exit error: exit code=%d\n", status.ExitStatus())
			case status.Signaled():
				logs.PrintDebugErr("Return exit error: signal code=%d\n", status.Signal())
			}
		} else {
			logs.PrintDebugErr("Return other error: %s\n", err)
		}
		return err
	}

	logs.PrintlnSuccess("mysql init success")
	return nil
}
