use serde::Deserialize;
use serde::Serialize;

#[derive(Debug, Clone, Deserialize, Serialize, PartialEq, thiserror::Error)]
pub enum Error {
    #[error("device not connected")]
    DeviceNotConnected,

    #[error("request timeout")]
    Timeout,
}
