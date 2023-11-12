# Vaccine Bot

This project is no longer deployed because some dipshit billionaire bought twitter and fucked around with the APIs and I couldn't be bothered fixing it.

Twitter bot: [@immunotron](https://twitter.com/immunotron)

A simple twitter bot written in Go that will tweet about boosting your immune system by getting vaccinated against diseases that can be prevented by immunisation.  

Uses AWS SAM CLI to package and deploy the code as a Lambda function that is invoked daily by a CloudWatch event. Runs unit tests and the `sam build` and `sam deploy` process via GitHub actions.