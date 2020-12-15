package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/eventbridge"

	// "github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/session"
	// "github.com/aws/aws-sdk-go/service/sns"

	"encoding/json"
	"fmt"
	"time"
)

type Client struct {
    eventBridgeClient *eventbridge.EventBridge
}

func main() {
    session := session.Must(session.NewSession())
    client := Client{eventbridge.New(session, aws.NewConfig().WithRegion("us-east-1"))}

    now := time.Now().Add(time.Minute) // change to directly call the sns endpoint
    client.ScheduleMessage("testFromMessage", "The timer started.", &now)
    alermTime := now.Add(time.Minute * 70)
    client.ScheduleMessage("testFromSdk", "Your laundry is done!", &alermTime)

    fmt.Println("Done!")
}



func (client *Client) ScheduleMessage(name string, message string, messageTime *time.Time) {
    bytes, err := json.Marshal(map[string]string{"message": message})
    if err != nil {
        panic(err)
    }
    client.PutRule(name, messageTime)
    client.PutTargets(name, string(bytes))
}

func (client *Client) PutRule(name string, messageTime *time.Time) {
    scheduleExperession := messageTime.UTC().Format("cron(04 15 2 01 ? 2006)")
    params := eventbridge.PutRuleInput{Name: &name, ScheduleExpression: &scheduleExperession}
    _, err := client.eventBridgeClient.PutRule(&params)
    if err != nil {
        panic(err)
    } else {
        fmt.Println("Event scheduled at " + messageTime.String())
    }
}

func (client *Client) PutTargets(name string, message string) {
    arn := "arn:aws:sns:us-east-1:050309447832:test"
    id := "putTargetId"
    target := eventbridge.Target{ Input: &message, Arn: &arn, Id: &id }
    targets := []*eventbridge.Target{&target}
    targetParams := eventbridge.PutTargetsInput{Rule: &name, Targets: targets }
    _, err := client.eventBridgeClient.PutTargets(&targetParams)
    if err != nil {
        panic(err)
    }
}
