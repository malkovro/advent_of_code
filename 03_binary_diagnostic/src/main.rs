use std::env;
use std::vec::Vec;

fn main() {
    let args: Vec<String> = env::args().collect();
    let filename = args[1].clone();
    let lines = read_into_vec(&filename);

    let diagnostic_numbers: Vec<u32> = lines
        .iter()
        .map(|x| u32::from_str_radix(x, 2).unwrap())
        .collect();
    let size = lines.first().unwrap().len();
    let base: usize = 2;
    let mut mask_counts = vec![0; size];
    for (bit, _e) in mask_counts.clone().iter().enumerate() {
        let mask = base.pow(bit.clone().try_into().unwrap()) as u32;
        for diagnostic_number in diagnostic_numbers.clone() {
            if diagnostic_number & mask > 0 {
                mask_counts[bit] += 1;
            }
        }
    }

    println!("{:?}", mask_counts);
    let threshold = &(lines.len() / 2);
    let gamma_binary = mask_counts
        .iter()
        .rev()
        .map(|count| {
            if count > threshold {
                "1".to_string()
            } else {
                "0".to_string()
            }
        })
        .collect::<Vec<String>>()
        .join("");
    let gamma = u32::from_str_radix(&gamma_binary, 2).unwrap();

    let epsilon: u32 = (2 as u32).pow(size.try_into().unwrap()) as u32 - gamma - 1;
    println!("{:?} * {:?}: {:?}", gamma, epsilon, gamma * epsilon);
}

fn read_into_vec(file_name: &str) -> Vec<String> {
    std::fs::read_to_string(file_name)
        .unwrap()
        .split('\n')
        .into_iter()
        .filter(|s| !s.is_empty())
        .map(String::from)
        .collect()
}
