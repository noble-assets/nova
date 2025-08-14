alias novad=./simapp/build/simd

for arg in "$@"
do
    case $arg in
        -r|--reset)
        rm -rf .nova
        shift
        ;;
    esac
done

if ! [ -f .nova/data/priv_validator_state.json ]; then
  novad init validator --chain-id "nova-1" --home .nova &> /dev/null

  novad config set client chain-id "nova-1" --home .nova
  novad config set client keyring-backend test --home .nova

  novad keys add validator --home .nova --keyring-backend test &> /dev/null
  novad genesis add-genesis-account validator 1000000ustake,1000000uusdn --home .nova --keyring-backend test
  novad keys add authority --recover --home .nova --keyring-backend test <<< "deny enhance title bind tunnel drill zebra daring hurt hedgehog outer suspect please suffer cinnamon able relief hen collect female capital jaguar page stand" &> /dev/null
  novad genesis add-genesis-account authority 9000000uusdn --home .nova --keyring-backend test

  TEMP=.nova/genesis.json
  jq '.app_state.hyperlane.mailboxes = [{"id":"0x68797065726c616e650000000000000000000000000000000000000000000000","owner":"noble1kg0mtdjvaqdnk5hw599na9mq66vs2c8zanes8t","message_sent":0,"message_received":0,"default_ism":"0x726f757465725f69736d00000000000000000000000000ff0000000000000000","default_hook":null,"required_hook":null,"local_domain":1313822273}]' .nova/config/genesis.json > $TEMP && mv $TEMP .nova/config/genesis.json
  jq '.app_state.nova.config.hook_address = "0xb19b36b1456E65E3A6D514D3F715f204BD59f431"' .nova/config/genesis.json > $TEMP && mv $TEMP .nova/config/genesis.json
  jq '.app_state.staking.params.bond_denom = "ustake"' .nova/config/genesis.json > $TEMP && mv $TEMP .nova/config/genesis.json
  jq '.app_state.warp.tokens = [{"id":"0x726f757465725f61707000000000000000000000000000010000000000000000","owner":"noble1kg0mtdjvaqdnk5hw599na9mq66vs2c8zanes8t","token_type":"HYP_TOKEN_TYPE_COLLATERAL","origin_mailbox":"0x68797065726c616e650000000000000000000000000000000000000000000000","origin_denom":"uusdn","collateral_balance":"0","ism_id":null},{"id":"0x726f757465725f61707000000000000000000000000000020000000000000001","owner":"noble1kg0mtdjvaqdnk5hw599na9mq66vs2c8zanes8t","token_type":"HYP_TOKEN_TYPE_SYNTHETIC","origin_mailbox":"0x68797065726c616e650000000000000000000000000000000000000000000000","origin_denom":"anoble","collateral_balance":"0","ism_id":null}]' .nova/config/genesis.json > $TEMP && mv $TEMP .nova/config/genesis.json
  jq '.app_state.warp.remote_routers = [{"token_id":"1","remote_router":{"receiver_domain":31337,"receiver_contract":"0x000000000000000000000000ed1db453c3156ff3155a97ad217b3087d5dc5f6e","gas":"0"}}]' .nova/config/genesis.json > $TEMP && mv $TEMP .nova/config/genesis.json
  jq '.consensus.params.abci.vote_extensions_enable_height = "5"' .nova/config/genesis.json > $TEMP && mv $TEMP .nova/config/genesis.json

  novad genesis gentx validator 1000000ustake --chain-id "nova-1" --home .nova --keyring-backend test &> /dev/null
  novad genesis collect-gentxs --home .nova &> /dev/null

  sed -i '' 's/timeout_commit = "5s"/timeout_commit = "1s"/g' .nova/config/config.toml
fi

novad start --home .nova --log_level "*:warn,nova:trace"
