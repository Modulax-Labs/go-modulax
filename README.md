<p align="center">
<img src="https://framerusercontent.com/images/AsXYYgiRRwm1tY0cMjzDtsP6xo.png" alt="Modulax Logo" width="150"/>
</p>

<h1 align="center">Modulax Core (go-modulax)</h1>

<p align="center">
<strong>The future-proof ledger. EVM today. Quantum-resistant tomorrow.</strong>
<br />
<br />
</p>

<div align="center">

[![Go Report Card](https://goreportcard.com/badge/github.com/ethereum/go-ethereum)](https://goreportcard.com/report/github.com/ethereum/go-ethereum)
[![Travis](https://app.travis-ci.com/ethereum/go-ethereum.svg?branch=master)](https://app.travis-ci.com/github/ethereum/go-ethereum)
[![Twitter](https://img.shields.io/twitter/follow/ModulaxOrg)](https://x.com/ModulaxOrg)


</div>

📖 About Modulax Core
Welcome to the heart of the Modulax network. This repository contains the official Golang implementation of the Modulax node software. It is a high-performance, open-source blockchain client designed from the ground up to address the impending threat of quantum computing while providing full EVM compatibility for developers today.

This is where the core logic lives: the consensus engine, the peer-to-peer networking, and our groundbreaking Post-Quantum EVM (PQ-EVM) implementation.

🚀 Key Features
🛡️ Post-Quantum Security: Implements quantum-resistant signature schemes directly within the EVM, providing a future-proof foundation for all smart contracts.

⚡ High-Performance PoS: A lightweight and efficient Proof-of-Stake consensus mechanism designed for fast transaction finality and low energy consumption.

🌐 EVM Compatibility: A fully compatible JSON-RPC API ensures that all existing Ethereum tools, libraries, and dApps work seamlessly with Modulax.

** modular Networking:** Built on libp2p for a robust and flexible peer-to-peer layer that is scalable and resilient.

🛠 Tech Stack
Primary Language: Go (Golang)

Networking: libp2p

Database: LevelDB

CLI: Cobra

🏁 Getting Started
To get a local node running for development or testing, follow these steps.

Prerequisites
Go (version 1.19 or later)

A C compiler (like GCC for LevelDB)

Installation & Running
Clone the repository:
`

git clone [https://github.com/Modulax-Protocol/go-modulax.git](https://github.com/Modulax-Protocol/go-modulax.git)
`

Navigate to the project directory:
`

cd go-modulax
`

Build the binary:
`

go build -o modulax ./cmd/modulax
`

Run your node:

./modulax run

🤝 How to Contribute
Contributions are the lifeblood of any open-source project. We welcome developers, security researchers, and enthusiasts to help us build the future-proof ledger.

Fork the repository.

Create a new branch (git checkout -b feature/AmazingFeature).

Commit your changes (git commit -m 'Add some AmazingFeature').

Push to the branch (git push origin feature/AmazingFeature).

Open a Pull Request.

Please read our CONTRIBUTING.md for more details on our code of conduct and the process for submitting pull requests.

🔗 Join Our Community
Stay up to date with the latest developments and connect with the team.

Website: modulax.org

Twitter: @modulaxorg

Discord: Join our Server

📜 License
Distributed under the MIT License. See LICENSE for more information.
