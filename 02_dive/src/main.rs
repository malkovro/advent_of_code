use std::env;
use  std::vec::Vec;

#[derive(Debug)]
struct Position {
    aim: i64,
    horizontal: i64,
    naive_depth: i64,
    depth: i64
}

fn main() {
    let args: Vec<String> = env::args().collect();
    let filename = args[1].clone();
    let lines = read_into_vec(&filename);

    let position = Position {
        aim: 0,
        horizontal: 0,
        naive_depth: 0,
        depth: 0
    };

    let position = lines.iter().fold(position, |position, instruction| {
        let split = instruction.split(" ").collect::<Vec<&str>>();
        let value = split[1].parse::<i64>().unwrap();
        match split[0] {
            "forward" => Position { horizontal: position.horizontal + value, aim: position.aim, depth: position.depth + position.aim * value, naive_depth: position.naive_depth },
            "up" => Position { horizontal: position.horizontal, aim: position.aim - value, depth: position.depth, naive_depth: position.naive_depth - value },
            "down" => Position { horizontal: position.horizontal, aim: position.aim + value, depth: position.depth, naive_depth: position.naive_depth + value },
            _ => panic!("Unexpected Instruction"),
        }
    });
    println!("{:?}", position);
    println!("Solution Part 1: {}", position.horizontal * position.naive_depth);
    println!("Solution Part 2: {}", position.horizontal * position.depth);
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

