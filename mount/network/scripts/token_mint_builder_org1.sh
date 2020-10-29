#!/bin/bash

HLF_PATH="/srv/test-networks/HLF2"

base="'{\"Args\":[\"issue\", \""
split="\", \""
end="\"]}'"
borrower=$1
tokenID=$2
issueDate=$3
maturityDate=$4
amount=$5
echo sudo ./comm_command_org1.sh -i -p $HLF_PATH -c "$base$borrower$split$tokenID$split$issueDate$split$maturityDate$split$amount$end"
sudo ./comm_command_org1.sh -i -p $HLF_PATH -c "$base$borrower$split$tokenID$split$issueDate$split$maturityDate$split$amount$end"