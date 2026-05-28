const { expect } = require("chai");
const { ethers } = require("hardhat");

describe("VerificationRegistry Contract", function () {

  let VerificationRegistry;
  let verificationRegistry;

  let owner;
  let validator1;
  let validator2;
  let validator3;
  let randomUser;

  // SETUP

  beforeEach(async function () {

    [
      owner,
      validator1,
      validator2,
      validator3,
      randomUser
    ] = await ethers.getSigners();

    VerificationRegistry =
      await ethers.getContractFactory(
        "VerificationRegistry"
      );

    verificationRegistry =
      await VerificationRegistry.deploy();

    await verificationRegistry.waitForDeployment();
  });

  // DEPLOYMENT TESTS

  describe("Deployment", function () {

    it("Should deploy successfully", async function () {

      expect(
        await verificationRegistry.validatorCount()
      ).to.equal(0);
    });

  });

  // VALIDATOR MANAGEMENT TESTS

  describe("Validator Management", function () {

    it("Should add validator successfully", async function () {

      await verificationRegistry.addValidator(
        validator1.address
      );

      expect(
        await verificationRegistry.validators(
          validator1.address
        )
      ).to.equal(true);

      expect(
        await verificationRegistry.validatorCount()
      ).to.equal(1);
    });

    it("Should remove validator successfully", async function () {

      await verificationRegistry.addValidator(
        validator1.address
      );

      await verificationRegistry.removeValidator(
        validator1.address
      );

      expect(
        await verificationRegistry.validators(
          validator1.address
        )
      ).to.equal(false);

      expect(
        await verificationRegistry.validatorCount()
      ).to.equal(0);
    });

    it("Should reject duplicate validator", async function () {

      await verificationRegistry.addValidator(
        validator1.address
      );

      await expect(
        verificationRegistry.addValidator(
          validator1.address
        )
      ).to.be.revertedWith(
        "Validator already exists"
      );
    });

    it("Should reject removing nonexistent validator", async function () {

      await expect(
        verificationRegistry.removeValidator(
          validator1.address
        )
      ).to.be.revertedWith(
        "Validator does not exist"
      );
    });

    it("Should reject non-owner validator addition", async function () {

      await expect(
        verificationRegistry
          .connect(randomUser)
          .addValidator(validator1.address)
      ).to.be.reverted;
    });

    it("Should reject invalid validator address", async function () {

      await expect(
        verificationRegistry.addValidator(
          ethers.ZeroAddress
        )
      ).to.be.revertedWith(
        "Invalid validator address"
      );
    });

  });

  // CASE CREATION TESTS

  describe("Verification Case Creation", function () {

    it("Should create verification case successfully", async function () {

      await verificationRegistry.createVerificationCase(
        1
      );

      const verificationCase =
        await verificationRegistry.verificationCases(1);

      expect(
        verificationCase.caseId
      ).to.equal(1);

      expect(
        verificationCase.verified
      ).to.equal(false);

      expect(
        verificationCase.approvalCount
      ).to.equal(0);
    });

    it("Should reject duplicate case creation", async function () {

      await verificationRegistry.createVerificationCase(
        1
      );

      await expect(
        verificationRegistry.createVerificationCase(1)
      ).to.be.revertedWith(
        "Case already exists"
      );
    });

    it("Should reject non-owner case creation", async function () {

      await expect(
        verificationRegistry
          .connect(randomUser)
          .createVerificationCase(1)
      ).to.be.reverted;
    });

  });

  // VERIFICATION TESTS

  describe("Case Verification", function () {

    beforeEach(async function () {

      await verificationRegistry.addValidator(
        validator1.address
      );

      await verificationRegistry.addValidator(
        validator2.address
      );

      await verificationRegistry.addValidator(
        validator3.address
      );

      await verificationRegistry.createVerificationCase(
        1
      );
    });

    it("Should allow validator to verify case", async function () {

      await verificationRegistry
        .connect(validator1)
        .verifyCase(1);

      const approvalCount =
        await verificationRegistry.getApprovalCount(1);

      expect(approvalCount).to.equal(1);
    });

    it("Should approve case after two validations", async function () {

      await verificationRegistry
        .connect(validator1)
        .verifyCase(1);

      await verificationRegistry
        .connect(validator2)
        .verifyCase(1);

      const isVerified =
        await verificationRegistry.isCaseVerified(1);

      expect(isVerified).to.equal(true);
    });

    it("Should reject unauthorized verifier", async function () {

      await expect(
        verificationRegistry
          .connect(randomUser)
          .verifyCase(1)
      ).to.be.revertedWith(
        "Not an authorized validator"
      );
    });

    it("Should reject double verification by same validator", async function () {

      await verificationRegistry
        .connect(validator1)
        .verifyCase(1);

      await expect(
        verificationRegistry
          .connect(validator1)
          .verifyCase(1)
      ).to.be.revertedWith(
        "Validator already approved this case"
      );
    });

    it("Should reject verification of approved case", async function () {

      await verificationRegistry
        .connect(validator1)
        .verifyCase(1);

      await verificationRegistry
        .connect(validator2)
        .verifyCase(1);

      await expect(
        verificationRegistry
          .connect(validator3)
          .verifyCase(1)
      ).to.be.revertedWith(
        "Case already approved"
      );
    });

    it("Should reject verification of nonexistent case", async function () {

      await expect(
        verificationRegistry
          .connect(validator1)
          .verifyCase(999)
      ).to.be.revertedWith(
        "Verification case does not exist"
      );
    });

  });

  // VERIFICATION HISTORY TESTS

  describe("Verification History", function () {

    beforeEach(async function () {

      await verificationRegistry.addValidator(
        validator1.address
      );

      await verificationRegistry.addValidator(
        validator2.address
      );

      await verificationRegistry.createVerificationCase(
        1
      );
    });

    it("Should store verification history", async function () {

      await verificationRegistry
        .connect(validator1)
        .verifyCase(1);

      await verificationRegistry
        .connect(validator2)
        .verifyCase(1);

      const history =
        await verificationRegistry.getVerificationHistory(1);

      expect(history.length).to.equal(2);

      expect(
        history[0].validator
      ).to.equal(validator1.address);

      expect(
        history[1].validator
      ).to.equal(validator2.address);
    });

  });

  // VIEW FUNCTION TESTS

  describe("View Functions", function () {

    beforeEach(async function () {

      await verificationRegistry.addValidator(
        validator1.address
      );

      await verificationRegistry.addValidator(
        validator2.address
      );

      await verificationRegistry.createVerificationCase(
        1
      );
    });

    it("Should return correct approval count", async function () {

      await verificationRegistry
        .connect(validator1)
        .verifyCase(1);

      const approvalCount =
        await verificationRegistry.getApprovalCount(1);

      expect(approvalCount).to.equal(1);
    });

    it("Should return false for unverified case", async function () {

      const isVerified =
        await verificationRegistry.isCaseVerified(1);

      expect(isVerified).to.equal(false);
    });

    it("Should return true after approval threshold", async function () {

      await verificationRegistry
        .connect(validator1)
        .verifyCase(1);

      await verificationRegistry
        .connect(validator2)
        .verifyCase(1);

      const isVerified =
        await verificationRegistry.isCaseVerified(1);

      expect(isVerified).to.equal(true);
    });

  });

});