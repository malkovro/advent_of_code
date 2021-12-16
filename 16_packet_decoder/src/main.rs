use std::env::args;
use std::iter::Iterator;

fn main() {
    let args: Vec<String> = args().collect();
    let filename = args[1].clone();
    let input = read_input(&filename);

    let res: Vec<char> = input
        .chars()
        .map(|h| {
            format!("{:0>4b}", h.to_digit(16).unwrap())
                .chars()
                .collect::<Vec<char>>()
        })
        .flatten()
        .collect();

    // println!("{:?}", res);
    let mut stream = res.into_iter();

    let (checksum, value, slurped) = read_packet(&mut stream);
    println!(
        "Value: {}, Checksum: {} after reading {} bits",
        value, checksum, slurped
    );
}

fn read_packet(stream: &mut std::vec::IntoIter<char>) -> (u32, u64, usize) {
    let mut checksum = 0;
    let mut slurped = 0;
    let version_bits = take_n(stream, 3).iter().collect::<String>();
    slurped += version_bits.len();
    // println!("Version: {}", version_bits);
    if version_bits.is_empty() {
        panic!("Nothing to see here!");
    }
    let version = u32::from_str_radix(&version_bits, 2).unwrap();
    checksum += version;
    let packet_type_id_bits = take_n(stream, 3).iter().collect::<String>();
    slurped += 3;
    let packet_type_id = u8::from_str_radix(&packet_type_id_bits, 2).unwrap();
    match packet_type_id {
        4 => {
            let (litteral, s) = read_litteral(stream);
            slurped += s;
            // println!(
            //     "Litteral Packet: {}",
            //     u64::from_str_radix(&litteral, 2).unwrap()
            // );
            return (
                checksum,
                u64::from_str_radix(&litteral, 2).unwrap(),
                slurped,
            );
        }
        _ => {
            let length_id_type = stream.next().unwrap();
            slurped += 1;
            if length_id_type == '1' {
                let number_subpacket_bits = take_n(stream, 11);
                slurped += 11;
                let number_subpacket =
                    u32::from_str_radix(&number_subpacket_bits.iter().collect::<String>(), 2)
                        .unwrap();

                // println!("Operator contains {} subpackets", number_subpacket);
                let subpckts: Vec<(u32, u64, usize)> =
                    (0..number_subpacket).map(|_| read_packet(stream)).collect();
                checksum += subpckts.iter().map(|(c, _, _)| c).sum::<u32>();
                slurped += subpckts.iter().map(|(_, _, s)| s).sum::<usize>();
                let values: Vec<u64> = subpckts.iter().map(|(_, v, _)| *v).collect();
                let value = operate(values, packet_type_id);
                return (checksum, value, slurped);
            } else {
                let subpackets_length_bits = take_n(stream, 15);
                slurped += 15;
                let length_subpackets =
                    u32::from_str_radix(&subpackets_length_bits.iter().collect::<String>(), 2)
                        .unwrap();
                // println!(
                //     "Operator contains subpackets of length {}",
                //     length_subpackets
                // );
                let slurped_total_target = slurped as u32 + length_subpackets;
                let mut values = vec![];
                while (slurped as u32) < slurped_total_target {
                    let (c, v, s) = read_packet(stream);
                    // println!("Read contained packet of {}bits", s);
                    checksum += c;
                    slurped += s;
                    values.push(v);
                }
                let value = operate(values, packet_type_id);
                return (checksum, value, slurped);
            }
        }
    }
}

fn operate(numbers: Vec<u64>, packet_type_id: u8) -> u64 {
    match packet_type_id {
        0 => numbers.iter().sum(),
        1 => numbers.iter().fold(1, |acc, x| acc * x),
        2 => *numbers.iter().min().unwrap(),
        3 => *numbers.iter().max().unwrap(),
        5 => {
            if numbers[0] > numbers[1] {
                1
            } else {
                0
            }
        }
        6 => {
            if numbers[0] < numbers[1] {
                1
            } else {
                0
            }
        }
        7 => {
            if numbers[0] == numbers[1] {
                1
            } else {
                0
            }
        }
        _ => {
            panic!("Unknown op")
        }
    }
}

fn read_litteral(stream: &mut std::vec::IntoIter<char>) -> (String, usize) {
    let mut litteral = String::new();
    let mut slurped = 0;
    loop {
        let chunk = take_n(stream, 5);
        slurped += 5;
        litteral += &chunk[1..].iter().collect::<String>();
        if chunk[0] == '0' {
            return (litteral, slurped);
        }
    }
}

fn take_n<T>(stream: &mut std::vec::IntoIter<T>, n: usize) -> Vec<T> {
    (0..n).fold(vec![], |mut v, _| {
        if let Some(el) = stream.next() {
            v.push(el);
        }
        v
    })
}

fn read_input(file_name: &str) -> String {
    std::fs::read_to_string(file_name)
        .unwrap()
        .split('\n')
        .into_iter()
        .filter(|s| !s.is_empty())
        .map(String::from)
        .next()
        .unwrap()
}
