// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

/*
 * VerificationRegistry.sol
 * --------------------------------------------------------
 * AidStream Verification Registry
 *
 * Responsibilities:
 * - Register validators
 * - Verify beneficiary cases
 * - Store verification approvals
 * - Maintain transparent verification history
 *
 * Security:
 * - Ownable access control
 * - Validator authorization checks
 * - Duplicate verification prevention
 */

import "@openzeppelin/contracts/access/Ownable.sol";

contract VerificationRegistry is Ownable {

    // STRUCTS

    struct VerificationCase {
        uint256 caseId;
        bool verified;
        uint256 approvalCount;
        uint256 createdAt;
    }

    struct VerificationRecord {
        address validator;
        uint256 timestamp;
    }

    // STATE VARIABLES

    mapping(address => bool) public validators;

    mapping(uint256 => VerificationCase) public verificationCases;

    mapping(uint256 => VerificationRecord[]) public verificationHistory;

    mapping(uint256 => mapping(address => bool))
        public hasVerified;

    uint256 public validatorCount;

    // EVENTS

    event ValidatorAdded(address indexed validator);

    event ValidatorRemoved(address indexed validator);

    event CaseCreated(uint256 indexed caseId);

    event CaseVerified(
        uint256 indexed caseId,
        address indexed validator
    );

    event CaseApproved(
        uint256 indexed caseId
    );

    // MODIFIERS

    modifier onlyValidator() {
        require(
            validators[msg.sender],
            "Not an authorized validator"
        );
        _;
    }

    modifier caseExists(uint256 _caseId) {
        require(
            verificationCases[_caseId].createdAt > 0,
            "Verification case does not exist"
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

        require(
            !validators[_validator],
            "Validator already exists"
        );

        validators[_validator] = true;

        validatorCount++;

        emit ValidatorAdded(_validator);
    }

    function removeValidator(
        address _validator
    ) external onlyOwner {

        require(
            validators[_validator],
            "Validator does not exist"
        );

        validators[_validator] = false;

        validatorCount--;

        emit ValidatorRemoved(_validator);
    }

    // CASE MANAGEMENT

    function createVerificationCase(
        uint256 _caseId
    ) external onlyOwner {

        require(
            verificationCases[_caseId].createdAt == 0,
            "Case already exists"
        );

        verificationCases[_caseId] = VerificationCase({
            caseId: _caseId,
            verified: false,
            approvalCount: 0,
            createdAt: block.timestamp
        });

        emit CaseCreated(_caseId);
    }

    // VERIFICATION LOGIC

    function verifyCase(
        uint256 _caseId
    )
        external
        onlyValidator
        caseExists(_caseId)
    {

        VerificationCase storage verificationCase =
            verificationCases[_caseId];

        require(
            !verificationCase.verified,
            "Case already approved"
        );

        require(
            !hasVerified[_caseId][msg.sender],
            "Validator already approved this case"
        );

        hasVerified[_caseId][msg.sender] = true;

        verificationCase.approvalCount++;

        verificationHistory[_caseId].push(
            VerificationRecord({
                validator: msg.sender,
                timestamp: block.timestamp
            })
        );

        emit CaseVerified(
            _caseId,
            msg.sender
        );

        /*
         * Approval Threshold
         * ------------------------------------------------
         * MVP Rule:
         * Case becomes verified after 2 approvals.
         *
         * This can later become configurable.
         */

        if (verificationCase.approvalCount >= 2) {

            verificationCase.verified = true;

            emit CaseApproved(_caseId);
        }
    }

    // VIEW FUNCTIONS

    function isCaseVerified(
        uint256 _caseId
    )
        external
        view
        caseExists(_caseId)
        returns (bool)
    {
        return verificationCases[_caseId].verified;
    }

    function getApprovalCount(
        uint256 _caseId
    )
        external
        view
        caseExists(_caseId)
        returns (uint256)
    {
        return verificationCases[_caseId].approvalCount;
    }

    function getVerificationHistory(
        uint256 _caseId
    )
        external
        view
        caseExists(_caseId)
        returns (VerificationRecord[] memory)
    {
        return verificationHistory[_caseId];
    }
}