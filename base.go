package yksdk

import (
    "github.com/tidwall/gjson"
    "strconv"
    "repo.gokuai.cn/golang/common/mcrypt"
    "net/url"
    "net/http"
    "io/ioutil"
    "encoding/json"
    "time"
)

type Config struct {
    Scheme string
    Host string
    UriPrefix string
    ClientId string
    ClientSecret string
}

type Error struct {
    ErrorCode int `json:"error_code"`
    ErrorMsg string `json:"error_msg"`
}

type Result struct {
    error *Error
    data *gjson.Result
}

type BaseSDK struct {
    config *Config
    ResponseStatusCode int
    ResponseBody       []byte
}

func (this *BaseSDK) ResponseToString() string {
    result := strconv.Itoa(this.ResponseStatusCode)
    if this.ResponseBody != nil {
        result += ":" + string(this.ResponseBody)
    }
    return result
}

func (this *BaseSDK) Post(api string, params map[string]string) (*Result, error) {
    this.ResponseStatusCode = 0
    this.ResponseBody = nil

    params["dateline"] = strconv.FormatInt(time.Now().Unix(), 10)
    sign := mcrypt.Signature(params, this.config.ClientSecret)
    params["sign"] = sign
    form := url.Values{}
    for k, v := range params {
        form.Set(k, v)
    }
    apiurl :=  this.config.Scheme + "://" + this.config.Host + this.config.UriPrefix + api
    resp, err := http.PostForm(apiurl, form)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    this.ResponseStatusCode = resp.StatusCode
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }
    result := &Result{}
    this.ResponseBody = body
    if resp.StatusCode == 200 {
        jsonResult := gjson.Parse(string(this.ResponseBody))
        result.data = &jsonResult
    } else {
        result.error = &Error{}
        err = json.Unmarshal(this.ResponseBody, result.error)
        if err != nil {
            return nil, err
        }
    }
    return result , nil
}

func (this *Result) IsOK() (bool, *Error) {
    if this.error == nil {
        return true, nil
    }
    return false, this.error
}

func (this *Result) GetResult() *gjson.Result {
    return this.data
}