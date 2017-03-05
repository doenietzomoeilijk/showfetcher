# showfetcher
A glue layer between showrss.info and Transmission, keeping track of shows and downloads.

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
That's pretty much it!
