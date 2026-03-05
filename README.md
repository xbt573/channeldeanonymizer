# channeldeanonymizer

Simple Telegram which deletes anonymous messages from channels (more on technical side later).

## Setup
1. Clone repository
```bash
git clone https://github.com/xbt573/channeldeanonymizer
```
2. Build program (`Dockerfile` and `docker-compose.yaml` are available for expert use)
```bash
go build
```
3. Launch
```
TOKEN="12345:dhiusdhfiusdhf" CHAT_ID=-10013376942
```
(arguments also available in command line, more at `channeldeanonymizer --help`)

Please note that channel IDs start with `-100`, so you should prepend it to your channel ID if it's not already there.

## Reason
Telegram allows to show channel post authors, but does not give technical solution to disable posting from 3rd party channels (including self, leading to anonymous post). This bot solves this issue by comparing current poster (available in message as author signature) with current list of administrators (lame, but works).

This bot lacks advanced features in some places, like config file or webhook, because it adds unnecessary complexion (at least for my usecase). If you want to see this functionality — feel free to send pull requests.