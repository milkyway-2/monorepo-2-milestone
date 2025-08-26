import React, { useState, useEffect } from 'react';
import BigNumber from 'bignumber.js';
import { CONTRACT_CONFIG } from '../config/contracts';
import { getWallets } from '@talismn/connect-wallets';

interface WalletState {
  isConnected: boolean;
  address: string | null;
  chainId: string | null;
  balance: string | null;
  walletType: 'metamask' | 'talisman' | null;
}

interface WalletConnectProps {
  onWalletStateChange?: (isConnected: boolean, address: string | null) => void;
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

  useEffect(() => {
    console.log('🔄 useEffect triggered - checking wallet connection on component mount');
    checkWalletConnection();
  }, []);

  // Call the callback whenever wallet state changes
  useEffect(() => {
    console.log('🔄 Wallet state changed:', { 
      isConnected: walletState.isConnected, 
      address: walletState.address,
      walletType: walletState.walletType 
    });
    
    if (onWalletStateChange) {
      console.log('📞 Calling onWalletStateChange callback...');
      onWalletStateChange(walletState.isConnected, walletState.address);
    }
  }, [walletState.isConnected, walletState.address, onWalletStateChange]);

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

  const checkWalletConnection = async () => {
    console.log('🔍 Starting wallet connection check...');
    const startTime = Date.now();
    
    // Check MetaMask (EVM) - only for Sepolia
    if (typeof window.ethereum !== 'undefined') {
      try {
        console.log('📡 Checking MetaMask connection...');
        const accounts = await window.ethereum.request({ method: 'eth_accounts' });
        if (accounts.length > 0) {
          const chainId = await window.ethereum.request({ method: 'eth_chainId' });
          console.log('✅ MetaMask connected with chainId:', chainId);
          // Only connect if on Sepolia
          if (chainId === SEPOLIA_CHAIN_ID) {
            await updateWalletState(accounts[0], chainId, 'metamask');
          }
        }
      } catch (err) {
        console.error('❌ Error checking MetaMask connection:', err);
      }
    }
    
    // Check Talisman using official SDK
    try {
      console.log('🔍 Checking for Talisman wallet using official SDK...');
      const installedWallets = getWallets().filter((wallet) => wallet.installed);
      console.log('📋 Installed wallets:', installedWallets.map(w => w.extensionName));
      
      const talismanWallet = installedWallets.find(
        (wallet) => wallet.extensionName === 'talisman',
      );
      
      if (talismanWallet) {
        console.log('✅ Talisman wallet found via SDK');
        // Check if already enabled
        try {
          const accounts = await talismanWallet.getAccounts();
          if (accounts.length > 0) {
            console.log('✅ Talisman already connected with accounts:', accounts.length);
            await updateWalletState(accounts[0].address, CONTRACT_CONFIG.PASEO_CHAIN_ID, 'talisman');
          }
        } catch (err) {
          console.log('📋 Talisman not yet enabled, will enable on connect');
        }
      } else {
        console.log('❌ Talisman wallet not found via SDK');
      }
    } catch (err) {
      console.error('❌ Error checking Talisman via SDK:', err);
    }
    
    const endTime = Date.now();
    console.log(`⏱️ Wallet connection check completed in ${endTime - startTime}ms`);
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
    console.log('🚀 Starting Talisman connection (Paseo)...');
    const startTime = Date.now();
    setIsConnecting(true);
    setError(null);

    try {
      console.log('🔍 Using Talisman Connect SDK...');
      
      // Get installed wallets
      const installedWallets = getWallets().filter((wallet) => wallet.installed);
      console.log('📋 Installed wallets:', installedWallets.map(w => w.extensionName));
      
      // Find Talisman wallet
      const talismanWallet = installedWallets.find(
        (wallet) => wallet.extensionName === 'talisman',
      );
      
      if (!talismanWallet) {
        setError('Talisman is not installed. Please install Talisman to use this feature.');
        return;
      }

      console.log('✅ Talisman wallet found, enabling...');
      
      // Enable the wallet
      await talismanWallet.enable('Milkyway2 Portal');
      console.log('✅ Talisman wallet enabled');
      
      // Get accounts
      const accounts = await talismanWallet.getAccounts();
      
      if (accounts.length === 0) {
        setError('No accounts found in Talisman. Please add an account to your wallet.');
        return;
      }

      console.log('✅ Talisman connected with accounts:', accounts.length);
      await updateWalletState(accounts[0].address, CONTRACT_CONFIG.PASEO_CHAIN_ID, 'talisman');
      
      const endTime = Date.now();
      console.log(`⏱️ Talisman connection completed in ${endTime - startTime}ms`);
      
    } catch (err: any) {
      console.error('❌ Error connecting Talisman:', err);
      const errorTime = Date.now();
      console.log(`⏱️ Talisman connection failed after ${errorTime - startTime}ms`);
      setError(err.message || 'Failed to connect Talisman (Paseo)');
    } finally {
      setIsConnecting(false);
    }
  };





  const disconnectWallet = () => {
    setWalletState({
      isConnected: false,
      address: null,
      chainId: null,
      balance: null,
      walletType: null,
    });
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
          <div style={{ 
            display: 'grid', 
            gridTemplateColumns: 'repeat(2, 1fr)', 
            gap: '1rem',
            maxWidth: '600px'
          }}>
            <button
              onClick={connectSepolia}
              disabled={isConnecting}
              style={{
                background: '#29b6f6',
                color: 'white',
                border: 'none',
                padding: '0.75rem 1rem',
                borderRadius: '8px',
                cursor: isConnecting ? 'not-allowed' : 'pointer',
                fontSize: '0.875rem',
                fontWeight: '600',
                opacity: isConnecting ? 0.6 : 1,
                width: '100%',
                minHeight: '48px',
                transition: 'all 0.2s ease',
              }}
            >
              Connect MetaMask (Sepolia)
            </button>
            <button
              onClick={connectPaseo}
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
                minHeight: '48px',
                transition: 'all 0.2s ease',
              }}
            >
              Connect Talisman (Paseo)
            </button>
          </div>
          <div style={{ marginTop: '1rem', fontSize: '0.875rem', color: '#666' }}>
            <p><strong>Note:</strong> MetaMask on Sepolia testnet and Talisman on Paseo are supported for smart contract interactions.</p>
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
    </div>
  );
}; 