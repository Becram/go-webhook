# go-webhook
This project is started as hobbyist to automate and have various system related functionality. Currently it performs 2 functionality

1. Sending emails based on the github pull request.
   The way this works it whenever there is pull request is merged a webhook is trigger to route `/github/webhook` with payload containing all the details in the payload. We  parse through the payload to get the version and app details to build a custom email template and send to all receivers
   It uses the `https://github.com/nikoksr/notify` for sending emails and slack alerts

2. Restart Kubernetes pods based on the payload.
   We have regular issue the rabbitmq losing the consumer worker frequently while the workers stalling without any issue.  Hence this services employs a webhook which will get the payload(currently from `alertmanager` prometheus-stack). Once the webhook is triggerd the service restarts the deployment based the on queue. The map of the queue to deployment is configured through  `config.yaml`

## Environment Variables

* `GH_SECRET`   Secret in webhook to authenticate you post in GH
* `SG_API_KEY`  Sengrid API Key
* `SLACK_TOKEN` Token for slack app

## Test TO-DO
