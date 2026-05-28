// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

/*
 * Escrow.sol
 * --------------------------------------------------------
 * AidStream Escrow Contract
 *
 * Responsibilities:
 * - Lock donated funds securely
 * - Release funds only after approval
 * - Prevent double withdrawals
 * - Maintain transparent transaction history
 *
 * Security:
 * - Ownable access control
 * - Reentrancy protection
 * - Validations for releases and deposits
 */

import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";

contract Escrow is Ownable, ReentrancyGuard {

    // STRUCTS

    struct EscrowCase {
        uint256 caseId;
        address payable beneficiary;
        uint256 totalFunds;
        bool verified;
        bool released;
        uint256 createdAt;
    }

    // STATE VARIABLES

    mapping(uint256 => EscrowCase) public escrowCases;

    mapping(uint256 => mapping(address => uint256)) public donorContributions;

    mapping(address => bool) public validators;

    // EVENTS

    event EscrowCreated(
        uint256 indexed caseId,
        address indexed beneficiary
    );

    event FundsDeposited(
        uint256 indexed caseId,
        address indexed donor,
        uint256 amount
    );

    event CaseVerified(
        uint256 indexed caseId,
        address indexed validator
    );

    event FundsReleased(
        uint256 indexed caseId,
        address indexed beneficiary,
        uint256 amount
    );

    event ValidatorAdded(address indexed validator);

    // MODIFIER

    modifier onlyValidator() {
        require(
            validators[msg.sender],
            "Not an authorized validator"
        );
        _;
    }

    modifier escrowExists(uint256 _caseId) {
        require(
            escrowCases[_caseId].beneficiary != address(0),
            "Escrow case does not exist"
        );
        _;
    }

    // CONSTRUCTOR

    constructor() Ownable(msg.sender) {}

    // VALIDATOR MANAGEMENT

    function addValidator(
        address _validator
    ) external onlyOwner {

        require(
            _validator != address(0),
            "Invalid validator address"
        );

        validators[_validator] = true;

        emit ValidatorAdded(_validator);
    }

    // ESCROW MANAGEMENT

    function createEscrow(
        uint256 _caseId,
        address payable _beneficiary
    ) external onlyOwner {

        require(
            _beneficiary != address(0),
            "Invalid beneficiary address"
        );

        require(
            escrowCases[_caseId].beneficiary == address(0),
            "Escrow already exists"
        );

        escrowCases[_caseId] = EscrowCase({
            caseId: _caseId,
            beneficiary: _beneficiary,
            totalFunds: 0,
            verified: false,
            released: false,
            createdAt: block.timestamp
        });

        emit EscrowCreated(
            _caseId,
            _beneficiary
        );
    }

    // DONATIONS

    function depositFunds(
        uint256 _caseId
    )
        external
        payable
        nonReentrant
        escrowExists(_caseId)
    {

        require(
            msg.value > 0,
            "Deposit amount must be greater than zero"
        );

        EscrowCase storage escrowCase = escrowCases[_caseId];

        require(
            !escrowCase.released,
            "Funds already released"
        );

        escrowCase.totalFunds += msg.value;

        donorContributions[_caseId][msg.sender] += msg.value;

        emit FundsDeposited(
            _caseId,
            msg.sender,
            msg.value
        );
    }

    // VERIFICATION

    function verifyEscrowCase(
        uint256 _caseId
    )
        external
        onlyValidator
        escrowExists(_caseId)
    {

        EscrowCase storage escrowCase = escrowCases[_caseId];

        require(
            !escrowCase.verified,
            "Case already verified"
        );

        escrowCase.verified = true;

        emit CaseVerified(
            _caseId,
            msg.sender
        );
    }

    // FUND RELEASE

    function releaseFunds(
        uint256 _caseId
    )
        external
        onlyOwner
        nonReentrant
        escrowExists(_caseId)
    {

        EscrowCase storage escrowCase = escrowCases[_caseId];

        require(
            escrowCase.verified,
            "Case is not verified"
        );

        require(
            !escrowCase.released,
            "Funds already released"
        );

        require(
            escrowCase.totalFunds > 0,
            "No funds available"
        );

        uint256 amount = escrowCase.totalFunds;

        escrowCase.released = true;

        (bool success, ) = escrowCase.beneficiary.call{
            value: amount
        }("");

        require(
            success,
            "Transfer failed"
        );

        emit FundsReleased(
            _caseId,
            escrowCase.beneficiary,
            amount
        );
    }

    // VIEW FUNCTIONS

    function getEscrowCase(
        uint256 _caseId
    )
        external
        view
        escrowExists(_caseId)
        returns (EscrowCase memory)
    {
        return escrowCases[_caseId];
    }

    function getContribution(
        uint256 _caseId,
        address _donor
    )
        external
        view
        returns (uint256)
    {
        return donorContributions[_caseId][_donor];
    }

    function getContractBalance()
        external
        view
        returns (uint256)
    {
        return address(this).balance;
    }
}