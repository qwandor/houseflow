pub mod fake;
pub mod smtp;

use async_trait::async_trait;
use houseflow_types::code::VerificationCode;
use lettre::Message;

#[derive(Debug, thiserror::Error)]
pub enum Error {
    #[error("smtp: {0}")]
    Smtp(#[from] lettre::transport::smtp::Error),
}

#[async_trait]
pub trait Mailer: Send + Sync {
    async fn send(&self, message: Message) -> Result<(), Error>;

    async fn send_verification_code(
        &self,
        address: &str,
        code: &VerificationCode,
    ) -> Result<(), Error> {
        let message = Message::builder()
            .from(self.from_address().parse().unwrap())
            .to(address.parse().unwrap())
            .subject("Your Houseflow account: Access from new device")
            .body(format!(
                "Your verification code: {}. It will be valid for next 30 minutes. Hurry up!",
                code
            ))
            .unwrap();
        self.send(message).await?;
        tracing::info!("Sent verification code to {}", address);
        Ok(())
    }

    fn from_address(&self) -> &str;
}

impl From<Error> for houseflow_types::errors::ServerError {
    fn from(val: Error) -> Self {
        houseflow_types::errors::InternalError::Mailer(val.to_string()).into()
    }
}
