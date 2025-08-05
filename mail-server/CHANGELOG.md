# Changelog
All notable changes to this project will be documented in this file. See [conventional commits](https://www.conventionalcommits.org/) for commit guidelines.

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