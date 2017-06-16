# unplugg-of-warcraft
Case study: trying to predict the prices of items in World of Warcraft's Auction House system, using [unplugg][unplugg-api].

This is a simple scrapper for [World of Warcraft Auction House API][wow-api] that fetches and dumps the data every 30 mins.

## setup

First, you'll need [golang](https://golang.org/doc/install) installed and configured.

You'll need an apikey, and to set it, along with the realm where you want data from, in the file `.env`.

Then, just build and run.
AH dump files will be available on the same dir as the executable.

## Notes

### AH data

```
{
"realms": [
	{"name":"Grim Batol","slug":"grim-batol"},
	{"name":"Aggra (PortuguÃªs)","slug":"aggra-portugues"}],
"auctions": [
	{"auc":1406686168,"item":32196,"owner":"Xishi","ownerRealm":"Grim Batol","bid":2989999,"buyout":2989999,"quantity":1,"timeLeft":"VERY_LONG","rand":0,"seed":0,"context":0},
	{"auc":1406686171,"item":78268,"owner":"Xishi","ownerRealm":"Grim Batol","bid":182709999,"buyout":182709999,"quantity":1,"timeLeft":"VERY_LONG","rand":0,"seed":0,"context":14},
	{"auc":1406817244,"item":128552,"owner":"Denzan","ownerRealm":"Aggra (PortuguÃªs)","bid":49390499,"buyout":51989999,"quantity":1,"timeLeft":"LONG","rand":0,"seed":0,"context":0},
```
[Example AH data](http://auction-api-eu.worldofwarcraft.com/auction-data/1878bff06a82775ebf6438e312cd2682/auctions.json)

### AH timeLeft
We don't know when an item on the AH was sold, cancelled or the auction simply expired.

We can instead assume that when an auction dissappers, from a previous dump to the new one:
  - if the previous lenght was `Short`, then consider it as `expired` 
  - else is considered as `sold`

We won't try to guess cancelled auctions since that should be a rare case and we have no good way to identify them.

About Auction lengths:

  - Short - Less than 30 minutes.
  - Medium - Between 30 minutes and 2 hours.
  - Long - Between 2 hours and 12 hours.
  - Very Long - Between 12 hours and 48 hours.




[wow-api]: https://dev.battle.net/io-docs
[unplugg-api]: https://github.com/whitesmith/unplugg-api
