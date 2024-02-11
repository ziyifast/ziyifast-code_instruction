package main

/*
	封装NewXXX函数
*/

type Protocol string

var (
	SmbProtocol Protocol = "smb"
	NfsProtocol Protocol = "nfs"
)

type IDownload interface {
	Download()
}

type SmbDownloader struct{}

func (s *SmbDownloader) Download() {
	println("smb download")
}

type NfsDownloader struct{}

func (n *NfsDownloader) Download() {
	println("nfs download")
}

func NewDownloader(t Protocol) IDownload {
	switch t {
	case SmbProtocol:
		return &SmbDownloader{}
	case NfsProtocol:
		return &NfsDownloader{}
	}
	return nil
}

func main() {
	//测试：根据协议类型，创建不同类型的下载器
	smbDownloader := NewDownloader(SmbProtocol)
	smbDownloader.Download()

	nfsDownloader := NewDownloader(NfsProtocol)
	nfsDownloader.Download()
}
