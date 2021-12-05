use regex::Regex;
use std::fmt::{Display, Error, Formatter};

#[derive(Debug)]
struct Point(i32, i32);

#[derive(Debug)]
struct Segment {
    start: Point,
    end: Point,
}

impl Point {
    fn cross(&self, vec: &Point) -> i32 {
        self.0 * (*vec).1 - self.1 * (*vec).0
    }

    fn dot(&self, vec: &Point) -> i32 {
        self.0 * (*vec).0 + self.1 * (*vec).1
    }
}

impl Segment {
    fn is_horizontal(&self) -> bool {
        self.start.1 == self.end.1
    }
    fn is_vertical(&self) -> bool {
        self.start.0 == self.end.0
    }
    fn contains(&self, p: Point) -> bool {
        let sp = Point(p.0 - self.start.0, p.1 - self.start.1);
        let se = Point(self.end.0 - self.start.0, self.end.1 - self.start.1);
        // println!("SP: {:?}\nSE: {:?} | Cross: {}, Dot: {}", sp, se, se.cross(&sp), se.dot(&sp));
        let se_sp_dot = se.dot(&sp);
        let contains = se.cross(&sp) == 0 && se_sp_dot >= 0 && se_sp_dot <= se.dot(&se);
        if contains {
            // println!("{:?} contains {:?}", self, p);
        }
        contains
    }
}

struct Grid(Vec<Vec<usize>>);

impl Display for Grid {
    fn fmt(&self, f: &mut Formatter) -> Result<(), Error> {
        let mut formatted = String::new();

        for line in &self.0[0..self.0.len() - 1] {
            formatted.push_str(&format!("{:?}", line));
            formatted.push_str("\n");
        }
        write!(f, "{}", formatted)
    }
}

pub fn detect(filename: &str) {
    let segments_def = read_into_vec(&filename);
    let segments = load_segments(segments_def);
    let max_x = segments
        .iter()
        .map(|s| {
            if s.start.0 > s.end.0 {
                s.start.0
            } else {
                s.end.0
            }
        })
        .max()
        .unwrap()
        + 1;
    let max_y = segments
        .iter()
        .map(|s| {
            if s.start.1 > s.end.1 {
                s.start.1
            } else {
                s.end.1
            }
        })
        .max()
        .unwrap()
        + 1;
    println!("Grid: {} x {}", max_x, max_y);
    let mut grid = vec![vec![0; max_x as usize]; max_y as usize];
    for x in 0..max_x {
        for y in 0..max_y {
            // println!("Checking: {},{}", x, y);
            grid[y as usize][x as usize] = segments
                .iter()
                .filter(|s| s.contains(Point(x as i32, y as i32)))
                .count()
        }
    }
    // println!("{}", Grid(grid.clone()));
    println!(
        "{:?}",
        grid.iter().fold(0, |acc, row| acc
            + (*row).iter().fold(0, |acc, cell| if *cell > 1 {
                acc + 1
            } else {
                acc
            }))
    );
}

fn load_segments(defs: Vec<String>) -> Vec<Segment> {
    let re = Regex::new(r"(\d+),(\d+) -> (\d+),(\d+)").unwrap();
    let segments: Vec<Segment> = defs
        .iter()
        .map(|def| {
            let caps = re.captures(def).unwrap();
            Segment {
                start: Point(
                    caps.get(1).unwrap().as_str().parse::<i32>().unwrap(),
                    caps.get(2).unwrap().as_str().parse::<i32>().unwrap(),
                ),
                end: Point(
                    caps.get(3).unwrap().as_str().parse::<i32>().unwrap(),
                    caps.get(4).unwrap().as_str().parse::<i32>().unwrap(),
                ),
            }
        })
        .collect();

    segments
        .into_iter()
        // .filter(|s| (*s).is_vertical() || (*s).is_horizontal()) // Uncomment for Part 2
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
