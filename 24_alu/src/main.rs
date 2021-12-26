use regex::Regex;
use std::collections::HashMap;
use std::env::args;

fn main() {
    let args: Vec<String> = args().collect();
    let filename = args[1].clone();
    let instructions = read_monad(&filename);

    println!(
        "Part 1: {}\nPart 2: {}",
        solve(0, &instructions, true, 0, &mut HashMap::new())
            .unwrap()
            .to_string()
            .chars()
            .rev()
            .collect::<String>(),
        solve(0, &instructions, false, 0, &mut HashMap::new())
            .unwrap()
            .to_string()
            .chars()
            .rev()
            .collect::<String>()
    );
}

fn solve(
    z: i64,
    instructions: &Vec<Vec<Instruction>>,
    highest: bool,
    layer: usize,
    m: &mut HashMap<(i64, usize), Option<i64>>,
) -> Option<i64> {
    if let Some(res) = m.get(&(z, layer)) {
        return *res;
    }

    let range = if highest {
        [9, 8, 7, 6, 5, 4, 3, 2, 1]
    } else {
        [1, 2, 3, 4, 5, 6, 7, 8, 9]
    };
    for d in range {
        let mut n_store = vec![d, 0, 0, z];
        for ins in instructions[layer].iter() {
            ins.apply(&mut n_store);
        }
        let n_z = n_store[3];
        if layer == instructions.len() - 1 {
            if n_z == 0 {
                m.insert((0, layer), Some(d));
                return Some(d);
            }
            continue;
        }
        if let Some(solution) = solve(n_z, &instructions, highest, layer + 1, m) {
            m.insert((z, layer), Some(d + 10 * solution));
            return Some(10 * solution + d);
        }
    }
    m.insert((z, layer), None);
    None
}

#[derive(Debug)]
enum InstArg {
    Number(i64),
    Var(usize),
}

#[derive(Debug)]
enum Instruction {
    Add(usize, InstArg),
    Mul(usize, InstArg),
    Div(usize, InstArg),
    Mod(usize, InstArg),
    Eql(usize, InstArg),
}

impl InstArg {
    fn value(&self, storage: &Vec<i64>) -> i64 {
        match self {
            InstArg::Number(v) => *v,
            InstArg::Var(i) => storage[*i],
        }
    }
}

impl Instruction {
    fn apply(&self, storage: &mut Vec<i64>) {
        // println!("Before: {:?}", storage);
        match self {
            Instruction::Add(i, arg) => {
                // println!("Adding {} to {}", arg.value(storage), i);
                storage[*i] += arg.value(storage);
            }

            Instruction::Mul(i, arg) => {
                // println!("Mukt {} to {}", arg.value(storage), i);
                storage[*i] *= arg.value(storage);
            }

            Instruction::Div(i, arg) => {
                // println!("Dividing {} by {}", arg.value(storage), i);
                storage[*i] /= arg.value(storage);
            }

            Instruction::Mod(i, arg) => {
                // println!("Moduling {} and {}", arg.value(storage), i);
                storage[*i] = storage[*i] % arg.value(storage);
            }

            Instruction::Eql(i, arg) => {
                // println!("Checking eq  of {} and {}", arg.value(storage), i);
                storage[*i] = if storage[*i] == arg.value(storage) {
                    1
                } else {
                    0
                }
            }
        }
    }
}

fn read_monad(filename: &str) -> Vec<Vec<Instruction>> {
    let re = Regex::new(r"(inp|add|mul|div|mod|eql) (.) (.+)").unwrap();
    std::fs::read_to_string(filename)
        .unwrap()
        .split('\n')
        .filter(|s| !s.is_empty())
        .fold(vec![], |mut acc, x| {
            if x == "inp w" {
                acc.push(vec![]);
                return acc;
            }
            let capture = re.captures_iter(x).next().unwrap();
            let variable = match &capture[2] {
                "w" => 0,
                "x" => 1,
                "y" => 2,
                "z" => 3,
                _ => unreachable!(),
            };
            let arg = match capture[3].parse::<i64>() {
                Ok(n) => InstArg::Number(n),
                _ => InstArg::Var(match &capture[3] {
                    "w" => 0,
                    "x" => 1,
                    "y" => 2,
                    "z" => 3,
                    _ => unreachable!(),
                }),
            };
            let instruction = match &capture[1] {
                "add" => Instruction::Add(variable, arg),
                "mul" => Instruction::Mul(variable, arg),
                "div" => Instruction::Div(variable, arg),
                "mod" => Instruction::Mod(variable, arg),
                "eql" => Instruction::Eql(variable, arg),
                _ => unreachable!(),
            };
            let len = acc.len();
            acc[len - 1].push(instruction);
            acc
        })
}
