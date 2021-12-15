use std::cmp::{Ord, Ordering};
use std::collections::BinaryHeap;
use std::env;

#[derive(Debug, PartialEq, Eq, Copy, Clone, Hash)]
struct Pos(usize, usize);

impl Ord for Pos {
    fn cmp(&self, other: &Self) -> Ordering {
        (other.0 + other.1).cmp(&(self.0 + self.1))
    }
}

impl PartialOrd for Pos {
    fn partial_cmp(&self, other: &Self) -> Option<Ordering> {
        Some(self.cmp(other))
    }
}

fn main() {
    let args: Vec<String> = env::args().collect();
    let filename = args[1].clone();
    let map_tile = read_into_vec(&filename);
    let map = build_real_map(map_tile); // Part 2
    let x_size = map[0].len();
    let y_size = map.len();
    let goal = Pos(x_size - 1, y_size - 1);
    let start = Pos(0, 0);
    let mut open_set = BinaryHeap::from([start]);
    let mut g_costs = vec![vec![u32::MAX; x_size]; y_size];
    g_costs[0][0] = 0;

    while let Some(current) = open_set.pop() {
        if current == goal {
            println!("We out: {}!", g_costs[y_size - 1][x_size - 1]);
            return;
        }
        for neighbor in neighbours(&current, &x_size, &y_size).iter() {
            let tent = g_costs[current.1][current.0] + map[neighbor.1][neighbor.0];
            if tent < g_costs[neighbor.1][neighbor.0] {
                g_costs[neighbor.1][neighbor.0] = tent;
                open_set.push(*neighbor);
            }
        }
    }
}

fn neighbours(pos: &Pos, x_size: &usize, y_size: &usize) -> Vec<Pos> {
    let mut ns = vec![];
    if pos.0 < x_size - 1 {
        ns.push(Pos(pos.0 + 1, pos.1));
    }
    if pos.1 < y_size - 1 {
        ns.push(Pos(pos.0, pos.1 + 1));
    }
    if pos.0 > 0 {
        ns.push(Pos(pos.0 - 1, pos.1));
    }
    if pos.1 > 0 {
        ns.push(Pos(pos.0, pos.1 - 1));
    }
    ns
}

fn build_real_map(map_tile: Vec<Vec<u32>>) -> Vec<Vec<u32>> {
    let x_size = map_tile[0].len();
    let y_size = map_tile.len();
    let mut map = vec![vec![0; x_size * 5]; y_size * 5];

    for ym_i in 0..=4 {
        for yi in 0..y_size {
            let y = ym_i * y_size + yi;
            for xm_i in 0..=4 {
                for xi in 0..x_size {
                    map[y][xm_i * x_size + xi] =
                        (map_tile[yi][xi] + ym_i as u32 + xm_i as u32 - 1) % 9 + 1
                }
            }
        }
    }
    map
}

fn read_into_vec(file_name: &str) -> Vec<Vec<u32>> {
    std::fs::read_to_string(file_name)
        .unwrap()
        .split('\n')
        .into_iter()
        .filter(|s| !s.is_empty())
        .map(|l| l.chars().map(|x| x.to_digit(10).unwrap()).collect())
        .collect()
}
