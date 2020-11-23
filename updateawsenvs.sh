# SQS input queue (payloads with woc internalids)
C_INPUT_QUEUE=""
# URL to GraphQL server
C_GQL_SERVER_URL=""
# Update completions table retries
C_GQL_SERVER_RETRYS=""
# Graphql server secret
C_GQL_SERVER_SECRET=""

ENVS="{\
INPUT_QUEUE=$C_INPUT_QUEUE,\
GQL_SERVER_URL=$C_GQL_SERVER_URL,\
GQL_SERVER_RETRYS=$C_GQL_SERVER_RETRYS,\
GQL_SERVER_SECRET=$C_GQL_SERVER_SECRET\
}"

aws lambda update-function-configuration \
    --function-name  goUpdateCompletion \
    --environment Variables=$ENVS
