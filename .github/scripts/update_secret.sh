#!/bin/bash

SECRET_VALUE=$(aws secretsmanager get-secret-value --secret-id /vaccine-bot/twitter-secrets --query "SecretString" --output text | jq -S .)
set +e

if [ "$SECRET_VALUE" == "" ]; then
    echo "Unable to describe secret, this may be because it does not exist. Attempting to create it..."
    aws secretsmanager create-secret --name /vaccine-bot/twitter-secrets --secret-string "${GITHUB_VALUE}"
    exit 0
fi


GITHUB_VALUE=$(echo ${TWITTER_SECRET} | jq -S .)

if [ "$SECRET_VALUE" != "$GITHUB_VALUE" ]; then
    echo "GitHub secret is different to AWS secret, updating..."
    aws secretsmanager update-secret --secret-id /vaccine-bot/twitter-secrets --secret-string "${GITHUB_VALUE}"
fi
