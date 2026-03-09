CREATE TABLE IF NOT EXISTS tcp_tunnels (
  id INT AUTO_INCREMENT PRIMARY KEY,
  listen_port INT NOT NULL,
  target_host VARCHAR(255) NOT NULL,
  target_port INT NOT NULL,
  enabled BOOLEAN DEFAULT TRUE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS tcp_tunnel_logs (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  tunnel_id INT,
  level ENUM('info','warning','error'),
  message TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Test row for tcp_tunnels
INSERT INTO tcp_tunnels (listen_port, target_host, target_port, enabled) VALUES (52000, 'httpforever.com', 80, TRUE);