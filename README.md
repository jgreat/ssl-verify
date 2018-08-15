# ssl-verify

Small utility to verify ssl certificate chains, before you install them.

## Installation

Download the appropriate binary for your system and place it in your `${PATH}`.

* [Linux]()
* [OSX]()
* [Windows]()

## Usage

```shell
NAME:
   ssl-verify - Verify your ssl certificates/chain is valid.

USAGE:
   ssl-verify [global options] [arguments...]

VERSION:
   0.0.1

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --ca value        Path to CA cert file.
   --cert value      Path to Certificate file. May include intermediate certificates.
   --key value       Path to Private Key file
   --hostname value  Certificate Common name for verification
   --port value      Https local port (default: "8443")
   --help, -h        show help
   --version, -v     print the version
```

## Example

Point `ssl-verify` at your certificate files and provide it the expected hostname. `ssl-verify` will start a small https webserver on localhost an verify the certificate chain.

* `Certificate Verified` - your certificate files are complete and should be ready for use.
* Various Errors - See [Errors](#errors) list below for troubleshooting.

```shell
ssl-verify --hostname rancher.myorg.org --cert ./tls.crt --key ./tls.key --ca ./cacerts.pem

2018/08/15 14:03:47 Server Started
2018/08/15 14:03:49 200 OK
2018/08/15 14:03:49 Certificate Verified
```

## Certificates

`--cert ./tls.crt` - File must have the certificate and may be followed by any Intermediate and Root certificates provided by your CA.

This file can be constructed by using `cat` to concatenate the `.crt` and `.ca-bundle` files together.

```shell
cat yourdomain.crt yourdomain.ca-bundle > tls.crt
```

You should end up with something like this:

```shell
-----BEGIN CERTIFICATE-----
MIIGY... Site Certificate
-----END CERTIFICATE-----
-----BEGIN CERTIFICATE-----
MIIGC... Intermediate Certificate 2
-----END CERTIFICATE-----
-----BEGIN CERTIFICATE-----
MIIGf... Intermediate Certificate 1
-----END CERTIFICATE-----
-----BEGIN CERTIFICATE-----
MIIGf... Root Certificate
-----END CERTIFICATE-----
```

`--key ./tls.key` - File must have the private key associated with the certificate.

```shell
-----BEGIN PRIVATE KEY-----
MIIEvQI...
-----END PRIVATE KEY-----
```

`--ca ./cacerts.pem` - This file is only required if you are using a Private "untrusted" CA. Include the CA cert here so the client can verify the chain.

```shell
----BEGIN CERTIFICATE-----
MIIGx... CA Certificate
-----END CERTIFICATE-----
```

## Errors

| Errors | Reason |
| --- | --- |
| x509: certificate is valid for xxxxx.example.com, not hostname.example.com | The `--hostname` value doesn't match the Subject Common Name in Cert |
| remote error: tls: bad certificate, x509: certificate signed by unknown authority | Certificate is from an "untrusted" Private CA. Make sure you include the Root Certificate and use the `--ca` option |
| tls: private key does not match public key | Wrong certificate for the Private Key file. |
| x509: certificate signed by unknown authority | Missing root and/or intermediate chain certificates. |