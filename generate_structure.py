import os

def diagnostic_create_directory_structure():
    """
    Create a directory structure for the 'masa-oracle' project.
    
    Directories and files are created based on a predefined list.
    """
    base_path = '.'
    
    # Print the current working directory for diagnostic purposes
    print(f"Current working directory: {os.getcwd()}")
    print(f"Attempting to create directories and files under: {os.path.join(os.getcwd(), base_path)}")
    
    paths = [
        '/domain/',
        '/domain/node/',
        '/domain/node/node.go.tmp',
        '/domain/node/node_registered.go.tmp',.
        '/domain/transaction/',
        '/domain/transaction/transaction.go.tmp',
        '/domain/transaction/transaction_processed.go.tmp',
        '/domain/stake/',
        '/domain/stake/stake.go.tmp',
        '/domain/stake/stake_increased.go.tmp',
        '/domain/stake/stake_decreased.go.tmp',
        '/domain/webhook/',
        '/domain/webhook/webhook_data.go.tmp',
        '/application/',
        '/application/node_service.go.tmp',
        '/application/transaction_service.go.tmp',
        '/application/stake_service.go.tmp',
        '/infrastructure/',
        '/infrastructure/libp2p/',
        '/infrastructure/libp2p/node_config.go.tmp',
        '/infrastructure/libp2p/peer_discovery.go.tmp',
        '/infrastructure/libp2p/transport.go.tmp',
        '/infrastructure/dht/',
        '/infrastructure/dht/dht_config.go.tmp',
        '/infrastructure/dht/node_registration.go.tmp',
        '/infrastructure/dht/node_discovery.go.tmp',
        '/infrastructure/ethereum/',
        '/infrastructure/ethereum/contracts/',
        '/infrastructure/ethereum/contracts/MasaToken.sol',
        '/infrastructure/ethereum/contracts/StakingContract.sol',
        '/infrastructure/ethereum/staking.go.tmp',
        '/infrastructure/ethereum/rewards.go.tmp',
        '/infrastructure/ethereum/truffle_config.go.tmp',
        '/infrastructure/db/',
        '/infrastructure/db/badger_config.go.tmp',
        '/infrastructure/db/data_schema.go.tmp',
        '/infrastructure/db/crud_operations.go.tmp',
        '/infrastructure/webhook/',
        '/infrastructure/webhook/api_server.go.tmp',
        '/infrastructure/webhook/data_propagation.go.tmp',
        '/infrastructure/security/',
        '/infrastructure/security/authentication.go.tmp',
        '/infrastructure/security/encryption.go.tmp',
        '/utils/',
        '/tests/'
    ]
    
    # Create directories and files
    for path in paths:
        full_path = os.path.join(base_path, path.lstrip('/'))  # Removing leading slash for compatibility
        if path.endswith('/'):
            os.makedirs(full_path, exist_ok=True)
            print(f"Created directory: {full_path}")
        else:
            with open(full_path, 'a', encoding='utf-8') as f:
                pass
            print(f"Created file: {full_path}")

# To run the diagnostic function:
diagnostic_create_directory_structure()
