const { ethers } = require("hardhat");

async function main() {
    const [deployer] = await ethers.getSigners();
    
    console.log("Deploying contracts with the account:", deployer.address);
    const balance = await deployer.provider.getBalance(deployer.address).toString();
    console.log("Account balance:", (await deployer.provider.getBalance(deployer.address)).toString());

    // Deploy GeneNFT contract
    const NTFToken = await ethers.getContractFactory("GeneNFT");
    const ntfToken = await NTFToken.deploy();
    console.log("NTF Token deployed at address:", ntfToken.target);

    // Deploy PostCovidStrokePrevention (PCSP) contract
    const PCSPToken = await ethers.getContractFactory("PostCovidStrokePrevention");
    const pcspToken = await PCSPToken.deploy();
    console.log("PCSP Token deployed at address:", pcspToken.target);

    // Deploy Controller contract with references to the previous contracts
    const Controller = await ethers.getContractFactory("Controller");
    const controller = await Controller.deploy(ntfToken.target, pcspToken.target);
    console.log("Controller deployed at address:", controller.target);

    // Transfer ownership of NTFToken and PCSPToken to the Controller contract
    await ntfToken.transferOwnership(controller.target);
    await pcspToken.transferOwnership(controller.target);
    console.log("Ownership of NTFToken and PCSPToken transferred to Controller.");
}

main().catch((error) => {
    console.error("Error during deployment:", error);
    process.exit(1);
});
