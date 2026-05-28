/*
 * deploy.js
 * --------------------------------------------------------
 * AidStream Smart Contracts Deployment Script
 *
 * Deploys:
 * - AidStream.sol
 * - Escrow.sol
 * - VerificationRegistry.sol
 *
 * Network:
 * - Local Hardhat
 * - Polygon Amoy Testnet
 *
 * Usage:
 * npx hardhat run scripts/deploy.js --network polygonAmoy
 */

const hre = require("hardhat");

async function main() {
    console.log("================================================");
    console.log("🚀 Starting AidStream Contract Deployment...");
    console.log("================================================\n");

    // GET DEPLOYER ACCOUNT

    const [deployer] = await hre.ethers.getSigners();

    console.log("Deploying contracts with account:");
    console.log(deployer.address);

    const balance = await hre.ethers.provider.getBalance(
        deployer.address
    );

    console.log(
        "Account balance:",
        hre.ethers.utils.formatEther(balance),
        "ETH\n"
    );

    // DEPLOY AIDSTREAM CONTRACT

    console.log("📦 Deploying AidStream contract...");

    const AidStream = await hre.ethers.getContractFactory(
        "AidStream"
    );

    const aidStream = await AidStream.deploy();

    await aidStream.waitForDeployment();

    const aidStreamAddress = await aidStream.getAddress();

    console.log(
        "✅ AidStream deployed at:",
        aidStreamAddress,
        "\n"
    );

    // DEPLOY ESCROW CONTRACT

    console.log("📦 Deploying Escrow contract...");

    const Escrow = await hre.ethers.getContractFactory(
        "Escrow"
    );

    const escrow = await Escrow.deploy();

    await escrow.waitForDeployment();

    const escrowAddress = await escrow.getAddress();

    console.log(
        "✅ Escrow deployed to:",
        escrowAddress,
        "\n"
    );

    // DEPLOY VERIFICATION REGISTRY

    console.log("📦 Deploying VerificationRegistry contract..."
    );

    const VerificationRegistry = await VerificationRegistry.deploy();

    await VerificationRegistry.waitForDeployment();

    const verificationRegistryAddress = await verificationRegistry.getAddress();

    console.log(
        "✅ VerificationRegistry deployed at:",
        verificationRegistryAddress,
        "\n"
    );

    // DEPLOYMENT SUMMARY
    
    console.log("================================================");
    console.log("🎉 DEPLOYMENT SUCCESSFUL");
    console.log("================================================\n");

    console.log("📌 Contract Addresses:\n");

    console.log(
        "AidStream:",
        aidStreamAddress
    );

    console.log(
        "Escrow:",
        escrowAddress
    );

    console.log(
        "VerificationRegistry:",
        verificationRegistryAddress
    );

    console.log("\n================================================");
}

// ERROR HANDLING

main().catch((error) => {

    console.error("❌ Deployment failed:", error);
    console.log(error);

    process.exitCode = 1;
});