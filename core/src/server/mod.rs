use crate::{Command, ServerCommandState};
use async_trait::async_trait;

use self::auth::RunAuthCommand;
use self::fulfillment::RunFulfillmentCommand;
use self::lighthouse::RunLighthouseCommand;

mod auth;
mod fulfillment;
mod lighthouse;

use clap::Clap;

#[derive(Clap)]
pub enum Service {
    Auth(RunAuthCommand),
    Lighthouse(RunLighthouseCommand),
    Fulfillment(RunFulfillmentCommand),
}

#[derive(Clap)]
pub enum ServerSubcommand {
    /// Run specific service
    Run {
        #[clap(subcommand)]
        service: Option<Service>,
    },
}

#[derive(Clap)]
pub struct ServerCommand {
    #[clap(subcommand)]
    pub subcommand: ServerSubcommand,
}

#[async_trait(?Send)]
impl Command<ServerCommandState> for ServerCommand {
    async fn run(&self, state: ServerCommandState) -> anyhow::Result<()> {
        match &self.subcommand {
            ServerSubcommand::Run { service } => match service {
                Some(service) => match service {
                    Service::Auth(cmd) => cmd.run(state).await,
                    Service::Lighthouse(cmd) => cmd.run(state).await,
                    Service::Fulfillment(cmd) => cmd.run(state).await,
                },
                None => {
                    let auth = RunAuthCommand {};
                    let lighthouse = RunLighthouseCommand {};
                    let fulfillment = RunFulfillmentCommand {};
                    tokio::select!(
                        result = auth.run(state.clone()) => result,
                        result = lighthouse.run(state.clone()) => result,
                        result = fulfillment.run(state.clone()) => result,
                    )
                }
            },
        }
    }
}
