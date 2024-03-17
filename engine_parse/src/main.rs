use anyhow::Result;
use csv::ReaderBuilder;
use std::{collections::HashSet, env, fs, u128};
use thiserror::Error;

#[derive(Error, Debug)]
enum VendorParseError {
    #[error("Vendor entry (index: {0}) is missing an id.")]
    MissingID(usize),
    #[error("Vendor entry (index: {0}, id: {1}) is missing a name.")]
    MissingName(usize, u64),
    #[error("Vendor \"{2}\" (index: {0}, id: {1}) is missing a contact.")]
    MissingContact(usize, u64, String),
    #[error("Vendor \"{2}\" (index: {0}, id: {1}) is missing an email.")]
    MissingEmail(usize, u64, String),
}

#[derive(Debug, PartialEq, Eq, Clone)]
struct Vendor {
    id: u64,
    name: String,
    contact: String,
    email: String,
}

impl From<Vendor> for json::JsonValue {
    fn from(vendor: Vendor) -> json::JsonValue {
        json::object! {
            id: vendor.id,
            name: vendor.name,
            contact: vendor.contact,
            email: vendor.email,
        }
    }
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

impl From<&MacBlock> for json::JsonValue {
    fn from(value: &MacBlock) -> Self {
        json::object! {
            name: value.name.to_owned(),
            prefix: value.prefix.to_owned(),
            block_type: value.block_type.to_owned(),
            last_update: value.last_update.to_owned(),
            private: value.private,
        }
    }
}

#[derive(Error, Debug)]
enum DeviceParseError {
    #[error("EngineID ({0:#x}) has unused format: {1}")]
    UnusedFormat(u128, u8),
}

#[derive(Debug, PartialEq, Eq)]
enum DeviceID<'a> {
    Mac {
        address: String,
        blocks: Vec<&'a MacBlock>,
    },
    IPv4([u8; 4]),
    IPv6([u16; 8]),
    Text(String),
    Octid([u8; 27]),
}

impl<'a> DeviceID<'a> {
    fn new(id: u128, overflow: u128, id_l: usize, mac_blocks: &'a [MacBlock]) -> Result<Self> {
        let format_mask = 0x0000_0000_ff00_0000_0000_0000_0000_0000;
        let id_mask = 0x0000_0000_00ff_ffff_ffff_ffff_ffff_ffff;
        Ok(match ((id & format_mask) >> 88) as u8 {
            1 => DeviceID::IPv4((((id_mask & id) >> 56) as u32).to_be_bytes()),
            2 => DeviceID::IPv6(
                (((id_mask & id) << 40) + (overflow >> 88))
                    .to_be_bytes()
                    .chunks_exact(2)
                    .map(|s| u16::from_be_bytes(s.try_into().unwrap()))
                    .collect::<Vec<_>>()
                    .try_into()
                    .expect("16 bytes should be 8 u16"),
            ),
            3 => {
                let id_str = (id_mask & id).to_be_bytes()[6..12]
                    .iter()
                    .map(|b| format!("{:#04x}", b)[2..].to_owned())
                    .collect::<Vec<_>>()
                    .join(":");
                let macs = mac_blocks
                    .iter()
                    .filter(|mb| id_str.starts_with(&mb.prefix.to_lowercase()))
                    .collect();
                DeviceID::Mac {
                    address: id_str,
                    blocks: macs,
                }
            }
            4 | 128..=255 => DeviceID::Text(
                format!("{:#x}{:#034x}", id_mask & id, overflow).replace("0x", "")[..id_l - 10]
                    .to_string(),
            ),
            5 => DeviceID::Octid(
                [&(id_mask & id).to_be_bytes()[5..], &overflow.to_be_bytes()].concat()[..28]
                    .try_into()?,
            ),
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
        id_l: usize,
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
                id: DeviceID::new(id, overflow, id_l, mac_blocks)?,
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
        EngineID::new(id, overflow, id_str.len(), vendors, mac_blocks)
    }
}

impl From<EngineID<'_>> for json::JsonValue {
    fn from(id: EngineID<'_>) -> json::JsonValue {
        match id {
            EngineID::CustomFormat { vendor, id } => json::object! {
                vendor: vendor.clone(),
                id: id
            },
            EngineID::CorrectFormat { vendor, id } => match id {
                DeviceID::IPv4(ip) => json::object! {
                    vendor: vendor.clone(),
                    ipv4: format!("{}.{}.{}.{}", ip[0], ip[1],ip[2], ip[3]),
                },
                DeviceID::IPv6(ip) => json::object! {
                    vendor: vendor.clone(),
                    ipv6: format!("{}:{}:{}:{}:{}:{}:{}:{}", ip[0], ip[1], ip[2], ip[3], ip[4], ip[5], ip[6], ip[7]),
                },
                DeviceID::Mac { address, blocks } => json::object! {
                    vendor: vendor.clone(),
                    mac_address: address,
                    mac_blocks: blocks,
                },
                DeviceID::Text(id) => json::object! {
                    vendor: vendor.clone(),
                    id: id,
                },
                DeviceID::Octid(bytes) => json::object! {
                    vendor: vendor.clone(),
                    octets: bytes.as_slice(),
                },
            },
        }
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
                .parse::<u64>()?;
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
    let mut args = env::args().skip(1);
    let data_path = args.next().unwrap_or("data/snmp_results.csv".to_owned());
    let out_path = args.next().unwrap_or("data/engine_ids.json".to_owned());
    let vendor_path = args.next().unwrap_or("data/enterprise-numbers".to_owned());
    let mac_path = args
        .next()
        .unwrap_or("data/mac-vendors-export.json".to_owned());
    let file = fs::read_to_string(vendor_path).expect("Missing the enterprise numbers file");
    let vendors = load_vendors(&file)?;
    let file = fs::read_to_string(mac_path).expect("Missing the file for the MAC blocks.");
    let mac_blocks = load_mac_blocks(&file)?;
    println!("Files loaded, loading data...\n");
    let ids = ReaderBuilder::new()
        .from_reader(
            fs::read_to_string(data_path)
                .unwrap_or_else(|_| {
                    panic!("No SNMP data file found at `{}`", "data/snmp_results.csv")
                })
                .as_bytes(),
        )
        .records()
        .filter_map(|res| {
            let rec = match res {
                Err(error) => {
                    println!("Something went wrong while parsing.({})", error);
                    return None;
                }
                Ok(val) => val,
            };
            match rec.get(4) {
                Some("Error") | None => None,
                Some(s) => Some(s.to_owned()),
            }
        })
        .collect::<HashSet<_>>()
        .iter()
        .filter_map(
            |id_str| match EngineID::from_string(id_str, &vendors, &mac_blocks) {
                // parse failed, log it
                Err(error) => {
                    println!("Error parsing id({}): {}", id_str.to_owned(), error);
                    None
                }
                // parse succeeded
                Ok(eng_id) => Some((id_str.to_owned(), eng_id)),
            },
        )
        .fold(json::JsonValue::new_object(), |mut o, (s, id)| {
            o.insert(&s, id)
                .expect("Should be able to insert into an object.");
            o
        });
    // TODO: extend data instead of just replace?
    fs::write(out_path, json::stringify(ids))?;
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
            Ok(_list) => (),
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
            24,
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

        let custom = EngineID::from_string("8000000901c0a80165", &vendors, &mac_blocks)
            .expect("Correct ID should parse (from str).");
        assert_eq!(
            custom,
            EngineID::CorrectFormat {
                vendor: &vendors[9],
                id: crate::DeviceID::IPv4([0xc0, 0xa8, 0x01, 0x65]),
            }
        );

        let custom = EngineID::from_string(
            "800000090200112233445566778899aabbccddeeff",
            &vendors,
            &mac_blocks,
        )
        .expect("Correct ID should parse (from str).");
        assert_eq!(
            custom,
            EngineID::CorrectFormat {
                vendor: &vendors[9],
                id: crate::DeviceID::IPv6([
                    0x11, 0x2233, 0x4455, 0x6677, 0x8899, 0xaabb, 0xccdd, 0xeeff
                ]),
            }
        );

        let custom = EngineID::from_string("8000000904afaaaa12345678", &vendors, &mac_blocks)
            .expect("Correct ID should parse (from str).");
        assert_eq!(
            custom,
            EngineID::CorrectFormat {
                vendor: &vendors[9],
                id: crate::DeviceID::Text("afaaaa12345678".to_string())
            }
        );

        let custom = EngineID::from_string("80000009030020cfae95c1a0", &vendors, &mac_blocks)
            .expect("Correct ID should parse (from str).");
        assert_eq!(
            custom,
            EngineID::CorrectFormat {
                vendor: &vendors[9],
                id: crate::DeviceID::Mac {
                    address: "20:cf:ae:95:c1:a0".to_owned(),
                    blocks: vec![&mac_blocks[42160]],
                }
            }
        );
    }
}
