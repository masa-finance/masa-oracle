import os

def create_directory_structure():
    base_path = 'masa-oracle'
    
    # Paths for directories and files
    paths = [
        '/domain/node',
        '/domain/node/node.go',
        '/domain/node/node_registered.go',
        '/domain/transaction',
        '/domain/transaction/transaction.go',
        '/domain/transaction/transaction_processed.go',
        '/domain/stake',
        '/domain/stake/stake.go',
        '/domain/stake/stake_increased.go',
        '/domain/stake/stake_decreased.go',
        '/domain/webhook',
        '/domain/webhook/webhook_data.go',
        '/application',
        '/application/node_service.go',
        '/application/transaction_service.go',
        '/application/stake_service.go',
        '/infrastructure/libp2p',
        '/infrastructure/libp2p/node_config.go',
        '/infrastructure/libp2p/peer_discovery.go',
        '/infrastructure/libp2p/transport.go',
        '/infrastructure/dht',
        '/infrastructure/dht/dht_config.go',
        '/infrastructure/dht/node_registration.go',
        '/infrastructure/dht/node_discovery.go',
        '/infrastructure/ethereum/contracts',
        '/infrastructure/ethereum/contracts/MasaToken.sol',
        '/infrastructure/ethereum/contracts/StakingContract.sol',
        '/infrastructure/ethereum/staking.go',
        '/infrastructure/ethereum/rewards.go',
        '/infrastructure/ethereum/truffle_config.go',
        '/infrastructure/db',
        '/infrastructure/db/badger_config.go',
        '/infrastructure/db/data_schema.go',
        '/infrastructure/db/crud_operations.go',
        '/infrastructure/webhook',
        '/infrastructure/webhook/api_server.go',
        '/infrastructure/webhook/data_propagation.go',
        '/infrastructure/security',
        '/infrastructure/security/authentication.go',
        '/infrastructure/security/encryption.go',
        '/utils',
        '/tests',
        'LICENSE',
        'README.md'
    ]
    
    # Create directories and files
    for path in paths:
        full_path = os.path.join(base_path, path)
        if path.endswith('/'):
            os.makedirs(full_path, exist_ok=True)
        else:
            open(full_path, 'a').close()

if __name__ == '__main__':
    create_directory_structure()
