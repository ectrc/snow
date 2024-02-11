package storage

import (
	"github.com/ectrc/snow/aid"
)

func GetDefaultEngine() []byte {
	return []byte(`
[OnlineSubsystemMcp.Xmpp]
bUseSSL=false
Protocol=ws
ServerAddr="ws://`+ aid.Config.API.Host + aid.Config.API.Port +`/?"

[OnlineSubsystemMcp.Xmpp Prod]
bUseSSL=false
Protocol=ws
ServerAddr="ws://`+ aid.Config.API.Host + aid.Config.API.Port +`/?"

[OnlineSubsystemMcp]
bUsePartySystemV2=true

[OnlineSubsystemMcp.OnlinePartySystemMcpAdapter]
bUsePartySystemV2=true

[XMPP]
bEnableWebsockets=true

[ConsoleVariables]
n.VerifyPeer=0
FortMatchmakingV2.ContentBeaconFailureCancelsMatchmaking=0
Fort.ShutdownWhenContentBeaconFails=0
FortMatchmakingV2.EnableContentBeacon=0

[/Script/Qos.QosRegionManager]
NumTestsPerRegion=5
PingTimeout=3.0`)
}

func GetDefaultGame() []byte {
	return []byte(`
[/Script/FortniteGame.FortGlobals]
bAllowLogout=false

[/Script/FortniteGame.FortChatManager]
bShouldRequestGeneralChatRooms=false
bShouldJoinGlobalChat=false
bShouldJoinFounderChat=false
bIsAthenaGlobalChatEnabled=false

[/Script/FortniteGame.FortOnlineAccount]
bEnableEulaCheck=false
bShouldCheckIfPlatformAllowed=false`)
}

func GetDefaultRuntime() []byte {
	return []byte(`
[/Script/FortniteGame.FortRuntimeOptions]
bEnableGlobalChat=true
bDisableGifting=false
bDisableGiftingPC=false
bDisableGiftingPS4=false
bDisableGiftingXB=false
!ExperimentalCohortPercent=ClearArray
+ExperimentalCohortPercent=(CohortPercent=100,ExperimentNum=20)`)
}