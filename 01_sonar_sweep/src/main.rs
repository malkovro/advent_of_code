use std::io::{self, BufRead};
use std::fs::File;
use std::path::Path;
use  std::vec::Vec;

fn main() {
    let filename = "input.txt";

    //let window_size = 1; // Part One
    let window_size  = 3; // Part Two
    let mut windows = Vec::new();

    let mut increasing_depth_count = 0;
    let mut previous_depth_sum = std::f64::INFINITY;

    if let Ok(lines) = read_lines(filename) {
        for line in lines {
            if let Ok(depth_s) = line {
                if windows.len() < window_size {
                    windows.push(Vec::new());
                }
                let depth = depth_s.parse::<f64>().unwrap();
                for i in 0..windows.len() {
                    windows[i].push(depth);
                    
                    if windows[i].len() == window_size {
                        let window_sum : f64 = windows[i].iter().sum();
    
                        windows[i] = Vec::new();
                        if previous_depth_sum < window_sum {
                            increasing_depth_count += 1
                        }
                        previous_depth_sum = window_sum
                    }
                }
            }
        }
    }
    print!("The count of increasing depth is: {}", increasing_depth_count);
}


fn read_lines<P>(filename: P) -> io::Result<io::Lines<io::BufReader<File>>>
where P: AsRef<Path>, {
    let file = File::open(filename)?;
    Ok(io::BufReader::new(file).lines())
}
