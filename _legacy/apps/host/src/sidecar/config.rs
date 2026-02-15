use serde::{Deserialize, Serialize};
use std::collections::HashMap;
use std::path::PathBuf;
use std::time::Duration;

/// Sidecar configuration
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct SidecarConfig {
    pub binary: PathBuf,
    pub args: Vec<String>,
    pub env: HashMap<String, String>,
    pub cwd: PathBuf,
    pub db_path: PathBuf,
    pub start_timeout: Duration,
    pub max_restarts: u32,
    pub initial_backoff_ms: u64,
    pub backoff_multiplier: f64,
}

impl SidecarConfig {
    pub fn new(binary: PathBuf, cwd: PathBuf, db_path: PathBuf) -> Self {
        Self {
            binary,
            args: vec![],
            env: HashMap::new(),
            cwd,
            db_path,
            start_timeout: Duration::from_secs(3),
            max_restarts: 10,
            initial_backoff_ms: 1000,
            backoff_multiplier: 2.0,
        }
    }
}

impl Default for SidecarConfig {
    fn default() -> Self {
        Self::new(
            find_pryx_core_binary().unwrap_or_else(|| PathBuf::from("pryx-core")),
            std::env::current_dir().unwrap_or_default(),
            PathBuf::from("pryx.db"),
        )
    }
}

pub fn find_pryx_core_binary() -> Option<PathBuf> {
    if let Ok(p) = std::env::var("PRYX_CORE_PATH") {
        let p = PathBuf::from(p);
        if p.exists() {
            return Some(p);
        }
    }

    if let Ok(exe) = std::env::current_exe() {
        if let Some(p) = search_ancestors(&exe) {
            return Some(p);
        }
    }
    if let Ok(cwd) = std::env::current_dir() {
        if let Some(p) = search_ancestors(&cwd) {
            return Some(p);
        }
    }

    Some(PathBuf::from("pryx-core"))
}

fn search_ancestors(start: &std::path::Path) -> Option<PathBuf> {
    for a in start.ancestors().take(8) {
        let c = a.join("apps").join("runtime").join("pryx-core");
        if c.exists() {
            return Some(c);
        }
        // Check dist/bin
        let d = a.join("dist").join("pryx-core");
        if d.exists() {
            return Some(d);
        }
    }
    None
}
