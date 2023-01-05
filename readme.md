# My HTTP test server

* Echo HTTP request received from client
* Get geographic information of an IP address

## List APIs

### [/]

Echo HTTP request received from client (start line, header, body).

Example response:

````text
echo req from 127.0.0.1:42428:

POST /?a=5 HTTP/1.1
Host: 127.0.0.1:20891
Accept: */*
Content-Length: 55
Content-Type: application/json
User-Agent: insomnia/6.5.3

{
	"username": "username0",
	"password": "password0"
}
````

### [/ip]

Return caller IP address in response body.

Example response: `127.0.0.1` 

### [/ip/geo]

Return caller IP address geographic information.

Example response:

````json
{
	"IP": "14.187.159.9",
	"Continent": "Asia",
	"ContinentCode": "AS",
	"Country": "Vietnam",
	"CountryCode": "VN",
	"City": "Ho Chi Minh City",
	"TimeZoneName": "Asia/Ho_Chi_Minh",
	"ISPName": "VNPT Corp"
}
````

### [/ip/geo/:ip]

Return geographic information of the IP in URL param.

Example response of [/ip/geo/216.58.221.238]:

````json
{
	"IP": "216.58.221.238",
	"Continent": "North America",
	"ContinentCode": "NA",
	"Country": "United States",
	"CountryCode": "US",
	"City": "Mountain View",
	"TimeZoneName": "America/Los_Angeles",
	"ISPName": "GOOGLE"
}
````

Can input a hostname [/ip/geo/lichess.org]:

````json
{
	"IP": "37.187.205.99",
	"Continent": "Europe",
	"ContinentCode": "EU",
	"Country": "France",
	"CountryCode": "FR",
	"City": "",
	"TimeZoneName": "Europe/Paris",
	"ISPName": "OVH SAS"
}
```` 

## Docker

Execute [build_run.sh](build_run.sh) to run a local container, default
listen on port 20891.  

## References

* MaxMind free [database](https://www.maxmind.com/en/accounts/404644/geoip/downloads).
* [Reader](https://github.com/oschwald/geoip2-golang) for  MaxMind data format