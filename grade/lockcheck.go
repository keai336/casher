package grade

import (
	"context"
	"github.com/juju/errors"
	"net/http"
	"time"
)
import lockcheck "github.com/keai336/MediaUnlockTest"

type CheckItem struct {
	Keyword string
	Check   func(client http.Client) lockcheck.Result
}

var LockTestDic = map[string]CheckItem{
	"Spotify":   {Keyword: "Spotify", Check: lockcheck.Spotify},
	"Bilibili ": {Keyword: "Bilibili", Check: lockcheck.BilibiliTW},
	"Youtube":   {Keyword: "Youtube", Check: lockcheck.YoutubeRegion},
	"Netflix":   {Keyword: "Netflix", Check: lockcheck.NetflixRegion},
	"Chatgpt":   {Keyword: "openai", Check: lockcheck.ChatGPT},
}

var Checktimedout int = 2
var ErrResult = lockcheck.Result{
	Status: -1,
	Region: "",
	Info:   "timeout",
	Err:    errors.New("timeout"),
}

func OneLockTest(locktest func(client http.Client) lockcheck.Result) lockcheck.Result {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(Checktimedout)*time.Second)
	defer cancel()

	// Make the call within the context
	rsChan := make(chan lockcheck.Result)
	go func() {
		rs := locktest(lockcheck.AutoHttpClient)
		rsChan <- rs
	}()

	select {
	case rs := <-rsChan:
		return rs
	case <-ctx.Done():
		return ErrResult
	}
}
