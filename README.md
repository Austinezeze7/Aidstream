AidStream

Transparent Aid Distribution Powered by Blockchain

AidStream is a blockchain-powered transparent aid distribution platform designed to ensure donations reach verified beneficiaries securely, transparently, and efficiently.

The platform enables donors, NGOs, and community verifiers to track aid distribution in real time using blockchain transaction records and transparent verification systems.

---

Problem Statement

Traditional aid and donation systems often face challenges such as:

- Lack of transparency
- Fraudulent beneficiaries
- Misuse of donated funds
- Poor accountability
- Limited visibility into how donations are used

AidStream addresses these challenges by leveraging blockchain technology, smart contracts, and community verification mechanisms to build trust between donors and beneficiaries.

---

Solution Overview

AidStream creates a transparent ecosystem where:

- Beneficiaries are verified before receiving aid
- Donations are recorded on the blockchain
- Donors can track transactions publicly
- Administrators can monitor fund disbursement
- Communities can verify real cases

Every donation transaction generates a blockchain transaction hash, ensuring that all records remain immutable and publicly verifiable.

---

Core Features

Authentication System

- Unified login system
- Role-based access control
- Donor and admin dashboards

---

Beneficiary Verification

- Community verification workflow
- NGO verification support
- Verification status tracking

---

Donation Tracking

- Blockchain transaction recording
- Transparent donation history
- Public transaction hashes

---

Smart Contract Escrow

- Funds held securely before release
- Controlled disbursement process
- Reduced risk of fraud

---

Transparency Dashboard

- Real-time statistics
- Donation analytics
- Active beneficiary tracking
- Public accountability system

---

Tech Stack

Frontend

- HTML5
- CSS3
- JavaScript

---

Backend

- Go (Golang)

---

Database

- Neon PostgreSQL

---

Blockchain

- Ethereum-compatible smart contracts
- MetaMask wallet integration

---

## Project Structure

Aidstream/
├── backend/           # Go backend APIs and business logic
├── database/          # PostgreSQL schema and migrations
├── frontend/          # HTML, CSS, JavaScript frontend
├── smart-contracts/   # Ethereum smart contracts
├── Docs/              # Project documentation
├── main.go            # Application entry point
└── README.md

---

Team Members

Member| Role| Responsibilities
Emma| Frontend Developer| UI/UX design, HTML, CSS, JavaScript
Edward| Backend Developer| Go APIs, authentication, backend logic
Philip| Database Engineer| Neon PostgreSQL setup and schema management
Paul| Blockchain Developer| Smart contracts, wallet integration, blockchain validation
Austine| Integration & Documentation Lead| Project integration, testing, GitHub management, documentation, presentation

---

System Workflow

1. Beneficiary Registration

A beneficiary case is submitted to the platform.

---

2. Verification Process

Admins or community verifiers validate the legitimacy of the case.

---

3. Donor Authentication

Users log in and connect their crypto wallets.

---

4. Donation Process

Donors send funds through blockchain transactions.

---

5. Blockchain Validation

Transactions are validated by blockchain network validators and permanently recorded.

---

6. Dashboard Updates

The system updates donation statistics and transaction records in real time.

---

Database Overview

The platform uses Neon PostgreSQL to store:

- Users
- Beneficiaries
- Donations
- Verification records
- Transaction hashes
- Wallet addresses

---

Blockchain Architecture

AidStream uses blockchain technology to:

- Record donation transactions
- Validate transaction authenticity
- Ensure transparency
- Prevent record tampering
- Improve donor trust

Each donation generates a unique transaction hash which can be publicly verified.

---

Smart Contract Functions

Donation Handling

donateToCase(caseId)

Fund Escrow

lockFunds(caseId)

Fund Release

releaseFunds(caseId)

---

Installation Guide

Clone Repository

git clone <repository-url>

---

Frontend Setup

Navigate to frontend files and open:

index.html

---

Backend Setup

Install Go dependencies and run:

go run main.go

---

Database Setup

1. Create a Neon PostgreSQL database
2. Configure environment variables
3. Run migrations

---

Future Improvements

- M-Pesa integration
- SMS notifications
- AI fraud detection
- Mobile application
- Multi-chain blockchain support
- Beneficiary geo-location tracking
- Real-time analytics dashboard

---

Hackathon Vision

AidStream aims to demonstrate how blockchain technology can improve transparency, accountability, and trust in humanitarian aid systems and charitable donations.

---

License

This project was developed for the Zone01 Hackathon.

---

Acknowledgements

Special thanks to:

- Zone01 Hackathon organizers
- Open-source contributors
- Blockchain developer communities

