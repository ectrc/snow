![1](https://github.com/ectrc/snow/assets/13946988/64c3b1ac-d308-4e5d-ad8d-2b7aead29195)

# Snow

Performance first, universal Fortnite backend written in Go.

## Features

- **Blazing Fast** Written in Go, snow is extremely fast and can handle any profile action in milliseconds.
- **Profile Changes** Snow keeps track of profile changes exactly like Fortnite does, meaning it is one-to-one with the game.

## Examples

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

- Implement the HTTP API for the game to communicate with the backend. This is the most important part of the project as it needs to handle thousands of requests per second. _Should I use Fiber?_
- Person Authentication for the game to determine if the person is valid or not. Fortnite uses JWT tokens for this which makes it easy to implement.
- Embed game assets into the backend e.g. Game XP Curve, Quest Data etc. _This would mean a single binary that can be run anywhere without the need of external files._
- Interact with external Buckets to save player data externally.
- A way to interact with persons outside of the game. This is mainly for a web app and other services to interact with the backend.
- Game Server Communication. This would mean a websocket server that communicates with the game servers to send and receive data.

## Contributing

Contributions are welcome! Please open an issue or pull request if you would like to contribute.
