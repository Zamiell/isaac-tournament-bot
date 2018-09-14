#!/bin/bash

cd "/root/go/src/github.com/Zamiell/isaac-tournament-bot/src"
GOPATH=/root/go /usr/local/go/bin/go install
if [ $? -eq 0 ]; then
        mv "/root/go/bin/src" "/root/go/bin/isaac-tournament-bot"
	supervisorctl restart isaac-tournament-bot
else
	echo "isaac-tournament-bot - Go compilation failed!"
fi
