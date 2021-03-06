# popcon-sc
[![Build Status](http://img.shields.io/travis/cs3238-tsuzu/popcon-sc/master.svg?style=flat-square)](https://travis-ci.org/cs3238-tsuzu/popcon-sc)
[![license](https://img.shields.io/github/license/mashape/apistatus.svg?style=flat-square)](./LICENSE)

## What is popcon-sc?
- Programming OPen(source) CONtest Server - SCalable
- popcon-sc is an open-source contest management system for competitive programming.
- Mainly, made to be used in my club activities.
- Main: Pure Go
- Web: Bootstrap3
- I'll make this more useful than the previous project, [popcon](https://github.com/cs3238-tsuzu/popcon) by using Docker.

## Features
- Scalable web server & judging system
- Easy to launch with Golang & Docker
- Support of multiple kinds of contests

## How to install
- Requirements: Docker
- $ cd path/to/somewhere


### For docker-compose
- $ cat > .env
- PP_TOKEN="your password"
- PP_DB_PASSWORD="your password"
- Ctrl-C
- $ wget "https://raw.githubusercontent.com/cs3238-tsuzu/popcon-sc/master/docker-compose.yml"
- $ docker-compose -f docker-compose.yml up -d
- $ docker-compose -f docker-compose.yml logs -f | grep Pass
- When you get admin's password, stop with Ctrl-C
- Access localhost:80 and signin

### For docker-swarm
- **Deprecated**
- This does not work well. I cannot understand why containers cannot get the other containers' true addresses.
- If your computer doesn't join a swarm network, $ docker swarm init
- $ wget "https://raw.githubusercontent.com/cs3238-tsuzu/popcon-sc/master/docker-compose-swarm.yml"
- $ export PP_TOKEN="your password"
- $ export PP_DB_PASSWORD="your password"
- $ docker stack deploy -f docker-compose-swarm.yml popcon
- $ docker service logs popcon_ppweb | grep Pass
- When you get admin's password, stop with Ctrl-C
- Access localhost:80 and signin

### After installation
- Launch ppjudge
- Refer to ppjudge/README.md

## Supported Platforms
- 64bit architecture
- Ubuntu 16.04 LTS
- macOS 10.12(with Docker for Mac)
- Windows 10 Pro (with Docker for Windows)
	- If you fail to install this on windows, read [FAQ](./doc/faq.md#for-windows)

## License
- Under the MIT License
- Copyright (c) 2017 Tsuzu
