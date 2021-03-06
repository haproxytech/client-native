package spoe

const basicConfig = `# _version=1
[ip-reputation]
spoe-agent iprep-agent
    messages check-client-ip
    option var-prefix iprep
    timeout hello 2s
    timeout idle  2m
    timeout processing 10ms
    use-backend agents
    log global
    option async

spoe-message check-client-ip
    args ip=src
    event on-client-session if ! { src -f /etc/haproxy/whitelist.lst }

spoe-group mygroup
    messages mymessage
`
