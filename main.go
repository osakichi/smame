package main

// WOPT 01
// WUART 80
// で使用

// 電力メーター：
// 機器番号：A15G592268
// 認証ID：000000BB040F00000000000000177824
// パスワード：CNWG6JIPFGJO
//EVENT 20 FE80:0000:0000:0000:021D:1291:0002:2B4F 0
//EPANDESC
//  Channel:3B
//  Channel Page:09
//  Pan ID:0257
//  Addr:38E08E00001C8257
//  LQI:C4
//  Side:0
//  PairID:00177824
//EVENT 22 FE80:0000:0000:0000:021D:1291:0002:2B4F 0

import (
	"bufio"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/tarm/serial"
)

type Command struct {
	Request  string
	Response string
}

func ReadWait(r *bufio.Reader, s string) error {
	for {
		buf, _, err := r.ReadLine()
		if err != nil {
			return err
		}
		hexdump(buf)
		//fmt.Println(string(buf))
		if strings.HasPrefix(string(buf), s) {
			break
		}
	}
	return nil
}

func main() {
	dongle := "/dev/ttyUSB0"

	c := &serial.Config{Name: dongle, Baud: 115200}
	port, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}
	defer port.Close()

	rd := bufio.NewReader(port)

	com := []Command{
		{"SKVER", "OK"},
		{"SKINFO", "OK"},
		{"SKAPPVER", "OK"},

		{"SKRESET", "OK"},
		{"SKSETRBID 000000BB040F00000000000000177824", "OK"},
		{"SKSETPWD C CNWG6JIPFGJO", "OK"},
		{"SKSCAN 2 FFFFFFFF 6 0", "EVENT 22"},
		{"SKSREG S2 3B", "OK"},   // CH設定
		{"SKSREG S3 0257", "OK"}, // PAN ID設定
		{"SKJOIN FE80:0000:0000:0000:3AE0:8E00:001C:8257", "EVENT 25"},

		{"SKSENDTO 1 FE80:0000:0000:0000:3AE0:8E00:001C:8257 0E1A 1 0 0006 HGWreq", "hoge"},
	}

	start := time.Now()

	for _, c := range com {
		_, err := port.Write([]byte(fmt.Sprintf("%s\r\n", c.Request)))
		if err != nil {
			log.Fatal(err)
		}

		err = ReadWait(rd, c.Response)
		if err != nil {
			log.Fatal(err)
		}
	}

	end := time.Now()
	fmt.Println(end.Sub(start).Seconds())
}

func hexdump(buf []byte) {
	fmt.Println(time.Now().Format("2006-01-02 15:04:05.000"))
	for x := 0; x < len(buf); x += 16 {
		hex := ""
		asc := ""
		for n := 0; n < 16; n++ {
			if x+n >= len(buf) {
				break
			}
			hex += fmt.Sprintf("%02x ", buf[x+n])
			if buf[x+n] >= 0x20 && buf[x+n] <= 0x7e {
				asc += fmt.Sprintf("%c", buf[x+n])
			} else {
				asc += "."
			}
		}
		fmt.Printf("    %-48s %s\n", hex, asc)
	}
}
