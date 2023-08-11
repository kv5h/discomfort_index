# discomfort_index

Calculates "Discomfort Index (a.k.a temperature–humidity index (THI))" of the current location
based on IP address.

## Prerequisite

Weather data is obtained via [Weather API | www.weatherapi.com](https://www.weatherapi.com/).

So you need to export API key as below:

```bash
export WEATHERAPI_API_KEY=*******************
```

## Run

Available via:

```
# cf.) Options
http://localhost:18080/di
```

Output example:

```json
{
    "city": "Tokyo",
    "feeling": "Hot a little",
    "humidity": 84,
    "index": 78.6112,
    "temperature": 27
}
```

## Options

Alternatively, you can overwrite some parameter.

```
Usage of ./discomfort_index:
  -apipath string
    	API path (default "/di")
  -ipaddress string
    	IP Address
  -port string
    	Port (default "18080")
```
