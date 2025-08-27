import React, { useState, useEffect } from 'react';
import BigNumber from 'bignumber.js';
import { CONTRACT_CONFIG } from '../config/contracts';
import { getWallets } from '@talismn/connect-wallets';
import { Modal } from '@talismn/connect-ui';

interface WalletState {
  isConnected: boolean;
  address: string | null;
  chainId: string | null;
  balance: string | null;
  walletType: 'metamask' | 'talisman' | null;
  unsubscribe?: () => void;
}

interface WalletConnectProps {
  onWalletStateChange?: (isConnected: boolean, address: string | null, walletType?: 'metamask' | 'talisman' | null) => void;
}

export const WalletConnect: React.FC<WalletConnectProps> = ({ onWalletStateChange }) => {
  const [walletState, setWalletState] = useState<WalletState>({
    isConnected: false,
    address: null,
    chainId: null,
    balance: null,
    walletType: null,
  });
  const [isConnecting, setIsConnecting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [isModalOpen, setIsModalOpen] = useState(false);

  // Call the callback whenever wallet state changes
  useEffect(() => {
    console.log('🔄 Wallet state changed:', { 
      isConnected: walletState.isConnected, 
      address: walletState.address,
      walletType: walletState.walletType 
    });
    
    if (onWalletStateChange) {
      console.log('📞 Calling onWalletStateChange callback...');
      onWalletStateChange(walletState.isConnected, walletState.address, walletState.walletType);
    }
  }, [walletState.isConnected, walletState.address, walletState.walletType, onWalletStateChange]);

  // Sepolia (Ethereum testnet) RPC + token decimals
  const SEPOLIA_RPC = 'https://rpc.sepolia.org';
  const SEPOLIA_DECIMALS = 18;
  const SEPOLIA_CHAIN_ID = '0xaa36a7';
  
  // Paseo EVM RPC + token decimals
  const PASEO_RPC = 'https://paseo-asset-hub-eth-rpc.polkadot.io';
  const PASEO_DECIMALS = 18;

  // Fetch balance from an EVM RPC
  const fetchEvmBalance = async (address: string, rpcUrl: string, decimals: number): Promise<string | null> => {
    try {
      const body = {
        jsonrpc: '2.0',
        method: 'eth_getBalance',
        params: [address, 'latest'],
        id: 1,
      };

      const res = await fetch(rpcUrl, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(body),
      });
      const json = await res.json();
      const hex = json?.result as string | undefined;
      if (!hex) {
        console.warn('eth_getBalance returned no result', json);
        return null;
      }
      // Convert hex to decimal and then to WND with 12 decimals
      const wei = new BigNumber(BigInt(hex).toString());
      const amount = wei.dividedBy(new BigNumber(10).pow(decimals));
      // Format with up to 6 fractional digits, trim trailing zeros
      return amount.toFormat(6);
    } catch (e) {
      console.error('Failed to fetch EVM balance:', e);
      return null;
    }
  };



  const updateWalletState = async (address: string, chainId: string, walletType: 'metamask' | 'talisman') => {
    console.log('🔄 Starting wallet state update...');
    const startTime = Date.now();
    
    try {
      console.log('📋 Updating wallet state with:', { address, chainId, walletType });

      // Immediately set basic state so UI updates without waiting for balance
      setWalletState({
        isConnected: true,
        address,
        chainId,
        balance: null,
        walletType,
      });

      // Fetch balance asynchronously to avoid hanging the UI
      (async () => {
        try {
          let fetchedBalance: string | null = null;
          if (walletType === 'metamask' && address && chainId === SEPOLIA_CHAIN_ID) {
            fetchedBalance = await Promise.race([
              fetchEvmBalance(address, SEPOLIA_RPC, SEPOLIA_DECIMALS),
              new Promise<string | null>(resolve => setTimeout(() => resolve(null), 4000)),
            ]);
          } else if (walletType === 'talisman' && address && chainId === CONTRACT_CONFIG.PASEO_CHAIN_ID) {
            fetchedBalance = await Promise.race([
              fetchEvmBalance(address, PASEO_RPC, PASEO_DECIMALS),
              new Promise<string | null>(resolve => setTimeout(() => resolve(null), 4000)),
            ]);
          }

          // Only update if the address/chainId are still current
          setWalletState(prev => {
            if (!prev.isConnected) return prev;
            if (prev.address !== address || prev.chainId !== chainId || prev.walletType !== walletType) return prev;
            return { ...prev, balance: fetchedBalance };
          });
        } catch (balanceErr) {
          console.warn('Failed to fetch balance (non-fatal):', balanceErr);
        }
      })();

      const endTime = Date.now();
      console.log(`⏱️ Wallet state update completed in ${endTime - startTime}ms`);
    } catch (err) {
      console.error('❌ Error updating wallet state:', err);
      const errorTime = Date.now();
      console.log(`⏱️ Wallet state update failed after ${errorTime - startTime}ms`);
    }
  };

  



  const connectSepolia = async () => {
    console.log('🚀 Starting MetaMask connection (Sepolia)...');
    const startTime = Date.now();
    setIsConnecting(true);
    setError(null);

    try {
      if (typeof window.ethereum === 'undefined') {
        setError('MetaMask is not installed. Please install MetaMask to use this feature.');
        return;
      }

      console.log('📡 Requesting MetaMask accounts...');
      const accounts = await window.ethereum.request({ method: 'eth_requestAccounts' });
      let currentChainId = await window.ethereum.request({ method: 'eth_chainId' });

      if (currentChainId !== SEPOLIA_CHAIN_ID) {
        console.log('🔄 Switching to Sepolia network...');
        try {
          await window.ethereum.request({
            method: 'wallet_switchEthereumChain',
            params: [{ chainId: SEPOLIA_CHAIN_ID }],
          });
        } catch (switchError: any) {
          if (switchError.code === 4902) {
            console.log('➕ Adding Sepolia to MetaMask...');
            await window.ethereum.request({
              method: 'wallet_addEthereumChain',
              params: [
                {
                  chainId: SEPOLIA_CHAIN_ID,
                  chainName: 'Sepolia',
                  nativeCurrency: { name: 'Sepolia ETH', symbol: 'ETH', decimals: 18 },
                  rpcUrls: [SEPOLIA_RPC],
                  blockExplorerUrls: ['https://sepolia.etherscan.io'],
                },
              ],
            });
            // After adding, attempt switch again for reliability
            await window.ethereum.request({
              method: 'wallet_switchEthereumChain',
              params: [{ chainId: SEPOLIA_CHAIN_ID }],
            });
          } else {
            throw switchError;
          }
        }
        // Re-read the chain after switch/add
        currentChainId = await window.ethereum.request({ method: 'eth_chainId' });
        if (currentChainId !== SEPOLIA_CHAIN_ID) {
          throw new Error('Failed to switch to Sepolia. Please approve the network switch in your wallet.');
        }
      }

      await updateWalletState(accounts[0], SEPOLIA_CHAIN_ID, 'metamask');
      const endTime = Date.now();
      console.log(`⏱️ MetaMask Sepolia connection completed in ${endTime - startTime}ms`);
    } catch (err: any) {
      console.error('❌ Error connecting MetaMask (Sepolia):', err);
      const errorTime = Date.now();
      console.log(`⏱️ MetaMask Sepolia connection failed after ${errorTime - startTime}ms`);
      setError(err.message || 'Failed to connect MetaMask (Sepolia)');
    } finally {
      setIsConnecting(false);
    }
  };

  const connectPaseo = async () => {
    console.log('🚀 Starting MetaMask connection (Paseo)...');
    const startTime = Date.now();
    setIsConnecting(true);
    setError(null);

    try {
      if (typeof window.ethereum === 'undefined') {
        setError('MetaMask is not installed. Please install MetaMask to use this feature.');
        return;
      }

      console.log('📡 Requesting MetaMask accounts...');
      const accounts = await window.ethereum.request({ method: 'eth_requestAccounts' });
      let currentChainId = await window.ethereum.request({ method: 'eth_chainId' });

      if (currentChainId !== CONTRACT_CONFIG.PASEO_CHAIN_ID) {
        console.log('🔄 Switching to Paseo network...');
        try {
          await window.ethereum.request({
            method: 'wallet_switchEthereumChain',
            params: [{ chainId: CONTRACT_CONFIG.PASEO_CHAIN_ID }],
          });
        } catch (switchError: any) {
          if (switchError.code === 4902) {
            console.log('➕ Adding Paseo to MetaMask...');
            await window.ethereum.request({
              method: 'wallet_addEthereumChain',
              params: [
                {
                  chainId: CONTRACT_CONFIG.PASEO_CHAIN_ID,
                  chainName: 'Paseo',
                  nativeCurrency: { name: 'Paseo PAS', symbol: 'PAS', decimals: 18 },
                  rpcUrls: [PASEO_RPC],
                  blockExplorerUrls: ['https://polkadot.js.org/apps/?rpc=wss://paseo-rpc.dwellir.com#/explorer'],
                },
              ],
            });
            // After adding, attempt switch again for reliability
            await window.ethereum.request({
              method: 'wallet_switchEthereumChain',
              params: [{ chainId: CONTRACT_CONFIG.PASEO_CHAIN_ID }],
            });
          } else {
            throw switchError;
          }
        }
        // Re-read the chain after switch/add
        currentChainId = await window.ethereum.request({ method: 'eth_chainId' });
        if (currentChainId !== CONTRACT_CONFIG.PASEO_CHAIN_ID) {
          throw new Error('Failed to switch to Paseo. Please approve the network switch in your wallet.');
        }
      }

      await updateWalletState(accounts[0], CONTRACT_CONFIG.PASEO_CHAIN_ID, 'metamask');
      const endTime = Date.now();
      console.log(`⏱️ MetaMask Paseo connection completed in ${endTime - startTime}ms`);
    } catch (err: any) {
      console.error('❌ Error connecting MetaMask (Paseo):', err);
      const errorTime = Date.now();
      console.log(`⏱️ MetaMask Paseo connection failed after ${errorTime - startTime}ms`);
      setError(err.message || 'Failed to connect MetaMask (Paseo)');
    } finally {
      setIsConnecting(false);
    }
  };







  const disconnectWallet = () => {
    // Call unsubscribe function if it exists
    if (walletState.unsubscribe) {
      walletState.unsubscribe();
    }
    
    setWalletState({
      isConnected: false,
      address: null,
      chainId: null,
      balance: null,
      walletType: null,
      unsubscribe: undefined,
    });
  };

  // Handle wallet connection using @talismn/connect-wallets
  const handleWalletConnect = async (wallet: any) => {
    try {
      console.log('🔗 Connecting to wallet:', wallet.title);
      setIsConnecting(true);
      setError(null);

      // Enable the wallet
      await wallet.enable('Milkyway2 Portal');
      console.log('✅ Wallet enabled:', wallet.title);

      // Subscribe to accounts
      const unsubscribe = await wallet.subscribeAccounts((accounts: any[]) => {
        if (accounts.length > 0) {
          console.log('✅ Accounts received:', accounts);
          
          // Determine wallet type and chain ID
          const walletType = wallet.extensionName === 'talisman' ? 'talisman' : 'metamask';
          const chainId = walletType === 'talisman' ? CONTRACT_CONFIG.PASEO_CHAIN_ID : SEPOLIA_CHAIN_ID;
          
          updateWalletState(accounts[0].address, chainId, walletType);
        }
      });

      // Store unsubscribe function for cleanup
      setWalletState(prev => ({ ...prev, unsubscribe }));

    } catch (err: any) {
      console.error('❌ Error connecting wallet:', err);
      setError(err.message || 'Failed to connect wallet');
    } finally {
      setIsConnecting(false);
    }
  };

  const formatAddress = (address: string) => {
    return `${address.slice(0, 6)}...${address.slice(-4)}`;
  };

  const getNetworkName = () => {
    if (walletState.walletType === 'metamask' && walletState.chainId === SEPOLIA_CHAIN_ID) {
      return 'Sepolia (EVM)';
    } else if (walletState.walletType === 'talisman' && walletState.chainId === CONTRACT_CONFIG.PASEO_CHAIN_ID) {
      return 'Paseo (EVM)';
    }
    return 'Unknown';
  };

  const getBalanceUnit = () => {
    if (walletState.walletType === 'metamask' && walletState.chainId === SEPOLIA_CHAIN_ID) {
      return 'ETH';
    } else if (walletState.walletType === 'talisman' && walletState.chainId === CONTRACT_CONFIG.PASEO_CHAIN_ID) {
      return 'PAS';
    }
    return '';
  };

  const getWalletName = () => {
    if (walletState.walletType === 'metamask') {
      return 'MetaMask';
    } else if (walletState.walletType === 'talisman') {
      return 'Talisman';
    }
    return 'Unknown';
  };

  return (
    <div style={{
      background: 'white',
      padding: '2rem',
      borderRadius: '8px',
      boxShadow: '0 2px 4px rgba(0,0,0,0.1)',
      marginBottom: '2rem'
    }}>
      <h3 style={{ marginBottom: '1rem', color: '#333' }}>Wallet Connection</h3>
      
      {error && (
        <div style={{
          background: '#f8d7da',
          color: '#721c24',
          padding: '0.75rem',
          borderRadius: '4px',
          marginBottom: '1rem',
          border: '1px solid #f5c6cb'
        }}>
          {error}
        </div>
      )}

      {!walletState.isConnected ? (
        <div>
          <p style={{ marginBottom: '1rem', color: '#666' }}>
            Connect your wallet to submit validator reports.
          </p>
          
          {/* MetaMask Sepolia Button */}
          <div style={{ marginBottom: '1rem' }}>
            <p style={{ marginBottom: '0.5rem', color: '#666', fontSize: '0.875rem' }}>
              <strong>For Sepolia Network:</strong>
            </p>
            <button
              onClick={connectSepolia}
              disabled={isConnecting}
              style={{
                background: '#f6851b',
                color: 'white',
                border: 'none',
                padding: '0.75rem 1rem',
                borderRadius: '8px',
                cursor: isConnecting ? 'not-allowed' : 'pointer',
                fontSize: '0.875rem',
                fontWeight: '600',
                opacity: isConnecting ? 0.6 : 1,
                width: '100%',
                maxWidth: '300px',
                minHeight: '48px',
                transition: 'all 0.2s ease',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                gap: '0.5rem',
              }}
            >
              <svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                <path d="M21.49 2L13.54 8.27L15.09 4.68L21.49 2Z" fill="#E2761B" stroke="#E2761B" strokeLinecap="round" strokeLinejoin="round"/>
                <path d="M2.51 2L10.36 8.39L8.91 4.68L2.51 2Z" fill="#E4761B" stroke="#E4761B" strokeLinecap="round" strokeLinejoin="round"/>
                <path d="M18.82 16.52L16.91 19.59L21.09 20.82L22.49 16.68L18.82 16.52Z" fill="#E4761B" stroke="#E4761B" strokeLinecap="round" strokeLinejoin="round"/>
                <path d="M1.59 16.68L2.99 20.82L7.17 19.59L5.26 16.52L1.59 16.68Z" fill="#E4761B" stroke="#E4761B" strokeLinecap="round" strokeLinejoin="round"/>
                <path d="M6.82 10.68L5.32 12.18L8.82 12.48L8.62 8.48L6.82 10.68Z" fill="#E4761B" stroke="#E4761B" strokeLinecap="round" strokeLinejoin="round"/>
                <path d="M17.18 10.68L15.32 8.38L15.18 12.48L18.68 12.18L17.18 10.68Z" fill="#E4761B" stroke="#E4761B" strokeLinecap="round" strokeLinejoin="round"/>
                <path d="M7.17 19.59L10.67 17.89L7.67 15.59L7.17 19.59Z" fill="#E4761B" stroke="#E4761B" strokeLinecap="round" strokeLinejoin="round"/>
                <path d="M13.33 17.89L16.83 19.59L16.33 15.59L13.33 17.89Z" fill="#E4761B" stroke="#E4761B" strokeLinecap="round" strokeLinejoin="round"/>
                <path d="M16.83 19.59L13.33 17.89L13.63 20.39L13.59 22.69L16.83 19.59Z" fill="#D05C15" stroke="#D05C15" strokeLinecap="round" strokeLinejoin="round"/>
                <path d="M7.17 19.59L10.41 22.69L10.37 20.39L10.67 17.89L7.17 19.59Z" fill="#D05C15" stroke="#D05C15" strokeLinecap="round" strokeLinejoin="round"/>
                <path d="M10.47 14.68L7.47 13.88L9.87 13.18L10.47 14.68Z" fill="#233447" stroke="#233447" strokeLinecap="round" strokeLinejoin="round"/>
                <path d="M13.53 14.68L14.13 13.18L16.53 13.88L13.53 14.68Z" fill="#233447" stroke="#233447" strokeLinecap="round" strokeLinejoin="round"/>
                <path d="M7.17 19.59L7.77 16.52L5.26 16.68L7.17 19.59Z" fill="#CD6116" stroke="#CD6116" strokeLinecap="round" strokeLinejoin="round"/>
                <path d="M16.23 16.52L16.83 19.59L18.74 16.68L16.23 16.52Z" fill="#CD6116" stroke="#CD6116" strokeLinecap="round" strokeLinejoin="round"/>
                <path d="M18.68 12.18L15.18 12.48L13.53 14.68L10.47 14.68L8.82 12.48L5.32 12.18L8.82 10.68L11.82 12.18L12.02 12.18L15.02 10.68L18.68 12.18Z" fill="#CD6116" stroke="#CD6116" strokeLinecap="round" strokeLinejoin="round"/>
                <path d="M5.32 12.18L8.82 12.48L8.82 10.68L11.82 12.18L12.02 12.18L15.02 10.68L8.82 10.68L5.32 12.18Z" fill="#E4751F" stroke="#E4751F" strokeLinecap="round" strokeLinejoin="round"/>
                <path d="M6.82 10.68L8.82 12.48L8.82 10.68L6.82 10.68Z" fill="#F6851B" stroke="#F6851B" strokeLinecap="round" strokeLinejoin="round"/>
                <path d="M15.18 10.68L17.18 12.48L15.18 10.68Z" fill="#F6851B" stroke="#F6851B" strokeLinecap="round" strokeLinejoin="round"/>
                <path d="M7.17 19.59L7.67 15.59L7.17 19.59Z" fill="#E4751F" stroke="#E4751F" strokeLinecap="round" strokeLinejoin="round"/>
                <path d="M16.33 15.59L16.83 19.59L16.33 15.59Z" fill="#E4751F" stroke="#E4751F" strokeLinecap="round" strokeLinejoin="round"/>
                <path d="M18.68 12.18L15.18 12.48L13.53 14.68L10.47 14.68L8.82 12.48L5.32 12.18L8.82 10.68L11.82 12.18L12.02 12.18L15.02 10.68L18.68 12.18Z" fill="#F6851B" stroke="#F6851F" strokeLinecap="round" strokeLinejoin="round"/>
              </svg>
              Connect MetaMask (Sepolia)
            </button>
          </div>

          {/* MetaMask Paseo Button */}
          <div style={{ marginBottom: '1rem' }}>
            <p style={{ marginBottom: '0.5rem', color: '#666', fontSize: '0.875rem' }}>
              <strong>For Paseo Network (MetaMask):</strong>
            </p>
            <button
              onClick={connectPaseo}
              disabled={isConnecting}
              style={{
                background: '#f6851b',
                color: 'white',
                border: 'none',
                padding: '0.75rem 1rem',
                borderRadius: '8px',
                cursor: isConnecting ? 'not-allowed' : 'pointer',
                fontSize: '0.875rem',
                fontWeight: '600',
                opacity: isConnecting ? 0.6 : 1,
                width: '100%',
                maxWidth: '300px',
                minHeight: '48px',
                transition: 'all 0.2s ease',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                gap: '0.5rem',
              }}
            >
              <svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                <path d="M21.49 2L13.54 8.27L15.09 4.68L21.49 2Z" fill="#E2761B" stroke="#E2761B" strokeLinecap="round" strokeLinejoin="round"/>
                <path d="M2.51 2L10.36 8.39L8.91 4.68L2.51 2Z" fill="#E4761B" stroke="#E4761B" strokeLinecap="round" strokeLinejoin="round"/>
                <path d="M18.82 16.52L16.91 19.59L21.09 20.82L22.49 16.68L18.82 16.52Z" fill="#E4761B" stroke="#E4761B" strokeLinecap="round" strokeLinejoin="round"/>
                <path d="M1.59 16.68L2.99 20.82L7.17 19.59L5.26 16.52L1.59 16.68Z" fill="#E4761B" stroke="#E4761B" strokeLinecap="round" strokeLinejoin="round"/>
                <path d="M6.82 10.68L5.32 12.18L8.82 12.48L8.62 8.48L6.82 10.68Z" fill="#E4761B" stroke="#E4761B" strokeLinecap="round" strokeLinejoin="round"/>
                <path d="M17.18 10.68L15.32 8.38L15.18 12.48L18.68 12.18L17.18 10.68Z" fill="#E4761B" stroke="#E4761B" strokeLinecap="round" strokeLinejoin="round"/>
                <path d="M7.17 19.59L10.67 17.89L7.67 15.59L7.17 19.59Z" fill="#E4761B" stroke="#E4761B" strokeLinecap="round" strokeLinejoin="round"/>
                <path d="M13.33 17.89L16.83 19.59L16.33 15.59L13.33 17.89Z" fill="#E4761B" stroke="#E4761B" strokeLinecap="round" strokeLinejoin="round"/>
                <path d="M16.83 19.59L13.33 17.89L13.63 20.39L13.59 22.69L16.83 19.59Z" fill="#D05C15" stroke="#D05C15" strokeLinecap="round" strokeLinejoin="round"/>
                <path d="M7.17 19.59L10.41 22.69L10.37 20.39L10.67 17.89L7.17 19.59Z" fill="#D05C15" stroke="#D05C15" strokeLinecap="round" strokeLinejoin="round"/>
                <path d="M10.47 14.68L7.47 13.88L9.87 13.18L10.47 14.68Z" fill="#233447" stroke="#233447" strokeLinecap="round" strokeLinejoin="round"/>
                <path d="M13.53 14.68L14.13 13.18L16.53 13.88L13.53 14.68Z" fill="#233447" stroke="#233447" strokeLinecap="round" strokeLinejoin="round"/>
                <path d="M7.17 19.59L7.77 16.52L5.26 16.68L7.17 19.59Z" fill="#CD6116" stroke="#CD6116" strokeLinecap="round" strokeLinejoin="round"/>
                <path d="M16.23 16.52L16.83 19.59L18.74 16.68L16.23 16.52Z" fill="#CD6116" stroke="#CD6116" strokeLinecap="round" strokeLinejoin="round"/>
                <path d="M18.68 12.18L15.18 12.48L13.53 14.68L10.47 14.68L8.82 12.48L5.32 12.18L8.82 10.68L11.82 12.18L12.02 12.18L15.02 10.68L18.68 12.18Z" fill="#CD6116" stroke="#CD6116" strokeLinecap="round" strokeLinejoin="round"/>
                <path d="M5.32 12.18L8.82 12.48L8.82 10.68L11.82 12.18L12.02 12.18L15.02 10.68L8.82 10.68L5.32 12.18Z" fill="#E4751F" stroke="#E4751F" strokeLinecap="round" strokeLinejoin="round"/>
                <path d="M6.82 10.68L8.82 12.48L8.82 10.68L6.82 10.68Z" fill="#F6851B" stroke="#F6851B" strokeLinecap="round" strokeLinejoin="round"/>
                <path d="M15.18 10.68L17.18 12.48L15.18 10.68Z" fill="#F6851B" stroke="#F6851B" strokeLinecap="round" strokeLinejoin="round"/>
                <path d="M7.17 19.59L7.67 15.59L7.17 19.59Z" fill="#E4751F" stroke="#E4751F" strokeLinecap="round" strokeLinejoin="round"/>
                <path d="M16.33 15.59L16.83 19.59L16.33 15.59Z" fill="#E4751F" stroke="#E4751F" strokeLinecap="round" strokeLinejoin="round"/>
                <path d="M18.68 12.18L15.18 12.48L13.53 14.68L10.47 14.68L8.82 12.48L5.32 12.18L8.82 10.68L11.82 12.18L12.02 12.18L15.02 10.68L18.68 12.18Z" fill="#F6851B" stroke="#F6851F" strokeLinecap="round" strokeLinejoin="round"/>
              </svg>
              Connect MetaMask (Paseo)
            </button>
          </div>
          
          {/* Wallet Selector for Paseo Network */}
          <div style={{ marginBottom: '1rem' }}>
            <p style={{ marginBottom: '0.5rem', color: '#666', fontSize: '0.875rem' }}>
              <strong>For Paseo Network (Talisman):</strong>
            </p>
            <button
              onClick={() => setIsModalOpen(true)}
              disabled={isConnecting}
              style={{
                background: '#6366f1',
                color: 'white',
                border: 'none',
                padding: '0.75rem 1rem',
                borderRadius: '8px',
                cursor: isConnecting ? 'not-allowed' : 'pointer',
                fontSize: '0.875rem',
                fontWeight: '600',
                opacity: isConnecting ? 0.6 : 1,
                width: '100%',
                maxWidth: '300px',
                minHeight: '48px',
                transition: 'all 0.2s ease',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                gap: '0.5rem',
              }}
            >
              <svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                <path d="M12 2L2 7L12 12L22 7L12 2Z" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
                <path d="M2 17L12 22L22 17" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
                <path d="M2 12L12 17L22 12" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
              </svg>
              Connect Talisman (Paseo)
            </button>
          </div>
          
          <div style={{ marginTop: '1rem', fontSize: '0.875rem', color: '#666' }}>
            <p><strong>Note:</strong> MetaMask on Sepolia/Paseo testnets and Talisman on Paseo are supported for smart contract interactions.</p>
          </div>
        </div>
      ) : (
        <div>
          <div style={{ marginBottom: '1rem' }}>
            <p style={{ margin: '0.5rem 0', color: '#666' }}>
              <strong>Connected:</strong> {formatAddress(walletState.address!)}
            </p>
            <p style={{ margin: '0.5rem 0', color: '#666' }}>
              <strong>Wallet:</strong> {getWalletName()}
            </p>
            <p style={{ margin: '0.5rem 0', color: '#666' }}>
              <strong>Network:</strong> {getNetworkName()}
            </p>
            {walletState.balance && (
              <p style={{ margin: '0.5rem 0', color: '#666' }}>
                <strong>Balance:</strong> {walletState.balance} {getBalanceUnit()}
              </p>
            )}
          </div>
          <button
            onClick={disconnectWallet}
            style={{
              background: '#dc3545',
              color: 'white',
              border: 'none',
              padding: '0.5rem 1rem',
              borderRadius: '4px',
              cursor: 'pointer',
              fontSize: '0.875rem',
            }}
          >
            Disconnect
          </button>
        </div>
      )}

      {/* Talisman Connect Modal */}
      <Modal
        title="Connect Wallet"
        isOpen={isModalOpen}
        handleClose={() => setIsModalOpen(false)}
        appId="milkyway2-portal"
      >
        <div style={{ padding: '1rem' }}>
          <h3 style={{ marginBottom: '1rem', color: '#333' }}>Select Wallet</h3>
          <div style={{ display: 'flex', flexDirection: 'column', gap: '0.75rem' }}>
            {getWallets()
              .filter((wallet) => wallet.installed && wallet.extensionName === 'talisman')
              .map((wallet) => (
                <button
                  key={wallet.extensionName}
                  onClick={async () => {
                    setIsModalOpen(false);
                    await handleWalletConnect(wallet);
                  }}
                  disabled={isConnecting}
                  style={{
                    background: '#f8f9fa',
                    color: '#333',
                    border: '1px solid #dee2e6',
                    padding: '1rem',
                    borderRadius: '8px',
                    cursor: isConnecting ? 'not-allowed' : 'pointer',
                    fontSize: '0.875rem',
                    fontWeight: '500',
                    opacity: isConnecting ? 0.6 : 1,
                    transition: 'all 0.2s ease',
                    display: 'flex',
                    alignItems: 'center',
                    gap: '0.75rem',
                    width: '100%',
                    textAlign: 'left',
                  }}
                  onMouseEnter={(e) => {
                    if (!isConnecting) {
                      e.currentTarget.style.background = '#e9ecef';
                      e.currentTarget.style.borderColor = '#adb5bd';
                    }
                  }}
                  onMouseLeave={(e) => {
                    if (!isConnecting) {
                      e.currentTarget.style.background = '#f8f9fa';
                      e.currentTarget.style.borderColor = '#dee2e6';
                    }
                  }}
                >
                  {wallet.logo && (
                    <img 
                      src={wallet.logo.src} 
                      alt={wallet.title}
                      style={{ width: '32px', height: '32px', borderRadius: '4px' }}
                    />
                  )}
                  <div>
                    <div style={{ fontWeight: '600', marginBottom: '0.25rem' }}>
                      {wallet.title}
                    </div>
                    <div style={{ fontSize: '0.75rem', color: '#6c757d' }}>
                      {wallet.installed ? 'Installed' : 'Not installed'}
                    </div>
                  </div>
                </button>
              ))}
          </div>
          
          {getWallets().filter((wallet) => wallet.extensionName === 'talisman').length === 0 && (
            <div style={{ 
              textAlign: 'center', 
              padding: '2rem', 
              color: '#6c757d',
              fontSize: '0.875rem'
            }}>
              <p>No supported wallets found.</p>
              <p>Please install Talisman wallet to connect.</p>
            </div>
          )}
        </div>
      </Modal>
    </div>
  );
}; 