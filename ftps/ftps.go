package ftps

import (
	"github.com/huangwenjimmy/gosml"
	"github.com/jlaffaye/ftp"
	"io"
	"time"
)

type Ftps struct {
	Urls *gosml.Urls
	Conn *ftp.ServerConn
}

func (f *Ftps) Login() *Ftps {
	var err error
	to := gosml.ConvertToInt(f.Urls.Params["timeout"])
	if to > 0 {
		f.Conn, err = ftp.DialTimeout(f.Urls.Host+":"+gosml.ConvertToString(f.Urls.Port), time.Duration(to)*time.Millisecond)
	} else {
		f.Conn, err = ftp.Connect(f.Urls.Host + ":" + gosml.ConvertToString(f.Urls.Port))
	}
	gosml.ThrowRuntime(err)
	err = f.Conn.Login(f.Urls.Username, f.Urls.Password)
	gosml.ThrowRuntime(err)
	return f
}
func (f *Ftps) Cd() *Ftps {
	err := f.Conn.ChangeDir(f.Urls.LastPath)
	gosml.ThrowRuntime(err)
	return f
}
func (f *Ftps) Ls() []string {
	ns, err := f.Conn.NameList(f.Urls.LastPath)
	gosml.ThrowRuntime(err)
	return ns
}
func (f *Ftps) Put(ir io.Reader) error {
	return f.Conn.Stor(f.Urls.Path, ir)
}
func (f *Ftps) Get() (io.Reader, error) {
	return f.Conn.Retr(f.Urls.Path)
}
func (f *Ftps) GetTo(iw io.Writer) error {
	resp, err := f.Conn.Retr(f.Urls.Path)
	defer resp.Close()
	io.Copy(iw, resp)
	return err
}
func (f *Ftps) DisConnect() error {
	return f.Conn.Quit()
}
func NewFtps(url string) *Ftps {
	urls := gosml.NewUrls(url)
	f := &Ftps{Urls: urls}
	return f
}
