# golang-redis-server


This is a basic implementation of a Redis server using RESP in Golang. So far it's fairly limited in terms of the subset of Redis commands it supports, but it's growing.

Current supported commands are:

* PING
* SET
* GET
* GETSET
* DEL
* EXISTS
* TIME
* EXPIRE

However, new commands can relatively easily be supported by simply adding a new routing function mapped to the name of the command within the setup of server.go.

This project currently has no external dependencies.
