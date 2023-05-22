# go-webhook

## Environment Variables

* `GH_SECRET`   Secret in webhook to authenticate you post in GH
* `SG_API_KEY`  Sengrid API Key
* `SG_FROM` Email to be be displayed for reciever
* `SG_TO_LIST` comma separated list of emails
* `SG_EMAIL_TMPL_FILE`    Location of the file email template file  `/work/static/email.pr.tpl.html`
* `SG_FROM_NAME` Email from name
* `PR_PREFIX` PR title to look for example `Bump arthur`
* `ALERT_SERVICE_LIST` lets of pr labels to be alerted "arthur,apollo"
* `POSTGRESQL_PASSWORD` "admin"
* `POSTGRESQL_DATABASE` "go-webhook"
* `POSTGRESQL_USERNAME` "admin"
* `POSTGRESQL_ADDRESS` "127.0.0.1:5432"
