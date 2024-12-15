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
	log.Printf("Listening to:%s, Serve on:\n", s.port)
	s.listener = listener
	// show local ip
	addrs, err := net.InterfaceAddrs()
	utils.HandleError(err, utils.ExitOnErr)
	addrIdx := 0
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil || ipnet.IP.To16() != nil {
				addrIdx++
				log.Printf("%d:\t%s\n", addrIdx, ipnet.IP.String())
			}
		}
	}

}

func (s *ReceiveHandler) serveConn(conn net.Conn) {
	log.Printf("--- New connection from %s ---", conn.RemoteAddr())
	defer conn.Close()
	start := time.Now()
	totalDataSize := int64(0)
	for {
		trans, err := protocols.ReceiveNextTrans(conn)
		utils.HandleError(err, utils.ExitOnErr)
		if trans.Head.Type == protocols.DisconnectType {
			break
		}
		if trans.Head.Type == protocols.NewFileType {
			dataSize, err := s.receiveNewFile(conn)
			utils.HandleError(err, utils.DoNothingOnErr)
			totalDataSize += dataSize
		}
	}
	dur := float32(time.Since(start).Microseconds()) / 1000
	avgSpeed := float32(totalDataSize) / (1024 * dur / 1000)
	log.Printf("--- Info: receive file complete, total time: %.2fms, avg speed %.2fKB/s ---\n", dur, avgSpeed)
}

func (s *ReceiveHandler) receiveNewFile(conn net.Conn) (int64, error) {
	start := time.Now()
	var err error

	// init file
	file, fileSize, err := s.initFile(conn)
	utils.HandleError(err)
	if file == nil || err != nil {
		log.Printf("WARN: init file failed, error:%s", err)
	}
	defer file.Close()

	// write file
	seq := 0
	dataSize := int64(0)
	var rcvdMd5 string
	md5Chk := md5.New()
	for {
		// read bytes from connection
		trans, err := protocols.ReceiveNextTrans(conn)
		utils.HandleError(err)
		if err != nil {
			return 0, err
		}

		// if is ended
		if trans.Head.Type == protocols.EndType {
			rcvdMd5 = hex.EncodeToString(trans.Content) // get remote md5Val
			break
		}

		// write bytes to file
		_, err = file.Write(trans.Content)
		utils.HandleError(err)
		if err != nil {
			return 0, err
		}
		md5Chk.Write(trans.Content)

		// show status
		seq++
		dataSize += int64(len(trans.Content))
		log.Printf("seq: %v, received %d/%dKB(%.2f%%)", seq, dataSize/1024, fileSize/1024, 100*float64(dataSize)/float64(fileSize))

	}

	// show end status
	dur := float32(time.Since(start).Microseconds()) / 1000
	avgSpeed := float32(dataSize) / (1024 * dur / 1000)
	md5Val := hex.EncodeToString(md5Chk.Sum(nil))
	log.Printf("--- Receive file complete ---")
	log.Printf("Filepath: %s", file.Name())
	log.Printf("MD5: %s", md5Val)
	log.Printf("Info: cost time: %.2fms, avg speed %.2fKB/s\n", dur, avgSpeed)
	if md5Val != rcvdMd5 {
		log.Printf("WARN: File md5Val is different, please check manually!")
		return 0, fmt.Errorf("file md5 not match fileMd5:%s, actualMd5:%s", md5Val, rcvdMd5)
	}
	return dataSize, nil
}

/**
接收文件名，创建目标文件
*/
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
