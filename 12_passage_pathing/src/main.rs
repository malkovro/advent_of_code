use std::env;

mod passage_pathing;

fn main() {
    let args: Vec<String> = env::args().collect();
    let filename = args[1].clone();
    let number_of_path_to_the_end = passage_pathing::count_paths(&filename);
    println!("Number of paths: {}", number_of_path_to_the_end);
}
