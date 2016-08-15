package main

import (
    "github.com/YoSmudge/anaconda"
    "fmt"
    "os"
    "strings"
    "regexp"
    "net/url"
    "log"
)


func UpdateName(api *anaconda.TwitterApi, tweet anaconda.Tweet, myname string, done_msg string) error {
    /// Get User Info
    self, err := api.GetSelf(nil)
    if err != nil {
        return err
    }

    /// Is reply to me?
    if tweet.InReplyToUserID != self.Id && strings.Index(tweet.Text, "@" + myname) == -1 {
        return nil
    }

    /// Create Pattern
    pattern, err := regexp.Compile(`update_name\s([^\s]|\\.){1,20}(\s|$)`)
    if err != nil {
        return err
    }

    /// Is update_name Query?
    result := pattern.Find([]byte(tweet.Text))
    if result == nil {
        return nil
    }

    /// exract updated_name
    updated_name := string(result[len("update_name "):])

    /// update_name
    v := url.Values{}
    v.Set("name", updated_name)
    _, err = api.PostAccountUpdateProfile(v)
    if err != nil {
        return err
    }

    /// post done message
    v = url.Values{}
    _, err = api.PostTweet(fmt.Sprintf(done_msg, updated_name), v)
    if err != nil {
        return err
    }

    return nil
}

func main() {
    if len(os.Args) < 7 {
        log.Println("6 arguments required. Consumer key, Consumer secret, Access token, Access token secret, @name, done_mesg")
    }
    anaconda.SetConsumerKey(os.Args[1])
    anaconda.SetConsumerSecret(os.Args[2])

    api := anaconda.NewTwitterApi(os.Args[3], os.Args[4])

    myname := os.Args[5]
    done_mesg := os.Args[6]

    v := url.Values{}
    v.Set("track", myname)
    stream := api.UserStream(v)

    for {
        t := <-stream.C
        switch tweet := t.(type) {
        case anaconda.Tweet:
            go func () {
                err := UpdateName(api, tweet, myname, done_mesg)
                if err != nil {
                    log.Println(err)
                }
            }()
        }
    }
}
