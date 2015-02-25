x509util
========

Go utility code for handling X.509 certificates.

Certificate Formats
===================

This library handles certificates using the PKCS standards. Key formats
supported include PKCS #1 for private keys and PKCS #8 for public keys.

Below is some information to convert keys to these formats.

Converting OpenSSH key formats to PKCS key formats
--------------------------------------------------

### Converting Private Keys from OpenSSH to PKCS #1

To decrypt OpenSSH Private Key to OpenSSL PKCS1 Private Key Format, run
the following command, assuming the standard id_rsa private key file
name

> openssl rsa -in id_rsa -out id_rsa.private.pkcs1

### Converting Public Keys from OpenSSH to PKCS #8

To convert OpenSSH Public Key to OpenSSL PKCS8 Public Key Format,
assuming the standard id_rsa.pub public key file name.

> ssh-keygen -e -m PKCS8 -f id_rsa.pub > id_rsa.public.pkcs8