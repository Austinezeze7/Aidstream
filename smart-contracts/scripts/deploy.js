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
 * Usage:
 * Local:
 * npx hardhat run scripts/deploy.js
 *
 * Polygon Amoy:
 * npx hardhat run scripts/deploy.js --network polygonAmoy
 */

const hre = require("hardhat");

async function main() {

    console.log("================================================");
    console.log("🚀 Starting AidStream Contract Deployment...");
    console.log("================================================\n");

    // =====================================================
    // GET DEPLOYER ACCOUNT
    // =====================================================

    const [deployer] = await hre.ethers.getSigners();

    console.log("Deploying contracts with account:");
    console.log(deployer.address);

    const balance = await hre.ethers.provider.getBalance(
        deployer.address
    );

    console.log(
        "Account balance:",
        hre.ethers.formatEther(balance),
        "POL\n"
    );

    // =====================================================
    // DEPLOY AIDSTREAM CONTRACT
    // =====================================================

    console.log("📦 Deploying AidStream contract...");

    const AidStream = await hre.ethers.getContractFactory(
        "AidStream"
    );

    const aidStream = await AidStream.deploy();

    // FIX:
    // waitForDeployment() MUST be awaited before using contract
    await aidStream.waitForDeployment();

    const aidStreamAddress = await aidStream.getAddress();

    console.log(
        "✅ AidStream deployed at:",
        aidStreamAddress,
        "\n"
    );

    // =====================================================
    // DEPLOY ESCROW CONTRACT
    // =====================================================

    console.log("📦 Deploying Escrow contract...");

    const Escrow = await hre.ethers.getContractFactory(
        "Escrow"
    );

    const escrow = await Escrow.deploy();

    await escrow.waitForDeployment();

    const escrowAddress = await escrow.getAddress();

    console.log(
        "✅ Escrow deployed at:",
        escrowAddress,
        "\n"
    );

    // =====================================================
    // DEPLOY VERIFICATION REGISTRY CONTRACT
    // =====================================================

    console.log("📦 Deploying VerificationRegistry contract...");

    const VerificationRegistry =
        await hre.ethers.getContractFactory(
            "VerificationRegistry"
        );

    const verificationRegistry =
        await VerificationRegistry.deploy();

    await verificationRegistry.waitForDeployment();

    const verificationRegistryAddress =
        await verificationRegistry.getAddress();

    console.log(
        "✅ VerificationRegistry deployed at:",
        verificationRegistryAddress,
        "\n"
    );

    // =====================================================
    // OPTIONAL INITIAL SETUP
    // =====================================================

    console.log("⚙️ Running initial setup...\n");

    /*
     * Example:
     * Add deployer as validator in AidStream
     */

    const tx1 = await aidStream.addValidator(
        deployer.address
    );

    await tx1.wait();

    console.log(
        "✅ Added deployer as AidStream validator"
    );

    /*
     * Add deployer as validator in Escrow
     */

    const tx2 = await escrow.addValidator(
        deployer.address
    );

    await tx2.wait();

    console.log(
        "✅ Added deployer as Escrow validator"
    );

    /*
     * Add deployer as validator in VerificationRegistry
     */

    const tx3 =
        await verificationRegistry.addValidator(
            deployer.address
        );

    await tx3.wait();

    console.log(
        "✅ Added deployer as VerificationRegistry validator\n"
    );

    // =====================================================
    // DEPLOYMENT SUMMARY
    // =====================================================

    console.log("================================================");
    console.log("🎉 DEPLOYMENT SUCCESSFUL");
    console.log("================================================\n");

    console.log("📌 Contract Addresses:\n");

    console.log("AidStream:");
    console.log(aidStreamAddress, "\n");

    console.log("Escrow:");
    console.log(escrowAddress, "\n");

    console.log("VerificationRegistry:");
    console.log(verificationRegistryAddress, "\n");

    console.log("================================================");

    // =====================================================
    // FRONTEND ENV VARIABLES
    // =====================================================

    console.log("\n📄 Copy these into your frontend .env file:\n");

    console.log(
        `NEXT_PUBLIC_AIDSTREAM_ADDRESS=${aidStreamAddress}`
    );

    console.log(
        `NEXT_PUBLIC_ESCROW_ADDRESS=${escrowAddress}`
    );

    console.log(
        `NEXT_PUBLIC_VERIFICATION_REGISTRY_ADDRESS=${verificationRegistryAddress}`
    );

    console.log("\n================================================");
}

// =========================================================
// ERROR HANDLING
// =========================================================

main().catch((error) => {

    console.error("\n❌ Deployment failed:\n");

    console.error(error);

    process.exitCode = 1;
});