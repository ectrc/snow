![1](https://github.com/ectrc/snow/assets/13946988/fc007f07-3878-46e7-b990-668fc3d758d0)

# Snow

Performance first, universal Fortnite private server backend written in Go.

## Features

- **Blazing Fast** Written in Go and build on Fast HTTP, snow is extremely fast and can handle any profile action in milliseconds with its caching system.
- **Profile Changes** Snow, automatically, keeps track of profile changes exactly like Fortnite does, meaning it is one-to-one with the game and never desyncs.
- **Universal Database** Easily add new storage methods that satisfy the `Storage` interface. This means you can use any database you want. _(example of how to do this coming soon)_

## What's next?

- Final Fortnite Operations, this includes: Battle Pass, Friends, XMPP and Gifting.
- Interact with external Services like Amazon S3 Buckets to save player data externally.
- A way to interact with accounts outside of the game. This is mainly for a web app and other services to interact with the backend.

## Version Support

### Supported

- **_Chapter 1 Season 2_** `Fortnite+Release-2.5-CL-3889387-Windows` I started with this build of the game as it requires more work to get working, this means snow can support _most_ versions of the game.
- **_Chapter 1 Season 5_** `Fortnite+Release-5.41-CL-4363240-Windows` This build was used to make sure challenges, variants and lobby backgrounds work.
- **_Chapter 1 Season 8_** `Fortnite+Release-8.51-CL-6165369-Windows` Fixed the invisible player bug caused by invalid account responses. Also fixed the issue with the item shop spamming the api.
- **_Chapter 2 Season 2_** `Fortnite+Release-12.41-CL-12905909-Windows` Item Shop length is correct, also Creative profile stopping login has also been fixed.
- **_Chapter 3 Season 1_** `Fortnite+Release-19.10-CL-Unknown-Windows` This is a very new build of fortnite that introfuces alot of different methods e.g. locker data is now stored as an item. Every MCP action is now fully working and tested. You need to start using easy anticheat otherwise this will not work.

### Broken

- **_Chapter 1 Season 4_** `Fortnite+Release-4.5-CL-4159770-Windows` Does not accept the Access Token for user authentication. I have some ideas why however not planned for a fix.

## How do I use this?

- **[Discord OAuth Setup Guide](oauth.md)** How to setup Discord OAuth for your backend. This enabled the ability to login to the web app with Discord.

It is becoming more and more possible to quickly setup with backend with each update. Currently, this backend is in a very early stage and it is not optimal to setup however there are guides on the discord server by community members which take you through a step by step process.

**I would recommend waiting until the backend is more stable and has more features before using it.**

## Contributing

Contributions are welcome! Please open an issue or pull request if you would like to contribute.
