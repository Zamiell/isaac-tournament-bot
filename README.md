isaac-tournament-bot
====================

Description
-----------

This is a [Discord](https://discordapp.com/) bot written in [Go](https://golang.org/) that helps to run video game tournaments and leagues by automatically interacting with [Challonge](http://challonge.com/). It is specifically tailored for [The Binding of Isaac: Afterbirth+](http://store.steampowered.com/app/570660/The_Binding_of_Isaac_Afterbirth/) racing leagues. It stores data in a [MariaDB](https://mariadb.org/) database.

<br />


Features
--------

The bot does many different things. You can get a sense of all of the features by taking a look at the [list of commands](https://github.com/Zamiell/isaac-tournament-bot/blob/master/src/command.go).

<br />



Install
-------

These instructions assume you are running Ubuntu 16.04 LTS. Some adjustment will be needed for Windows installations.

- Install Go:
  - `sudo add-apt-repository ppa:longsleep/golang-backports` (if you don't do this, it will install a version of Go that is very old)
  - `sudo apt update`
  - `sudo apt install golang-go -y`
  - `mkdir "$HOME/go"`
  - `export GOPATH=$HOME/go && echo 'export GOPATH=$HOME/go' >> ~/.profile`
  - `export PATH=$PATH:$GOPATH/bin && echo 'export PATH=$PATH:$GOPATH/bin' >> ~/.profile`
- Install [MariaDB](https://mariadb.org/) and set up a user:
  - `sudo apt install mariadb-server -y`
  - `sudo mysql_secure_installation`
    - Follow the prompts.
  - `sudo mysql -u root -p`
    - `CREATE DATABASE isaac;`
    - `CREATE USER 'isaacuser'@'localhost' IDENTIFIED BY '1234567890';` (change the password to something else)
    - `GRANT ALL PRIVILEGES ON isaac.* to 'isaacuser'@'localhost';`
    - `FLUSH PRIVILEGES;`
- Clone the server:
  - `mkdir -p "$GOPATH/src/github.com/Zamiell"`
  - `cd "$GOPATH/src/github.com/Zamiell/"`
  - `git clone https://github.com/Zamiell/isaac-tournament-bot.git` (or clone a fork, if you are doing development work)
  - `cd isaac-tournament-bot`
- Download and install all of the Go dependencies:
  - `cd src` (this is where all of the Go source code lives)
  - `go get ./...` (it is normal for this to take a very long time)
  - `cd ..`
- Set up environment variables:
  - `cp .env_template .env`
  - `nano .env` (fill in the values)
- Import the database schema:
  - `mysql -uisaacuser -p < install/database_schema.sql`

<br />



Run
---

- `cd "$GOPATH/src/github.com/Zamiell/isaac-tournament-bot"`
- `go run src/*.go`

<br />



Compile / Build
---------------

- `cd "$GOPATH/src/github.com/Zamiell/isaac-tournament-bot/src"`
- `go install`
- `mv "$GOPATH/bin/src" "$GOPATH/bin/isaac-tournament-bot"` (the binary is called `src` by default, since the name of the directory is `src`)

<br />



Install as a service (optional)
-------------------------------

- Install Supervisor:
  - `apt install supervisor`
  - `systemctl enable supervisor` (this is needed due to [a quirk in Ubuntu 16.04](http://unix.stackexchange.com/questions/281774/ubuntu-server-16-04-cannot-get-supervisor-to-start-automatically))
- Copy the configuration files:
  - `cp "/root/isaac-tournament-bot/install/supervisord/supervisord.conf" "/etc/supervisor/supervisord.conf"`
  - `cp "/root/isaac-tournament-bot/install/supervisord/isaac-tournament-bot.conf" "/etc/supervisor/conf.d/isaac-tournament-bot.conf"`
- Start it: `systemctl start supervisor`

Later, to manage the service:

- Start it: `supervisorctl start isaac-tournament-bot`
- Stop it: `supervisorctl stop isaac-tournament-bot`
- Restart it: `supervisorctl restart isaac-tournament-bot`

<br />
