#include <bits/stdc++.h>

using namespace std;


vector<string> vendors;

ifstream v_in("vendors.txt");
ifstream c_in("enterprise-numbers");
ofstream fout("../../input/mapping.csv");

int min(int x, int y, int z) { return min(min(x, y), z); }

int editDistDP(string str1, string str2, int m, int n)
{
    // Create a table to store results of subproblems
    int dp[m + 1][n + 1];

    // Fill d[][] in bottom up manner
    for (int i = 0; i <= m; i++) {
        for (int j = 0; j <= n; j++) {
            // If first string is empty, only option is to
            // insert all characters of second string
            if (i == 0)
                dp[i][j] = j; // Min. operations = j

            // If second string is empty, only option is to
            // remove all characters of second string
            else if (j == 0)
                dp[i][j] = i; // Min. operations = i

            // If last characters are same, ignore last char
            // and recur for remaining string
            else if (str1[i - 1] == str2[j - 1])
                dp[i][j] = dp[i - 1][j - 1];

            // If the last character is different, consider
            // all possibilities and find the minimum
            else
                dp[i][j]
                    = 1
                      + min(dp[i][j - 1], // Insert
                            dp[i - 1][j], // Remove
                            dp[i - 1][j - 1]); // Replace
        }
    }

    return dp[m][n];
}

string tolower(string &s){
    for (char &c : s) {
        if (c >= 'A' && c <= 'Z'){
            c += ('a' - 'A');
        }
    }

    return s;
}

string get_company(const string &s){
    string company = s;
    company.erase(0, company.find(",") + 1);
    company.erase(company.find(","), company.size());
    company.erase(0, 1);
    company.pop_back();

    return company;
}

string get_email(const string &s){
    string email = s;
    email.erase(0, email.find(",") + 1);
    email.erase(0, email.find(",") + 1);
    email.erase(0, email.find(",") + 1);
    email.erase(0, 1);
    email.pop_back();

    if (email == "---none---") {
        return "";
    }

    if(email.find("yahoo") != string::npos){
        return "";
    }

    email.erase(0, email.find("&") + 1);
    // email.erase(email.find("."), email.size());
    return email;
}

vector<string> get_combinations (const string &s, char separator){
    vector<string> split;
    vector<string> combinations;

    if (s == "") {
        return combinations;
    }

    string aux = s;
    while (aux.find(separator) != string::npos) {
        split.push_back(aux.substr(0, aux.find(separator)));
        aux.erase(0, aux.find(separator) + 1);
    }

    split.push_back(aux);

    for (int i = 0; i < split.size(); i++) {
        string current = "";
        for (int j = i; j < split.size(); j++) {
            if (current == "")
                current += split[j];
            else 
                current += separator + split[j];
            combinations.push_back(current);
        }
    }

    if(separator == ' '){
        string acronim = "";
        for (int i = 0; i < split.size(); i++) {
            acronim += split[i][0];
        }

        combinations.push_back(acronim);
    }

    return combinations;
}

bool search(const string &s, const vector<string> &v){
    int left = 0;
    int right = v.size() - 1;;

    while (left <= right) {
        int mid = (left + right) / 2;

        if (v[mid] == s) {
            return true;
        }

        if (v[mid] < s) {
            left = mid + 1;
        } else {
            right = mid - 1;
        }
    }

    return false;
}

int main(){
    string aux;
    while(getline(v_in, aux)) {
        vendors.push_back(aux);
    }

    sort(vendors.begin(), vendors.end());

    fout << "\"Company\",\"Vendor\"\n";

    int i = 0;
    time_t timer = time(NULL);
    while (getline(c_in, aux)) {
        i++;
        if(i <= 2){
            continue;
        }



        string company = get_company(aux);
        string email = get_email(aux);

        if(company.find("Juniper") != string::npos) {
            cout << "Juniper Networks, Inc.\n";
        }

        vector<string> combinationsCompany = get_combinations(company, ' ');
        vector<string> combinationsEmail = get_combinations(email, '.');

        
        vector<string> variants = combinationsCompany;
        variants.insert(variants.end(), combinationsEmail.begin(), combinationsEmail.end());


        string match = "NONE";
        int minDist = 1000000;

        for (string &v : variants) {
            if (search(v, vendors)){
                match = v;
                break;
            }
        }

        if(match != "NONE"){
            fout << "\"" << company << "\",\"" << match << "\"\n";
        }

        if(!(i % 1000)) {
            cout << i << '\n';
        }
    }

    fout << "}";


    cout << difftime(time(NULL), timer) << '\n';
    
    return 0;
}