import React, { useState, useEffect, useCallback } from 'react';
import { ethers } from 'ethers';
import { WalletConnect } from '../components/WalletConnect';
import { getContractAddress, CONTRACT_CONFIG } from '../config/contracts';

// Polkadot Brand Colors
const POLKADOT_COLORS = {
  primary: '#E6007A', // Polkadot Pink
  secondary: '#07FFFF', // Cyan
  violet: '#7916F3',
  white: '#FFFFFF',
  black: '#000000',
  storm200: '#DCE2E9', // Light gray
  storm400: '#AEB7CB', // Medium gray
  storm700: '#6E7391', // Dark gray
  pinkLight: '#FF2670', // Lighter pink variation
};

interface Message {
  validator_address: string;
  nominator_address: string;
  msgText: string;
}

export const Messages: React.FC = () => {
  const [walletState, setWalletState] = useState({
    isConnected: false,
    address: null as string | null,
  });
  const [messages, setMessages] = useState<Message[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [totalCount, setTotalCount] = useState(0);
  const [currentPage, setCurrentPage] = useState(1);
  const [messagesPerPage] = useState(10);

  // Wallet connection callback
  const handleWalletStateChange = useCallback((isConnected: boolean, address: string | null) => {
    console.log('Wallet state changed:', { isConnected, address });
    setWalletState({ isConnected, address });
  }, []);

  // Smart contract ABI (extracted from ValidatorReport.tsx)
  const contractABI = [
    { "inputs": [], "stateMutability": "nonpayable", "type": "constructor" },
    {
      "anonymous": false,
      "inputs": [
        { "indexed": false, "internalType": "string", "name": "validator", "type": "string" },
        { "indexed": false, "internalType": "string", "name": "nominator", "type": "string" },
        { "indexed": false, "internalType": "string", "name": "msgText", "type": "string" }
      ],
      "name": "MessageStored",
      "type": "event"
    },
    {
      "inputs": [ { "internalType": "uint256", "name": "index", "type": "uint256" } ],
      "name": "getMessage",
      "outputs": [
        {
          "components": [
            { "internalType": "string", "name": "validator_address", "type": "string" },
            { "internalType": "string", "name": "nominator_address", "type": "string" },
            { "internalType": "string", "name": "msgText", "type": "string" }
          ],
          "internalType": "struct OracleVerifiedDelegation.Message",
          "name": "",
          "type": "tuple"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "getMessageCount",
      "outputs": [ { "internalType": "uint256", "name": "", "type": "uint256" } ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [ { "internalType": "uint256", "name": "", "type": "uint256" } ],
      "name": "messages",
      "outputs": [
        { "internalType": "string", "name": "validator_address", "type": "string" },
        { "internalType": "string", "name": "nominator_address", "type": "string" },
        { "internalType": "string", "name": "msgText", "type": "string" }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "oracleAddress",
      "outputs": [ { "internalType": "address", "name": "", "type": "address" } ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        { "internalType": "string", "name": "validator_address", "type": "string" },
        { "internalType": "string", "name": "nominator_address", "type": "string" },
        { "internalType": "string", "name": "msgText", "type": "string" },
        { "internalType": "bytes", "name": "signature", "type": "bytes" }
      ],
      "name": "submitMessage",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        { "internalType": "string", "name": "validator_address", "type": "string" },
        { "internalType": "string", "name": "nominator_address", "type": "string" },
        { "internalType": "string", "name": "msgText", "type": "string" }
      ],
      "name": "submitMessageUnverified",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    }
  ];

  const getContractAddressFromProvider = async (provider: ethers.providers.Web3Provider): Promise<string> => {
    const chainId = await provider.send('eth_chainId', []);
    return getContractAddress(chainId);
  };

  const getProvider = async (): Promise<ethers.providers.Web3Provider> => {
    const injectedEthereum: any = (window as any).ethereum;
    const talismanEth: any = (window as any).talismanEth;
    const candidates = [injectedEthereum, talismanEth].filter(Boolean);
    
    for (const c of candidates) {
      try {
        const req = c.request
          ? (method: string, params?: any[]) => c.request({ method, params })
          : (method: string, params?: any[]) => c.send(method, params ?? []);
        const accs = await req('eth_requestAccounts');
        if (Array.isArray(accs) && accs.length > 0) {
          return new ethers.providers.Web3Provider(c, 'any');
        }
      } catch (_) {}
    }
    throw new Error('No EVM accounts available. Please connect MetaMask/Talisman (EVM) and approve access.');
  };

  const fetchMessages = useCallback(async () => {
    if (!walletState.isConnected) {
      setError('Please connect your wallet first');
      return;
    }

    setIsLoading(true);
    setError(null);

    try {
      const provider = await getProvider();
      const contractAddress = await getContractAddressFromProvider(provider);
      const contract = new ethers.Contract(contractAddress, contractABI, provider);

      // Get total message count
      const count = await contract.getMessageCount();
      const totalMessages = count.toNumber();
      setTotalCount(totalMessages);

      if (totalMessages === 0) {
        setMessages([]);
        return;
      }

      // Calculate pagination
      const startIndex = (currentPage - 1) * messagesPerPage;
      const endIndex = Math.min(startIndex + messagesPerPage, totalMessages);

      // Fetch messages for current page
      const messagePromises = [];
      for (let i = startIndex; i < endIndex; i++) {
        messagePromises.push(contract.getMessage(i));
      }

      const messageResults = await Promise.all(messagePromises);
      const fetchedMessages: Message[] = messageResults.map((msg: any) => ({
        validator_address: msg.validator_address,
        nominator_address: msg.nominator_address,
        msgText: msg.msgText,
      }));

      setMessages(fetchedMessages);

    } catch (err: any) {
      console.error('Failed to fetch messages:', err);
      setError(err.message || 'Failed to fetch messages from smart contract');
    } finally {
      setIsLoading(false);
    }
  }, [walletState.isConnected, currentPage, messagesPerPage]);

  useEffect(() => {
    if (walletState.isConnected) {
      fetchMessages();
    }
  }, [fetchMessages]);

  const formatAddress = (address: string) => {
    if (address.length > 20) {
      return `${address.slice(0, 10)}...${address.slice(-8)}`;
    }
    return address;
  };

  const handleRefresh = () => {
    fetchMessages();
  };

  const handlePageChange = (newPage: number) => {
    setCurrentPage(newPage);
  };

  const totalPages = Math.ceil(totalCount / messagesPerPage);

  return (
    <div style={{ maxWidth: '1200px', margin: '0 auto', padding: '2rem' }}>
      <div style={{
        background: `linear-gradient(135deg, ${POLKADOT_COLORS.white} 0%, ${POLKADOT_COLORS.storm200} 100%)`,
        padding: '2rem',
        borderRadius: '12px',
        boxShadow: `0 4px 20px rgba(230, 0, 122, 0.1)`,
        border: `1px solid ${POLKADOT_COLORS.storm200}`,
        marginBottom: '2rem'
      }}>
        <h1 style={{ 
          marginBottom: '1rem', 
          color: POLKADOT_COLORS.black,
          background: `linear-gradient(135deg, ${POLKADOT_COLORS.primary} 0%, ${POLKADOT_COLORS.violet} 100%)`,
          WebkitBackgroundClip: 'text',
          WebkitTextFillColor: 'transparent',
          backgroundClip: 'text',
          textAlign: 'center'
        }}>
          Smart Contract Messages
        </h1>
        
        <p style={{ 
          marginBottom: '2rem', 
          color: POLKADOT_COLORS.storm700, 
          textAlign: 'center',
          fontSize: '1.1rem'
        }}>
          View all messages submitted to the validator reporting smart contract.
        </p>

        {!walletState.isConnected && (
          <div style={{
            padding: '1.5rem',
            background: `linear-gradient(135deg, ${POLKADOT_COLORS.pinkLight} 0%, ${POLKADOT_COLORS.primary} 100%)`,
            borderRadius: '8px',
            color: POLKADOT_COLORS.white,
            textAlign: 'center',
            marginBottom: '2rem'
          }}>
            <strong>Wallet Connection Required</strong><br />
            Please connect your wallet to view messages from the smart contract.
          </div>
        )}

        <WalletConnect onWalletStateChange={handleWalletStateChange} />

        {error && (
          <div style={{
            padding: '1rem',
            background: '#fee',
            borderRadius: '8px',
            color: '#c33',
            marginBottom: '2rem',
            border: '1px solid #fcc'
          }}>
            <strong>Error:</strong> {error}
          </div>
        )}

        {walletState.isConnected && (
          <>
            <div style={{
              display: 'flex',
              justifyContent: 'space-between',
              alignItems: 'center',
              marginBottom: '2rem',
              padding: '1rem',
              background: POLKADOT_COLORS.white,
              borderRadius: '8px',
              border: `1px solid ${POLKADOT_COLORS.storm200}`
            }}>
              <div>
                <strong>Total Messages:</strong> {totalCount}
              </div>
              <button
                onClick={handleRefresh}
                disabled={isLoading}
                style={{
                  padding: '0.5rem 1rem',
                  background: `linear-gradient(135deg, ${POLKADOT_COLORS.primary} 0%, ${POLKADOT_COLORS.pinkLight} 100%)`,
                  color: POLKADOT_COLORS.white,
                  border: 'none',
                  borderRadius: '6px',
                  cursor: isLoading ? 'not-allowed' : 'pointer',
                  fontWeight: '600',
                  opacity: isLoading ? 0.6 : 1
                }}
              >
                {isLoading ? 'Loading...' : 'Refresh'}
              </button>
            </div>

            {isLoading ? (
              <div style={{
                textAlign: 'center',
                padding: '3rem',
                color: POLKADOT_COLORS.storm700
              }}>
                Loading messages...
              </div>
            ) : messages.length === 0 ? (
              <div style={{
                textAlign: 'center',
                padding: '3rem',
                color: POLKADOT_COLORS.storm700,
                background: POLKADOT_COLORS.white,
                borderRadius: '8px',
                border: `1px solid ${POLKADOT_COLORS.storm200}`
              }}>
                No messages found in the smart contract.
              </div>
            ) : (
              <>
                <div style={{ marginBottom: '2rem' }}>
                  {messages.map((message, index) => (
                    <div
                      key={index}
                      style={{
                        background: POLKADOT_COLORS.white,
                        padding: '1.5rem',
                        borderRadius: '8px',
                        border: `1px solid ${POLKADOT_COLORS.storm200}`,
                        marginBottom: '1rem',
                        boxShadow: '0 2px 8px rgba(0, 0, 0, 0.05)'
                      }}
                    >
                      <div style={{
                        display: 'grid',
                        gridTemplateColumns: '1fr 1fr',
                        gap: '1rem',
                        marginBottom: '1rem'
                      }}>
                        <div>
                          <strong style={{ color: POLKADOT_COLORS.primary }}>Validator:</strong><br />
                          <span style={{ fontFamily: 'monospace', fontSize: '0.9rem' }}>
                            {formatAddress(message.validator_address)}
                          </span>
                        </div>
                        <div>
                          <strong style={{ color: POLKADOT_COLORS.violet }}>Nominator:</strong><br />
                          <span style={{ fontFamily: 'monospace', fontSize: '0.9rem' }}>
                            {formatAddress(message.nominator_address)}
                          </span>
                        </div>
                      </div>
                      <div>
                        <strong style={{ color: POLKADOT_COLORS.storm700 }}>Message:</strong><br />
                        <div style={{
                          background: POLKADOT_COLORS.storm200,
                          padding: '1rem',
                          borderRadius: '6px',
                          marginTop: '0.5rem',
                          whiteSpace: 'pre-wrap',
                          wordBreak: 'break-word'
                        }}>
                          {message.msgText}
                        </div>
                      </div>
                    </div>
                  ))}
                </div>

                {totalPages > 1 && (
                  <div style={{
                    display: 'flex',
                    justifyContent: 'center',
                    gap: '0.5rem',
                    marginTop: '2rem'
                  }}>
                    <button
                      onClick={() => handlePageChange(currentPage - 1)}
                      disabled={currentPage === 1}
                      style={{
                        padding: '0.5rem 1rem',
                        background: currentPage === 1 ? POLKADOT_COLORS.storm400 : POLKADOT_COLORS.white,
                        color: currentPage === 1 ? POLKADOT_COLORS.white : POLKADOT_COLORS.storm700,
                        border: `1px solid ${POLKADOT_COLORS.storm400}`,
                        borderRadius: '6px',
                        cursor: currentPage === 1 ? 'not-allowed' : 'pointer',
                        fontWeight: '600'
                      }}
                    >
                      Previous
                    </button>
                    
                    <span style={{
                      padding: '0.5rem 1rem',
                      background: POLKADOT_COLORS.white,
                      border: `1px solid ${POLKADOT_COLORS.storm400}`,
                      borderRadius: '6px',
                      fontWeight: '600'
                    }}>
                      Page {currentPage} of {totalPages}
                    </span>
                    
                    <button
                      onClick={() => handlePageChange(currentPage + 1)}
                      disabled={currentPage === totalPages}
                      style={{
                        padding: '0.5rem 1rem',
                        background: currentPage === totalPages ? POLKADOT_COLORS.storm400 : POLKADOT_COLORS.white,
                        color: currentPage === totalPages ? POLKADOT_COLORS.white : POLKADOT_COLORS.storm700,
                        border: `1px solid ${POLKADOT_COLORS.storm400}`,
                        borderRadius: '6px',
                        cursor: currentPage === totalPages ? 'not-allowed' : 'pointer',
                        fontWeight: '600'
                      }}
                    >
                      Next
                    </button>
                  </div>
                )}
              </>
            )}
          </>
        )}
      </div>
    </div>
  );
};
