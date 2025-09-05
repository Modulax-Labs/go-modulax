<p align="center">
<img src="https://www.google.com/search?q=https://i.imgur.com/example.png" alt="Modulax Logo" width="150"/>
</p>

<h1 align="center">Modulax - Core Protocol (go-modulax)</h1>

<p align="center">
<strong>Official Golang implementation of the Modulax L1 blockchain, including the Post-Quantum EVM (PQ-EVM), consensus, and networking layers.</strong>
<br />
<br />
<a href="https://modulax.com">Website</a>
·
<a href="https://www.google.com/search?q=https://docs.modulax.com">Documentation</a>
·
<a href="https://www.google.com/search?q=https://twitter.com/ModulaxProtocol">Twitter</a>
·
<a href="https://www.google.com/search?q=https://discord.gg/modulax">Discord</a>
</p>

About This Repository
Modulax Core is the backbone of the entire Modulax network. Written in Golang for maximum performance and concurrency, this repository contains the official implementation of our node software. It is the engine that processes transactions, executes smart contracts within our unique Post-Quantum EVM, maintains the security of the ledger through our Proof-of-Stake consensus, and communicates with other nodes across the globe.

🚀 Key Features
Post-Quantum EVM Implementation: The core innovation, designed for future-proof security while maintaining full compatibility with the standard EVM toolchain.

Proof-of-Stake (PoS) Consensus Engine: An efficient and decentralized consensus mechanism designed for high throughput and energy efficiency.

Full & Archive Node Sync Modes: Run the node in different modes to serve various network needs, from lightweight validation to full historical state archival.

EVM-Compatible JSON-RPC API: A comprehensive API that is fully compatible with Ethereum standards, ensuring seamless integration with existing tools and libraries.

🛠 Tech Stack
Golang: For building a high-performance, concurrent, and reliable distributed system.

PostgreSQL: For advanced indexing of blockchain data for high-performance queries.

gRPC: For efficient, high-performance communication between internal node services.

Libp2p: For the modular peer-to-peer networking stack.

🏁 Getting Started
To get a local copy up and running, please follow these steps.

Prerequisites
Go (version 1.21 or later)

A C compiler (like GCC)

Installation & Running
Clone the repo

git clone [https://github.com/Modulax-Protocol/go-modulax.git](https://github.com/Modulax-Protocol/go-modulax.git)

Navigate to the directory

cd go-modulax

Install dependencies

go mod tidy

Build the binary

go build -o modulax ./cmd/modulax

Run the node

./modulax run

🤝 Contributing
Contributions are what make the open-source community such an amazing place to learn, inspire, and create. Any contributions you make are greatly appreciated.

Please see our CONTRIBUTING.md file for details on our code of conduct and the process for submitting pull requests to us.

📜 License
Distributed under the MIT License. See LICENSE for more information.
