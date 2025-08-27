import React, { useState, useCallback } from 'react';
import { WalletConnect } from '../components/WalletConnect';
import { ValidatorReport } from '../components/ValidatorReport';

export const Reports: React.FC = () => {
  const [walletState, setWalletState] = useState({
    isConnected: false,
    address: null as string | null,
    walletType: null as 'metamask' | 'talisman' | null,
  });

  // This callback will be called by WalletConnect when the wallet state changes
  const handleWalletStateChange = useCallback((isConnected: boolean, address: string | null, walletType?: 'metamask' | 'talisman' | null) => {
    console.log('Wallet state changed:', { isConnected, address, walletType });
    setWalletState({ isConnected, address, walletType: walletType || null });
  }, []);

  return (
    <div style={{ maxWidth: '800px', margin: '0 auto', padding: '2rem' }}>
      <h1 style={{ marginBottom: '2rem', color: '#333', textAlign: 'center' }}>
        Submit Report
      </h1>
      
      <p style={{ 
        marginBottom: '2rem', 
        color: '#666', 
        textAlign: 'center',
        fontSize: '1.1rem'
      }}>
        Connect your wallet and submit reports about validator behavior.
      </p>

      <WalletConnect onWalletStateChange={handleWalletStateChange} />
      <ValidatorReport 
        isWalletConnected={walletState.isConnected}
        walletAddress={walletState.address}
        walletType={walletState.walletType}
      />

      <div style={{ 
        marginTop: '3rem', 
        padding: '2rem', 
        background: '#f8f9fa', 
        borderRadius: '8px',
        border: '1px solid #dee2e6'
      }}>
        <h3 style={{ marginBottom: '1rem', color: '#333' }}>About Submit Report</h3>
        <div style={{ fontSize: '0.875rem', color: '#666', lineHeight: '1.6' }}>
          <p style={{ marginBottom: '1rem' }}>
            This feature allows users to submit reports about validator behavior on the Westend testnet. 
            Reports can include issues such as:
          </p>
          <ul style={{ marginBottom: '1rem', paddingLeft: '1.5rem' }}>
            <li><strong>Performance Issues:</strong> Validators with poor performance or frequent downtime</li>
            <li><strong>Malicious Behavior:</strong> Validators engaging in harmful activities</li>
            <li><strong>Inactive Validators:</strong> Validators that are not participating in consensus</li>
            <li><strong>Other Issues:</strong> Any other problematic behavior</li>
          </ul>
          <p>
            Reports are submitted to the blockchain and can be used by the community to identify 
            and address problematic validators. The smart contract for this feature is not yet deployed.
          </p>
        </div>
      </div>
    </div>
  );
}; 