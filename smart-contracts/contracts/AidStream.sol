// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

/*
 * AidStream.sol
 * ------------------------------------------
 * Main donation tracking contract for AidStream.
 *
 * Features:
 * - Create beneficiary aid cases
 * - Donate to verified cases
 * - Track donations transparently
 * - Store donor history
 * - Emit blockchain events
 *
 * NOTE:
 * Escrow release logic can later be delegated
 * to Escrow.sol for better modularity.
 */

 import "@openzeppelin/contracts/access/Ownable.sol";
 import "@openzeppelin/contracts/security/ReentrancyGuard.sol";

 contract AidStream is Ownable, ReentrancyGuard {
    // STRUCTS
    
    struct Case {
        uint256 id;
        string title;
        string description;
        address payable beneficiary;
        uint256 goalAmount;
        uint256 donatedAmount;
        bool verified;
        bool fundsReleased;
        bool active;
        uint256 createdAt;
    }

    struct Donation {
        address donor;
        uint256 amount;
        uint256 timestamp;
    }

    // STATE VARIABLES

    uint256 public caseCount;

    mapping(uint256 => Case) public cases;

    mapping(uint256 => Donation[]) public caseDonations;

    mapping(address => uint256[]) public donorHistory;

    mapping(address => bool) public validators;

    // EVENTS

    event CaseCreated(
        uint256 indexed caseId,
        string title,
        address indexed beneficiary,
        uint256 goalAmount
    );

    event CaseVerified(
        uint256 indexed caseId,
        address indexed validator
    );

    event DonationReceived(
        uint256 indexed caseId,
        address indexed donor,
        uint256 amount
    );

    event FundsReleased(
        uint256 indexed caseId,
        address indexed beneficiary,
        uint256 amount
    );

    event ValidatorAdded(address indexed validator);

    // MODIFIERS

    modifier onlyValidator() {
        require(validators[msg.sender], "Not an authorized validator");
        _;
    }

    modifier caseExists(uint256 _caseId) {
        require(_caseId > 0 && _caseId <= caseCount, "Case does not exist");
        _;
    }

    modifier caseIsActive(uint256 _caseId) {
        require(cases[_caseId].active, "Case is not active");
        _;
    }

    // CONSTRUCTOR

    constructor() Ownable(msg.sender) {}

    // ADMIN FUNCTIONS

    function addValidator(address _validator) external onlyOwner {
        require(_validator != address(0), "Invalid address");

        validators[_validator] = true;

        emit ValidatorAdded(_validator);
    }

    // CASE MANAGEMENT

    function createCase(
        string memory _title,
        string memory _description,
        address payable _beneficiary,
        uint256 _goalAmount
    ) external {
        require(bytes(_title).length > 0, "Title required");
        require(_beneficiary != address(0), "Invalid beneficiary");
        require(_goalAmount > 0, "Goal amount must be greater than zero");

        caseCount++;

        cases[caseCount] = Case({
            id: caseCount,
            title: _title,
            description: _description,
            beneficiary: _beneficiary,
            goalAmount: _goalAmount,
            donatedAmount: 0,
            verified: false,
            active: true,
            createdAt: block.timestamp
        });

        emit CaseCreated(
            caseCount,
            _title,
            _beneficiary,
            _goalAmount
        );
    }

    // VERIFICATION

    function verifyCase(
        uint256 _caseId
    )
        external
        onlyValidator
        caseExists(_caseId)
    {
        Case storage aidCase = cases[_caseId];

        require(!aidCase.verified, "Case already verified");

        aidCase.verified = true;

        emit CaseVerified(_caseId, msg.sender);
    }

    // DONATIONS 

    function donateToCase(
        uint256 _caseId
    )
        external
        payable
        nonReentrant
        caseExists(_caseId)
        caseIsActive(_caseId)
    {
        require(msg.value > 0, "Donation amount must be greater than zero");

        Case storage aidCase = cases[_caseId];

        require(aidCase.verified, "Case is not verified");

        aidCase.donatedAmount += msg.value;

        caseDonations[_caseId].push(
            Donation({
                donor: msg.sender,
                amount: msg.value,
                timestamp: block.timestamp
            })
        );

        donorHistory[msg.sender].push(_caseId);

        emit DonationReceived(
            _caseId,
            msg.sender,
            msg.value
        );
    }

    // FUND RELEASE 

    function releaseFunds(
        uint256 _caseId
    )
        external
        onlyOwner
        nonReentrant
        caseExists(_caseId)
    {
        Case storage aidCase = cases[_caseId];

        require(aidCase.verified, "Case is not verified");
        require(!aidCase.fundsReleased, "Funds already released");
        require(aidCase.donatedAmount > 0, "No funds available");

        uint256 amount = aidCase.donatedAmount;

        aidCase.fundsReleased = true;
        aidCase.active = false;

        (bool success, ) = aidCase.beneficiary.call{value: amount}("");

        require(success, "Transfer failed");

        emit FundsReleased(
            _caseId,
            aidCase.beneficiary,
            amount
        );
    }

    // VIEW FUNCTIONS

    function getCase(
        uint256 _caseId
    )
        external
        view 
        caseExists(_caseId)
        returns (Case memory)
    {
        return cases[_caseId];
    }

    function getCaseDonations(
        uint256 _caseId
    )
        external
        view
        caseExists(_caseId)
        returns (Donation[] memory)
    {
        return caseDonations[_caseId];
    }

    function getDonorHistory(
        address _donor
    )
        external
        view
        returns (uint256[] memory)
    {
        return donorHistory[_donor];
    }

    function getContractBalance()
        external
        view
        returns (uint256)
    {
        return address(this).balance;
    }
 }