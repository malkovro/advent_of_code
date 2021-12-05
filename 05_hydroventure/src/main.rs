use std::env;
use std::vec::Vec;

mod hydroventure;

fn main() {
    let args: Vec<String> = env::args().collect();
    let filename = args[1].clone();
    hydroventure::detect(&filename);
}


