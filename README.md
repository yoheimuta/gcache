# gcache

[![wercker status](https://app.wercker.com/status/5b8af1fb600ac896b7de2194aabaca6b/m "wercker status")](https://app.wercker.com/project/bykey/5b8af1fb600ac896b7de2194aabaca6b)

`gcache` is a http server that manages [groupcache](https://github.com/golang/groupcache) for the `Redis database`.

It enables your project built by `except golang` to introduce groupcache via http protocol. When each many process caches in many instances(like application servers) loads indivisualy data from the Redis database, the problems are the `Redis database loads like CPU usage %`  and `wasteful memory(which is the same data) for many processes in many instance`.

(Refering to and modifying a part of groupcache README) Whereas process caches just says "Sorry, cache miss"(ex. because of origin data updated), often resulting in a `thundering herd of Redis database (or whatever) loads` from an unbounded number of clients (which has resulted in several fun outages), groupcache coordinates cache fills such that only one load in one process of an entire replicated set of processes populates the cache, then multiplexes the loaded value to all callers.

groupcache has no cluster membership and member failure detection which are necessary to use it in multiple instances. gcache uses [memberlist](https://github.com/hashicorp/memberlist) to solve this issue.

### Usage

Create the initial gcache cluster.

```sh
$ GROUPCACHE_ADDR=172.20.20.11 gcache
```

Join an existing cluster by specifying at least one known member.

```sh
$ GROUPCACHE_ADDR=172.20.20.12 JOIN_TO=172.20.20.11 gcache
```

Ask for data. Request URL pathname has a protocol to define `[optional prefix]-[arguments number]-[return type]-[command]-[command arguments]`

```sh
curl "http://172.20.20.11:3000/1417475105-4-str-HGET-INFO-1"
// or
curl "http://172.20.20.12:3000/1417475105-4-str-HGET-INFO-1"
```

### Install

```sh
$ git clone https://github.com/yoheimuta/gcache.git && cd gcache
$ make install
```

### Run

```sh
$ .godeps/bin/gcache
```

### Test

```sh
$ make test
// or
$ make test DEBUG=1
```

### Try multiple gcache servers on vagrant

Install gcache according to above instruction, `before` setting up vagrant environment.

Then, create vagrant hosts named n1(node to run redis-server), n2(node to run gcache), and n3(node to run gcache joining to n2).

```sh
$ vagrant up
```

Run gcache in n2.

```sh
vagrant ssh n2
cd /vagrant/
GROUPCACHE_ADDR=172.20.20.11 bin/gcache
```

Run gcache and join to n2 in n3.

```sh
vagrant ssh n3
cd /vagrant/
$ GROUPCACHE_ADDR=172.20.20.12 JOIN_TO=172.20.20.11 bin/gcache
```

Try a http request to n2 or n3 gcache server from n1.

```sh
vagrant ssh n1
curl "http://172.20.20.11:3000/1417475105-4-str-HGET-INFO-1"
// or
curl "http://127.20.20.12:3000/1417475105-4-str-HGET-INFO-1"
```
