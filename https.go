package sml

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"net"
	"strings"
	"io"
	"os"
	"time"
	"reflect"
	"crypto/tls"
	"bytes"
	"encoding/base64"
	"mime/multipart"
)

type UpFile struct{
	FormName string
	FileName string
	Input io.Reader
}

type Https struct{
	url string
	insecureSkipVerify bool
	httpClient *http.Client
	method string
	requestType string
	bodyWriter *multipart.Writer
	body io.Reader
	upBodyBuffer *bytes.Buffer
	connectTimeout time.Duration  
	requestHeader http.Header
	rwTimeout time.Duration    
	queryParam url.Values
	requestBody interface{}
	request *http.Request
	Response *http.Response
}
func newHttps(urlStr string) *Https{
	https:= &Https{method:http.MethodGet,requestHeader:http.Header{},queryParam:url.Values{},url:urlStr,connectTimeout:5000*time.Millisecond,insecureSkipVerify:true}
	return https.Form()
}

func NewGetHttps(urlStr string) *Https{
	https:=newHttps(urlStr).Method(http.MethodGet)
	return https
}
func NewPostFormHttps(urlStr string) *Https{
	https:=newHttps(urlStr).Method(http.MethodPost).Form()
	return https
}
func NewPostBodyHttps(urlStr string) *Https{
	https:=newHttps(urlStr).Method(http.MethodPost).Json()
	return https
}
func (https *Https) Multipart() *Https{
	https.requestType="multipart"
	https.upBodyBuffer=&bytes.Buffer{}
    https.bodyWriter= multipart.NewWriter(https.upBodyBuffer)
    contentType := https.bodyWriter.FormDataContentType()
	return https.Header("Content-Type", contentType)
}
func (https *Https) Form() *Https{
	https.requestType="form"
	return https.Header("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
}
func (https *Https) Json()  *Https{
	https.requestType="json"
	return https.Header("Content-Type","application/json; charset=UTF-8")
}

func (https *Https) Method(method string) *Https{
	https.method=method
	return https
}

func (https *Https) Param(name string,value string) *Https{
	https.queryParam.Add(name,value)
	return https
}
func (https *Https) Header(name string,value string) *Https{
	https.requestHeader.Set(name,value)
	return https
}
func (https *Https) ConnectTimeout(connectTimeout time.Duration) *Https{
	https.connectTimeout=connectTimeout;
	return https
}
func (https *Https) RWTimeout(rwTimeout time.Duration) *Https{
	https.rwTimeout=rwTimeout;
	return https
}
//string []byte  io.reader support
func (https *Https) Body(body interface{}) *Https{
	kind:=reflect.ValueOf(body).Kind()
	switch kind{
		case reflect.String: {
			https.body=strings.NewReader(reflect.ValueOf(body).String())
		}
		case reflect.Slice: {
			https.body=bytes.NewReader(reflect.ValueOf(body).Bytes())
		}
		default :{
			if(https.requestType!="multipart"){
				if readerBody,ok:=body.(io.Reader);ok{
					https.body=readerBody
				}
			}
		}
		
	}
	return https
}
func (https *Https) UpFile(upFiles ... *UpFile) *Https{
	for _,upFile:=range upFiles{
		iwriter,_:=https.bodyWriter.CreateFormFile(upFile.FormName, upFile.FileName)
		io.Copy(iwriter,upFile.Input)
	}
	return https
}
func (https *Https) clientHandler() *http.Client{
	c := &http.Client{
        Transport: &http.Transport{
            TLSClientConfig: https.tlsConfig(),
            Dial:https.timeoutDialer(),
        },
    }
	return c
}
func (https *Https) urlHandler(){
	 if((https.method==http.MethodGet||https.method==http.MethodDelete)||(https.method==http.MethodPost&&https.requestType!="form"&&https.requestType!="multipart") ){
		 queryString:=https.queryParam.Encode()
		 if len(queryString)>2{
			 if strings.Contains(https.url,"?"){
				 https.url=https.url+"&"+queryString
			 }else{
				  https.url=https.url+"?"+queryString
			 }
		 }
	 }else if(https.method==http.MethodPost&&https.requestType=="form"){
		 https.Body(https.queryParam.Encode())
	 }else if(https.requestType=="multipart"){
		 for name,value:=range https.queryParam{
			 https.bodyWriter.WriteField(name, value[0])
		 }
	 }
}
func (https *Https) Auth(authType string,credentials string) *Https{
	https.Header("Authorization",authType+" "+credentials)
	return https
}
func (https *Https) BasicAuth(credentials string) *Https{
	return https.Auth("Basic",base64.StdEncoding.EncodeToString([]byte(credentials)))
}
func (https *Https) headerHandler(){
	for name,value:=range https.requestHeader{
		https.request.Header.Set(name,value[0])
	}
}

func (https *Https) Execute() error{
	client:=https.clientHandler()
	https.urlHandler()
	var err1 error
	https.request,err1=http.NewRequest(https.method,https.url,nil)
	if err1!=nil{
		return err1
	}
	https.headerHandler()
	switch https.method{
		case http.MethodGet,http.MethodDelete:{
			
		}
		case http.MethodPost,http.MethodPut:{
			if(https.requestType=="multipart"){
				https.bodyWriter.Close()
				https.request.Body=ioutil.NopCloser(https.upBodyBuffer)
			}else{
				https.request.Body=ioutil.NopCloser(https.body)
			}
		}
	}
	resp,err:=client.Do(https.request)
	https.Response=resp
	return err
}
func (https *Https) GetBodyString() string{
	return string(https.GetBodyBytes())
}

func (https *Https) GetBodyBytes() []byte{
	defer https.Response.Body.Close()
	data,err:=ioutil.ReadAll(https.Response.Body)
	if err!=nil{
		panic(err)
	}
	return data
}

func (https *Https) GetBodyTo(writer io.Writer) int64{
	body:=https.Response.Body
	defer body.Close()
	l,err:=io.Copy(writer,body)
	if err!=nil{
		panic(err)
	}
	return l
}
//file not close 
func (https *Https) GetBodyToFile(file *os.File) {
	body:=https.Response.Body
	defer body.Close()
	bs:=make([]byte,512)
	for{
		n,err:=body.Read(bs)
		file.Write(bs[:n])
        if(err!=nil&&err==io.EOF){
	        break
        }else if(err!=nil){
        	panic( err)
        }	
	}
}
func (https *Https) GetBodyTo1(){
	
}


func (https *Https) timeoutDialer() func(net, addr string) (c net.Conn, err error) {
    return func(netw, addr string) (net.Conn, error) {
	    conn, err := net.DialTimeout(netw, addr,https.connectTimeout)
	    if err != nil {
	         return nil, err
	    }
	    if(https.rwTimeout>time.Millisecond){
	        conn.SetDeadline(time.Now().Add(https.rwTimeout))
	    }
        return conn, nil
    }
}
func (https *Https) tlsConfig() *tls.Config{
	if strings.HasPrefix(https.url,"https"){
		return &tls.Config{InsecureSkipVerify: https.insecureSkipVerify}
	}
	return nil
}
