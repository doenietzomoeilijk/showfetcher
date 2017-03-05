# showfetcher
A glue layer between [showrss.info](http://showrss.info) and [Transmission](https://transmissionbt.com/), keeping track of shows and downloads.

## Usage

Make sure there's a `config.json` file present, looking like this:

```json
{
    "feed_url": "http://showrss.info/user/xxx.rss?magnets=true&namespaces=true&name=clean&quality=null&re=null",
    "incomplete_dir": "/path/to/incomplete/torrents",
    "transmission_rpc_url": "http://your.transmission.host:9091/transmission/rpc",
    "shows": [
        {
            "title": "A Series, Matching the ShowRSS Series Name",
            "search_string": "SearchString is not used at this moment",    
            "location": "/path/to/complete/torrents/for/this/series/"
        }
    ]
}
```

For the feed URL, make sure you enable namespaces and clean names, quality and repacks is up to you.

That's pretty much it!

```bash
go build .
./showfetcher
```
