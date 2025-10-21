# Changelog
All notable changes to this project will be documented in this file. See [conventional commits](https://www.conventionalcommits.org/) for commit guidelines.

- - -
## mail-server/v0.3.0 - 2025-10-21

- - -

## shared/v0.1.0-rc-47833be - 2025-10-21
#### Bug Fixes
- use amqp to send email to the user instead of sync grpc call - (68cc0fe) - Brandon Guigo
#### Features
- send email to user when he joins the waiting list - (861f34b) - Brandon Guigo

- - -

## shared/v0.0.2-rc-cf87e90 - 2025-10-21
#### Bug Fixes
- golint - (af1f6bc) - Brandon Guigo
- linter - (705f203) - Brandon Guigo
- handle the case of permanent failure when sending an email - (7127640) - Brandon Guigo
- migrate reset password to internal mail-server via grpc - (ce9b636) - Brandon Guigo
- retry only failed mails + send no reply grpc implem + add back missing Dockerfiles - (6abc995) - Brandon Guigo
- setup rpc for mail server to be able to send emails as noreply - (1ee3502) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** bump versions for mail-server@mail-server/v0.2.0 mail@mail/v0.2.0 [skip ci] - (a98b0da) - GitHub Actions

- - -

## mail-server/v0.2.0 - 2025-09-25
#### Features
- upgrade the smtp connection to TLS while sending an email - (447ecf4) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** bump versions for mail-server@mail-server/v0.1.1 [skip ci] - (f111062) - GitHub Actions

- - -

## mail-server/v0.1.1 - 2025-09-25
#### Bug Fixes
- get the send mail id from the message body instead of the message header - (29f16cc) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** bump RC versions for mail-server@mail-server/v0.1.1-rc-fdeb76e [skip ci] - (f7b0a9f) - GitHub Actions
- **(release)** bump versions for auth@auth/v0.11.0 grpc@grpc/v0.2.0 mail-server@mail-server/v0.1.0 mail@mail/v0.1.0 productivity@productivity/v0.11.0 shared@shared/v0.0.1 [skip ci] - (bcf942a) - GitHub Actions

- - -

## mail-server/v0.1.1-rc-fdeb76e - 2025-09-25
#### Bug Fixes
- get the send mail id from the message body instead of the message header - (29f16cc) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** bump versions for auth@auth/v0.11.0 grpc@grpc/v0.2.0 mail-server@mail-server/v0.1.0 mail@mail/v0.1.0 productivity@productivity/v0.11.0 shared@shared/v0.0.1 [skip ci] - (bcf942a) - GitHub Actions

- - -

## mail-server/v0.1.0 - 2025-09-23
#### Bug Fixes
- send email model + controller issue - (b64805f) - Brandon Guigo
- error when sending a message through amqp - (d41b2ed) - Brandon Guigo
- tests and add a script to run all tests - (c0fce5a) - Brandon Guigo
- refactor mail-server - (79e6072) - Brandon Guigo
- move utils into shared library - (3b5b0a6) - Brandon Guigo
- move services into a shared directory - (f9de3a3) - Brandon Guigo
- linter - (ec1ed75) - Brandon Guigo
- move amqp and age encryption utils to their services - (c1e9f09) - Brandon Guigo
- unit tests for mail-server - (5d02b4c) - Brandon Guigo
- linter - (bfb2e19) - Brandon Guigo
- linter - (5b4820d) - Brandon Guigo
- worker only treated one message (double ack) - (d5f6c8b) - Brandon Guigo
- parsing of the mail payload - (cced535) - Brandon Guigo
- retry finally works with rabbitmq expiration - (b895a83) - Brandon Guigo
- worker processing registers message - (66563c2) - Brandon Guigo
- some error with retry publishing - (dd151eb) - Brandon Guigo
- revert: "feat: start of the implementation of the gRPC calls to manage the sending emails" - (8b91950) - Brandon Guigo
- revert: "feat: setup the cron for send email" - (930010e) - Brandon Guigo
- tests - (36b4421) - Brandon Guigo
- completely disable auth on receive email server - (e184b76) - Brandon Guigo
- centralize dockerfiles and allow build with grpc in the monorepo - (120804e) - Brandon Guigo
- remove maizzle dockerfile code for mail-server - (9fb9ee7) - Brandon Guigo
- linter for mail and mail-server - (b52c88f) - Brandon Guigo
- grpc linter + mail in test ci/cd + fix error in smtp server test - (6af5aed) - Brandon Guigo
- receive email without attachement leads to no content disposition - (6e6e90a) - Brandon Guigo
#### Features
- make grpc call work - (9de7e32) - Brandon Guigo
- add grpc client to mail server - (075ac2d) - Brandon Guigo
- send email is working without grpc calls - (bce3a06) - Brandon Guigo
- continue work on sending the emails - (bc6a224) - Brandon Guigo
- setup dkim signing - (aee8834) - Brandon Guigo
- get the mx record and setup the sending loop for each recipients - (7ca3a1c) - Brandon Guigo
- setup the structure of the send methods - (7d26440) - Brandon Guigo
- setup amqp worker to listen to retry worker too - (be4bc66) - Brandon Guigo
- make amqp consumer / producer configuration totally via - (819472d) - Brandon Guigo
- setup the send worker for emails - (2a4aa0b) - Brandon Guigo
- start of the implementation of the gRPC calls to manage the sending emails - (1d434a6) - Brandon Guigo
- setup the cron for send email - (a2879b2) - Brandon Guigo
- add bruno collections and fix errors - (3f3c000) - Brandon Guigo
- upload file and email to storage with transactions for all recipients - (5b7dd6e) - Brandon Guigo
- update grpc to use latest version + configure dev docker compose to use go workspaces + add grpc to get public key - (77f90ca) - Brandon Guigo
- make the smtp server handle anonymous auth mecanism + use emersion packages inside the test script - (162d9f9) - Brandon Guigo
- ack the message when processing is done - (79df841) - Brandon Guigo
- parse the newly added amqp message - (a581449) - Brandon Guigo
- add smtp server tests and mocks - (8c82ef8) - Brandon Guigo
- make the smtp server catch all the infos of the sender before sending amqp - (9b6bbd1) - Brandon Guigo
- configure rspamd in docker - (b80a568) - Brandon Guigo
- parse email attachements in the mail server - (87fe245) - Brandon Guigo
- configure amqp producer in mail server - (d23f581) - Brandon Guigo
- configure amqp consumer and producer - (353f891) - Brandon Guigo
- fix middleware and add mail api - (b2f1de8) - Brandon Guigo
- setup mail-server - (395efe7) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.1 grpc@grpc/v0.2.0-rc.1 mail-server@mail-server/v0.1.0-rc.1 mail@mail/v0.1.0-rc.1 productivity@productivity/v0.11.0-rc.1 shared@shared/v0.0.1-rc.1 [skip ci] - (c3e349d) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.20 grpc@grpc/v0.2.0-rc.20 mail-server@mail-server/v0.1.0-rc.20 mail@mail/v0.1.0-rc.20 productivity@productivity/v0.11.0-rc.20 shared@shared/v0.0.1-rc.11 [skip ci] - (5b4e664) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.19 grpc@grpc/v0.2.0-rc.19 mail-server@mail-server/v0.1.0-rc.19 mail@mail/v0.1.0-rc.19 productivity@productivity/v0.11.0-rc.19 shared@shared/v0.0.1-rc.10 [skip ci] - (b643ac0) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.18 grpc@grpc/v0.2.0-rc.18 mail-server@mail-server/v0.1.0-rc.18 mail@mail/v0.1.0-rc.18 productivity@productivity/v0.11.0-rc.18 shared@shared/v0.0.1-rc.9 [skip ci] - (d941278) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.17 grpc@grpc/v0.2.0-rc.17 mail-server@mail-server/v0.1.0-rc.17 mail@mail/v0.1.0-rc.17 productivity@productivity/v0.11.0-rc.17 shared@shared/v0.0.1-rc.8 [skip ci] - (c2fe363) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.16 grpc@grpc/v0.2.0-rc.16 mail-server@mail-server/v0.1.0-rc.16 mail@mail/v0.1.0-rc.16 productivity@productivity/v0.11.0-rc.16 shared@shared/v0.0.1-rc.7 [skip ci] - (f32a990) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.15 grpc@grpc/v0.2.0-rc.15 mail-server@mail-server/v0.1.0-rc.15 mail@mail/v0.1.0-rc.15 productivity@productivity/v0.11.0-rc.15 shared@shared/v0.0.1-rc.6 [skip ci] - (6bd52c5) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.14 grpc@grpc/v0.2.0-rc.14 mail-server@mail-server/v0.1.0-rc.14 mail@mail/v0.1.0-rc.14 productivity@productivity/v0.11.0-rc.14 shared@shared/v0.0.1-rc.5 [skip ci] - (1a356e4) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.13 grpc@grpc/v0.2.0-rc.13 mail-server@mail-server/v0.1.0-rc.13 mail@mail/v0.1.0-rc.13 productivity@productivity/v0.11.0-rc.13 shared@shared/v0.0.1-rc.4 [skip ci] - (ebdffca) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.12 grpc@grpc/v0.2.0-rc.12 mail-server@mail-server/v0.1.0-rc.12 mail@mail/v0.1.0-rc.12 productivity@productivity/v0.11.0-rc.12 shared@shared/v0.0.1-rc.3 [skip ci] - (cc97a37) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.11 grpc@grpc/v0.2.0-rc.11 mail-server@mail-server/v0.1.0-rc.11 mail@mail/v0.1.0-rc.11 productivity@productivity/v0.11.0-rc.11 shared@shared/v0.0.1-rc.2 [skip ci] - (9392d2d) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.10 grpc@grpc/v0.2.0-rc.10 mail-server@mail-server/v0.1.0-rc.10 mail@mail/v0.1.0-rc.10 productivity@productivity/v0.11.0-rc.10 shared@shared/v0.0.1-rc.1 [skip ci] - (4afa958) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.9 grpc@grpc/v0.2.0-rc.9 mail-server@mail-server/v0.1.0-rc.9 mail@mail/v0.1.0-rc.9 productivity@productivity/v0.11.0-rc.9 [skip ci] - (829f6b8) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.8 grpc@grpc/v0.2.0-rc.8 mail-server@mail-server/v0.1.0-rc.8 mail@mail/v0.1.0-rc.8 productivity@productivity/v0.11.0-rc.8 [skip ci] - (3b1612b) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.7 grpc@grpc/v0.2.0-rc.7 mail-server@mail-server/v0.1.0-rc.7 mail@mail/v0.1.0-rc.7 productivity@productivity/v0.11.0-rc.7 [skip ci] - (de4b8fc) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.6 grpc@grpc/v0.2.0-rc.6 mail-server@mail-server/v0.1.0-rc.6 mail@mail/v0.1.0-rc.6 productivity@productivity/v0.11.0-rc.6 [skip ci] - (4e388ee) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.5 grpc@grpc/v0.2.0-rc.5 mail-server@mail-server/v0.1.0-rc.5 mail@mail/v0.1.0-rc.5 productivity@productivity/v0.11.0-rc.5 [skip ci] - (e3f1da7) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.4 grpc@grpc/v0.2.0-rc.4 mail-server@mail-server/v0.1.0-rc.4 mail@mail/v0.1.0-rc.4 productivity@productivity/v0.11.0-rc.4 [skip ci] - (58c5f8c) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.3 grpc@grpc/v0.2.0-rc.3 mail-server@mail-server/v0.1.0-rc.3 mail@mail/v0.1.0-rc.3 productivity@productivity/v0.11.0-rc.3 [skip ci] - (3391a8b) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.2 grpc@grpc/v0.2.0-rc.2 mail-server@mail-server/v0.1.0-rc.2 mail@mail/v0.1.0-rc.2 productivity@productivity/v0.11.0-rc.2 [skip ci] - (0efc6ff) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.1 grpc@grpc/v0.2.0-rc.1 mail-server@mail-server/v0.1.0-rc.1 mail@mail/v0.1.0-rc.1 productivity@productivity/v0.11.0-rc.1 [skip ci] - (a3f0371) - GitHub Actions
- add mail and mail-server to cog.toml - (68847a1) - Brandon Guigo

- - -

## mail-server/v0.1.0-rc.21 - 2025-09-23
#### Bug Fixes
- send email model + controller issue - (b64805f) - Brandon Guigo
- error when sending a message through amqp - (d41b2ed) - Brandon Guigo
- tests and add a script to run all tests - (c0fce5a) - Brandon Guigo
- refactor mail-server - (79e6072) - Brandon Guigo
- move utils into shared library - (3b5b0a6) - Brandon Guigo
- move services into a shared directory - (f9de3a3) - Brandon Guigo
- linter - (ec1ed75) - Brandon Guigo
- move amqp and age encryption utils to their services - (c1e9f09) - Brandon Guigo
- unit tests for mail-server - (5d02b4c) - Brandon Guigo
- linter - (bfb2e19) - Brandon Guigo
- linter - (5b4820d) - Brandon Guigo
- worker only treated one message (double ack) - (d5f6c8b) - Brandon Guigo
- parsing of the mail payload - (cced535) - Brandon Guigo
- retry finally works with rabbitmq expiration - (b895a83) - Brandon Guigo
- worker processing registers message - (66563c2) - Brandon Guigo
- some error with retry publishing - (dd151eb) - Brandon Guigo
- revert: "feat: start of the implementation of the gRPC calls to manage the sending emails" - (8b91950) - Brandon Guigo
- revert: "feat: setup the cron for send email" - (930010e) - Brandon Guigo
- tests - (36b4421) - Brandon Guigo
- completely disable auth on receive email server - (e184b76) - Brandon Guigo
- centralize dockerfiles and allow build with grpc in the monorepo - (120804e) - Brandon Guigo
- remove maizzle dockerfile code for mail-server - (9fb9ee7) - Brandon Guigo
- linter for mail and mail-server - (b52c88f) - Brandon Guigo
- grpc linter + mail in test ci/cd + fix error in smtp server test - (6af5aed) - Brandon Guigo
- receive email without attachement leads to no content disposition - (6e6e90a) - Brandon Guigo
#### Features
- make grpc call work - (9de7e32) - Brandon Guigo
- add grpc client to mail server - (075ac2d) - Brandon Guigo
- send email is working without grpc calls - (bce3a06) - Brandon Guigo
- continue work on sending the emails - (bc6a224) - Brandon Guigo
- setup dkim signing - (aee8834) - Brandon Guigo
- get the mx record and setup the sending loop for each recipients - (7ca3a1c) - Brandon Guigo
- setup the structure of the send methods - (7d26440) - Brandon Guigo
- setup amqp worker to listen to retry worker too - (be4bc66) - Brandon Guigo
- make amqp consumer / producer configuration totally via - (819472d) - Brandon Guigo
- setup the send worker for emails - (2a4aa0b) - Brandon Guigo
- start of the implementation of the gRPC calls to manage the sending emails - (1d434a6) - Brandon Guigo
- setup the cron for send email - (a2879b2) - Brandon Guigo
- add bruno collections and fix errors - (3f3c000) - Brandon Guigo
- upload file and email to storage with transactions for all recipients - (5b7dd6e) - Brandon Guigo
- update grpc to use latest version + configure dev docker compose to use go workspaces + add grpc to get public key - (77f90ca) - Brandon Guigo
- make the smtp server handle anonymous auth mecanism + use emersion packages inside the test script - (162d9f9) - Brandon Guigo
- ack the message when processing is done - (79df841) - Brandon Guigo
- parse the newly added amqp message - (a581449) - Brandon Guigo
- add smtp server tests and mocks - (8c82ef8) - Brandon Guigo
- make the smtp server catch all the infos of the sender before sending amqp - (9b6bbd1) - Brandon Guigo
- configure rspamd in docker - (b80a568) - Brandon Guigo
- parse email attachements in the mail server - (87fe245) - Brandon Guigo
- configure amqp producer in mail server - (d23f581) - Brandon Guigo
- configure amqp consumer and producer - (353f891) - Brandon Guigo
- fix middleware and add mail api - (b2f1de8) - Brandon Guigo
- setup mail-server - (395efe7) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.20 grpc@grpc/v0.2.0-rc.20 mail-server@mail-server/v0.1.0-rc.20 mail@mail/v0.1.0-rc.20 productivity@productivity/v0.11.0-rc.20 shared@shared/v0.0.1-rc.11 [skip ci] - (5b4e664) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.19 grpc@grpc/v0.2.0-rc.19 mail-server@mail-server/v0.1.0-rc.19 mail@mail/v0.1.0-rc.19 productivity@productivity/v0.11.0-rc.19 shared@shared/v0.0.1-rc.10 [skip ci] - (b643ac0) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.18 grpc@grpc/v0.2.0-rc.18 mail-server@mail-server/v0.1.0-rc.18 mail@mail/v0.1.0-rc.18 productivity@productivity/v0.11.0-rc.18 shared@shared/v0.0.1-rc.9 [skip ci] - (d941278) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.17 grpc@grpc/v0.2.0-rc.17 mail-server@mail-server/v0.1.0-rc.17 mail@mail/v0.1.0-rc.17 productivity@productivity/v0.11.0-rc.17 shared@shared/v0.0.1-rc.8 [skip ci] - (c2fe363) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.16 grpc@grpc/v0.2.0-rc.16 mail-server@mail-server/v0.1.0-rc.16 mail@mail/v0.1.0-rc.16 productivity@productivity/v0.11.0-rc.16 shared@shared/v0.0.1-rc.7 [skip ci] - (f32a990) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.15 grpc@grpc/v0.2.0-rc.15 mail-server@mail-server/v0.1.0-rc.15 mail@mail/v0.1.0-rc.15 productivity@productivity/v0.11.0-rc.15 shared@shared/v0.0.1-rc.6 [skip ci] - (6bd52c5) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.14 grpc@grpc/v0.2.0-rc.14 mail-server@mail-server/v0.1.0-rc.14 mail@mail/v0.1.0-rc.14 productivity@productivity/v0.11.0-rc.14 shared@shared/v0.0.1-rc.5 [skip ci] - (1a356e4) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.13 grpc@grpc/v0.2.0-rc.13 mail-server@mail-server/v0.1.0-rc.13 mail@mail/v0.1.0-rc.13 productivity@productivity/v0.11.0-rc.13 shared@shared/v0.0.1-rc.4 [skip ci] - (ebdffca) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.12 grpc@grpc/v0.2.0-rc.12 mail-server@mail-server/v0.1.0-rc.12 mail@mail/v0.1.0-rc.12 productivity@productivity/v0.11.0-rc.12 shared@shared/v0.0.1-rc.3 [skip ci] - (cc97a37) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.11 grpc@grpc/v0.2.0-rc.11 mail-server@mail-server/v0.1.0-rc.11 mail@mail/v0.1.0-rc.11 productivity@productivity/v0.11.0-rc.11 shared@shared/v0.0.1-rc.2 [skip ci] - (9392d2d) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.10 grpc@grpc/v0.2.0-rc.10 mail-server@mail-server/v0.1.0-rc.10 mail@mail/v0.1.0-rc.10 productivity@productivity/v0.11.0-rc.10 shared@shared/v0.0.1-rc.1 [skip ci] - (4afa958) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.9 grpc@grpc/v0.2.0-rc.9 mail-server@mail-server/v0.1.0-rc.9 mail@mail/v0.1.0-rc.9 productivity@productivity/v0.11.0-rc.9 [skip ci] - (829f6b8) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.8 grpc@grpc/v0.2.0-rc.8 mail-server@mail-server/v0.1.0-rc.8 mail@mail/v0.1.0-rc.8 productivity@productivity/v0.11.0-rc.8 [skip ci] - (3b1612b) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.7 grpc@grpc/v0.2.0-rc.7 mail-server@mail-server/v0.1.0-rc.7 mail@mail/v0.1.0-rc.7 productivity@productivity/v0.11.0-rc.7 [skip ci] - (de4b8fc) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.6 grpc@grpc/v0.2.0-rc.6 mail-server@mail-server/v0.1.0-rc.6 mail@mail/v0.1.0-rc.6 productivity@productivity/v0.11.0-rc.6 [skip ci] - (4e388ee) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.5 grpc@grpc/v0.2.0-rc.5 mail-server@mail-server/v0.1.0-rc.5 mail@mail/v0.1.0-rc.5 productivity@productivity/v0.11.0-rc.5 [skip ci] - (e3f1da7) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.4 grpc@grpc/v0.2.0-rc.4 mail-server@mail-server/v0.1.0-rc.4 mail@mail/v0.1.0-rc.4 productivity@productivity/v0.11.0-rc.4 [skip ci] - (58c5f8c) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.3 grpc@grpc/v0.2.0-rc.3 mail-server@mail-server/v0.1.0-rc.3 mail@mail/v0.1.0-rc.3 productivity@productivity/v0.11.0-rc.3 [skip ci] - (3391a8b) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.2 grpc@grpc/v0.2.0-rc.2 mail-server@mail-server/v0.1.0-rc.2 mail@mail/v0.1.0-rc.2 productivity@productivity/v0.11.0-rc.2 [skip ci] - (0efc6ff) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.1 grpc@grpc/v0.2.0-rc.1 mail-server@mail-server/v0.1.0-rc.1 mail@mail/v0.1.0-rc.1 productivity@productivity/v0.11.0-rc.1 [skip ci] - (a3f0371) - GitHub Actions
- add mail and mail-server to cog.toml - (68847a1) - Brandon Guigo

- - -

## mail-server/v0.1.0-rc.20 - 2025-09-20
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.19 grpc@grpc/v0.2.0-rc.19 mail-server@mail-server/v0.1.0-rc.19 mail@mail/v0.1.0-rc.19 productivity@productivity/v0.11.0-rc.19 shared@shared/v0.0.1-rc.10 [skip ci] - (d988b4d) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.18 grpc@grpc/v0.2.0-rc.18 mail-server@mail-server/v0.1.0-rc.18 mail@mail/v0.1.0-rc.18 productivity@productivity/v0.11.0-rc.18 shared@shared/v0.0.1-rc.9 [skip ci] - (636c138) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.17 grpc@grpc/v0.2.0-rc.17 mail-server@mail-server/v0.1.0-rc.17 mail@mail/v0.1.0-rc.17 productivity@productivity/v0.11.0-rc.17 shared@shared/v0.0.1-rc.8 [skip ci] - (ced0e46) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.16 grpc@grpc/v0.2.0-rc.16 mail-server@mail-server/v0.1.0-rc.16 mail@mail/v0.1.0-rc.16 productivity@productivity/v0.11.0-rc.16 shared@shared/v0.0.1-rc.7 [skip ci] - (b1020c3) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.15 grpc@grpc/v0.2.0-rc.15 mail-server@mail-server/v0.1.0-rc.15 mail@mail/v0.1.0-rc.15 productivity@productivity/v0.11.0-rc.15 shared@shared/v0.0.1-rc.6 [skip ci] - (ed4481d) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.14 grpc@grpc/v0.2.0-rc.14 mail-server@mail-server/v0.1.0-rc.14 mail@mail/v0.1.0-rc.14 productivity@productivity/v0.11.0-rc.14 shared@shared/v0.0.1-rc.5 [skip ci] - (ccad083) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.13 grpc@grpc/v0.2.0-rc.13 mail-server@mail-server/v0.1.0-rc.13 mail@mail/v0.1.0-rc.13 productivity@productivity/v0.11.0-rc.13 shared@shared/v0.0.1-rc.4 [skip ci] - (0e70765) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.12 grpc@grpc/v0.2.0-rc.12 mail-server@mail-server/v0.1.0-rc.12 mail@mail/v0.1.0-rc.12 productivity@productivity/v0.11.0-rc.12 shared@shared/v0.0.1-rc.3 [skip ci] - (68e71b6) - GitHub Actions

- - -

## shared/v0.0.1-rc.3 - 2025-09-20
#### Bug Fixes
- send email model + controller issue - (8b159dc) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.11 grpc@grpc/v0.2.0-rc.11 mail-server@mail-server/v0.1.0-rc.11 mail@mail/v0.1.0-rc.11 productivity@productivity/v0.11.0-rc.11 shared@shared/v0.0.1-rc.2 [skip ci] - (8c59737) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.10 grpc@grpc/v0.2.0-rc.10 mail-server@mail-server/v0.1.0-rc.10 mail@mail/v0.1.0-rc.10 productivity@productivity/v0.11.0-rc.10 shared@shared/v0.0.1-rc.1 [skip ci] - (e44c01a) - GitHub Actions

- - -

## shared/v0.0.1-rc.1 - 2025-09-20
#### Bug Fixes
- error when sending a message through amqp - (e396868) - Brandon Guigo
- tests and add a script to run all tests - (57a16fe) - Brandon Guigo
- refactor mail-server - (fe43346) - Brandon Guigo
- move utils into shared library - (f0f414b) - Brandon Guigo
- move services into a shared directory - (d430c40) - Brandon Guigo
- linter - (4f38427) - Brandon Guigo
- move amqp and age encryption utils to their services - (e56a537) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.9 grpc@grpc/v0.2.0-rc.9 mail-server@mail-server/v0.1.0-rc.9 mail@mail/v0.1.0-rc.9 productivity@productivity/v0.11.0-rc.9 [skip ci] - (5e5d565) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.8 grpc@grpc/v0.2.0-rc.8 mail-server@mail-server/v0.1.0-rc.8 mail@mail/v0.1.0-rc.8 productivity@productivity/v0.11.0-rc.8 [skip ci] - (34a482a) - GitHub Actions

- - -

## productivity/v0.11.0-rc.8 - 2025-09-20
#### Bug Fixes
- unit tests for mail-server - (556eacb) - Brandon Guigo
- linter - (dc7b49b) - Brandon Guigo
- linter - (e472637) - Brandon Guigo
- worker only treated one message (double ack) - (9da89af) - Brandon Guigo
- parsing of the mail payload - (0cba80b) - Brandon Guigo
- retry finally works with rabbitmq expiration - (bc34795) - Brandon Guigo
- worker processing registers message - (0796c81) - Brandon Guigo
- some error with retry publishing - (37f859d) - Brandon Guigo
- revert: "feat: start of the implementation of the gRPC calls to manage the sending emails" - (ba1a891) - Brandon Guigo
- revert: "feat: setup the cron for send email" - (b492a51) - Brandon Guigo
#### Features
- make grpc call work - (c685348) - Brandon Guigo
- add grpc client to mail server - (e51a369) - Brandon Guigo
- send email is working without grpc calls - (fe51726) - Brandon Guigo
- continue work on sending the emails - (176f195) - Brandon Guigo
- setup dkim signing - (0cf38ea) - Brandon Guigo
- get the mx record and setup the sending loop for each recipients - (4a8d8d9) - Brandon Guigo
- setup the structure of the send methods - (d0b53ed) - Brandon Guigo
- setup amqp worker to listen to retry worker too - (9266adf) - Brandon Guigo
- make amqp consumer / producer configuration totally via - (5fcc202) - Brandon Guigo
- setup the send worker for emails - (b324e62) - Brandon Guigo
- start of the implementation of the gRPC calls to manage the sending emails - (ef92c62) - Brandon Guigo
- setup the cron for send email - (2f2b18e) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.7 grpc@grpc/v0.2.0-rc.7 mail-server@mail-server/v0.1.0-rc.7 mail@mail/v0.1.0-rc.7 productivity@productivity/v0.11.0-rc.7 [skip ci] - (24990d2) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.6 grpc@grpc/v0.2.0-rc.6 mail-server@mail-server/v0.1.0-rc.6 mail@mail/v0.1.0-rc.6 productivity@productivity/v0.11.0-rc.6 [skip ci] - (e000019) - GitHub Actions

- - -

## productivity/v0.11.0-rc.6 - 2025-09-20
#### Bug Fixes
- tests - (cd68156) - Brandon Guigo
- completely disable auth on receive email server - (5f71329) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.5 grpc@grpc/v0.2.0-rc.5 mail-server@mail-server/v0.1.0-rc.5 mail@mail/v0.1.0-rc.5 productivity@productivity/v0.11.0-rc.5 [skip ci] - (d55f603) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.4 grpc@grpc/v0.2.0-rc.4 mail-server@mail-server/v0.1.0-rc.4 mail@mail/v0.1.0-rc.4 productivity@productivity/v0.11.0-rc.4 [skip ci] - (f890081) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.3 grpc@grpc/v0.2.0-rc.3 mail-server@mail-server/v0.1.0-rc.3 mail@mail/v0.1.0-rc.3 productivity@productivity/v0.11.0-rc.3 [skip ci] - (ebe0561) - GitHub Actions

- - -

## productivity/v0.11.0-rc.3 - 2025-09-20
#### Bug Fixes
- centralize dockerfiles and allow build with grpc in the monorepo - (95b2dfc) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.2 grpc@grpc/v0.2.0-rc.2 mail-server@mail-server/v0.1.0-rc.2 mail@mail/v0.1.0-rc.2 productivity@productivity/v0.11.0-rc.2 [skip ci] - (03832fc) - GitHub Actions

- - -

## productivity/v0.11.0-rc.2 - 2025-09-20
#### Bug Fixes
- remove maizzle dockerfile code for mail-server - (18cb64e) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.1 grpc@grpc/v0.2.0-rc.1 mail-server@mail-server/v0.1.0-rc.1 mail@mail/v0.1.0-rc.1 productivity@productivity/v0.11.0-rc.1 [skip ci] - (873e29b) - GitHub Actions

- - -

## productivity/v0.11.0-rc.1 - 2025-09-20
#### Bug Fixes
- linter for mail and mail-server - (401774e) - Brandon Guigo
- grpc linter + mail in test ci/cd + fix error in smtp server test - (3b67a14) - Brandon Guigo
- receive email without attachement leads to no content disposition - (f8e0ac5) - Brandon Guigo
#### Features
- add bruno collections and fix errors - (b6b303c) - Brandon Guigo
- upload file and email to storage with transactions for all recipients - (3e69645) - Brandon Guigo
- update grpc to use latest version + configure dev docker compose to use go workspaces + add grpc to get public key - (370934b) - Brandon Guigo
- make the smtp server handle anonymous auth mecanism + use emersion packages inside the test script - (e807664) - Brandon Guigo
- ack the message when processing is done - (080b6fb) - Brandon Guigo
- parse the newly added amqp message - (e38584a) - Brandon Guigo
- add smtp server tests and mocks - (8da1f07) - Brandon Guigo
- make the smtp server catch all the infos of the sender before sending amqp - (4b5fe3c) - Brandon Guigo
- configure rspamd in docker - (32aa09a) - Brandon Guigo
- parse email attachements in the mail server - (5d54950) - Brandon Guigo
- configure amqp producer in mail server - (4cb56c6) - Brandon Guigo
- configure amqp consumer and producer - (325623a) - Brandon Guigo
- fix middleware and add mail api - (b2f1de8) - Brandon Guigo
- setup mail-server - (395efe7) - Brandon Guigo
#### Miscellaneous Chores
- add mail and mail-server to cog.toml - (541e204) - Brandon Guigo

- - -

## mail-server/v0.1.0-rc.19 - 2025-09-16
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.18 grpc@grpc/v0.2.0-rc.18 mail-server@mail-server/v0.1.0-rc.18 mail@mail/v0.1.0-rc.18 productivity@productivity/v0.11.0-rc.18 shared@shared/v0.0.1-rc.9 [skip ci] - (636c138) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.17 grpc@grpc/v0.2.0-rc.17 mail-server@mail-server/v0.1.0-rc.17 mail@mail/v0.1.0-rc.17 productivity@productivity/v0.11.0-rc.17 shared@shared/v0.0.1-rc.8 [skip ci] - (ced0e46) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.16 grpc@grpc/v0.2.0-rc.16 mail-server@mail-server/v0.1.0-rc.16 mail@mail/v0.1.0-rc.16 productivity@productivity/v0.11.0-rc.16 shared@shared/v0.0.1-rc.7 [skip ci] - (b1020c3) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.15 grpc@grpc/v0.2.0-rc.15 mail-server@mail-server/v0.1.0-rc.15 mail@mail/v0.1.0-rc.15 productivity@productivity/v0.11.0-rc.15 shared@shared/v0.0.1-rc.6 [skip ci] - (ed4481d) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.14 grpc@grpc/v0.2.0-rc.14 mail-server@mail-server/v0.1.0-rc.14 mail@mail/v0.1.0-rc.14 productivity@productivity/v0.11.0-rc.14 shared@shared/v0.0.1-rc.5 [skip ci] - (ccad083) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.13 grpc@grpc/v0.2.0-rc.13 mail-server@mail-server/v0.1.0-rc.13 mail@mail/v0.1.0-rc.13 productivity@productivity/v0.11.0-rc.13 shared@shared/v0.0.1-rc.4 [skip ci] - (0e70765) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.12 grpc@grpc/v0.2.0-rc.12 mail-server@mail-server/v0.1.0-rc.12 mail@mail/v0.1.0-rc.12 productivity@productivity/v0.11.0-rc.12 shared@shared/v0.0.1-rc.3 [skip ci] - (68e71b6) - GitHub Actions

- - -

## shared/v0.0.1-rc.3 - 2025-09-16
#### Bug Fixes
- send email model + controller issue - (8b159dc) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.11 grpc@grpc/v0.2.0-rc.11 mail-server@mail-server/v0.1.0-rc.11 mail@mail/v0.1.0-rc.11 productivity@productivity/v0.11.0-rc.11 shared@shared/v0.0.1-rc.2 [skip ci] - (8c59737) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.10 grpc@grpc/v0.2.0-rc.10 mail-server@mail-server/v0.1.0-rc.10 mail@mail/v0.1.0-rc.10 productivity@productivity/v0.11.0-rc.10 shared@shared/v0.0.1-rc.1 [skip ci] - (e44c01a) - GitHub Actions

- - -

## shared/v0.0.1-rc.1 - 2025-09-16
#### Bug Fixes
- error when sending a message through amqp - (e396868) - Brandon Guigo
- tests and add a script to run all tests - (57a16fe) - Brandon Guigo
- refactor mail-server - (fe43346) - Brandon Guigo
- move utils into shared library - (f0f414b) - Brandon Guigo
- move services into a shared directory - (d430c40) - Brandon Guigo
- linter - (4f38427) - Brandon Guigo
- move amqp and age encryption utils to their services - (e56a537) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.9 grpc@grpc/v0.2.0-rc.9 mail-server@mail-server/v0.1.0-rc.9 mail@mail/v0.1.0-rc.9 productivity@productivity/v0.11.0-rc.9 [skip ci] - (5e5d565) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.8 grpc@grpc/v0.2.0-rc.8 mail-server@mail-server/v0.1.0-rc.8 mail@mail/v0.1.0-rc.8 productivity@productivity/v0.11.0-rc.8 [skip ci] - (34a482a) - GitHub Actions

- - -

## productivity/v0.11.0-rc.8 - 2025-09-16
#### Bug Fixes
- unit tests for mail-server - (556eacb) - Brandon Guigo
- linter - (dc7b49b) - Brandon Guigo
- linter - (e472637) - Brandon Guigo
- worker only treated one message (double ack) - (9da89af) - Brandon Guigo
- parsing of the mail payload - (0cba80b) - Brandon Guigo
- retry finally works with rabbitmq expiration - (bc34795) - Brandon Guigo
- worker processing registers message - (0796c81) - Brandon Guigo
- some error with retry publishing - (37f859d) - Brandon Guigo
- revert: "feat: start of the implementation of the gRPC calls to manage the sending emails" - (ba1a891) - Brandon Guigo
- revert: "feat: setup the cron for send email" - (b492a51) - Brandon Guigo
#### Features
- make grpc call work - (c685348) - Brandon Guigo
- add grpc client to mail server - (e51a369) - Brandon Guigo
- send email is working without grpc calls - (fe51726) - Brandon Guigo
- continue work on sending the emails - (176f195) - Brandon Guigo
- setup dkim signing - (0cf38ea) - Brandon Guigo
- get the mx record and setup the sending loop for each recipients - (4a8d8d9) - Brandon Guigo
- setup the structure of the send methods - (d0b53ed) - Brandon Guigo
- setup amqp worker to listen to retry worker too - (9266adf) - Brandon Guigo
- make amqp consumer / producer configuration totally via - (5fcc202) - Brandon Guigo
- setup the send worker for emails - (b324e62) - Brandon Guigo
- start of the implementation of the gRPC calls to manage the sending emails - (ef92c62) - Brandon Guigo
- setup the cron for send email - (2f2b18e) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.7 grpc@grpc/v0.2.0-rc.7 mail-server@mail-server/v0.1.0-rc.7 mail@mail/v0.1.0-rc.7 productivity@productivity/v0.11.0-rc.7 [skip ci] - (24990d2) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.6 grpc@grpc/v0.2.0-rc.6 mail-server@mail-server/v0.1.0-rc.6 mail@mail/v0.1.0-rc.6 productivity@productivity/v0.11.0-rc.6 [skip ci] - (e000019) - GitHub Actions

- - -

## productivity/v0.11.0-rc.6 - 2025-09-16
#### Bug Fixes
- tests - (cd68156) - Brandon Guigo
- completely disable auth on receive email server - (5f71329) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.5 grpc@grpc/v0.2.0-rc.5 mail-server@mail-server/v0.1.0-rc.5 mail@mail/v0.1.0-rc.5 productivity@productivity/v0.11.0-rc.5 [skip ci] - (d55f603) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.4 grpc@grpc/v0.2.0-rc.4 mail-server@mail-server/v0.1.0-rc.4 mail@mail/v0.1.0-rc.4 productivity@productivity/v0.11.0-rc.4 [skip ci] - (f890081) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.3 grpc@grpc/v0.2.0-rc.3 mail-server@mail-server/v0.1.0-rc.3 mail@mail/v0.1.0-rc.3 productivity@productivity/v0.11.0-rc.3 [skip ci] - (ebe0561) - GitHub Actions

- - -

## productivity/v0.11.0-rc.3 - 2025-09-16
#### Bug Fixes
- centralize dockerfiles and allow build with grpc in the monorepo - (95b2dfc) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.2 grpc@grpc/v0.2.0-rc.2 mail-server@mail-server/v0.1.0-rc.2 mail@mail/v0.1.0-rc.2 productivity@productivity/v0.11.0-rc.2 [skip ci] - (03832fc) - GitHub Actions

- - -

## productivity/v0.11.0-rc.2 - 2025-09-16
#### Bug Fixes
- remove maizzle dockerfile code for mail-server - (18cb64e) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.1 grpc@grpc/v0.2.0-rc.1 mail-server@mail-server/v0.1.0-rc.1 mail@mail/v0.1.0-rc.1 productivity@productivity/v0.11.0-rc.1 [skip ci] - (873e29b) - GitHub Actions

- - -

## productivity/v0.11.0-rc.1 - 2025-09-16
#### Bug Fixes
- linter for mail and mail-server - (401774e) - Brandon Guigo
- grpc linter + mail in test ci/cd + fix error in smtp server test - (3b67a14) - Brandon Guigo
- receive email without attachement leads to no content disposition - (f8e0ac5) - Brandon Guigo
#### Features
- add bruno collections and fix errors - (b6b303c) - Brandon Guigo
- upload file and email to storage with transactions for all recipients - (3e69645) - Brandon Guigo
- update grpc to use latest version + configure dev docker compose to use go workspaces + add grpc to get public key - (370934b) - Brandon Guigo
- make the smtp server handle anonymous auth mecanism + use emersion packages inside the test script - (e807664) - Brandon Guigo
- ack the message when processing is done - (080b6fb) - Brandon Guigo
- parse the newly added amqp message - (e38584a) - Brandon Guigo
- add smtp server tests and mocks - (8da1f07) - Brandon Guigo
- make the smtp server catch all the infos of the sender before sending amqp - (4b5fe3c) - Brandon Guigo
- configure rspamd in docker - (32aa09a) - Brandon Guigo
- parse email attachements in the mail server - (5d54950) - Brandon Guigo
- configure amqp producer in mail server - (4cb56c6) - Brandon Guigo
- configure amqp consumer and producer - (325623a) - Brandon Guigo
- fix middleware and add mail api - (b2f1de8) - Brandon Guigo
- setup mail-server - (395efe7) - Brandon Guigo
#### Miscellaneous Chores
- add mail and mail-server to cog.toml - (541e204) - Brandon Guigo

- - -

## mail-server/v0.1.0-rc.18 - 2025-09-15
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.17 grpc@grpc/v0.2.0-rc.17 mail-server@mail-server/v0.1.0-rc.17 mail@mail/v0.1.0-rc.17 productivity@productivity/v0.11.0-rc.17 shared@shared/v0.0.1-rc.8 [skip ci] - (ced0e46) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.16 grpc@grpc/v0.2.0-rc.16 mail-server@mail-server/v0.1.0-rc.16 mail@mail/v0.1.0-rc.16 productivity@productivity/v0.11.0-rc.16 shared@shared/v0.0.1-rc.7 [skip ci] - (b1020c3) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.15 grpc@grpc/v0.2.0-rc.15 mail-server@mail-server/v0.1.0-rc.15 mail@mail/v0.1.0-rc.15 productivity@productivity/v0.11.0-rc.15 shared@shared/v0.0.1-rc.6 [skip ci] - (ed4481d) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.14 grpc@grpc/v0.2.0-rc.14 mail-server@mail-server/v0.1.0-rc.14 mail@mail/v0.1.0-rc.14 productivity@productivity/v0.11.0-rc.14 shared@shared/v0.0.1-rc.5 [skip ci] - (ccad083) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.13 grpc@grpc/v0.2.0-rc.13 mail-server@mail-server/v0.1.0-rc.13 mail@mail/v0.1.0-rc.13 productivity@productivity/v0.11.0-rc.13 shared@shared/v0.0.1-rc.4 [skip ci] - (0e70765) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.12 grpc@grpc/v0.2.0-rc.12 mail-server@mail-server/v0.1.0-rc.12 mail@mail/v0.1.0-rc.12 productivity@productivity/v0.11.0-rc.12 shared@shared/v0.0.1-rc.3 [skip ci] - (68e71b6) - GitHub Actions

- - -

## shared/v0.0.1-rc.3 - 2025-09-15
#### Bug Fixes
- send email model + controller issue - (8b159dc) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.11 grpc@grpc/v0.2.0-rc.11 mail-server@mail-server/v0.1.0-rc.11 mail@mail/v0.1.0-rc.11 productivity@productivity/v0.11.0-rc.11 shared@shared/v0.0.1-rc.2 [skip ci] - (8c59737) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.10 grpc@grpc/v0.2.0-rc.10 mail-server@mail-server/v0.1.0-rc.10 mail@mail/v0.1.0-rc.10 productivity@productivity/v0.11.0-rc.10 shared@shared/v0.0.1-rc.1 [skip ci] - (e44c01a) - GitHub Actions

- - -

## shared/v0.0.1-rc.1 - 2025-09-15
#### Bug Fixes
- error when sending a message through amqp - (e396868) - Brandon Guigo
- tests and add a script to run all tests - (57a16fe) - Brandon Guigo
- refactor mail-server - (fe43346) - Brandon Guigo
- move utils into shared library - (f0f414b) - Brandon Guigo
- move services into a shared directory - (d430c40) - Brandon Guigo
- linter - (4f38427) - Brandon Guigo
- move amqp and age encryption utils to their services - (e56a537) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.9 grpc@grpc/v0.2.0-rc.9 mail-server@mail-server/v0.1.0-rc.9 mail@mail/v0.1.0-rc.9 productivity@productivity/v0.11.0-rc.9 [skip ci] - (5e5d565) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.8 grpc@grpc/v0.2.0-rc.8 mail-server@mail-server/v0.1.0-rc.8 mail@mail/v0.1.0-rc.8 productivity@productivity/v0.11.0-rc.8 [skip ci] - (34a482a) - GitHub Actions

- - -

## productivity/v0.11.0-rc.8 - 2025-09-15
#### Bug Fixes
- unit tests for mail-server - (556eacb) - Brandon Guigo
- linter - (dc7b49b) - Brandon Guigo
- linter - (e472637) - Brandon Guigo
- worker only treated one message (double ack) - (9da89af) - Brandon Guigo
- parsing of the mail payload - (0cba80b) - Brandon Guigo
- retry finally works with rabbitmq expiration - (bc34795) - Brandon Guigo
- worker processing registers message - (0796c81) - Brandon Guigo
- some error with retry publishing - (37f859d) - Brandon Guigo
- revert: "feat: start of the implementation of the gRPC calls to manage the sending emails" - (ba1a891) - Brandon Guigo
- revert: "feat: setup the cron for send email" - (b492a51) - Brandon Guigo
#### Features
- make grpc call work - (c685348) - Brandon Guigo
- add grpc client to mail server - (e51a369) - Brandon Guigo
- send email is working without grpc calls - (fe51726) - Brandon Guigo
- continue work on sending the emails - (176f195) - Brandon Guigo
- setup dkim signing - (0cf38ea) - Brandon Guigo
- get the mx record and setup the sending loop for each recipients - (4a8d8d9) - Brandon Guigo
- setup the structure of the send methods - (d0b53ed) - Brandon Guigo
- setup amqp worker to listen to retry worker too - (9266adf) - Brandon Guigo
- make amqp consumer / producer configuration totally via - (5fcc202) - Brandon Guigo
- setup the send worker for emails - (b324e62) - Brandon Guigo
- start of the implementation of the gRPC calls to manage the sending emails - (ef92c62) - Brandon Guigo
- setup the cron for send email - (2f2b18e) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.7 grpc@grpc/v0.2.0-rc.7 mail-server@mail-server/v0.1.0-rc.7 mail@mail/v0.1.0-rc.7 productivity@productivity/v0.11.0-rc.7 [skip ci] - (24990d2) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.6 grpc@grpc/v0.2.0-rc.6 mail-server@mail-server/v0.1.0-rc.6 mail@mail/v0.1.0-rc.6 productivity@productivity/v0.11.0-rc.6 [skip ci] - (e000019) - GitHub Actions

- - -

## productivity/v0.11.0-rc.6 - 2025-09-15
#### Bug Fixes
- tests - (cd68156) - Brandon Guigo
- completely disable auth on receive email server - (5f71329) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.5 grpc@grpc/v0.2.0-rc.5 mail-server@mail-server/v0.1.0-rc.5 mail@mail/v0.1.0-rc.5 productivity@productivity/v0.11.0-rc.5 [skip ci] - (d55f603) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.4 grpc@grpc/v0.2.0-rc.4 mail-server@mail-server/v0.1.0-rc.4 mail@mail/v0.1.0-rc.4 productivity@productivity/v0.11.0-rc.4 [skip ci] - (f890081) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.3 grpc@grpc/v0.2.0-rc.3 mail-server@mail-server/v0.1.0-rc.3 mail@mail/v0.1.0-rc.3 productivity@productivity/v0.11.0-rc.3 [skip ci] - (ebe0561) - GitHub Actions

- - -

## productivity/v0.11.0-rc.3 - 2025-09-15
#### Bug Fixes
- centralize dockerfiles and allow build with grpc in the monorepo - (95b2dfc) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.2 grpc@grpc/v0.2.0-rc.2 mail-server@mail-server/v0.1.0-rc.2 mail@mail/v0.1.0-rc.2 productivity@productivity/v0.11.0-rc.2 [skip ci] - (03832fc) - GitHub Actions

- - -

## productivity/v0.11.0-rc.2 - 2025-09-15
#### Bug Fixes
- remove maizzle dockerfile code for mail-server - (18cb64e) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.1 grpc@grpc/v0.2.0-rc.1 mail-server@mail-server/v0.1.0-rc.1 mail@mail/v0.1.0-rc.1 productivity@productivity/v0.11.0-rc.1 [skip ci] - (873e29b) - GitHub Actions

- - -

## productivity/v0.11.0-rc.1 - 2025-09-15
#### Bug Fixes
- linter for mail and mail-server - (401774e) - Brandon Guigo
- grpc linter + mail in test ci/cd + fix error in smtp server test - (3b67a14) - Brandon Guigo
- receive email without attachement leads to no content disposition - (f8e0ac5) - Brandon Guigo
#### Features
- add bruno collections and fix errors - (b6b303c) - Brandon Guigo
- upload file and email to storage with transactions for all recipients - (3e69645) - Brandon Guigo
- update grpc to use latest version + configure dev docker compose to use go workspaces + add grpc to get public key - (370934b) - Brandon Guigo
- make the smtp server handle anonymous auth mecanism + use emersion packages inside the test script - (e807664) - Brandon Guigo
- ack the message when processing is done - (080b6fb) - Brandon Guigo
- parse the newly added amqp message - (e38584a) - Brandon Guigo
- add smtp server tests and mocks - (8da1f07) - Brandon Guigo
- make the smtp server catch all the infos of the sender before sending amqp - (4b5fe3c) - Brandon Guigo
- configure rspamd in docker - (32aa09a) - Brandon Guigo
- parse email attachements in the mail server - (5d54950) - Brandon Guigo
- configure amqp producer in mail server - (4cb56c6) - Brandon Guigo
- configure amqp consumer and producer - (325623a) - Brandon Guigo
- fix middleware and add mail api - (b2f1de8) - Brandon Guigo
- setup mail-server - (395efe7) - Brandon Guigo
#### Miscellaneous Chores
- add mail and mail-server to cog.toml - (541e204) - Brandon Guigo

- - -

## mail-server/v0.1.0-rc.17 - 2025-09-14
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.16 grpc@grpc/v0.2.0-rc.16 mail-server@mail-server/v0.1.0-rc.16 mail@mail/v0.1.0-rc.16 productivity@productivity/v0.11.0-rc.16 shared@shared/v0.0.1-rc.7 [skip ci] - (b1020c3) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.15 grpc@grpc/v0.2.0-rc.15 mail-server@mail-server/v0.1.0-rc.15 mail@mail/v0.1.0-rc.15 productivity@productivity/v0.11.0-rc.15 shared@shared/v0.0.1-rc.6 [skip ci] - (ed4481d) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.14 grpc@grpc/v0.2.0-rc.14 mail-server@mail-server/v0.1.0-rc.14 mail@mail/v0.1.0-rc.14 productivity@productivity/v0.11.0-rc.14 shared@shared/v0.0.1-rc.5 [skip ci] - (ccad083) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.13 grpc@grpc/v0.2.0-rc.13 mail-server@mail-server/v0.1.0-rc.13 mail@mail/v0.1.0-rc.13 productivity@productivity/v0.11.0-rc.13 shared@shared/v0.0.1-rc.4 [skip ci] - (0e70765) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.12 grpc@grpc/v0.2.0-rc.12 mail-server@mail-server/v0.1.0-rc.12 mail@mail/v0.1.0-rc.12 productivity@productivity/v0.11.0-rc.12 shared@shared/v0.0.1-rc.3 [skip ci] - (68e71b6) - GitHub Actions

- - -

## shared/v0.0.1-rc.3 - 2025-09-14
#### Bug Fixes
- send email model + controller issue - (8b159dc) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.11 grpc@grpc/v0.2.0-rc.11 mail-server@mail-server/v0.1.0-rc.11 mail@mail/v0.1.0-rc.11 productivity@productivity/v0.11.0-rc.11 shared@shared/v0.0.1-rc.2 [skip ci] - (8c59737) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.10 grpc@grpc/v0.2.0-rc.10 mail-server@mail-server/v0.1.0-rc.10 mail@mail/v0.1.0-rc.10 productivity@productivity/v0.11.0-rc.10 shared@shared/v0.0.1-rc.1 [skip ci] - (e44c01a) - GitHub Actions

- - -

## shared/v0.0.1-rc.1 - 2025-09-14
#### Bug Fixes
- error when sending a message through amqp - (e396868) - Brandon Guigo
- tests and add a script to run all tests - (57a16fe) - Brandon Guigo
- refactor mail-server - (fe43346) - Brandon Guigo
- move utils into shared library - (f0f414b) - Brandon Guigo
- move services into a shared directory - (d430c40) - Brandon Guigo
- linter - (4f38427) - Brandon Guigo
- move amqp and age encryption utils to their services - (e56a537) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.9 grpc@grpc/v0.2.0-rc.9 mail-server@mail-server/v0.1.0-rc.9 mail@mail/v0.1.0-rc.9 productivity@productivity/v0.11.0-rc.9 [skip ci] - (5e5d565) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.8 grpc@grpc/v0.2.0-rc.8 mail-server@mail-server/v0.1.0-rc.8 mail@mail/v0.1.0-rc.8 productivity@productivity/v0.11.0-rc.8 [skip ci] - (34a482a) - GitHub Actions

- - -

## productivity/v0.11.0-rc.8 - 2025-09-14
#### Bug Fixes
- unit tests for mail-server - (556eacb) - Brandon Guigo
- linter - (dc7b49b) - Brandon Guigo
- linter - (e472637) - Brandon Guigo
- worker only treated one message (double ack) - (9da89af) - Brandon Guigo
- parsing of the mail payload - (0cba80b) - Brandon Guigo
- retry finally works with rabbitmq expiration - (bc34795) - Brandon Guigo
- worker processing registers message - (0796c81) - Brandon Guigo
- some error with retry publishing - (37f859d) - Brandon Guigo
- revert: "feat: start of the implementation of the gRPC calls to manage the sending emails" - (ba1a891) - Brandon Guigo
- revert: "feat: setup the cron for send email" - (b492a51) - Brandon Guigo
#### Features
- make grpc call work - (c685348) - Brandon Guigo
- add grpc client to mail server - (e51a369) - Brandon Guigo
- send email is working without grpc calls - (fe51726) - Brandon Guigo
- continue work on sending the emails - (176f195) - Brandon Guigo
- setup dkim signing - (0cf38ea) - Brandon Guigo
- get the mx record and setup the sending loop for each recipients - (4a8d8d9) - Brandon Guigo
- setup the structure of the send methods - (d0b53ed) - Brandon Guigo
- setup amqp worker to listen to retry worker too - (9266adf) - Brandon Guigo
- make amqp consumer / producer configuration totally via - (5fcc202) - Brandon Guigo
- setup the send worker for emails - (b324e62) - Brandon Guigo
- start of the implementation of the gRPC calls to manage the sending emails - (ef92c62) - Brandon Guigo
- setup the cron for send email - (2f2b18e) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.7 grpc@grpc/v0.2.0-rc.7 mail-server@mail-server/v0.1.0-rc.7 mail@mail/v0.1.0-rc.7 productivity@productivity/v0.11.0-rc.7 [skip ci] - (24990d2) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.6 grpc@grpc/v0.2.0-rc.6 mail-server@mail-server/v0.1.0-rc.6 mail@mail/v0.1.0-rc.6 productivity@productivity/v0.11.0-rc.6 [skip ci] - (e000019) - GitHub Actions

- - -

## productivity/v0.11.0-rc.6 - 2025-09-14
#### Bug Fixes
- tests - (cd68156) - Brandon Guigo
- completely disable auth on receive email server - (5f71329) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.5 grpc@grpc/v0.2.0-rc.5 mail-server@mail-server/v0.1.0-rc.5 mail@mail/v0.1.0-rc.5 productivity@productivity/v0.11.0-rc.5 [skip ci] - (d55f603) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.4 grpc@grpc/v0.2.0-rc.4 mail-server@mail-server/v0.1.0-rc.4 mail@mail/v0.1.0-rc.4 productivity@productivity/v0.11.0-rc.4 [skip ci] - (f890081) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.3 grpc@grpc/v0.2.0-rc.3 mail-server@mail-server/v0.1.0-rc.3 mail@mail/v0.1.0-rc.3 productivity@productivity/v0.11.0-rc.3 [skip ci] - (ebe0561) - GitHub Actions

- - -

## productivity/v0.11.0-rc.3 - 2025-09-14
#### Bug Fixes
- centralize dockerfiles and allow build with grpc in the monorepo - (95b2dfc) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.2 grpc@grpc/v0.2.0-rc.2 mail-server@mail-server/v0.1.0-rc.2 mail@mail/v0.1.0-rc.2 productivity@productivity/v0.11.0-rc.2 [skip ci] - (03832fc) - GitHub Actions

- - -

## productivity/v0.11.0-rc.2 - 2025-09-14
#### Bug Fixes
- remove maizzle dockerfile code for mail-server - (18cb64e) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.1 grpc@grpc/v0.2.0-rc.1 mail-server@mail-server/v0.1.0-rc.1 mail@mail/v0.1.0-rc.1 productivity@productivity/v0.11.0-rc.1 [skip ci] - (873e29b) - GitHub Actions

- - -

## productivity/v0.11.0-rc.1 - 2025-09-14
#### Bug Fixes
- linter for mail and mail-server - (401774e) - Brandon Guigo
- grpc linter + mail in test ci/cd + fix error in smtp server test - (3b67a14) - Brandon Guigo
- receive email without attachement leads to no content disposition - (f8e0ac5) - Brandon Guigo
#### Features
- add bruno collections and fix errors - (b6b303c) - Brandon Guigo
- upload file and email to storage with transactions for all recipients - (3e69645) - Brandon Guigo
- update grpc to use latest version + configure dev docker compose to use go workspaces + add grpc to get public key - (370934b) - Brandon Guigo
- make the smtp server handle anonymous auth mecanism + use emersion packages inside the test script - (e807664) - Brandon Guigo
- ack the message when processing is done - (080b6fb) - Brandon Guigo
- parse the newly added amqp message - (e38584a) - Brandon Guigo
- add smtp server tests and mocks - (8da1f07) - Brandon Guigo
- make the smtp server catch all the infos of the sender before sending amqp - (4b5fe3c) - Brandon Guigo
- configure rspamd in docker - (32aa09a) - Brandon Guigo
- parse email attachements in the mail server - (5d54950) - Brandon Guigo
- configure amqp producer in mail server - (4cb56c6) - Brandon Guigo
- configure amqp consumer and producer - (325623a) - Brandon Guigo
- fix middleware and add mail api - (b2f1de8) - Brandon Guigo
- setup mail-server - (395efe7) - Brandon Guigo
#### Miscellaneous Chores
- add mail and mail-server to cog.toml - (541e204) - Brandon Guigo

- - -

## mail-server/v0.1.0-rc.16 - 2025-09-12
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.15 grpc@grpc/v0.2.0-rc.15 mail-server@mail-server/v0.1.0-rc.15 mail@mail/v0.1.0-rc.15 productivity@productivity/v0.11.0-rc.15 shared@shared/v0.0.1-rc.6 [skip ci] - (ed4481d) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.14 grpc@grpc/v0.2.0-rc.14 mail-server@mail-server/v0.1.0-rc.14 mail@mail/v0.1.0-rc.14 productivity@productivity/v0.11.0-rc.14 shared@shared/v0.0.1-rc.5 [skip ci] - (ccad083) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.13 grpc@grpc/v0.2.0-rc.13 mail-server@mail-server/v0.1.0-rc.13 mail@mail/v0.1.0-rc.13 productivity@productivity/v0.11.0-rc.13 shared@shared/v0.0.1-rc.4 [skip ci] - (0e70765) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.12 grpc@grpc/v0.2.0-rc.12 mail-server@mail-server/v0.1.0-rc.12 mail@mail/v0.1.0-rc.12 productivity@productivity/v0.11.0-rc.12 shared@shared/v0.0.1-rc.3 [skip ci] - (68e71b6) - GitHub Actions

- - -

## shared/v0.0.1-rc.3 - 2025-09-12
#### Bug Fixes
- send email model + controller issue - (8b159dc) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.11 grpc@grpc/v0.2.0-rc.11 mail-server@mail-server/v0.1.0-rc.11 mail@mail/v0.1.0-rc.11 productivity@productivity/v0.11.0-rc.11 shared@shared/v0.0.1-rc.2 [skip ci] - (8c59737) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.10 grpc@grpc/v0.2.0-rc.10 mail-server@mail-server/v0.1.0-rc.10 mail@mail/v0.1.0-rc.10 productivity@productivity/v0.11.0-rc.10 shared@shared/v0.0.1-rc.1 [skip ci] - (e44c01a) - GitHub Actions

- - -

## shared/v0.0.1-rc.1 - 2025-09-12
#### Bug Fixes
- error when sending a message through amqp - (e396868) - Brandon Guigo
- tests and add a script to run all tests - (57a16fe) - Brandon Guigo
- refactor mail-server - (fe43346) - Brandon Guigo
- move utils into shared library - (f0f414b) - Brandon Guigo
- move services into a shared directory - (d430c40) - Brandon Guigo
- linter - (4f38427) - Brandon Guigo
- move amqp and age encryption utils to their services - (e56a537) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.9 grpc@grpc/v0.2.0-rc.9 mail-server@mail-server/v0.1.0-rc.9 mail@mail/v0.1.0-rc.9 productivity@productivity/v0.11.0-rc.9 [skip ci] - (5e5d565) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.8 grpc@grpc/v0.2.0-rc.8 mail-server@mail-server/v0.1.0-rc.8 mail@mail/v0.1.0-rc.8 productivity@productivity/v0.11.0-rc.8 [skip ci] - (34a482a) - GitHub Actions

- - -

## productivity/v0.11.0-rc.8 - 2025-09-12
#### Bug Fixes
- unit tests for mail-server - (556eacb) - Brandon Guigo
- linter - (dc7b49b) - Brandon Guigo
- linter - (e472637) - Brandon Guigo
- worker only treated one message (double ack) - (9da89af) - Brandon Guigo
- parsing of the mail payload - (0cba80b) - Brandon Guigo
- retry finally works with rabbitmq expiration - (bc34795) - Brandon Guigo
- worker processing registers message - (0796c81) - Brandon Guigo
- some error with retry publishing - (37f859d) - Brandon Guigo
- revert: "feat: start of the implementation of the gRPC calls to manage the sending emails" - (ba1a891) - Brandon Guigo
- revert: "feat: setup the cron for send email" - (b492a51) - Brandon Guigo
#### Features
- make grpc call work - (c685348) - Brandon Guigo
- add grpc client to mail server - (e51a369) - Brandon Guigo
- send email is working without grpc calls - (fe51726) - Brandon Guigo
- continue work on sending the emails - (176f195) - Brandon Guigo
- setup dkim signing - (0cf38ea) - Brandon Guigo
- get the mx record and setup the sending loop for each recipients - (4a8d8d9) - Brandon Guigo
- setup the structure of the send methods - (d0b53ed) - Brandon Guigo
- setup amqp worker to listen to retry worker too - (9266adf) - Brandon Guigo
- make amqp consumer / producer configuration totally via - (5fcc202) - Brandon Guigo
- setup the send worker for emails - (b324e62) - Brandon Guigo
- start of the implementation of the gRPC calls to manage the sending emails - (ef92c62) - Brandon Guigo
- setup the cron for send email - (2f2b18e) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.7 grpc@grpc/v0.2.0-rc.7 mail-server@mail-server/v0.1.0-rc.7 mail@mail/v0.1.0-rc.7 productivity@productivity/v0.11.0-rc.7 [skip ci] - (24990d2) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.6 grpc@grpc/v0.2.0-rc.6 mail-server@mail-server/v0.1.0-rc.6 mail@mail/v0.1.0-rc.6 productivity@productivity/v0.11.0-rc.6 [skip ci] - (e000019) - GitHub Actions

- - -

## productivity/v0.11.0-rc.6 - 2025-09-12
#### Bug Fixes
- tests - (cd68156) - Brandon Guigo
- completely disable auth on receive email server - (5f71329) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.5 grpc@grpc/v0.2.0-rc.5 mail-server@mail-server/v0.1.0-rc.5 mail@mail/v0.1.0-rc.5 productivity@productivity/v0.11.0-rc.5 [skip ci] - (d55f603) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.4 grpc@grpc/v0.2.0-rc.4 mail-server@mail-server/v0.1.0-rc.4 mail@mail/v0.1.0-rc.4 productivity@productivity/v0.11.0-rc.4 [skip ci] - (f890081) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.3 grpc@grpc/v0.2.0-rc.3 mail-server@mail-server/v0.1.0-rc.3 mail@mail/v0.1.0-rc.3 productivity@productivity/v0.11.0-rc.3 [skip ci] - (ebe0561) - GitHub Actions

- - -

## productivity/v0.11.0-rc.3 - 2025-09-12
#### Bug Fixes
- centralize dockerfiles and allow build with grpc in the monorepo - (95b2dfc) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.2 grpc@grpc/v0.2.0-rc.2 mail-server@mail-server/v0.1.0-rc.2 mail@mail/v0.1.0-rc.2 productivity@productivity/v0.11.0-rc.2 [skip ci] - (03832fc) - GitHub Actions

- - -

## productivity/v0.11.0-rc.2 - 2025-09-12
#### Bug Fixes
- remove maizzle dockerfile code for mail-server - (18cb64e) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.1 grpc@grpc/v0.2.0-rc.1 mail-server@mail-server/v0.1.0-rc.1 mail@mail/v0.1.0-rc.1 productivity@productivity/v0.11.0-rc.1 [skip ci] - (873e29b) - GitHub Actions

- - -

## productivity/v0.11.0-rc.1 - 2025-09-12
#### Bug Fixes
- linter for mail and mail-server - (401774e) - Brandon Guigo
- grpc linter + mail in test ci/cd + fix error in smtp server test - (3b67a14) - Brandon Guigo
- receive email without attachement leads to no content disposition - (f8e0ac5) - Brandon Guigo
#### Features
- add bruno collections and fix errors - (b6b303c) - Brandon Guigo
- upload file and email to storage with transactions for all recipients - (3e69645) - Brandon Guigo
- update grpc to use latest version + configure dev docker compose to use go workspaces + add grpc to get public key - (370934b) - Brandon Guigo
- make the smtp server handle anonymous auth mecanism + use emersion packages inside the test script - (e807664) - Brandon Guigo
- ack the message when processing is done - (080b6fb) - Brandon Guigo
- parse the newly added amqp message - (e38584a) - Brandon Guigo
- add smtp server tests and mocks - (8da1f07) - Brandon Guigo
- make the smtp server catch all the infos of the sender before sending amqp - (4b5fe3c) - Brandon Guigo
- configure rspamd in docker - (32aa09a) - Brandon Guigo
- parse email attachements in the mail server - (5d54950) - Brandon Guigo
- configure amqp producer in mail server - (4cb56c6) - Brandon Guigo
- configure amqp consumer and producer - (325623a) - Brandon Guigo
- fix middleware and add mail api - (b2f1de8) - Brandon Guigo
- setup mail-server - (395efe7) - Brandon Guigo
#### Miscellaneous Chores
- add mail and mail-server to cog.toml - (541e204) - Brandon Guigo

- - -

## mail-server/v0.1.0-rc.15 - 2025-09-12
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.14 grpc@grpc/v0.2.0-rc.14 mail-server@mail-server/v0.1.0-rc.14 mail@mail/v0.1.0-rc.14 productivity@productivity/v0.11.0-rc.14 shared@shared/v0.0.1-rc.5 [skip ci] - (ccad083) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.13 grpc@grpc/v0.2.0-rc.13 mail-server@mail-server/v0.1.0-rc.13 mail@mail/v0.1.0-rc.13 productivity@productivity/v0.11.0-rc.13 shared@shared/v0.0.1-rc.4 [skip ci] - (0e70765) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.12 grpc@grpc/v0.2.0-rc.12 mail-server@mail-server/v0.1.0-rc.12 mail@mail/v0.1.0-rc.12 productivity@productivity/v0.11.0-rc.12 shared@shared/v0.0.1-rc.3 [skip ci] - (68e71b6) - GitHub Actions

- - -

## shared/v0.0.1-rc.3 - 2025-09-12
#### Bug Fixes
- send email model + controller issue - (8b159dc) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.11 grpc@grpc/v0.2.0-rc.11 mail-server@mail-server/v0.1.0-rc.11 mail@mail/v0.1.0-rc.11 productivity@productivity/v0.11.0-rc.11 shared@shared/v0.0.1-rc.2 [skip ci] - (8c59737) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.10 grpc@grpc/v0.2.0-rc.10 mail-server@mail-server/v0.1.0-rc.10 mail@mail/v0.1.0-rc.10 productivity@productivity/v0.11.0-rc.10 shared@shared/v0.0.1-rc.1 [skip ci] - (e44c01a) - GitHub Actions

- - -

## shared/v0.0.1-rc.1 - 2025-09-12
#### Bug Fixes
- error when sending a message through amqp - (e396868) - Brandon Guigo
- tests and add a script to run all tests - (57a16fe) - Brandon Guigo
- refactor mail-server - (fe43346) - Brandon Guigo
- move utils into shared library - (f0f414b) - Brandon Guigo
- move services into a shared directory - (d430c40) - Brandon Guigo
- linter - (4f38427) - Brandon Guigo
- move amqp and age encryption utils to their services - (e56a537) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.9 grpc@grpc/v0.2.0-rc.9 mail-server@mail-server/v0.1.0-rc.9 mail@mail/v0.1.0-rc.9 productivity@productivity/v0.11.0-rc.9 [skip ci] - (5e5d565) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.8 grpc@grpc/v0.2.0-rc.8 mail-server@mail-server/v0.1.0-rc.8 mail@mail/v0.1.0-rc.8 productivity@productivity/v0.11.0-rc.8 [skip ci] - (34a482a) - GitHub Actions

- - -

## productivity/v0.11.0-rc.8 - 2025-09-12
#### Bug Fixes
- unit tests for mail-server - (556eacb) - Brandon Guigo
- linter - (dc7b49b) - Brandon Guigo
- linter - (e472637) - Brandon Guigo
- worker only treated one message (double ack) - (9da89af) - Brandon Guigo
- parsing of the mail payload - (0cba80b) - Brandon Guigo
- retry finally works with rabbitmq expiration - (bc34795) - Brandon Guigo
- worker processing registers message - (0796c81) - Brandon Guigo
- some error with retry publishing - (37f859d) - Brandon Guigo
- revert: "feat: start of the implementation of the gRPC calls to manage the sending emails" - (ba1a891) - Brandon Guigo
- revert: "feat: setup the cron for send email" - (b492a51) - Brandon Guigo
#### Features
- make grpc call work - (c685348) - Brandon Guigo
- add grpc client to mail server - (e51a369) - Brandon Guigo
- send email is working without grpc calls - (fe51726) - Brandon Guigo
- continue work on sending the emails - (176f195) - Brandon Guigo
- setup dkim signing - (0cf38ea) - Brandon Guigo
- get the mx record and setup the sending loop for each recipients - (4a8d8d9) - Brandon Guigo
- setup the structure of the send methods - (d0b53ed) - Brandon Guigo
- setup amqp worker to listen to retry worker too - (9266adf) - Brandon Guigo
- make amqp consumer / producer configuration totally via - (5fcc202) - Brandon Guigo
- setup the send worker for emails - (b324e62) - Brandon Guigo
- start of the implementation of the gRPC calls to manage the sending emails - (ef92c62) - Brandon Guigo
- setup the cron for send email - (2f2b18e) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.7 grpc@grpc/v0.2.0-rc.7 mail-server@mail-server/v0.1.0-rc.7 mail@mail/v0.1.0-rc.7 productivity@productivity/v0.11.0-rc.7 [skip ci] - (24990d2) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.6 grpc@grpc/v0.2.0-rc.6 mail-server@mail-server/v0.1.0-rc.6 mail@mail/v0.1.0-rc.6 productivity@productivity/v0.11.0-rc.6 [skip ci] - (e000019) - GitHub Actions

- - -

## productivity/v0.11.0-rc.6 - 2025-09-12
#### Bug Fixes
- tests - (cd68156) - Brandon Guigo
- completely disable auth on receive email server - (5f71329) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.5 grpc@grpc/v0.2.0-rc.5 mail-server@mail-server/v0.1.0-rc.5 mail@mail/v0.1.0-rc.5 productivity@productivity/v0.11.0-rc.5 [skip ci] - (d55f603) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.4 grpc@grpc/v0.2.0-rc.4 mail-server@mail-server/v0.1.0-rc.4 mail@mail/v0.1.0-rc.4 productivity@productivity/v0.11.0-rc.4 [skip ci] - (f890081) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.3 grpc@grpc/v0.2.0-rc.3 mail-server@mail-server/v0.1.0-rc.3 mail@mail/v0.1.0-rc.3 productivity@productivity/v0.11.0-rc.3 [skip ci] - (ebe0561) - GitHub Actions

- - -

## productivity/v0.11.0-rc.3 - 2025-09-12
#### Bug Fixes
- centralize dockerfiles and allow build with grpc in the monorepo - (95b2dfc) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.2 grpc@grpc/v0.2.0-rc.2 mail-server@mail-server/v0.1.0-rc.2 mail@mail/v0.1.0-rc.2 productivity@productivity/v0.11.0-rc.2 [skip ci] - (03832fc) - GitHub Actions

- - -

## productivity/v0.11.0-rc.2 - 2025-09-12
#### Bug Fixes
- remove maizzle dockerfile code for mail-server - (18cb64e) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.1 grpc@grpc/v0.2.0-rc.1 mail-server@mail-server/v0.1.0-rc.1 mail@mail/v0.1.0-rc.1 productivity@productivity/v0.11.0-rc.1 [skip ci] - (873e29b) - GitHub Actions

- - -

## productivity/v0.11.0-rc.1 - 2025-09-12
#### Bug Fixes
- linter for mail and mail-server - (401774e) - Brandon Guigo
- grpc linter + mail in test ci/cd + fix error in smtp server test - (3b67a14) - Brandon Guigo
- receive email without attachement leads to no content disposition - (f8e0ac5) - Brandon Guigo
#### Features
- add bruno collections and fix errors - (b6b303c) - Brandon Guigo
- upload file and email to storage with transactions for all recipients - (3e69645) - Brandon Guigo
- update grpc to use latest version + configure dev docker compose to use go workspaces + add grpc to get public key - (370934b) - Brandon Guigo
- make the smtp server handle anonymous auth mecanism + use emersion packages inside the test script - (e807664) - Brandon Guigo
- ack the message when processing is done - (080b6fb) - Brandon Guigo
- parse the newly added amqp message - (e38584a) - Brandon Guigo
- add smtp server tests and mocks - (8da1f07) - Brandon Guigo
- make the smtp server catch all the infos of the sender before sending amqp - (4b5fe3c) - Brandon Guigo
- configure rspamd in docker - (32aa09a) - Brandon Guigo
- parse email attachements in the mail server - (5d54950) - Brandon Guigo
- configure amqp producer in mail server - (4cb56c6) - Brandon Guigo
- configure amqp consumer and producer - (325623a) - Brandon Guigo
- fix middleware and add mail api - (b2f1de8) - Brandon Guigo
- setup mail-server - (395efe7) - Brandon Guigo
#### Miscellaneous Chores
- add mail and mail-server to cog.toml - (541e204) - Brandon Guigo

- - -

## mail-server/v0.1.0-rc.14 - 2025-09-12
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.13 grpc@grpc/v0.2.0-rc.13 mail-server@mail-server/v0.1.0-rc.13 mail@mail/v0.1.0-rc.13 productivity@productivity/v0.11.0-rc.13 shared@shared/v0.0.1-rc.4 [skip ci] - (0e70765) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.12 grpc@grpc/v0.2.0-rc.12 mail-server@mail-server/v0.1.0-rc.12 mail@mail/v0.1.0-rc.12 productivity@productivity/v0.11.0-rc.12 shared@shared/v0.0.1-rc.3 [skip ci] - (68e71b6) - GitHub Actions

- - -

## shared/v0.0.1-rc.3 - 2025-09-12
#### Bug Fixes
- send email model + controller issue - (8b159dc) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.11 grpc@grpc/v0.2.0-rc.11 mail-server@mail-server/v0.1.0-rc.11 mail@mail/v0.1.0-rc.11 productivity@productivity/v0.11.0-rc.11 shared@shared/v0.0.1-rc.2 [skip ci] - (8c59737) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.10 grpc@grpc/v0.2.0-rc.10 mail-server@mail-server/v0.1.0-rc.10 mail@mail/v0.1.0-rc.10 productivity@productivity/v0.11.0-rc.10 shared@shared/v0.0.1-rc.1 [skip ci] - (e44c01a) - GitHub Actions

- - -

## shared/v0.0.1-rc.1 - 2025-09-12
#### Bug Fixes
- error when sending a message through amqp - (e396868) - Brandon Guigo
- tests and add a script to run all tests - (57a16fe) - Brandon Guigo
- refactor mail-server - (fe43346) - Brandon Guigo
- move utils into shared library - (f0f414b) - Brandon Guigo
- move services into a shared directory - (d430c40) - Brandon Guigo
- linter - (4f38427) - Brandon Guigo
- move amqp and age encryption utils to their services - (e56a537) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.9 grpc@grpc/v0.2.0-rc.9 mail-server@mail-server/v0.1.0-rc.9 mail@mail/v0.1.0-rc.9 productivity@productivity/v0.11.0-rc.9 [skip ci] - (5e5d565) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.8 grpc@grpc/v0.2.0-rc.8 mail-server@mail-server/v0.1.0-rc.8 mail@mail/v0.1.0-rc.8 productivity@productivity/v0.11.0-rc.8 [skip ci] - (34a482a) - GitHub Actions

- - -

## productivity/v0.11.0-rc.8 - 2025-09-12
#### Bug Fixes
- unit tests for mail-server - (556eacb) - Brandon Guigo
- linter - (dc7b49b) - Brandon Guigo
- linter - (e472637) - Brandon Guigo
- worker only treated one message (double ack) - (9da89af) - Brandon Guigo
- parsing of the mail payload - (0cba80b) - Brandon Guigo
- retry finally works with rabbitmq expiration - (bc34795) - Brandon Guigo
- worker processing registers message - (0796c81) - Brandon Guigo
- some error with retry publishing - (37f859d) - Brandon Guigo
- revert: "feat: start of the implementation of the gRPC calls to manage the sending emails" - (ba1a891) - Brandon Guigo
- revert: "feat: setup the cron for send email" - (b492a51) - Brandon Guigo
#### Features
- make grpc call work - (c685348) - Brandon Guigo
- add grpc client to mail server - (e51a369) - Brandon Guigo
- send email is working without grpc calls - (fe51726) - Brandon Guigo
- continue work on sending the emails - (176f195) - Brandon Guigo
- setup dkim signing - (0cf38ea) - Brandon Guigo
- get the mx record and setup the sending loop for each recipients - (4a8d8d9) - Brandon Guigo
- setup the structure of the send methods - (d0b53ed) - Brandon Guigo
- setup amqp worker to listen to retry worker too - (9266adf) - Brandon Guigo
- make amqp consumer / producer configuration totally via - (5fcc202) - Brandon Guigo
- setup the send worker for emails - (b324e62) - Brandon Guigo
- start of the implementation of the gRPC calls to manage the sending emails - (ef92c62) - Brandon Guigo
- setup the cron for send email - (2f2b18e) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.7 grpc@grpc/v0.2.0-rc.7 mail-server@mail-server/v0.1.0-rc.7 mail@mail/v0.1.0-rc.7 productivity@productivity/v0.11.0-rc.7 [skip ci] - (24990d2) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.6 grpc@grpc/v0.2.0-rc.6 mail-server@mail-server/v0.1.0-rc.6 mail@mail/v0.1.0-rc.6 productivity@productivity/v0.11.0-rc.6 [skip ci] - (e000019) - GitHub Actions

- - -

## productivity/v0.11.0-rc.6 - 2025-09-12
#### Bug Fixes
- tests - (cd68156) - Brandon Guigo
- completely disable auth on receive email server - (5f71329) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.5 grpc@grpc/v0.2.0-rc.5 mail-server@mail-server/v0.1.0-rc.5 mail@mail/v0.1.0-rc.5 productivity@productivity/v0.11.0-rc.5 [skip ci] - (d55f603) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.4 grpc@grpc/v0.2.0-rc.4 mail-server@mail-server/v0.1.0-rc.4 mail@mail/v0.1.0-rc.4 productivity@productivity/v0.11.0-rc.4 [skip ci] - (f890081) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.3 grpc@grpc/v0.2.0-rc.3 mail-server@mail-server/v0.1.0-rc.3 mail@mail/v0.1.0-rc.3 productivity@productivity/v0.11.0-rc.3 [skip ci] - (ebe0561) - GitHub Actions

- - -

## productivity/v0.11.0-rc.3 - 2025-09-12
#### Bug Fixes
- centralize dockerfiles and allow build with grpc in the monorepo - (95b2dfc) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.2 grpc@grpc/v0.2.0-rc.2 mail-server@mail-server/v0.1.0-rc.2 mail@mail/v0.1.0-rc.2 productivity@productivity/v0.11.0-rc.2 [skip ci] - (03832fc) - GitHub Actions

- - -

## productivity/v0.11.0-rc.2 - 2025-09-12
#### Bug Fixes
- remove maizzle dockerfile code for mail-server - (18cb64e) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.1 grpc@grpc/v0.2.0-rc.1 mail-server@mail-server/v0.1.0-rc.1 mail@mail/v0.1.0-rc.1 productivity@productivity/v0.11.0-rc.1 [skip ci] - (873e29b) - GitHub Actions

- - -

## productivity/v0.11.0-rc.1 - 2025-09-12
#### Bug Fixes
- linter for mail and mail-server - (401774e) - Brandon Guigo
- grpc linter + mail in test ci/cd + fix error in smtp server test - (3b67a14) - Brandon Guigo
- receive email without attachement leads to no content disposition - (f8e0ac5) - Brandon Guigo
#### Features
- add bruno collections and fix errors - (b6b303c) - Brandon Guigo
- upload file and email to storage with transactions for all recipients - (3e69645) - Brandon Guigo
- update grpc to use latest version + configure dev docker compose to use go workspaces + add grpc to get public key - (370934b) - Brandon Guigo
- make the smtp server handle anonymous auth mecanism + use emersion packages inside the test script - (e807664) - Brandon Guigo
- ack the message when processing is done - (080b6fb) - Brandon Guigo
- parse the newly added amqp message - (e38584a) - Brandon Guigo
- add smtp server tests and mocks - (8da1f07) - Brandon Guigo
- make the smtp server catch all the infos of the sender before sending amqp - (4b5fe3c) - Brandon Guigo
- configure rspamd in docker - (32aa09a) - Brandon Guigo
- parse email attachements in the mail server - (5d54950) - Brandon Guigo
- configure amqp producer in mail server - (4cb56c6) - Brandon Guigo
- configure amqp consumer and producer - (325623a) - Brandon Guigo
- fix middleware and add mail api - (b2f1de8) - Brandon Guigo
- setup mail-server - (395efe7) - Brandon Guigo
#### Miscellaneous Chores
- add mail and mail-server to cog.toml - (541e204) - Brandon Guigo

- - -

## mail-server/v0.1.0-rc.13 - 2025-09-11
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.12 grpc@grpc/v0.2.0-rc.12 mail-server@mail-server/v0.1.0-rc.12 mail@mail/v0.1.0-rc.12 productivity@productivity/v0.11.0-rc.12 shared@shared/v0.0.1-rc.3 [skip ci] - (68e71b6) - GitHub Actions

- - -

## shared/v0.0.1-rc.3 - 2025-09-11
#### Bug Fixes
- send email model + controller issue - (8b159dc) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.11 grpc@grpc/v0.2.0-rc.11 mail-server@mail-server/v0.1.0-rc.11 mail@mail/v0.1.0-rc.11 productivity@productivity/v0.11.0-rc.11 shared@shared/v0.0.1-rc.2 [skip ci] - (8c59737) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.10 grpc@grpc/v0.2.0-rc.10 mail-server@mail-server/v0.1.0-rc.10 mail@mail/v0.1.0-rc.10 productivity@productivity/v0.11.0-rc.10 shared@shared/v0.0.1-rc.1 [skip ci] - (e44c01a) - GitHub Actions

- - -

## shared/v0.0.1-rc.1 - 2025-09-11
#### Bug Fixes
- error when sending a message through amqp - (e396868) - Brandon Guigo
- tests and add a script to run all tests - (57a16fe) - Brandon Guigo
- refactor mail-server - (fe43346) - Brandon Guigo
- move utils into shared library - (f0f414b) - Brandon Guigo
- move services into a shared directory - (d430c40) - Brandon Guigo
- linter - (4f38427) - Brandon Guigo
- move amqp and age encryption utils to their services - (e56a537) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.9 grpc@grpc/v0.2.0-rc.9 mail-server@mail-server/v0.1.0-rc.9 mail@mail/v0.1.0-rc.9 productivity@productivity/v0.11.0-rc.9 [skip ci] - (5e5d565) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.8 grpc@grpc/v0.2.0-rc.8 mail-server@mail-server/v0.1.0-rc.8 mail@mail/v0.1.0-rc.8 productivity@productivity/v0.11.0-rc.8 [skip ci] - (34a482a) - GitHub Actions

- - -

## productivity/v0.11.0-rc.8 - 2025-09-11
#### Bug Fixes
- unit tests for mail-server - (556eacb) - Brandon Guigo
- linter - (dc7b49b) - Brandon Guigo
- linter - (e472637) - Brandon Guigo
- worker only treated one message (double ack) - (9da89af) - Brandon Guigo
- parsing of the mail payload - (0cba80b) - Brandon Guigo
- retry finally works with rabbitmq expiration - (bc34795) - Brandon Guigo
- worker processing registers message - (0796c81) - Brandon Guigo
- some error with retry publishing - (37f859d) - Brandon Guigo
- revert: "feat: start of the implementation of the gRPC calls to manage the sending emails" - (ba1a891) - Brandon Guigo
- revert: "feat: setup the cron for send email" - (b492a51) - Brandon Guigo
#### Features
- make grpc call work - (c685348) - Brandon Guigo
- add grpc client to mail server - (e51a369) - Brandon Guigo
- send email is working without grpc calls - (fe51726) - Brandon Guigo
- continue work on sending the emails - (176f195) - Brandon Guigo
- setup dkim signing - (0cf38ea) - Brandon Guigo
- get the mx record and setup the sending loop for each recipients - (4a8d8d9) - Brandon Guigo
- setup the structure of the send methods - (d0b53ed) - Brandon Guigo
- setup amqp worker to listen to retry worker too - (9266adf) - Brandon Guigo
- make amqp consumer / producer configuration totally via - (5fcc202) - Brandon Guigo
- setup the send worker for emails - (b324e62) - Brandon Guigo
- start of the implementation of the gRPC calls to manage the sending emails - (ef92c62) - Brandon Guigo
- setup the cron for send email - (2f2b18e) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.7 grpc@grpc/v0.2.0-rc.7 mail-server@mail-server/v0.1.0-rc.7 mail@mail/v0.1.0-rc.7 productivity@productivity/v0.11.0-rc.7 [skip ci] - (24990d2) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.6 grpc@grpc/v0.2.0-rc.6 mail-server@mail-server/v0.1.0-rc.6 mail@mail/v0.1.0-rc.6 productivity@productivity/v0.11.0-rc.6 [skip ci] - (e000019) - GitHub Actions

- - -

## productivity/v0.11.0-rc.6 - 2025-09-11
#### Bug Fixes
- tests - (cd68156) - Brandon Guigo
- completely disable auth on receive email server - (5f71329) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.5 grpc@grpc/v0.2.0-rc.5 mail-server@mail-server/v0.1.0-rc.5 mail@mail/v0.1.0-rc.5 productivity@productivity/v0.11.0-rc.5 [skip ci] - (d55f603) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.4 grpc@grpc/v0.2.0-rc.4 mail-server@mail-server/v0.1.0-rc.4 mail@mail/v0.1.0-rc.4 productivity@productivity/v0.11.0-rc.4 [skip ci] - (f890081) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.3 grpc@grpc/v0.2.0-rc.3 mail-server@mail-server/v0.1.0-rc.3 mail@mail/v0.1.0-rc.3 productivity@productivity/v0.11.0-rc.3 [skip ci] - (ebe0561) - GitHub Actions

- - -

## productivity/v0.11.0-rc.3 - 2025-09-11
#### Bug Fixes
- centralize dockerfiles and allow build with grpc in the monorepo - (95b2dfc) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.2 grpc@grpc/v0.2.0-rc.2 mail-server@mail-server/v0.1.0-rc.2 mail@mail/v0.1.0-rc.2 productivity@productivity/v0.11.0-rc.2 [skip ci] - (03832fc) - GitHub Actions

- - -

## productivity/v0.11.0-rc.2 - 2025-09-11
#### Bug Fixes
- remove maizzle dockerfile code for mail-server - (18cb64e) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.1 grpc@grpc/v0.2.0-rc.1 mail-server@mail-server/v0.1.0-rc.1 mail@mail/v0.1.0-rc.1 productivity@productivity/v0.11.0-rc.1 [skip ci] - (873e29b) - GitHub Actions

- - -

## productivity/v0.11.0-rc.1 - 2025-09-11
#### Bug Fixes
- linter for mail and mail-server - (401774e) - Brandon Guigo
- grpc linter + mail in test ci/cd + fix error in smtp server test - (3b67a14) - Brandon Guigo
- receive email without attachement leads to no content disposition - (f8e0ac5) - Brandon Guigo
#### Features
- add bruno collections and fix errors - (b6b303c) - Brandon Guigo
- upload file and email to storage with transactions for all recipients - (3e69645) - Brandon Guigo
- update grpc to use latest version + configure dev docker compose to use go workspaces + add grpc to get public key - (370934b) - Brandon Guigo
- make the smtp server handle anonymous auth mecanism + use emersion packages inside the test script - (e807664) - Brandon Guigo
- ack the message when processing is done - (080b6fb) - Brandon Guigo
- parse the newly added amqp message - (e38584a) - Brandon Guigo
- add smtp server tests and mocks - (8da1f07) - Brandon Guigo
- make the smtp server catch all the infos of the sender before sending amqp - (4b5fe3c) - Brandon Guigo
- configure rspamd in docker - (32aa09a) - Brandon Guigo
- parse email attachements in the mail server - (5d54950) - Brandon Guigo
- configure amqp producer in mail server - (4cb56c6) - Brandon Guigo
- configure amqp consumer and producer - (325623a) - Brandon Guigo
- fix middleware and add mail api - (b2f1de8) - Brandon Guigo
- setup mail-server - (395efe7) - Brandon Guigo
#### Miscellaneous Chores
- add mail and mail-server to cog.toml - (541e204) - Brandon Guigo

- - -

## mail-server/v0.1.0-rc.12 - 2025-09-10
#### Bug Fixes
- send email model + controller issue - (8b159dc) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.11 grpc@grpc/v0.2.0-rc.11 mail-server@mail-server/v0.1.0-rc.11 mail@mail/v0.1.0-rc.11 productivity@productivity/v0.11.0-rc.11 shared@shared/v0.0.1-rc.2 [skip ci] - (8c59737) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.10 grpc@grpc/v0.2.0-rc.10 mail-server@mail-server/v0.1.0-rc.10 mail@mail/v0.1.0-rc.10 productivity@productivity/v0.11.0-rc.10 shared@shared/v0.0.1-rc.1 [skip ci] - (e44c01a) - GitHub Actions

- - -

## shared/v0.0.1-rc.1 - 2025-09-10
#### Bug Fixes
- error when sending a message through amqp - (e396868) - Brandon Guigo
- tests and add a script to run all tests - (57a16fe) - Brandon Guigo
- refactor mail-server - (fe43346) - Brandon Guigo
- move utils into shared library - (f0f414b) - Brandon Guigo
- move services into a shared directory - (d430c40) - Brandon Guigo
- linter - (4f38427) - Brandon Guigo
- move amqp and age encryption utils to their services - (e56a537) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.9 grpc@grpc/v0.2.0-rc.9 mail-server@mail-server/v0.1.0-rc.9 mail@mail/v0.1.0-rc.9 productivity@productivity/v0.11.0-rc.9 [skip ci] - (5e5d565) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.8 grpc@grpc/v0.2.0-rc.8 mail-server@mail-server/v0.1.0-rc.8 mail@mail/v0.1.0-rc.8 productivity@productivity/v0.11.0-rc.8 [skip ci] - (34a482a) - GitHub Actions

- - -

## productivity/v0.11.0-rc.8 - 2025-09-10
#### Bug Fixes
- unit tests for mail-server - (556eacb) - Brandon Guigo
- linter - (dc7b49b) - Brandon Guigo
- linter - (e472637) - Brandon Guigo
- worker only treated one message (double ack) - (9da89af) - Brandon Guigo
- parsing of the mail payload - (0cba80b) - Brandon Guigo
- retry finally works with rabbitmq expiration - (bc34795) - Brandon Guigo
- worker processing registers message - (0796c81) - Brandon Guigo
- some error with retry publishing - (37f859d) - Brandon Guigo
- revert: "feat: start of the implementation of the gRPC calls to manage the sending emails" - (ba1a891) - Brandon Guigo
- revert: "feat: setup the cron for send email" - (b492a51) - Brandon Guigo
#### Features
- make grpc call work - (c685348) - Brandon Guigo
- add grpc client to mail server - (e51a369) - Brandon Guigo
- send email is working without grpc calls - (fe51726) - Brandon Guigo
- continue work on sending the emails - (176f195) - Brandon Guigo
- setup dkim signing - (0cf38ea) - Brandon Guigo
- get the mx record and setup the sending loop for each recipients - (4a8d8d9) - Brandon Guigo
- setup the structure of the send methods - (d0b53ed) - Brandon Guigo
- setup amqp worker to listen to retry worker too - (9266adf) - Brandon Guigo
- make amqp consumer / producer configuration totally via - (5fcc202) - Brandon Guigo
- setup the send worker for emails - (b324e62) - Brandon Guigo
- start of the implementation of the gRPC calls to manage the sending emails - (ef92c62) - Brandon Guigo
- setup the cron for send email - (2f2b18e) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.7 grpc@grpc/v0.2.0-rc.7 mail-server@mail-server/v0.1.0-rc.7 mail@mail/v0.1.0-rc.7 productivity@productivity/v0.11.0-rc.7 [skip ci] - (24990d2) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.6 grpc@grpc/v0.2.0-rc.6 mail-server@mail-server/v0.1.0-rc.6 mail@mail/v0.1.0-rc.6 productivity@productivity/v0.11.0-rc.6 [skip ci] - (e000019) - GitHub Actions

- - -

## productivity/v0.11.0-rc.6 - 2025-09-10
#### Bug Fixes
- tests - (cd68156) - Brandon Guigo
- completely disable auth on receive email server - (5f71329) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.5 grpc@grpc/v0.2.0-rc.5 mail-server@mail-server/v0.1.0-rc.5 mail@mail/v0.1.0-rc.5 productivity@productivity/v0.11.0-rc.5 [skip ci] - (d55f603) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.4 grpc@grpc/v0.2.0-rc.4 mail-server@mail-server/v0.1.0-rc.4 mail@mail/v0.1.0-rc.4 productivity@productivity/v0.11.0-rc.4 [skip ci] - (f890081) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.3 grpc@grpc/v0.2.0-rc.3 mail-server@mail-server/v0.1.0-rc.3 mail@mail/v0.1.0-rc.3 productivity@productivity/v0.11.0-rc.3 [skip ci] - (ebe0561) - GitHub Actions

- - -

## productivity/v0.11.0-rc.3 - 2025-09-10
#### Bug Fixes
- centralize dockerfiles and allow build with grpc in the monorepo - (95b2dfc) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.2 grpc@grpc/v0.2.0-rc.2 mail-server@mail-server/v0.1.0-rc.2 mail@mail/v0.1.0-rc.2 productivity@productivity/v0.11.0-rc.2 [skip ci] - (03832fc) - GitHub Actions

- - -

## productivity/v0.11.0-rc.2 - 2025-09-10
#### Bug Fixes
- remove maizzle dockerfile code for mail-server - (18cb64e) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.1 grpc@grpc/v0.2.0-rc.1 mail-server@mail-server/v0.1.0-rc.1 mail@mail/v0.1.0-rc.1 productivity@productivity/v0.11.0-rc.1 [skip ci] - (873e29b) - GitHub Actions

- - -

## productivity/v0.11.0-rc.1 - 2025-09-10
#### Bug Fixes
- linter for mail and mail-server - (401774e) - Brandon Guigo
- grpc linter + mail in test ci/cd + fix error in smtp server test - (3b67a14) - Brandon Guigo
- receive email without attachement leads to no content disposition - (f8e0ac5) - Brandon Guigo
#### Features
- add bruno collections and fix errors - (b6b303c) - Brandon Guigo
- upload file and email to storage with transactions for all recipients - (3e69645) - Brandon Guigo
- update grpc to use latest version + configure dev docker compose to use go workspaces + add grpc to get public key - (370934b) - Brandon Guigo
- make the smtp server handle anonymous auth mecanism + use emersion packages inside the test script - (e807664) - Brandon Guigo
- ack the message when processing is done - (080b6fb) - Brandon Guigo
- parse the newly added amqp message - (e38584a) - Brandon Guigo
- add smtp server tests and mocks - (8da1f07) - Brandon Guigo
- make the smtp server catch all the infos of the sender before sending amqp - (4b5fe3c) - Brandon Guigo
- configure rspamd in docker - (32aa09a) - Brandon Guigo
- parse email attachements in the mail server - (5d54950) - Brandon Guigo
- configure amqp producer in mail server - (4cb56c6) - Brandon Guigo
- configure amqp consumer and producer - (325623a) - Brandon Guigo
- fix middleware and add mail api - (b2f1de8) - Brandon Guigo
- setup mail-server - (395efe7) - Brandon Guigo
#### Miscellaneous Chores
- add mail and mail-server to cog.toml - (541e204) - Brandon Guigo

- - -

## mail-server/v0.1.0-rc.11 - 2025-09-10
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.10 grpc@grpc/v0.2.0-rc.10 mail-server@mail-server/v0.1.0-rc.10 mail@mail/v0.1.0-rc.10 productivity@productivity/v0.11.0-rc.10 shared@shared/v0.0.1-rc.1 [skip ci] - (e44c01a) - GitHub Actions

- - -

## shared/v0.0.1-rc.1 - 2025-09-10
#### Bug Fixes
- error when sending a message through amqp - (e396868) - Brandon Guigo
- tests and add a script to run all tests - (57a16fe) - Brandon Guigo
- refactor mail-server - (fe43346) - Brandon Guigo
- move utils into shared library - (f0f414b) - Brandon Guigo
- move services into a shared directory - (d430c40) - Brandon Guigo
- linter - (4f38427) - Brandon Guigo
- move amqp and age encryption utils to their services - (e56a537) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.9 grpc@grpc/v0.2.0-rc.9 mail-server@mail-server/v0.1.0-rc.9 mail@mail/v0.1.0-rc.9 productivity@productivity/v0.11.0-rc.9 [skip ci] - (5e5d565) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.8 grpc@grpc/v0.2.0-rc.8 mail-server@mail-server/v0.1.0-rc.8 mail@mail/v0.1.0-rc.8 productivity@productivity/v0.11.0-rc.8 [skip ci] - (34a482a) - GitHub Actions

- - -

## productivity/v0.11.0-rc.8 - 2025-09-10
#### Bug Fixes
- unit tests for mail-server - (556eacb) - Brandon Guigo
- linter - (dc7b49b) - Brandon Guigo
- linter - (e472637) - Brandon Guigo
- worker only treated one message (double ack) - (9da89af) - Brandon Guigo
- parsing of the mail payload - (0cba80b) - Brandon Guigo
- retry finally works with rabbitmq expiration - (bc34795) - Brandon Guigo
- worker processing registers message - (0796c81) - Brandon Guigo
- some error with retry publishing - (37f859d) - Brandon Guigo
- revert: "feat: start of the implementation of the gRPC calls to manage the sending emails" - (ba1a891) - Brandon Guigo
- revert: "feat: setup the cron for send email" - (b492a51) - Brandon Guigo
#### Features
- make grpc call work - (c685348) - Brandon Guigo
- add grpc client to mail server - (e51a369) - Brandon Guigo
- send email is working without grpc calls - (fe51726) - Brandon Guigo
- continue work on sending the emails - (176f195) - Brandon Guigo
- setup dkim signing - (0cf38ea) - Brandon Guigo
- get the mx record and setup the sending loop for each recipients - (4a8d8d9) - Brandon Guigo
- setup the structure of the send methods - (d0b53ed) - Brandon Guigo
- setup amqp worker to listen to retry worker too - (9266adf) - Brandon Guigo
- make amqp consumer / producer configuration totally via - (5fcc202) - Brandon Guigo
- setup the send worker for emails - (b324e62) - Brandon Guigo
- start of the implementation of the gRPC calls to manage the sending emails - (ef92c62) - Brandon Guigo
- setup the cron for send email - (2f2b18e) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.7 grpc@grpc/v0.2.0-rc.7 mail-server@mail-server/v0.1.0-rc.7 mail@mail/v0.1.0-rc.7 productivity@productivity/v0.11.0-rc.7 [skip ci] - (24990d2) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.6 grpc@grpc/v0.2.0-rc.6 mail-server@mail-server/v0.1.0-rc.6 mail@mail/v0.1.0-rc.6 productivity@productivity/v0.11.0-rc.6 [skip ci] - (e000019) - GitHub Actions

- - -

## productivity/v0.11.0-rc.6 - 2025-09-10
#### Bug Fixes
- tests - (cd68156) - Brandon Guigo
- completely disable auth on receive email server - (5f71329) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.5 grpc@grpc/v0.2.0-rc.5 mail-server@mail-server/v0.1.0-rc.5 mail@mail/v0.1.0-rc.5 productivity@productivity/v0.11.0-rc.5 [skip ci] - (d55f603) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.4 grpc@grpc/v0.2.0-rc.4 mail-server@mail-server/v0.1.0-rc.4 mail@mail/v0.1.0-rc.4 productivity@productivity/v0.11.0-rc.4 [skip ci] - (f890081) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.3 grpc@grpc/v0.2.0-rc.3 mail-server@mail-server/v0.1.0-rc.3 mail@mail/v0.1.0-rc.3 productivity@productivity/v0.11.0-rc.3 [skip ci] - (ebe0561) - GitHub Actions

- - -

## productivity/v0.11.0-rc.3 - 2025-09-10
#### Bug Fixes
- centralize dockerfiles and allow build with grpc in the monorepo - (95b2dfc) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.2 grpc@grpc/v0.2.0-rc.2 mail-server@mail-server/v0.1.0-rc.2 mail@mail/v0.1.0-rc.2 productivity@productivity/v0.11.0-rc.2 [skip ci] - (03832fc) - GitHub Actions

- - -

## productivity/v0.11.0-rc.2 - 2025-09-10
#### Bug Fixes
- remove maizzle dockerfile code for mail-server - (18cb64e) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.1 grpc@grpc/v0.2.0-rc.1 mail-server@mail-server/v0.1.0-rc.1 mail@mail/v0.1.0-rc.1 productivity@productivity/v0.11.0-rc.1 [skip ci] - (873e29b) - GitHub Actions

- - -

## productivity/v0.11.0-rc.1 - 2025-09-10
#### Bug Fixes
- linter for mail and mail-server - (401774e) - Brandon Guigo
- grpc linter + mail in test ci/cd + fix error in smtp server test - (3b67a14) - Brandon Guigo
- receive email without attachement leads to no content disposition - (f8e0ac5) - Brandon Guigo
#### Features
- add bruno collections and fix errors - (b6b303c) - Brandon Guigo
- upload file and email to storage with transactions for all recipients - (3e69645) - Brandon Guigo
- update grpc to use latest version + configure dev docker compose to use go workspaces + add grpc to get public key - (370934b) - Brandon Guigo
- make the smtp server handle anonymous auth mecanism + use emersion packages inside the test script - (e807664) - Brandon Guigo
- ack the message when processing is done - (080b6fb) - Brandon Guigo
- parse the newly added amqp message - (e38584a) - Brandon Guigo
- add smtp server tests and mocks - (8da1f07) - Brandon Guigo
- make the smtp server catch all the infos of the sender before sending amqp - (4b5fe3c) - Brandon Guigo
- configure rspamd in docker - (32aa09a) - Brandon Guigo
- parse email attachements in the mail server - (5d54950) - Brandon Guigo
- configure amqp producer in mail server - (4cb56c6) - Brandon Guigo
- configure amqp consumer and producer - (325623a) - Brandon Guigo
- fix middleware and add mail api - (b2f1de8) - Brandon Guigo
- setup mail-server - (395efe7) - Brandon Guigo
#### Miscellaneous Chores
- add mail and mail-server to cog.toml - (541e204) - Brandon Guigo

- - -

## mail-server/v0.1.0-rc.10 - 2025-08-12
#### Bug Fixes
- error when sending a message through amqp - (e396868) - Brandon Guigo
- tests and add a script to run all tests - (57a16fe) - Brandon Guigo
- refactor mail-server - (fe43346) - Brandon Guigo
- move utils into shared library - (f0f414b) - Brandon Guigo
- move services into a shared directory - (d430c40) - Brandon Guigo
- linter - (4f38427) - Brandon Guigo
- move amqp and age encryption utils to their services - (e56a537) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.9 grpc@grpc/v0.2.0-rc.9 mail-server@mail-server/v0.1.0-rc.9 mail@mail/v0.1.0-rc.9 productivity@productivity/v0.11.0-rc.9 [skip ci] - (5e5d565) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.8 grpc@grpc/v0.2.0-rc.8 mail-server@mail-server/v0.1.0-rc.8 mail@mail/v0.1.0-rc.8 productivity@productivity/v0.11.0-rc.8 [skip ci] - (34a482a) - GitHub Actions

- - -

## productivity/v0.11.0-rc.8 - 2025-08-12
#### Bug Fixes
- unit tests for mail-server - (556eacb) - Brandon Guigo
- linter - (dc7b49b) - Brandon Guigo
- linter - (e472637) - Brandon Guigo
- worker only treated one message (double ack) - (9da89af) - Brandon Guigo
- parsing of the mail payload - (0cba80b) - Brandon Guigo
- retry finally works with rabbitmq expiration - (bc34795) - Brandon Guigo
- worker processing registers message - (0796c81) - Brandon Guigo
- some error with retry publishing - (37f859d) - Brandon Guigo
- revert: "feat: start of the implementation of the gRPC calls to manage the sending emails" - (ba1a891) - Brandon Guigo
- revert: "feat: setup the cron for send email" - (b492a51) - Brandon Guigo
#### Features
- make grpc call work - (c685348) - Brandon Guigo
- add grpc client to mail server - (e51a369) - Brandon Guigo
- send email is working without grpc calls - (fe51726) - Brandon Guigo
- continue work on sending the emails - (176f195) - Brandon Guigo
- setup dkim signing - (0cf38ea) - Brandon Guigo
- get the mx record and setup the sending loop for each recipients - (4a8d8d9) - Brandon Guigo
- setup the structure of the send methods - (d0b53ed) - Brandon Guigo
- setup amqp worker to listen to retry worker too - (9266adf) - Brandon Guigo
- make amqp consumer / producer configuration totally via - (5fcc202) - Brandon Guigo
- setup the send worker for emails - (b324e62) - Brandon Guigo
- start of the implementation of the gRPC calls to manage the sending emails - (ef92c62) - Brandon Guigo
- setup the cron for send email - (2f2b18e) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.7 grpc@grpc/v0.2.0-rc.7 mail-server@mail-server/v0.1.0-rc.7 mail@mail/v0.1.0-rc.7 productivity@productivity/v0.11.0-rc.7 [skip ci] - (24990d2) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.6 grpc@grpc/v0.2.0-rc.6 mail-server@mail-server/v0.1.0-rc.6 mail@mail/v0.1.0-rc.6 productivity@productivity/v0.11.0-rc.6 [skip ci] - (e000019) - GitHub Actions

- - -

## productivity/v0.11.0-rc.6 - 2025-08-12
#### Bug Fixes
- tests - (cd68156) - Brandon Guigo
- completely disable auth on receive email server - (5f71329) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.5 grpc@grpc/v0.2.0-rc.5 mail-server@mail-server/v0.1.0-rc.5 mail@mail/v0.1.0-rc.5 productivity@productivity/v0.11.0-rc.5 [skip ci] - (d55f603) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.4 grpc@grpc/v0.2.0-rc.4 mail-server@mail-server/v0.1.0-rc.4 mail@mail/v0.1.0-rc.4 productivity@productivity/v0.11.0-rc.4 [skip ci] - (f890081) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.3 grpc@grpc/v0.2.0-rc.3 mail-server@mail-server/v0.1.0-rc.3 mail@mail/v0.1.0-rc.3 productivity@productivity/v0.11.0-rc.3 [skip ci] - (ebe0561) - GitHub Actions

- - -

## productivity/v0.11.0-rc.3 - 2025-08-12
#### Bug Fixes
- centralize dockerfiles and allow build with grpc in the monorepo - (95b2dfc) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.2 grpc@grpc/v0.2.0-rc.2 mail-server@mail-server/v0.1.0-rc.2 mail@mail/v0.1.0-rc.2 productivity@productivity/v0.11.0-rc.2 [skip ci] - (03832fc) - GitHub Actions

- - -

## productivity/v0.11.0-rc.2 - 2025-08-12
#### Bug Fixes
- remove maizzle dockerfile code for mail-server - (18cb64e) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.1 grpc@grpc/v0.2.0-rc.1 mail-server@mail-server/v0.1.0-rc.1 mail@mail/v0.1.0-rc.1 productivity@productivity/v0.11.0-rc.1 [skip ci] - (873e29b) - GitHub Actions

- - -

## productivity/v0.11.0-rc.1 - 2025-08-12
#### Bug Fixes
- linter for mail and mail-server - (401774e) - Brandon Guigo
- grpc linter + mail in test ci/cd + fix error in smtp server test - (3b67a14) - Brandon Guigo
- receive email without attachement leads to no content disposition - (f8e0ac5) - Brandon Guigo
#### Features
- add bruno collections and fix errors - (b6b303c) - Brandon Guigo
- upload file and email to storage with transactions for all recipients - (3e69645) - Brandon Guigo
- update grpc to use latest version + configure dev docker compose to use go workspaces + add grpc to get public key - (370934b) - Brandon Guigo
- make the smtp server handle anonymous auth mecanism + use emersion packages inside the test script - (e807664) - Brandon Guigo
- ack the message when processing is done - (080b6fb) - Brandon Guigo
- parse the newly added amqp message - (e38584a) - Brandon Guigo
- add smtp server tests and mocks - (8da1f07) - Brandon Guigo
- make the smtp server catch all the infos of the sender before sending amqp - (4b5fe3c) - Brandon Guigo
- configure rspamd in docker - (32aa09a) - Brandon Guigo
- parse email attachements in the mail server - (5d54950) - Brandon Guigo
- configure amqp producer in mail server - (4cb56c6) - Brandon Guigo
- configure amqp consumer and producer - (325623a) - Brandon Guigo
- fix middleware and add mail api - (b2f1de8) - Brandon Guigo
- setup mail-server - (395efe7) - Brandon Guigo
#### Miscellaneous Chores
- add mail and mail-server to cog.toml - (541e204) - Brandon Guigo

- - -

## mail-server/v0.1.0-rc.9 - 2025-08-12
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.8 grpc@grpc/v0.2.0-rc.8 mail-server@mail-server/v0.1.0-rc.8 mail@mail/v0.1.0-rc.8 productivity@productivity/v0.11.0-rc.8 [skip ci] - (34a482a) - GitHub Actions

- - -

## productivity/v0.11.0-rc.8 - 2025-08-12
#### Bug Fixes
- unit tests for mail-server - (556eacb) - Brandon Guigo
- linter - (dc7b49b) - Brandon Guigo
- linter - (e472637) - Brandon Guigo
- worker only treated one message (double ack) - (9da89af) - Brandon Guigo
- parsing of the mail payload - (0cba80b) - Brandon Guigo
- retry finally works with rabbitmq expiration - (bc34795) - Brandon Guigo
- worker processing registers message - (0796c81) - Brandon Guigo
- some error with retry publishing - (37f859d) - Brandon Guigo
- revert: "feat: start of the implementation of the gRPC calls to manage the sending emails" - (ba1a891) - Brandon Guigo
- revert: "feat: setup the cron for send email" - (b492a51) - Brandon Guigo
#### Features
- make grpc call work - (c685348) - Brandon Guigo
- add grpc client to mail server - (e51a369) - Brandon Guigo
- send email is working without grpc calls - (fe51726) - Brandon Guigo
- continue work on sending the emails - (176f195) - Brandon Guigo
- setup dkim signing - (0cf38ea) - Brandon Guigo
- get the mx record and setup the sending loop for each recipients - (4a8d8d9) - Brandon Guigo
- setup the structure of the send methods - (d0b53ed) - Brandon Guigo
- setup amqp worker to listen to retry worker too - (9266adf) - Brandon Guigo
- make amqp consumer / producer configuration totally via - (5fcc202) - Brandon Guigo
- setup the send worker for emails - (b324e62) - Brandon Guigo
- start of the implementation of the gRPC calls to manage the sending emails - (ef92c62) - Brandon Guigo
- setup the cron for send email - (2f2b18e) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.7 grpc@grpc/v0.2.0-rc.7 mail-server@mail-server/v0.1.0-rc.7 mail@mail/v0.1.0-rc.7 productivity@productivity/v0.11.0-rc.7 [skip ci] - (24990d2) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.6 grpc@grpc/v0.2.0-rc.6 mail-server@mail-server/v0.1.0-rc.6 mail@mail/v0.1.0-rc.6 productivity@productivity/v0.11.0-rc.6 [skip ci] - (e000019) - GitHub Actions

- - -

## productivity/v0.11.0-rc.6 - 2025-08-12
#### Bug Fixes
- tests - (cd68156) - Brandon Guigo
- completely disable auth on receive email server - (5f71329) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.5 grpc@grpc/v0.2.0-rc.5 mail-server@mail-server/v0.1.0-rc.5 mail@mail/v0.1.0-rc.5 productivity@productivity/v0.11.0-rc.5 [skip ci] - (d55f603) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.4 grpc@grpc/v0.2.0-rc.4 mail-server@mail-server/v0.1.0-rc.4 mail@mail/v0.1.0-rc.4 productivity@productivity/v0.11.0-rc.4 [skip ci] - (f890081) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.3 grpc@grpc/v0.2.0-rc.3 mail-server@mail-server/v0.1.0-rc.3 mail@mail/v0.1.0-rc.3 productivity@productivity/v0.11.0-rc.3 [skip ci] - (ebe0561) - GitHub Actions

- - -

## productivity/v0.11.0-rc.3 - 2025-08-12
#### Bug Fixes
- centralize dockerfiles and allow build with grpc in the monorepo - (95b2dfc) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.2 grpc@grpc/v0.2.0-rc.2 mail-server@mail-server/v0.1.0-rc.2 mail@mail/v0.1.0-rc.2 productivity@productivity/v0.11.0-rc.2 [skip ci] - (03832fc) - GitHub Actions

- - -

## productivity/v0.11.0-rc.2 - 2025-08-12
#### Bug Fixes
- remove maizzle dockerfile code for mail-server - (18cb64e) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.1 grpc@grpc/v0.2.0-rc.1 mail-server@mail-server/v0.1.0-rc.1 mail@mail/v0.1.0-rc.1 productivity@productivity/v0.11.0-rc.1 [skip ci] - (873e29b) - GitHub Actions

- - -

## productivity/v0.11.0-rc.1 - 2025-08-12
#### Bug Fixes
- linter for mail and mail-server - (401774e) - Brandon Guigo
- grpc linter + mail in test ci/cd + fix error in smtp server test - (3b67a14) - Brandon Guigo
- receive email without attachement leads to no content disposition - (f8e0ac5) - Brandon Guigo
#### Features
- add bruno collections and fix errors - (b6b303c) - Brandon Guigo
- upload file and email to storage with transactions for all recipients - (3e69645) - Brandon Guigo
- update grpc to use latest version + configure dev docker compose to use go workspaces + add grpc to get public key - (370934b) - Brandon Guigo
- make the smtp server handle anonymous auth mecanism + use emersion packages inside the test script - (e807664) - Brandon Guigo
- ack the message when processing is done - (080b6fb) - Brandon Guigo
- parse the newly added amqp message - (e38584a) - Brandon Guigo
- add smtp server tests and mocks - (8da1f07) - Brandon Guigo
- make the smtp server catch all the infos of the sender before sending amqp - (4b5fe3c) - Brandon Guigo
- configure rspamd in docker - (32aa09a) - Brandon Guigo
- parse email attachements in the mail server - (5d54950) - Brandon Guigo
- configure amqp producer in mail server - (4cb56c6) - Brandon Guigo
- configure amqp consumer and producer - (325623a) - Brandon Guigo
- fix middleware and add mail api - (b2f1de8) - Brandon Guigo
- setup mail-server - (395efe7) - Brandon Guigo
#### Miscellaneous Chores
- add mail and mail-server to cog.toml - (541e204) - Brandon Guigo

- - -

## mail-server/v0.1.0-rc.8 - 2025-08-12
#### Bug Fixes
- unit tests for mail-server - (556eacb) - Brandon Guigo
- linter - (dc7b49b) - Brandon Guigo
- linter - (e472637) - Brandon Guigo
- worker only treated one message (double ack) - (9da89af) - Brandon Guigo
- parsing of the mail payload - (0cba80b) - Brandon Guigo
- retry finally works with rabbitmq expiration - (bc34795) - Brandon Guigo
- worker processing registers message - (0796c81) - Brandon Guigo
- some error with retry publishing - (37f859d) - Brandon Guigo
- revert: "feat: start of the implementation of the gRPC calls to manage the sending emails" - (ba1a891) - Brandon Guigo
- revert: "feat: setup the cron for send email" - (b492a51) - Brandon Guigo
#### Features
- make grpc call work - (c685348) - Brandon Guigo
- add grpc client to mail server - (e51a369) - Brandon Guigo
- send email is working without grpc calls - (fe51726) - Brandon Guigo
- continue work on sending the emails - (176f195) - Brandon Guigo
- setup dkim signing - (0cf38ea) - Brandon Guigo
- get the mx record and setup the sending loop for each recipients - (4a8d8d9) - Brandon Guigo
- setup the structure of the send methods - (d0b53ed) - Brandon Guigo
- setup amqp worker to listen to retry worker too - (9266adf) - Brandon Guigo
- make amqp consumer / producer configuration totally via - (5fcc202) - Brandon Guigo
- setup the send worker for emails - (b324e62) - Brandon Guigo
- start of the implementation of the gRPC calls to manage the sending emails - (ef92c62) - Brandon Guigo
- setup the cron for send email - (2f2b18e) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.7 grpc@grpc/v0.2.0-rc.7 mail-server@mail-server/v0.1.0-rc.7 mail@mail/v0.1.0-rc.7 productivity@productivity/v0.11.0-rc.7 [skip ci] - (24990d2) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.6 grpc@grpc/v0.2.0-rc.6 mail-server@mail-server/v0.1.0-rc.6 mail@mail/v0.1.0-rc.6 productivity@productivity/v0.11.0-rc.6 [skip ci] - (e000019) - GitHub Actions

- - -

## productivity/v0.11.0-rc.6 - 2025-08-12
#### Bug Fixes
- tests - (cd68156) - Brandon Guigo
- completely disable auth on receive email server - (5f71329) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.5 grpc@grpc/v0.2.0-rc.5 mail-server@mail-server/v0.1.0-rc.5 mail@mail/v0.1.0-rc.5 productivity@productivity/v0.11.0-rc.5 [skip ci] - (d55f603) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.4 grpc@grpc/v0.2.0-rc.4 mail-server@mail-server/v0.1.0-rc.4 mail@mail/v0.1.0-rc.4 productivity@productivity/v0.11.0-rc.4 [skip ci] - (f890081) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.3 grpc@grpc/v0.2.0-rc.3 mail-server@mail-server/v0.1.0-rc.3 mail@mail/v0.1.0-rc.3 productivity@productivity/v0.11.0-rc.3 [skip ci] - (ebe0561) - GitHub Actions

- - -

## productivity/v0.11.0-rc.3 - 2025-08-12
#### Bug Fixes
- centralize dockerfiles and allow build with grpc in the monorepo - (95b2dfc) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.2 grpc@grpc/v0.2.0-rc.2 mail-server@mail-server/v0.1.0-rc.2 mail@mail/v0.1.0-rc.2 productivity@productivity/v0.11.0-rc.2 [skip ci] - (03832fc) - GitHub Actions

- - -

## productivity/v0.11.0-rc.2 - 2025-08-12
#### Bug Fixes
- remove maizzle dockerfile code for mail-server - (18cb64e) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.1 grpc@grpc/v0.2.0-rc.1 mail-server@mail-server/v0.1.0-rc.1 mail@mail/v0.1.0-rc.1 productivity@productivity/v0.11.0-rc.1 [skip ci] - (873e29b) - GitHub Actions

- - -

## productivity/v0.11.0-rc.1 - 2025-08-12
#### Bug Fixes
- linter for mail and mail-server - (401774e) - Brandon Guigo
- grpc linter + mail in test ci/cd + fix error in smtp server test - (3b67a14) - Brandon Guigo
- receive email without attachement leads to no content disposition - (f8e0ac5) - Brandon Guigo
#### Features
- add bruno collections and fix errors - (b6b303c) - Brandon Guigo
- upload file and email to storage with transactions for all recipients - (3e69645) - Brandon Guigo
- update grpc to use latest version + configure dev docker compose to use go workspaces + add grpc to get public key - (370934b) - Brandon Guigo
- make the smtp server handle anonymous auth mecanism + use emersion packages inside the test script - (e807664) - Brandon Guigo
- ack the message when processing is done - (080b6fb) - Brandon Guigo
- parse the newly added amqp message - (e38584a) - Brandon Guigo
- add smtp server tests and mocks - (8da1f07) - Brandon Guigo
- make the smtp server catch all the infos of the sender before sending amqp - (4b5fe3c) - Brandon Guigo
- configure rspamd in docker - (32aa09a) - Brandon Guigo
- parse email attachements in the mail server - (5d54950) - Brandon Guigo
- configure amqp producer in mail server - (4cb56c6) - Brandon Guigo
- configure amqp consumer and producer - (325623a) - Brandon Guigo
- fix middleware and add mail api - (b2f1de8) - Brandon Guigo
- setup mail-server - (395efe7) - Brandon Guigo
#### Miscellaneous Chores
- add mail and mail-server to cog.toml - (541e204) - Brandon Guigo

- - -

## mail-server/v0.1.0-rc.7 - 2025-08-08
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.6 grpc@grpc/v0.2.0-rc.6 mail-server@mail-server/v0.1.0-rc.6 mail@mail/v0.1.0-rc.6 productivity@productivity/v0.11.0-rc.6 [skip ci] - (e000019) - GitHub Actions

- - -

## productivity/v0.11.0-rc.6 - 2025-08-08
#### Bug Fixes
- tests - (cd68156) - Brandon Guigo
- completely disable auth on receive email server - (5f71329) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.5 grpc@grpc/v0.2.0-rc.5 mail-server@mail-server/v0.1.0-rc.5 mail@mail/v0.1.0-rc.5 productivity@productivity/v0.11.0-rc.5 [skip ci] - (d55f603) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.4 grpc@grpc/v0.2.0-rc.4 mail-server@mail-server/v0.1.0-rc.4 mail@mail/v0.1.0-rc.4 productivity@productivity/v0.11.0-rc.4 [skip ci] - (f890081) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.3 grpc@grpc/v0.2.0-rc.3 mail-server@mail-server/v0.1.0-rc.3 mail@mail/v0.1.0-rc.3 productivity@productivity/v0.11.0-rc.3 [skip ci] - (ebe0561) - GitHub Actions

- - -

## productivity/v0.11.0-rc.3 - 2025-08-08
#### Bug Fixes
- centralize dockerfiles and allow build with grpc in the monorepo - (95b2dfc) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.2 grpc@grpc/v0.2.0-rc.2 mail-server@mail-server/v0.1.0-rc.2 mail@mail/v0.1.0-rc.2 productivity@productivity/v0.11.0-rc.2 [skip ci] - (03832fc) - GitHub Actions

- - -

## productivity/v0.11.0-rc.2 - 2025-08-08
#### Bug Fixes
- remove maizzle dockerfile code for mail-server - (18cb64e) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.1 grpc@grpc/v0.2.0-rc.1 mail-server@mail-server/v0.1.0-rc.1 mail@mail/v0.1.0-rc.1 productivity@productivity/v0.11.0-rc.1 [skip ci] - (873e29b) - GitHub Actions

- - -

## productivity/v0.11.0-rc.1 - 2025-08-08
#### Bug Fixes
- linter for mail and mail-server - (401774e) - Brandon Guigo
- grpc linter + mail in test ci/cd + fix error in smtp server test - (3b67a14) - Brandon Guigo
- receive email without attachement leads to no content disposition - (f8e0ac5) - Brandon Guigo
#### Features
- add bruno collections and fix errors - (b6b303c) - Brandon Guigo
- upload file and email to storage with transactions for all recipients - (3e69645) - Brandon Guigo
- update grpc to use latest version + configure dev docker compose to use go workspaces + add grpc to get public key - (370934b) - Brandon Guigo
- make the smtp server handle anonymous auth mecanism + use emersion packages inside the test script - (e807664) - Brandon Guigo
- ack the message when processing is done - (080b6fb) - Brandon Guigo
- parse the newly added amqp message - (e38584a) - Brandon Guigo
- add smtp server tests and mocks - (8da1f07) - Brandon Guigo
- make the smtp server catch all the infos of the sender before sending amqp - (4b5fe3c) - Brandon Guigo
- configure rspamd in docker - (32aa09a) - Brandon Guigo
- parse email attachements in the mail server - (5d54950) - Brandon Guigo
- configure amqp producer in mail server - (4cb56c6) - Brandon Guigo
- configure amqp consumer and producer - (325623a) - Brandon Guigo
- fix middleware and add mail api - (b2f1de8) - Brandon Guigo
- setup mail-server - (395efe7) - Brandon Guigo
#### Miscellaneous Chores
- add mail and mail-server to cog.toml - (541e204) - Brandon Guigo

- - -

## mail-server/v0.1.0-rc.6 - 2025-08-07
#### Bug Fixes
- tests - (cd68156) - Brandon Guigo
- completely disable auth on receive email server - (5f71329) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.5 grpc@grpc/v0.2.0-rc.5 mail-server@mail-server/v0.1.0-rc.5 mail@mail/v0.1.0-rc.5 productivity@productivity/v0.11.0-rc.5 [skip ci] - (d55f603) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.4 grpc@grpc/v0.2.0-rc.4 mail-server@mail-server/v0.1.0-rc.4 mail@mail/v0.1.0-rc.4 productivity@productivity/v0.11.0-rc.4 [skip ci] - (f890081) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.3 grpc@grpc/v0.2.0-rc.3 mail-server@mail-server/v0.1.0-rc.3 mail@mail/v0.1.0-rc.3 productivity@productivity/v0.11.0-rc.3 [skip ci] - (ebe0561) - GitHub Actions

- - -

## productivity/v0.11.0-rc.3 - 2025-08-07
#### Bug Fixes
- centralize dockerfiles and allow build with grpc in the monorepo - (95b2dfc) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.2 grpc@grpc/v0.2.0-rc.2 mail-server@mail-server/v0.1.0-rc.2 mail@mail/v0.1.0-rc.2 productivity@productivity/v0.11.0-rc.2 [skip ci] - (03832fc) - GitHub Actions

- - -

## productivity/v0.11.0-rc.2 - 2025-08-07
#### Bug Fixes
- remove maizzle dockerfile code for mail-server - (18cb64e) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.1 grpc@grpc/v0.2.0-rc.1 mail-server@mail-server/v0.1.0-rc.1 mail@mail/v0.1.0-rc.1 productivity@productivity/v0.11.0-rc.1 [skip ci] - (873e29b) - GitHub Actions

- - -

## productivity/v0.11.0-rc.1 - 2025-08-07
#### Bug Fixes
- linter for mail and mail-server - (401774e) - Brandon Guigo
- grpc linter + mail in test ci/cd + fix error in smtp server test - (3b67a14) - Brandon Guigo
- receive email without attachement leads to no content disposition - (f8e0ac5) - Brandon Guigo
#### Features
- add bruno collections and fix errors - (b6b303c) - Brandon Guigo
- upload file and email to storage with transactions for all recipients - (3e69645) - Brandon Guigo
- update grpc to use latest version + configure dev docker compose to use go workspaces + add grpc to get public key - (370934b) - Brandon Guigo
- make the smtp server handle anonymous auth mecanism + use emersion packages inside the test script - (e807664) - Brandon Guigo
- ack the message when processing is done - (080b6fb) - Brandon Guigo
- parse the newly added amqp message - (e38584a) - Brandon Guigo
- add smtp server tests and mocks - (8da1f07) - Brandon Guigo
- make the smtp server catch all the infos of the sender before sending amqp - (4b5fe3c) - Brandon Guigo
- configure rspamd in docker - (32aa09a) - Brandon Guigo
- parse email attachements in the mail server - (5d54950) - Brandon Guigo
- configure amqp producer in mail server - (4cb56c6) - Brandon Guigo
- configure amqp consumer and producer - (325623a) - Brandon Guigo
- fix middleware and add mail api - (b2f1de8) - Brandon Guigo
- setup mail-server - (395efe7) - Brandon Guigo
#### Miscellaneous Chores
- add mail and mail-server to cog.toml - (541e204) - Brandon Guigo

- - -

## mail-server/v0.1.0-rc.5 - 2025-08-05
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.4 grpc@grpc/v0.2.0-rc.4 mail-server@mail-server/v0.1.0-rc.4 mail@mail/v0.1.0-rc.4 productivity@productivity/v0.11.0-rc.4 [skip ci] - (f890081) - GitHub Actions
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.3 grpc@grpc/v0.2.0-rc.3 mail-server@mail-server/v0.1.0-rc.3 mail@mail/v0.1.0-rc.3 productivity@productivity/v0.11.0-rc.3 [skip ci] - (ebe0561) - GitHub Actions

- - -

## productivity/v0.11.0-rc.3 - 2025-08-05
#### Bug Fixes
- centralize dockerfiles and allow build with grpc in the monorepo - (95b2dfc) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.2 grpc@grpc/v0.2.0-rc.2 mail-server@mail-server/v0.1.0-rc.2 mail@mail/v0.1.0-rc.2 productivity@productivity/v0.11.0-rc.2 [skip ci] - (03832fc) - GitHub Actions

- - -

## productivity/v0.11.0-rc.2 - 2025-08-05
#### Bug Fixes
- remove maizzle dockerfile code for mail-server - (18cb64e) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.1 grpc@grpc/v0.2.0-rc.1 mail-server@mail-server/v0.1.0-rc.1 mail@mail/v0.1.0-rc.1 productivity@productivity/v0.11.0-rc.1 [skip ci] - (873e29b) - GitHub Actions

- - -

## productivity/v0.11.0-rc.1 - 2025-08-05
#### Bug Fixes
- linter for mail and mail-server - (401774e) - Brandon Guigo
- grpc linter + mail in test ci/cd + fix error in smtp server test - (3b67a14) - Brandon Guigo
- receive email without attachement leads to no content disposition - (f8e0ac5) - Brandon Guigo
#### Features
- add bruno collections and fix errors - (b6b303c) - Brandon Guigo
- upload file and email to storage with transactions for all recipients - (3e69645) - Brandon Guigo
- update grpc to use latest version + configure dev docker compose to use go workspaces + add grpc to get public key - (370934b) - Brandon Guigo
- make the smtp server handle anonymous auth mecanism + use emersion packages inside the test script - (e807664) - Brandon Guigo
- ack the message when processing is done - (080b6fb) - Brandon Guigo
- parse the newly added amqp message - (e38584a) - Brandon Guigo
- add smtp server tests and mocks - (8da1f07) - Brandon Guigo
- make the smtp server catch all the infos of the sender before sending amqp - (4b5fe3c) - Brandon Guigo
- configure rspamd in docker - (32aa09a) - Brandon Guigo
- parse email attachements in the mail server - (5d54950) - Brandon Guigo
- configure amqp producer in mail server - (4cb56c6) - Brandon Guigo
- configure amqp consumer and producer - (325623a) - Brandon Guigo
- fix middleware and add mail api - (b2f1de8) - Brandon Guigo
- setup mail-server - (395efe7) - Brandon Guigo
#### Miscellaneous Chores
- add mail and mail-server to cog.toml - (541e204) - Brandon Guigo

- - -

## mail-server/v0.1.0-rc.4 - 2025-08-05
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.3 grpc@grpc/v0.2.0-rc.3 mail-server@mail-server/v0.1.0-rc.3 mail@mail/v0.1.0-rc.3 productivity@productivity/v0.11.0-rc.3 [skip ci] - (ebe0561) - GitHub Actions

- - -

## productivity/v0.11.0-rc.3 - 2025-08-05
#### Bug Fixes
- centralize dockerfiles and allow build with grpc in the monorepo - (95b2dfc) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.2 grpc@grpc/v0.2.0-rc.2 mail-server@mail-server/v0.1.0-rc.2 mail@mail/v0.1.0-rc.2 productivity@productivity/v0.11.0-rc.2 [skip ci] - (03832fc) - GitHub Actions

- - -

## productivity/v0.11.0-rc.2 - 2025-08-05
#### Bug Fixes
- remove maizzle dockerfile code for mail-server - (18cb64e) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.1 grpc@grpc/v0.2.0-rc.1 mail-server@mail-server/v0.1.0-rc.1 mail@mail/v0.1.0-rc.1 productivity@productivity/v0.11.0-rc.1 [skip ci] - (873e29b) - GitHub Actions

- - -

## productivity/v0.11.0-rc.1 - 2025-08-05
#### Bug Fixes
- linter for mail and mail-server - (401774e) - Brandon Guigo
- grpc linter + mail in test ci/cd + fix error in smtp server test - (3b67a14) - Brandon Guigo
- receive email without attachement leads to no content disposition - (f8e0ac5) - Brandon Guigo
#### Features
- add bruno collections and fix errors - (b6b303c) - Brandon Guigo
- upload file and email to storage with transactions for all recipients - (3e69645) - Brandon Guigo
- update grpc to use latest version + configure dev docker compose to use go workspaces + add grpc to get public key - (370934b) - Brandon Guigo
- make the smtp server handle anonymous auth mecanism + use emersion packages inside the test script - (e807664) - Brandon Guigo
- ack the message when processing is done - (080b6fb) - Brandon Guigo
- parse the newly added amqp message - (e38584a) - Brandon Guigo
- add smtp server tests and mocks - (8da1f07) - Brandon Guigo
- make the smtp server catch all the infos of the sender before sending amqp - (4b5fe3c) - Brandon Guigo
- configure rspamd in docker - (32aa09a) - Brandon Guigo
- parse email attachements in the mail server - (5d54950) - Brandon Guigo
- configure amqp producer in mail server - (4cb56c6) - Brandon Guigo
- configure amqp consumer and producer - (325623a) - Brandon Guigo
- fix middleware and add mail api - (b2f1de8) - Brandon Guigo
- setup mail-server - (395efe7) - Brandon Guigo
#### Miscellaneous Chores
- add mail and mail-server to cog.toml - (541e204) - Brandon Guigo

- - -

## mail-server/v0.1.0-rc.3 - 2025-08-05
#### Bug Fixes
- centralize dockerfiles and allow build with grpc in the monorepo - (95b2dfc) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.2 grpc@grpc/v0.2.0-rc.2 mail-server@mail-server/v0.1.0-rc.2 mail@mail/v0.1.0-rc.2 productivity@productivity/v0.11.0-rc.2 [skip ci] - (03832fc) - GitHub Actions

- - -

## productivity/v0.11.0-rc.2 - 2025-08-05
#### Bug Fixes
- remove maizzle dockerfile code for mail-server - (18cb64e) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.1 grpc@grpc/v0.2.0-rc.1 mail-server@mail-server/v0.1.0-rc.1 mail@mail/v0.1.0-rc.1 productivity@productivity/v0.11.0-rc.1 [skip ci] - (873e29b) - GitHub Actions

- - -

## productivity/v0.11.0-rc.1 - 2025-08-05
#### Bug Fixes
- linter for mail and mail-server - (401774e) - Brandon Guigo
- grpc linter + mail in test ci/cd + fix error in smtp server test - (3b67a14) - Brandon Guigo
- receive email without attachement leads to no content disposition - (f8e0ac5) - Brandon Guigo
#### Features
- add bruno collections and fix errors - (b6b303c) - Brandon Guigo
- upload file and email to storage with transactions for all recipients - (3e69645) - Brandon Guigo
- update grpc to use latest version + configure dev docker compose to use go workspaces + add grpc to get public key - (370934b) - Brandon Guigo
- make the smtp server handle anonymous auth mecanism + use emersion packages inside the test script - (e807664) - Brandon Guigo
- ack the message when processing is done - (080b6fb) - Brandon Guigo
- parse the newly added amqp message - (e38584a) - Brandon Guigo
- add smtp server tests and mocks - (8da1f07) - Brandon Guigo
- make the smtp server catch all the infos of the sender before sending amqp - (4b5fe3c) - Brandon Guigo
- configure rspamd in docker - (32aa09a) - Brandon Guigo
- parse email attachements in the mail server - (5d54950) - Brandon Guigo
- configure amqp producer in mail server - (4cb56c6) - Brandon Guigo
- configure amqp consumer and producer - (325623a) - Brandon Guigo
- fix middleware and add mail api - (b2f1de8) - Brandon Guigo
- setup mail-server - (395efe7) - Brandon Guigo
#### Miscellaneous Chores
- add mail and mail-server to cog.toml - (541e204) - Brandon Guigo

- - -

## mail-server/v0.1.0-rc.2 - 2025-08-05
#### Bug Fixes
- remove maizzle dockerfile code for mail-server - (18cb64e) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** bump RC versions for auth@auth/v0.11.0-rc.1 grpc@grpc/v0.2.0-rc.1 mail-server@mail-server/v0.1.0-rc.1 mail@mail/v0.1.0-rc.1 productivity@productivity/v0.11.0-rc.1 [skip ci] - (873e29b) - GitHub Actions

- - -

## productivity/v0.11.0-rc.1 - 2025-08-05
#### Bug Fixes
- linter for mail and mail-server - (401774e) - Brandon Guigo
- grpc linter + mail in test ci/cd + fix error in smtp server test - (3b67a14) - Brandon Guigo
- receive email without attachement leads to no content disposition - (f8e0ac5) - Brandon Guigo
#### Features
- add bruno collections and fix errors - (b6b303c) - Brandon Guigo
- upload file and email to storage with transactions for all recipients - (3e69645) - Brandon Guigo
- update grpc to use latest version + configure dev docker compose to use go workspaces + add grpc to get public key - (370934b) - Brandon Guigo
- make the smtp server handle anonymous auth mecanism + use emersion packages inside the test script - (e807664) - Brandon Guigo
- ack the message when processing is done - (080b6fb) - Brandon Guigo
- parse the newly added amqp message - (e38584a) - Brandon Guigo
- add smtp server tests and mocks - (8da1f07) - Brandon Guigo
- make the smtp server catch all the infos of the sender before sending amqp - (4b5fe3c) - Brandon Guigo
- configure rspamd in docker - (32aa09a) - Brandon Guigo
- parse email attachements in the mail server - (5d54950) - Brandon Guigo
- configure amqp producer in mail server - (4cb56c6) - Brandon Guigo
- configure amqp consumer and producer - (325623a) - Brandon Guigo
- fix middleware and add mail api - (b2f1de8) - Brandon Guigo
- setup mail-server - (395efe7) - Brandon Guigo
#### Miscellaneous Chores
- add mail and mail-server to cog.toml - (541e204) - Brandon Guigo

- - -

## mail-server/v0.1.0-rc.1 - 2025-08-05
#### Bug Fixes
- linter for mail and mail-server - (401774e) - Brandon Guigo
- grpc linter + mail in test ci/cd + fix error in smtp server test - (3b67a14) - Brandon Guigo
- receive email without attachement leads to no content disposition - (f8e0ac5) - Brandon Guigo
#### Features
- add bruno collections and fix errors - (b6b303c) - Brandon Guigo
- upload file and email to storage with transactions for all recipients - (3e69645) - Brandon Guigo
- update grpc to use latest version + configure dev docker compose to use go workspaces + add grpc to get public key - (370934b) - Brandon Guigo
- make the smtp server handle anonymous auth mecanism + use emersion packages inside the test script - (e807664) - Brandon Guigo
- ack the message when processing is done - (080b6fb) - Brandon Guigo
- parse the newly added amqp message - (e38584a) - Brandon Guigo
- add smtp server tests and mocks - (8da1f07) - Brandon Guigo
- make the smtp server catch all the infos of the sender before sending amqp - (4b5fe3c) - Brandon Guigo
- configure rspamd in docker - (32aa09a) - Brandon Guigo
- parse email attachements in the mail server - (5d54950) - Brandon Guigo
- configure amqp producer in mail server - (4cb56c6) - Brandon Guigo
- configure amqp consumer and producer - (325623a) - Brandon Guigo
- fix middleware and add mail api - (b2f1de8) - Brandon Guigo
- setup mail-server - (395efe7) - Brandon Guigo
#### Miscellaneous Chores
- add mail and mail-server to cog.toml - (541e204) - Brandon Guigo

- - -

Changelog generated by [cocogitto](https://github.com/cocogitto/cocogitto).