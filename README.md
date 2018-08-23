# go-discord

Discord Bots for The Rest of Us

## About

go-discord is a golang framework for modular, D-I-Y discord bots that
can be setup and run on a home computer.

## Features

- Easy to setup, easier to deploy.
- Modular design allowing for easy integration of commands into the
  bot
- Highly customizable: make the bot fit specific needs.
- Universal: easy to run at home on any computer.
- Daemonized: tmux need not apply.
- Private: No data is collected as opposed to many big, centralized
  bots.
- Powerful moderation tools built in 
- Permission-based commands already sync up with moderator and
  administrative roles
- Easy to develop API with builtin features like standard logging and
  flag parsing

## Setting up

### For the first time (macOS and linux):

**macOS Users:** note that you must have the xcode terminal tools. If
you don't want to install xcode, simply run
    
    xcode-select --install

1. Go into `cmd/handler/map.go` and import any packages you may want.
   See the comment for more info about importing third-party plugins.
2. Go to [discord's application
   page](http://discordapp.com/developers/applications/) and create
   a new application, select "bot", and configure its name, icon, etc
   - If you don't know what permissions to give the bot, select
     "administrator"
   - Remember to save once you're done!
3. Go back to the bot and run `sudo setup.sh`. Provide it with the
   necessary information. The bot will compile, install itself as
   a daemon, then start.
   - If you don't know what port to use, just use any number over
     2,000 and under 65,000. One port per bot.
   - If you need to stop the bot daemon, use `sudo <bot> stop`, `sudo
     <bot> remove` un-installs the daemon, and `sudo <bot> start`
     starts it again. replace <bot> with the name of your bot.

That's it. If you need to update the bot, un-install the daemon before
hand, make your changes, then run setup.sh again. The same goes if you
want to create another bot.

### For the first time (Windows 10)

Setup your built-in ubuntu emulator, then follow the steps above.

## Issues

If you encounter issues with go-discord or its first-party plugins,
please feel free to file an issue
[here](https://github.com/whitman-colm/go-discord/issues/)

## Sensible API Documentation

I promise this is coming.

# EOF.
