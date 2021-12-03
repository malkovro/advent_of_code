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
    let gamma = find_gamma_rating(diagnostic_numbers.clone(), size);
    let epsilon = (2 as u32).pow(size as u32) - 1 - gamma;
    let oxy_rating = find_oxygen_generator_rating(diagnostic_numbers.clone(), size);
    let co2_rating = find_co2_scrubber_rating(diagnostic_numbers, size);
    println!(
        "gamma({:?}) * epsilon({:?}): power({:?})\noxygen_rating({:?}) * co2_rating({:?}): life_support_rating({:?})",
        gamma,
        epsilon,
        gamma * epsilon,
        oxy_rating,
        co2_rating,
        oxy_rating * co2_rating
    );
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

fn find_gamma_rating(diagnostic_numbers: Vec<u32>, size: usize) -> u32 {
    let base: usize = 2;
    u32::from_str_radix(
        &(0..size)
            .rev()
            .fold("".to_string(), |binary_rep: String, bit: usize| {
                let mask = base.pow(bit.clone().try_into().unwrap()) as u32;
                let high_levels_count = diagnostic_numbers
                    .iter()
                    .filter(|number| *number & mask > 0)
                    .count();
                let bit_value = if high_levels_count * 2 >= diagnostic_numbers.len() {
                    "1"
                } else {
                    "0"
                };
                binary_rep + bit_value
            }),
        2,
    )
    .unwrap()
}

fn find_oxygen_generator_rating(mut diagnostic_numbers: Vec<u32>, size: usize) -> u32 {
    let base: usize = 2;
    for bit in (0..size).rev() {
        let mask = base.pow(bit.clone().try_into().unwrap()) as u32;
        let candidates_count = diagnostic_numbers.len();
        let high_levels_count = diagnostic_numbers
            .iter()
            .filter(|number| *number & mask > 0)
            .count();
        if high_levels_count * 2 >= candidates_count {
            diagnostic_numbers = diagnostic_numbers
                .iter()
                .filter(|number| *number & mask > 0)
                .map(|x| *x)
                .collect();
        } else {
            diagnostic_numbers = diagnostic_numbers
                .iter()
                .filter(|number| *number & mask == 0)
                .map(|x| *x)
                .collect();
        }
        if diagnostic_numbers.len() == 1 {
            return diagnostic_numbers[0];
        }
    }
    panic!("No diagnostic_number seemed to match!");
}

fn find_co2_scrubber_rating(mut diagnostic_numbers: Vec<u32>, size: usize) -> u32 {
    let base: usize = 2;
    for bit in (0..size).rev() {
        let mask = base.pow(bit.clone().try_into().unwrap()) as u32;
        let candidates_count = diagnostic_numbers.len();
        let high_levels_count = diagnostic_numbers
            .iter()
            .filter(|number| *number & mask > 0)
            .count();
        if high_levels_count * 2 >= candidates_count {
            diagnostic_numbers = diagnostic_numbers
                .iter()
                .filter(|number| *number & mask == 0)
                .map(|x| *x)
                .collect();
        } else {
            diagnostic_numbers = diagnostic_numbers
                .iter()
                .filter(|number| *number & mask > 0)
                .map(|x| *x)
                .collect();
        }
        if diagnostic_numbers.len() == 1 {
            return diagnostic_numbers[0];
        }
    }
    panic!("No diagnostic_number seemed to match!");
}
