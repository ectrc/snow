# Configure Discord OAuth

## Getting your OAuth Credentials

![image](https://github.com/ectrc/snow/assets/13946988/8eeab2bc-6056-4329-8c8a-40c908f399d9)

![image](https://github.com/ectrc/snow/assets/13946988/8bdfb31e-9669-47a5-a76e-80c42e01bb84)

Part of the file `config.ini` should look like this:

```ini
[discord]
; discord id of the bot
id="1234567890..."
; oauth2 client secret
secret="abcdefg..."
; discord bot token
token="OTK...."
```

Replace the values with your own, save and rebuild to apply the changes.

## Setup the bot

Add the correct redirects to your discord application:

![image](https://github.com/ectrc/snow/assets/13946988/73fa37b8-3cc2-4b35-85bc-14e4121c6ebd)

This will be from the `config.ini` file:

```ini
[api]
port=":3000"
host="http://localhost"
```

Make sure to add `/snow/discord` to the end of the redirect url.

## Inviting the bot

Generate an invite link for the bot with the following permissions:

![image](https://github.com/ectrc/snow/assets/13946988/04364150-0a93-42a7-a1b8-743b25f49ee9)
![image](https://github.com/ectrc/snow/assets/13946988/90ad3429-ca22-43f6-b426-4d0a6aa83d7c)

The invite link should look like this:

```url
https://discord.com/api/oauth2/authorize?client_id=CLIENT_ID&permissions=34816&redirect_uri=CALLBACK_URL&scope=bot+applications.commands
```
