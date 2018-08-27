# Moderation                                                            
                                                                                     
The moderation package is a powerful built in set of tools to moderate  
and filter a server. It requires a considerably deep integration        
within the bot to work as efficiently as it does (separate database,    
special handler, etc) and should not be replicated, in fact I feel      
rather dirty about creating it.

## AutoMod

The moderation package comes with builtin automatic functions to
censor problematic phrases and remove spam.

## Commands

### mod

mod takes appropriate moderation action against users. This is defined
by the order `warn -> mute -> mute -> kick -> kick -> tempban ->
permaban`

  Usage:
  mod <@user> [--reason|-r] [--duration|-t]

  Permissions:


### warn

Warn issues an official warning of user's behavior and logs it. The
warn decays after 24 hours.

  Usage:
  warn <@user> [--reason|-r]

