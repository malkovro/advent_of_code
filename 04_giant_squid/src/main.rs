use std::env;
use std::vec::Vec;

mod bingo;

fn main() {
    let args: Vec<String> = env::args().collect();
    let filename = args[1].clone();
    bingo::play(&filename, false);
}


