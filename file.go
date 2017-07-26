package yksdk

type FileSDK struct {
    Base BaseSDK
}

func NewFileSDK(config *Config) *FileSDK {
    sdk := &FileSDK{}
    sdk.Base.config = config
    return sdk
}

func (this *FileSDK) Post(api string, params map[string]string) (*Result, error) {
    params["org_client_id"] = this.Base.config.ClientId
    return this.Base.Post(api, params)
}