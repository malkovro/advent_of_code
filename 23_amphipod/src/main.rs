use std::collections::{HashMap, BinaryHeap};

fn right_configuration(maze: &Vec<Vec<char>>) -> bool {
    maze[2..(maze.len() - 1)]
        .iter()
        .all(|l| l[3] == 'A' && l[5] == 'B' && l[7] == 'C' && l[9] == 'D')
}

type Point = (usize, usize);

trait MazeMover {
    fn can_move_out(&self, maze: &Vec<Vec<char>>) -> Vec<(Point, i64)>;
    fn can_move_into_room(&self, maze: &Vec<Vec<char>>) -> Option<(Point, i64)>;
    fn target_room_j(&self, maze: &Vec<Vec<char>>) -> Option<usize>;
    fn target_room_i(&self, maze: &Vec<Vec<char>>) -> Option<usize>;
    fn unitary_cost(&self, maze: &Vec<Vec<char>>) -> i64;
}

impl MazeMover for Point {
    fn can_move_out(&self, maze: &Vec<Vec<char>>) -> Vec<(Point, i64)> {
        let destinations = [1, 2, 4, 6, 8, 10, 11];
        let mut alt = vec![];
        if (1..self.0).any(|i| maze[i][self.1] != '.') {
            return alt;
        }
        let goal = match maze[self.0][self.1] {
            'A' => 3,
            'B' => 5,
            'C' => 7,
            'D' => 9,
            _ => unreachable!(),
        };
        if self.1 == goal
            && (self.0..(maze.len() - 1)).all(|i| maze[i][self.1] == maze[self.0][self.1])
        {
            return alt;
        }

        let ucost = self.unitary_cost(maze);
        let cost = (self.0 as i64 - 1) * ucost;

        for d in destinations {
            let move_authorized = if self.1 < d {
                (self.1..(d + 1)).all(|j| maze[1][j] == '.')
            } else {
                (d..self.1).all(|j| maze[1][j] == '.')
            };
            if move_authorized {
                alt.push(((1, d), cost + ((d as i64 - self.1 as i64).abs() * ucost)));
            }
        }

        alt
    }

    fn can_move_into_room(&self, maze: &Vec<Vec<char>>) -> Option<(Point, i64)> {
        if let Some(target_j) = self.target_room_j(maze) {
            if let Some(target_i) = self.target_room_i(maze) {
                let cost = self.unitary_cost(maze)
                    * (target_i as i64 - self.0 as i64 + (self.1 as i64 - target_j as i64).abs());
                return Some(((target_i, target_j), cost));
            }
        }
        None
    }

    fn target_room_i(&self, maze: &Vec<Vec<char>>) -> Option<usize> {
        let target_j = self.target_room_j(maze)?;
        let col_size = maze.len();

        let target_i = (2..(col_size - 1))
            .take_while(|&i| maze[i][target_j] == '.')
            .last();

        if target_i.is_some()
            && target_i.unwrap() < col_size - 2
            && ((target_i.unwrap() + 1)..(col_size - 1))
                .any(|i| maze[i][target_j] != maze[self.0][self.1])
        {
            return None;
        }

        return target_i;
    }

    fn target_room_j(&self, maze: &Vec<Vec<char>>) -> Option<usize> {
        let target_j = match maze[self.0][self.1] {
            'A' => 3,
            'B' => 5,
            'C' => 7,
            'D' => 9,
            _ => return None,
        };
        let mut j_range = if self.1 < target_j {
            (self.1 + 1)..(target_j + 1)
        } else {
            (target_j)..(self.1)
        };

        if j_range.all(|j| maze[self.0][j] == '.') {
            return Some(target_j);
        }

        None
    }

    fn unitary_cost(&self, maze: &Vec<Vec<char>>) -> i64 {
        match maze[self.0][self.1] {
            'A' => 1,
            'B' => 10,
            'C' => 100,
            'D' => 1000,
            _ => unreachable!(),
        }
    }
}

fn next_moves(maze: &Vec<Vec<char>>) -> Vec<(i64, Vec<Vec<char>>)> {
    let col_size = maze.len();
    let row_size = maze[1].len();
    let mut mvs = vec![];

    for j in 1..(row_size - 1) {
        let amphipod = match maze[1][j] {
            '.' => continue,
            '#' => continue,
            any => any,
        };
        if let Some(((ni, nj), cost)) = (1, j).can_move_into_room(maze) {
            let mut cloned_maze = maze.clone();
            cloned_maze[ni][nj] = amphipod;
            cloned_maze[1][j] = '.';
            mvs.push((cost, cloned_maze));
        }
    }
    for i in 2..(col_size - 1) {
        for j in [3, 5, 7, 9] {
            let amphipod = match maze[i][j] {
                '.' => continue,
                '#' => continue,
                any => any,
            };
            for ((ni, nj), cost) in (i, j).can_move_out(maze) {
                let mut cloned_maze = maze.clone();
                cloned_maze[ni][nj] = amphipod;
                cloned_maze[i][j] = '.';
                mvs.push((cost, cloned_maze));
            }
        }
    }
    mvs
}

fn print_map(map: &Vec<Vec<char>>) {
    for l in map.iter() {
        for e in l.iter() {
            print!("{}", e);
        }
        println!("");
    }
}

fn shortest_path(maze: &Vec<Vec<char>>) -> i64 {
    let mut dist = HashMap::new();
    let mut q = BinaryHeap::new();
    let mut it = 0;
    q.push((0, maze.clone()));
    while let Some((cost, m)) = q.pop() {
        it += 1;
        if right_configuration(&m) {
            println!(
                "Finished after iterating {} - PriorityQueue size {} - HashMap size {}",
                it,
                q.len(),
                dist.len()
            );
            return -cost;
        }
        if let Some(&c) = dist.get(&m) {
            if -cost > c {
                continue;
            }
        }
        let moves = next_moves(&m);
        for (neighbor_cost, m) in moves {
            let next_cost = -cost + neighbor_cost;
            let &c = dist.get(&m).unwrap_or(&i64::MAX);
            if c > next_cost {
                dist.insert(m.clone(), next_cost);
                q.push((-next_cost, m));
            }
        }
    }
    unreachable!()
}

fn main() {
    let mut map: Vec<Vec<char>> = read_map("input.txt");
    let p1 = shortest_path(&map);
    map.insert(3, "  #D#B#A#C#  ".chars().collect());
    map.insert(3, "  #D#C#B#A#  ".chars().collect());
    let p2 = shortest_path(&map);
    println!("Part 1: {}\nPart 2: {}", p1, p2);
}

fn read_map(filename: &str) -> Vec<Vec<char>> {
    std::fs::read_to_string(filename)
        .unwrap()
        .split('\n')
        .filter(|s| !s.is_empty())
        .map(|l| l.chars().collect())
        .collect()
}
