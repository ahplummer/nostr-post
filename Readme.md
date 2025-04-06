# Go setup
* Set up .env file:
```
export CMC_API=<your key>
export NSEC=<your nsec>
```
* Source .env file: `source .env`
* Build: `make build`
* Run: `make run`


# Bash example

## Source API calls

* Set up .env file:
```
export CMC_API=<your key>
```
* Source .env file: `source .env`

* Fear and Greed Index: Get API key first from CoinMarketCap (free):
```
curl --header "X-CMC_PRO_API_KEY: $CMC_API" https://pro-api.coinmarketcap.com/v3/fear-and-greed/latest 
```
* Get Dad joke: 
```
curl -s -H \"Accept: application/json\" https://icanhazdadjoke.com/ 
```
* Get Block Height:
```
curl -sSL "https://mempool.space/api/blocks/tip/height"
```
* Get Price:
```
curl --header "X-CMC_PRO_API_KEY: $CMC_API" "https://pro-api.coinmarketcap.com/v3/cryptocurrency/quotes/latest?slug=bitcoin" | jq -r '.data[0].quote[0].price'
```

## Nak commands
* [Source](https://github.com/fiatjaf/nak?tab=readme-ov-file)
* Generate private key, and save:
```
nak key generate > private_key
```
* Open incognito window, point to [Snort](https://snort.social), and log in with that private key.
* Navigate to "Settings, Export keys", and take note of the npub and nsec, putting nsec into `nsec_key` file.
* Use private nsec key to send a message:
```
NOSTR_SECRET_KEY=$(cat nsec_key) bash -c 'nak event --sec 02 -c "good morning" --tag t=gm nostr-pub.wellorder.net relay.damus.io'
```
* View event while signed in on Snort.social with same nsec.

