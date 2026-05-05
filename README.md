### What is this?

This is the repo of a boot.dev guided project with the goal of reconstructing the http server package from scratch, or at least a very basic version of it.

### Which technologies/methods were used?
- TCP connections being read from to get http requests and written to to send http responses.
- Work with various versions of io.Reader & io.Writer interfaces.
- Just generally lots of bytedata parsing and manipulation.

### Is it finished?
The guided project course is finished, but I will likely tinker with this a little bit more since I think there are ways of making the API/UX of the server package more pleasant and there's a few experiments I might run, but for the most part this should be done. 