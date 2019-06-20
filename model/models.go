package model

// Result struct
type PushAppResult struct {
	App                    App
	PushSucced             bool `json:"pushSucced"`
	PushElapsed            int  `json:"pushElapsed"`
	HTTPVerificationSucced bool `json:"httpVerificationSucced"`
}

// App struct
type App struct {
	Name   string `json:"name"`
	Domain string `json:"domain"`
}
