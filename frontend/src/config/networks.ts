import type { NetworkId } from '../types';

// The default network to use when no network is specified
export const DefaultNetwork: NetworkId = 'paseo';

// All supported networks
export const NetworkList = {
  local: {
    name: 'local',
    endpoints: {
      rpc: {
        'Local Node': 'ws://127.0.0.1:9944',
      },
    },
    unit: 'DOT',
    units: 10,
    ss58: 0,
    defaultFeeReserve: BigInt(1000000000),
    meta: {
      hubChain: 'statemint',
      peopleChain: 'people-polkadot',
      stakingChain: 'polkadot',
      subscanBalanceChainId: 'polkadot',
      supportOperators: true,
    },
  },
  polkadot: {
    name: 'polkadot',
    endpoints: {
      rpc: {
        'Automata 1RPC': 'wss://1rpc.io/dot',
        Dwellir: 'wss://polkadot-rpc.dwellir.com',
        'IBP-GeoDNS1': 'wss://rpc.ibp.network/polkadot',
        'IBP-GeoDNS2': 'wss://rpc.dotters.network/polkadot',
        LuckyFriday: 'wss://rpc-polkadot.luckyfriday.io',
        OnFinality: 'wss://polkadot.api.onfinality.io/public-ws',
        Stakeworld: 'wss://dot-rpc.stakeworld.io',
        Parity: 'wss://rpc.polkadot.io',
      },
    },
    unit: 'DOT',
    units: 10,
    ss58: 0,
    defaultFeeReserve: BigInt(1000000000),
    meta: {
      hubChain: 'statemint',
      peopleChain: 'people-polkadot',
      stakingChain: 'polkadot',
      subscanBalanceChainId: 'polkadot',
      supportOperators: true,
    },
  },
  kusama: {
    name: 'kusama',
    endpoints: {
      rpc: {
        'Automata 1RPC': 'wss://1rpc.io/ksm',
        Dwellir: 'wss://kusama-rpc.dwellir.com',
        'IBP-GeoDNS1': 'wss://rpc.ibp.network/kusama',
        'IBP-GeoDNS2': 'wss://rpc.dotters.network/kusama',
        LuckyFriday: 'wss://rpc-kusama.luckyfriday.io',
        OnFinality: 'wss://kusama.api.onfinality.io/public-ws',
        Stakeworld: 'wss://ksm-rpc.stakeworld.io',
      },
    },
    unit: 'KSM',
    units: 12,
    ss58: 2,
    defaultFeeReserve: BigInt(50000000000),
    meta: {
      hubChain: 'statemine',
      peopleChain: 'people-kusama',
      stakingChain: 'kusama',
      subscanBalanceChainId: 'kusama',
      supportOperators: true,
    },
  },
  westend: {
    name: 'westend',
    endpoints: {
      rpc: {
        Dwellir: 'wss://westend-rpc.dwellir.com',
        'IBP-GeoDNS1': 'wss://rpc.ibp.network/westend',
        'IBP-GeoDNS2': 'wss://rpc.dotters.network/westend',
        LuckyFriday: 'wss://rpc-westend.luckyfriday.io',
        OnFinality: 'wss://westend.api.onfinality.io/public-ws',
        Stakeworld: 'wss://wnd-rpc.stakeworld.io',
      },
    },
    unit: 'WND',
    units: 12,
    ss58: 42,
    defaultFeeReserve: BigInt(100000000000),
    meta: {
      hubChain: 'westmint',
      stakingChain: 'westmint',
      peopleChain: 'people-westend',
      subscanBalanceChainId: 'assethub-westend',
      supportOperators: true,
    },
  },
  paseo: {
    name: 'paseo',
    endpoints: {
      rpc: {
        'Asset Hub Dotters': 'wss://asset-hub-paseo.dotters.network',
        'Asset Hub Dwellir': 'wss://asset-hub-paseo-rpc.n.dwellir.com',
        'Asset Hub IBP': 'wss://sys.ibp.network/asset-hub-paseo',
        'Asset Hub Stakeworld': 'wss://pas-rpc.stakeworld.io/assethub',
        'Asset Hub Turboflakes': 'wss://sys.turboflakes.io/asset-hub-paseo',
      },
    },
    unit: 'PAS',
    units: 10,
    ss58: 0,
    defaultFeeReserve: BigInt(1000000000),
    meta: {
      hubChain: 'paseo-asset-hub',
      stakingChain: 'paseo',
      peopleChain: 'people-paseo',
      subscanBalanceChainId: 'paseo',
      supportOperators: true,
    },
  },
};

// Get a random RPC endpoint for a network
export const getRandomRpcEndpoint = (network: NetworkId): string => {
  const networkConfig = NetworkList[network as keyof typeof NetworkList];
  if (!networkConfig) {
    throw new Error(`Network ${network} not supported`);
  }
  const endpoints = Object.values(networkConfig.endpoints.rpc);
  const randomIndex = Math.floor(Math.random() * endpoints.length);
  return endpoints[randomIndex] as string;
};

// Get all RPC endpoints for a network
export const getRpcEndpoints = (network: NetworkId): Record<string, string> => {
  const networkConfig = NetworkList[network as keyof typeof NetworkList];
  if (!networkConfig) {
    throw new Error(`Network ${network} not supported`);
  }
  return networkConfig.endpoints.rpc;
};
