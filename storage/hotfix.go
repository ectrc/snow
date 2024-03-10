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
PingTimeout=3.0

[/Script/Qos.QosRegionManager]
NumTestsPerRegion=5
PingTimeout=3.0
!RegionDefinitions=ClearArray
+RegionDefinitions=(DisplayName=NSLOCTEXT("MMRegion", "Europe", "Europe"), RegionId="EU", bEnabled=true, bVisible=true, bAutoAssignable=true)
+RegionDefinitions=(DisplayName=NSLOCTEXT("MMRegion", "North America", "North America"), RegionId="NA", bEnabled=true, bVisible=true, bAutoAssignable=true)
+RegionDefinitions=(DisplayName=NSLOCTEXT("MMRegion", "Oceania", "Oceania"), RegionId="OCE", bEnabled=true, bVisible=true, bAutoAssignable=true)
!DatacenterDefinitions=ClearArray
+DatacenterDefinitions=(Id="DE", RegionId="EU", bEnabled=true, Servers[0]=(Address="142.132.145.234", Port=22222))
+DatacenterDefinitions=(Id="VA", RegionId="NA", bEnabled=true, Servers[0]=(Address="69.10.34.38", Port=22222))
+DatacenterDefinitions=(Id="SYD", RegionId="OCE", bEnabled=true, Servers[0]=(Address="139.99.209.91", Port=22222))
!Datacenters=ClearArray
+Datacenters=(DisplayName=NSLOCTEXT("MMRegion", "Europe", "Europe"), RegionId="EU", bEnabled=true, bVisible=true, bBeta=false, Servers[0]=(Address="142.132.145.234", Port=22222))
+Datacenters=(DisplayName=NSLOCTEXT("MMRegion", "North America", "North America"), RegionId="NA", bEnabled=true, bVisible=true, bBeta=false, Servers[0]=(Address="69.10.34.38", Port=22222))
+Datacenters=(DisplayName=NSLOCTEXT("MMRegion", "Oceania", "Oceania"), RegionId="OCE", bEnabled=true, bVisible=true, bBeta=false, Servers[0]=(Address="139.99.209.91", Port=22222))

[LwsWebSocket]
bDisableCertValidation=true
bDisableDomainWhitelist=true

[/Script/Engine.NetworkSettings]
n.VerifyPeer=false

[WinHttpWebSocket]
bDisableCertValidation=true
bDisableDomainWhitelist=true`

	if aid.Config.Fortnite.Season <= 2 {
		str += `
[OnlineSubsystemMcp.Xmpp]
bUsePlainTextAuth=true
bUseSSL=false
Protocol=tcp
ServerAddr="`+ aid.Config.API.XMPP.Host + aid.Config.API.XMPP.Port + `"
ServerPort=`+ realPort + `

[OnlineSubsystemMcp.Xmpp Prod]
bUsePlainTextAuth=true
bUseSSL=false
Protocol=tcp
ServerAddr="`+ aid.Config.API.XMPP.Host + aid.Config.API.XMPP.Port + `"
ServerPort=`+ realPort
	} else {
		str += `
[OnlineSubsystemMcp.Xmpp]
bUseSSL=false
Protocol=ws
ServerAddr="ws://`+ aid.Config.API.XMPP.Host + aid.Config.API.XMPP.Port +`/?SNOW_SOCKET_CONNECTION"

[OnlineSubsystemMcp.Xmpp Prod]
bUseSSL=false
Protocol=ws
ServerAddr="ws://`+ aid.Config.API.XMPP.Host + aid.Config.API.XMPP.Port +`/?SNOW_SOCKET_CONNECTION"`
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
CI="http://127.0.0.1:3000/purchase"
GameDev="http://127.0.0.1:3000/purchase"
Stage="http://127.0.0.1:3000/purchase"
Prod="http://127.0.0.1:3000/purchase"
UEPlatform="FNGame"

[/Script/FortniteGame.FortTextHotfixConfig]
+TextReplacements=(Category=Game, bIsMinimalPatch=True, Namespace="", Key="68ADE44C49B20BFF78677799BE68B0EE", NativeString="FORTNITEMARES", LocalizedStrings=(("en","BOOST PERKS")))
+TextReplacements=(Category=Game, bIsMinimalPatch=True, Namespace="", Key="BE6B17BD456F3F13EEB2998AF91DC717", NativeString="THANKS FOR PLAYING!", LocalizedStrings=(("en","THANKS FOR SUPPORTING SNOW!")))

[/Script/FortniteGame.FortGameInstance]
!FrontEndPlaylistData=ClearArray
+FrontEndPlaylistData=(PlaylistName=Playlist_DefaultSolo, PlaylistAccess=(bEnabled=True, bIsDefaultPlaylist=True, bVisibleWhenDisabled=True, bDisplayAsNew=False, CategoryIndex=0, bDisplayAsLimitedTime=False, DisplayPriority=0))
+FrontEndPlaylistData=(PlaylistName=Playlist_DefaultDuo, PlaylistAccess=(bEnabled=True, bIsDefaultPlaylist=True, bVisibleWhenDisabled=True, bDisplayAsNew=False, CategoryIndex=0, bDisplayAsLimitedTime=False, DisplayPriority=1))
+FrontEndPlaylistData=(PlaylistName=Playlist_DefaultSquad, PlaylistAccess=(bEnabled=True, bIsDefaultPlaylist=True, bVisibleWhenDisabled=True, bDisplayAsNew=False, CategoryIndex=0, bDisplayAsLimitedTime=False, DisplayPriority=2))
+FrontEndPlaylistData=(PlaylistName=Playlist_ShowdownAlt_Solo, PlaylistAccess=(bEnabled=True, bIsDefaultPlaylist=False, bVisibleWhenDisabled=True, bDisplayAsNew=False, CategoryIndex=1, bDisplayAsLimitedTime=False, DisplayPriority=0))
+FrontEndPlaylistData=(PlaylistName=Playlist_ShowdownAlt_Duos, PlaylistAccess=(bEnabled=True, bIsDefaultPlaylist=False, bVisibleWhenDisabled=True, bDisplayAsNew=False, CategoryIndex=1, bDisplayAsLimitedTime=False, DisplayPriority=1))
`)}

func GetDefaultRuntime() []byte {return []byte(`
[/Script/FortniteGame.FortRuntimeOptions]
!DisabledFrontendNavigationTabs=ClearArray
;+DisabledFrontendNavigationTabs=(TabName="AthenaChallenges",TabState=EFortRuntimeOptionTabState::Hidden)
;+DisabledFrontendNavigationTabs=(TabName="Showdown",TabState=EFortRuntimeOptionTabState::Hidden)
;+DisabledFrontendNavigationTabs=(TabName="AthenaStore",TabState=EFortRuntimeOptionTabState::Hidden)

[/Script/FortniteGame.FortRuntimeOptions]
bForceBRMode=True
bSkipSubgameSelect=True
bEnableInGameMatchmaking=True

bEnableGlobalChat=true
bDisableGifting=false
bDisableGiftingPC=false
bDisableGiftingPS4=false
bDisableGiftingXB=false
!ExperimentalCohortPercent=ClearArray
+ExperimentalCohortPercent=(CohortPercent=100,ExperimentNum=20)`)}