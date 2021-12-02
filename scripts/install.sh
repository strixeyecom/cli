#!/bin/sh

# download latest version of strixeyed
curl -s https://api.github.com/repos/strixeyecom/cli/releases/latest \
| grep browser_download_url  \
| grep amd64.tar.gz         \
| cut -d : -f 2,3 \
| tr -d \" \
| tr -d , \
| xargs -n1 wget -qO- \
| tar zxf - -C /usr/bin


chmod +x /usr/bin/strixeye

strixeye agent install --user-api-token=$USER_API_TOKEN --agent-id=$AGENT_ID
systemctl enable strixeyed
systemctl start strixeyed