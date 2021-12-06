pub fn population(t0: Vec<usize>, days: usize) -> Vec<u128> {
    let mut census: Vec<u128> = vec![0; 9];
    for t in t0 {
        census[t] += 1;
    }
    for i in 0..days {
        // census.rotate_left(1); // Useless visual way to debug the timer of each group
        // census[6] += census[8];
        let about_new_parents_index = i % 9;
        let new_parents_resetted_index = (i + 7) % 9;
        census[new_parents_resetted_index] += census[about_new_parents_index];
    }
    census
}

pub fn fishes_at_t0(filename: &str) -> Vec<usize> {
    let fish_line = read_into_vec(&filename);
    fish_line[0]
        .split(',')
        .map(String::from)
        .map(|n| n.parse::<usize>().unwrap())
        .collect()
}

fn read_into_vec(file_name: &str) -> Vec<String> {
    std::fs::read_to_string(file_name)
        .unwrap()
        .split('\n')
        .into_iter()
        .filter(|s| !s.is_empty())
        .map(String::from)
        .collect()
}
