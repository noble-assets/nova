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

  novad keys add validator --home .nova --keyring-backend test &> /dev/null
  novad genesis add-genesis-account validator 1000000ustake --home .nova --keyring-backend test
  novad keys add authority --recover --home .nova --keyring-backend test <<< "deny enhance title bind tunnel drill zebra daring hurt hedgehog outer suspect please suffer cinnamon able relief hen collect female capital jaguar page stand" &> /dev/null
  novad genesis add-genesis-account authority 10000000uusdn --home .nova --keyring-backend test

  TEMP=.nova/genesis.json
  jq '.app_state.staking.params.bond_denom = "ustake"' .nova/config/genesis.json > $TEMP && mv $TEMP .nova/config/genesis.json
  jq '.consensus.params.abci.vote_extensions_enable_height = "5"' .nova/config/genesis.json > $TEMP && mv $TEMP .nova/config/genesis.json

  novad genesis gentx validator 1000000ustake --chain-id "nova-1" --home .nova --keyring-backend test &> /dev/null
  novad genesis collect-gentxs --home .nova &> /dev/null

  sed -i '' 's/timeout_commit = "5s"/timeout_commit = "1s"/g' .nova/config/config.toml
fi

novad start --home .nova --log_level "*:warn,nova:trace"
