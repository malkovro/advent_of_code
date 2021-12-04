use regex::Regex;
use std::collections::HashMap;
use std::vec::Vec;

#[derive(Debug)]
struct Cell {
    number: i16,
    drawn: bool,
}

#[derive(Debug)]
struct Grid {
    rows: Vec<Vec<Cell>>,
    qa: HashMap<i16, (usize, usize)>,
    won: bool,
}

impl Grid {
    fn mark_drawn(&mut self, number: i16) -> Option<(usize, usize)> {
        match self.qa.get(&number) {
            Some((i, j)) => {
                self.rows[*j][*i].drawn = true;
                Some((*i, *j))
            }
            None => None,
        }
    }

    fn bingo(&self, i: usize, j: usize) -> bool {
        // println!("Checking bingo with {}, {}", i, j);
        (0..self.rows.len()).fold(true, |acc, j| acc && self.rows[j][i].drawn)
            || (0..self.rows[j].len()).fold(true, |acc, i| acc && self.rows[j][i].drawn)
    }

    fn unmarked_sum(&self) -> i16 {
        self.qa.iter().fold(0, |acc, (key, (i, j))| {
            if !self.rows[*j][*i].drawn {
                acc + key
            } else {
                acc
            }
        })
    }
}

pub fn play(filename: &str, silent: bool) -> i32 {
    let lines = read_into_vec(&filename);

    let mut grids: Vec<Grid> = lines
        .iter()
        .enumerate()
        .filter(|(index, _)| *index != 0)
        .map(|(_, value)| value.to_string())
        .collect::<Vec<String>>()
        .iter()
        .map(construct_grid)
        .collect();

    let drawn_numbers: Vec<i16> = lines[0]
        .split(',')
        .into_iter()
        .map(|s| s.parse::<i16>().unwrap())
        .collect();
    // println!("{:?}", drawn_numbers);
    for number in drawn_numbers {
        if !silent {
            println!("Drawing numberrr {}!!!", number);
        }
        for grid_index in 0..grids.len() {
            if grids[grid_index].won {
                continue;
            }
            let optional_pos = grids[grid_index].mark_drawn(number);
            if let Some((i, j)) = optional_pos {
                if grids[grid_index].bingo(i, j) {
                    grids[grid_index].won = true;
                    let unmarked_sum = grids[grid_index].unmarked_sum();
                    if !silent {
                        println!(
                            "We've got a winner! Score: {}",
                            unmarked_sum as i32 * number as i32
                        );
                    }
                    if grids.iter().fold(true, |acc, g| acc && g.won) {
                        if !silent {
                            println!("Everybody won, it's safe to take this board to let pulpik win!");
                        }
                        return unmarked_sum as i32 * number as i32
                    }
                }
            }
        }
    }
    0
}

fn construct_grid(lines: &String) -> Grid {
    let rows: Vec<Vec<Cell>> = lines
        .split('\n')
        .into_iter()
        .map(|g| {
            Regex::new(r"(\s+)")
                .unwrap()
                .split(g.trim())
                .map(|l| Cell {
                    number: l.trim().parse::<i16>().unwrap(),
                    drawn: false,
                })
                .collect()
        })
        .collect();

    let mut qa = HashMap::new();
    for j in 0..rows.len() {
        let row = &rows[j];
        for i in 0..row.len() {
            qa.insert(row[i].number, (i, j));
        }
    }
    Grid {
        rows: rows,
        qa: qa,
        won: false,
    }
}

fn read_into_vec(file_name: &str) -> Vec<String> {
    std::fs::read_to_string(file_name)
        .unwrap()
        .split("\n\n")
        .into_iter()
        .filter(|s| !s.is_empty())
        .map(String::from)
        .collect()
}
