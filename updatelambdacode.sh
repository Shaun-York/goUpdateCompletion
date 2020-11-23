# AWS Account number
AWS_ACCOUNT=""
# AWS Role with proper permissions
AWS_ROLE=""
GOARCH="amd64"
GOOS="linux"

GOARCH=amd64 GOOS=linux go build -ldflags="-s -w"

zip -r goUpdateCompletion.zip goUpdateCompletion

aws lambda update-function-code \
    --function-name goUpdateCompletion \
    --zip-file fileb://goUpdateCompletion.zip