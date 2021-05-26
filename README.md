# gClient
Generic Pub/Sub client for 0hub - it's a simulator 🤖

## Motivation

Recently I published one fast, light-weight pub/sub system `pub0sub` --- powered by async I/O. I'm interested in collecting some statistics of its performance in real-world, so I'm writing these simulators & visualisation tool.

I plan to collect data over a timespan & visualise the result to gain deeper insight into performance from end-user's perspective.

You can also use it for `"how to use pub0sub ?"`

## Prerequisite

This simulation depends on `0hub` - a pub/sub server system built using `pub0sub` library. We need to set it up first.

> You might want to check [`0hub`](https://github.com/itzmeanjan/pub0sub#hub)

First clone `pub0sub` into your machine

```bash
cd
git clone https://github.com/itzmeanjan/pub0sub.git
```

Now build & run `0hub`

```bash
cd pub0sub
make hub # with default config
```

But you may want to run `0hub` with different config

```bash
# you're in `pub0sub`

make build_hub
./0hub -help # do check it
```

And probably run `0hub` with

```bash
./0hub -addr 127.0.0.1 -port 13000 -capacity 8192
```

**The server is up & running now.**

Using this simulator you're well capable of running lots of publishers/ subscribers on same machine **( or differente machines )** --- each of them will connect to `0hub` using its own socket connection & use same through out its lifetime.

You'll require to increase max open file descriptor count **( system wide & per process )**, or you'll see `too many open files`

> **Disclaimer :** I only ran this simulation on MacOS & GNU/Linux.

## Usage

There're 3 components

- Simulator
    - [Publisher](#publisher)
    - [Subscriber](#subscriber)
- Visualiser

We'll go through each of them

### Publisher

`gclient-pub` -- a binary read for use

```bash
make build_pub
```

Check help **( definitely )**

```bash
./gclient-pub -help
```

May be run with/ anyhow you please

```bash
./gclient-pub -addr 127.0.0.1 -port 13000 -topic a -topic b -topic c -repeat 0 -client 64 -delay 500ms
```

## Subscriber

Another binary ready to use -- `gclient-sub`

```bash
make build_sub
```

Check help **( definitely )**

```bash
./gclient-sub -help
```

You might want to run it with


```bash
./gclient-sub -addr 127.0.0.1 -port 13000 -topic a -topic b -topic c -client 64
```
