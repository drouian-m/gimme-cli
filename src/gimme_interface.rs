use std::collections::HashMap;

use reqwest::blocking::Client;
use reqwest::blocking::multipart::Form;
use simple_error::bail;

pub fn create_package(cdn_url: &str, pkg_name: &str, pkg_version: &str, pkg_archive_path: &str, gimme_token: &str) -> Result<(), Box<dyn std::error::Error>> {
  let form = Form::new()
    .text("name", pkg_name.to_string())
    .text("version", pkg_version.to_string())
    .file(
      "file",
      pkg_archive_path,
    )
    .unwrap();

  let client = Client::new();
  let resp = client
    .post(format!("{}/packages", cdn_url))
    .multipart(form)
    .header("Authorization", gimme_token)
    .send()?;

  if resp.status().is_success() {
    Ok(())
  } else {
    let resp_body = resp.json::<HashMap<String, String>>()?;
    let error = resp_body.get("error").unwrap();
    bail!("Error while uploading package. Cause : {}", error);
  }
}

#[cfg(test)]
mod tests {
  use httpmock::prelude::*;
  use serde_json::json;

  use crate::gimme_interface;

  #[test]
  fn already_exists_error() {
    let server = MockServer::start();

    let cdn_mock = server.mock(|when, then| {
      when.path("/packages");
      then.status(409).json_body(json!({ "error": "the package already exists" }));
    });

    assert!(gimme_interface::create_package(
      server.url("").as_str(),
      "test",
      "0.0.1",
      "resources/test/test.zip",
      "xxxx").is_err());
    cdn_mock.assert();
  }

  #[test]
  fn publish_success() {
    let server = MockServer::start();

    let cdn_mock = server.mock(|when, then| {
      when.path("/packages");
      then.status(201);
    });

    assert!(gimme_interface::create_package(
      server.url("").as_str(),
      "test",
      "0.0.1",
      "resources/test/test.zip",
      "xxxx").is_ok());
    cdn_mock.assert();
  }
}