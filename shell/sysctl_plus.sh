cat >>/etc/sysctl.conf <<EOL 
net.ipv4.tcp_tw_reuse = 1  
net.ipv4.tcp_tw_recycle = 1  
net.ipv4.tcp_fin_timeout = 30  
net.ipv4.tcp_keepalive_time = 1200  
net.ipv4.tcp_max_syn_backlog = 8192 
net.ipv4.tcp_max_tw_buckets = 5000  
net.ipv4.tcp_orphan_retries = 3
net.core.default_qdisc = fq
net.ipv4.tcp_congestion_control = bbr
EOL

/sbin/sysctl --system

/sbin/sysctl -p 
