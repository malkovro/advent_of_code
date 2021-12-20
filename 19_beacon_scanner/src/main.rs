use itertools::Itertools;
use regex::Regex;
use std::collections::HashSet;
use std::env::args;

struct Orientation {
    points: Vec<[i32; 3]>,
    signs: [[i32; 3]; 8],
    sign_index: usize,
    rotations: [[usize; 3]; 6],
    rotation_index: usize,
}

impl Orientation {
    fn new(points: Vec<[i32; 3]>) -> Self {
        Self {
            points,
            signs: [
                [1, 1, 1],
                [1, 1, -1],
                [1, -1, 1],
                [1, -1, -1],
                [-1, 1, 1],
                [-1, 1, -1],
                [-1, -1, 1],
                [-1, -1, -1],
            ],
            sign_index: 0,
            rotations: [
                [0, 1, 2],
                [0, 2, 1],
                [1, 0, 2],
                [1, 2, 0],
                [2, 0, 1],
                [2, 1, 0],
            ],
            rotation_index: 0,
        }
    }
}

impl Iterator for Orientation {
    type Item = Vec<[i32; 3]>;

    fn next(&mut self) -> Option<Self::Item> {
        if self.sign_index >= self.signs.len() {
            return None;
        }
        let sign = self.signs[self.sign_index];
        let rotation = self.rotations[self.rotation_index];
        let points = self
            .points
            .iter()
            .map(|p| {
                [
                    p[rotation[0]] * sign[0],
                    p[rotation[1]] * sign[1],
                    p[rotation[2]] * sign[2],
                ]
            })
            .collect();
        if self.rotation_index == self.rotations.len() - 1 {
            self.rotation_index = 0;
            self.sign_index += 1;
        } else {
            self.rotation_index += 1;
        }
        Some(points)
    }
}

fn substract(source: &[i32], target: &[i32]) -> [i32; 3] {
    [
        target[0] - source[0],
        target[1] - source[1],
        target[2] - source[2],
    ]
}

fn add(source: &[i32], target: &[i32]) -> [i32; 3] {
    [
        target[0] + source[0],
        target[1] + source[1],
        target[2] + source[2],
    ]
}
fn mergeable(
    scan_a: &Vec<[i32; 3]>,
    scan_b: &Vec<[i32; 3]>,
) -> Option<([i32; 3], Vec<[i32; 3]>)> {
    for i in 0..scan_a.len() {
        let pos_a = scan_a
            .iter()
            .map(|p| substract(&scan_a[i], p))
            .collect::<HashSet<_>>();
        for j in 0..scan_b.len() {
            let pos_b = scan_b
                .iter()
                .map(|p| substract(&scan_b[j], p))
                .collect::<Vec<_>>();
            let overlaps = pos_b
                .iter()
                .map(|p| pos_a.contains(p) as usize)
                .sum::<usize>();
            if overlaps >= 12 {
                let mut merged = scan_a.iter().map(|e| *e).collect::<HashSet<_>>();
                for p in &pos_b {
                    merged.insert(add(&scan_a[i], p));
                }
                let scanner_b_pos = substract(&scan_b[j], &scan_a[i]);
                return Some((scanner_b_pos, merged.into_iter().collect()));
            }
        }
    }
    None
}


fn main() {
    let args: Vec<String> = args().collect();
    let filename = args[1].clone();
    let mut scan_zones = read_input(&filename);
    let mut beacons = scan_zones[0].clone().unwrap();
    let mut scanners = vec![[0, 0, 0]];

    while scanners.len() != scan_zones.len() {
        for i in 1..scan_zones.len() {
            if let Some(scan_zone) = scan_zones[i].clone() {
                for o in Orientation::new(scan_zone.clone()) {
                    if let Some((scanner_pos, new_beacons)) = mergeable(&beacons, &o) {
                        beacons = new_beacons;
                        scanners.push(scanner_pos);
                        scan_zones[i] = None;
                        break;
                    }
                }
            }
        }
    }

    let max_distance = scanners.into_iter().combinations(2).map(|scanner_pair| {
        let dist = substract(&scanner_pair[0], &scanner_pair[1]);
        dist[0].abs() + dist[1].abs() + dist[2].abs()
    }).max().unwrap();

    println!("Counted {} beacons. Farthest scanners are {} away of each other", beacons.len(), max_distance);
}

fn read_input(filename: &str) -> Vec<Option<Vec<[i32; 3]>>> {
    let re = Regex::new(r"(-?\d+),(-?\d+),(-?\d+)").unwrap();

    std::fs::read_to_string(filename)
        .unwrap()
        .split('\n')
        .into_iter()
        .fold((vec![], vec![]), |(mut scans, mut curr_scan), l| {
            if l.contains("scanner") {
                return (scans, curr_scan);
            }
            if l.is_empty() {
                scans.push(Some(curr_scan));
                return (scans, vec![]);
            }
            for cap in re.captures_iter(l) {
                curr_scan.push([
                    cap[1].parse::<i32>().unwrap(),
                    cap[2].parse::<i32>().unwrap(),
                    cap[3].parse::<i32>().unwrap(),
                ]);
            }
            return (scans, curr_scan);
        })
        .0
}
