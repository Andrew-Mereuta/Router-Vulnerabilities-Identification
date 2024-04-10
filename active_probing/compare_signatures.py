"""
Script used to compare two vendor-specific signature datasets.
"""

def parse_signature_set(signature_set):
    return [line.rstrip().split(",") for (line_index, line) in enumerate(signature_set) if line_index != 0]


signatures_1 = open("signatures/signatures.csv", "r")
signature_set_1 = parse_signature_set(signatures_1)
signatures_2 = open("signatures/signatures_old.csv", "r")
signature_set_2 = parse_signature_set(signatures_2)

number_of_exact_match_signatures = 0
number_of_matches_1_unit_distance = 0
for entry in signature_set_1:
    if any(",".join(s[:-1]).startswith(",".join(entry[:-1])) for s in signature_set_2):
        number_of_exact_match_signatures += 1
        continue

    if any(",".join(s[:-2]).startswith(",".join(entry[:-2])) for s in signature_set_2):
        print(",".join(entry))
        number_of_matches_1_unit_distance += 1

print("Number of exact matches between the two sets:" + str(number_of_exact_match_signatures))
print("Number of matches disregarding udp_size between the two sets:" + str(number_of_matches_1_unit_distance))


