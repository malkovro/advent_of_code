use std::collections::HashMap;
use std::collections::HashSet;
use termion::{color, style};

use regex::Regex;

#[derive(Debug, PartialEq, Eq, Hash, Clone, Copy)]
struct Point {
    x: u32,
    y: u32,
}

pub fn count_active_points(filename: &str) -> usize {
    let input = read_into_vec(&filename);
    let mut drawing_map = true;
    let mut map = HashMap::new();
    let mut points: Vec<Point> = vec![];
    let mut x_max: u32 = 0;
    let mut y_max: u32 = 0;
    let y_fold_regex = Regex::new(r"fold along y=(\d+)").unwrap();
    let x_fold_regex = Regex::new(r"fold along x=(\d+)").unwrap();
    let mut number_of_active_points = 0;
    let mut first_fold = true;
    for line in input.iter() {
        if line.is_empty() {
            drawing_map = false;
            continue;
        }
        if drawing_map {
            let coord: Vec<String> = line.split(',').map(String::from).collect();
            let point = Point {
                x: coord[0].parse::<u32>().unwrap(),
                y: coord[1].parse::<u32>().unwrap(),
            };

            if point.x > x_max {
                x_max = point.x;
            }
            if point.y > y_max {
                y_max = point.y;
            }
            map.insert(point, true);
            points.push(point);
        } else {
            if let Some(y_match) = y_fold_regex.captures(line) {
                let y_index_match = y_match.get(1).unwrap().as_str();
                let y_index = y_index_match.parse::<u32>().unwrap();
                points = points
                    .iter()
                    .map(|p| {
                        if p.y > y_index {
                            Point {
                                x: p.x,
                                y: y_index - (p.y - y_index),
                            }
                        } else {
                            p.clone()
                        }
                    })
                    .collect();
                y_max = y_index;
            }
            if let Some(x_match) = x_fold_regex.captures(line) {
                let x_index_match = x_match.get(1).unwrap().as_str();
                let x_index = x_index_match.parse::<u32>().unwrap();
                points = points
                    .iter()
                    .map(|p| {
                        if p.x > x_index {
                            Point {
                                x: x_index - (p.x - x_index),
                                y: p.y,
                            }
                        } else {
                            p.clone()
                        }
                    })
                    .collect();
                x_max = x_index;
            }
            if first_fold {
                first_fold = false;
                number_of_active_points = points
                    .iter()
                    .fold(HashSet::new(), |mut set, p| {
                        set.insert(p);
                        set
                    })
                    .len();
            }
        }
    }
    for y in 0..y_max {
        for x in 0..x_max {
            let p = points.iter().find(|p| p.x == x && p.y == y);
            if p.is_none() {
                print!(" ");
            } else {
                print!("{}{}{}", color::Bg(color::Green), 'X', style::Reset)
            }
        }
        println!("");
    }
    number_of_active_points
}

fn read_into_vec(file_name: &str) -> Vec<String> {
    std::fs::read_to_string(file_name)
        .unwrap()
        .split('\n')
        .into_iter()
        // .filter(|s| !s.is_empty())
        .map(String::from)
        .collect()
}
