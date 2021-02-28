package taillog

import (
	"fmt"
	"github.com/hpcloud/tail"
	"time"
)

func main() {
	filename := `./test.log`
	config := tail.Config{
		ReOpen: true,  // 当文件大于设定的值后重新打开
		Follow: true,  // 是否跟随
		Location:&tail.SeekInfo{Offset: 0, Whence: 2},  // 从文件哪个地方读
		MustExist: false,  // 允许文件不存在
		Poll: true,  // 轮询
	}

	// 打开文件，开始读取数据
	tails, err := tail.TailFile(filename,config)
	if err != nil{
		fmt.Printf("tail %s failed, err:%v\n", filename, err)
		return
	}

	var (
		msg *tail.Line
		ok bool
	)
	for{
		msg,ok = <- tails.Lines  // chan tail.Line
		if !ok{
			fmt.Printf("tail file failed, filename:%s\n",tails.Filename)
			time.Sleep(time.Second)
			continue
		}
		fmt.Println("line:", msg.Text)
	}
}
