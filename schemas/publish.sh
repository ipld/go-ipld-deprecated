#!/bin/sh
hash=$(ipfs add -r -q . | tail -n1)
echo "$hash" >>versions
echo "published $hash"
