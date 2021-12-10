use rand::seq::SliceRandom;
use rand::thread_rng;
use std::collections::HashMap;
use std::env;
use termion::{color, cursor, style};

#[derive(PartialEq, Eq, Clone, Copy, Debug, Hash)]
struct Point {
    x: i32,
    y: i32,
}

fn build_smoke_map(input: Vec<String>) -> HashMap<Point, u32> {
    let mut smoke_map: HashMap<Point, u32> = HashMap::new();

    for (y, line) in input.iter().enumerate() {
        for (x, v) in line.chars().map(|v| v.to_digit(10).unwrap()).enumerate() {
            smoke_map.insert(
                Point {
                    x: x as i32,
                    y: y as i32,
                },
                v,
            );
        }
    }
    smoke_map
}

fn detect_low_points(map: &HashMap<Point, u32>) -> Vec<Point> {
    map.keys()
        .filter(|point| is_low_zone(map, point, &vec![]))
        .map(|x| *x)
        .collect()
}

fn main() {
    let args: Vec<String> = env::args().collect();
    let filename = args[1].clone();

    let lines = read_into_vec(&filename);
    let line_count = lines.len();
    let smoke_map = build_smoke_map(lines);
    let line_size = smoke_map.len() / line_count;
    let low_points = detect_low_points(&smoke_map);
    let colors: Vec<color::AnsiValue> = get_colors();

    println!(
        "Sum of Risk Levels: {}",
        low_points
            .iter()
            .map(|p| smoke_map.get(p).unwrap() + 1)
            .sum::<u32>()
    );

    let mut basin: Vec<Vec<Point>> = vec![];
    for p in low_points.into_iter() {
        basin.push(find_zone(
            vec![p],
            &smoke_map,
            basin.clone(),
            line_size,
            line_count,
            &colors
        ));
    }

    let mut basin_size: Vec<usize> = basin.iter().map(|z| z.len()).collect();
    basin_size.sort();

    let basinity_score = basin_size.iter().rev().take(3).fold(1, |acc, v| v * acc);

    print_map(line_size, line_count, &smoke_map, &basin, &colors);
    println!("Basinity Score: {}", basinity_score);
}

fn find_zone(
    points: Vec<Point>,
    map: &HashMap<Point, u32>,
    mut others: Vec<Vec<Point>>,
    x_size: usize,
    y_size: usize,
    colors: &Vec<color::AnsiValue>
) -> Vec<Point> {
    let mut points_in_zone = vec![];
    let mut added_points: Vec<Point> = points;
    loop {
        added_points = added_points.iter().fold(vec![], |mut zone, point| {
            let current_position = Point {
                x: point.x,
                y: point.y,
            };
            if !points_in_zone.contains(&current_position) {
                points_in_zone.push(current_position);
                for neighbour in find_neighbours(point).into_iter() {
                    if is_low_zone(map, &neighbour, &points_in_zone) {
                        zone.push(neighbour);
                    }
                }
            }
            zone
        });
        if added_points.len() == 0 {
            return points_in_zone;
        }
        others.push(points_in_zone.clone());
        // print_map(x_size, y_size, map, &others, colors);
        // print!("{}", cursor::Up(y_size as u16));
    }
}

fn find_neighbours(point: &Point) -> Vec<Point> {
    vec![
        Point {
            x: point.x - 1,
            y: point.y,
        },
        Point {
            x: point.x,
            y: point.y - 1,
        },
        Point {
            x: point.x + 1,
            y: point.y,
        },
        Point {
            x: point.x,
            y: point.y + 1,
        },
    ]
}

fn is_low_zone(map: &HashMap<Point, u32>, point: &Point, exclude: &Vec<Point>) -> bool {
    let option_value = map.get(&point);

    if option_value.is_none() {
        return false;
    }
    let v = option_value.unwrap();
    if v == &9 {
        return false;
    }
    find_neighbours(point).iter().fold(true, |acc, np| {
        let mut neighbour_ok = true;
        if let Some(neighbour_value) = map.get(&np) {
            neighbour_ok = exclude.contains(&np) || *neighbour_value >= *v;
        }
        acc && neighbour_ok
    })
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

fn get_colors() -> Vec<color::AnsiValue> {
    let mut rng = thread_rng();
    let mut y = (0..100)
        .map(|i| color::AnsiValue((i * 255 / 100) as u8))
        .collect::<Vec<color::AnsiValue>>();
    y.shuffle(&mut rng);
    y
}


fn print_map(
    x_size: usize,
    y_size: usize,
    smoke_map: &HashMap<Point, u32>,
    basins: &Vec<Vec<Point>>,
    colors: &Vec<color::AnsiValue>
) {
    for y in 0..y_size {
        for x in 0..x_size {
            let point = Point {
                x: x as i32,
                y: y as i32,
            };
            let value = smoke_map.get(&point).unwrap();
            let containing_basin = basins.iter().position(|zone| zone.contains(&point));
            if !containing_basin.is_none() {
                let color = colors[containing_basin.unwrap() % 100];
                print!("{}{}{}", color::Fg(color), value, style::Reset);
            } else if value == &9 {
                print!("{}{}{}", color::Fg(color::Black), value, style::Reset);
            } else {
                print!("{}", value);
            }
        }
        print!("\n");
    }
}
