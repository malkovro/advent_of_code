use std::collections::HashMap;
use std::env;

#[derive(Debug, PartialEq, Eq, Hash, Copy, Clone)]
struct Pairing(char, char);

fn main() {
    let args: Vec<String> = env::args().collect();
    let filename = args[1].clone();

    let input = read_into_vec(&filename);
    let (mut hash, rules) = read_input(input);

    for _ in 0..40 {
        polymerize(&mut hash, &rules);
    }

    let mut elements_count = count_elements(&hash);
   elements_count.sort();
    println!("{:?}", elements_count[elements_count.len() -1] - elements_count[0]);
}

fn count_elements(hash: &HashMap<Pairing, usize>) -> Vec<usize> {
    let elements = hash.keys().fold(vec![], |mut elements, pairing| {
        if !elements.contains(&pairing.0) {
            elements.push(pairing.0);
        }
        if !elements.contains(&pairing.1) {
            elements.push(pairing.1);
        }
        elements
    });

    elements
        .iter()
        .map(|el| {
            let count_as_start = hash
                .iter()
                .filter(|(k, _)| k.0 == *el)
                .map(|(_, v)| v)
                .sum();
            let count_as_end = hash
                .iter()
                .filter(|(k, _)| k.1 == *el)
                .map(|(_, v)| v)
                .sum();
            std::cmp::max(count_as_start, count_as_end)
        })
        .collect()
 
}
fn polymerize(polymer_hash: &mut HashMap<Pairing, usize>, rules: &HashMap<Pairing, char>) {
    let pcis = rules.iter().fold(vec![], |mut pcis, (p, i)| {
        if let Some((_, count)) = polymer_hash.remove_entry(p) {
            pcis.push((p, count, i));
        }
        pcis
    });

    for (pairing, pairing_count, injection) in pcis {
        let injection_pair_a = Pairing(pairing.0, *injection);
        let injection_pair_b = Pairing(*injection, pairing.1);

        let pair_a_in_polymer = polymer_hash.entry(injection_pair_a).or_insert(0);
        *pair_a_in_polymer += pairing_count;
        let pair_b_in_polymer = polymer_hash.entry(injection_pair_b).or_insert(0);
        *pair_b_in_polymer += pairing_count;
    }
}

fn read_input(input: Vec<String>) -> (HashMap<Pairing, usize>, HashMap<Pairing, char>) {
    let mut hash = HashMap::new();

    let polymer_chain_length = input[0].len();
    for i in 0..(polymer_chain_length - 1) {
        let pairing_c = hash
            .entry(Pairing(
                input[0].chars().nth(i).unwrap(),
                input[0].chars().nth(i + 1).unwrap(),
            ))
            .or_insert(0);
        *pairing_c += 1
    }

    let rules =
        input[1]
            .split('\n')
            .filter(|s| !s.is_empty())
            .fold(HashMap::new(), |mut hash, l| {
                let line = String::from(l);
                hash.insert(
                    Pairing(l.chars().nth(0).unwrap(), line.chars().nth(1).unwrap()),
                    line.chars().nth(6).unwrap(),
                );
                hash
            });
    (hash, rules)
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
