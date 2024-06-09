# TEP-BFT
This is the official code for the paper entitled "Parallel Byzantine Fault Tolerance Consensus Based on Trusted Execution Environments".

# Project Requirements

1. Using the Gin framework to build a web server and initially implement API calls to start FISCO binary files.

2. Retrieve the trusted remote attestation report from /dev in the gramine container and print it.

3. Interact with the authoritative blockchain to save the trusted attestation to the authoritative blockchain.

4. Connect to the FISCO node under this web server and develop business blockchains using gosdk.

    ```
    go fisco question -->
    authoritative blockchain query local block or question target go fisco -->
    authoritative blockchain resp yes or no -->
    go fisco get success resp
    ```

