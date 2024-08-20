// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import "@openzeppelin/contracts/utils/Counters.sol";
import "./NFT.sol";
import "./Token.sol";

contract Controller {
    using Counters for Counters.Counter;

    //
    // STATE VARIABLES
    //
    Counters.Counter private _sessionIdCounter;
    GeneNFT public geneNFT;
    PostCovidStrokePrevention public pcspToken;

    struct UploadSession {
        uint256 id;
        address user;
        string proof;
        bool confirmed;
    }

    struct DataDoc {
        string id;
        string hashContent;
    }

    mapping(uint256 => UploadSession) sessions;
    mapping(string => DataDoc) docs;
    mapping(string => bool) docSubmits;
    mapping(uint256 => string) nftDocs;

    //
    // EVENTS
    //
    event UploadData(string docId, uint256 sessionId);

    constructor(address nftAddress, address pcspAddress) {
        geneNFT = GeneNFT(nftAddress);
        pcspToken = PostCovidStrokePrevention(pcspAddress);
    }

    function uploadData(string memory docId) public returns (uint256) {
        // TODO: Implement this method: to start an uploading gene data session. The doc id is used to identify a unique gene profile. Also should check if the doc id has been submited to the system before. This method return the session id
    }

    function confirm(
        string memory docId,
        string memory contentHash,
        string memory proof,
        uint256 sessionId,
        uint256 riskScore
    ) public {
        // TODO: Implement this method: The proof here is used to verify that the result is returned from a valid computation on the gene data. For simplicity, we will skip the proof verification in this implementation. The gene data's owner will receive a NFT as a ownership certicate for his/her gene profile.

        // TODO: Verify proof, we can skip this step

        // TODO: Update doc content

        // TODO: Mint NFT 

        // TODO: Reward PCSP token based on risk stroke

        // TODO: Close session
    }

    function getSession(uint256 sessionId) public view returns(UploadSession memory) {
        return sessions[sessionId];
    }

    function getDoc(string memory docId) public view returns(DataDoc memory) {
        return docs[docId];
    }
}
