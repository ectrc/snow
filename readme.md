![1](https://github.com/ectrc/snow/assets/13946988/fc007f07-3878-46e7-b990-668fc3d758d0)

# Snow

Performance first, universal Fortnite private server backend written in Go.

## Features

- **Blazingly Fast** Written in Go and build on Fast HTTP, it is extremely fast and can handle any profile action in milliseconds with its caching system.
- **Profile Changes** Automatically keeps track of profile changes exactly so any external changes are displayed in-game on the next action.
- **Universal Database** It is possible to add new database types to satisfy your needs. Currently, it only supports `postgresql`.

## What's next?

- Friends, Gifting, XMPP, Matchmaker and Battle Pass support.
- Interact with external Services like Amazon S3 Buckets to save player data externally.
- A way to interact with accounts outside of the game. This is mainly for a web app and other services to interact with the backend.

## Version Support

### Supported

- **_Chapter 1 Season 2_** `Fortnite+Release-2.5-CL-3889387-Windows`
- **_Chapter 1 Season 5_** `Fortnite+Release-5.41-CL-4363240-Windows`
- **_Chapter 1 Season 8_** `Fortnite+Release-8.51-CL-6165369-Windows`
- **_Chapter 2 Season 2_** `Fortnite+Release-12.41-CL-12905909-Windows`
- **_Chapter 3 Season 1_** `Fortnite+Release-19.10-CL-Unknown-Windows`

### Not Supported

- **_Chapter 1 Season 4_** `Fortnite+Release-4.5-CL-4159770-Windows` I cannot get JWT Tokens to correcly work. I need to supplt a KID for the JWT Token to work however I cannot find a way to get the KID from the game. If you know how to get the KID from the game please open an issue or pull request.

## Support

- **[Github Wiki](https://github.com/ectrc/snow/wiki)** View all of the setup guides and usefull information on how to setup the backend.
- **[Discord Server](discord.gg/kBefMZA4Qp)** Get help from community members on anything else!

## Contributing

Contributions are welcome! Please open an issue or pull request if you would like to contribute. Make sure to follow the same format (2 space indents) and style!

### Commits

Keep commits on a per feature level e.g. do not commit 17 files at once under the name `add`, rather commit every add or change with the format below so that it is easy to track and understand any commits to the repository.#

- **Feature** `feat: commit message here`
- **Refactor** `refact: commit message here`