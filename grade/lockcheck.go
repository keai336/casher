package grade

import (
	"context"
	"github.com/juju/errors"
	"net/http"
	"time"
)
import lockcheck "github.com/keai336/MediaUnlockTest"

var LockTestDic = map[string]func(client http.Client) lockcheck.Result{
	"Spotify":   lockcheck.Spotify,
	"Bilibili ": lockcheck.BilibiliTW,
	"Youtube":   lockcheck.YoutubeRegion,
	"Netflix":   lockcheck.NetflixRegion,
	"Chatgpt":   lockcheck.ChatGPT,
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
