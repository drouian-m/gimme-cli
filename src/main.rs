use std::env;

use clap::Parser;

mod gimme_interface;

#[derive(Parser)]
#[clap(about, version, author)]
struct Args {
  #[clap(short, long)]
  name: String,

  #[clap(short, long)]
  version: String,

  #[clap(short, long)]
  file: String,
}

fn main() -> Result<(), Box<dyn std::error::Error>> {
  let args = Args::parse();

  let cdn_url = env::var("GIMME_CDN_URL").expect("Please set GIMME_CDN_URL environment variable with your Gimme server url instance");
  let gimme_token = env::var("GIMME_TOKEN").expect("Please set GIMME_TOKEN environment variable with a valid access token");
  gimme_interface::create_package(
    cdn_url.as_str(),
    args.name.as_str(),
    args.version.as_str(),
    args.file.as_str(),
    gimme_token.as_str(),
  )?;

  println!("Package {}@{} successfully published.", args.name, args.version);
  Ok(())
}
