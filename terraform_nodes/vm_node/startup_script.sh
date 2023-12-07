#! /bin/bash
useradd -m -s /bin/bash masa
adduser masa sudo
echo "masa ALL=(ALL:ALL) NOPASSWD: ALL" >> /etc/sudoers
source /etc/sudoers
cd /home/masa
curl -OL https://go.dev/dl/go1.21.4.linux-amd64.tar.gz
sha256sum go1.21.4.linux-amd64.tar.gz
sudo tar -C /usr/local -xvf go1.21.4.linux-amd64.tar.gz
echo "export PATH=$PATH:/usr/local/go/bin" >> /home/masa/.profile
source /home/masa/.profile

cd /home/masa/ && sudo curl -O -L https://github.com/actions/runner/releases/download/v2.311.0/actions-runner-linux-x64-2.311.0.tar.gz
sudo apt install -y jq
token=$(curl -L \
    -X POST \
    -H "Accept: application/vnd.github+json" \
    -H "Authorization: Bearer ${pat_token}"\
    -H "X-GitHub-Api-Version: 2022-11-28" \
    https://api.github.com/repos/${github_name}/${repo_name}/actions/runners/registration-token | jq -r .token)
mkdir /home/masa/actions-runner && cd /home/masa/actions-runner

#remove runner
cat <<EOF > removal_script.sh
#!/bin/bash
removal_token=\$(curl -L \
  -X POST \
  -H "Accept: application/vnd.github+json" \
  -H "Authorization: Bearer ${pat_token}" \
  -H "X-GitHub-Api-Version: 2022-11-28" \
  "https://api.github.com/repos/${github_name}/${repo_name}/actions/runners/remove-token" | jq -r .token)
sudo /home/masa/actions-runner/svc.sh uninstall
/home/masa/actions-runner/config.sh remove --token \$removal_token
EOF
chmod +x removal_script.sh

#install runner
cd /home/masa/actions-runner
tar xzf /home/masa/actions-runner-linux-x64-2.311.0.tar.gz -C /home/masa/actions-runner
sudo chown -R masa:masa /home/masa/actions-runner
sudo -u masa ./config.sh --url https://github.com/${github_name}/${repo_name} --token $token --name ${name} --labels ${name} --runnergroup default --work _work
sudo -u masa sudo /home/masa/actions-runner/svc.sh install
sudo -u masa sudo /home/masa/actions-runner/svc.sh start