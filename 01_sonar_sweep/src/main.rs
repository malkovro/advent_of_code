use std::io::{self, BufRead};
use std::fs::File;
use std::path::Path;

fn main() {
    let filename = "input.txt";

    let mut increasing_depth_count = 0;
    let mut previous_depth = std::f64::INFINITY;

    // File hosts must exist in current path before this produces output
    if let Ok(lines) = read_lines(filename) {
        // Consumes the iterator, returns an (Optional) String
        for line in lines {
            if let Ok(depth_s) = line {
                let depth = depth_s.parse::<f64>().unwrap();
                if previous_depth < depth {
                    increasing_depth_count += 1
                }
                previous_depth = depth
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
