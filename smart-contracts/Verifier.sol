// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

contract OracleVerifiedDelegation {
    struct Message {
        string validator_address;
        string nominator_address;
        string msgText;
    }

    Message[] public messages;

    address public immutable oracleAddress;

    event MessageStored(string validator, string nominator, string msgText);

    constructor() {
        // Ethereum address derived from public key:
        // 048bca88da17bef374aaade8b23c6f2820d1e1a07463eeff55b86d7c81d6d4eac023b9920ad0a6e9240c9cc006d15319a00bd15a86c290a2900ffee26fd53c0c4f
        oracleAddress = 0xC4Ba2d2F1b3d93D77211cA9Ef5FF1FDD7cA9389E;
    }

    function submitMessage(
        string memory validator_address,
        string memory nominator_address,
        string memory msgText,
        bytes memory signature
    ) public {
        // Step 1: Check that msg.sender is the nominator
        string memory senderAsHex = toAsciiString(msg.sender);
        require(
            keccak256(bytes(senderAsHex)) == keccak256(bytes(nominator_address)),
            "msg.sender does not match nominator_address"
        );

        // Step 2: Rebuild message hash
        bytes32 messageHash = keccak256(
            abi.encodePacked(validator_address, nominator_address, msgText)
        );
        bytes32 ethSignedMessageHash = toEthSignedMessageHash(messageHash);

        // Step 3: Recover signer from signature
        address recovered = recoverSigner(ethSignedMessageHash, signature);
        require(recovered == oracleAddress, "Signature not from oracle");

        // Step 4: Persist the message
        messages.push(
            Message({
                validator_address: validator_address,
                nominator_address: nominator_address,
                msgText: msgText
            })
        );

        emit MessageStored(validator_address, nominator_address, msgText);
    }

    /// @notice Stores a message without any sender or signature validation
    function submitMessageUnverified(
        string memory validator_address,
        string memory nominator_address,
        string memory msgText
    ) public {
        messages.push(
            Message({
                validator_address: validator_address,
                nominator_address: nominator_address,
                msgText: msgText
            })
        );

        emit MessageStored(validator_address, nominator_address, msgText);
    }

    // Helpers

    function toEthSignedMessageHash(bytes32 hash) internal pure returns (bytes32) {
        return keccak256(
            abi.encodePacked("\x19Ethereum Signed Message:\n32", hash)
        );
    }

    function recoverSigner(bytes32 ethSignedMessageHash, bytes memory signature) internal pure returns (address) {
        require(signature.length == 65, "Invalid signature length");

        bytes32 r;
        bytes32 s;
        uint8 v;

        assembly {
            r := mload(add(signature, 0x20))
            s := mload(add(signature, 0x40))
            v := byte(0, mload(add(signature, 0x60)))
        }

        if (v < 27) v += 27;

        return ecrecover(ethSignedMessageHash, v, r, s);
    }

    function toAsciiString(address x) internal pure returns (string memory) {
        bytes memory s = new bytes(40);
        for (uint i = 0; i < 20; i++) {
            bytes1 b = bytes1(uint8(uint(uint160(x)) / (2**(8*(19 - i)))));
            bytes1 hi = bytes1(uint8(b) / 16);
            bytes1 lo = bytes1(uint8(b) - 16 * uint8(hi));
            s[2*i] = char(hi);
            s[2*i+1] = char(lo);
        }
        return string(s);
    }

    function char(bytes1 b) internal pure returns (bytes1 c) {
        if (uint8(b) < 10) return bytes1(uint8(b) + 48);
        else return bytes1(uint8(b) + 87);
    }

    // Optional getter
    function getMessage(uint index) public view returns (Message memory) {
        return messages[index];
    }

    function getMessageCount() public view returns (uint) {
        return messages.length;
    }
}
