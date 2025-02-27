cat <<-EOF > /etc/docker/daemon.json
{
  "registry-mirrors": [
        "https://hub.geekery.cn/",
        "https://docker.1panel.live",
        "https://ghcr.geekery.cn"
        ]
}
EOF
systemctl daemon-reload
systemctl restart docker

docker pull dockerpull.org/consul:latest