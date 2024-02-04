![1](https://github.com/ectrc/snow/assets/13946988/fc007f07-3878-46e7-b990-668fc3d758d0)

# Snow

Performance first, universal Fortnite private server backend written in Go.

## Features

- **Single File** It will embed all of the external files inside of one executable! This allows the backend to be ran anywhere with no setup _(after initial config)_!
- **Blazingly Fast** Written in Go and built upon Fast HTTP, it is extremely fast and can handle any profile action in milliseconds with its caching system.
- **Profile Changes** Automatically keeps track of profile changes exactly so any external changes are displayed in-game on the next action.
- **Universal Database** It is possible to add new database types to satisfy your needs. Currently, it only supports `postgresql`.

## What's next?

- Final niche actions like `RefundMtxPurchase`, `SetAffiliateName` and more.
- Interact with external services like Amazon S3 or Cloudflare R2 to save player data externally.
- Seed Randomization for the Item Shop instead of a random number generator. This will allow for the same item shop to be generated for each day without the need for a database or keeping the server running.
- Integrating a matchmaker with a hoster to smartly put players into games and know when servers become available.
- Refactor the XMPP solution to use [melium/xmpp](https://github.com/mellium/xmpp).

### Missing Features

> Use this as a roadmap to see what is coming next! It is easier to list what is missing than what is included as the server is almost feature complete.

- **User Settings** An Amazon S3 bucket will be required to store user settings.
- **Party System V2** Ccurrently it relies on the automatic XMPP solution which is very hard to keep track of.
- **Battle Pass** For support, I will need the Battle Pass Storefront ID for every build. I am yet to think of a solution for this.
- **Event Tracking** This will come from the game server and will send down every event for every player. With this it is possible to track quests, player stats, and more. This will need a game server to be built specifically for Snow. Features reliant on this are **Leaderboards**, **Challenges**, **Stats**, **Pins** and **Career**.
- **Save The World** Really?

## Supported MCP Actions

> These are request made from Fortnite to the server to perform actions on the profile.

`QueryProfile`, `ClientQuestLogin`, `MarkItemSeen`, `SetItemFavoriteStatusBatch`, `EquipBattleRoyaleCustomization`, `SetBattleRoyaleBanner`, `SetCosmeticLockerSlot`, `SetCosmeticLockerBanner`, `SetCosmeticLockerName`, `CopyCosmeticLoadout`, `DeleteCosmeticLoadout`, `PurchaseCatalogEntry`, `GiftCatalogEntry`, `RemoveGiftBox`

## Support

- **[Discord Server](https://discord.gg/kBefMZA4Qp)** Get help from community members!

## Contributing

Contributions are welcome! Please open an issue or pull request if you would like to contribute. Make sure to follow the same format (2 space indents) and style! Make sure to keep commits concise and readable e.g. do not change formating to mess up code review!
