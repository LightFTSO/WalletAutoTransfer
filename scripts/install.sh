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
    cp crypto_auto_transfer.service crypto_auto_transfer.service.temp;
    sed -i "s@{{BINARY_FILENAME}}@\/usr\/bin\/$BINARY_FILENAME@g" crypto_auto_transfer.service.temp;
    sed -i "s@{{USER}}@$USER@g" crypto_auto_transfer.service.temp;
    sudo cp crypto_auto_transfer.service.temp /etc/systemd/system/crypto_auto_transfer.service;
    rm crypto_auto_transfer.service.temp;
    echo "Initializing service... $(if [ $TELEGRAM_NOTIFICATIONS_ENABLED -eq "1" ];then echo "You should receive a Telegram notification when ready";fi)";
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