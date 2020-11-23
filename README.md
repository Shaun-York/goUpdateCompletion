# goUpdateCompletion AWS Lambda 3 of 3
Consume sqs messages. Update completions table with work order completion internalid. rm completed

AWSAPI Gateway -> SQS -> [goUpdateTaskQty](https://github.com/Shaun-York/goUpdateTaskQty) -> SQS -> [goUpdateNetSuite](https://github.com/Shaun-York/goUpdateNetSuite) -> SQS -> [goUpdateCompletion](https://github.com/Shaun-York/goUpdateCompletion)
