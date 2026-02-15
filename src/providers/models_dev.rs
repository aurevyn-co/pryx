use reqwest::blocking::Client;
use serde::Deserialize;

const MODELS_DEV_URL: &str = "https://models.dev/api.json";

#[derive(Debug, Deserialize)]
pub struct ProviderCatalog {
    #[serde(flatten)]
    pub providers: std::collections::HashMap<String, ProviderInfo>,
}

impl ProviderCatalog {
    pub fn get(&self, id: &str) -> Option<&ProviderInfo> {
        self.providers.get(id)
    }
}

#[derive(Debug, Deserialize)]
pub struct ProviderInfo {
    pub id: String,
    pub name: String,
    #[serde(default)]
    pub api: Option<String>,
    #[serde(default)]
    pub env: Vec<String>,
    #[serde(default)]
    pub models: std::collections::HashMap<String, ModelInfo>,
}

#[derive(Debug, Deserialize)]
pub struct ModelInfo {
    pub id: String,
    pub name: String,
    #[serde(default)]
    pub family: Option<String>,
    #[serde(default)]
    pub tool_call: Option<bool>,
    #[serde(default)]
    pub reasoning: Option<bool>,
    #[serde(default)]
    pub attachment: Option<bool>,
    #[serde(default)]
    pub temperature: Option<bool>,
    #[serde(default)]
    pub cost: Option<ModelCost>,
    #[serde(default)]
    pub limit: Option<ModelLimit>,
    #[serde(default)]
    pub modalities: Option<ModelModalities>,
}

#[derive(Debug, Deserialize)]
pub struct ModelCost {
    #[serde(default)]
    pub input: Option<f64>,
    #[serde(default)]
    pub output: Option<f64>,
    #[serde(default)]
    pub cache_read: Option<f64>,
    #[serde(default)]
    pub cache_write: Option<f64>,
}

#[derive(Debug, Deserialize)]
pub struct ModelLimit {
    #[serde(default)]
    pub context: Option<u64>,
    #[serde(default)]
    pub output: Option<u64>,
}

#[derive(Debug, Deserialize)]
pub struct ModelModalities {
    #[serde(default)]
    pub input: Option<Vec<String>>,
    #[serde(default)]
    pub output: Option<Vec<String>>,
}

pub struct ModelsDevClient {
    client: Client,
    catalog: Option<ProviderCatalog>,
}

impl ModelsDevClient {
    pub fn new() -> Self {
        Self {
            client: Client::builder()
                .timeout(std::time::Duration::from_secs(30))
                .build()
                .unwrap_or_default(),
            catalog: None,
        }
    }

    pub fn fetch_catalog(&mut self) -> anyhow::Result<&ProviderCatalog> {
        if let Some(ref catalog) = self.catalog {
            return Ok(catalog);
        }

        tracing::info!("Fetching models.dev catalog from {}", MODELS_DEV_URL);

        let response = self.client.get(MODELS_DEV_URL).send()?;
        tracing::info!("Response status: {}", response.status());

        let text = response.text()?;
        tracing::info!("Response body length: {} bytes", text.len());

        let catalog: ProviderCatalog = serde_json::from_str(&text)?;
        self.catalog = Some(catalog);
        Ok(self.catalog.as_ref().unwrap())
    }

    pub fn get_provider(&self, id: &str) -> Option<&ProviderInfo> {
        self.catalog.as_ref()?.providers.get(id)
    }

    pub fn list_providers(&self) -> Vec<&ProviderInfo> {
        self.catalog
            .as_ref()
            .map(|c| c.providers.values().collect())
            .unwrap_or_default()
    }

    pub fn clear_cache(&mut self) {
        self.catalog = None;
    }
}

impl Default for ModelsDevClient {
    fn default() -> Self {
        Self::new()
    }
}
