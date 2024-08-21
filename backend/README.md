# Genomic End-User Flow Simulation

This will execute the end-user flow as defined in the cmd/main.go file, demonstrating the interaction between various services in a simulated Genimoc application.

## Flow Overview

1. User Registration: A new user is registered, generating an Ethereum key pair.
2. User Authentication: The user is authenticated using their Ethereum address.
3. Gene Data Encryption: The user's public key encrypts gene data using the Trusted Execution Environment (TEE) service.
4. Gene Data Signing: The user's private key signs the encrypted gene data using the Trusted Execution Environment (TEE) service.
5. Data Storage: The encrypted data, along with its signature and hash, is securely stored.
6. Signature Verification: The stored signature is verified to ensure the integrity of the data.
7. Risk Score Calculation: A risk score is calculated based on the gene data.
8. Blockchain Upload: The gene data is uploaded to the blockchain.
9. Transaction Confirmation: The transaction is confirmed, minting NFTs and rewarding tokens.
10. Data Retrieval: The user retrieves and decrypts the original gene data.

## How to Run

To build and run the project, use the following commands:

```bash
make build
./genomic-be
