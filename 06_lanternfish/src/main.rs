use std::env;
use std::vec::Vec;

mod fishlife;

fn main() {
    let args: Vec<String> = env::args().collect();
    let filename = args[1].clone();
    let epoch = args[2].clone().parse::<usize>().unwrap();
    let fishes_at_t0 = fishlife::fishes_at_t0(&filename);
    let population = fishlife::population(fishes_at_t0, epoch);
    println!(
        "Number of fishes after {} days: {:?}",
        epoch,
        population.iter().sum::<u128>()
    );
}


