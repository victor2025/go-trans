package handlers

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"go-trans/protocols"
	"go-trans/utils"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"
	"time"
)

type SendHandler struct {
	addr      string
	port      string
	path      string
	baseDir   string
	sliceSize uint16
}

func NewSendHandler(addr, port, path string) *SendHandler {
	return &SendHandler{
		addr:      addr,
		port:      port,
		path:      path,
		sliceSize: 1024 * 4, // max 2^16-1
	}
}

func (s *SendHandler) Handle() {
	log.Printf("--- Send mode ---")

	// connect server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", s.addr, s.port))
	utils.HandleError(err, utils.ExitOnErr)
	defer conn.Close()
	log.Printf("Connected to %s:%s", s.addr, s.port)
	// open local file
	absPath, err := filepath.Abs(s.path)
	var fileOrDirName string
	s.baseDir, fileOrDirName = filepath.Split(absPath)
	if utils.IsDir(absPath) {
		// send the whole dir
		s.walkAndSendDir(conn, absPath, fileOrDirName)
	} else {
		// send single file
		s.sendFile(conn, fileOrDirName)
	}
}

func (s *SendHandler) walkAndSendDir(conn net.Conn, dirPath string, dirPrefix string) {
	files, _ := ioutil.ReadDir(dirPath)

	for _, onefile := range files {
		if onefile.IsDir() {
			log.Printf("Processing dir %s\n", onefile.Name())
			//fmt.Println(tmpPrefix, onefile.Name(), "目录:")
			s.walkAndSendDir(conn, onefile.Name(), dirPrefix+"/"+onefile.Name())
		} else {
			s.sendFile(conn, dirPrefix+"/"+onefile.Name())
		}
	}
}

/**
发送单个文件
*/
func (s *SendHandler) sendFile(conn net.Conn, fileRelativePath string) {
	start := time.Now()
	// send filename first
	file, err := os.Open(s.baseDir + fileRelativePath)
	utils.HandleError(err, utils.ExitOnErr)
	defer file.Close()
	log.Printf("Start transferring: %v", fileRelativePath)
	bytes := protocols.StrTransMsg(fileRelativePath).Bytes()
	_, err = conn.Write(bytes)
	utils.HandleError(err, utils.ExitOnErr)

	// send file size
	stat, _ := file.Stat()
	fileSize := stat.Size()
	bytes = protocols.NumTransMsg(fileSize).Bytes()
	log.Printf("Total size : %d bytes", fileSize)
	_, err = conn.Write(bytes)
	utils.HandleError(err, utils.ExitOnErr)

	// read file and send
	buf := make([]byte, s.sliceSize)
	seq := 0
	dataSize := 0
	md5Chk := md5.New()
	for {
		// read from file
		n, err := file.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			utils.HandleError(err)
		}

		// send to conn
		trans := protocols.ByteTransMsg(buf[:n])
		_, err = conn.Write(trans.Bytes())
		utils.HandleError(err, utils.ExitOnErr)
		md5Chk.Write(buf[:n])

		// control transmit speed
		// time.Sleep(time.Millisecond * 10000)

		// show status
		seq++
		dataSize += n
		log.Printf("seq: %v, sent %d/%dKB(%.2f%%)", seq, dataSize/1024, fileSize/1024, 100*float64(dataSize)/float64(fileSize))

		// is end
		if n < int(s.sliceSize) {
			break
		}

	}
	// send end flag
	bytes = protocols.EndTransMsg(md5Chk.Sum(nil)).Bytes() // transmit md5
	_, err = conn.Write(bytes)
	utils.HandleError(err, utils.ExitOnErr)

	// show end status
	dur := float32(time.Since(start).Microseconds()) / 1000
	avgSpeed := float32(dataSize) / (1024 * dur / 1000)
	md5 := hex.EncodeToString(md5Chk.Sum(nil))
	log.Printf("--- Send file complete ---")
	log.Printf("Filepath: %s", file.Name())
	log.Printf("MD5: %s", md5)
	log.Printf("Info: total time: %.2fms, avg speed: %.2fKB/s\n", dur, avgSpeed)
}
