package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const infFie = "backinfo.ginfo"

type copytime struct {
	month time.Month
	date  int
}

type dir struct {
	disk, home, company string
}

type wkdir struct {
	category string
	path     dir
	dtsTime  copytime
	stdTime  copytime
}

func (wd *wkdir) write() string {
	return fmt.Sprintln(wd.category, wd.path.company, wd.path.disk,
		wd.path.home, wd.dtsTime.month, wd.dtsTime.date,
		wd.stdTime.month, wd.stdTime.date)
}

func main() {

	//Holds the total backup floders,Max 5
	var wd = make([]wkdir, 8)
	var swkdir []string

	file, err := os.Open(infFie)
	defer file.Close()
	if err != nil {
		log.Fatal("Open file failure!Quit...")
	}
	r := bufio.NewReader(file)
	var index int
	for {
		str, err := r.ReadString('\n')
		if err != nil {
			log.Println(err)
			break
		}
		swkdir = strings.Fields(str)
		if len(swkdir) < 8 {
			continue
		}
		wd[index].category = swkdir[0]
		wd[index].path.company = swkdir[1]
		wd[index].path.disk = func() string {
			k := []byte(swkdir[2])
			ret, _ := os.Getwd()
			k[0] = []byte(ret)[0]
			return string(k)
		}()
		wd[index].path.home = swkdir[3]
		wd[index].dtsTime.month = mon(swkdir[4])
		wd[index].dtsTime.date, err = strconv.Atoi(swkdir[5])
		wd[index].dtsTime.month = mon(swkdir[6])
		wd[index].dtsTime.date, err = strconv.Atoi(swkdir[7])
		index++
	}

		notes := `
			*******************************
			文件备份处理      
			1)公    司 ---->移动硬盘
			2)移动硬盘 ---->公    司
			3)宿    舍 ---->移动硬盘
			4)移动硬盘 ---->宿    舍
			*******************************
			`
		fmt.Println(notes)

		buf := bufio.NewReader(os.Stdin)
		choice, _ := buf.ReadByte()
		switch string(choice) {
		case "1":
			func() {
				for _, v := range wd {
					backUp(v.path.company, v.path.disk)
				}
			}()
		case "2":
			func() {
				for _, v := range wd {
					backUp(v.path.disk, v.path.company)
				}
			}()
		case "3":
			for _, v := range wd {
				backUp(v.path.home, v.path.disk)
			}
		case "4":
			for _, v := range wd {
				backUp(v.path.disk, v.path.home)
			}
		default:
			fmt.Println("Default is running")
			
			//Write the ginfo
			f, err := os.Create(infFie)
			if err != nil {
				log.Fatal("Open file for write failure!Quit...")
			}
			defer f.Close()
			//Todo
				f.WriteString("Last copy time")
			f.WriteString("Category Company Disk Home Month-Date Month-Date")
			for _, v := range wd {
				f.WriteString(v.write())
		
		}

	}

}

//convert string to time.Month
func mon(month string) time.Month {
	switch month {
	case "January":
		return time.January
	case "February":
		return time.February
	case "March":
		return time.March
	case "April":
		return time.April
	case "May":
		return time.May
	case "June":
		return time.June
	case "July":
		return time.July
	case "August":
		return time.August
	case "September":
		return time.September
	case "October":
		return time.October
	case "November":
		return time.November
	case "December":
		return time.December
	default:
		return time.January
	}
}

func backUp(src, des string) {
	cmd := exec.Command("cmd", "/c", "robocopy", "/mir", src, des)
	outf, _ := cmd.CombinedOutput()
	fmt.Println(string(outf))
}
