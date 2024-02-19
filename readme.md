![1](https://github.com/ectrc/snow/assets/13946988/fc007f07-3878-46e7-b990-668fc3d758d0)

# Snow

> Performance first, feature-rich universal Fortnite private server backend written in Go.

## Overview

- **Single File** It will embed all of the external files inside of one executable! This allows the backend to be ran anywhere with no setup _(after initial config)_!
- **Blazingly Fast** Written in Go and built upon Fast HTTP, it is extremely fast and can handle any profile action in milliseconds with its caching system.
- **Automatic Profile Changes** Automatically keeps track of profile changes exactly so any external changes are displayed in-game on the next action.
- **Universal Database** It is possible to add new database types to satisfy your needs. Currently, it only supports `postgresql`.

## What's up next?

- Purchasing the **Battle Pass**. This will require the Battle Pass Storefront ID for every build. I am yet to think of a solution for this.
- Interaction with a Game Server to handle **Event Tracking** for player statistics and challenges. This will be a very large task as a new specialised game server will need to be created.
- After the game server addition, a **Matchmaking System** will be added to match players together for a game. It will use a bin packing algorithm to ensure that games are filled as much as possible.

And once battle royale is completed ...

- **Save The World**

## Feature List

- **XMPP** For interacting with friends, parties and gifting.
- **Friends** On every build, this will allow for adding, removing and blocking friends.
- **Party System V2** This replaces the legacy xmpp driven party system.
- **Automatic Item Shop** Will automatically update the item shop for the day, for all builds.
- **Gifting** Of any item shop entry to any friend.
- **Locker Loadouts** On seasons 12 onwards, this allows for the saving and loading of multiple locker presets.
- **Item Refunding** Of previous shop purchases, will use a refund ticket if refunded in time.
- **V-Bucks Purchasing** Buy V-Bucks and Starter Packs right from the in-game store!
- **Client Settings Storage** Uses amazon buckets to store client settings.
- **Giftable Bundles** Players can recieve bundles, e.g. Twitch Prime, and gift them to friends.
- **Support A Creator 5%** Use any display name and each purchase will give them 5% of the vbucks spent.
- **Discord Bot** Very useful to control players, their inventory and their settings

## Supported MCP Actions

`QueryProfile`, `ClientQuestLogin`, `MarkItemSeen`, `SetItemFavoriteStatusBatch`, `EquipBattleRoyaleCustomization`, `SetBattleRoyaleBanner`, `SetCosmeticLockerSlot`, `SetCosmeticLockerBanner`, `SetCosmeticLockerName`, `CopyCosmeticLoadout`, `DeleteCosmeticLoadout`, `PurchaseCatalogEntry`, `GiftCatalogEntry`, `RemoveGiftBox`, `RefundMtxPurchase`, `SetAffiliateName`, `SetReceiveGiftsEnabled`

## Support

- **[Discord Server](https://discord.gg/kBefMZA4Qp)** Get help from community members!

## Contributing

Contributions are welcome! Please open an issue or pull request if you would like to contribute. Make sure to follow the same formatting and to keep commits concise and readable e.g. do not change line indents to mess up code review!
