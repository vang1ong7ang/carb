# Carb

Carb is a tool for p2p connections.

TODO

# Example of Config

```json
{
    "SK": "<THE PRIVATE KEY>",
    "LISTEN": "/ip4/0.0.0.0/tcp/5517",
    "RELAYMODE": [],
    "PEERS": {
        "<IDS OF PEERS>": "<ADDRESSES OF PEERS>",
    },
    "PROTOCOLS": [
        {
            "ID": "Ping",
            "CONFIG": {
                "PrintLog": true
            }
        },
        {
            "ID": "None",
            "CONFIG": {}
        }
    ],
    "CLIENTS": [
        {
            "id": "HeartBeater",
            "config": {
                "TimeInterval": 1000000000,
                "Target": "<ID OF TARGET>",
                "PrintLog": true
            }
        },
        {
            "id": "TCPListener",
            "config": {
                "Target": "<ID OF TARGET>",
                "ListenAddress": "127.0.0.1:1080",
                "PrintLog": true
            }
        }
    ]
}
```