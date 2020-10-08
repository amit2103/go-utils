# go-utils
A utility library for Go , aims to provide easy to use utilities for developers

1.  The concurrency package simply provides some utilities. Check the corresponding tests for usage examples

2.  The rate limiter package now provides a simple counting rate limiter which works without channel and tickers.  This avoids unnecessary ticks 
even when the system is idle.

3.  The cache package provides a simple LRU cache which can be used to cache data.

4.  The collections package contains the heap datastructure which can be used for heap sort and also provides reference heap impl. Will also add
other implementations here.

