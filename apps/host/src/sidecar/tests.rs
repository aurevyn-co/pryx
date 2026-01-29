#[cfg(test)]
mod tests {
    use crate::sidecar::permissions::{PermissionDialogConfig, PermissionManager};
    use crate::sidecar::SidecarConfig;
    use std::path::PathBuf;

    #[test]
    fn test_default_config() {
        let config = SidecarConfig::default();
        // sidecar_port/grpc_port are likely dynamically assigned or in env,
        // SidecarConfig struct has: binary, args, env, cwd, db_path, etc.
        // Check that binary ends with "pryx-core" (it may be an absolute path)
        assert!(
            config.binary.to_string_lossy().ends_with("pryx-core")
                || config.binary.as_os_str() == "pryx-core"
        );
        assert_eq!(config.db_path.to_string_lossy(), "pryx.db");
    }

    #[test]
    fn test_config_validation() {
        let mut config = SidecarConfig::default();
        config.binary = PathBuf::from("custom-bin");
        assert_eq!(config.binary.to_string_lossy(), "custom-bin");
    }

    #[test]
    fn test_permission_manager_initial_state() {
        let config = PermissionDialogConfig::default();
        let manager = PermissionManager::new(config);

        let pending = manager.list_pending();
        assert!(pending.is_empty());
    }
}
