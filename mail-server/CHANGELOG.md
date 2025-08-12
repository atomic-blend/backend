# Changelog
All notable changes to this project will be documented in this file. See [conventional commits](https://www.conventionalcommits.org/) for commit guidelines.

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