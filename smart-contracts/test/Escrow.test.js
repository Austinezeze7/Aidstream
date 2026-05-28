const { expect } = require("chai");
const { ethers } = require("hardhat");

describe("Escrow Contract", function () {

  let Escrow;
  let escrow;

  let owner;
  let validator;
  let donor;
  let donor2;
  let beneficiary;

  // SETUP

  beforeEach(async function () {

    [
      owner,
      validator,
      donor,
      donor2,
      beneficiary
    ] = await ethers.getSigners();

    Escrow = await ethers.getContractFactory(
      "Escrow"
    );

    escrow = await Escrow.deploy();

    await escrow.waitForDeployment();
  });

  // DEPLOYMENT TESTS

  describe("Deployment", function () {

    it("Should deploy successfully", async function () {

      const balance =
        await escrow.getContractBalance();

      expect(balance).to.equal(0);
    });

  });

  // VALIDATOR TESTS

  describe("Validator Management", function () {

    it("Should add validator successfully", async function () {

      await escrow.addValidator(
        validator.address
      );

      expect(
        await escrow.validators(
          validator.address
        )
      ).to.equal(true);
    });

    it("Should reject non-owner validator addition", async function () {

      await expect(
        escrow
          .connect(donor)
          .addValidator(validator.address)
      ).to.be.reverted;
    });

  });

  // ESCROW CREATION TESTS

  describe("Escrow Creation", function () {

    it("Should create escrow successfully", async function () {

      await escrow.createEscrow(
        1,
        beneficiary.address
      );

      const escrowCase =
        await escrow.getEscrowCase(1);

      expect(
        escrowCase.caseId
      ).to.equal(1);

      expect(
        escrowCase.beneficiary
      ).to.equal(
        beneficiary.address
      );

      expect(
        escrowCase.totalFunds
      ).to.equal(0);

      expect(
        escrowCase.verified
      ).to.equal(false);
    });

    it("Should reject duplicate escrow creation", async function () {

      await escrow.createEscrow(
        1,
        beneficiary.address
      );

      await expect(
        escrow.createEscrow(
          1,
          beneficiary.address
        )
      ).to.be.revertedWith(
        "Escrow already exists"
      );
    });

    it("Should reject invalid beneficiary", async function () {

      await expect(
        escrow.createEscrow(
          1,
          ethers.ZeroAddress
        )
      ).to.be.revertedWith(
        "Invalid beneficiary address"
      );
    });

  });

  // DONATION TESTS

  describe("Deposits", function () {

    beforeEach(async function () {

      await escrow.createEscrow(
        1,
        beneficiary.address
      );
    });

    it("Should deposit funds successfully", async function () {

      await escrow
        .connect(donor)
        .depositFunds(1, {
          value: ethers.parseEther("1")
        });

      const escrowCase =
        await escrow.getEscrowCase(1);

      expect(
        escrowCase.totalFunds
      ).to.equal(
        ethers.parseEther("1")
      );
    });

    it("Should track donor contributions", async function () {

      await escrow
        .connect(donor)
        .depositFunds(1, {
          value: ethers.parseEther("2")
        });

      const contribution =
        await escrow.getContribution(
          1,
          donor.address
        );

      expect(contribution).to.equal(
        ethers.parseEther("2")
      );
    });

    it("Should allow multiple donors", async function () {

      await escrow
        .connect(donor)
        .depositFunds(1, {
          value: ethers.parseEther("1")
        });

      await escrow
        .connect(donor2)
        .depositFunds(1, {
          value: ethers.parseEther("3")
        });

      const escrowCase =
        await escrow.getEscrowCase(1);

      expect(
        escrowCase.totalFunds
      ).to.equal(
        ethers.parseEther("4")
      );
    });

    it("Should reject zero deposit", async function () {

      await expect(
        escrow
          .connect(donor)
          .depositFunds(1, {
            value: 0
          })
      ).to.be.revertedWith(
        "Deposit amount must be greater than zero"
      );
    });

    it("Should reject deposits to nonexistent escrow", async function () {

      await expect(
        escrow
          .connect(donor)
          .depositFunds(999, {
            value: ethers.parseEther("1")
          })
      ).to.be.revertedWith(
        "Escrow case does not exist"
      );
    });

  });

  // VERIFICATION TESTS

  describe("Verification", function () {

    beforeEach(async function () {

      await escrow.addValidator(
        validator.address
      );

      await escrow.createEscrow(
        1,
        beneficiary.address
      );
    });

    it("Should verify escrow case successfully", async function () {

      await escrow
        .connect(validator)
        .verifyEscrowCase(1);

      const escrowCase =
        await escrow.getEscrowCase(1);

      expect(
        escrowCase.verified
      ).to.equal(true);
    });

    it("Should reject unauthorized validator", async function () {

      await expect(
        escrow
          .connect(donor)
          .verifyEscrowCase(1)
      ).to.be.revertedWith(
        "Not an authorized validator"
      );
    });

    it("Should reject double verification", async function () {

      await escrow
        .connect(validator)
        .verifyEscrowCase(1);

      await expect(
        escrow
          .connect(validator)
          .verifyEscrowCase(1)
      ).to.be.revertedWith(
        "Case already verified"
      );
    });

  });

  // FUND RELEASE TESTS

  describe("Fund Release", function () {

    beforeEach(async function () {

      await escrow.addValidator(
        validator.address
      );

      await escrow.createEscrow(
        1,
        beneficiary.address
      );

      await escrow
        .connect(validator)
        .verifyEscrowCase(1);

      await escrow
        .connect(donor)
        .depositFunds(1, {
          value: ethers.parseEther("2")
        });
    });

    it("Should release funds successfully", async function () {

      const balanceBefore =
        await ethers.provider.getBalance(
          beneficiary.address
        );

      await escrow.releaseFunds(1);

      const balanceAfter =
        await ethers.provider.getBalance(
          beneficiary.address
        );

      expect(
        balanceAfter
      ).to.be.gt(balanceBefore);
    });

    it("Should reject release for unverified case", async function () {

      await escrow.createEscrow(
        2,
        beneficiary.address
      );

      await escrow
        .connect(donor)
        .depositFunds(2, {
          value: ethers.parseEther("1")
        });

      await expect(
        escrow.releaseFunds(2)
      ).to.be.revertedWith(
        "Case is not verified"
      );
    });

    it("Should reject double release", async function () {

      await escrow.releaseFunds(1);

      await expect(
        escrow.releaseFunds(1)
      ).to.be.revertedWith(
        "Funds already released"
      );
    });

    it("Should reject non-owner release", async function () {

      await expect(
        escrow
          .connect(donor)
          .releaseFunds(1)
      ).to.be.reverted;
    });

    it("Should reject release with no funds", async function () {

      await escrow.createEscrow(
        3,
        beneficiary.address
      );

      await escrow
        .connect(validator)
        .verifyEscrowCase(3);

      await expect(
        escrow.releaseFunds(3)
      ).to.be.revertedWith(
        "No funds available"
      );
    });

  });

  // VIEW FUNCTION TESTS

  describe("View Functions", function () {

    beforeEach(async function () {

      await escrow.createEscrow(
        1,
        beneficiary.address
      );

      await escrow
        .connect(donor)
        .depositFunds(1, {
          value: ethers.parseEther("1")
        });
    });

    it("Should return contract balance", async function () {

      const balance =
        await escrow.getContractBalance();

      expect(balance).to.equal(
        ethers.parseEther("1")
      );
    });

    it("Should return donor contribution", async function () {

      const contribution =
        await escrow.getContribution(
          1,
          donor.address
        );

      expect(contribution).to.equal(
        ethers.parseEther("1")
      );
    });

  });

});