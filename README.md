# experimental go thing

Just figuring out how golang works.

This is a bot that looks for uses of the word "blockchain" on twitter,
and if it finds a tweet, replies with "@user That's Numberwang!"

The idea is based off this tweet
https://twitter.com/jpwarren/status/968712815656697857

Numberwang is this https://www.youtube.com/watch?v=xmBCh76_qWE


## Details
* If the user has been tweeted at before, they should not be tweeted at again
* RTs instances are not valid(?)


## TODOs
* LoggerWorker
* DB+tweet worker
* Are RT and replies potentially forever tweetable?
* Separate concerns into separate workers
* Ensure no races or multitweets because something is slow
