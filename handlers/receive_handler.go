package handlers

import (
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"go-trans/protocols"
	"go-trans/utils"
	"log"
	"net"
	"os"
	"path/filepath"
	"time"
)

type ReceiveHandler struct {
	port     string
	basePath string
	listener net.Listener
}

func NewReceiveHandler(port, basePath string) *ReceiveHandler {
	path, err := filepath.Abs(basePath)
	utils.HandleError(err, utils.ExitOnErr)
	return &ReceiveHandler{
		port:     port,
		basePath: path + "/",
	}
}

func (s *ReceiveHandler) Handle() {
	// 开启服务
	s.startServer()
	defer s.listener.Close() // 退出时关闭服务

	// 循环处理请求
	for {
		isNormal := true
		// wait for connection
		conn, err := s.listener.Accept()
		utils.HandleError(err, func() { isNormal = false })
		if !isNormal {
			continue
		}

		// serve connection
		go s.serveConn(conn)

	}
}

func (s *ReceiveHandler) startServer() {
	log.Printf("--- Receive mode ---")
	// start tcp listen
	listener, err := net.Listen("tcp", ":"+s.port)
	utils.HandleError(err, utils.ExitOnErr)
	log.Printf("Listening :%s\n", s.port)
	s.listener = listener
}

func (s *ReceiveHandler) serveConn(conn net.Conn) {
	log.Printf("--- New connection ---")
	start := time.Now()
	defer conn.Close()
	var err error

	// init file
	file, fileSize, err := s.initFile(conn)
	utils.HandleError(err)
	if file != nil {
		defer file.Close()
	}
	if err != nil {
		return
	}

	// write file
	seq := 0
	dataSize := 0
	var rcvdMd5 string
	md5Chk := md5.New()
	for {
		// read bytes from connection
		trans, err := protocols.ReceiveNextTrans(conn)
		utils.HandleError(err)
		if err != nil {
			return
		}

		// if is end
		if trans.Head.Type == protocols.EndType {
			rcvdMd5 = hex.EncodeToString(trans.Content) // get remote md5
			break
		}

		// write bytes to file
		_, err = file.Write(trans.Content)
		utils.HandleError(err)
		if err != nil {
			return
		}
		md5Chk.Write(trans.Content)

		// show status
		seq++
		dataSize += len(trans.Content)
		log.Printf("seq: %v, received %d/%dKB(%.2f%%)", seq, dataSize/1024, fileSize/1024, 100*float64(dataSize)/float64(fileSize))

	}

	// show end status
	dur := float32(time.Since(start).Microseconds()) / 1000
	avgSpeed := float32(dataSize) / (1024 * dur / 1000)
	md5 := hex.EncodeToString(md5Chk.Sum(nil))
	log.Printf("--- Receive file complete ---")
	log.Printf("Filepath: %s", file.Name())
	log.Printf("MD5: %s", md5)
	log.Printf("Info: total time: %.2fms, avg speed %.2fKB/s\n", dur, avgSpeed)
	if md5 != rcvdMd5 {
		log.Printf("WARN: File md5 is different, please check manually!")
	}
}

// return file, file size, error
func (r *ReceiveHandler) initFile(conn net.Conn) (*os.File, int64, error) {
	var err error

	// read filename trans msg
	trans, err := protocols.ReceiveNextTrans(conn)
	utils.HandleError(err)
	if err != nil {
		return nil, 0, err
	}
	if trans.Head.Type != protocols.StrType {
		return nil, 0, fmt.Errorf("error get file")
	}

	// receive filename
	filename := string(trans.Content)
	log.Printf("Receiving file: %s", filename)
	path, err := filepath.Abs(r.basePath + filename)
	utils.HandleError(err)
	log.Printf("Saving file to: %s", path)
	dir, _ := filepath.Split(path)
	os.MkdirAll(dir, 0777)

	// create file
	file, err := os.Create(path)
	utils.HandleError(err)
	if err != nil {
		return nil, 0, err
	}

	// read file size trans msg
	trans, err = protocols.ReceiveNextTrans(conn)
	utils.HandleError(err)
	if err != nil {
		return nil, 0, err
	}
	if trans.Head.Type != protocols.NumType {
		return nil, 0, fmt.Errorf("error get file size")
	}
	fileSize, _ := binary.Varint(trans.Content)
	log.Printf("Total size : %d bytes", fileSize)

	return file, fileSize, nil
}
