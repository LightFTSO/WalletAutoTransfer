#! /bin/bash

set -e;

CONFIG_ROOT=$HOME/.lightftso

mkdir -p $CONFIG_ROOT;

CONFIG_FILE=$CONFIG_ROOT/config.env
if [ -f $CONFIG_FILE ]; then
   CONFIRM=$(gum confirm "The configuration file already exists, remove and start again?" --default=false);
   if $CONFIRM; then
        rm $CONFIG_FILE;
    else 
        exit 0;
    fi
fi

gum style --border normal --margin "1" --padding "1 2" --border-foreground 202 "Hello there $(gum style --foreground 222 'logoperf')! We're going to set up the service.
Press enter whenever you're ready";

read;
clear;

touch $CONFIG_FILE;

gum style --foreground 222 "First, we're need the network we're going to work in, Flare is selected by default";
NETWORK=$(gum choose "Flare" "Songbird" "Coston", "Other");

if [ ]

echo "NETWORK=$NETWORK" >> $CONFIG_FILE;

clear;

gum style --foreground 222 "Now, we're going to select a RPC provider for the network, the default choice will work but it might be slow, be rate limited or might fail sometimes for whatever reason.
You can enter the same URL you used when you set up the network in Metamask"

FLARE_DEFAULT_RPCURL=https://flare-api.flare.network/ext/bc/C/rpc
SONGBIRD_DEFAULT_RPCURL=https://songbird-api.flare.network/ext/bc/C/rpc
COSTON_DEFAULT_RPCURL=https://coston-api.flare.network/ext/bc/C/rpc

RPC_URL=null

case $NETWORK in
    Flare)
        RPC_URL=$FLARE_DEFAULT_RPCURL
        NAT=FLR
    ;;
    Songbird)
        RPC_URL=$SONGBIRD_DEFAULT_RPCURL
        NAT=SGB
    ;;
    Coston)
        RPC_URL=$COSTON_DEFAULT_RPCURL
        NAT=CFLR
    ;;
    Other)
        $NETWORK=$(gum input --placeholder "Enter the network name here");
        RPC_URL=$COSTON_DEFAULT_RPCURL
        NAT=$(gum input --placeholder "Enter the network's native token here (e.g. ETH for Ethereum)");
    ;;
esac

RPC_URL=$(gum input --value "$RPC_URL" --placeholder $RPC_URL);

echo "RPC_URL=$RPC_URL" >> $CONFIG_FILE;
echo "NAT=$NAT" >> $CONFIG_FILE;
clear;

## Wallets
SOURCE_URL=https://github.com/LightFTSO/WalletAutoTransfer
gum style --border normal --padding "2 2" "$(gum style --foreground 2 'Great!') Now we need the address and the private key of the origin wallet. 

The private key is needed to sign the transactions to send the received funds to another wallet.
You can check the source code of this program at $SOURCE_URL to check that it isn't doing anything nefarious with it."

gum format "Please write (or copy/paste) the address of the origin wallet (the one you want to move the funds from)"
ORIGIN_WALLET_ADDRESS=$(gum input --placeholder "0xabcd1234...")

gum format "Ok, now write (or copy/paste) the private key of the origin wallet (the one you want to move the funds from). 
If you only have the mnemonic (12 or 24 words from a Metamask wallet for example), please contact me on Twitter or search in the web how to obtain it.
We could use the mnemonic words here but that would be a bit more insecure"
ORIGIN_WALLET_PKEY=$(gum input --placeholder "0x1234abcd...")

echo "ORIGIN_WALLET_ADDRESS=$ORIGIN_WALLET_ADDRESS" >> $CONFIG_FILE;
echo "ORIGIN_WALLET_PKEY=$ORIGIN_WALLET_PKEY" >> $CONFIG_FILE;

clear;

gum style --border normal --padding "2 2" "$(gum style --foreground 21 'Almost done here...')
Now we need the address of the destination wallet (the one you're transfering the funds to)"

DESTINATION_WALLET_ADDRESS=$(gum input --placeholder "0xabcd1234...")

echo "DESTINATION_WALLET_ADDRESS=$DESTINATION_WALLET_ADDRESS" >> $CONFIG_FILE;

clear;

# Telegram notifications
TELEGRAM_BOT_API_KEY=null
TELEGRAM_BOT_CHAT_ID=-1
TELEGRAM_NOTIFICATIONS_ENABLED=0
function enableTelegram(){
    gum style --border normal --padding "2 2" "Please create a new Telegram Bot, just talk to BotFather ($(echo '{{ Bold "https://t.me/botfather" }}' | gum format -t template))
    Once you've created a bot and received your authentication token, please enter it below.
    If you need more help, just go to https://core.telegram.org/bots#6-botfather"

    TELEGRAM_BOT_API_KEY=$(gum input --placeholder "Your Telegram bot's auth token here")
    echo "TELEGRAM_BOT_API_KEY=$TELEGRAM_BOT_API_KEY" >> $CONFIG_FILE;

    gum style --border normal --padding "2 2" "Now we're going to get your chat ID.
    When the message 'Bot ready' appears below, send any message to the bot
    It will reply on Telegram with your chat ID, please enter it below";

    TELEGRAM_BOT_CHAT_ID=$(go run main.go -init-telegram)
    echo "Your chat ID is: $TELEGRAM_BOT_CHAT_ID"
    TELEGRAM_NOTIFICATIONS_ENABLED=1
}
function saveTelegramEnv(){
    echo "TELEGRAM_BOT_API_KEY=$TELEGRAM_BOT_API_KEY" >> $CONFIG_FILE;
    echo "TELEGRAM_BOT_CHAT_ID=$TELEGRAM_BOT_CHAT_ID" >> $CONFIG_FILE;
    echo "TELEGRAM_NOTIFICATIONS_ENABLED=$TELEGRAM_NOTIFICATIONS_ENABLED" >> $CONFIG_FILE;
}
if gum confirm "Would you like to enable $(gum style --foreground 25 'Telegram') notifications?"; then
    enableTelegram;
    saveTelegramEnv;
else
    saveTelegramEnv;
fi
gum style --border normal --padding "2 2" "$(gum style --foreground 2 'Perfect!');

These are the settings you chose:
$(echo '{{ Bold "Network" }}' | gum format -t template): $NETWORK
$(echo '{{ Bold "RPC Url" }}' | gum format -t template): $RPC_URL
$(echo '{{ Bold "Origin wallet address" }}' | gum format -t template): $ORIGIN_WALLET_ADDRESS
$(echo '{{ Bold "Origin wallet private key" }}' | gum format -t template): **************
$(echo '{{ Bold "Destination wallet address:" }}' | gum format -t template): $DESTINATION_WALLET_ADDRESS
$(echo '{{ Bold "Telegram notifications" }}' | gum format -t template): $(if [ $TELEGRAM_NOTIFICATIONS_ENABLED -eq "1" ];then echo Enabled;else echo Disabled;fi)
$(echo '{{ Bold "Telegram Bot Token" }}' | gum format -t template): **************
$(echo '{{ Bold "Chat ID" }}' | gum format -t template): $TELEGRAM_BOT_CHAT_ID";



# install
BASE_PATH=$( cd "$( dirname "${BASH_SOURCE[0]}" )"; cd .. && pwd )
BINARY_FILENAME=flare_auto_transfer
BUILD_DIR=$BASE_PATH/build
USER=$(whoami)
function install(){
    echo;
    echo "Compiling...";
    go mod tidy;
    mkdir -p build && go build -o $BUILD_DIR/$BINARY_FILENAME && chmod +x $BUILD_DIR/$BINARY_FILENAME && sleep 1;
    echo "Installing... Please enter your password if prompted:";
    sudo cp $BUILD_DIR/$BINARY_FILENAME /usr/bin/$BINARY_FILENAME && sleep 1;
    echo "Creating service... Please enter your password if prompted:";
    sudo cp $BUILD_DIR/$BINARY_FILENAME /usr/bin/$BINARY_FILENAME && sleep 1;
    echo "Initializing service... $(if [ $TELEGRAM_NOTIFICATIONS_ENABLED -eq "1" ];then echo "You should receive a Telegram notification when ready";fi)";
    cp crypto_auto_transfer.service crypto_auto_transfer.service.temp;
    sed -i "s@{{BINARY_FILENAME}}@\/usr\/bin\/$BINARY_FILENAME@g" crypto_auto_transfer.service.temp;
    sed -i "s@{{USER}}@$USER@g" crypto_auto_transfer.service.temp;
    sudo cp crypto_auto_transfer.service.temp /etc/systemd/system/crypto_auto_transfer.service;
    rm crypto_auto_transfer.service.temp;
    sudo systemctl daemon-reload && sudo systemctl enable crypto_auto_transfer.service;
    sudo systemctl start crypto_auto_transfer.service;

    gum style --border normal --margin "1" --padding "1 2" --border-foreground 202 "$(gum style --foreground 222 'Finished!')! If everything worked out correctly, you can test the service by sending some funds to it.
    They should be automatically transferred to the address you chose here.
    
    A few notes:
    * This only works with native currencies, e.g. FLR, SGB and CFLR. ERC-20 tokens, NFTs or others will not be transferred.
    * On rare occasions the service might encounter an error (especially if there's a problem communicating with the RPC provider). The service will restart automatically during these times.
    * You can check the status by issuing this command on the terminal: 'sudo systemctl status crypto_auto_transfer.service'
    * And you can check the logs with: sudo journalctl '-xeu crypto_auto_transfer.service -f'";
}
install;
