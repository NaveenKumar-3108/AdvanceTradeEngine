## AdvanceTradeEngine
 high-performance trade engine

## Prerequisites
- **Go:** v1.23+
- **Git**
- **Docker** and **Docker Compose**

## Getting Started

### 1. Clone the Repository

```bash
git clone 
```
### 2. Configure the Database

Update `toml/config.toml` with your kafka&redis credentials: 

```toml
[kafka]
broker = "192.168.99.100:9092"
topic1="orders.v1"
topic2="trades.v1"

[redis]
address  = "192.168.99.100:6379"
password = ""
db       = 0
```
Update `docker-compose.yml` with your kafka & redis credentials: 

### 3. Initialize & Run
```bash
docker-compose up -d
```
### server for api and engine for kafka
```bash
go mod tidy
go run ./cmd/server
go run ./cmd/engine
```
### For load test
```bash
go run ./loadgen/server
```
App will start on: `http://localhost:8080`

## 4. LOG file 
To track execution flow and for debugging purpose, the application generates logs during runtime.
LOCATION : log/logfile19042025.12.20.03.610924298.txt

## Author
Maintained by NaveenKumar A. 

## Sample API Requests 
1. API: To Place order
Method: POST
Route: /api/order
Host: localhost:8080

sample request:
http://localhost:8080/api/order

**Sample Request** (Success):
```json
{
  "pair": "BTC/USD",
  "side": "SELL",
  "price":1440,
  "quantity": 2,
  "type":"LIMIT"
}
```

**Sample Response** (Success):
```json
{
    "order_id": "e2622504-35b9-4991-9757-a43a87353f14",
    "status": "S",
    "msg": "Order placed"
}
```
**Sample Response** (error):
```json
{
    "order_id": "",
    "status": "E",
    "msg": "json: cannot unmarshal number into Go struct field Order.side of type string"
}
```

2. API: To Fetch OrderBook
Method: GET
Route: /api/getOrderBook
Host: localhost:8080

sample request:
http://localhost:8080/api/getOrderBook?pair=BTC/USD


**Sample Response** (Success):
```json
{
    "orderbook": {
        "asks": [
            {
                "id": "0285dd2c-1286-4a80-a282-1ff2849a1871",
                "pair": "BTC/USD",
                "price": 1440,
                "quantity": 2,
                "side": "SELL",
                "timestamp": "0001-01-01T00:00:00Z",
                "type": "LIMIT"
            },
        ],
        "bids": []
    },
    "status": "S",
    "msg": ""
}
```

**Sample Response** (error):
```json
{
    "orderbook": null,
    "status": "E",
    "msg": "Missing pair"
}
```
3. API: To Fetch Candles
Method: GET
Route: /api/getOrderBook
Host: localhost:8080

sample request:
http://localhost:8080/api/getCandles?pair=BTC/USD&interval=1s

Interval:
1s or 1m    

**Sample Response** (Success):
```json
{
    "candles": [
        {
            "close": 20839.577609725944,
            "high": 20839.577609725944,
            "low": 20839.577609725944,
            "open": 20839.577609725944,
            "start": "2025-09-07T02:58:57-07:00",
            "volume": 0.2814473134409916
        }
        ],
    "status": "S",
    "msg": ""
}

```

**Sample Response** (error):
```json
{
    "candles": null,
    "status": "E",
    "msg": "Missing pair or interval"
}
```

4. API: To Fetch Trades
Method: GET
Route:/api/getTrades
Host: localhost:8080

sample request:
http://localhost:8080/api/getTrades?pair=BTC/USD


**Sample Response** (Success):
```json
{
   {
    "status": "S",
    "trades": [
        {
            "buy_order_id": "buy1",
            "pair": "BTC/USD",
            "price": 100,
            "quantity": 1,
            "sell_order_id": "sell1",
            "timestamp": "2025-09-07T11:36:21.6551645-07:00"
        },
    ]
   }
}

```

**Sample Response** (error):
```json
{
    "status": "E",
    "msg": "Missing pair"
}
```

**Benchmark Result**
Benchmark TPS: 321.82 orders/sec (N=1, elapsed=3.1073ms)
goos: windows
goarch: amd64
pkg: AdvanceTradeEngine/processor
cpu: Intel(R) Core(TM) i5-8350U CPU @ 1.70GHz
BenchmarkAddOrderTPS-8   	2025/09/07 11:46:27  Connected to Redis

Benchmark TPS: 348.16 orders/sec (N=100, elapsed=287.2211ms)
2025/09/07 11:46:27  Connected to Redis

Benchmark TPS: 250.25 orders/sec (N=417, elapsed=1.6663201s)
     417	   3996478 ns/op	   78651 B/op	     467 allocs/op
PASS
ok  	AdvanceTradeEngine/processor	2.676s