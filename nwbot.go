package main

import (
  "fmt"
  "os"

  "github.com/dghubble/go-twitter/twitter"
  "github.com/dghubble/oauth1"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("example")
var (
    // DefaultFormatter is the default formatter used and is only the message.
    DefaultFormatter = logging.MustStringFormatter("%{message}")
    // GlogFormatter mimics the glog format
    GlogFormatter = logging.MustStringFormatter("%{level:.1s}%{time:0102 15:04:05.999999} %{pid} %{shortfile}] %{message}")
)
// Example format string. Everything except the message has a custom color
// which is dependent on the log level. Many fields have a custom output
// formatting too, eg. the time returns the hour down to the milli second.
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)



func main() {
  _setup_logger()
  client := _setup_twitter_client()

//TODO Rewrite to use stream
  search, resp, err := client.Search.Tweets(&twitter.SearchTweetParams{
    Query: "blockchain",
  })
  log.Info("\n")

  if search != nil {
    //TODO save tweet ID to db
    //TODO send tweet if not in db
    //TODO send tweet in worker
    for _, t := range search.Statuses {
      if needToTweetAt(t.User.ID) {
          saveUIDThen(t.User.ID, func() {
                                   sendTweetTo(t.IDStr, t.User.ScreenName, client)
                                 })
      }

    }
  }
  if err != nil {
    log.Debug(resp.Body)
    log.Error(err)
  }
}

func saveUIDThen(uid int64, afterSave func()) int64 {
  uid = uid
  afterSave()
  return uid
}

func needToTweetAt(uid int64) bool {
  return true
}

func sendTweetTo(tid, username string, client *twitter.Client) {

  responseStr := "@" + username + " That's Numberwang!"//, &twitter.StatusUpdateParams{InReplyToStatusID: t.IDStr})
  fmt.Println(responseStr, tid)
  // // tweet, resp, err := client.Statuses.Update(responseStr, &twitter.StatusUpdateParams{InReplyToStatusID: tid})
  //
  // if err != nil {
  //   log.Debug(resp.Body)
  //   log.Error(err)
  // }
  // log.Info(tweet)
}

func _setup_twitter_client() *twitter.Client {
  config := oauth1.NewConfig("ZnLFucEj6ylvkgI5qkXuZLFug", "Efi61amab0I9ZUQEt5g6fiwdpBwFbaIDB8OnIl2LDDvZNGck1S")
  token := oauth1.NewToken("970375160317362176-01dg7jAAjUKPAAIPBa8QSE0zLyth5Jc", "IecJbUhttV1OLmSWyGpaA7V1UgInsfJDa0PfKWtOqBPjE")
  httpClient := config.Client(oauth1.NoContext, token)
  client := twitter.NewClient(httpClient)

  return client;
}

func _setup_logger() {
  backend1 := logging.NewLogBackend(os.Stderr, "", 0)
	backend1Leveled := logging.AddModuleLevel(backend1)
	backend1Formatter := logging.NewBackendFormatter(backend1, format)
  logging.SetBackend(backend1Leveled, backend1Formatter)
}
