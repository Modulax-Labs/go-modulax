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

---

## 📖 About Modulax Core
Welcome to the heart of the **Modulax Network**.  
This repository contains the official **Golang implementation** of the Modulax node software.  

Modulax Core is a **high-performance, open-source blockchain client** designed from the ground up to:
- Address the **impending threat of quantum computing**  
- Provide **full EVM compatibility** for developers today  

Here you’ll find the **consensus engine**, **peer-to-peer networking**, and our groundbreaking **Post-Quantum EVM (PQ-EVM)** implementation.

---

## 🚀 Key Features

- 🛡️ **Post-Quantum Security**  
  Quantum-resistant signature schemes integrated directly into the EVM.  

- ⚡ **High-Performance PoS**  
  Lightweight Proof-of-Stake consensus with fast finality & low energy use.  

- 🌐 **EVM Compatibility**  
  Full JSON-RPC support, compatible with existing Ethereum tools & dApps.  

- 🔗 **Modular Networking**  
  Built on **libp2p** for scalability, robustness, and resilience.  

---

## 🛠 Tech Stack

| Component      | Technology |
|----------------|------------|
| **Language**   | Go (Golang) |
| **Networking** | libp2p |
| **Database**   | LevelDB |
| **CLI**        | Cobra |

---

## 🏁 Getting Started

### Prerequisites
- [Go](https://go.dev/dl/) **v1.19+**
- C compiler (e.g. GCC) for LevelDB

### Installation & Running

Clone the repo:
```bash
git clone https://github.com/Modulax-Protocol/go-modulax.git
cd go-modulax
