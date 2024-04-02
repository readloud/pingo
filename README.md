# Pingo

ping in Go.

<br />

## Usage

```sh
pingo example.com

# sends pingo 5 times
pingo -c 5 example.com

# 2 seconds per pingo
pingo -i 2
```

<br />

## Installation

- **go install**

    If you have go and want to get the Pingo for the latest version, run 'go install'.

    ```sh
    go install github.com/hideckies/pingo@latest
    ```

- **Binary**

    Download a prebuilt binary from [release page](https://github.com/hideckies/pingo/releases).

- **git clone**

    ```sh
    git clone https://github.com/hideckies/pingo.git
    cd pingo
    go get ; go build
    ```

<br />

## Capabilities

If you feel annoying to 'sudo' every time you run, it encourages to set the capabilities as follow.

```sh
setcap cap_net_raw+ep ./pingo
```

<br />

## Unprivileged Ping

If you want to the unprivileged (UDP) ping in Linux, set the following sysctl command.

```sh
sudo sysctl -w net.ipv4.ping_group_range="0 2147483647"
```

Then run the pingo with the "-u" flag.

```sh
pingo -u example.com
```
