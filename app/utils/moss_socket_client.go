package utils

// @Author: 陈健航
// @Date: 2021/1/10 21:14
// @Description: stanford moss代码查重的客户端，根据脚本写的，原脚本：http://moss.stanford.edu/general/scripts/mossnet
// 可直接使用，不建议修改

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/url"
	"strings"
)

type Stage int

const (
	disconnected Stage = iota
	awaitingInitialization
	awaitingLanguage
	awaitingFiles
	awaitingQuery
	awaitingResults
	awaitingEnd
)

type MossClient struct {
	currentStage       Stage  // 当前状态
	addr               string // 地址
	userId             string // 用户id
	language           string // 语言
	setId              int
	optM               int64
	optD               int
	optX               int
	optN               int
	optC               string
	ResultURL          *url.URL // 结果url
	supportedLanguages []string
	conn               *net.TCPConn
}

// NewMossClient 构造函数
// @params language_enum 语言
// @params userId 用户id
// @return *MossClient
// @return error
// @date 2021-01-11 17:03:26
func NewMossClient(language string, userId string) (*MossClient, error) {
	// 支持的语言
	supportedLanguages := []string{"c", "cc", "java", "ml", "pascal", "ada", "lisp", "schema", "haskell", "fortran",
		"ascii", "vhdl", "perl", "matlab", "python", "mips", "prolog", "spice", "vb", "csharp", "modula2", "a8086",
		"javascript", "plsql"}
	isContain := false
	for _, supportedLanguage := range supportedLanguages {
		if language == supportedLanguage {
			isContain = true
			break
		}
	}
	if isContain {
		return &MossClient{
			currentStage:       disconnected,
			setId:              1,
			optM:               10,
			optD:               1,
			optX:               0,
			optN:               250,
			optC:               "",
			supportedLanguages: supportedLanguages,
			addr:               "moss.stanford.edu:7690",
			userId:             userId,
			language:           language,
		}, nil
	} else {
		return nil, errors.New("MOSS Server does not recognize this programming language_enum")
	}
}

// Close 关闭
// @receiver c
// @return error
// @date 2021-01-11 17:04:12
func (c *MossClient) Close() (err error) {
	c.currentStage = disconnected
	if err = c.sendCommand("end\n"); err != nil {
		return err
	}
	if err = c.conn.Close(); err != nil {
		return err
	}
	return nil
}

// connect 连接斯坦福的moss
// @receiver c
// @return error
// @date 2021-01-11 22:38:17
func (c *MossClient) connect() (err error) {
	if c.currentStage != disconnected {
		return errors.New("client is already connected")
	} else {
		tcpAddr, err := net.ResolveTCPAddr("tcp4", c.addr)
		if err != nil {
			return err
		}
		if c.conn, err = net.DialTCP("tcp", nil, tcpAddr); err != nil {
			return err
		}
		c.currentStage = awaitingInitialization
	}
	return nil
}

// Run 启动
// @receiver c
// @return error
// @date 2021-01-11 22:38:28
func (c *MossClient) Run() (err error) {
	if err = c.connect(); err != nil {
		return err
	}
	if err = c.sendInitialization(); err != nil {
		return err
	}
	if err = c.sendLanguage(); err != nil {
		return err
	}
	return nil
}

// sendCommand 发送命令行
// @receiver c
// @params objects
// @return error
// @date 2021-01-11 22:38:40
func (c *MossClient) sendCommand(objects ...interface{}) (err error) {
	commandStrings := make([]string, 0, len(objects))

	for var5 := 0; var5 < len(objects); var5++ {
		o := objects[var5]
		s := fmt.Sprintf("%v", o)
		commandStrings = append(commandStrings, s)
	}
	if err := c.sendCommandStrings(commandStrings); err != nil {
		return err
	}
	return nil
}

// sendCommandStrings 发送命令指令序列
// @receiver c
// @params stringSlice
// @return error
// @date 2021-01-11 22:38:52
func (c *MossClient) sendCommandStrings(stringSlice []string) (err error) {
	if len(stringSlice) > 0 {
		//slice转字符串,空格分隔
		s := strings.Join(stringSlice, " ")
		s += "\n"
		if _, err := c.conn.Write([]byte(s)); err != nil {
			return errors.New("failed to send command: " + err.Error())
		}
		return nil
	} else {
		return errors.New("failed to send command because it was empty")
	}
}

// sendInitialization 输出化序列值
// @receiver c
// @return error
// @date 2021-01-11 22:39:11
func (c *MossClient) sendInitialization() (err error) {
	if c.currentStage != awaitingInitialization {
		return errors.New("cannot send initialization. Client is either already initialized or not connected yet")
	}
	if err = c.sendCommand("moss", c.userId); err != nil {
		return nil
	}
	if err = c.sendCommand("directory", c.optD); err != nil {
		return nil
	}
	if err = c.sendCommand("X", c.optX); err != nil {
		return nil
	}
	if err = c.sendCommand("maxmatches", c.optM); err != nil {
		return nil
	}
	if err = c.sendCommand("show", c.optN); err != nil {
		return nil
	}
	c.currentStage = awaitingLanguage
	return nil
}

// sendLanguage 发送语言
// @receiver c
// @return error
// @date 2021-01-11 22:39:27
func (c *MossClient) sendLanguage() error {
	if c.currentStage != awaitingLanguage {
		return errors.New("language_enum already sent or client is not initialized yet")
	}
	if err := c.sendCommand("language", c.language); err != nil {
		return err
	}
	b := make([]byte, 1024)
	n, err := c.conn.Read(b)
	if err != nil {
		return err
	}
	//print(n)
	serverString := string(b[:n])
	if len(serverString) > 0 && strings.HasPrefix(serverString, "yes") {
		c.currentStage = awaitingFiles
	} else {
		return errors.New("MOSS Server does not recognize this programming language")
	}
	return nil
}

// SendQuery 查询结果
// @receiver c
// @return error
// @date 2021-01-11 22:40:00
func (c *MossClient) SendQuery() (err error) {
	if c.currentStage != awaitingQuery {
		return errors.New("cannot send query at this time. Connection is either not initialized or already closed")
	} else if c.setId == 1 {
		return errors.New("you did not upload any files yet")
	} else {
		if err = c.sendCommand(fmt.Sprintf("%s %d %s", "query", 0, c.optC)); err != nil {
			return nil
		}
		c.currentStage = awaitingResults
		//serverByte := &bytes.Buffer{}
		b := make([]byte, 1024)
		n, err := c.conn.Read(b)
		if err != nil {
			return err
		}
		result := string(b[:n-1])
		if len(result) > 0 && strings.HasPrefix(strings.ToLower(result), "http") {
			if c.ResultURL, err = url.Parse(strings.Trim(result, " ")); err != nil {
				return err
			}
			c.currentStage = awaitingEnd
		} else {
			return errors.New("MOSS submission failed. The server did not return a valid URL with detection results")
		}
	}
	return nil
}

// UploadFile 上传代码源文件
// @receiver c
// @params filePath
// @params isBaseFile
// @return error
// @date 2021-01-11 22:40:09
func (c *MossClient) UploadFile(filePath string, isBaseFile bool) (err error) {
	if c.currentStage != awaitingFiles && c.currentStage != awaitingQuery {
		return errors.New("cannot upload file. Client is either not initialized properly or the connection is already closed")
	}
	fileBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	setID := 0
	if !isBaseFile {
		setID = c.setId
		c.setId++
	}
	filename := strings.ReplaceAll(filePath, "\\", "/")
	uploadString := fmt.Sprintf("file %d %s %d %s\n", setID, c.language, len(fileBytes), filename)
	if _, err = c.conn.Write([]byte(uploadString)); err != nil {
		return err
	}
	if _, err = c.conn.Write(fileBytes); err != nil {
		return err
	}
	c.currentStage = awaitingQuery
	return nil
}
