# banyan

### 打包程序

GOOS=linux GOARCH=amd64 go build -o bin/banyan_linux

### 配置信息

```azure
# 图片目录
mkdir bucket/i
# 视频目录
mkdir bucket/v 
# 音频目录
mkdir bucket/a
# 文件目录
mkdir bucket/f
```

### centos 部署

touch /etc/systemd/system/banyan.service

```azure
[Unit]
Description=Banyan Golang Application
After=network.target

[Service]
ExecStart=/path/to/deployment/directory/banyan_app
WorkingDirectory=/path/to/deployment/directory
Restart=always

[Install]
WantedBy=multi-user.target
```

启动服务
sudo systemctl start banyan
设置开机启动
sudo systemctl enable banyan

查看状态
sudo systemctl status banyan

刷新服务配置
sudo systemctl daemon-reload

查看日志
journalctl -u banyan -f