# DiscordTemplate
This repository contains the skeleton code required to build a new Discord bot in Go. It includes predefined structs for managing commands, loggers, databases, and Discord interactions. This should allow you to customize and build a scalable bot quickly and easily.

## Installation
- Clone this Repository
```
git clone git@github.com:brandonliau/DiscordTemplate.git
```

## Create and Configure a Discord Bot
- Go to the [Discord Developer Portal](https://discord.com/developers/applications)
- Create application with `New Application`
- Enable all settings under `Privileged Gateway Intents` in the `Bot` section

## Inviting the Bot
- Under `OAuth2`, select the `URL Generator`
- Select `bot` in the `scopes` section
- Select `Administrator` in the `Bot Permissions` section
- Go to the `Generated URL` link and select a server to add the bot

## Configuring the Project
- Select `reset token` in the `Bot` section
- Copy the token and paste it into the `config.yml`

## Creating a Daemon Service
This method of deployment is optional. Another popular method of deployment includes containerizing the application with [Docker](https://www.docker.com/) and deploying it on a cloud service provider such as [AWS](https://aws.amazon.com/free) or [GCP](https://cloud.google.com/gcp).

- Create a daemon service
```
sudo nano /etc/systemd/system/{service_name}.service
```
- Copy the following into the service file
```
[Unit]
Description={description}
After=multi-user.target
[Service]
Type=simple
Restart=always
ExecStart=/root/{service_name}/{service_name}
WorkingDirectory=/root/{service_name}
StandardOutput=append:/root/{service_name}/log.log
StandardError=append:/root/{service_name}/log.log
[Install]
WantedBy=multi-user.target
```
- Replace `{description}` with a short description of your service
- Replace `{service_name}` with the name of your service or application
- Paths for `WorkingDirectory`, `StandardOutput`, and `StandardError` may differ depending on your system configuration

## Systemctl Commands
- Reload service configuration
```
sudo systemctl daemon-reload
```
- Enable the service
```
sudo systemctl enable {service_name}.service
```
- Disable the service
```
sudo systemctl disable {service_name}.service
```
- Start the service
```
sudo systemctl start {service_name}.service
```
- Stop the service
```
sudo systemctl stop {service_name}.service
```
- Restart the service
```
sudo systemctl restart {service_name}.service
```
- View service status
```
sudo systemctl status {service_name}.service
```
- For more information on daemons, visit the following
    - https://medium.com/p/f0cc55a42267
    - https://github.com/torfsen/python-systemd-tutorial
