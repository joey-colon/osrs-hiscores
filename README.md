# OSRS Hiscores
A light wrapper in Golang to interact with Old School RuneScape Hiscores

## Getting started
To create an instance, start by importing the package:
`go get https://github.com/joey-colon/osrs-hiscores`

and import it as the following:
`import "https://github.com/joey-colon/osrs-hiscores"`

Next, we can spin up a hiscores via:
```
h, err := hiscores.NewHiscores()
```

## Supported methods:
* hiscores.GetPlayer(rsn string) (*Player, error)
* hiscores.GetPlayerSkillLevel(rsn string, skill string) (int64, error)
* hiscores.GetPlayerSkillXp(rsn string, skill string) (int64, error)
* hiscores.GetPlayerSkillRank(rsn string, skill string) (int64, error)
