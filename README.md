# go-tls-pkcs11
Golang TLS client/server scaffold with PKCS11 card access

Tested with the following hardware:

* **Reader**: ACR1281 1S Dual Reader
* **Card**: Feitian PKI FTCOS/PK-01C EnterSafe

# Card setup for the test:

```
> pkcs15-init --erase-card`
Using reader with a card: ACS ACR1281 1S Dual Reader 00 00

> pkcs15-tool --dump
Using reader with a card: ACS ACR1281 1S Dual Reader 00 00
PKCS#15 Card [(null)]:
        Version        : 0
        Serial number  : XXXX
        Manufacturer ID: entersafe
        Flags          : 

> pkcs15-init --create-pkcs15 --profile pkcs15+onepin --label "Test token 1-6/1-8" --pin 123456 --puk 12345678
Using reader with a card: ACS ACR1281 1S Dual Reader 00 00

> pkcs15-tool --dump                                                                                                
Using reader with a card: ACS ACR1281 1S Dual Reader 00 00
PKCS#15 Card [PaKr test token 1-6/1-8]:
        Version        : 0
        Serial number  : XXX
        Manufacturer ID: EnterSafe
        Last update    : 20201123184743Z
        Flags          : EID compliant

PIN [User PIN]
        Object Flags   : [0x03], private, modifiable
        ID             : 01
        Flags          : [0x32], local, initialized, needs-padding
        Length         : min_len:4, max_len:16, stored_len:16
        Pad char       : 0x00
        Reference      : 1 (0x01)
        Type           : ascii-numeric
        Path           : 3f005015


> pkcs11-tool --keypairgen --key-type rsa:2048 --label "Test RSA key" --id 42 --login
Using slot 1 with a present token (0x4)
Logging in to "PaKr test token 1-... (User PIN)".
Please enter User PIN: 123456
Key pair generated:
Private Key Object; RSA 
  label:      Test RSA key
  ID:         42
  Usage:      decrypt, sign, unwrap
  Access:     none
Public Key Object; RSA 2048 bits
  label:      Test RSA key
  ID:         42
  Usage:      encrypt, verify, wrap
  Access:     none
```

Then use `XCA` to generate the following data:
* Generate 4k RSA key & certificate for CA
* Generate TLS client+server certificate for RSA key on-card and sign it with CA
* Copy all certificates to card

Card should have the following objects:

```
> pkcs11-tool --list-objects -l       
Using slot 1 with a present token (0x4)
Logging in to "PaKr test token 1-... (User PIN)".
Please enter User PIN: 123456
Private Key Object; RSA 
  label:      test
  ID:         42
  Usage:      decrypt, sign, unwrap
  Access:     none
Public Key Object; RSA 2048 bits
  label:      Test RSA key
  ID:         42
  Usage:      encrypt, verify, wrap
  Access:     none
Certificate Object; type = X.509 cert
  label:      test
  subject:    DN: O=test, OU=client, CN=0001
  ID:         42
Certificate Object; type = X.509 cert
  label:      CA
  subject:    DN: O=test, OU=admin, CN=CA
  ID:         e228878a552a4313
Public Key Object; RSA 4096 bits
  label:      CA
  ID:         e228878a552a4313
  Usage:      encrypt, verify
  Access:     local
```

the first 3 objects (all of ID 42) are private/public/certificate of user. The last two are public key and cerificate of CA.
