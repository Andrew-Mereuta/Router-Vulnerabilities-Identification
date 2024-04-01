import re

rgx1 = r"<cpe-23:cpe23-item name="
rgx2 = r"cpe:2.3:[^:]*:([^:]*):[^:]*:[^:]*:[^:]*:[^:]*:[^:]*:[^:]*:[^:]*:[^:]*:[^:]*"

def main():

    vendors = set()

    with open("official-cpe-dictionary_v2.3.xml", "r") as f:
        for i, l in enumerate(f):
            if re.search(rgx1, l):
                ans = re.search(rgx2, l)
                if not ans:
                    print(l)
                    print("!!!!!!!!!!")
                    break

                vendors.add(ans.group(1))
    
            # cpe:2.3:a:zzzcms:zzzphp:2.1.0:*:*:*:*:*:*:*
            # cpe:2.3:a:0kims:snarkjs:0.1.7:*:*:*:*:*:*:*
    
    
    print(len(vendors))

    with open("vendors.txt", "w") as f:
        for val in vendors:
            f.write(val + "\n")
    

if __name__ == "__main__":
    main()