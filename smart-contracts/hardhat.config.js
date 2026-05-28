require("@nomicfoundation/hardhat-toolbox");
require("@nomicfoundation/hardhat-ignition");
require("dotenv").config();

/*
 * hardhat.config.js
 * --------------------------------------------------------
 * AidStream Hardhat Configuration
 *
 * Features:
 * - Solidity compiler setup
 * - Polygon Amoy configuration
 * - Local Hardhat network
 * - Environment variable support
 * - Ignition deployment support
 *
 * Required .env variables:
 *
 * PRIVATE_KEY=your_wallet_private_key
 * RPC_URL=your_polygon_amoy_rpc_url
 * ETHERSCAN_API_KEY=your_polygonscan_api_key
 */

module.exports = {

    // SOLIDITY CONFIGURATION

    solidity: {
        version: "0.8.20",

        settings: {
            optimizer: {
                enabled: true,
                runs: 200
            }
        }
    },

    // NETWORK CONFIGURATION

    networks: {

        // Local Hardhat Network
        hardhat: {},

        // Polygon Amoy Testnet
        polygonAmoy: {
            url: process.env.RPC_URL || "",
            accounts:
                process.env.PRIVATE_KEY !== undefined
                    ? [process.env.PRIVATE_KEY]
                    : [],
            chainId: 80002
        }
    },

    // ETHERSCAN VERIFICATION

    etherscan: {
        apiKey: process.env.ETHERSCAN_API_KEY || ""
    },

    // PATH CONFIGURATION

    paths: {
        sources: "./contracts",
        tests: "./test",
        cache: "./cache",
        artifacts: "./artifacts"
    },

    // MOCHA TEST CONFIGURATION
    
    mocha: {
        timeout: 40000
    }
};