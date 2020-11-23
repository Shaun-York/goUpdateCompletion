package main

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"goUpdateCompletion/completion"
	"goUpdateCompletion/gql"
	"goUpdateCompletion/sqssrv"
)

// Handle - update completion with workordercompletion id
func Handle(ctx context.Context, event events.SQSEvent) (string, error) {
	var failed error
	var msgID string 

	for _, v := range event.Records {
		msgID = v.MessageId
		compl := &completion.Completion {
			ReceiptHandle: v.ReceiptHandle,
			MessageBody: v.Body,
		}

		marerr := json.Unmarshal([]byte(v.Body), &compl)
		if marerr != nil {
			failed = marerr
			break
		}
		
		gqlerr := gql.UpdateCompletion(compl)
		if gqlerr != nil {
			failed = gqlerr 
			break
		}

		sqs := &sqssrv.CompletionsToNetSuite{}

		sqserr := sqs.GetSrv()
		if (sqserr != nil) {
			failed = sqserr
			break
		}

		delmsg := compl.SqsDelMsg()
		delerr := sqs.DelMsg(delmsg)

		if delerr != nil {
			failed = delerr
			break
		}
	}
	return msgID, failed
}

func main() {
	lambda.Start(Handle)
}