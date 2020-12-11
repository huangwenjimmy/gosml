package gosml

import (
	"net/url"
	"strings"
)

type Urls struct {
	Url      string
	Schema   string
	Host     string
	Username string
	Password string
	IsFile   bool
	Path     string
	PrePath  string
	LastPath string
	Query    string
	Port     int
	Params   map[string]string
}

func (this *Urls) initUrl() {
	vs := strings.SplitN(this.Url, "://", 2)
	if len(vs) < 2 {
		panic("url [" + this.Url + "] err")
	}
	this.Schema = vs[0]
	sso := vs[1]
	splitIndex := strings.Index(sso, "/")
	if splitIndex > -1 {
		uri := SubStr(sso, splitIndex+1, len(sso))
		queryIndex := strings.Index(uri, "?")
		if queryIndex == -1 {
			this.Path = uri
		} else {
			this.Path = SubStr(uri, 0, queryIndex)
			this.Query = SubStr(uri, queryIndex+1, len(uri))
			this.Params = this.decodeQuery()
		}
		lastIndex := strings.LastIndex(this.Path, "/")
		this.LastPath = SubStr(this.Path, lastIndex+1, len(this.Path))
		this.IsFile = strings.Index(this.LastPath, ".") > -1
		if this.IsFile {
			if this.LastPath != this.Path {
				this.PrePath = SubStr(this.Path, 0, strings.LastIndex(this.Path, "/"))
			} else {
				this.PrePath = ""
			}
		} else {
			this.PrePath = this.Path
		}
	}
	hostIndex := len(sso)
	if splitIndex > -1 {
		hostIndex = splitIndex
	}
	hostInfos := SubStr(sso, 0, hostIndex)
	authIndex := strings.Index(hostInfos, "@")
	addr := hostInfos
	if authIndex > -1 {
		addr = SubStr(hostInfos, authIndex+1, len(hostInfos))
		auth := SubStr(hostInfos, 0, authIndex)
		as := strings.SplitN(auth, ":", 2)
		if len(as) == 1 {
			this.Password = as[0]
		} else {
			this.Username = as[0]
			this.Password = as[1]
		}
	}
	this.Host = addr
	if strings.Contains(addr, ":") {
		as := strings.SplitN(addr, ":", 2)
		this.Port = int(ConvertToInt(as[1]))
		this.Host = as[0]
	}
}
func (this *Urls) decodeQuery() map[string]string {
	result := make(map[string]string, 0)
	kv := strings.Split(this.Query, "&")
	for _, kvs := range kv {
		pkv := strings.SplitN(kvs, "=", 2)
		if len(pkv) == 2 {
			result[pkv[0]], _ = url.QueryUnescape(pkv[1])
		}
	}
	return result
}

func NewUrls(url string) *Urls {
	fs := &Urls{Url: url}
	fs.initUrl()
	return fs
}
