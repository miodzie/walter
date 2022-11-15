# walter
A highly configurable, multi-tenant, modable IRC/Discord/Slack/Anything bot.

#### Purpose
I've written multiple bots for Slack, IRC, and Discord. Typically, with similar 
features that I want across all those platforms. With walter, I can write a feature 
once, and have all the fun things I want everywhere configured to my liking.

### Getting Started
walter is in an active early stage, with some early alpha stages available. However, 
expect _many_ breaking changes in the future.

#### Installing
`go install github.com/miodzie/walter/cmd/...@latest`

Then start `walter`, it will create a default configuration in `~/.walter/config.toml`

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


### TODOs
- [ ] Write module README.mds.
- [ ] Initial README.md

## License
walter is licensed under Apache 2.0 as found in the [LICENSE file](LICENSE).
