# Naka-Kai

Naka-Kai is a discord bot to play song from YouTube. Successor of my [Naka](https://github.com/rl404/naka).

[Naka](https://en.wikipedia.org/wiki/Japanese_cruiser_Naka)'s name is taken from japanese cuiser. Also, [exists](https://kancolle.fandom.com/wiki/Naka) in Kantai Collection games and anime as a fleet's idol. To live up its name, this bot will 'sing' for you like an idol.

## Features

- Play song from youtube url.
- Search song from youtube.
- Song queue system.
  - Pause song.
  - Resume song.
  - Stop song.
  - Next song.
  - Previous song.
  - Skip/jump to specific song in queue.
  - Remove 1 or more song from queue.
  - Delete queue.

## Requirement

- [Discord bot](https://discordpy.readthedocs.io/en/latest/discord.html) and its token
- [Go](https://golang.org/)
- [Youtube API key](https://developers.google.com/youtube/v3/getting-started)
- [MySQL](https://www.mysql.com/) / [PostgreSQL](https://www.postgresql.org/)
- [Docker](https://docker.com) + [Docker compose](https://docs.docker.com/compose/) (optional)

## Steps

1. Git clone this repo.

```bash
git clone github.com/rl404/naka
```

2. Rename `.env.sample` to `.env` and modify according to your configuration.

| Env                             |     Default      | Description                                                           |
| ------------------------------- | :--------------: | --------------------------------------------------------------------- |
| `NAKA_KAI_DISCORD_TOKEN`        |                  | Discord bot token                                                     |
| `NAKA_KAI_DISCORD_DELETE_TIME`  |       `0`        | Interval to auto delete discord message (in seconds). `0` to disable. |
| `NAKA_KAI_DISCORD_QUEUE_LIMIT`  |       `20`       | Song queue count limit                                                |
| `NAKA_KAI_DB_DIALECT`           |     `mysql`      | Database dialect (`mysql`/`postgresql`)                               |
| `NAKA_KAI_DB_ADDRESS`           | `localhost:3306` | Database address                                                      |
| `NAKA_KAI_DB_NAME`              |    `naka-kai`    | Database name                                                         |
| `NAKA_KAI_DB_USER`              |      `root`      | Database username                                                     |
| `NAKA_KAI_DB_PASSWORD`          |                  | Database password                                                     |
| `NAKA_KAI_DB_MAX_CONN_OPEN`     |       `10`       | Database max open connection                                          |
| `NAKA_KAI_DB_MAX_CONN_IDLE`     |       `10`       | Database max idle connection                                          |
| `NAKA_KAI_DB_MAX_CONN_LIFETIME` |       `1m`       | Database max connection lifetime                                      |
| `NAKA_KAI_YOUTUBE_KEY`          |                  | YouTube API key                                                       |

3. Migrate the database.

```bash
make migrate

# or using docker
make docker-migrate
```

4. Run the bot.

```bash
make

# or using docker
make docker
# to stop docker
make docker-stop
```

5. Invite the bot to your server.
6. Join a voice channel.
7. Try `/play https://www.youtube.com/watch?v=dQw4w9WgXcQ` from the bot's slash command.
8. Have fun.

## License

MIT License

Copyright (c) 2024 Axel
