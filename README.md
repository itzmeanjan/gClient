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

Using this simulator you're well capable of running lots of publishers/ subscribers on same machine **( or differente machines, talking over LAN )** --- each of them will connect to `0hub` using its own socket connection & use same through out its lifetime.

You'll require to increase max open file descriptor count **( system wide & per process )**, or you'll see `too many open files`

> **Disclaimer :** I only ran this simulation on MacOS & GNU/Linux.

## Result

Here's result of running this simulation on MacOS x86_64 + 8GB RAM + 8 cores, where 

```
subscriber count = 32
subscriber count = 64
topic = 4

Each publisher publishes each message on each topic
Each subscriber subscribes to each of 4 topics
```

Most messages are queued for publishing in `< 5ms`.

![pub_perf](./sc/pub_out.png)

Each queued message is eventually processed by `0hub` - sending it to destination in async fashion.

![sub_perf](./sc/sub_out.png)

## Usage

There're 3 components

- Simulator
    - [Publisher](#publisher)
    - [Subscriber](#subscriber)
- [Visualiser](#visualiser)

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

Log generation is by-default enabled, which generates some new files in **CWD**. If you see **N**-many new log files, it's due to you've set `client` cli option to **N**.

```bash
find . -name 'log_pub_*.csv' # you may try
```

File content might look like

publisher-sent-at-timestamp `( t1 )` | publisher-completed-queueing-at-timestamp `( t2 )` | topic-name
--- | --- | ---
1622035122511 | 1622035122512 | na
1622035123013 | 1622035123013 | na
1622035123514 | 1622035123514 | na
1622035124017 | 1622035124017 | na

> You can ignore `topic-name` column.

These timestamps will satisfy happens before relation

```
t2 >= t1
```

File content to be consumed by visualiser script. See [below](#visualiser).

---

You can disable log generation

```bash
./gclient-pub -topic a -repeat 0 -out false # enabled by default
```

---

### Subscriber

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

When log generation is enabled **( default )** you'll see some new log files created in **CWD**. There'll will be **N**-many log files generated, where **N = concurrent subscriber count**

```bash
find . -name 'log_sub_*.csv' # you may try
```

They're nothing but append only logs in CSV format, recording

publisher-sent-at-timestamp `( t1 )` | receiver-received-at-timestamp `( t2 )` | topic-name
--- | --- | ---
1622019788183 | 1622019788194 | a
1622019788187 | 1622019788195 | b
1622019788190 | 1622019788199 | c
1622019788192 | 1622019788199 | a

So definitely

```
t2 > t1 # satisfying happens before relation
```

These files to be consumed by visualiser for generating plots demostrating system performance i.e. **how long does it generally take for a message to reach destination from when it was sent ?**

---

You can disable log generation

```bash
./gclient-sub -topic a -out false # enabled by default
```

---

### Visualiser

Running above simulation collects performance log for

- Publisher(s) : **Time taken to publish each message**
- Subscriber(s) : **Message reception delay**

---

**How ?** 

Each message published carries millisecond precision unix timestamp serialised in **8 bytes**, which is when message was actually sent from publisher process.

Using this subscriber can deduce how long did it take for a message to arrive to destination from when it was actually prepared on publisher process.

This is what's logged by subscriber [simulator](#subscriber).

On otherside publisher process logs how long did it take for a message to be queued for publishing. This donotes queueing message doesn't really guarantee it's already delivered, rather it says it's queued for delivery.

**It's async** ⭐️.

> All timestamp in milliseconds.

---

Assuming you've already run simulator, you must have logs generated

```bash
find . -name 'log_*.csv'
```

Now you can stop simualtor & start visualiser tool

```bash
cd visualiser
python3 -m venv venv # for first time
source venv/bin/activate

pip install -r requirements.txt # install dependencies

python main.py -h # do check it
python main.py ./.. 12

deactivate
```

If everything goes well, it'll generate two bar charts depicting performance in terms of system **delay in reception/ latency in publishing** in CWD.


```bash
find . -name '*.png' # file names are self-explanatory
```

Yes, that's it.
