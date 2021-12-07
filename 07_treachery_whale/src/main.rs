use std::env;

fn main() {
    let args: Vec<String> = env::args().collect();
    let filename = args[1].clone();
    let crab_positions: Vec<i32> = read_into_vec(&filename)[0]
        .split(',')
        .map(|n| n.parse::<i32>().unwrap())
        .collect();

    let higher_crab = *crab_positions.iter().min().unwrap();
    let deeper_crab = *crab_positions.iter().max().unwrap();

    let sum_abs_diff = (higher_crab..(deeper_crab + 1))
        .map(|p| crab_positions.iter().map(|c| i32::abs(c - p)).sum::<i32>());

    let sum_arythm_diff = (higher_crab..(deeper_crab + 1)).map(|p| {
        crab_positions
            .iter()
            .map(|c| {
                let diff = i32::abs(c - p);
                diff * (diff + 1) / 2
            })
            .sum::<i32>()
    });

    println!(
        "The price for the crabs to set in position: (Part 1) {:?} / (Part 2) {:?}",
        sum_abs_diff.min(),
        sum_arythm_diff.min()
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
