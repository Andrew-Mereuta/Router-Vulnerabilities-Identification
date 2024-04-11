#!/bin/bash
"""
A shell script that runs LFP with the new signatures ips.
"""


# python3 setup.py

# The names of the input CSV files
FILE1="signatures/signatures_default.csv"
FILE2="signatures/signatures_new.csv"

MERGED_FILE="signatures/merged_signatures.csv"

# Check if the output file already exists
if [ -f "$MERGED_FILE" ]; then
    echo "$MERGED_FILE already exists. Removing it."
    rm $MERGED_FILE
fi

# Copy the header from the first file to the merged file
head -n 1 $FILE1 > $MERGED_FILE

# Skip the header from the first file and append its content to the merged file
tail -n +2 $FILE1 >> $MERGED_FILE

# Skip the header from the second file and append its content to the merged file
tail -n +2 $FILE2 >> $MERGED_FILE

echo "Files have been merged into $MERGED_FILE"

yes | cp -rf signatures/merged_signatures.csv ../lfp/signatures.csv

# ./probe.sh output/test_ips.txt output/probes_demo.json
./probe.sh output/unknown_unique_ips.txt output/probes_demo.json

python3 ../lfp/analysis.py output/probes_demo.json