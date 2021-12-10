use std::collections::HashMap;
use std::env;

trait SyntaxElement {
    fn is_closing(&self) -> bool;
    fn closes(&self, el: Self) -> bool;
}

impl SyntaxElement for char {
    fn is_closing(&self) -> bool {
        self == &')' || self == &']' || self == &'>' || self == &'}'
    }

    fn closes(&self, el: Self) -> bool {
        match self {
            &')' => el == '(',
            &']' => el == '[',
            &'}' => el == '{',
            &'>' => el == '<',
            _ => panic!("Unexpected char {}", self),
        }
    }
}

fn main() {
    let args: Vec<String> = env::args().collect();
    let filename = args[1].clone();
    let lines = read_into_vec(&filename);

    let mut illegal_chars: HashMap<char, usize> = HashMap::new();
    let mut incomplete_lines_scores = vec![];

    'outer: for line in lines.iter() {
        let mut stack: Vec<char> = vec![];
        for c in line.chars() {
            if c.is_closing() {
                let last_el = stack.pop();
                if last_el.is_none() || !c.closes(last_el.unwrap()) {
                    let counter = illegal_chars.entry(c).or_insert(0);
                    *counter += 1;
                    continue 'outer;
                }
            } else {
                stack.push(c);
            }
        }
        // We are here so we have an incomplete line
        let mut line_score: u64 = 0;
        while stack.len() > 0 {
            line_score *= 5;
            match stack.pop().unwrap() {
                '(' => line_score += 1,
                '[' => line_score += 2,
                '{' => line_score += 3,
                '<' => line_score += 4,
                _ => panic!("Unexpected awaiting close char"),
            }
        }
        incomplete_lines_scores.push(line_score);
    }

    incomplete_lines_scores.sort();

    let score = 3 * *illegal_chars.entry(')').or_insert(0)
        + 57 * *illegal_chars.entry(']').or_insert(0)
        + 1197 * *illegal_chars.entry('}').or_insert(0)
        + 25137 * *illegal_chars.entry('>').or_insert(0);
    println!("Score of illegal chars: {}", score);
    println!(
        "Autocomplete score: {}",
        incomplete_lines_scores[incomplete_lines_scores.len() / 2]
    );
}

fn read_into_vec(filename: &str) -> Vec<String> {
    std::fs::read_to_string(filename)
        .unwrap()
        .split('\n')
        .into_iter()
        .filter(|s| !s.is_empty())
        .map(String::from)
        .collect()
}
