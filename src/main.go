package main

import (
	"github.com/aws/aws-sdk-go/service/eventbridge"
	"github.com/aws/aws-sdk-go/aws"
	// "github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/session"
	// "github.com/aws/aws-sdk-go/service/sns"

	"fmt"
	"os"
	"time"
)

func main() {
    session := session.Must(session.NewSession())
    client := eventbridge.New(session, aws.NewConfig().WithRegion("us-east-1"))
    name := "testFromSdk"
    putRule(client, &name)
    putTargets(client, &name)
    fmt.Println("Done!")
}

func putRule(client *eventbridge.EventBridge, name *string) {
    now := time.Now().UTC().Add(time.Minute * 70)
    scheduleExperession := now.Format("cron(04 15 2 01 ? 2006)")
    params := eventbridge.PutRuleInput{Name: name, ScheduleExpression: &scheduleExperession}
    _, error := client.PutRule(&params)
    if error != nil {
        fmt.Println(error)
        os.Exit(1)
    } else {
        fmt.Println("Event scheduled at " + now.String())
    }
}

func putTargets(client *eventbridge.EventBridge, name *string) {
    input := "{\"Message\": \"Your laundry is done!\"}"
    arn := "arn:aws:sns:us-east-1:050309447832:test"
    id := "putTargetId"
    target := eventbridge.Target{ Input: &input, Arn: &arn, Id: &id }
    targets := []*eventbridge.Target{&target}
    targetParams := eventbridge.PutTargetsInput{Rule: name, Targets: targets }
    _, error := client.PutTargets(&targetParams)
    if error != nil {
        fmt.Println(error)
        os.Exit(1)
    }
}
