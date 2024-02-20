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
[OnlineSubsystemMcp.OnlinePaymentServiceMcp Fortnite]
Domain="launcher-website-prod.ak.epicgames.com"
BasePath="/logout?redirectUrl=https%3A%2F%2Fwww.unrealengine.com%2Fid%2Flogout%3FclientId%3Dxyza7891KKDWlczTxsyy7H3ExYgsNT4Y%26responseType%3Dcode%26redirectUrl%3Dhttps%253A%252F%252Ftesting-site.neonitedev.live%252Fid%252Flogin%253FredirectUrl%253Dhttps%253A%252F%252Ftesting-site.neonitedev.live%252Fpurchase%252Facquire&path="

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
PingTimeout=3.0

[LwsWebSocket]
bDisableCertValidation=true
bDisableDomainWhitelist=true

[WinHttpWebSocket]
bDisableDomainWhitelist=true`

	if aid.Config.Fortnite.Season <= 2 {
		str += `
		
[OnlineSubsystemMcp.Xmpp]
bUsePlainTextAuth=true
bUseSSL=false
Protocol=tcp
ServerAddr="`+ aid.Config.API.Host + `"
ServerPort=`+ realPort + `

[OnlineSubsystemMcp.Xmpp Prod]
bUsePlainTextAuth=true
bUseSSL=false
Protocol=tcp
ServerAddr="`+ aid.Config.API.Host + `"
ServerPort=`+ realPort
	} else {
		str += `
[OnlineSubsystemMcp.Xmpp]
bUsePlainTextAuth=true
bUseSSL=false
Protocol=ws
ServerAddr="ws://`+ aid.Config.API.Host + aid.Config.API.Port +`/?SNOW_SOCKET_CONNECTION"

[OnlineSubsystemMcp.Xmpp Prod]
bUsePlainTextAuth=true
bUseSSL=false
Protocol=ws
ServerAddr="ws://`+ aid.Config.API.Host + aid.Config.API.Port +`/?SNOW_SOCKET_CONNECTION"`
	}

	return []byte(str)
}

func GetDefaultGame() []byte {return []byte(`
[/Script/FortniteGame.FortGlobals]
bAllowLogout=false

[/Script/FortniteGame.FortChatManager]
bShouldRequestGeneralChatRooms=false
bShouldJoinGlobalChat=false
bShouldJoinFounderChat=false
bIsAthenaGlobalChatEnabled=false

[/Script/FortniteGame.FortOnlineAccount]
bEnableEulaCheck=false
bShouldCheckIfPlatformAllowed=false

[EpicPurchaseFlow]
bUsePaymentWeb=false
CI="http://localhost:5173/purchase"
GameDev="http://localhost:5173/purchase"
Stage="http://127.0.0.1:5173/purchase"
Prod="http://127.0.0.1:5173/purchase"
UEPlatform="FNGame"

`)}

func GetDefaultRuntime() []byte {return []byte(`
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
+ExperimentalCohortPercent=(CohortPercent=100,ExperimentNum=20)`)}