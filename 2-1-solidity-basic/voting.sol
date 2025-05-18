/*
 创建一个名为Voting的合约，包含以下功能：,
,
一个mapping来存储候选人的得票数,
一个vote函数，允许用户投票给某个候选人,
一个getVotes函数，返回某个候选人的得票数,
一个resetVotes函数，重置所有候选人的得票数
*/

// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract Voting {
    mapping(string => uint) private votes;
    string[] private candidateList;

    function vote(string memory candidate) public {
        require(bytes(candidate).length > 0, "Candidate name cannot be empty");
        if (votes[candidate] == 0 && !contains(candidate)) {
            candidateList.push(candidate);
        }
        votes[candidate] += 1;
    }

    function getVotes(string memory candidate) public view returns (uint) {
        return votes[candidate];
    }

    function resetVotes() public {
        for (uint i = 0; i < candidateList.length; i++) {
            votes[candidateList[i]] = 0;
        }
    }

    // 辅助函数：检查候选人是否已存在
    function contains(string memory candidate) private view returns (bool) {
        for (uint i = 0; i < candidateList.length; i++) {
            if (keccak256(bytes(candidateList[i])) == keccak256(bytes(candidate))) {
                return true;
            }
        }
        return false;
    }


    /*
        反转字符串 (Reverse String),
        题目描述：反转一个字符串。输入 "abcde"，输出 "edcba"
    */

    function reverseString(string memory input) public pure returns (string memory) {
        bytes memory strBytes = bytes(input);
        uint length = strBytes.length;

        // 如果字符串长度小于等于1，直接返回原字符串
        if (length <= 1) {
            return input;
        }

        // 创建一个新的字节数组用于存储反转后的字符串
        bytes memory reversed = new bytes(length);

        // 反转字符串
        for (uint i = 0; i < length; i++) {
            reversed[i] = strBytes[length - 1 - i];
        }

        return string(reversed);
    }


    /*
        用 solidity 实现整数转罗马数字;
        题目描述在 https://leetcode.cn/problems/roman-to-integer/description/3.
    */

    mapping(bytes1 => int256) private symbolValues;

    constructor() {
        symbolValues[bytes1('I')] = 1;
        symbolValues[bytes1('V')] = 5;
        symbolValues[bytes1('X')] = 10;
        symbolValues[bytes1('L')] = 50;
        symbolValues[bytes1('C')] = 100;
        symbolValues[bytes1('D')] = 500;
        symbolValues[bytes1('M')] = 1000;
    }

    function romanToInt(string memory s) public view returns (int256 ans) {
        bytes memory b = bytes(s);
        uint n = b.length;
        ans = 0;
        for (uint i = 0; i < n; i++) {
            int256 v = symbolValues[b[i]];
            if (i+1 < n && v < symbolValues[b[i+1]]) ans -= v;
            else ans += v;
        }
    }

    /*
        合并两个有序数组 (Merge Sorted Array),
        题目描述：将两个有序数组合并为一个有序数组。
    */

    function mergeSortedArrays(int[] memory arr1, int[] memory arr2) public pure returns (int[] memory) {
        uint len1 = arr1.length;
        uint len2 = arr2.length;
        uint i = 0;
        uint j = 0;
        uint k = 0;

        int[] memory merged = new int[](len1 + len2);

        while (i < len1 && j < len2) {
            if (arr1[i] <= arr2[j]) {
                merged[k] = arr1[i];
                i++;
            } else {
                merged[k] = arr2[j];
                j++;
            }
            k++;
        }

        while (i < len1) {
            merged[k] = arr1[i];
            i++;
            k++;
        }

        while (j < len2) {
            merged[k] = arr2[j];
            j++;
            k++;
        }

        return merged;
    }

    /*
        二分查找 (Binary Search), 题目描述：在一个有序数组中查找目标值。
    */

    function binarySearch(int[] memory arr, int target) public pure returns (int) {
        int left = 0;
        int right = int(arr.length) - 1;

        while (left <= right) {
            int mid = (left + right) / 2;
            if (arr[uint(mid)] == target) {
                return mid;
            } else if (arr[uint(mid)] < target) {
                left = mid + 1;
            } else {
                right = mid - 1;
            }
        }

        return -1;
    }


}



