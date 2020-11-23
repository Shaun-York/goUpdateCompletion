package sqssrv

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

// CompletionsToNetSuite stiff
type CompletionsToNetSuite struct {
    srv *sqs.SQS
}

// GetSrv return sqs
func (c *CompletionsToNetSuite) GetSrv() error {
    region := os.Getenv("AWS_REGION")
    awsSession, sesserr := session.NewSession(&aws.Config{
        Region: aws.String(region)},
    )
    if sesserr != nil {
        return sesserr
    }
    c.srv = sqs.New(awsSession)
    return nil
}
// DelMsg delete message from queue
func (c *CompletionsToNetSuite) DelMsg(delmsg *sqs.DeleteMessageInput) error {
	_, delerr := c.srv.DeleteMessage(delmsg)
    if (delerr != nil) {
         return delerr
    }
    return nil
}
