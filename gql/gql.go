package gql

import (
	"context"
	"errors"
	"log"
	"math"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"goUpdateCompletion/completion"
	"goUpdateCompletion/gqldocs"

	"github.com/machinebox/graphql"
)

//HTTPClientSettings HTTPClientSettings
type HTTPClientSettings struct {
    Connect          time.Duration
    ConnKeepAlive    time.Duration
    ExpectContinue   time.Duration
    IdleConn         time.Duration
    MaxAllIdleConns  int
    MaxHostIdleConns int
    ResponseHeader   time.Duration
    TLSHandshake     time.Duration
}
//SetupClient SetupClient
func SetupClient(httpSettings *HTTPClientSettings) *http.Client {
    tr := &http.Transport{
        ResponseHeaderTimeout: httpSettings.ResponseHeader,
        Proxy:                 http.ProxyFromEnvironment,
        DialContext: (&net.Dialer{
            KeepAlive: httpSettings.ConnKeepAlive,
            DualStack: true,
            Timeout:   httpSettings.Connect,
        }).DialContext,
        MaxIdleConns:          httpSettings.MaxAllIdleConns,
        IdleConnTimeout:       httpSettings.IdleConn,
        TLSHandshakeTimeout:   httpSettings.TLSHandshake,
        MaxIdleConnsPerHost:   httpSettings.MaxHostIdleConns,
        ExpectContinueTimeout: httpSettings.ExpectContinue,
    }

    return &http.Client{
        Transport: tr,
    }
}

//RandMs (attempt) => Math.round(((2 ** (attempt - 1)) * 64) + (Math.random() * 100))
func RandMs(n int64) int64 {
    v := float64(n)
    rand.Seed(time.Now().UnixNano())
    x := v - 1
    ms := math.Pow(2, x) * 64 + float64(rand.Intn(100 - 0 + 1) + 1)
    return int64(math.Round(ms))
}

func waitABit(n int64) {
	ms := RandMs(n)
	log.Printf("Waiting for %dms...",ms)
	duration := time.Duration(ms) * time.Millisecond
	time.Sleep(duration)
}

// UpdateCompletion Update Completion row with workordercompletionID
func UpdateCompletion(c *completion.Completion) error {
	var failed error

	gqlclient := graphql.NewClient(os.Getenv("GQL_SERVER_URL"), 
		graphql.WithHTTPClient(
			SetupClient(
				&HTTPClientSettings{
					Connect:          5 * time.Second,
    				ExpectContinue:   1 * time.Second,
    				IdleConn:         90 * time.Second,
    				ConnKeepAlive:    30 * time.Second,
    				MaxAllIdleConns:  100,
    				MaxHostIdleConns: 10,
    				ResponseHeader:   5 * time.Second,
    				TLSHandshake:     5 * time.Second,
				},
			),
		),
	)

	request := graphql.NewRequest(gqldocs.AddWocID)
	request.Var("id", c.ID)
	request.Var("workordercompletion_id", c.WorkorderCompletionID)

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("x-hasura-admin-secret", os.Getenv("GQL_SERVER_SECRET"))

	retries, rterr := strconv.ParseInt(os.Getenv("GQL_SERVER_RETRYS"), 10, 64)
	if rterr != nil {
		return rterr
	}

	var n int64

	for {
		for {
			n = n + 1
			log.Printf("try %d",n)

			if n >= retries {
				failed = errors.New("updateCompletion failed, ran out of trys")
				break
			}

			resp := &gqldocs.UpdatedCompletionByPk{}

			reqerr := gqlclient.Run(context.Background(), request, &resp)

			if reqerr != nil {
				failed = reqerr
				break
			}

			if c.WorkorderCompletionID == resp.Compl.WorkorderCompletionID {
				break
			}

			if c.WorkorderCompletionID != resp.Compl.WorkorderCompletionID {
				failed = errors.New("updating Completions row failed")
				break
			}
	
			if (c.LastCompletion) {
				delresp := &gqldocs.DeleteTaskByPk{}
				taskpk := c.WorkorderID+c.OperationSequence

				rmreq := graphql.NewRequest(gqldocs.DeletedCompleted)
				rmreq.Header.Add("Content-Type", "application/json")
				rmreq.Header.Add("x-hasura-admin-secret", os.Getenv("GQL_SERVER_SECRET"))
				rmreq.Var("internalid", taskpk)

				delcomplerr := gqlclient.Run(context.Background(), rmreq, &delresp)

				if delcomplerr != nil {
					failed = delcomplerr
					break
				}
			}
		}
		break
	}
	return failed
}
	
