extern crate nom;

use itertools::Itertools;
use nom::{
    branch::alt, bytes::complete::tag, character::complete, combinator::map, sequence::delimited,
    sequence::separated_pair, IResult,
};
use std::env;
use std::ops::Add;

#[derive(Debug, Clone)]
enum Snail {
    Regular(u64),
    Pair(Box<Snail>, Box<Snail>),
}

impl Add for Snail {
    type Output = Self;

    fn add(self, other: Self) -> Self {
        Self::Pair(Box::new(self), Box::new(other))
    }
}

impl Snail {
    fn reduce(self) -> Self {
        let mut snail = self;
        loop {
            let (next_snail, res) = snail.explode(0);
            snail = next_snail;
            if res.is_some() {
                continue;
            };
            let (next_snail, splat) = snail.split();
            snail = next_snail;
            if !splat {
                break;
            }
        }
        snail
    }

    fn split(self) -> (Self, bool) {
        match self {
            Self::Regular(n) if n >= 10 => (
                Self::Pair(
                    Box::new(Self::Regular(n / 2)),
                    Box::new(Self::Regular((n + 1) / 2)),
                ),
                true,
            ),
            Self::Pair(l, r) => {
                let (l_split, l_was_split) = l.split();
                if l_was_split {
                    (Self::Pair(Box::new(l_split), r), true)
                } else {
                    let (r_split, r_was_split) = r.split();
                    (
                        Self::Pair(Box::new(l_split), Box::new(r_split)),
                        r_was_split,
                    )
                }
            }
            _ => (self, false),
        }
    }

    fn explode(self, depth: usize) -> (Self, Option<(Option<u64>, Option<u64>)>) {
        if let Self::Regular(_) = self {
            return (self, None);
        } else if let Self::Pair(left, right) = self {
            if depth >= 4 {
                match (*left, *right) {
                    (Self::Regular(lv), Self::Regular(rv)) => {
                        return (Self::Regular(0), Some((Some(lv), Some(rv))))
                    }
                    _ => panic!("How did we reached a nested levl of more than 4 😱"),
                }
            }

            let (exploded_left, exploded_pair) = left.explode(depth + 1);
            if let Some(pair) = exploded_pair {
                match pair {
                    (exploded_pair_left, Some(r)) => {
                        return (
                            Self::Pair(
                                Box::new(exploded_left),
                                Box::new(right.add_to_left_end_leaf(r)),
                            ),
                            Some((exploded_pair_left, None)),
                        )
                    }
                    _ => return (Self::Pair(Box::new(exploded_left), right), Some(pair)),
                }
            }

            let original_left = exploded_left; // If we are here it means that actually nothing exploded on the left, setting a new var so the naming reflects that
            let (exploded_right, exploded_pair) = right.explode(depth + 1);
            if let Some(pair) = exploded_pair {
                match pair {
                    (Some(l), exploded_pair_right) => {
                        return (
                            Self::Pair(
                                Box::new(original_left.add_to_right_end_leaf(l)),
                                Box::new(exploded_right),
                            ),
                            Some((None, exploded_pair_right)),
                        )
                    }
                    _ => {
                        return (
                            Self::Pair(Box::new(original_left), Box::new(exploded_right)),
                            Some(pair),
                        )
                    }
                }
            }
            let original_right = exploded_right; // If we are here it means that actually nothing exploded on the right, setting a new var so the naming reflects that
            return (
                Self::Pair(Box::new(original_left), Box::new(original_right)),
                None,
            );
        }
        panic!("Oops! Neither Regular nor pair, something went wrong!");
    }

    fn add_to_left_end_leaf(self, v: u64) -> Self {
        match self {
            Self::Regular(n) => Self::Regular(n + v),
            Self::Pair(l, r) => Self::Pair(Box::new(l.add_to_left_end_leaf(v)), r),
        }
    }

    fn add_to_right_end_leaf(self, v: u64) -> Self {
        match self {
            Snail::Regular(n) => Snail::Regular(n + v),
            Snail::Pair(l, r) => Snail::Pair(l, Box::new(r.add_to_right_end_leaf(v))),
        }
    }

    fn magnitude(&self) -> u64 {
        match self {
            Self::Regular(n) => *n,
            Self::Pair(l, r) => 3 * l.magnitude() + 2 * r.magnitude(),
        }
    }
}

trait ReduceSummable {
    fn reduced_sum(self) -> Snail;
}

impl ReduceSummable for Vec<Snail> {
    fn reduced_sum(self) -> Snail {
        self.into_iter()
            .reduce(|acc, snail| (acc + snail).reduce())
            .unwrap()
    }
}

fn main() {
    let args: Vec<String> = env::args().collect();
    let filename = args[1].clone();
    let snails = read_into_vec(&filename);

    let snail_sum = snails.clone().reduced_sum();

    let best_sum_magnitude = snails
        .into_iter()
        .permutations(2)
        .map(|pairs| pairs.reduced_sum().magnitude())
        .max()
        .unwrap();

    println!("Magnitude of sum {:?}", snail_sum.magnitude());
    println!("Best Magnitude of sum {:?}", best_sum_magnitude);
}

fn read_into_vec(file_name: &str) -> Vec<Snail> {
    std::fs::read_to_string(file_name)
        .unwrap()
        .split('\n')
        .into_iter()
        .filter(|s| !s.is_empty())
        .map(|l| parse_snail(l).unwrap().1)
        .collect()
}

fn parse_snail(snail_str: &str) -> IResult<&str, Snail> {
    alt((
        map(complete::u64, Snail::Regular),
        map(
            delimited(
                tag("["),
                separated_pair(parse_snail, tag(","), parse_snail),
                tag("]"),
            ),
            |(left, right)| Snail::Pair(Box::new(left), Box::new(right)),
        ),
    ))(snail_str)
}
