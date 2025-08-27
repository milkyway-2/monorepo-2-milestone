import React, { useState } from 'react';
import { ethers } from 'ethers';
import { getContractAddress, getNetworkName, CONTRACT_CONFIG } from '../config/contracts';

interface ValidatorReport {
  validatorAddress: string;
  message: string;
}

interface ValidatorReportProps {
  isWalletConnected: boolean;
  walletAddress: string | null;
  walletType: 'metamask' | 'talisman' | null;
}

interface VerificationResult {
  success: boolean;
  signature?: string;
  message?: string;
  fullResponse?: any;
}

export const ValidatorReport: React.FC<ValidatorReportProps> = ({ isWalletConnected, walletAddress, walletType }) => {
  const [report, setReport] = useState<ValidatorReport>({
    validatorAddress: '',
    message: '',
  });
  const [isVerifying, setIsVerifying] = useState(false);
  const [isSubmittingOnChain, setIsSubmittingOnChain] = useState(false);
  const [isSimpleSubmitting, setIsSimpleSubmitting] = useState(false);
  const [submitStatus, setSubmitStatus] = useState<'idle' | 'success' | 'error'>('idle');
  const [verificationResult, setVerificationResult] = useState<VerificationResult | null>(null);
  const [errorMessage, setErrorMessage] = useState<string | null>(null);

  // Update report fields
  const handleInputChange = (field: keyof ValidatorReport, value: string) => {
    setReport(prev => ({ ...prev, [field]: value }));
  };

  // Basic validation for required fields
  const validateReport = (): boolean => {
    if (!report.validatorAddress.trim()) {
      setErrorMessage('Validator address is required');
      return false;
    }
    if (!report.message.trim()) {
      setErrorMessage('Message is required');
      return false;
    }
    return true;
  };

  // Verification logic (hidden in UI but present)
  const verifyReport = async () => {
    if (!isWalletConnected || !walletAddress) {
      setErrorMessage('Please connect your wallet first');
      return;
    }
    if (!validateReport()) {
      return;
    }
    setIsVerifying(true);
    setErrorMessage(null);

    try {
      const response = await fetch('http://localhost:4001/verify', { // Backend verification endpoint
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          validator_address: report.validatorAddress,
          nominator_address: walletAddress,
          msg: report.message,
        }),
      });

      const result = await response.json();

      if (response.ok) {
        setVerificationResult({
          success: true,
          signature: result.signature,
          message: 'Report verified successfully!',
          fullResponse: result,
        });
        setSubmitStatus('success');
      } else {
        setVerificationResult({
          success: false,
          message: result.error || 'Verification failed',
          fullResponse: result,
        });
        setSubmitStatus('error');
      }
    } catch (error) {
      console.error('Error verifying report:', error);
      setVerificationResult({
        success: false,
        message: 'Failed to connect to verification service',
      });
      setSubmitStatus('error');
    } finally {
      setIsVerifying(false);
    }
  };

  // Simulated on-chain submission (placeholder)
  const submitOnChain = async () => {
    if (!verificationResult?.success) {
      setErrorMessage('Please verify the report first');
      return;
    }
    setIsSubmittingOnChain(true);
    setErrorMessage(null);

    try {
      // Simulate delay for on-chain submission
      await new Promise(resolve => setTimeout(resolve, 3000));
      console.log('Submitting verified report on-chain:', {
        validatorAddress: report.validatorAddress,
        nominatorAddress: walletAddress,
        message: report.message,
        signature: verificationResult.signature,
        timestamp: new Date().toISOString(),
      });

      setSubmitStatus('success');
      setVerificationResult(null);
      setReport({ validatorAddress: '', message: '' });
    } catch (error) {
      console.error('Error submitting on-chain:', error);
      setSubmitStatus('error');
      setErrorMessage('Failed to submit on-chain. Please try again.');
    } finally {
      setIsSubmittingOnChain(false);
    }
  };

  // Simple direct submit to smart contract (without verification)
  const simpleSubmit = async () => {
    if (!isWalletConnected || !walletAddress) {
      setErrorMessage('Please connect your wallet first');
      return;
    }
    if (!validateReport()) {
      return;
    }
    setIsSimpleSubmitting(true);
    setErrorMessage(null);
    let nativeSymbol = 'ETH';

    try {
      // Contract ABI with needed methods and events
      const contractABI = [
        { inputs: [], stateMutability: 'nonpayable', type: 'constructor' },
        {
          anonymous: false,
          inputs: [
            { indexed: false, internalType: 'string', name: 'validator', type: 'string' },
            { indexed: false, internalType: 'string', name: 'nominator', type: 'string' },
            { indexed: false, internalType: 'string', name: 'msgText', type: 'string' }
          ],
          name: 'MessageStored',
          type: 'event'
        },
        {
          inputs: [{ internalType: 'uint256', name: 'index', type: 'uint256' }],
          name: 'getMessage',
          outputs: [{
            components: [
              { internalType: 'string', name: 'validator_address', type: 'string' },
              { internalType: 'string', name: 'nominator_address', type: 'string' },
              { internalType: 'string', name: 'msgText', type: 'string' }
            ],
            internalType: 'struct OracleVerifiedDelegation.Message',
            name: '',
            type: 'tuple'
          }],
          stateMutability: 'view',
          type: 'function'
        },
        {
          inputs: [],
          name: 'getMessageCount',
          outputs: [{ internalType: 'uint256', name: '', type: 'uint256' }],
          stateMutability: 'view',
          type: 'function'
        },
        {
          inputs: [{ internalType: 'uint256', name: '', type: 'uint256' }],
          name: 'messages',
          outputs: [
            { internalType: 'string', name: 'validator_address', type: 'string' },
            { internalType: 'string', name: 'nominator_address', type: 'string' },
            { internalType: 'string', name: 'msgText', type: 'string' }
          ],
          stateMutability: 'view',
          type: 'function'
        },
        {
          inputs: [],
          name: 'oracleAddress',
          outputs: [{ internalType: 'address', name: '', type: 'address' }],
          stateMutability: 'view',
          type: 'function'
        },
        {
          inputs: [
            { internalType: 'string', name: 'validator_address', type: 'string' },
            { internalType: 'string', name: 'nominator_address', type: 'string' },
            { internalType: 'string', name: 'msgText', type: 'string' },
            { internalType: 'bytes', name: 'signature', type: 'bytes' }
          ],
          name: 'submitMessage',
          outputs: [],
          stateMutability: 'nonpayable',
          type: 'function'
        },
        {
          inputs: [
            { internalType: 'string', name: 'validator_address', type: 'string' },
            { internalType: 'string', name: 'nominator_address', type: 'string' },
            { internalType: 'string', name: 'msgText', type: 'string' }
          ],
          name: 'submitMessageUnverified',
          outputs: [],
          stateMutability: 'nonpayable',
          type: 'function'
        }
      ];

      // Use the correct provider based on wallet type
      let rawProvider: any = null;
      let accounts: string[] = [];

      console.log('🔍 Setting up provider for wallet type:', walletType);

      if (walletType === 'metamask') {
        if (typeof window.ethereum === 'undefined') {
          throw new Error('MetaMask is not installed. Please install MetaMask to use this feature.');
        }
        rawProvider = window.ethereum;
        console.log('🔍 Requesting MetaMask accounts...');
        accounts = await window.ethereum.request({ method: 'eth_requestAccounts' });
        console.log('✅ MetaMask accounts received:', accounts);
      } else if (walletType === 'talisman') {
        // Use Talisman Connect SDK to get the provider
        const { getWallets } = await import('@talismn/connect-wallets');
        const installedWallets = getWallets().filter((wallet) => wallet.installed);
        const talismanWallet = installedWallets.find(
          (wallet) => wallet.extensionName === 'talisman',
        );
        
        if (!talismanWallet) {
          throw new Error('Talisman is not installed. Please install Talisman to use this feature.');
        }
        
        console.log('🔍 Enabling Talisman wallet...');
        // Ensure Talisman is enabled for this app
        try {
          await talismanWallet.enable('Milkyway2 Portal');
          console.log('✅ Talisman wallet enabled');
        } catch (err) {
          console.log('⚠️ Talisman already enabled or enabling failed:', err);
        }
        
        // Get the EVM provider from Talisman
        console.log('🔍 Getting Talisman accounts...');
        const talismanAccounts = await talismanWallet.getAccounts();
        console.log('✅ Talisman accounts:', talismanAccounts);
        
        if (talismanAccounts.length === 0) {
          throw new Error('No accounts found in Talisman. Please add an account to your wallet.');
        }
        
        // For Talisman, we need to use the EVM provider and request authorization
        if (typeof window.ethereum !== 'undefined') {
          rawProvider = window.ethereum;
          console.log('🔍 Requesting EVM authorization for Talisman...');
          // Request authorization for the specific account
          try {
            const authorizedAccounts = await window.ethereum.request({ 
              method: 'eth_requestAccounts' 
            });
            console.log('✅ EVM authorization successful:', authorizedAccounts);
            if (authorizedAccounts.length > 0) {
              accounts = authorizedAccounts;
            } else {
              accounts = [talismanAccounts[0].address];
            }
          } catch (err) {
            console.log('⚠️ EVM authorization request failed, using existing account:', err);
            accounts = [talismanAccounts[0].address];
          }
        } else {
          throw new Error('Talisman EVM provider not available.');
        }
      } else {
        throw new Error('Unsupported wallet type. Please connect MetaMask or Talisman.');
      }

      if (!rawProvider || accounts.length === 0) {
        throw new Error('No EVM accounts available. Please connect MetaMask/Talisman (EVM) and approve access.');
      }

      console.log('🔍 Creating provider with:', { 
        walletType, 
        rawProvider: rawProvider?.constructor?.name,
        rawProviderType: typeof rawProvider,
        accounts,
        accountCount: accounts.length 
      });
      
      const provider = new ethers.providers.Web3Provider(rawProvider, 'any');
      const providerRequest = async (method: string, params?: any[]) => provider.send(method, params ?? []);
      const account = accounts[0];
      const signer = provider.getSigner(account);
      
      console.log('✅ Provider created successfully with account:', account);
      console.log('🔍 Signer details:', {
        signerType: signer.constructor?.name,
        signerAddress: await signer.getAddress(),
        providerNetwork: await provider.getNetwork(),
      });

      // Determine network and contract address
      const finalChainId = await providerRequest('eth_chainId');
      console.log('🔍 Detected chain ID:', finalChainId);
      
      let contractAddress = '';
      let networkName = '';
      nativeSymbol = 'ETH';

      if (finalChainId === CONTRACT_CONFIG.WESTEND_ASSET_HUB_CHAIN_ID) {
        contractAddress = CONTRACT_CONFIG.WESTEND_ASSET_HUB_CONTRACT_ADDRESS;
        networkName = getNetworkName(finalChainId);
        nativeSymbol = 'WND';
      } else if (finalChainId === CONTRACT_CONFIG.SEPOLIA_CHAIN_ID) {
        contractAddress = CONTRACT_CONFIG.SEPOLIA_CONTRACT_ADDRESS;
        networkName = getNetworkName(finalChainId);
        nativeSymbol = 'ETH';
      } else if (finalChainId === CONTRACT_CONFIG.PASEO_CHAIN_ID) {
        contractAddress = CONTRACT_CONFIG.PASEO_CONTRACT_ADDRESS;
        networkName = getNetworkName(finalChainId);
        nativeSymbol = 'PAS';
      } else {
        throw new Error('Unsupported EVM network. Please switch to a supported network.');
      }
      
      // Verify the wallet is on the correct network for the wallet type
      if (walletType === 'metamask' && finalChainId !== CONTRACT_CONFIG.SEPOLIA_CHAIN_ID && finalChainId !== CONTRACT_CONFIG.PASEO_CHAIN_ID) {
        throw new Error(`MetaMask should be on Sepolia network (${CONTRACT_CONFIG.SEPOLIA_CHAIN_ID}) or Paseo network (${CONTRACT_CONFIG.PASEO_CHAIN_ID}), but is on ${finalChainId}. Please switch networks.`);
      } else if (walletType === 'talisman' && finalChainId !== CONTRACT_CONFIG.PASEO_CHAIN_ID) {
        throw new Error(`Talisman should be on Paseo network (${CONTRACT_CONFIG.PASEO_CHAIN_ID}), but is on ${finalChainId}. Please switch networks.`);
      }
      
      console.log('✅ Network verification passed:', { walletType, finalChainId, contractAddress, networkName });

      const contract = new ethers.Contract(contractAddress, contractABI, signer);

      const hasVerified = !!verificationResult?.success && !!verificationResult?.signature;
      const methodName = hasVerified ? 'submitMessage' : 'submitMessageUnverified';
      const methodArgs = hasVerified
        ? [report.validatorAddress, walletAddress, report.message, verificationResult!.signature]
        : [report.validatorAddress, walletAddress, report.message];

      // Preflight check to detect reverts and get revert reasons
      try {
        // @ts-ignore dynamic method access
        await contract.callStatic[methodName](...methodArgs);
      } catch (preflightError: any) {
        const msg: string =
          preflightError?.data?.message ||
          preflightError?.message ||
          'Contract call reverted during preflight check';

        if (finalChainId === CONTRACT_CONFIG.SEPOLIA_CHAIN_ID) {
          // Abort on Sepolia if revert reason is present
          throw new Error(`Preflight failed: ${msg}`);
        } else if (finalChainId === CONTRACT_CONFIG.PASEO_CHAIN_ID) {
          // Abort on Paseo if revert reason is present
          throw new Error(`Preflight failed: ${msg}`);
        } else {
          console.warn('Preflight eth_call failed on this network (continuing):', msg);
        }
      }

      // Estimate gas
      let gasLimit: ethers.BigNumber | undefined;
      try {
        // @ts-ignore dynamic access
        gasLimit = await contract.estimateGas[methodName](...methodArgs);
      } catch (e) {
        if (finalChainId === CONTRACT_CONFIG.SEPOLIA_CHAIN_ID) {
          throw new Error('Gas estimation failed. The contract may be reverting (e.g., Unauthorized).');
        } else if (finalChainId === CONTRACT_CONFIG.PASEO_CHAIN_ID) {
          throw new Error('Gas estimation failed. The contract may be reverting (e.g., Unauthorized).');
        }
        console.warn('estimateGas failed on this network, using default gas limit');
        gasLimit = ethers.BigNumber.from('500000');
      }

      // Fetch gas price
      let gasPrice: ethers.BigNumber | undefined;
      try {
        gasPrice = await provider.getGasPrice();
      } catch (e) {
        console.warn('getGasPrice failed, using 1 wei as fallback');
        gasPrice = ethers.BigNumber.from(1);
      }

      console.log('🚀 Submitting to smart contract via ethers:', {
        contractAddress,
        method: methodName,
        args: methodArgs,
        gasLimit: gasLimit?.toString(),
        gasPrice: gasPrice?.toString(),
        network: networkName,
        chainId: finalChainId,
        signerAddress: await signer.getAddress(),
      });

      console.log('🔐 About to send transaction - this should trigger wallet signing popup...');
      
      // Send transaction
      // @ts-ignore dynamic method access
      const txResponse = await contract[methodName](...methodArgs, { gasLimit, gasPrice });
      console.log('✅ Transaction sent successfully:', txResponse.hash);
      
      console.log('⏳ Waiting for transaction confirmation...');
      await txResponse.wait();
      console.log('✅ Transaction confirmed on blockchain');

      setSubmitStatus('success');
      setReport({ validatorAddress: '', message: '' });
    } catch (error: any) {
      console.error('Error submitting to smart contract:', error);
      setSubmitStatus('error');

      let errorMsg = 'Failed to submit to smart contract. Please try again.';
      if (error.code === 4001) {
        errorMsg = 'Transaction was rejected by user.';
      } else if (error.code === -32603) {
        errorMsg = 'Internal JSON-RPC error. The contract might not be deployed at this address or it reverted. Check console for details.';
      } else if (error.message?.includes('insufficient funds')) {
        errorMsg = `Insufficient funds for gas. Please add some ${nativeSymbol} to your account.`;
      } else if (error.message?.includes('Preflight failed')) {
        errorMsg = error.message;
      } else if (error.message?.toLowerCase().includes('gas estimation failed')) {
        errorMsg = 'Gas estimation failed. The contract may be reverting (e.g., Unauthorized). Check contract permissions.';
      } else if (error.message?.toLowerCase().includes('network')) {
        errorMsg = 'Network error. Please ensure you are connected to a supported EVM network.';
      } else if (error.message) {
        errorMsg = error.message;
      }

      setErrorMessage(errorMsg);
    } finally {
      setIsSimpleSubmitting(false);
    }
  };

  // UI if wallet not connected
  if (!isWalletConnected) {
    return (
      <div style={{
        background: 'white',
        padding: '2rem',
        borderRadius: '8px',
        boxShadow: '0 2px 4px rgba(0,0,0,0.1)',
        textAlign: 'center'
      }}>
        <h3 style={{ marginBottom: '1rem', color: '#333' }}>Submit Validator Report</h3>
        <p style={{ color: '#666' }}>
          Please connect your wallet to a supported network to submit validator reports.
        </p>
      </div>
    );
  }

  // Main form UI
  return (
    <div style={{
      background: 'white',
      padding: '2rem',
      borderRadius: '8px',
      boxShadow: '0 2px 4px rgba(0,0,0,0.1)'
    }}>
      <h3 style={{ marginBottom: '1rem', color: '#333' }}>Submit Validator Report</h3>
      <p style={{ marginBottom: '2rem', color: '#666' }}>
        Report suspicious or problematic validator behavior.
      </p>

      {submitStatus === 'success' && (
        <div style={{
          background: '#d4edda',
          color: '#155724',
          padding: '0.75rem',
          borderRadius: '4px',
          marginBottom: '1rem',
          border: '1px solid #c3e6cb'
        }}>
          {verificationResult?.success
            ? 'Report verified and submitted successfully!'
            : 'Report submitted successfully! Your report has been recorded on the blockchain.'}
        </div>
      )}

      {errorMessage && (
        <div style={{
          background: '#f8d7da',
          color: '#721c24',
          padding: '0.75rem',
          borderRadius: '4px',
          marginBottom: '1rem',
          border: '1px solid #f5c6cb'
        }}>
          {errorMessage}
        </div>
      )}

      {verificationResult && (
        <div style={{
          background: verificationResult.success ? '#d4edda' : '#f8d7da',
          color: verificationResult.success ? '#155724' : '#721c24',
          padding: '0.75rem',
          borderRadius: '4px',
          marginBottom: '1rem',
          border: `1px solid ${verificationResult.success ? '#c3e6cb' : '#f5c6cb'}`
        }}>
          <div style={{ marginBottom: '0.5rem', fontWeight: 'bold' }}>
            {verificationResult.message}
          </div>
          {verificationResult.fullResponse && (
            <div style={{
              marginTop: '1rem',
              background: '#f8f9fa',
              padding: '1rem',
              borderRadius: '4px',
              border: '1px solid #dee2e6'
            }}>
              <div style={{ marginBottom: '0.5rem', fontWeight: 'bold', fontSize: '0.875rem' }}>
                API Response:
              </div>
              <pre style={{
                background: '#f8f9fa',
                padding: '0.75rem',
                borderRadius: '4px',
                border: '1px solid #dee2e6',
                fontSize: '0.75rem',
                overflow: 'auto',
                whiteSpace: 'pre-wrap',
                wordBreak: 'break-word',
                margin: 0,
                fontFamily: 'monospace'
              }}>
                {JSON.stringify(verificationResult.fullResponse, null, 2)}
              </pre>
            </div>
          )}
        </div>
      )}

      {/* Input fields */}
      <div>
        <div style={{ marginBottom: '1rem' }}>
          <label style={{ display: 'block', marginBottom: '0.5rem', fontWeight: 'bold' }}>
            Validator Address *
          </label>
          <input
            type="text"
            value={report.validatorAddress}
            onChange={e => handleInputChange('validatorAddress', e.target.value)}
            placeholder="Enter validator address (0x...)"
            style={{
              width: '100%',
              padding: '0.75rem',
              border: '1px solid #ddd',
              borderRadius: '4px',
              fontSize: '0.875rem'
            }}
          />
        </div>

        <div style={{ marginBottom: '2rem' }}>
          <label style={{ display: 'block', marginBottom: '0.5rem', fontWeight: 'bold' }}>
            Message *
          </label>
          <textarea
            value={report.message}
            onChange={e => handleInputChange('message', e.target.value)}
            placeholder="Enter your report message..."
            rows={6}
            style={{
              width: '100%',
              padding: '0.75rem',
              border: '1px solid #ddd',
              borderRadius: '4px',
              fontSize: '0.875rem',
              resize: 'vertical'
            }}
          />
        </div>

        <div style={{ display: 'flex', gap: '1rem', flexWrap: 'wrap' }}>
          {/* Verify button is hidden but code preserved */}
          {/*
          <button
            type="button"
            onClick={verifyReport}
            disabled={isVerifying}
            style={{
              background: '#007bff',
              color: 'white',
              border: 'none',
              padding: '0.75rem 1.5rem',
              borderRadius: '4px',
              cursor: isVerifying ? 'not-allowed' : 'pointer',
              fontSize: '1rem',
              fontWeight: 'bold',
              opacity: isVerifying ? 0.6 : 1,
            }}
          >
            {isVerifying ? 'Verifying...' : 'Verify Report'}
          </button>
          */}

          <button
            type="button"
            onClick={simpleSubmit}
            disabled={isSimpleSubmitting || isVerifying}
            style={{
              background: '#dc3545',
              color: 'white',
              border: 'none',
              padding: '0.75rem 1.5rem',
              borderRadius: '4px',
              cursor: isSimpleSubmitting ? 'not-allowed' : 'pointer',
              fontSize: '1rem',
              fontWeight: 'bold',
              opacity: isSimpleSubmitting ? 0.6 : 1,
            }}
          >
            {isSimpleSubmitting ? 'Submitting...' : 'SimpleSubmit'}
          </button>

          {verificationResult?.success && (
            <button
              type="button"
              onClick={submitOnChain}
              disabled={isSubmittingOnChain}
              style={{
                background: '#28a745',
                color: 'white',
                border: 'none',
                padding: '0.75rem 1.5rem',
                borderRadius: '4px',
                cursor: isSubmittingOnChain ? 'not-allowed' : 'pointer',
                fontSize: '1rem',
                fontWeight: 'bold',
                opacity: isSubmittingOnChain ? 0.6 : 1,
              }}
            >
              {isSubmittingOnChain ? 'Submitting...' : 'Submit On-Chain'}
            </button>
          )}
        </div>
      </div>

      {/* Instructions */}
      <div style={{ marginTop: '2rem', padding: '1rem', background: '#f8f9fa', borderRadius: '4px' }}>
        <h4 style={{ marginBottom: '0.5rem', color: '#333' }}>Instructions:</h4>
        <div style={{ fontSize: '0.875rem', color: '#666', lineHeight: '1.6' }}>
          <p style={{ marginBottom: '0.5rem' }}>
            <strong>1. SimpleSubmit:</strong> Directly submits the report to the smart contract using submitMessageUnverified
          </p>
          <p style={{ marginBottom: '0.5rem' }}>
            <strong>2. Submit On-Chain:</strong> Submits the verified report to the blockchain (not implemented yet)
          </p>
          <p style={{ marginBottom: '0.5rem', fontSize: '0.8rem', color: '#888', fontStyle: 'italic' }}>
            <strong>Note:</strong> Verify Report functionality is available in the code but hidden from the UI
          </p>
        </div>
      </div>
    </div>
  );
};
