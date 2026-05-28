# AidStream: Technical Design & Workflow Documentation

---

# 1. Executive Overview

**AidStream** is a decentralized humanitarian aid platform designed to restore trust in charitable giving. By leveraging **Blockchain Technology** and an **Escrow-based distribution model**, the platform ensures that donations are not only promised but verifiably delivered to legitimate beneficiaries through a transparent and immutable audit trail.

## Core Value Proposition

### Immutability

Every donation transaction is permanently recorded on the blockchain.

### Accountability

Funds remain in escrow until predefined verification conditions are fulfilled.

### Efficiency

Direct donor-to-beneficiary tracking is enabled through the Go backend and blockchain smart contracts.

---

# 2. System Architecture

AidStream follows a **decoupled architecture**, separating the user interface, business logic layer, database, and decentralized trust layer.

## 2.1 Component Breakdown

| Layer                  | Technology               | Responsibility                                                                         |
| ---------------------- | ------------------------ | -------------------------------------------------------------------------------------- |
| **Presentation Layer** | HTML5 / JavaScript / CSS | Handles user interaction, MetaMask wallet connection, and data visualization           |
| **Logic (API) Layer**  | Go (Golang)              | Manages authentication, API routing, database operations, and blockchain communication |
| **Persistence Layer**  | Neon PostgreSQL          | Stores off-chain data such as user profiles, beneficiary cases, and metadata           |
| **Trust Layer**        | Solidity                 | Executes escrow management and transaction validation through smart contracts          |

---

# 3. Data & Transaction Flow

This section outlines the lifecycle of a **Beneficiary Case** and a **Donation Transaction** within the AidStream ecosystem.

---

# 3.1 Beneficiary Lifecycle

## Step 1: Case Creation

An admin or humanitarian organization creates a beneficiary case through the **Admin Dashboard**.

## Step 2: Off-chain Storage

The Go backend validates the submitted information and stores the case details in the **Neon PostgreSQL** database.

## Step 3: Verification

A secondary verification process confirms the legitimacy of the case before it becomes publicly visible on the `cases.html` page.

---

# 3.2 Donation Lifecycle

## Step 1: Authentication

The donor logs in and connects their MetaMask wallet using `wallet.js`.

## Step 2: Donation Initiation

The donor selects a beneficiary case and specifies the donation amount.

## Step 3: Smart Contract Execution

The `AidStream.sol` contract initiates an escrow transaction. Funds are transferred to the smart contract address instead of an organization’s private wallet.

## Step 4: Backend Synchronization

Once the blockchain transaction hash is generated, the Go backend (`donation_controller.go`) captures the hash and updates donation records in the database.

## Step 5: Transaction Validation

The `blockchain.go` service monitors the blockchain network for confirmation events and updates the donor dashboard in real time.

---

# 4. Module Responsibilities (Team Alignment)

## Frontend — Emma

Responsible for maintaining the frontend responsiveness through `frontend/js/api.js`.

### Key Responsibilities

* Ensure UI reacts dynamically to backend updates
* Display accurate donation metrics
* Differentiate clearly between:

  * **Total Contributed**
  * **Transaction Status**

---

## Backend & API — Edward

Responsible for server-side routing, authentication, and security enforcement.

### Key Responsibilities

* API endpoint management
* Request validation
* Access control through `middleware/auth.go`
* Restrict verification panel access to authorized administrators only

---

## Database Management — Philip

Responsible for maintaining relational integrity within the Neon PostgreSQL database.

### Critical Tables

* **Donations**

  * Linked to blockchain transaction hashes
* **Beneficiaries**

  * Linked to verification statuses

---

## Blockchain Development — Paul

Responsible for smart contract development and escrow logic.

### Key Responsibilities

* Development of:

  * `AidStream.sol`
  * `Escrow.sol`
* Ensure secure release and traceability of funds
* Validate transactions through blockchain event logic

---

## Integration — Austine

Responsible for connecting backend services with blockchain events and smart contracts.

### Key Responsibilities

* Synchronize Go services with blockchain activity
* Ensure donation events are properly logged
* Handle blockchain event listeners and confirmations

---

# 5. Security & Verification Measures

## Wallet Validation

Only transactions signed by the connected MetaMask wallet are accepted and processed.

## Escrow Protection

Funds remain locked in escrow until smart contract conditions are satisfied, reducing the risk of fraud or unauthorized withdrawals.

## Data Integrity

* Sensitive and operational metadata is stored securely in PostgreSQL
* Financial proof and transaction records remain permanently stored on-chain

---

# 6. Project Roadmap

## Phase 1 — Core Infrastructure *(Current Phase)*

* Implementation of front end and backend side of the server 

## Phase 2 — Real-Time Monitoring

* Integration of real-time blockchain transaction monitoring
* Live donor dashboard updates

## Phase 3 — AI Fraud Detection

* AI-powered verification and fraud detection for beneficiary applications

## Phase 4 — Hybrid Fiat-Crypto Donations

* Integration of M-Pesa for seamless fiat and crypto donation support

---

# Conclusion

AidStream combines blockchain transparency, escrow security, and real-time backend synchronization to create a trustworthy humanitarian aid ecosystem. By separating concerns across frontend, backend, database, and blockchain layers, the platform ensures scalability, accountability, and operational efficiency while maintaining donor confidence.
