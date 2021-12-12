use std::collections::{HashMap, HashSet};
use std::env;

fn main() {
    let args: Vec<String> = env::args().collect();
    let filename = args[1].clone();
    let graph = read_graph(&filename);
    let mut expl_paths: Vec<Expl> = vec![Expl {
        nodes: HashSet::new(),
        head: Node("start".to_string(), false),
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

        let last = exploration.head.clone();
        if last.is_end() {
            number_of_path_to_the_end += 1;
            continue;
        }
        for next_node in graph.get(&last).unwrap().iter() {
            if next_node.is_big() || !exploration.nodes.contains(&next_node) {
                let mut visited_set = exploration.nodes.clone();
                visited_set.insert(last.clone());
                let next_expl = Expl {
                    nodes: visited_set,
                    head: next_node.clone(),
                    scvt: exploration.scvt,
                };
                expl_paths.push(next_expl);
                continue;
            }

            if !exploration.scvt && !next_node.is_start_or_finish() {
                let mut visited_set = exploration.nodes.clone();
                visited_set.insert(last.clone());
                let next_expl = Expl {
                    nodes: visited_set,
                    head: next_node.clone(),
                    scvt: true,
                };
                expl_paths.push(next_expl);
                continue;
            }
        }
    }

    println!("Number of paths: {}", number_of_path_to_the_end);
}

#[derive(Debug, PartialEq, Eq, Clone, Hash)]
struct Node(String, bool);

impl Node {
    fn is_big(&self) -> bool {
        self.1
    }

    fn is_start(&self) -> bool {
        self.0 == "start".to_string()
    }
    fn is_end(&self) -> bool {
        self.0 == "end".to_string()
    }
    fn is_start_or_finish(&self) -> bool {
        self.is_start() || self.is_end()
    }
}

struct Expl {
    nodes: HashSet<Node>,
    head: Node,
    /// Whether a Small Cavity has been visited twice already on this exploration
    scvt: bool,
}

fn read_graph(filename: &str) -> HashMap<Node, Vec<Node>> {
    std::fs::read_to_string(filename)
        .unwrap()
        .split('\n')
        .into_iter()
        .filter(|s| !s.is_empty())
        .fold(HashMap::new(), |mut map, segment| {
            let seg_str = String::from(segment);
            let splitted: Vec<&str> = seg_str.split("-").collect();
            let node_a = Node(
                String::from(splitted[0]),
                splitted[0].chars().next().unwrap().is_uppercase(),
            );
            let node_b = Node(
                String::from(splitted[1]),
                splitted[1].chars().next().unwrap().is_uppercase(),
            );
            let destinations_a = map.entry(node_a.clone()).or_insert(vec![]);
            (*destinations_a).push(node_b.clone());
            let destinations_b = map.entry(node_b).or_insert(vec![]);
            (*destinations_b).push(node_a);
            map
        })
}
