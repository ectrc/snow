![1](https://github.com/ectrc/snow/assets/13946988/fc007f07-3878-46e7-b990-668fc3d758d0)

# Snow

Performance first, universal Fortnite private server backend written in Go.

## Features

- **Single File** It will embed all of the external files inside of one executable! This allows the backend to be ran anywhere with no setup _(after initial config)_!
- **Blazingly Fast** Written in Go and built upon Fast HTTP, it is extremely fast and can handle any profile action in milliseconds with its caching system.
- **Profile Changes** Automatically keeps track of profile changes exactly so any external changes are displayed in-game on the next action.
- **Universal Database** It is possible to add new database types to satisfy your needs. Currently, it only supports `postgresql`.

## What's next?

- More profile actions like `RefundMtxPurchase` and `CopyCosmeticLoadout`.
- Integrating matchmaking with a hoster to smartly put players into games and know when servers become available.
- Interact with external services like Amazon S3 or Cloudflare R2 to save player data externally.
- Refactor the XMPP solution to use [melium/xmpp](https://github.com/mellium/xmpp).

## Version Support

### Supported

- **_Chapter 1 Season 2_** `Fortnite+Release-2.5-CL-3889387-Windows`
- **_Chapter 1 Season 5_** `Fortnite+Release-5.41-CL-4363240-Windows`
- **_Chapter 1 Season 8_** `Fortnite+Release-8.51-CL-6165369-Windows`
- **_Chapter 2 Season 2_** `Fortnite+Release-12.41-CL-12905909-Windows`
- **_Chapter 3 Season 1_** `Fortnite+Release-19.10-CL-Unknown-Windows`

### Not Supported

- **_Chapter 1 Season 4_** `Fortnite+Release-4.5-CL-4159770-Windows` I cannot get JWT Tokens to correctly work. I need to supply a key id header for the JWT Token to work however I cannot find the proper kids. If you know how to get the kid from the game please open an issue or pull request.

## Support

- **[Discord Server](https://discord.gg/kBefMZA4Qp)** Get help from community members!

## Contributing

Contributions are welcome! Please open an issue or pull request if you would like to contribute. Make sure to follow the same format (2 space indents) and style! Make sure to keep commits concise and readable e.g. do not change formating to mess up code review!
