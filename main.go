package main

import (
    "github.com/YoSmudge/anaconda"
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
    pattern, err := regexp.Compile(`update_name\s+[^\s]{1,20}(\s|$)`)
    if err != nil {
        return err
    }
    prefix, err := regexp.Compile(`update_name\s+`)
    if err != nil {
        return err
    }

    /// Is update_name Query?
    result := pattern.FindString(tweet.Text)
    if len(result) == 0 {
        return nil
    }
    prefix_result := prefix.FindString(tweet.Text)

    /// exract updated_name
    updated_name := result[len(prefix_result):]

    /// update_name
    v := url.Values{}
    v.Set("name", updated_name)
    _, err = api.PostAccountUpdateProfile(v)
    if err != nil {
        return err
    }
    log.Println("update_name to ", updated_name, " by ", tweet.User.Id, "'s tweet ", tweet.Text)

    /// post done message
    v = url.Values{}
    _, err = api.PostTweet(strings.Replace(done_msg, "{updated}", updated_name, -1), v)
    if err != nil {
        return err
    }


    return nil
}

func main() {
    if len(os.Args) < 7 {
        log.Println("6 arguments required. Consumer key, Consumer secret, Access token, Access token secret, @name, done_mesg")
        return
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
            go func() {
                err := UpdateName(api, tweet, myname, done_mesg)
                if err != nil {
                    log.Println(err)
                }
            }()
        }
    }
}
