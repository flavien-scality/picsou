package main

import (
	"github.com/scality/Picsou/pkg/stats"
	"github.com/aws/aws-sdk-go/aws/session"
)

func main() {
	sess, err := session.NewSession()
	if err != nil {
		panic(err)
	}
	_ = stats.New(sess, stats.Regions)
}
