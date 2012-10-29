# Usage
Distru is currently in the testing and development stages, not really intended for showing off as of yet. However, there are some pieces which may be interesting for demos.

### WebUI
When an instance of Distru is running, it serves the primary search page on port 9048. If you are running Distru locally, point your browser at [port 9048 locally](http://localhost:9048). If someone else is hosting it, go to their IP (if it is an IPv6 address, enclose it in square brackets,) and append `:9048` to it. For example, to reach a Distru webui on [example.com](http://example.com), you would visit [example.com:9048](http://example.com:9048).

### JSON Index
Distru serves its more intimate details on port 9049. These can't be accessed directly via the webui. At the moment, the best tool for looking at these is `telnet`. Distru uses a human-readable syntax for requesting and serving indexes, so you can make the `telnet` request directly.

To just get the Distru index owned by a Distru server, perform the command `telnet <hostnameorip> 9049`, then type `distru json` and hit enter. This will immediately flood your terminal with that server's JSON-encoded index. It is recommended that you direct the output of telnet to a file, using `telnet <hostnameorip> 9049 > index.json`. The file `index.json` can then be viewed in any text editor or viewer. 
