package main

import (
	"fmt"
	"runtime"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	//"github.com/readline"
	//"os"
	"strings"
	"bufio"
	"bytes"
	//"io"
)

var workers = runtime.NumCPU()

func main() {

	runtime.GOMAXPROCS(workers)
	fmt.Printf("hello\n")

	//readline.ReadLine(os.Stdin, func(line string){fmt.Printf("Line[%s]\n", line)})
	str := "aaa\nbbb\r\nccc\ndddd"
	r := strings.NewReader(str)
	buf := bufio.NewReader(r)
	line, err := buf.ReadBytes('\n')
	for err == nil {
		line = bytes.TrimRight(line,"\n")
		if len(line) > 0 {
			if line[len(line)-1] == 13 {
				line = bytes.TrimRight(line, "\r")
			}
			fmt.Printf("line[%s]\n", line)
		}
		line, err = buf.ReadBytes('\n')
	}
	fmt.Printf("out line[%s]\n", line)

	db, err := sql.Open("mysql", "gerryyang:@(127.0.0.1:3306)/db_conf");
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	stmtOut, err := db.Prepare("SELECT Fservice_type from xxx_config where Fservice_code = ?")
	if err != nil {
		panic(err.Error())
	}
	defer stmtOut.Close()

	var rec string
	err = stmtOut.QueryRow("-gerry_once_cgi").Scan(&rec)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("rec[%s]\n", rec)


}

