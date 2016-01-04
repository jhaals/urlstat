# Monitor URLs for state changes

Installation

        ~ go get github.com/jhaals/urlstat

Usage

        ~ urlstat --help
        usage: urlstat [<flags>] <urls>...

        Monitor URLs for state changes

        Flags:
          --help           Show context-sensitive help (also try --help-long and
                           --help-man).
        -d, --diff         Display content diff
        -c, --content      Check for changes in content
        -i, --interval=1s  Interval between checks.
        -t, --timeout=1s   HTTP GET timeout.
          --version        Show application version.

        Args:
        <urls>  URLs to monitor

Example

        ~ urlstat http://127.0.1:8000/index.html --content
        Checking http://127.0.1:8000/index.html with 1s interval. Inital Status 404
        2016/01/04 14:50:28 http://127.0.1:8000/index.html status code changed from 404 to 200
        2016/01/04 14:50:28 http://127.0.1:8000/index.html content changed
        2016/01/04 14:50:44 http://127.0.1:8000/index.html content changed
