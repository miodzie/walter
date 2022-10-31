# seras
A highly configurable, multi-tenant, modable IRC/Discord/Slack/Anything bot.

#### Purpose
I've written multiple bots for Slack, IRC, and Discord. Typically, with similar 
features that I want across all those platforms. With seras, I can write a feature 
once, and have all the fun things I want everywhere configured to my liking.

### Getting Started
seras is in an active early stage, with some early alpha stages available. However, 
expect _many_ breaking changes in the future.

#### Installing
`go install github.com/miodzie/seras/cmd/...@latest`

Then start `seras`, it will create a default configuration in `~/.seras/config.toml`

Edit it to your liking, and then start the bot again.

### Features
* **Connection Adapters**
  - [X] IRC
  - [X] Discord
  - [ ] Telegram
  - [ ] Slack
* [X] Multi-Tenant
* [X] Configurable Storage
* [X] Global-level Logging module
* [ ] Interactive CLI w/ bubbletea [WIP]

### Modules

* **Message Logging**
  * [X] Console driver
  * [X] Multi-writer driver
  * [ ] File driver
  * [ ] ElasticSearch


* **RSS Feed Listener**
  - [X] MVP
  - [X] Ignore words.
  - [X] Item.Description character limit on String()
  - [X] Better Message formatting (Impl MessageFormatter)
  - [X] NotificationFormatter for RSS. For IRC, it'd be pretty annoying to have a giant wall of text.

* **sed**
  - [X] MVP


* **dong**
  - [X] MVP
  - [ ] Load initial dongs from a web scrape, prepacked sqlite.db, or CSV.


* **Art bot** (more fun for IRC)
  - [X] MVP
  - [ ] User added art.
  - [ ] Discord support. (Only one user is needed)


* **Word triggers**
   - [ ] Canned replies that trigger on specific words.
   - [ ] Percentage change for trigger and cool down timers. 


* **Help module**
  - [ ] Allow other modules to register their available command lists.


* Legacy:
  - [ ] Point system.


* **Moderating**
  * [ ] Profanity list
  * [ ] Meme Repost bot

### TODOs
- [ ] Write module README.mds.
- [ ] Initial README.md

## License
seras is licensed under Apache 2.0 as found in the [LICENSE file](LICENSE).
