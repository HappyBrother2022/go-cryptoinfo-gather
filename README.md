# go-cryptoinfo-gather
Crypto info gather

## Usage
All options are not mandatory even if you do not use related APIs.

```bash
go run main.go \
-logpath=[LOG_PATH] \
-targetSymbol=[TOKEN_SYMBOL] \
-targetAddr=[TOKEN_CONTRACT_ADDR] \
-cmcApikey=[CMC_API_KEY] \
-coinsuper:accesskey=[COINSUPER_ACCESS_KEY] \
-coinsuper:secretkey=[COINSUPER_SECRET_KEY] \
-kucoin:accesskey=[KUCOIN_ACCESS_KEY] \
-kucoin:secretkey=[KUCOIN_SECRET_KEY]
```
...

## License
MIT