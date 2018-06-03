# Yolochain

This is the code for my first blockchain application. Brushing up my Golang and also on blockchain concepts.

The tutorial I'm following here is:

* [Basic Blockchain Setup](https://medium.com/@mycoralhealth/code-your-own-blockchain-in-less-than-200-lines-of-go-e296282bcffc)
* [Proof of Work](https://medium.com/@mycoralhealth/part-2-networking-code-your-own-blockchain-in-less-than-200-lines-of-go-17fe1dad46e1)
* [Mining a new block](https://medium.com/@mycoralhealth/code-your-own-blockchain-mining-algorithm-in-go-82c6a71aba1f)
* [Proof of Stake](https://medium.com/@mycoralhealth/code-your-own-proof-of-stake-blockchain-in-go-610cd99aa658)

TODOs:
- [ ] Clustered environment. Currently only 1 master node is capable of mining the nodes and orchestrating the winners for POS.
- [ ] REST API to read the blockchain given a Hash.
- [ ] REST API to write to the blockchain.
- [ ] Add context to the function calls once the blockchain is clustered

# License

[MIT](License)