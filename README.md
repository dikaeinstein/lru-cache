# lru-cache

[![Build Status](https://travis-ci.com/dikaeinstein/lru-cache.svg?branch=master)](https://travis-ci.com/dikaeinstein/lru-cache)
[![Coverage Status](https://coveralls.io/repos/github/dikaeinstein/lru-cache/badge.svg?branch=master)](https://coveralls.io/github/dikaeinstein/lru-cache?branch=master)

A LRU replacement cache

It comes with an in-memory store that is safe for concurrent use.

You can always roll out your own custom store implementation once it implements the `Store` interface.
