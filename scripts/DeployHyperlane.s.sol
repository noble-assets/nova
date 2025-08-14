// SPDX-License-Identifier: Apache-2.0
pragma solidity ^0.8.13;

import {Script, console} from "forge-std/Script.sol";
import {Mailbox} from "@hyperlane/Mailbox.sol";
import {NoopIsm} from "@hyperlane/isms/NoopIsm.sol";
import {MerkleTreeHook} from "@hyperlane/hooks/MerkleTreeHook.sol";
import {HypNative} from "@hyperlane/token/HypNative.sol";
import {TransparentUpgradeableProxy} from "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol";

contract DeployHyperlaneScript is Script {
    // Anvil Default Account #1 - 0x59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d
    address constant PROXY_ADMIN = 0x70997970C51812dc3A010C7d01b50e0d17dc79C8;

    function setUp() public {}

    function run() public {
        vm.startBroadcast();

        Mailbox mailboxImplementation = new Mailbox(uint32(block.chainid));
        TransparentUpgradeableProxy mailboxProxy = new TransparentUpgradeableProxy(
            address(mailboxImplementation), PROXY_ADMIN, ""
        );
        Mailbox mailbox = Mailbox(address(mailboxProxy));
        console.log("Mailbox: %s", address(mailbox));

        MerkleTreeHook hook = new MerkleTreeHook(address(mailbox));
        console.log("MerkleTreeHook: %s", address(hook));
        MerkleTreeHook tempHook = new MerkleTreeHook(address(mailbox));
        mailbox.initialize(msg.sender, address(new NoopIsm()), address(tempHook), address(hook));

        HypNative nativeImplementation = new HypNative(1, address(mailbox));
        TransparentUpgradeableProxy nativeProxy = new TransparentUpgradeableProxy(
            address(nativeImplementation),
            PROXY_ADMIN,
            abi.encodeWithSelector(HypNative.initialize.selector, address(0), address(0), msg.sender)
        );
        HypNative native = HypNative(payable(address(nativeProxy)));
        console.log("HypNative: %s", address(native));

        native.enrollRemoteRouter(
            1313822273, // console.log(parseInt('0x'+Buffer.from('NOVA').toString('hex')))
            0x726f757465725f61707000000000000000000000000000020000000000000001
        );

        vm.stopBroadcast();
    }
}
