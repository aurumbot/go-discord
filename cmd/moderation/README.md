# Moderation                                                            
                                                                                     
The moderation package is a powerful built in set of tools to moderate  
and filter a server. It requires a considerably deep integration        
within the bot to work as efficiently as it does (separate database,    
special handler, etc) and should not be replicated, in fact I feel      
rather dirty about creating it.

## General Moderation Commands

### mod

mod takes appropriate moderation action against users. This is defined
by the order `warn -> mute -> mute -> kick -> kick -> tempban ->
permaban`

    Usage:
    mod <@user> [--reason|-r] [--duration|-t]

    Permissions:
    PermissionManageNicknames OR PermissionMuteMembers OR 
    PermissionKickMembers OR PermissionBanMembers

Each permission is needed to warn, mute, kick, and ban. The command
will not prevent a user who can warn from warning, but will from
muting.

### warn

Warn issues an official warning of user's behavior and logs it.

    Usage:
    warn <@user> [--reason|-r]

    Permissions:
    PermissionManageNicknames

### mute

mute silences a user for a given period of time (default
1 hour) and logs it.

    Usage
    warn <@user> [--reason|-r] [--duration|-t]

    Permissions:
    PermissionMuteMembers

### restore

restore unmutes a user and logs it.

    Usage
    restore <@user> [--reason|-r]

    Permissions:
    PermissionMuteMembers

### kick

kick removes a user from the server and logs it. The user is not
re-invited

    Usage:
    kick <@user> [--reason|-r]

    Permissions:
    PermissionKickMembers

### tempban

tempban bans a user for a period of time (default 1 day) and logs it.
The user is given a 1-use invite which never expires.

    Usage:
    tempban <@user> [--reason|-r] [--duration|-t]

    Permissions:
    PermissionKickMembers

### permaban

permaban bans a user and logs it.

    Usage:
    permaban <@user> [--reason|-r]

    Permissions:
    PermissionBanMembers

### pardon

pardon unbans a user and logs it.

    Usage:
    pardon <@user> [--reason|-r]

    Permissions:
    PermissionBanMembers

## AutoMod

The moderation package comes with builtin automatic functions to
censor problematic phrases and remove spam.

### Censor Messages

The moderation package comes with a builtin censor for inappropriate
chat as defined in the config.json. It comes in two modes: blacklist
and whitelist.

#### Whitelist Censorship

Whitelist censorship removes messages containing inappropriate phrases
separated by a space or punctuation, but not if contained in other
words. For example, if one were to ban the word `foo`, then
    Hello, foo.
would be blocked, but not
    Hello, foobar.

#### Blacklist Censorship

Blacklist censorship is the harsher alternative. It removes all
messages containing inappropriate phrases, regardless of how it may be
separated, except in words defined by a list of acceptable phrases.
For example, if one were to ban the word `foo`, and allow the word
`foobaz`, then
    Hello, foo.
    Hello, foobar.
would be blocked, but not
    Hello, foobaz.

### Anti-spam

Along with censorship, moderation package also comes with anti-spam
systems. This removes messages sent too quickly, containing more than
x line returns (as defined in the config.json), containing "zaglo"
characters, or containing more than x mentions. Violation results in
the offending messages being removed and a mute being issued.

### Command `censor`:

censor is the config command for this subset of moderation tools.

    Usage:
    censor <list|spam> <config options>
    censor <list> [--type|-t <whitelist|blacklist>] [--ban <word1,
                  word2>] [--unban <word1, word2>] [--allow <word1,
                  word2>] [--unallow <word1, word2>]

    censor <spam> [--setmute <time*>]

    *time example: 2m, 10s

    Permissions:
    PermissionManageServer

## Logs

The bot logs all server events to a database 

TODO: Write this.

# EOF
