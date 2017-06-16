# unplugg-of-warcraft
Case study: trying to predict the prices of items in World of Warcraft's Auction House system, using [unplugg][unplugg-api].

This is a simple scrapper for [World of Warcraft Auction House API][wow-api] that fetches and dumps the data every 30 mins.

# setup

First, you'll need [golang](https://golang.org/doc/install) installed and configured.

You'll need an apikey, and to set it, along with the realm where you want data from, in the file `.env`.

Then, just build and run.
AH dump files will be available on the same dir as the executable.


[wow-api]: https://dev.battle.net/io-docs
[unplugg-api]: https://github.com/whitesmith/unplugg-api
