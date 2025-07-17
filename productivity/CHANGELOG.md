# Changelog
All notable changes to this project will be documented in this file. See [conventional commits](https://www.conventionalcommits.org/) for commit guidelines.

- - -
## productivity-v0.9.0-rc.1 - 2025-07-17
#### Bug Fixes
- remove old cog.toml - (cab208f) - Brandon Guigo
- setup cocogitto to work in a monorepo setup - (9f5e362) - Brandon Guigo
- linter - (6d2f946) - Brandon Guigo
- refactor cron to be gRPC ready - (74ac6d2) - Brandon Guigo
- linter - (4a451f2) - Brandon Guigo
- add missing github.com to go modules - (4f38ac7) - Brandon Guigo
- change module name of productivity go module - (798aaa1) - Brandon Guigo
- tag create test - (dec70a9) - Brandon Guigo
- create habit test - (d2d6413) - Brandon Guigo
- validator test error - (755bba6) - Brandon Guigo
- jwt + cicd - (08f7672) - Brandon Guigo
- linter - (fef5452) - Brandon Guigo
#### Features
- convert cron to use gRPC call to get user device tokens - (ad76ffe) - Brandon Guigo
- configure auth grpc client into productivity - (92a12df) - Brandon Guigo
- add tests for new repo methods - (d4516db) - Brandon Guigo
- delete folder and time entry too - (89d3492) - Brandon Guigo
- delete tag too - (a139276) - Brandon Guigo
- delete notes too - (6c55701) - Brandon Guigo
- delete habits too - (8012678) - Brandon Guigo
- delete tasks in productivity gRPCserver - (b5b0661) - Brandon Guigo
- configure gRPC server in productivity api - (8aa9aab) - Brandon Guigo
- add rc-microservice - (71aacb7) - Brandon Guigo
- basic config of docker file - (350dc36) - Brandon Guigo
- replace subscription with isSubcribed from token - (697b981) - Brandon Guigo
- add auth to productivity - (d8265ed) - Brandon Guigo
- rename imports and clean files - (a6f6851) - Brandon Guigo

- - -

## productivity-v0.9.0 - 2025-07-17
#### Bug Fixes
- linter - (6d2f946) - Brandon Guigo
- refactor cron to be gRPC ready - (74ac6d2) - Brandon Guigo
- linter - (4a451f2) - Brandon Guigo
- add missing github.com to go modules - (4f38ac7) - Brandon Guigo
- change module name of productivity go module - (798aaa1) - Brandon Guigo
- tag create test - (dec70a9) - Brandon Guigo
- create habit test - (d2d6413) - Brandon Guigo
- validator test error - (755bba6) - Brandon Guigo
- jwt + cicd - (08f7672) - Brandon Guigo
- linter - (fef5452) - Brandon Guigo
#### Features
- convert cron to use gRPC call to get user device tokens - (ad76ffe) - Brandon Guigo
- configure auth grpc client into productivity - (92a12df) - Brandon Guigo
- add tests for new repo methods - (d4516db) - Brandon Guigo
- delete folder and time entry too - (89d3492) - Brandon Guigo
- delete tag too - (a139276) - Brandon Guigo
- delete notes too - (6c55701) - Brandon Guigo
- delete habits too - (8012678) - Brandon Guigo
- delete tasks in productivity gRPCserver - (b5b0661) - Brandon Guigo
- configure gRPC server in productivity api - (8aa9aab) - Brandon Guigo
- add rc-microservice - (71aacb7) - Brandon Guigo
- basic config of docker file - (350dc36) - Brandon Guigo
- replace subscription with isSubcribed from token - (697b981) - Brandon Guigo
- add auth to productivity - (d8265ed) - Brandon Guigo
- rename imports and clean files - (a6f6851) - Brandon Guigo

- - -

## 0.8.1 - 2025-07-11
#### Bug Fixes
- use buildx to build amd and arm images - (24ff873) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** 0.8.0 [skip ci] - (a3bce9f) - GitHub Actions

- - -

## 0.8.0 - 2025-07-10
#### Bug Fixes
- check that user have the right to patch note / task - (749af27) - Brandon Guigo
- linter issues - (a9386ae) - Brandon Guigo
- check for force patch - (ddfeb50) - Brandon Guigo
- skip conflict check if patch have the force boolean - (82a3da5) - Brandon Guigo
- return a valid conflicted item entity - (a6cca10) - Brandon Guigo
- make the patch task test working - (4d29e11) - Brandon Guigo
- most of the tests - (cb5eaf8) - Brandon Guigo
- add tests except create - (b2911fa) - Brandon Guigo
- task creation bug - (3a21bdb) - Brandon Guigo
- support for booleans - (6125f6f) - Brandon Guigo
- make the update date work - (64b9a2b) - Brandon Guigo
- camelCase error - (6167bbd) - Brandon Guigo
- wrong date in backend for patch conflict check - (b921b38) - Brandon Guigo
- remove broken updateBulk to setup patch - (3526c90) - Brandon Guigo
- remove unnecessary logs in cron - (bb9ea1a) - Brandon Guigo
- only update when the task have some changes - (06c2478) - Brandon Guigo
- bulk update request type - (1c4f663) - Brandon Guigo
- rename conflicted item fields - (450ec7b) - Brandon Guigo
- tests and linter - (62f8b25) - Brandon Guigo
#### Features
- add notes patch endpoint - (dd1655d) - Brandon Guigo
- add methods for create / delete - (cb5d922) - Brandon Guigo
- update patch for task - (0378c9d) - Brandon Guigo
- start of the patch method for task - (b7014db) - Brandon Guigo
- add skipped list in bulk update - (fc604eb) - Brandon Guigo
- add bulk update method to task - (c386a82) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** 0.7.0 [skip ci] - (66bae0c) - GitHub Actions

- - -

## 0.7.0 - 2025-06-25
#### Features
- add deleted note boolean - (3c6419f) - Brandon Guigo
- add notes controller - (ae298b9) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** 0.6.3 [skip ci] - (49d1f62) - GitHub Actions

- - -

## 0.6.3 - 2025-06-19
#### Bug Fixes
- make CORS configurable by env var - (d678f72) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** 0.6.2 [skip ci] - (8ef7b36) - GitHub Actions

- - -

## 0.6.2 - 2025-06-19
#### Bug Fixes
- be able to use a env provided full connection string as db creds - (16f01cb) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** 0.6.1 [skip ci] - (5e2cd7a) - GitHub Actions

- - -

## 0.6.1 - 2025-06-19
#### Bug Fixes
- add authMechanism option to mongo - (f6f5e5d) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** 0.6.0 [skip ci] - (4f58c03) - GitHub Actions

- - -

## 0.6.0 - 2025-06-19
#### Bug Fixes
- refactor mongodb service - (9bbb20b) - Brandon Guigo
- linter issues - (30ac4ec) - Brandon Guigo
- add the tests for subscription secure endpoint - (6e52448) - Brandon Guigo
- revenue cat model in user - (9aeca97) - Brandon Guigo
- add revenue cat payload model to parse content from revenuecat - (7305fde) - Brandon Guigo
- linter and naming - (9f79bae) - Brandon Guigo
- change type of duration to String since it's encrypted - (421d7bd) - Brandon Guigo
#### Features
- check subscription when the user add a tag or a habit - (81024f7) - Brandon Guigo
- add subscription utils to check if a user request is needing subscription or not - (d077214) - Brandon Guigo
- store the purchase inside the user object to be ready to dispatch to app at next getUser - (c16c30f) - Brandon Guigo
- setup webhooks controller with static token security - (a8d110d) - Brandon Guigo
- add fields to time_entry - (310ff31) - Brandon Guigo
- move time entries to the dedicated collection - (11f0906) - Brandon Guigo
- add task id to time entry model - (5fa9c5a) - Brandon Guigo
- add controller for time entries without any objects associated - (414e66c) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** 0.5.0 [skip ci] - (642739d) - GitHub Actions
- fix example .env - (7fd5de0) - Brandon Guigo

- - -

## 0.5.0 - 2025-05-23
#### Bug Fixes
- linter - (0801a26) - Brandon Guigo
- add missing folder_id in update - (bafff02) - Brandon Guigo
- update of the folder data + naming for json - (6a7fa9f) - Brandon Guigo
- add the tests for the folder controller - (5e10283) - Brandon Guigo
- linter - (008c95c) - Brandon Guigo
- delete didn't work - (ad5b601) - Brandon Guigo
- UT and add timer and pomo boolean for future usage - (e5dd181) - Brandon Guigo
- update test - (3f240cd) - Brandon Guigo
- upate priority in the repo - (3e967ed) - Brandon Guigo
- add priority field to task entity - (c1e8529) - Brandon Guigo
#### Features
- add folder controller - (e55f806) - Brandon Guigo
- add folders to task entity - (88d21d1) - Brandon Guigo
- add and delete of time entry works - (01da644) - Brandon Guigo
- add the controller endpoints to manage time entries - (4879e27) - Brandon Guigo
- add endpoints to manage time entries on a task - (c678f04) - Brandon Guigo
- add time entries field to task - (209119e) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** 0.4.9 [skip ci] - (5f26e10) - GitHub Actions

- - -

## 0.4.9 - 2025-05-14
#### Bug Fixes
- add missing headers to cors - (7f2b7b7) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** 0.4.8 [skip ci] - (843298c) - GitHub Actions

- - -

## 0.4.8 - 2025-05-14
#### Bug Fixes
- add right image into update step - (a4d691f) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** 0.4.7 [skip ci] - (fda60e3) - GitHub Actions

- - -

## 0.4.7 - 2025-05-14
#### Bug Fixes
- use github actions made for aws to release - (438c284) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** 0.4.6 [skip ci] - (fb640dc) - GitHub Actions

- - -

## 0.4.6 - 2025-05-14
#### Bug Fixes
- update aws ecs service name - (9e16060) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** 0.4.5 [skip ci] - (ac946e3) - GitHub Actions

- - -

## 0.4.5 - 2025-05-14
#### Bug Fixes
- add cors - (e679ed2) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** 0.4.4 [skip ci] - (b8bd9da) - CircleCI
- add cicd for github - (0a3f44a) - Brandon Guigo

- - -

## 0.4.4 - 2025-04-25
#### Bug Fixes
- wrong mongo uri composition - (d44dd37) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** 0.4.3 [skip ci] - (d9cc8cb) - CircleCI

- - -

## 0.4.3 - 2025-04-25
#### Bug Fixes
- retry writes was only set when true, set the param every time it's defined - (7745c3d) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** 0.4.2 [skip ci] - (c1d655d) - CircleCI

- - -

## 0.4.2 - 2025-04-25
#### Bug Fixes
- retry writes parameter badly set - (0eb17b7) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** 0.4.1 [skip ci] - (c493074) - CircleCI

- - -

## 0.4.1 - 2025-04-25
#### Bug Fixes
- use findOrCreate in register endpoint - (bc4a025) - Brandon Guigo
- add logs to the register so we know when data's missing - (094f7f8) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** 0.4.0 [skip ci] - (5712ab1) - CircleCI

- - -

## 0.4.0 - 2025-04-25
#### Bug Fixes
- check that user is owner of the tag for update - (96cc6b8) - Brandon Guigo
- linter issues - (ec38cd1) - Brandon Guigo
- test for get backup key reset - (44de1d6) - Brandon Guigo
- bug with required and false booleans - (be7ec20) - Brandon Guigo
- unauthenticated start reset password in independant table instead of in authenticated user controller - (0d24a0d) - Brandon Guigo
- add tests for user utils - (da5a4fc) - Brandon Guigo
- add tests for start reset password - (2f8264f) - Brandon Guigo
- linter issues - (a9113ae) - Brandon Guigo
- add salt in update password payload - (2165047) - Brandon Guigo
- remove unused keySalt in register - (adc761e) - Brandon Guigo
- linter - (001b482) - Brandon Guigo
- change type of tag in task entity model - (8fe9122) - Brandon Guigo
- delete tag from tasks when tag is deleted - (b56a4aa) - Brandon Guigo
#### Features
- add maizzle email generation to cicd - (e26586e) - Brandon Guigo
- get backup key endpoint + some adjustments - (48f7d9c) - Brandon Guigo
- reset pwd endpoints the right wat, in public auth - (d5cecc7) - Brandon Guigo
- reset user data is mnemonic is lost - (8cd5f52) - Brandon Guigo
- add confirm reset password - (1eeb1c2) - Brandon Guigo
- send email with code via resend - (9b008bb) - Brandon Guigo
- add the templating of the email content from maizzle + store reset code in database - (1bd5f48) - Brandon Guigo
- add maizzle for emails + start reset pwd endpoint - (972d06c) - Brandon Guigo
- add update_password endpoint - (4247832) - Brandon Guigo
- add tags route to main application and update tag model JSON keys - (31f2c96) - Brandon Guigo
- add optional list of tags to a task - (649ccb8) - Brandon Guigo
- add tag controller - (2247d04) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** 0.3.10 [skip ci] - (a907634) - CircleCI
- fix the release changelog content [skip ci] - (5204d22) - Brandon Guigo

- - -

## 0.3.10 - 2025-04-15
#### Bug Fixes
- changelog format - (21c92db) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** 0.3.9 [skip ci] - (ba0453c) - CircleCI

- - -

## 0.3.9 - 2025-04-15
#### Bug Fixes
- last fixes on gh release script - (0d20072) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** 0.3.8 [skip ci] - (1d380aa) - CircleCI

- - -

## 0.3.8 - 2025-04-15
#### Bug Fixes
- script to parse changelog - (9c1923b) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** 0.3.7 [skip ci] - (b342d4d) - CircleCI

- - -

## 0.3.7 - 2025-04-15
#### Bug Fixes
- add logs - (66f29f8) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** 0.3.6 [skip ci] - (87c9e9b) - CircleCI

- - -

## 0.3.6 - 2025-04-15
#### Bug Fixes
- test tag in script - (ca44593) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** 0.3.5 [skip ci] - (3a300ef) - CircleCI

- - -

## 0.3.5 - 2025-04-15
#### Bug Fixes
- remove comments - (fec1e52) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** 0.3.4 [skip ci] - (ae3513e) - CircleCI

- - -

## 0.3.4 - 2025-04-15
#### Bug Fixes
- use valid ghr install command - (ff19e0b) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** 0.3.3 [skip ci] - (655e7d2) - CircleCI

- - -

## 0.3.3 - 2025-04-15
#### Bug Fixes
- boostrap go before doing the release - (9931acc) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** 0.3.2 [skip ci] - (b8c3b95) - CircleCI

- - -

## 0.3.2 - 2025-04-15
#### Bug Fixes
- create github release in cicd - (546c532) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** 0.3.1 [skip ci] - (c96530f) - CircleCI
- add license [skip ci] - (46846e8) - Brandon Guigo

- - -

## 0.3.1 - 2025-04-15
#### Bug Fixes
- unit tests - (451d6cc) - Brandon Guigo
- populate roles in update profile before returning updated user - (ec0e3fb) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** 0.3.0 [skip ci] - (b4c282f) - CircleCI

- - -

## 0.3.0 - 2025-04-10
#### Bug Fixes
- add missing tests - (24de254) - Brandon Guigo
- linter - (0aa3d89) - Brandon Guigo
- retrigger circleci - (0d0b79c) - Brandon Guigo
- models and tests - (cd8702f) - Brandon Guigo
#### Features
- habit reminder notification - (5b05872) - Brandon Guigo
- add days of months field to habit - (725915a) - Brandon Guigo
- add duration into habit entity - (e31ad18) - Brandon Guigo
- implement habit entry management with add, edit, and delete functionalities - (1404384) - Brandon Guigo
- add habit controller and model - (e45fe36) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** 0.2.23 [skip ci] - (2ab3af1) - CircleCI
- replace wrong repo url [skip ci] - (8667388) - Brandon Guigo
- add contirbuting and code of conduct [skip ci] - (d93fcf6) - Brandon Guigo

- - -

## 0.2.23 - 2025-04-02
#### Bug Fixes
- typo in verify deploy - (7d880b5) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** 0.2.22 [skip ci] - (c12579f) - CircleCI

- - -

## 0.2.22 - 2025-04-02
#### Bug Fixes
- remove newline - (770b30d) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** 0.2.21 [skip ci] - (b399a31) - CircleCI

- - -

## 0.2.21 - 2025-04-02
#### Bug Fixes
- tag not propagated between steps, so empty in resource upgrade - (88fdc65) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** 0.2.20 [skip ci] - (43529ec) - CircleCI

- - -

## 0.2.20 - 2025-04-02
#### Bug Fixes
- missing $ in region - (b1cf9dc) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** 0.2.19 [skip ci] - (a7e028a) - CircleCI

- - -

## 0.2.19 - 2025-04-02
#### Bug Fixes
- add region to pipeline - (415c688) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** 0.2.18 [skip ci] - (39ebd2f) - CircleCI

- - -

## 0.2.18 - 2025-04-02
#### Bug Fixes
- deploy to ecs - (7cad098) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** 0.2.17 [skip ci] - (4c0ed32) - CircleCI

- - -

## 0.2.17 - 2025-04-02
#### Bug Fixes
- add missing withProjectId back - (60f9436) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** 0.2.16 [skip ci] - (e7d18bc) - CircleCI

- - -

## 0.2.16 - 2025-04-02
#### Bug Fixes
- add curl to dockerfile - (259b6ee) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** 0.2.15 [skip ci] - (72d034b) - CircleCI

- - -

## 0.2.15 - 2025-04-02
#### Bug Fixes
- add the possibility to set a TLS CA cert for mongo - (dd62ac4) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** 0.2.14 [skip ci] - (7f79b2c) - CircleCI

- - -

## 0.2.14 - 2025-04-01
#### Bug Fixes
- revert fix: mongo ssl in main also - (85a6ad0) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** 0.2.13 [skip ci] - (d4dc346) - CircleCI

- - -

## 0.2.13 - 2025-04-01
#### Bug Fixes
- mongo ssl in main also - (b7a343a) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** 0.2.12 [skip ci] - (2e46156) - CircleCI

- - -

## 0.2.12 - 2025-04-01
#### Bug Fixes
- add support for optional ssl, tls and retryWrites - (bbf4654) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** 0.2.11 [skip ci] - (02f22ed) - CircleCI

- - -

## 0.2.11 - 2025-03-31
#### Bug Fixes
- handling of nil values - (d72c2ee) - Brandon Guigo
- enhance checks for task notifications - (4634412) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** 0.2.10 [skip ci] - (59395ef) - CircleCI

- - -

## 0.2.10 - 2025-03-31
#### Bug Fixes
- dockerfile for CA cert causing API error in k8s - (a4037cb) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** 0.2.9 [skip ci] - (37c04f6) - CircleCI

- - -

## 0.2.9 - 2025-03-31
#### Bug Fixes
- add some logs - (c2f7a51) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** 0.2.8 [skip ci] - (7ad73df) - CircleCI

- - -

## 0.2.8 - 2025-03-31
#### Bug Fixes
- remove project id to exactly match doc - (d73da97) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** 0.2.7 [skip ci] - (63d92fd) - CircleCI

- - -

## 0.2.7 - 2025-03-31
#### Bug Fixes
- test withCredentials - (c98e58a) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** 0.2.6 [skip ci] - (90f5936) - CircleCI

- - -

## 0.2.6 - 2025-03-31
#### Bug Fixes
- revert to original - (bfc695b) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** 0.2.5 [skip ci] - (bbc1ff5) - CircleCI

- - -

## 0.2.5 - 2025-03-31
#### Bug Fixes
- add firebase sac to options - (c885769) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** 0.2.4 [skip ci] - (5e6b36b) - CircleCI

- - -

## 0.2.4 - 2025-03-31
#### Bug Fixes
- test getting the sac directly from env - (d462fae) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** 0.2.3 [skip ci] - (7dc8178) - CircleCI

- - -

## 0.2.3 - 2025-03-31
#### Bug Fixes
- read the google application credentials from disk and pass it to the client as byte array - (f635638) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** 0.2.2 [skip ci] - (464ccc3) - CircleCI

- - -

## 0.2.2 - 2025-03-31
#### Bug Fixes
- add some logs - (9f7f9fb) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** 0.2.1 [skip ci] - (3be0754) - CircleCI

- - -

## 0.2.1 - 2025-03-29
#### Bug Fixes
- remove buggy prints - (4e18ece) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** 0.2.0 [skip ci] - (22e366b) - CircleCI

- - -

## 0.2.0 - 2025-03-28
#### Bug Fixes
- machine change for pr build - (0900a66) - Brandon Guigo
- unit test bis - (999475d) - Brandon Guigo
- UT - (1a60d28) - Brandon Guigo
- add test for regex_utils and comment for linter - (4fdf1c1) - Brandon Guigo
- codeQL injection warning - (05b9ec7) - Brandon Guigo
- linter - (5173592) - Brandon Guigo
- refresh token mock for jwtin the test - (9b4942c) - Brandon Guigo
- populate user roles after device update - (2ad1453) - Brandon Guigo
- linter - (346b615) - Brandon Guigo
- trigger stuck CICD - (25de501) - Brandon Guigo
- change CreatedAt and UpdatedAt fields to use primitive.DateTime - (7486c9b) - Brandon Guigo
- right yq path [skip ci] - (a55825e) - Brandon Guigo
#### Features
- enhance task due notification logging and update FCM multicast message structure - (23709f2) - Brandon Guigo
- add notif payloads + send multicast to user device when task is due - (73644c3) - Brandon Guigo
- add reminders field to TaskEntity and update related tests - (18c53ba) - Brandon Guigo
- update device information handling and add DeviceTimezone field - (66371a7) - Brandon Guigo
- implement cron job for task due notifications and user retrieval - (586d043) - Brandon Guigo
- add device update functionality and related tests - (43ffa2c) - Brandon Guigo
- add user profile update functionality and related tests - (92351f1) - Brandon Guigo
- mnemonic and user salt - (65178f5) - Brandon Guigo
- store the keyset to backup data key with a seed phrase - (1f49747) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** 0.1.14 [skip ci] - (6011c03) - CircleCI
- add codecov in cicd - (3c42a00) - Brandon Guigo

- - -

## 0.1.14 - 2025-03-18
#### Bug Fixes
- manifest file name - (d5bc5ac) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** 0.1.13 [skip ci] - (9c6578f) - CircleCI

- - -

## 0.1.13 - 2025-03-18
#### Bug Fixes
- bad path - (e28cd1a) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** 0.1.12 [skip ci] - (88bcd81) - CircleCI

- - -

## 0.1.12 - 2025-03-18
#### Bug Fixes
- order of the steps in the script - (b174f62) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** 0.1.11 [skip ci] - (e2601ce) - CircleCI

- - -

## 0.1.11 - 2025-03-18
#### Bug Fixes
- try manual clone in script with PAT - (60e7fa6) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** 0.1.10 [skip ci] - (fe21109) - CircleCI

- - -

## 0.1.10 - 2025-03-18
#### Bug Fixes
- try with a netrc - (141c217) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** 0.1.9 [skip ci] - (706aaf6) - CircleCI

- - -

## 0.1.9 - 2025-03-18
#### Bug Fixes
- try to use relative url so auth might works ? - (c530ef8) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** 0.1.8 [skip ci] - (9c1adde) - CircleCI

- - -

## 0.1.8 - 2025-03-18
#### Bug Fixes
- move submodule init - (24e6e90) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** 0.1.7 [skip ci] - (6702307) - CircleCI

- - -

## 0.1.7 - 2025-03-18
#### Bug Fixes
- try new submodule clone - (54436c1) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** 0.1.6 [skip ci] - (728b33b) - CircleCI

- - -

## 0.1.6 - 2025-03-18
#### Bug Fixes
- yq install - (019a317) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** 0.1.5 [skip ci] - (b7cbb2e) - CircleCI

- - -

## 0.1.5 - 2025-03-18
#### Bug Fixes
- try to release new version of API with circleCI - (c605775) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** 0.1.4 [skip ci] - (a9c9377) - CircleCI
- add infra submodule [skip ci] - (667b4f0) - Brandon Guigo

- - -

## 0.1.4 - 2025-03-18
#### Bug Fixes
- add database name to url - (1ae21a5) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** 0.1.3 [skip ci] - (579b81d) - CircleCI

- - -

## 0.1.3 - 2025-03-18
#### Bug Fixes
- add some logs - (8fea7e5) - Brandon Guigo
- parse mongo config from separate env vars - (1abcbd3) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** 0.1.2 [skip ci] - (9f7251c) - CircleCI

- - -

## 0.1.2 - 2025-03-15
#### Bug Fixes
- remove panic the env file is not present + remove tag_prefix - (e2f5f68) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** v0.1.1 [skip ci] - (08ecdc8) - CircleCI
- move the skip ci tag [skip ci] - (16f33da) - Brandon Guigo

- - -

## v0.1.1 - 2025-03-11
#### Bug Fixes
- typo - (ce781a6) - Brandon Guigo
#### Miscellaneous Chores
- **(release)** v0.1.0 [skip ci] - (8cadd52) - CircleCI

- - -

## v0.1.0 - 2025-03-11
#### Bug Fixes
- disable bump commit + push tag + fix docker - (f0107fa) - Brandon Guigo
- cache cocogitto with version - (1bb46c3) - Brandon Guigo
- try to fix git push - (1ed3b37) - Brandon Guigo
- setup git before cog + machine type - (b3d88a3) - Brandon Guigo
- push using PAT - (02eaf8a) - Brandon Guigo
- configure git email and name locally - (55e8b5c) - Brandon Guigo
- add job dependencies + fix cog check - (18ce7c3) - Brandon Guigo
- fix indentation - (ba6445b) - Brandon Guigo
- specify docker small - (c545164) - Brandon Guigo
- resource_class to small - (4c6d351) - Brandon Guigo
- push the latest tag in local git - (de11417) - Brandon Guigo
- rename the github pat var - (3bbedd3) - Brandon Guigo
- add main pipeline - (a435f97) - Brandon Guigo
- remove swagger - (caca2d7) - Brandon Guigo
- swagger main.go path - (d678760) - Brandon Guigo
- build cmd - (cc34693) - Brandon Guigo
- change mongo version in test - (bffd23b) - Brandon Guigo
- linting issues - (d74d1c7) - Brandon Guigo
- change project path - (b9ee01d) - Brandon Guigo
- change go docker version - (eca58e2) - Brandon Guigo
- add runner - (77b8140) - Brandon Guigo
#### Features
- add Dockerfile for multi-stage build - (e1e6eee) - Brandon Guigo
#### Miscellaneous Chores
- **(version)** v0.1.0 - (2efe075) - CircleCI
- **(version)** v0.1.0 - (6d0d65e) - CircleCI
- **(version)** v0.1.0 - (e7685b5) - CircleCI
- **(version)** v0.1.0 - (0b1562f) - CircleCI
- add cocogitto config file - (ab0c398) - Brandon Guigo
- add CircleCI test pipeline - (6a3ed34) - Brandon Guigo
#### Revert
- "fix: specify docker small" - (7f15953) - Brandon Guigo

- - -

## v0.1.0 - 2025-03-11
#### Bug Fixes
- cache cocogitto with version - (1bb46c3) - Brandon Guigo
- try to fix git push - (1ed3b37) - Brandon Guigo
- setup git before cog + machine type - (b3d88a3) - Brandon Guigo
- push using PAT - (02eaf8a) - Brandon Guigo
- configure git email and name locally - (55e8b5c) - Brandon Guigo
- add job dependencies + fix cog check - (18ce7c3) - Brandon Guigo
- fix indentation - (ba6445b) - Brandon Guigo
- specify docker small - (c545164) - Brandon Guigo
- resource_class to small - (4c6d351) - Brandon Guigo
- push the latest tag in local git - (de11417) - Brandon Guigo
- rename the github pat var - (3bbedd3) - Brandon Guigo
- add main pipeline - (a435f97) - Brandon Guigo
- remove swagger - (caca2d7) - Brandon Guigo
- swagger main.go path - (d678760) - Brandon Guigo
- build cmd - (cc34693) - Brandon Guigo
- change mongo version in test - (bffd23b) - Brandon Guigo
- linting issues - (d74d1c7) - Brandon Guigo
- change project path - (b9ee01d) - Brandon Guigo
- change go docker version - (eca58e2) - Brandon Guigo
- add runner - (77b8140) - Brandon Guigo
#### Features
- add Dockerfile for multi-stage build - (e1e6eee) - Brandon Guigo
#### Miscellaneous Chores
- **(version)** v0.1.0 - (6d0d65e) - CircleCI
- **(version)** v0.1.0 - (e7685b5) - CircleCI
- **(version)** v0.1.0 - (0b1562f) - CircleCI
- add cocogitto config file - (ab0c398) - Brandon Guigo
- add CircleCI test pipeline - (6a3ed34) - Brandon Guigo
#### Revert
- "fix: specify docker small" - (7f15953) - Brandon Guigo

- - -

## v0.1.0 - 2025-03-11
#### Bug Fixes
- try to fix git push - (1ed3b37) - Brandon Guigo
- setup git before cog + machine type - (b3d88a3) - Brandon Guigo
- push using PAT - (02eaf8a) - Brandon Guigo
- configure git email and name locally - (55e8b5c) - Brandon Guigo
- add job dependencies + fix cog check - (18ce7c3) - Brandon Guigo
- fix indentation - (ba6445b) - Brandon Guigo
- specify docker small - (c545164) - Brandon Guigo
- resource_class to small - (4c6d351) - Brandon Guigo
- push the latest tag in local git - (de11417) - Brandon Guigo
- rename the github pat var - (3bbedd3) - Brandon Guigo
- add main pipeline - (a435f97) - Brandon Guigo
- remove swagger - (caca2d7) - Brandon Guigo
- swagger main.go path - (d678760) - Brandon Guigo
- build cmd - (cc34693) - Brandon Guigo
- change mongo version in test - (bffd23b) - Brandon Guigo
- linting issues - (d74d1c7) - Brandon Guigo
- change project path - (b9ee01d) - Brandon Guigo
- change go docker version - (eca58e2) - Brandon Guigo
- add runner - (77b8140) - Brandon Guigo
#### Features
- add Dockerfile for multi-stage build - (e1e6eee) - Brandon Guigo
#### Miscellaneous Chores
- **(version)** v0.1.0 - (e7685b5) - CircleCI
- **(version)** v0.1.0 - (0b1562f) - CircleCI
- add cocogitto config file - (ab0c398) - Brandon Guigo
- add CircleCI test pipeline - (6a3ed34) - Brandon Guigo
#### Revert
- "fix: specify docker small" - (7f15953) - Brandon Guigo

- - -

## v0.1.0 - 2025-03-11
#### Bug Fixes
- try to fix git push - (1ed3b37) - Brandon Guigo
- setup git before cog + machine type - (b3d88a3) - Brandon Guigo
- push using PAT - (02eaf8a) - Brandon Guigo
- configure git email and name locally - (55e8b5c) - Brandon Guigo
- add job dependencies + fix cog check - (18ce7c3) - Brandon Guigo
- fix indentation - (ba6445b) - Brandon Guigo
- specify docker small - (c545164) - Brandon Guigo
- resource_class to small - (4c6d351) - Brandon Guigo
- push the latest tag in local git - (de11417) - Brandon Guigo
- rename the github pat var - (3bbedd3) - Brandon Guigo
- add main pipeline - (a435f97) - Brandon Guigo
- remove swagger - (caca2d7) - Brandon Guigo
- swagger main.go path - (d678760) - Brandon Guigo
- build cmd - (cc34693) - Brandon Guigo
- change mongo version in test - (bffd23b) - Brandon Guigo
- linting issues - (d74d1c7) - Brandon Guigo
- change project path - (b9ee01d) - Brandon Guigo
- change go docker version - (eca58e2) - Brandon Guigo
- add runner - (77b8140) - Brandon Guigo
#### Features
- add Dockerfile for multi-stage build - (e1e6eee) - Brandon Guigo
#### Miscellaneous Chores
- **(version)** v0.1.0 - (0b1562f) - CircleCI
- add cocogitto config file - (ab0c398) - Brandon Guigo
- add CircleCI test pipeline - (6a3ed34) - Brandon Guigo
#### Revert
- "fix: specify docker small" - (7f15953) - Brandon Guigo

- - -

## v0.1.0 - 2025-03-11
#### Bug Fixes
- try to fix git push - (1ed3b37) - Brandon Guigo
- setup git before cog + machine type - (b3d88a3) - Brandon Guigo
- push using PAT - (02eaf8a) - Brandon Guigo
- configure git email and name locally - (55e8b5c) - Brandon Guigo
- add job dependencies + fix cog check - (18ce7c3) - Brandon Guigo
- fix indentation - (ba6445b) - Brandon Guigo
- specify docker small - (c545164) - Brandon Guigo
- resource_class to small - (4c6d351) - Brandon Guigo
- push the latest tag in local git - (de11417) - Brandon Guigo
- rename the github pat var - (3bbedd3) - Brandon Guigo
- add main pipeline - (a435f97) - Brandon Guigo
- remove swagger - (caca2d7) - Brandon Guigo
- swagger main.go path - (d678760) - Brandon Guigo
- build cmd - (cc34693) - Brandon Guigo
- change mongo version in test - (bffd23b) - Brandon Guigo
- linting issues - (d74d1c7) - Brandon Guigo
- change project path - (b9ee01d) - Brandon Guigo
- change go docker version - (eca58e2) - Brandon Guigo
- add runner - (77b8140) - Brandon Guigo
#### Features
- add Dockerfile for multi-stage build - (e1e6eee) - Brandon Guigo
#### Miscellaneous Chores
- add cocogitto config file - (ab0c398) - Brandon Guigo
- add CircleCI test pipeline - (6a3ed34) - Brandon Guigo
#### Revert
- "fix: specify docker small" - (7f15953) - Brandon Guigo

- - -

Changelog generated by [cocogitto](https://github.com/cocogitto/cocogitto).