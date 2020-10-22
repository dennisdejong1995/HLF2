#!/bin/bash

HLF_PATH="/srv/test-networks/HLF2"

base="'{\"Args\":[\"issue\", \""
split="\", \""
end="]\}'"
issuer=$1
paperNumber=$2
issueDate=$3
maturityDate=$4
amount=$5
echo "$base""$issuer""$split""$paperNumber""$split""$issueDate""$split""$maturityDate""$split""$amount""$end"
# sudo ./comm_command_org1.sh -i -p "$HLF_PATH" -c "$base""$issuer""$split""$paperNumber""$split""$issueDate""$split""$maturityDate""$split""$amount""$end"
