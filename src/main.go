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
    "encoding/json"
)

func main() {
    session := session.Must(session.NewSession())
    client := eventbridge.New(session, aws.NewConfig().WithRegion("us-east-1"))

    now := time.Now().Add(time.Minute) // change to directly call the sns endpoint
    scheduleMessage(client, "testFromMessage", "The timer started.", &now)
    alermTime := now.Add(time.Minute * 70)
    scheduleMessage(client, "testFromSdk", "Your laundry is done!", &alermTime)

    fmt.Println("Done!")
}

func scheduleMessage(client *eventbridge.EventBridge, name string, message string, messageTime *time.Time) {
    bytes, err := json.Marshal(map[string]string{"message": message})
    if err != nil {
        panic(err)
    }
    putRule(client, name, messageTime)
    putTargets(client, name, string(bytes))
}

func putRule(client *eventbridge.EventBridge, name string, messageTime *time.Time) {
    scheduleExperession := messageTime.UTC().Format("cron(04 15 2 01 ? 2006)")
    params := eventbridge.PutRuleInput{Name: &name, ScheduleExpression: &scheduleExperession}
    _, error := client.PutRule(&params)
    if error != nil {
        fmt.Println(error)
        os.Exit(1)
    } else {
        fmt.Println("Event scheduled at " + messageTime.String())
    }
}

func putTargets(client *eventbridge.EventBridge, name string, message string) {
    arn := "arn:aws:sns:us-east-1:050309447832:test"
    id := "putTargetId"
    target := eventbridge.Target{ Input: &message, Arn: &arn, Id: &id }
    targets := []*eventbridge.Target{&target}
    targetParams := eventbridge.PutTargetsInput{Rule: &name, Targets: targets }
    _, error := client.PutTargets(&targetParams)
    if error != nil {
        fmt.Println(error)
        os.Exit(1)
    }
}
