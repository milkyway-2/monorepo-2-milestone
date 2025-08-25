echo "=== Ethereum Address (raw + checksum) ==="
if ! ETH_INFO=$(
  PUB_UNCOMP_HEX="$PUB_UNCOMP_HEX" python3 - <<'PY'
import os, binascii, sys
from Crypto.Hash import keccak  # requires: pip install pycryptodome

pub_hex = os.environ.get("PUB_UNCOMP_HEX","").strip()
if not pub_hex:
    print("ERROR: Missing PUB_UNCOMP_HEX", file=sys.stderr); sys.exit(1)
if not pub_hex.startswith("04"):
    print("ERROR: Public key is not uncompressed (missing 04 prefix)", file=sys.stderr); sys.exit(1)

# drop 0x04 and keccak256 the raw X||Y
xy_hex = pub_hex[2:]
data = binascii.unhexlify(xy_hex)

k = keccak.new(digest_bits=256); k.update(data)
digest = k.digest()
addr_bytes = digest[-20:]
addr_hex = "0x" + binascii.hexlify(addr_bytes).decode()

def keccak256_bytes(b: bytes) -> bytes:
    k2 = keccak.new(digest_bits=256); k2.update(b); return k2.digest()

# EIP-55 checksum
def to_checksum_address(addr: str) -> str:
    s = addr.lower().replace("0x","")
    kh = keccak256_bytes(s.encode()).hex()
    out = "0x"
    for i, ch in enumerate(s):
        out += ch.upper() if ch.isalpha() and int(kh[i], 16) >= 8 else ch
    return out

print(addr_hex)
print(to_checksum_address(addr_hex))
PY
); then
  echo "Failed to compute address. See error above."
  exit 1
fi

ETH_RAW=$(echo "$ETH_INFO" | sed -n '1p')
ETH_CHECKSUM=$(echo "$ETH_INFO" | sed -n '2p')
echo "Original:    $ETH_RAW"
echo "Checksummed: $ETH_CHECKSUM"
echo ""
