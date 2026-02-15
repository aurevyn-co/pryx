use serde::{Deserialize, Serialize};
use serde_json::Value;

/// Sidecar process state
#[derive(Debug, Clone, Serialize, Deserialize)]
pub enum SidecarState {
    Stopped,
    Starting,
    Running,
    Crashed { attempts: u32 },
    Restarting { backoff_ms: u64 },
    Stopping,
}

/// RPC Request from Sidecar
#[derive(Debug, Deserialize)]
pub struct RpcRequest {
    #[allow(dead_code)]
    pub jsonrpc: String,
    pub method: String,
    pub params: Value,
    pub id: u64,
}

/// RPC Response to Sidecar
#[derive(Debug, Serialize)]
pub struct RpcResponse {
    pub jsonrpc: String,
    pub result: Value,
    pub id: u64,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct SidecarStatus {
    pub state: SidecarState,
    pub pid: Option<u32>,
    pub port: Option<u16>,
    pub uptime_secs: Option<f64>,
    pub crash_count: u32,
    pub started_at: Option<String>,
}

#[derive(Debug, thiserror::Error)]
pub enum SidecarError {
    #[error("Failed to spawn sidecar binary '{binary}': {reason}")]
    SpawnFailed { binary: String, reason: String },

    #[error("Sidecar process not running")]
    NoChild,

    #[error("Sidecar process not running: {0}")]
    ProcessNotRunning(String),

    #[error("Port discovery failed: {0}")]
    PortDiscoveryFailed(std::io::Error),

    #[error("Serialization error: {0}")]
    Serialization(String),

    #[error("IO error: {0}")]
    Io(#[from] std::io::Error),

    #[error("Operation cancelled")]
    Cancelled,
}
