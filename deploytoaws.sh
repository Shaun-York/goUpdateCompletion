# aws-cli needs to be install and authenticated 
# AWS Account number
AWS_ACCOUNT=""
# AWS role with proper privileges 
AWS_ROLE=""
GOARCH="amd64"
GOOS="linux"

GOARCH=amd64 GOOS=linux go build .

zip -r goUpdateCompletion.zip ./goUpdateCompletion

aws lambda create-function \
    --function-name goUpdateCompletion \
    --runtime go1.x \
    --zip-file fileb://goUpdateCompletion.zip \
    --handler goUpdateCompletion \
    --role arn:aws:iam::$AWS_ACCOUNT:role/service-role/$AWS_ROLE