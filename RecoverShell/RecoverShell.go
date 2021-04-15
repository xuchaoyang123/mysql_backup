package RecoverShell

//
//import (
//	"time"
//	"check_list/config"
//	"os/exec"
//	"context"
//	"fmt"
//	"strings"
//	"os"
//)
//
//package execShell
//
//import (
//"check_list/config"
//"context"
//"fmt"
//"os"
//"os/exec"
//"strings"
//"time"
//)
//
//func ExecShell() (err error) {
//
//	//定义备份的文件显示的日期格式
//	t := time.Now()
//	FormatTimes := time.Now().Format("2006-01-02")
//
//	//读取json内mysql项里的内容
//	mysql := config.ConfMap["mysql"]
//
//
//
//
//	//生成[全备]拼接命令
//	FullCmd := exec.CommandContext(context.TODO(), "/bin/bash", "-c", fmt.Sprintf(`innobackupex --defaults-file=%s --host=%s --port=%s --user=%s --password="%s"  --no-lock --no-timestamp  %s_%s_innobackup_full`, mysql["mycnf"], mysql["host"], mysql["port"], mysql["user"], mysql["password"], strings.Join([]string{mysql["save_path"], FormatTimes}, "/"), mysql["port"]))
//	fmt.Println(FullCmd)
//	//fmt.Println(sprintf)
//
//	//生成[增量备]拼接命令
//	IncrementCmd := exec.CommandContext(context.TODO(), "/bin/bash", "-c", fmt.Sprintf(`innobackupex --defaults-file=%s --host=%s --port=%s --user=%s --password="%s"  --no-lock --no-timestamp  %s_%s_innobackup_increment`, mysql["mycnf"], mysql["host"], mysql["port"], mysql["user"], mysql["password"], strings.Join([]string{mysql["save_path"], FormatTimes}, "/"), mysql["port"]))
//	//fmt.Println(IncrementCmd)
//
//	//拼接[全备]目录和文件名
//	Innobackup_Full := strings.Join([]string{strings.Join([]string{mysql["save_path"], FormatTimes}, "/"), mysql["port"], "innobackup_full"}, "_")
//	fmt.Println(Innobackup_Full)
//	//拼接[增量备]目录和文件名
//	Innobackup_Increment := strings.Join([]string{strings.Join([]string{mysql["save_path"], FormatTimes}, "/"), mysql["port"], "innobackup_increment"}, "_")
//	fmt.Println(Innobackup_Increment)
//
//	//全备
//	//执行命令前判断前目录是否存在
//	full, err := os.Stat(Innobackup_Full)
//
//	//判断是否有错误
//	if err == nil {
//		fmt.Printf("%s 全备目录已存在,文件大小是%v:KB\n", full.Name(), full.Size()/1024)
//	} else if int(t.Weekday()) == 0 {
//
//		/*如果是星期日,则进行全备,数字0=星期日,1-6=其他星期时间*/
//		//[执行全备命令]
//		if _, err := FullCmd.CombinedOutput(); err != nil {
//			fmt.Println("[error] :调用innobackupex全备命令时出错,err= ", err)
//			return err
//		} else {
//			//fmt.Println(v)
//		}
//
//	}
//
//	//增量
//	//执行命令前判断前目录是否存在
//	inr, err := os.Stat(Innobackup_Increment)
//
//	//判断是否有错误
//	if err == nil {
//		fmt.Printf("%s 增量备目录已存在,文件大小是%v:KB\n", inr.Name(), inr.Size()/1024)
//	} else if int(t.Weekday()) != 0 {
//		//[执行增量备命令]
//		if _, err := IncrementCmd.CombinedOutput(); err != nil {
//			fmt.Println("[error] :调用innobackupex增量备命令时出错,err= ", err)
//			return err
//		} else {
//			//fmt.Println(v)
//		}
//
//	}
//
//	//
//	//
//	////执行命令前判断前目录是否存在
//	//info, err := os.Stat(Innobackup_Full)
//	//
//	////判断是否有错误
//	//if err == nil {
//	//	fmt.Printf("%s 目录已存在,文件大小是%v:KB\n", info.Name(), info.Size()/1024)
//	//} else if int(t.Weekday()) == 0 {
//	//
//	//	/*如果是星期日,则进行全备,数字0=星期日,1-6=其他星期时间*/
//	//	//[执行全备命令]
//	//	if _, err := FullCmd.CombinedOutput(); err != nil {
//	//		fmt.Println("[error] :调用innobackupex全备命令时出错,err= ", err)
//	//		return err
//	//	} else {
//	//		//fmt.Println(v)
//	//	}
//	//
//	//	/*如果不是星期日,则进行增量备,数字0=星期日,1-6=其他星期时间*/
//	//} else if int(t.Weekday()) != 0 {
//	//	//[执行增量备命令]
//	//	if _, err := IncrementCmd.CombinedOutput(); err != nil {
//	//		fmt.Println("[error] :调用innobackupex增量备命令时出错,err= ", err)
//	//		return err
//	//	} else {
//	//		//fmt.Println(v)
//	//	}
//	//
//	//}
//
//	//处理恐慌中断.
//	defer func() {
//		if errs := recover(); errs != nil {
//			err = errs.(error)
//		}
//	}()
//
//	return
//
//}
