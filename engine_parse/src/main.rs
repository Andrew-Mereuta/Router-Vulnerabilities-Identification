use anyhow::{Ok, Result};
use csv::ReaderBuilder;
use std::{fs, u128};
use thiserror::Error;

#[derive(Error, Debug)]
enum VendorParseError {
    #[error("Vendor entry (index: {0}) is missing an id.")]
    MissingID(usize),
    #[error("Vendor entry (index: {0}, id: {1}) is missing a name.")]
    MissingName(usize, u128),
    #[error("Vendor \"{2}\" (index: {0}, id: {1}) is missing a contact.")]
    MissingContact(usize, u128, String),
    #[error("Vendor \"{2}\" (index: {0}, id: {1}) is missing an email.")]
    MissingEmail(usize, u128, String),
}

#[derive(Debug, PartialEq, Eq)]
struct Vendor {
    id: u128,
    name: String,
    contact: String,
    email: String,
}

#[derive(Error, Debug)]
enum MacParseError {
    #[error("JSON is not a list, got: {0}")]
    JSONNotList(json::JsonValue),
    #[error("Entry field '{0}' has wrong type in {1}.")]
    IncorrectType(String, json::JsonValue),
}

#[derive(Debug, PartialEq, Eq)]
struct MacBlock {
    prefix: String,
    name: String,
    private: bool,
    block_type: String,
    last_update: String,
}

#[derive(Error, Debug)]
enum DeviceParseError {
    #[error("EngineID ({0:#x}) has unused format: {1}")]
    UnusedFormat(u128, u8),
}

#[derive(Debug, PartialEq, Eq)]
enum DeviceID<'a> {
    Mac {
        address: u128,
        blocks: &'a [MacBlock],
    },
    IPv4(u128),
    IPv6(u128),
    Text(u128),
    Octid(u128),
}

impl DeviceID<'_> {
    fn new(id: u128, overflow: u128, mac_blocks: &[MacBlock]) -> Result<Self> {
        let format_mask = 0x0000_0000_ff00_0000_0000_0000;
        let id_mask = 0x0000_0000_00ff_ffff_ffff_ffff;
        Ok(match ((id & format_mask) >> 56) as u8 {
            1 => DeviceID::IPv4(id_mask & id),
            2 => DeviceID::IPv6(id_mask & id),
            3 => {
                let id_str = 0;
                let macs = mac_blocks.iter().filter(|mb| unimplemented!("oops!"));
                unimplemented!("oh oh!")
            }
            n => Err(DeviceParseError::UnusedFormat(id, n))?,
        })
    }
}

#[derive(Error, Debug)]
enum EngineParseError {
    #[error("EngineID ({0}) is too long ({1} > 64).")]
    IdTooLong(String, usize),
}

#[derive(Debug, PartialEq, Eq)]
enum EngineID<'a> {
    CustomFormat {
        vendor: &'a Vendor,
        id: u64,
    },
    CorrectFormat {
        vendor: &'a Vendor,
        id: DeviceID<'a>,
    },
}

impl<'a> EngineID<'a> {
    fn new(
        id: u128,
        overflow: u128,
        vendors: &'a [Vendor],
        mac_blocks: &'a [MacBlock],
    ) -> Result<Self> {
        let conform_mask = 0x8000_0000_0000_0000_0000_0000_0000_0000;
        let vendor_mask = 0x7fff_ffff_0000_0000_0000_0000_0000_0000;
        let vendor = &vendors[((vendor_mask & id) >> 96) as usize];
        if (conform_mask & id) == 0 {
            let id_mask = 0x0000_0000_ffff_ffff_ffff_ffff_0000_0000;
            Ok(EngineID::CustomFormat {
                vendor,
                id: ((id_mask & id) >> 32) as u64,
            })
        } else {
            Ok(EngineID::CorrectFormat {
                vendor,
                id: DeviceID::new(id, overflow, mac_blocks)?,
            })
        }
    }

    fn from_string(
        id_str: &str,
        vendors: &'a [Vendor],
        mac_blocks: &'a [MacBlock],
    ) -> Result<Self> {
        if id_str.len() > 64 {
            Err(EngineParseError::IdTooLong(id_str.to_owned(), id_str.len()))?
        }
        let padded_str = format!("{:0<64}", id_str);
        let id = u128::from_str_radix(&padded_str[..32], 16)?;
        let overflow = u128::from_str_radix(&padded_str[32..], 16)?;
        EngineID::new(id, overflow, vendors, mac_blocks)
    }
}

fn load_vendors(file: &str) -> Result<Vec<Vendor>> {
    ReaderBuilder::new()
        .from_reader(file.as_bytes())
        .records()
        .enumerate()
        .map(|(i, result)| {
            let record = result?;
            let id = record
                .get(0)
                .ok_or(VendorParseError::MissingID(i))?
                .parse::<u128>()?;
            let name = record.get(1).ok_or(VendorParseError::MissingName(i, id))?;
            let contact =
                record
                    .get(2)
                    .ok_or(VendorParseError::MissingContact(i, id, name.to_owned()))?;
            let email =
                record
                    .get(3)
                    .ok_or(VendorParseError::MissingEmail(i, id, name.to_owned()))?;
            Ok(Vendor {
                id,
                name: name.to_owned(),
                contact: contact.to_owned(),
                email: email.to_owned(),
            })
        })
        .collect()
}

fn load_mac_blocks(file: &str) -> Result<Vec<MacBlock>> {
    match json::parse(file)? {
        json::JsonValue::Array(mut list) => list
            .iter_mut()
            .map(|e| {
                let tmp = e.clone();
                let prefix = e["macPrefix"]
                    .take_string()
                    .ok_or(MacParseError::IncorrectType(
                        "prefix".to_owned(),
                        tmp.clone(),
                    ))?;
                let name = e["vendorName"]
                    .take_string()
                    .ok_or(MacParseError::IncorrectType("name".to_owned(), tmp.clone()))?;
                let private = match e["private"].take() {
                    json::JsonValue::Boolean(b) => b,
                    _ => Err(MacParseError::IncorrectType(
                        "private".to_owned(),
                        tmp.clone(),
                    ))?,
                };
                let block_type =
                    e["blockType"]
                        .take_string()
                        .ok_or(MacParseError::IncorrectType(
                            "block-type".to_owned(),
                            tmp.clone(),
                        ))?;
                let last_update =
                    e["lastUpdate"]
                        .take_string()
                        .ok_or(MacParseError::IncorrectType(
                            "last-update".to_owned(),
                            tmp.clone(),
                        ))?;
                Ok(MacBlock {
                    prefix,
                    name,
                    private,
                    block_type,
                    last_update,
                })
            })
            .collect(),
        e => Err(MacParseError::JSONNotList(e).into()),
    }
}

fn main() -> Result<()> {
    let file =
        fs::read_to_string("data/enterprise-numbers").expect("Missing the enterprise numbers file");
    let _ = load_vendors(&file)?;
    let file = fs::read_to_string("data/mac-vendors-export.json")
        .expect("Missing the file for the MAC blocks.");
    let _ = load_mac_blocks(&file)?;
    println!("Files loaded, loading data...\n");
    Ok(())
}

#[cfg(test)]
mod tests {
    use core::panic;
    use std::fs;

    use crate::{load_mac_blocks, load_vendors, EngineID};

    #[test]
    fn test_load_vendors() {
        let file = fs::read_to_string("data/enterprise-numbers")
            .expect("Missing the enterprise numbers file");
        let vendors = load_vendors(&file);
        match vendors {
            // check no error loading
            Err(error) => panic!("Got an error: {}", error),
            // check order
            Ok(list) => assert!(list.iter().enumerate().all(|(i, e)| i == e.id as usize)),
        }
    }

    #[test]
    fn test_load_mac_blocks() {
        let file = fs::read_to_string("data/mac-vendors-export.json")
            .expect("Missing the file for the MAC blocks.");
        let mac_blocks = load_mac_blocks(&file);
        match mac_blocks {
            Err(error) => panic!("Got an error: {}", error),
            Ok(list) => (),
        }
    }

    #[test]
    fn test_parse_id() {
        let file = fs::read_to_string("data/enterprise-numbers")
            .expect("Missing the enterprise numbers file");
        let vendors = load_vendors(&file).expect("See `test_load_vendors`.");
        let file = fs::read_to_string("data/mac-vendors-export.json")
            .expect("Missing the file for the MAC blocks.");
        let mac_blocks = load_mac_blocks(&file).expect("See `test_load_mac_blocks`.");

        let custom = EngineID::new(
            0x0000_0009_abaf_aaaa_1234_5678_0000_0000,
            0,
            &vendors,
            &mac_blocks,
        )
        .expect("Correct ID should parse.");
        assert_eq!(
            custom,
            EngineID::CustomFormat {
                vendor: &vendors[9],
                id: 0xabaf_aaaa_1234_5678
            }
        );

        let custom = EngineID::from_string("00000009abafaaaa12345678", &vendors, &mac_blocks)
            .expect("Correct ID should parse (from str).");
        assert_eq!(
            custom,
            EngineID::CustomFormat {
                vendor: &vendors[9],
                id: 0xabaf_aaaa_1234_5678
            }
        );
    }
}
