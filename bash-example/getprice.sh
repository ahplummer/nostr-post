#!/bin/bash

#PREREQs:
# jq - JSON parser, installed on system.
# curl - HTTP client; most distros have this.
# nak - nostr client; executable that is downloaded and renamed `nak` in current directory

#PRE:
# CMC_API: coinmarketcap API key stored in "cmc_key" in current directory
# NSEC: secret nostr key located in "nsec_key" in current directory

# Easy way to generate new nsec key:
# Generate private key, and save: "nak key generate > private_key"
# Open incognito window, point to [Snort](https://snort.social), and log in with that private key.
# Navigate to "Settings, Export keys", and take note of the npub and nsec, putting nsec into `nsec_key` file.

# Get the script (working) directory. This becomes important if you do cron jobs, where the working directory will lose context.
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
CMC_API=$(cat "$SCRIPT_DIR"/cmc_key)

# Current Price, with some fancy formatting
number=`curl -s --header "X-CMC_PRO_API_KEY: "$CMC_API"" "https://pro-api.coinmarketcap.com/v3/cryptocurrency/quotes/latest?slug=bitcoin"  | jq -r '.data[0].quote[0].price'`
if (( $(echo "$number >= 1000" | bc -l) )); then
  value=$(echo "$number / 1000" | bc -l)
  formatted=$(printf "%.1fK" "$value")
else
  formatted=$(printf "%.2f" "$number")
fi

# Get the sats per dollar
sats=$(echo "$number / 100000000" | bc -l)
satsperdollar=$(echo "1 / $sats" | bc -l)
satsrounded=$(printf "%.0f\n" "$satsperdollar")

# Current Fear and Greed.
fear=`curl -s --header "X-CMC_PRO_API_KEY: $CMC_API" https://pro-api.coinmarketcap.com/v3/fear-and-greed/latest | jq -r '.data.value'`
classification=`curl -s --header "X-CMC_PRO_API_KEY: $CMC_API" https://pro-api.coinmarketcap.com/v3/fear-and-greed/latest | jq -r '.data.value_classification'`

# Block Height
height=`curl -sSL "https://mempool.space/api/blocks/tip/height"`

# Build message
message=`echo The block height is "$height", and the current price for bitcoin is "$formatted" USD. This means that you can get "$satsrounded" sats for one dollar. The fear/greed index is "$fear": "$classification". Buy some ₿, it may catch on.`

# Send message to nostr relays; please adjust relays accordingly.
NOSTR_SECRET_KEY=$(cat "$SCRIPT_DIR"/nsec_key) bash -c ''"$SCRIPT_DIR"'/nak event --sec "'"$NOSTR_SECRET_KEY"'" -c "'"$message"'" --tag t=bitcoin nos.lol nostr-pub.wellorder.net relay.damus.io http://umbrel.local:4848'