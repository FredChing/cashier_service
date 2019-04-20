package lib

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"regexp"
	"time"
)

var numre *regexp.Regexp

func init() {
	numre, _ = regexp.Compile("[^0-9]")
}

func GrandNum(l int) string {
	h := md5.New()
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	fmt.Fprintf(h, "%d.%d.%f", r.Intn(100), r.Int63(), r.Float64())
	str := hex.EncodeToString(h.Sum(nil))
	str = numre.ReplaceAllString(str, "")

	n := len(str)
	if n < l {
		str = str + GrandNum(l-n)
	} else {
		str = str[0:l]
	}

	return str
}

func MD5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}