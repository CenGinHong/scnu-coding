package utils

import (
	"bytes"
	"github.com/dimchansky/utfbom"
	"golang.org/x/text/encoding/simplifiedchinese"
	"io"
	"unicode/utf8"
)

// @Author: 陈健航
// @Date: 2021/4/29 15:37
// @Description:

// RemoveBom 移除bom头和将gbk转utf8
// @params r
// @return transformedReader
// @date 2021-05-03 00:00:39
func RemoveBom(r io.Reader) (transformedReader io.Reader, err error) {
	// 去掉bom
	transformedReader = utfbom.SkipOnly(r)
	// 转utf8
	if all, err := io.ReadAll(transformedReader); err != nil {
		return nil, err
	} else if !utf8.Valid(all) {
		transformedReader = simplifiedchinese.GB18030.NewDecoder().Reader(transformedReader)
	}
	return transformedReader, nil
}

// WriteBom 写入bom头
// @params file
// @date 2021-05-08 16:44:21
func WriteBom(file *bytes.Buffer) {
	file.WriteString("\xEF\xBB\xBF")
}
