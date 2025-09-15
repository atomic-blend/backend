# Self-hosted runner

## Deploy a self-hosted runner

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
sudo apt install build-essential
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
8. Start the svc
```
./svc.sh start
```