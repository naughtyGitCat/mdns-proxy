# mdns-proxy
proxy local mdns service to remote, by a simple dns server

## attention
some linux distros run `systemd-resolved` by default. 
`systemd-resolved` service occupy the `:53` port which is dns default port.
in macos or somewhere, if does not support custom dns port, you need stop `systemd-resolved` service, free port `:53`

```bash
systemctl stop systemd-resolved  # temp stop
systemctl disable systemd-resolved # disable start after reboot
```