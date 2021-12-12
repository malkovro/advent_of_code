use std::collections::{HashMap, HashSet};

#[derive(Debug, PartialEq, Eq, Hash, Clone, Copy)]
enum Cave {
    Start,
    End,
    Big(u32),
    Small(u32),
}
impl Cave {
    fn is_big(&self) -> bool {
        matches!(*self, Cave::Big(_))
    }
}
pub fn count_paths(filename: &str) -> usize {
    let graph = read_graph(&filename);
    let mut expl_paths: Vec<Expl> = vec![Expl {
        nodes: HashSet::new(),
        head: Cave::Start,
        scvt: false,
    }];
    let mut number_of_path_to_the_end = 0;
    loop {
        let exploration_w = expl_paths.pop();
        if exploration_w.is_none() {
            // Dead end...
            break;
        }
        let exploration = exploration_w.unwrap();

        let last = exploration.head;
        if last == Cave::End {
            number_of_path_to_the_end += 1;
            continue;
        }
        for next_node in graph.get(&last).unwrap().iter() {
            if next_node.is_big() || !exploration.nodes.contains(&next_node) {
                let mut visited_set = exploration.nodes.clone();
                visited_set.insert(last);
                let next_expl = Expl {
                    nodes: visited_set,
                    head: *next_node,
                    scvt: exploration.scvt,
                };
                expl_paths.push(next_expl);
                continue;
            }

            if !exploration.scvt && next_node != &Cave::Start && next_node != &Cave::End {
                let mut visited_set = exploration.nodes.clone();
                visited_set.insert(last);
                let next_expl = Expl {
                    nodes: visited_set,
                    head: *next_node,
                    scvt: true,
                };
                expl_paths.push(next_expl);
                continue;
            }
        }
    }

    number_of_path_to_the_end
}

#[derive(Debug)]
struct Expl {
    nodes: HashSet<Cave>,
    head: Cave,
    /// Whether a Small Cavity has been visited twice already on this exploration
    scvt: bool,
}

fn str_to_cave<'a>(
    str_cave_map: &mut HashMap<&'a str, Cave>,
    max_value: &mut u32,
    value: &'a str,
) -> Cave {
    match value {
        "start" => Cave::Start,
        "end" => Cave::End,
        _ => match str_cave_map.get(value) {
            Some(k) => *k,
            None => {
                *max_value <<= 1;
                let cave = if value.chars().next().unwrap().is_uppercase() {
                    Cave::Big(*max_value)
                } else {
                    Cave::Small(*max_value)
                };
                str_cave_map.insert(value, cave);
                cave
            }
        },
    }
}

fn read_graph(filename: &str) -> HashMap<Cave, Vec<Cave>> {
    let mut max_value: u32 = 1;
    let mut input_map: HashMap<&str, Cave> = HashMap::new();

    std::fs::read_to_string(filename)
        .unwrap()
        .split('\n')
        .into_iter()
        .filter(|s| !s.is_empty())
        .fold(HashMap::new(), |mut map, segment| {
            let splitted: Vec<&str> = segment.split("-").collect::<Vec<&str>>().clone();
            let cave_a = str_to_cave(&mut input_map, &mut max_value, splitted[0]);
            let cave_b = str_to_cave(&mut input_map, &mut max_value, splitted[1]);
            let destinations_b = map.entry(cave_b).or_insert(vec![]);
            (*destinations_b).push(cave_a);
            let destinations_a = map.entry(cave_a).or_insert(vec![]);
            (*destinations_a).push(cave_b);
            map
        })
}
