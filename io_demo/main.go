package main

import (
	"context"
	"fmt"
	"github.com/aobco/log"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"golang.org/x/sync/semaphore"
	"io"
	"os"
	"time"
	"xsky.com/ocpf/common/osg"
)

/*
	通过io.Copy实现可中断的流复制
*/
var (
	ak       = "xxxx"
	sk       = "xxxxxxx"
	endpoint = "http://xx.xx.xx.xx:xxx"
	bucket   = "test-bucket"
	key      = "d_xp/2G/2G.txt"
)

func main() {
	s3Client := osg.Client.GetS3Client(ak, sk, endpoint)
	ctx, cancelFunc := context.WithCancel(context.Background())
	object, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	go func() {
		time.Sleep(time.Second * 10)
		cancelFunc()
		log.Infof("canceled...")
	}()
	if err != nil {
		log.Errorf("%v", err)
		return
	}
	body := object.Body
	defer body.Close()
	file, err := os.Create("/Users/ziyi2/GolandProjects/MyTest/demo_home/io_demo/target.txt")
	if err != nil {
		log.Errorf("%v", err)
		return
	}
	defer file.Close()
	_, err = FileService.Copy(ctx, file, body)
	if err != nil {
		log.Errorf("%v", err)
		return
	}

}

type fileService struct {
	sem *semaphore.Weighted
}

var FileService = &fileService{
	sem: semaphore.NewWeighted(1),
}

type IoCopyCancelledErr struct {
	errMsg string
}

func (e *IoCopyCancelledErr) Error() string {
	return fmt.Sprintf("io copy error, %s", e.errMsg)
}

func NewIoCopyCancelledErr(msg string) *IoCopyCancelledErr {
	return &IoCopyCancelledErr{
		errMsg: msg,
	}
}

type readerFunc func(p []byte) (n int, err error)

func (rf readerFunc) Read(p []byte) (n int, err error) { return rf(p) }

//通过ctx实现可中断的流拷贝
// Copy closable copy
func (s *fileService) Copy(ctx context.Context, dst io.Writer, src io.Reader) (int64, error) {
	// Copy will call the Reader and Writer interface multiple time, in order
	// to copy by chunk (avoiding loading the whole file in memory).
	// I insert the ability to cancel before read time as it is the earliest
	// possible in the call process.
	size, err := io.Copy(dst, readerFunc(func(p []byte) (int, error) {
		select {
		// if context has been canceled
		case <-ctx.Done():
			// stop process and propagate "context canceled" error
			return 0, NewIoCopyCancelledErr(ctx.Err().Error())
		default:
			// otherwise just run default io.Reader implementation
			return src.Read(p)
		}
	}))
	return size, err
}
