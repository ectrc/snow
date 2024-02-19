package storage

import (
	"fmt"
	"strconv"

	"github.com/ectrc/snow/aid"
)

func GetDefaultEngine() []byte {
	portNumber, err := strconv.Atoi(aid.Config.API.Port[1:])
	if err != nil {
		return nil
	}
	portNumber++
	realPort := fmt.Sprintf("%d", portNumber)

	str := `
[XMPP]
bEnableWebsockets=true

[OnlineSubsystem]
bHasVoiceEnabled=true

[ConsoleVariables]
n.VerifyPeer=0
FortMatchmakingV2.ContentBeaconFailureCancelsMatchmaking=0
Fort.ShutdownWhenContentBeaconFails=0
FortMatchmakingV2.EnableContentBeacon=0

[/Script/Qos.QosRegionManager]
NumTestsPerRegion=5
PingTimeout=3.0`

	if aid.Config.Fortnite.Season <= 2 {
		str += `
		
[OnlineSubsystemMcp.Xmpp]
bUseSSL=false
Protocol=tcp
ServerAddr="`+ aid.Config.API.Host + `"
ServerPort=`+ realPort + `

[OnlineSubsystemMcp.Xmpp Prod]
bUseSSL=false
Protocol=tcp
ServerAddr="`+ aid.Config.API.Host + `"
ServerPort=`+ realPort
	} else {
		str += `
[OnlineSubsystemMcp.Xmpp]
bUseSSL=false
Protocol=ws
ServerAddr="ws://`+ aid.Config.API.Host + aid.Config.API.Port +`/?"

[OnlineSubsystemMcp.Xmpp Prod]
bUseSSL=false
Protocol=ws
ServerAddr="ws://`+ aid.Config.API.Host + aid.Config.API.Port +`/?"`
	}

	return []byte(str)
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
!DisabledFrontendNavigationTabs=ClearArray
;+DisabledFrontendNavigationTabs=(TabName="AthenaChallenges",TabState=EFortRuntimeOptionTabState::Hidden)
;+DisabledFrontendNavigationTabs=(TabName="Showdown",TabState=EFortRuntimeOptionTabState::Hidden)
;+DisabledFrontendNavigationTabs=(TabName="AthenaStore",TabState=EFortRuntimeOptionTabState::Hidden)

bEnableGlobalChat=true
bDisableGifting=false
bDisableGiftingPC=false
bDisableGiftingPS4=false
bDisableGiftingXB=false
!ExperimentalCohortPercent=ClearArray
+ExperimentalCohortPercent=(CohortPercent=100,ExperimentNum=20)`)
}