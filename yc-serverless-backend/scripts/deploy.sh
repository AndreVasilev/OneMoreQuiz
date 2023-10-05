#!/bin/sh

question_function_id="$QUESTION_FUNCTION_ID"
question_function_name="$QUESTION_FUNCTION_NAME"
question_entrypoint=main.QuestionHandler

user_function_id="$USER_FUNCTION_ID"
user_function_name="$USER_FUNCTION_NAME"
user_entrypoint=main.UserHandler

timestamp=$(date +%s)
archive=func-$timestamp.zip
service_account_id="$SERVICE_ACCOUNT_ID"
bot_token="$BOT_TOKEN"
ydb_connection_url="$YDB_CONNECTION_URL"

echo "Prepare archive"

zip -r $archive ./ -x .*/ -x .* -x scripts/ -x ./scripts/* -x .idea/* -x .idea/.* -x *.zip

echo "Deploy Question function"

yc serverless function version create \
  --function-name=$question_function_name \
  --runtime golang119 \
  --entrypoint $question_entrypoint \
  --memory 128m \
  --execution-timeout 3s \
  --source-path $archive \
  --service-account-id $service_account_id \
  --environment="YDB_CONNECTION_URL=$ydb_connection_url,BOT_TOKEN=$bot_token" \

yc serverless function invoke $question_function_id -d '{"queryStringParameters": {"integration": "raw"}}'

echo "Deploy User function"

yc serverless function version create \
  --function-name=$user_function_name \
  --runtime golang119 \
  --entrypoint $user_entrypoint \
  --memory 128m \
  --execution-timeout 3s \
  --source-path $archive \
  --service-account-id $service_account_id \
  --environment="YDB_CONNECTION_URL=$ydb_connection_url,BOT_TOKEN=$bot_token" \

yc serverless function invoke $user_function_id -d '{"queryStringParameters": {"integration": "raw"}}'
