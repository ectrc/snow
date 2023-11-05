![1](https://github.com/ectrc/snow/assets/13946988/fc007f07-3878-46e7-b990-668fc3d758d0)

# Snow

Performance first, universal Fortnite private server backend written in Go.

## Features

- **Blazing Fast** Written in Go and build on Fast HTTP, snow is extremely fast and can handle any profile action in milliseconds with its caching system.
- **Profile Changes** Snow, semi-automatically, keeps track of profile changes exactly like Fortnite does, meaning it is one-to-one with the game.
- **Universal Database** Easily add new storage methods that satisfy the `Storage` interface. This means you can use any database you want. _(example of how to do this coming soon)_

## Examples of Person Structures

### Quests

```golang
schedule := person.NewItem("ChallengeBundleSchedule:Paid_1", 1)
user.AthenaProfile.Items.AddItem(schedule)

bundle := person.NewItem("ChallengeBundle:Daily_1", 1)
user.AthenaProfile.Items.AddItem(bundle)

quest := person.NewQuest("Quest:Quest_2", bundle.ID, schedule.ID)
quest.AddObjective("quest_objective_eliminateplayers", 0)
user.AthenaProfile.Quests.AddQuest(quest)

daily := person.NewDailyQuest("Quest:Quest_3")
daily.AddObjective("quest_objective_place_top10", 0)
user.AthenaProfile.Quests.AddQuest(daily)
```

### Profile Changes

```golang
snapshot := user.CommonCoreProfile.Snapshot()
{
  vbucks := user.CommonCoreProfile.Items.GetItemByTemplateID("Currency:MtxPurchased")
  vbucks.Quantity = 200
  vbucks.Favorite = true

  user.CommonCoreProfile.Items.DeleteItem(user.CommonCoreProfile.Items.GetItemByTemplateID("Token:CampaignAccess").ID)
  user.CommonCoreProfile.Items.AddItem(person.NewItem("Token:ReceiveMtxCurrency", 1))
}
user.CommonCoreProfile.Diff(snapshot)
```

## What's next?

- Every endpoint that is used by Fortnite. This includes all MCP Operations, extracting data from the telemetry and even the Party Service (maybe party v2).
- Automatic storefront that uses previous data to generate a storefront. This would use item shops from history to make sure there are no blank spots in the storefront. Also battle pass tiers etc.
- Embed game assets into the backend e.g. Battle Pass, Quest Data etc. _This would mean a single binary that can be run anywhere without the need of external files._
- Interact with external Services like Amazon S3 Buckets to save player data externally.
- A way to interact with accounts outside of the game. This is mainly for a web app and other services to interact with the backend.

## Known Supported Versions

- **_Chapter 1 Season 2_** `Fortnite+Release-2.5-CL-3889387-Windows` I started with this build of the game as it requires more work to get working, this means snow can support _most_ versions of the game.
- **_Chapter 1 Season 5_** `Fortnite+Release-5.41-CL-4363240-Windows` This build was used to make sure challenges, variants and lobby backgrounds work.

## How do I use this?

It is _technically_ possible to clone the repository and host the backend. However, this is in a very early stage and does not have a lot of features. _For example, you cannot create an account externally._ **I would recommend waiting until the backend is more stable and has more features before using it.**

## Contributing

Contributions are welcome! Please open an issue or pull request if you would like to contribute.
