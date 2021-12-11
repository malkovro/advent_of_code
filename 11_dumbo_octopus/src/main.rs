use std::collections::{HashMap, HashSet};
use std::env;
use termion::{color, cursor, style};

#[derive(PartialEq, Eq, Clone, Copy, Debug, Hash)]
struct Point {
    x: i32,
    y: i32,
}

fn build_octopus_map(input: Vec<String>) -> HashMap<Point, u32> {
    let mut map: HashMap<Point, u32> = HashMap::new();

    for (y, line) in input.iter().enumerate() {
        for (x, v) in line.chars().map(|v| v.to_digit(10).unwrap()).enumerate() {
            map.insert(
                Point {
                    x: x as i32,
                    y: y as i32,
                },
                v,
            );
        }
    }
    map
}

fn main() {
    let args: Vec<String> = env::args().collect();
    let filename = args[1].clone();
    let steps = args[2].clone().parse::<usize>().unwrap();

    let lines = read_into_vec(&filename);
    let line_count = lines.len();
    let mut map = build_octopus_map(lines);
    let line_size = map.len() / line_count;
    let mut total_flashes = 0;
    let mut iteration = 0;
    let mut first_simultaneous_flashing_at = 0;

    println!("Part 1 - Total flashes count after {} steps: ??", steps);
    println!("Part 2 - Simultaneous flashing on iteration: ??");
    print_map(&map, &line_size, &line_count);

    loop {
        iteration += 1;
        map = life_goes_on(map, &line_size, &line_count);
        print!("{}", cursor::Up(line_count as u16));
        print_map(&map, &line_size, &line_count);
        std::thread::sleep(std::time::Duration::from_millis(20));
        let flashes = count_flashes(&map, &line_size, &line_count);
        total_flashes += flashes;
        if flashes == map.keys().len() && first_simultaneous_flashing_at == 0 {
            first_simultaneous_flashing_at = iteration;
            println!(
                "{}Part 2 - Simultaneous flashing on iteration: {}{}",
                cursor::Up(line_count as u16 + 1),
                first_simultaneous_flashing_at,
                cursor::Down(line_count as u16)
            );

            if iteration > steps {
                return;
            }
        }
        if iteration == steps {
            println!(
                "{}Part 1 - Total flashes count after {} steps: {}{}",
                cursor::Up(line_count as u16 + 2),
                steps,
                total_flashes,
                cursor::Down(line_count as u16 + 1)
            );
            if first_simultaneous_flashing_at != 0 {
                return;
            }
        }
    }
}

fn count_flashes(map: &HashMap<Point, u32>, line_size: &usize, line_count: &usize) -> usize {
    (0..*line_count).fold(0, |c, y| {
        c + (0..*line_size)
            .filter(|x| {
                map.get(&Point {
                    x: *x as i32,
                    y: y as i32,
                })
                .unwrap()
                    == &0
            })
            .count()
    })
}

fn life_goes_on(
    mut map: HashMap<Point, u32>,
    line_size: &usize,
    line_count: &usize,
) -> HashMap<Point, u32> {
    let mut to_inc = (0..*line_count).fold(vec![], |mut acc, y| {
        for x in 0..*line_size {
            acc.push(Point {
                x: x as i32,
                y: y as i32,
            });
        }
        acc
    });

    let mut flashing = HashSet::new();
    loop {
        let point_w = to_inc.pop();
        if point_w.is_none() {
            return map;
        }
        let point = point_w.unwrap();

        if flashing.contains(&point) {
            continue;
        }
        let energy = map.entry(point).or_insert(0);
        *energy = (*energy + 1) % 10;
        if energy == &0 {
            flashing.insert(point);
            let mut neighbours = find_neighbours(&point, &map);
            to_inc.append(&mut neighbours);
        }
    }
}

fn find_neighbours(point: &Point, map: &HashMap<Point, u32>) -> Vec<Point> {
    let neighbours = vec![
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
        Point {
            x: point.x - 1,
            y: point.y - 1,
        },
        Point {
            x: point.x + 1,
            y: point.y + 1,
        },
        Point {
            x: point.x - 1,
            y: point.y + 1,
        },
        Point {
            x: point.x + 1,
            y: point.y - 1,
        },
    ];

    neighbours
        .into_iter()
        .filter(|x| (*map).contains_key(x))
        .collect()
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

fn print_map(map: &HashMap<Point, u32>, x_size: &usize, y_size: &usize) {
    for y in 0..*y_size {
        for x in 0..*x_size {
            let point = Point {
                x: x as i32,
                y: y as i32,
            };
            let value = map.get(&point).unwrap();
            if value == &0 {
                print!("{}{}{}", color::Fg(color::Green), "💡", style::Reset);
            } else {
                print!(" {}", value);
            }
        }
        print!("\n");
    }
}
