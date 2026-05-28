const { expect } = require("chai");
const { ethers } = require("hardhat");

describe("AidStream Contract", function () {

  let AidStream;
  let aidStream;

  let owner;
  let validator;
  let validator2;
  let donor;
  let donor2;
  let beneficiary;

  // SETUP

  beforeEach(async function () {

    [
      owner,
      validator,
      validator2,
      donor,
      donor2,
      beneficiary
    ] = await ethers.getSigners();

    AidStream = await ethers.getContractFactory(
      "AidStream"
    );

    aidStream = await AidStream.deploy();

    await aidStream.waitForDeployment();
  });

  // DEPLOYMENT TESTS

  describe("Deployment", function () {

    it("Should deploy successfully", async function () {

      expect(
        await aidStream.caseCount()
      ).to.equal(0);
    });

  });

  // VALIDATOR TESTS

  describe("Validator Management", function () {

    it("Should add validator successfully", async function () {

      await aidStream.addValidator(
        validator.address
      );

      expect(
        await aidStream.validators(
          validator.address
        )
      ).to.equal(true);
    });

    it("Should reject non-owner validator addition", async function () {

      await expect(
        aidStream
          .connect(donor)
          .addValidator(validator.address)
      ).to.be.reverted;
    });

  });

  // CASE CREATION TESTS

  describe("Case Creation", function () {

    it("Should create a beneficiary case", async function () {

      await aidStream.createCase(
        "Medical Support",
        "Emergency surgery support",
        beneficiary.address,
        ethers.parseEther("5")
      );

      const aidCase =
        await aidStream.getCase(1);

      expect(aidCase.id).to.equal(1);

      expect(aidCase.title).to.equal(
        "Medical Support"
      );

      expect(aidCase.beneficiary).to.equal(
        beneficiary.address
      );

      expect(aidCase.goalAmount).to.equal(
        ethers.parseEther("5")
      );

      expect(aidCase.verified).to.equal(false);
    });

    it("Should reject invalid beneficiary", async function () {

      await expect(
        aidStream.createCase(
          "Invalid Case",
          "Invalid beneficiary",
          ethers.ZeroAddress,
          ethers.parseEther("1")
        )
      ).to.be.revertedWith(
        "Invalid beneficiary"
      );
    });

    it("Should reject zero goal amount", async function () {

      await expect(
        aidStream.createCase(
          "Zero Goal",
          "Invalid goal amount",
          beneficiary.address,
          0
        )
      ).to.be.revertedWith(
        "Goal amount must be greater than zero"
      );
    });

  });

  // VERIFICATION TESTS

  describe("Case Verification", function () {

    beforeEach(async function () {

      await aidStream.addValidator(
        validator.address
      );

      await aidStream.createCase(
        "School Fees",
        "Support tuition fees",
        beneficiary.address,
        ethers.parseEther("3")
      );
    });

    it("Should verify a case", async function () {

      await aidStream
        .connect(validator)
        .verifyCase(1);

      const aidCase =
        await aidStream.getCase(1);

      expect(
        aidCase.verified
      ).to.equal(true);
    });

    it("Should reject unauthorized verifier", async function () {

      await expect(
        aidStream
          .connect(donor)
          .verifyCase(1)
      ).to.be.revertedWith(
        "Not an authorized validator"
      );
    });

    it("Should reject double verification", async function () {

      await aidStream
        .connect(validator)
        .verifyCase(1);

      await expect(
        aidStream
          .connect(validator)
          .verifyCase(1)
      ).to.be.revertedWith(
        "Case already verified"
      );
    });

  });

  // DONATION TESTS

  describe("Donations", function () {

    beforeEach(async function () {

      await aidStream.addValidator(
        validator.address
      );

      await aidStream.createCase(
        "Hospital Bills",
        "Critical treatment support",
        beneficiary.address,
        ethers.parseEther("10")
      );

      await aidStream
        .connect(validator)
        .verifyCase(1);
    });

    it("Should donate successfully", async function () {

      await aidStream
        .connect(donor)
        .donateToCase(1, {
          value: ethers.parseEther("1")
        });

      const aidCase =
        await aidStream.getCase(1);

      expect(
        aidCase.donatedAmount
      ).to.equal(
        ethers.parseEther("1")
      );
    });

    it("Should record donor history", async function () {

      await aidStream
        .connect(donor)
        .donateToCase(1, {
          value: ethers.parseEther("2")
        });

      const donorHistory =
        await aidStream.getDonorHistory(
          donor.address
        );

      expect(
        donorHistory.length
      ).to.equal(1);

      expect(
        donorHistory[0]
      ).to.equal(1);
    });

    it("Should reject zero donation", async function () {

      await expect(
        aidStream
          .connect(donor)
          .donateToCase(1, {
            value: 0
          })
      ).to.be.revertedWith(
        "Donation amount must be greater than zero"
      );
    });

    it("Should reject donation to unverified case", async function () {

      await aidStream.createCase(
        "Flood Relief",
        "Emergency flood support",
        beneficiary.address,
        ethers.parseEther("5")
      );

      await expect(
        aidStream
          .connect(donor2)
          .donateToCase(2, {
            value: ethers.parseEther("1")
          })
      ).to.be.revertedWith(
        "Case is not verified"
      );
    });

  });

  // FUND RELEASE TESTS

  describe("Fund Release", function () {

    beforeEach(async function () {

      await aidStream.addValidator(
        validator.address
      );

      await aidStream.createCase(
        "Food Relief",
        "Community feeding program",
        beneficiary.address,
        ethers.parseEther("5")
      );

      await aidStream
        .connect(validator)
        .verifyCase(1);

      await aidStream
        .connect(donor)
        .donateToCase(1, {
          value: ethers.parseEther("2")
        });
    });

    it("Should release funds successfully", async function () {

      const balanceBefore =
        await ethers.provider.getBalance(
          beneficiary.address
        );

      await aidStream.releaseFunds(1);

      const balanceAfter =
        await ethers.provider.getBalance(
          beneficiary.address
        );

      expect(
        balanceAfter
      ).to.be.gt(balanceBefore);
    });

    it("Should reject double release", async function () {

      await aidStream.releaseFunds(1);

      await expect(
        aidStream.releaseFunds(1)
      ).to.be.revertedWith(
        "Funds already released"
      );
    });

    it("Should reject non-owner release", async function () {

      await expect(
        aidStream
          .connect(donor)
          .releaseFunds(1)
      ).to.be.reverted;
    });

  });

  // VIEW FUNCTION TESTS

  describe("View Functions", function () {

    beforeEach(async function () {

      await aidStream.addValidator(
        validator.address
      );

      await aidStream.createCase(
        "Education Support",
        "Scholarship funding",
        beneficiary.address,
        ethers.parseEther("4")
      );

      await aidStream
        .connect(validator)
        .verifyCase(1);

      await aidStream
        .connect(donor)
        .donateToCase(1, {
          value: ethers.parseEther("1")
        });
    });

    it("Should return case donations", async function () {

      const donations =
        await aidStream.getCaseDonations(1);

      expect(
        donations.length
      ).to.equal(1);

      expect(
        donations[0].donor
      ).to.equal(donor.address);
    });

    it("Should return contract balance", async function () {

      const balance =
        await aidStream.getContractBalance();

      expect(balance).to.equal(
        ethers.parseEther("1")
      );
    });

  });

});