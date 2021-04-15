package execShell

import (
	"check_list/config"
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var (
	List map[int64]string = map[int64]string{} //增量备份获取最新依赖的全量信息,用到的map
	Keys                  = make([]int64, 0)   //增量备份获取最新依赖的全量信息,用到的列表

	List002 map[int]string = map[int]string{} //删除全备份天数,用到的map
	Keys002                = make([]int, 0)   //删除全备份天数,用到的列表

	List003 map[int]string = map[int]string{} //删除增备份天数,用到的map
	Keys003                = make([]int, 0)   //删除增备份天数,用到的列表

	List004 map[int]string = map[int]string{} //删除日志备份天数,用到的map
	Keys004                = make([]int, 0)   //删除日志备份天数,用到的列表

	List02 map[int64]string = map[int64]string{} //删除全备份个数,用到的map
	Keys02                  = make([]int64, 0)   //删除全备份个数,用到的列表

	List03 map[int64]string = map[int64]string{} //删除增备份个数,用到的map
	Keys03                  = make([]int64, 0)   //删除增备份个数,用到的列表

	List04 map[int64]string = map[int64]string{} //删除日志备份个数,用到的map
	Keys04                  = make([]int64, 0)   //删除日志备份个数,用到的列表

	mysql       = config.InitConfig()             //获取配置文件信息
	FormatTimes = time.Now().Format("2006-01-02") //定义备份的文件显示的日期格式
	t           = time.Now()                      //定义赋值当前时间
	strRet      = mysql["save_path"]              //获取目录

	Full_Slave_Backup_Date = mysql["full_slave_backup_date"]
	Inc_Slave_Backup_Date  = mysql["inc_slave_backup_date"]
	Log_Slave_Backup_Date  = mysql["log_slave_backup_date"]

	Full_Slave_Backup_Count = mysql["full_slave_backup_count"]
	Inc_Slave_Backup_Count  = mysql["inc_slave_backup_count"]
	Log_Slave_Backup_Count  = mysql["log_slave_backup_count"]
)

//总调度
func ExecShell() (err error) {

	//判断当天是星期几,是不是满足备份条件
	if mysql["full_bk_date"] == "" {
		fmt.Printf("[warning]: 😔 [full_bk_date] 参数未设置,请设置一个有效参数...\n")
	} else {
		full_bk_date := strings.Split(mysql["full_bk_date"], ",")
		for _, v := range full_bk_date {
			//判断当天是不是满足备份条件
			if strings.Contains(v, strconv.Itoa(int(t.Weekday()))) == true {

				FUllBak()

			}
		}
	}

	//根据定义的时间进行增量备份
	if mysql["inc_bk_date"] == "" {
		fmt.Printf("[warning]: 😔 [inc_bk_date] 参数未设置,请设置一个有效参数...\n")
	} else {
		inc_bk_date := strings.Split(mysql["inc_bk_date"], ",")

		for _, v := range inc_bk_date {
			if strings.Contains(v, strconv.Itoa(int(t.Weekday()))) == true {
				InrBak()
			}
		}
	}
	//执行清理策略
	//如果任意一个参数里是空的 就以保留份数的方法来清理文件
	switch {
	case (Full_Slave_Backup_Date == "" || Inc_Slave_Backup_Date == "" || Log_Slave_Backup_Date == "") && (Full_Slave_Backup_Count == "" || Inc_Slave_Backup_Count == "" || Log_Slave_Backup_Count == ""):
		fmt.Printf("[warning]: 😔 请使用一种清理策略,另外一种全部进行注释...\n")
	case Full_Slave_Backup_Date == "" || Inc_Slave_Backup_Date == "" || Log_Slave_Backup_Date == "":
		CountCleanFile()
		//fmt.Println("CountCleanFile")
	case Full_Slave_Backup_Count == "" || Inc_Slave_Backup_Count == "" || Log_Slave_Backup_Count == "":
		DateCleanFile()
		//fmt.Println("DateCleanFile")

	default:

	}

	//log backup
	LogBak()

	return
}

//【全备】 整套操作
func FUllBak() (err error) {

	//【全备】 整套操作
	if mysql["mycnf"] == "" || mysql["host"] == "" || mysql["port"] == "" || mysql["user"] == "" || mysql["password"] == "" || mysql["save_path"] == "" {

		fmt.Printf("[warning]: 😔 数据库链接信息配置不完整,请检查配置信息..\n")

	} else {
		//生成[全备]拼接命令
		FullCmd := exec.CommandContext(context.TODO(), "/bin/bash", "-c", fmt.Sprintf(`innobackupex --defaults-file=%s --host=%s --port=%s --user=%s --password="%s"  --no-lock --no-timestamp  %s_%s_innobackup_full`, mysql["mycnf"], mysql["host"], mysql["port"], mysql["user"], mysql["password"], strings.Join([]string{mysql["save_path"], FormatTimes}, "/"), mysql["port"]))
		//fmt.Println(FullCmd)

		//拼接[全备]目录和文件名
		Innobackup_Full := strings.Join([]string{strings.Join([]string{mysql["save_path"], FormatTimes}, "/"), mysql["port"], "innobackup_full"}, "_")
		//fmt.Println(Innobackup_Full)

		//全备
		//执行命令前判断前目录是否存在
		full, err := os.Stat(Innobackup_Full)

		//判断是否有错误
		if err == nil {
			fmt.Printf("[warning]: 😔 全量备份文件:[%s]已存在,备份文件创建时间为:%s\n", full.Name(), full.ModTime().Format("2006-01-02 15:04:05"))
		} else {

			//[执行全备命令]
			BeginTimes := t
			fmt.Printf("[info]: 开始进行全量备份..\n")
			if _, err := FullCmd.CombinedOutput(); err != nil {
				fmt.Printf("[error]: 😭 调用innobackupex全备命令时出错,err=%s \n", err)
				return err
			} else {
				//fmt.Println(v)

			}
			Elapsed := time.Since(BeginTimes)
			fmt.Printf("[info]: 全量备份已完成..   耗时%s\n", Elapsed)

		}

	}

	return
}

//【增量备】 整套操作
func InrBak() (err error) {
	if mysql["mycnf"] == "" || mysql["host"] == "" || mysql["port"] == "" || mysql["user"] == "" || mysql["password"] == "" || mysql["save_path"] == "" {

		fmt.Println("[warning]: 😔 数据库链接信息配置不完整,请检查配置信息..")

	} else {

		//拼接出增量备份文件目录
		increment := mysql["save_path"] + "/" + FormatTimes + "_" + mysql["port"] + "_innobackup_increment"

		////【搜索过滤备份目录下,包含有全备标识的文件,判断出哪个是最新的】
		//01.获取文件或目录相关信息
		File, _ := ioutil.ReadDir(strRet)

		//02.过滤包含全备内容的目录信息追加到map中
		for _, Values := range File {
			if strings.Contains(Values.Name(), "innobackup_full") == true {
				List[Values.ModTime().Unix()] = Values.Name()
			}
		}

		//如果执行增量备份的时候发现没有全备,会进行提示,然后自动执行一次全备,在做增量备份
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("[warning]: 😔 执行增量备份时异常,因为全量备份不存在,现在先自动开始进行全量备份..\n")
				FUllBak()
				InrBak()
			}

		}()

		//03.将map的key转换为 列表中,然后准备进行 循环排序 求最新
		for k, _ := range List {
			Keys = append(Keys, k)
		}
		//04.排序算法来排序时间最大的是哪个文件
		for i := 0; i < len(Keys)-1; i++ {
			for j := i + 1; j < len(Keys); j++ {
				if Keys[j] > Keys[i] {
					Keys[i], Keys[j] = Keys[j], Keys[i]
				}
			}
		}
		//fmt.Println(List[Keys[0]])

		//通过当前星期-当天星期的天数=0 然后拼接出 上次全备份时的目录
		//full := mysql["save_path"] + "/" + t.AddDate(0, 0, -int(t.Weekday())).Format("2006-01-02") + "_" + mysql["port"] + "_innobackup_full"

		//拼接依赖的全备路径内容,通过上面求出的最新全备内容进行拼接
		full := mysql["save_path"] + "/" + List[Keys[0]]

		//fmt.Println(full)
		//生成[增量备]拼接命令
		IncrementCmd := exec.CommandContext(context.TODO(), "/bin/bash", "-c", fmt.Sprintf(`innobackupex --defaults-file=%s --host=%s --port=%s --user=%s --password="%s"  --no-lock --no-timestamp  --incremental-basedir=%s --incremental %s`, mysql["mycnf"], mysql["host"], mysql["port"], mysql["user"], mysql["password"], full, increment))

		//fmt.Println(IncrementCmd)
		//fmt.Printf("%T\n", IncrementCmd)

		//拼接[增量备]目录和文件名
		Innobackup_Increment := strings.Join([]string{strings.Join([]string{mysql["save_path"], FormatTimes}, "/"), mysql["port"], "innobackup_increment"}, "_")
		//fmt.Println(Innobackup_Increment)

		//如果执行增量备份的时候发现没有全备,则停止备份抛出错误提示
		//fmt.Println(full)

		_, err = os.Stat(full)
		if err != nil {
			fmt.Printf("[error]: 😭 执行增量备份失败,因为增量备份不存在,请检查备份数据..\n")
		} else {

			//增量
			//执行命令前判断前目录是否存在
			inr, err := os.Stat(Innobackup_Increment)

			//判断是否有错误
			if err == nil {
				fmt.Printf("[warning]: 😔 增量备份文件:[%s]已存在,备份文件创建时间为:%s\n", inr.Name(), inr.ModTime().Format("2006-01-02 15:04:05"))
			} else {

				fmt.Printf("[info]: 开始进行增量备份..\n")
				//[执行增量备命令]
				BeginTimes := t
				if _, err := IncrementCmd.CombinedOutput(); err != nil {
					fmt.Printf("[error]: 😭 调用innobackupex增量备命令时出错,err=%s\n", err)
					return err
				} else {
					//fmt.Println(v)
				}
				Elapsed := time.Since(BeginTimes)
				fmt.Printf("[info]: 增量备份已完成..   耗时%s\n", Elapsed)
			}

		}

	}

	return
}

//保留几天
func DateCleanFile() {

	Full_Slave_Backup_Date, _ := strconv.Atoi(Full_Slave_Backup_Date)
	Inc_Slave_Backup_Date, _ := strconv.Atoi(Inc_Slave_Backup_Date)
	Log_Slave_Backup_Date, _ := strconv.Atoi(Log_Slave_Backup_Date)

	//NowTime, _ := strconv.Atoi(t.Format("20060102"))

	File, _ := ioutil.ReadDir(strRet)
	for _, Values := range File {

		switch {

		//全量备份删除逻辑
		case strings.Contains(Values.Name(), "innobackup_full") == true:
			fuu, _ := strconv.Atoi(Values.ModTime().Format("20060102"))
			List002[fuu] = Values.Name()

			//	//增量备份删除逻辑
		case strings.Contains(Values.Name(), "innobackup_increment") == true:
			inr, _ := strconv.Atoi(Values.ModTime().Format("20060102"))
			List003[inr] = Values.Name()

		default:
		}
	}

	//-----------------全量-------------------------
	for k, _ := range List002 {
		Keys002 = append(Keys002, k)
	}
	//04.排序算法来排序时间最大的是哪个文件
	for i := 0; i < len(Keys002)-1; i++ {
		for j := i + 1; j < len(Keys002); j++ {
			if Keys002[j] > Keys002[i] {
				Keys002[i], Keys002[j] = Keys002[j], Keys002[i]
			}
		}
	}
	//算除了保留意外的文件个数是哪个
	if len(Keys002) > Full_Slave_Backup_Date {
		for _, v := range Keys002[Full_Slave_Backup_Date:] {
			//过滤删除最旧的那份,保留最新的N份
			//fmt.Println(List002[v])
			err := os.RemoveAll(strRet + "/" + List002[v])
			if err != nil {
				// 删除失败
				fmt.Printf("[error]: 😭 删除备份文件时出错,err=%s\n", err)
			} else {
				//删除成功
				fmt.Printf("[info]: 备份文件:[%s] 已删除完毕..\n", strRet+"/"+List002[v])
			}
		}
	}

	//-----------------增量-------------------------
	for k, _ := range List003 {
		Keys003 = append(Keys003, k)
	}
	//04.排序算法来排序时间最大的是哪个文件
	for i := 0; i < len(Keys003)-1; i++ {
		for j := i + 1; j < len(Keys003); j++ {
			if Keys003[j] > Keys003[i] {
				Keys003[i], Keys003[j] = Keys003[j], Keys003[i]
			}
		}
	}
	//算除了保留意外的文件个数是哪个
	if len(Keys003) > Inc_Slave_Backup_Date {
		for _, v := range Keys003[Inc_Slave_Backup_Date:] {
			//过滤删除最旧的那份,保留最新的N份
			//fmt.Println(List003[v])
			err := os.RemoveAll(strRet + "/" + List003[v])
			if err != nil {
				// 删除失败
				fmt.Printf("[error]: 😭 删除备份文件时出错,err=%s\n", err)
			} else {
				//删除成功
				fmt.Printf("[info]: 备份文件:[%s] 已删除完毕..\n", strRet+"/"+List003[v])
			}
		}
	}

	//-----------------日志-------------------------
	LogRet := mysql["save_path"] + "/innobackup_logfile/"
	File01, _ := ioutil.ReadDir(LogRet)
	for _, logs01 := range File01 {
		LogTime, _ := strconv.Atoi(logs01.ModTime().Format("20060102"))
		//fmt.Println("aaaaaaaaaa", LogTime)

		//求出最大的一天是那天
		for _, logs01 := range File01 {
			logs001, _ := strconv.Atoi(logs01.ModTime().Format("20060102"))
			List004[logs001] = logs01.Name()

		}
		for k, _ := range List004 {
			Keys004 = append(Keys004, k)
		}
		//04.排序算法来排序时间最大的是哪个文件
		for i := 0; i < len(Keys004)-1; i++ {
			for j := i + 1; j < len(Keys004); j++ {
				if Keys004[j] > Keys004[i] {
					Keys004[i], Keys004[j] = Keys004[j], Keys004[i]
				}
			}
		}

		if Keys004[0]-Log_Slave_Backup_Date >= LogTime {
			//fmt.Println(logs01.Name())
			err := os.RemoveAll(LogRet + logs01.Name())
			if err != nil {
				// 删除失败
				fmt.Printf("[error]: 😭 删除Binlog时出错,err=%s\n", err)
			} else {
				//删除成功
				fmt.Printf("[info]: Binlog备份:[%s] 已删除完毕..\n", LogRet+logs01.Name())
			}
		}
	}

	//NowTime, _ := strconv.Atoi(t.Format("20060102"))
	//LogRet := mysql["save_path"] + "/innobackup_logfile/"
	//File01, _ := ioutil.ReadDir(LogRet)
	//for _, logs01 := range File01 {
	//	LogTime, _ := strconv.Atoi(logs01.ModTime().Format("20060102"))
	//
	//	//求出最大的一天是那天
	//
	//	if NowTime-Log_Slave_Backup_Date >= LogTime {
	//		fmt.Println(logs01.Name())
	//		//err := os.RemoveAll(LogRet + logs01.Name())
	//		//if err != nil {
	//		//	// 删除失败
	//		//	fmt.Printf("[error]: 😭 删除Binlog时出错,err=%s\n", err)
	//		//} else {
	//		//	//删除成功
	//		//	fmt.Printf("[info]: Binlog备份:[%s] 已删除完毕..\n", LogRet+logs01.Name())
	//		//}
	//	}
	//}

}

////保留几天
//func DateCleanFile() {
//
//	Full_Slave_Backup_Date, _ := strconv.Atoi(Full_Slave_Backup_Date)
//	//Inc_Slave_Backup_Date, _ := strconv.Atoi(Inc_Slave_Backup_Date)
//	//Log_Slave_Backup_Date, _ := strconv.Atoi(Log_Slave_Backup_Date)
//
//	//清理备份功能
//	NowTime, _ := strconv.Atoi(t.Format("20060102"))
//
//	File, _ := ioutil.ReadDir(strRet)
//	for _, Values := range File {
//
//		switch {
//
//		//全量备份删除逻辑
//		//01.过滤文件类型
//		//case strings.Contains(Values.Name(), "innobackup_full") == true:
//		//	FileTime, _ := strconv.Atoi(Values.ModTime().Format("20060102"))
//		//	//02.判断文件是否过期
//		//	if NowTime-Full_Slave_Backup_Date > FileTime {
//		//		//03.删除过期文件
//		//		err := os.RemoveAll(strRet + "/" + Values.Name())
//		//		if err != nil {
//		//			// 删除失败
//		//			fmt.Printf("[error]: 😭 删除备份文件时出错,err=%s\n", err)
//		//		} else {
//		//			//删除成功
//		//			fmt.Printf("[info]: 备份文件:[%s] 已删除完毕..\n", strRet+"/"+Values.Name())
//		//		}
//		//	}
//
//		case strings.Contains(Values.Name(), "innobackup_full") == true:
//			FileTime, _ := strconv.Atoi(Values.ModTime().Format("20060102"))
//
//			//02.判断文件是否过期
//			if NowTime-Full_Slave_Backup_Date > FileTime {
//				//03.删除过期文件
//				//err := os.RemoveAll(strRet + "/" + Values.Name())
//				//if err != nil {
//				//	// 删除失败
//				//	fmt.Printf("[error]: 😭 删除备份文件时出错,err=%s\n", err)
//				//} else {
//				//	//删除成功
//				//	fmt.Printf("[info]: 备份文件:[%s] 已删除完毕..\n", strRet+"/"+Values.Name())
//				//}
//			}
//
//			//	//增量备份删除逻辑
//			//	//01.过滤文件类型
//			//case strings.Contains(Values.Name(), "innobackup_increment") == true:
//			//	FileTime, _ := strconv.Atoi(Values.ModTime().Format("20060102"))
//			//	//02.判断文件是否过期
//			//	if NowTime-Inc_Slave_Backup_Date > FileTime {
//			//		err := os.RemoveAll(strRet + "/" + Values.Name())
//			//		if err != nil {
//			//			// 删除失败
//			//			fmt.Printf("[error]: 😭 删除备份文件时出错,err=%s\n", err)
//			//		} else {
//			//			//删除成功
//			//			fmt.Printf("[info]: 备份文件:[%s] 已删除完毕..\n", strRet+"/"+Values.Name())
//			//		}
//			//	}
//			//default:
//		}
//	}
//
//	//	//日志备份删除逻辑
//	//	LogRet := mysql["save_path"] + "/innobackup_logfile/"
//	//	File01, _ := ioutil.ReadDir(LogRet)
//	//	for _, logs01 := range File01 {
//	//		LogTime, _ := strconv.Atoi(logs01.ModTime().Format("20060102"))
//	//		//02.判断文件是否过期
//	//		if NowTime-Log_Slave_Backup_Date >= LogTime {
//	//			err := os.RemoveAll(LogRet + logs01.Name())
//	//			if err != nil {
//	//				// 删除失败
//	//				fmt.Printf("[error]: 😭 删除Binlog时出错,err=%s\n", err)
//	//			} else {
//	//				//删除成功
//	//				fmt.Printf("[info]: Binlog备份:[%s] 已删除完毕..\n", LogRet+logs01.Name())
//	//			}
//	//		}
//	//	}
//}

//保留几份
func CountCleanFile() {

	File, _ := ioutil.ReadDir(strRet)
	for _, Values := range File {

		switch {

		//全量备份删除逻辑
		case strings.Contains(Values.Name(), "innobackup_full") == true:
			List02[Values.ModTime().Unix()] = Values.Name()

			//增量备份删除逻辑
		case strings.Contains(Values.Name(), "innobackup_increment") == true:
			List03[Values.ModTime().Unix()] = Values.Name()

		default:

		}
	}

	//全量删除几份判断要删除的文件
	for k, _ := range List02 {
		Keys02 = append(Keys02, k)
	}

	//04.排序算法来排序时间最大的是哪个文件
	for i := 0; i < len(Keys02)-1; i++ {
		for j := i + 1; j < len(Keys02); j++ {
			if Keys02[j] < Keys02[i] {
				Keys02[i], Keys02[j] = Keys02[j], Keys02[i]
			}
		}
	}

	//算除了保留意外的文件个数是哪个
	fsbc, _ := strconv.Atoi(Full_Slave_Backup_Count)
	i := len(Keys02) - fsbc

	if len(Keys02) > fsbc {

		//过滤删除最旧的那份,保留最新的N份
		for _, v := range Keys02[:i] {

			//fmt.Println(List02[v])
			err := os.RemoveAll(strRet + "/" + List02[v])
			if err != nil {
				// 删除失败
				fmt.Printf("[error]: 😭 删除备份文件时出错,err=%s\n", err)
			} else {
				//删除成功
				fmt.Printf("[info]: 备份文件:[%s] 已删除完毕..\n", strRet+"/"+List02[v])
			}

		}
	}

	//增量删除几份判断要删除的文件

	for k, _ := range List03 {
		Keys03 = append(Keys03, k)
	}

	//04.排序算法来排序时间最大的是哪个文件
	for i := 0; i < len(Keys03)-1; i++ {
		for j := i + 1; j < len(Keys03); j++ {
			if Keys03[j] < Keys03[i] {
				Keys03[i], Keys03[j] = Keys03[j], Keys03[i]
			}
		}
	}

	//算除了保留意外的文件个数是哪个
	isbc, _ := strconv.Atoi(Inc_Slave_Backup_Count)
	if len(Keys03) > isbc {

		i03 := len(Keys03) - isbc
		//过滤删除最旧的那份,保留最新的N份
		for _, v := range Keys03[:i03] {

			//fmt.Println(List02[v])
			err := os.RemoveAll(strRet + "/" + List03[v])
			if err != nil {
				// 删除失败
				fmt.Printf("[error]: 😭 删除备份文件时出错,err=%s\n", err)
			} else {
				//删除成功
				fmt.Printf("[info]: 备份文件:[%s] 已删除完毕..\n", strRet+"/"+List03[v])
			}

		}
	}

	//日志备份删除逻辑

	LogRet := mysql["save_path"] + "/innobackup_logfile/"
	File01, _ := ioutil.ReadDir(LogRet)

	for _, logs01 := range File01 {
		split := strings.Split(logs01.Name(), ".")
		logs001, _ := strconv.ParseInt(split[1], 10, 64)
		List04[logs001] = logs01.Name()

	}

	for k, _ := range List04 {
		Keys04 = append(Keys04, k)
	}

	//04.排序算法来排序时间最大的是哪个文件
	for i := 0; i < len(Keys04)-1; i++ {
		for j := i + 1; j < len(Keys04); j++ {
			if Keys04[j] < Keys04[i] {
				Keys04[i], Keys04[j] = Keys04[j], Keys04[i]
			}
		}
	}
	lsbc, _ := strconv.Atoi(Log_Slave_Backup_Count) //算除了保留意外的文件个数是哪个
	//fmt.Println(List04)

	if len(Keys04) > lsbc {
		i04 := len(Keys04) - lsbc
		//过滤删除最旧的那份,保留最新的N份
		for _, v := range Keys04[:i04] {
			//fmt.Println(List04[v])
			err := os.RemoveAll(LogRet + List04[v])
			if err != nil {
				// 删除失败
				fmt.Printf("[error]: 😭 删除Binlog时出错,err=%s\n", err)
			} else {
				//删除成功
				fmt.Printf("[info]: Binlog备份:[%s] 已删除完毕..\n", LogRet+List04[v])
			}

		}
	}

}

//通过访问数据库获取binlog相关信息
func ExecSql() (binlog, logfile string) {
	//"用户名:密码@[连接方式](主机名:端口号)/数据库名"

	conn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", mysql["user"], mysql["password"], mysql["host"], mysql["port"], "mysql")
	db, _ := sql.Open("mysql", conn) // 设置连接数据库的参数

	defer db.Close() //关闭数据库
	err := db.Ping() //连接数据库
	if err != nil {
		fmt.Println("[error]: 😭 数据库连接失败,请检查数据库链接信息...")
		return
	}

	variable_value, _ := db.Query("select variable_value from performance_schema.global_variables where variable_name='log_bin'") //获取所有数据

	for variable_value.Next() { //循环显示所有的数据
		variable_value.Scan(&binlog)
		//fmt.Println(values)
	}

	FILE_NAME, _ := db.Query("select SUBSTRING_INDEX(FILE_NAME,'/',-1)as logfile from performance_schema.file_instances where EVENT_NAME='wait/io/file/sql/binlog' order by 1 desc limit 1") //获取所有数据

	for FILE_NAME.Next() { //循环显示所有的数据
		FILE_NAME.Scan(&logfile)
		//fmt.Println(values)
	}

	return binlog, logfile
}

//日志备份
func LogBak() {

	/*
				1. 先要判断是否开启了binlog,如果没有开启抛出异常。
		        2. 获取binlog的名字 传到参数中引用,binlog的文件名字可以根据 配置中去查看
				3 . 第一次执行将命令后台执行,再次执行需要判断是否有这个进程,如果有不做任何操作,没有再次执行。
	*/

	////获取datadir目录位置
	//var Str2 []string
	//file, _ := os.OpenFile(mysql["mycnf"], os.O_RDONLY, 0666)
	//defer file.Close()
	////创建其带缓冲的读取器
	//reader := bufio.NewReader(file)
	//for {
	//	str, err := reader.ReadString('\n') //一行一行的读取
	//	if err == nil {
	//		//过滤#注释的内容,过滤bin-log开头的内容
	//		if strings.HasPrefix(str, "#") == false {
	//			if strings.HasPrefix(str, "datadir") {
	//				Str2 = strings.Split(str, "=")
	//
	//			}
	//
	//		}
	//	} else {
	//		//如果错误是已到文件末尾,就输出信息
	//		if err == io.EOF {
	//			break //读取完毕后 跳出循环,继续往下走.
	//		}
	//		return //异常错误,就退出
	//
	//	}
	//
	//}
	//fmt.Println(Str2[1]) // /data/mysql/mysql3306/data
	//
	//Mycnf, _ := ioutil.ReadFile(mysql["mycnf"])
	//Mycnf1 := string(Mycnf)
	//
	////1. 先要判断是否开启了binlog,如果没有开启抛出异常。
	//if (strings.Contains(Mycnf1, "log-bin") != true || strings.Contains(Mycnf1, "#log-bin")) && (strings.Contains(Mycnf1, "log_bin") != true || strings.Contains(Mycnf1, "#log_bin")) {
	//	fmt.Println("[error]: 😭 该数据库没有开启binlog,无法进行binlog备份...")
	//
	//} else {
	//	// 2. 获取binlog的名字 传到参数中引用,binlog的文件名字可以根据 配置中去查看;
	//	_, logfile := ExecSql()
	//
	//	LogCmd := exec.CommandContext(context.TODO(), "/bin/bash", "-c", fmt.Sprintf(`mysqlbinlog  --raw --read-from-remote-server --stop-never --host=%s  --port=%s --user=%s   --password=%s  %s &`, mysql["host"], mysql["port"], mysql["user"], mysql["password"], logfile))
	//	fmt.Println(LogCmd)

	//1. 先要判断是否开启了binlog,如果没有开启抛出异常。
	logpath := mysql["save_path"] + "/" + "innobackup_logfile/"
	binlog, logfile := ExecSql()

	//2.执行命令前 先判断是否有这个进程存在,如果存在不会重复执行
	var PID string
	cmd := exec.Command("/bin/bash", "-c", `ps aux | grep mysqlbinlog | grep -v "grep" |wc -l`)
	output, _ := cmd.Output()
	fields := strings.Fields(string(output))
	for _, v := range fields {
		PID = v
	}

	if binlog != "ON" {
		fmt.Println("[error]: 😭 无法获取binlog信息,binlog备份失败...")
	} else {

		_, err := os.Stat(logpath)
		if PID == "0" {
			if err != nil {
				os.Mkdir(logpath, os.ModePerm)
				LogCmd := exec.CommandContext(context.TODO(), "/bin/bash", "-c", fmt.Sprintf(`mysqlbinlog  --raw --read-from-remote-server --stop-never --host=%s  --port=%s --user=%s   --password=%s  --result-file=%s  %s &`, mysql["host"], mysql["port"], mysql["user"], mysql["password"], logpath, logfile))
				LogCmd.Start()
				fmt.Println("[info]: BinlogServer已开始运行...")

			} else {
				LogCmd := exec.CommandContext(context.TODO(), "/bin/bash", "-c", fmt.Sprintf(`mysqlbinlog  --raw --read-from-remote-server --stop-never --host=%s  --port=%s --user=%s   --password=%s  --result-file=%s  %s &`, mysql["host"], mysql["port"], mysql["user"], mysql["password"], logpath, logfile))
				LogCmd.Start()
				fmt.Println("[info]: BinlogServer已开始运行...")

				//if PID == "0" {
				//	fmt.Printf("[warning]: 😔 BinlogServer运行失败,执行命令内容为[%s]..\n", LogCmd)
				//}
			}
		}

	}

}
