package storage

import "github.com/ectrc/snow/aid"

func GetDefaultEngine() []byte {
	return []byte(`
[OnlineSubsystemMcp]
bUsePartySystemV2=false

[OnlineSubsystemMcp.OnlinePartySystemMcpAdapter]
bUsePartySystemV2=false
	
[XMPP]
bEnableWebsockets=true

[OnlineSubsystemMcp.Xmpp]
bUseSSL=false
ServerAddr="ws://` + aid.Config.API.Host + `/socket/presence"
ServerPort=` + aid.Config.API.Port + `

[OnlineSubsystemMcp.Xmpp Prod]
bUseSSL=false
ServerAddr="ws://` + aid.Config.API.Host + `/socket/presence"
ServerPort=` + aid.Config.API.Port + ``)
}