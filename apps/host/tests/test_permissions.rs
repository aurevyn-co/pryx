use pryx_host::sidecar::{find_pryx_core_binary, SidecarConfig, SidecarProcess, SidecarStatus};

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_sidecar_config_defaults() {
        let config = SidecarConfig::default();
        assert_eq!(config.binary_path, "pryx-core");
        assert_eq!(config.max_restarts, 10);
        assert_eq!(config.initial_backoff_ms, 1000);
        assert_eq!(config.max_backoff_ms, 30000);
        assert_eq!(config.backoff_multiplier, 2.0);
        assert_eq!(config.port_discovery_timeout_secs, 10);
    }

    #[test]
    fn test_sidecar_status_stopped() {
        let status = SidecarStatus::Stopped;
        assert!(matches!(status, SidecarStatus::Stopped));
    }

    #[test]
    fn test_sidecar_status_running() {
        let status = SidecarStatus::Running;
        assert!(matches!(status, SidecarStatus::Running));
    }

    #[test]
    fn test_find_pryx_core_binary_fallback() {
        let path = find_pryx_core_binary();
        assert!(path.contains("pryx-core"));
    }
}
