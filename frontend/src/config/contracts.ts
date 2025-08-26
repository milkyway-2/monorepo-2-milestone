// Contract configuration
export const CONTRACT_CONFIG = {
  // Contract addresses for different networks
  WESTEND_ASSET_HUB_CONTRACT_ADDRESS: '0x42245eAe30399974e89D9DE9602403F23e980993',
  SEPOLIA_CONTRACT_ADDRESS: '0xC8CF29d9D1595a3588AD36E6349A0E9a5b632720',
  PASEO_CONTRACT_ADDRESS: '0x0000000000000000000000000000000000000000', // TODO: Add actual Paseo contract address
  
  // Chain IDs
  WESTEND_ASSET_HUB_CHAIN_ID: '0x190f1b45', // 420420421
  SEPOLIA_CHAIN_ID: '0xaa36a7',
  PASEO_CHAIN_ID: '0x1b4', // TODO: Verify this is the correct Paseo EVM chain ID
  
  // Oracle address for smart contract verification
  ORACLE_ADDRESS: '0x6c6Fa8CEeF6AbB97dCd75a6e390386E4B49A5e09',
} as const;

// Helper function to get contract address by chain ID
export const getContractAddress = (chainId: string): string => {
  switch (chainId) {
    case CONTRACT_CONFIG.WESTEND_ASSET_HUB_CHAIN_ID:
      return CONTRACT_CONFIG.WESTEND_ASSET_HUB_CONTRACT_ADDRESS;
    case CONTRACT_CONFIG.SEPOLIA_CHAIN_ID:
      return CONTRACT_CONFIG.SEPOLIA_CONTRACT_ADDRESS;
    case CONTRACT_CONFIG.PASEO_CHAIN_ID:
      return CONTRACT_CONFIG.PASEO_CONTRACT_ADDRESS;
    default:
      throw new Error('Unsupported EVM network. Please switch to a supported network.');
  }
};

// Helper function to get network name by chain ID
export const getNetworkName = (chainId: string): string => {
  switch (chainId) {
    case CONTRACT_CONFIG.WESTEND_ASSET_HUB_CHAIN_ID:
      return 'EVM';
    case CONTRACT_CONFIG.SEPOLIA_CHAIN_ID:
      return 'Sepolia';
    case CONTRACT_CONFIG.PASEO_CHAIN_ID:
      return 'Paseo';
    default:
      return 'Unknown';
  }
};
