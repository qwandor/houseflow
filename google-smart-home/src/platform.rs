//! Enums automatically generated from JSON schema.

use schemafy::schemafy;
use serde::Deserialize;
use serde::Serialize;

schemafy!(
    root: Command
    "smart-home-schema/platform/commands.schema.json"
);
schemafy!(
    root: Error
    "smart-home-schema/platform/errors.schema.json"
);
schemafy!(
    root: Intent
    "smart-home-schema/platform/intents.schema.json"
);
schemafy!(
    root: Trait
    "smart-home-schema/platform/traits.schema.json"
);
schemafy!(
    root: Type
    "smart-home-schema/platform/types.schema.json"
);
