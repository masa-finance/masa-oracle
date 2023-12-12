// SPDX-License-Identifier: MIT
pragma solidity ^0.8.7;

import "./OracleNodeStakingContract.sol";

contract NodeDataConsensus {
    OracleNodeStakingContract public stakingContract;

    struct NodeData {
        string multiaddr;
        string peerId;
        uint256 lastJoined;
        uint256 lastLeft;
        uint256 lastUpdated;
        uint256 currentUptime;
        uint256 accumulatedUptime;
        uint256 activity;
    }

    // Mapping from period to node data submissions
    mapping(uint256 => mapping(address => NodeData)) public nodeDataSubmissions;

    // Mapping from period to the number of submissions
    mapping(uint256 => uint256) public submissionCount;

    // Mapping from period to consensus data
    mapping(uint256 => NodeData) public consensusData;

    // Required number of nodes for consensus
    uint256 public constant CONSENSUS_THRESHOLD = 5;

    // Event to be emitted when node data is submitted
    event NodeDataSubmitted(address indexed node, uint256 period, NodeData data);

    // Event to be emitted when consensus is reached
    event ConsensusReached(uint256 period, NodeData data);

    constructor(address _stakingContract) {
        stakingContract = OracleNodeStakingContract(_stakingContract);
    }

    function submitNodeData(uint256 period, NodeData calldata data) external {
        require(stakingContract.balanceOf(msg.sender) > 0, "Node is not staked");

        // Check if the node has already submitted data for this period
        require(nodeDataSubmissions[period][msg.sender].lastUpdated == 0, "Data already submitted for this period");

        nodeDataSubmissions[period][msg.sender] = data;
        submissionCount[period]++;

        emit NodeDataSubmitted(msg.sender, period, data);

        // Check if consensus is reached
        if (submissionCount[period] == CONSENSUS_THRESHOLD) {
            consensusData[period] = data;
            emit ConsensusReached(period, data);
        }
    }

    // Any additional functions to handle consensus mechanism, data retrieval, that we may want.
}