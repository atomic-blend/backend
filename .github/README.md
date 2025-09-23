# Self-hosted runner

## Deploy a self-hosted runner (DEBIAN)

1. Create a debian VM
2. create a actions user with a password
```
sudo adduser actions
su actions
```
3. Go through the GitHub Actions Runner registration process under the actions user
4. Go back to root, and install the svc with the actions user : 
```
exit
cd /home/actions/actions-runner/
./svc.sh install actions
```
5. Install buildessentials:
```
sudo apt install build-essential jq 
```
6. Install buf with this : 
```
# replace the version with what you want
curl -sSL "https://github.com/bufbuild/buf/releases/download/v1.55.1/buf-linux-x86_64" -o /usr/local/bin/buf
chmod +x /usr/local/bin/buf
buf --version
```
7. Add actions user to the docker group:
```
sudo usermod -aG docker actions
```
8. Install the GitHub cli:
```
(type -p wget >/dev/null || (sudo apt update && sudo apt install wget -y)) \
	&& sudo mkdir -p -m 755 /etc/apt/keyrings \
	&& out=$(mktemp) && wget -nv -O$out https://cli.github.com/packages/githubcli-archive-keyring.gpg \
	&& cat $out | sudo tee /etc/apt/keyrings/githubcli-archive-keyring.gpg > /dev/null \
	&& sudo chmod go+r /etc/apt/keyrings/githubcli-archive-keyring.gpg \
	&& sudo mkdir -p -m 755 /etc/apt/sources.list.d \
	&& echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/githubcli-archive-keyring.gpg] https://cli.github.com/packages stable main" | sudo tee /etc/apt/sources.list.d/github-cli.list > /dev/null \
	&& sudo apt update \
	&& sudo apt install gh -y
```
8. Start the svc
```
./svc.sh start
```

## Deploy a Windows runner

1. Install flutter for windows and the required Visual Studio installation and packages.
Follow the instructions on how to install on the Flutter website.
Run a `flutter doctor` to ensure that everything is fine.

2. Install rust with rustup.

3. Install Inno Setup

4. Open Powershell as administrator and run
```
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope LocalMachine
```

5. Open the Services app, locate the GitHub Actions Runner Service.
Make it start as boot, and in the Log On section, entre your current account username and password. 
Save and restart the service

## Deploy a MacOS runner

1. Install the latest version of XCode

2. Follow the registration process on GitHub

3. Install the svc:
```
./svc.sh install
```

4. Edit the service
```
vi ~/Library/LaunchAgents/actions.runner.atomic-blend.macbook-brandon.plist
```

5. Change `SessionCreate` to `true`

6. Start the service
```
./svc.sh start
```