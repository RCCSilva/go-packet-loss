# Packet Loss

Simple project to test impacts of packet loss and delays using TCP and UDP protocols.

## Commands

- Shows current rules

```shell
sudo tc qdisc show
```

- Adds 50% of packet loss

```shell
sudo tc qdisc replace dev lo root netem loss 50%
```

- Removes pack loss

```shell
sudo tc qdisc del dev lo root netem loss 50%
```

- Simulate packet delay

```shell
sudo tc qdisc replace dev lo root netem delay 2s 1.5s
```

- Remove packet delay

```shell
sudo tc qdisc del dev lo root netem delay 2s 1.5s
```