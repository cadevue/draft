# IF3230 Consensus Protocol: Raft
This program is a simple simulation of the Raft Consensus Protocol. The program is implemented in golang.

## How to Run 🚀
### Running the Server 🖥️
- As first node (automatically become leader)
    ```bash
    go run server/server.go <ip server> <port server>
    ```
- As follower node

    ```bash 
    go run server/server.go <ip server> <port server> <ip leader server> <port leader server>
    ```

### Running the Client 👨🏻‍💼
#### Method 1: Web Interface 🧩
```bash 
go run web-client/web_client.go <ip server> <port server>
```

#### Method 2: The Terminal </>
```bash
go run client/client.go <ip server> <port server>
```

> Notes: Make sure the leader server address is an active leader node, if not the terminal will prompt the correct leader address

## Available Command
1. Ping: checking server connection.
    ```bash
    ping
    ```
    the server will response with:
    ```bash
    Pong!
    ```

2. Get: get a value with a key.
    ```bash
    get <key>
    ```
    the server will response with:
    ```bash
    <value>
    ```
    If the key doesn't exist, it will return empty string.

3. Set: set new value for specific key.
    ```bash
    set <key> <value>
    ```
    the server will response with:
    ```bash
    OK
    ```
    if key doesn't exist new pair will be created. This command only available for leader node, if not it will return error message and a correct leader address. This command will be logged.

4. Strln: get value's length for specific key.
    ```bash
    strln <key>
    ```
    the server will response with:
    ```bash
    <value length>
    ```
    if key doesn't exist it will return -1.

5. Del: delete value for specific key.
    ```bash
    del <key>
    ```
    the server will response with:
    ```bash
    <deleted value>
    ```
    if key doesn't exist it will return empty string. This command only available for leader node, if not it will return error message and a correct leader address. This command will be logged.

6. Append: append value for specific key.
    ```bash
    append <key> <value>
    ```
    the server will response with:
    ```bash
    OK
    ```
    if key doesn't exist new pair will be created. This command only available for leader node, if not it will return error message and a correct leader address. This command will be logged.

7. Switch: switch server port.
    ```bash
    switch <ip server> <port server>
    ```
    the server will response with:
    ```bash
    OK
    ```
8. Request Log: get server's log.
    ```bash
    request_log
    ```
    the server will response with:
    ```bash
    <server log>
    ```
    
## Contributors

<table>
<tr>
    <td align="center" colspan="3">Made by Team dRaft</td>
</tr>
<tr>
    <td align="center">NIM</td>
    <td align="center">Nama</td>
    <td align="center">Username</td>
</tr>
    <td align="center">13521050</td>
    <td align="center">Naufal Syifa Firdaus</td>
    <td align="center"><a href=https://github.com/nomsf>nomsf</a></td>
</tr>
    <td align="center">13521069</td>
    <td align="center">Louis Caesa Kesuma</td>
    <td align="center"><a href=https://github.com/Ainzw0rth>Ainzw0rth</a></td>
</tr>
</tr>
    <td align="center">13521077</td>
    <td align="center">Husnia Munzayana</td>
    <td align="center"><a href=https://github.com/munzayanahusn>munzayanahusn</a></td>
</tr>
</tr>
    <td align="center">13521085</td>
    <td align="center">Addin Munawwar Yusuf</td>
    <td align="center"><a href=https://github.com/cadevue>moonawar</a></td>
</tr>
</tr>
    <td align="center">13521088</td>
    <td align="center">Puti Nabilla Aidira</td>
    <td align="center"><a href=https://github.com/Putinabillaa>Putinabillaa</a></td>
</tr>
</table>
