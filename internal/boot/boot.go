package boot

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/gogf/gf/v2/frame/g"
)

func urlParentPath(url string) string {
	url = strings.TrimRight(url, "/")
	suburl := strings.Split(url, "/")
	suburl = suburl[0 : len(suburl)-1]
	newurl := strings.Join(suburl, "/") + "/"
	return newurl
}

func emailPrefix(email interface{}) string {
	if email == nil {
		return ""
	}
	if reflect.TypeOf(email).Kind() != reflect.String {
		return ""
	}

	emailStr, ok := email.(string)
	if !ok {
		return ""
	}

	if emailStr == "" {
		return ""
	}

	subemail := strings.Split(emailStr, "@")
	if len(subemail) > 0 {
		return subemail[0]
	}

	return ""
}

func shortName(fullname string) string {
	g.Log().Debug(context.Background(), "shortName function called with fullname:", fullname)

	if fullname == "" {
		return fullname
	}
	fullname_split := strings.Split(fullname, ":")
	short_name := fullname_split[0]
	return short_name
}

func timeDiffNow(timestamp_int int) string {
	timestamp := int64(timestamp_int)
	if timestamp < 1 {
		return "None"
	}
	timeNow := time.Now().Unix()
	timediff := timeNow - timestamp
	if timediff < 60 {
		return fmt.Sprint(timediff) + " secs ago"
	} else if timediff < 3600 {
		return fmt.Sprint(timediff/60) + " mins ago"
	} else if timediff < 259200 {
		return fmt.Sprint(timediff/3600) + " hours ago"
	} else if timediff > 0 {
		return fmt.Sprint(timediff/86400) + " days ago"
	}
	return fmt.Sprint(timestamp)
}

func timeToStr(timestamp_int int) string {
	timestamp := int64(timestamp_int)
	tm := time.Unix(timestamp, 0)
	timestamp_str := fmt.Sprint(tm.Format("2006-01-02 15:04:05 Mon -0700"))
	return fmt.Sprint(timestamp_str)
}

func init() {
	g.Log().Debug(context.Background(), "Initializing template functions...")

	g.View().BindFunc("urlParentPath", urlParentPath)
	g.View().BindFunc("emailPrefix", emailPrefix)
	g.View().BindFunc("shortName", shortName)
	g.View().BindFunc("timeDiffNow", timeDiffNow)
	g.View().BindFunc("timeToStr", timeToStr)

	g.Log().Debug(context.Background(), "Template functions initialized")
}

func InitView() {
	view := g.View()
	view.BindFunc("urlParentPath", urlParentPath)
	view.BindFunc("emailPrefix", emailPrefix)
	view.BindFunc("shortName", shortName)
	view.BindFunc("timeDiffNow", timeDiffNow)
	view.BindFunc("timeToStr", timeToStr)
}
