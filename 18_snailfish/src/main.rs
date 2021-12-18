extern crate nom;

use itertools::Itertools;
use std::env;

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
struct Node(u32, u32);

fn flat_reduced_sum(snails: &Vec<Vec<Node>>) -> Vec<Node> {
    let len = snails.len();
    let mut snail_sum_flat = snails[0].clone();
    for i in 0..(len - 1) {
        let mut snail_to_add = snails[1 + i]
            .clone()
            .into_iter()
            .map(|mut x| {
                x.1 += 1;
                x
            })
            .collect::<Vec<Node>>();
        snail_sum_flat = snail_sum_flat
            .clone()
            .into_iter()
            .map(|mut x| {
                x.1 += 1;
                x
            })
            .collect::<Vec<Node>>();
        snail_sum_flat.append(&mut snail_to_add);
        reduce(&mut snail_sum_flat)
    }
    snail_sum_flat
}

fn main() {
    let args: Vec<String> = env::args().collect();
    let filename = args[1].clone();
    let flat_snails = flat_trees(&filename);
    let snail_sum_flat = flat_reduced_sum(&flat_snails);
    println!(
        "Magnitude of the Sum of snails: {:?}",
        magnitude(&snail_sum_flat)
    );

    let best_sum_magnitude = flat_snails
        .into_iter()
        .permutations(2)
        .map(|pairs| magnitude(&flat_reduced_sum(&pairs)))
        .max()
        .unwrap();

    println!("Best Magnitude of sum {:?}", best_sum_magnitude);
}

fn flat_trees(filename: &str) -> Vec<Vec<Node>> {
    std::fs::read_to_string(filename)
        .unwrap()
        .split('\n')
        .into_iter()
        .filter(|s| !s.is_empty())
        .map(|l| {
            l.chars()
                .fold((vec![], 0), |(mut nodes, depth), c| match c {
                    '[' => (nodes, depth + 1),
                    ']' => (nodes, depth - 1),
                    ',' => (nodes, depth),
                    _ => {
                        nodes.push(Node(c.to_digit(10).unwrap(), depth));
                        (nodes, depth)
                    }
                })
                .0
        })
        .collect()
}

fn magnitude(tree: &Vec<Node>) -> u32 {
    let mut copy = tree.clone();
    loop {
        let deepest = copy.iter().map(|x| x.1).max().unwrap();
        let len = copy.len();
        for i in 0..(len - 1) {
            if copy[i].1 == deepest && copy[i + 1].1 == deepest {
                copy[i].1 -= 1;
                copy[i].0 = 3 * copy[i].0 + 2 * copy[i + 1].0;
                copy.remove(i + 1);
                break;
            }
        }
        if copy.len() == 1 {
            return copy[0].0;
        }
    }
}

fn reduce(tree: &mut Vec<Node>) {
    while explode(tree) || split(tree) {}
}

fn explode(tree: &mut Vec<Node>) -> bool {
    let len = tree.len();

    for i in 0..len {
        if tree[i].1 > 4 {
            let left = tree[i].0;
            let right = tree[i + 1].0;
            tree[i].0 = 0;
            tree[i].1 -= 1;
            if i > 0 {
                tree[i - 1].0 += left;
            }

            if i < len - 2 {
                tree[i + 2].0 += right;
            }
            tree.remove(i + 1);
            return true;
        }
    }
    return false;
}

fn split(tree: &mut Vec<Node>) -> bool {
    let len = tree.len();
    for i in 0..len {
        if tree[i].0 >= 10 {
            let left = tree[i].0 / 2;
            let right = (tree[i].0 + 1) / 2;
            tree[i].0 = left;
            tree[i].1 += 1;
            tree.insert(i + 1, Node(right, tree[i].1));
            return true;
        }
    }
    false
}
