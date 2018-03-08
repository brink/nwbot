package main

import (
  "fmt"
  "os"
  "regexp"
  "time"
  "math/rand"
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
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

var (
  session, merr = mgo.Dial("localhost")
  client = _setup_twitter_client()
)

func main() {
  if merr != nil {
      panic(merr)
  }
  defer session.Close()
  _setup_logger()

params := &twitter.StreamFilterParams{
    Track: []string{"blockchain"},
    StallWarnings: twitter.Bool(true),
}
stream, err := client.Streams.Filter(params)
  // stream, err := client.Streams.Sample(params)
  log.Info("\n")

  if stream != nil {
    defer stream.Stop()
    demux := twitter.NewSwitchDemux()
    s1 := rand.NewSource(time.Now().UnixNano())
    r1 := rand.New(s1)

    demux.Tweet = func(t *twitter.Tweet) {
      txt := t.Text
      match, _ := regexp.MatchString("( |#)blockchain ", txt)
      if match {
        if needToTweetAt(t.User.ID) {
          if r1.Intn(100) < 30 {
            saveUIDThen(t.User.ID,
                      func() {
                         go sendTweetTo(t.IDStr, t.User.ScreenName)
                      })
          }
        } else {
          table := session.DB("nwbot").C("dupes")
            tim := time.Now()
          err := table.Insert(&User{t.User.ID, tim.Format("2006-01-02T15:04:05.999999-07:00")})
          if err != nil {
            log.Fatal(err)
          }
        }
      }
    }
    for message := range stream.Messages {
      demux.Handle(message)
    }
  }
  if err != nil {
    log.Error(err)
  }
}

type User struct {
        UserID  int64
        Time  string
}

func saveUIDThen(uid int64, afterSave func()) int64 {
  table := session.DB("nwbot").C("users")
  t := time.Now()
  err := table.Insert(&User{uid, t.Format("2006-01-02T15:04:05.999999-07:00")})
  if err != nil {
    log.Fatal(err)
  }
  afterSave()
  return uid
}

func needToTweetAt(uid int64) bool {
  result := User{}
  table := session.DB("nwbot").C("users")
  merr = table.Find(bson.M{"userid": uid}).One(&result)
  if merr != nil {
    // log.Fatal(merr)
  }
  return result.UserID != uid

}

func sendTweetTo(tid, username string) {
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
